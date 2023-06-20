package ekm

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
	"time"

	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	specqbft "github.com/bloxapp/ssv-spec/qbft"
	spectypes "github.com/bloxapp/ssv-spec/types"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/jamiealquiza/tachymeter"
	"github.com/prysmaticlabs/go-bitfield"
	"github.com/sourcegraph/conc/pool"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/ssv/logging"
	"github.com/bloxapp/ssv/networkconfig"
	"github.com/bloxapp/ssv/protocol/v2/types"
	"github.com/bloxapp/ssv/utils/threshold"
)

const (
	sk1Str = "3548db63ab5701878daf25fa877638dc7809778815b9d9ecd5369da33ca9e64f"
	pk1Str = "a8cb269bd7741740cfe90de2f8db6ea35a9da443385155da0fa2f621ba80e5ac14b5c8f65d23fd9ccc170cc85f29e27d"
	sk2Str = "66dd37ae71b35c81022cdde98370e881cff896b689fa9136917f45afce43fd3b"
	pk2Str = "8796fafa576051372030a75c41caafea149e4368aebaca21c9f90d9974b3973d5cee7d7874e4ec9ec59fb2c8945b3e01"
)

func testKeyManager(t *testing.T) spectypes.KeyManager {
	threshold.Init()

	logger := logging.TestLogger(t)

	db, err := getBaseStorage(logger)
	require.NoError(t, err)

	km, err := NewETHKeyManagerSigner(logger, db, networkconfig.TestNetwork, true)
	require.NoError(t, err)

	sk1 := &bls.SecretKey{}
	require.NoError(t, sk1.SetHexString(sk1Str))

	sk2 := &bls.SecretKey{}
	require.NoError(t, sk2.SetHexString(sk2Str))

	require.NoError(t, km.AddShare(sk1))
	require.NoError(t, km.AddShare(sk2))

	return km
}

func TestSlashing(t *testing.T) {
	km := testKeyManager(t)

	sk1 := &bls.SecretKey{}
	require.NoError(t, sk1.SetHexString(sk1Str))
	require.NoError(t, km.AddShare(sk1))

	currentSlot := km.(*ethKeyManagerSigner).storage.Network().EstimatedCurrentSlot()
	currentEpoch := km.(*ethKeyManagerSigner).storage.Network().EstimatedEpochAtSlot(currentSlot)

	highestTarget := currentEpoch + minimalAttSlashingProtectionEpochDistance + 1
	highestSource := highestTarget - 1
	highestProposal := currentSlot + minimalBlockSlashingProtectionSlotDistance + 1

	attestationData := &phase0.AttestationData{
		Slot:            30,
		Index:           1,
		BeaconBlockRoot: [32]byte{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2},
		Source: &phase0.Checkpoint{
			Epoch: highestSource,
			Root:  [32]byte{},
		},
		Target: &phase0.Checkpoint{
			Epoch: highestTarget,
			Root:  [32]byte{},
		},
	}

	var beaconBlock = &bellatrix.BeaconBlock{
		Slot:          highestProposal,
		ProposerIndex: 0,
		ParentRoot: phase0.Root{
			0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
			0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		},
		StateRoot: phase0.Root{
			0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
			0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
		},
		Body: &bellatrix.BeaconBlockBody{
			RANDAOReveal: phase0.BLSSignature{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
			ETH1Data: &phase0.ETH1Data{
				DepositRoot: phase0.Root{
					0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4a, 0x4b, 0x4c, 0x4d, 0x4e, 0x4f,
					0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58, 0x59, 0x5a, 0x5b, 0x5c, 0x5d, 0x5e, 0x5f,
				},
				DepositCount: 0,
				BlockHash: []byte{
					0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c, 0x6d, 0x6e, 0x6f,
					0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7a, 0x7b, 0x7c, 0x7d, 0x7e, 0x7f,
				},
			},
			Graffiti: [32]byte{
				0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
				0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			},
			ProposerSlashings: []*phase0.ProposerSlashing{},
			AttesterSlashings: []*phase0.AttesterSlashing{},
			Attestations:      []*phase0.Attestation{},
			Deposits:          []*phase0.Deposit{},
			VoluntaryExits:    []*phase0.SignedVoluntaryExit{},
			SyncAggregate: &altair.SyncAggregate{
				SyncCommitteeBits: bitfield.Bitvector512{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				},
				SyncCommitteeSignature: phase0.BLSSignature{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				},
			},
			ExecutionPayload: &bellatrix.ExecutionPayload{
				ParentHash: phase0.Hash32{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				},
				FeeRecipient: bellatrix.ExecutionAddress{},
				StateRoot: [32]byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				},
				ReceiptsRoot: [32]byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				},
				LogsBloom: [256]byte{},
				PrevRandao: [32]byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				},
				BlockNumber: 0,
				GasLimit:    0,
				GasUsed:     0,
				Timestamp:   0,
				ExtraData:   nil,
				BaseFeePerGas: [32]byte{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				},
				BlockHash: phase0.Hash32{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
					0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
				},
				Transactions: []bellatrix.Transaction{},
			},
		},
	}

	t.Run("sign once", func(t *testing.T) {
		_, sig, err := km.(*ethKeyManagerSigner).SignBeaconObject(attestationData, phase0.Domain{}, sk1.GetPublicKey().Serialize(), spectypes.DomainAttester)
		require.NoError(t, err)
		require.NotNil(t, sig)
		require.NotEqual(t, [32]byte{}, sig)
	})
	t.Run("slashable sign, fail", func(t *testing.T) {
		_, sig, err := km.(*ethKeyManagerSigner).SignBeaconObject(attestationData, phase0.Domain{}, sk1.GetPublicKey().Serialize(), spectypes.DomainAttester)
		require.EqualError(t, err, "slashable attestation (HighestAttestationVote), not signing")
		require.Equal(t, [32]byte{}, sig)
	})

	t.Run("sign once", func(t *testing.T) {
		_, sig, err := km.(*ethKeyManagerSigner).SignBeaconObject(beaconBlock, phase0.Domain{}, sk1.GetPublicKey().Serialize(), spectypes.DomainProposer)
		require.NoError(t, err)
		require.NotNil(t, sig)
		require.NotEqual(t, [32]byte{}, sig)
	})
	t.Run("slashable sign, fail", func(t *testing.T) {
		_, sig, err := km.(*ethKeyManagerSigner).SignBeaconObject(beaconBlock, phase0.Domain{}, sk1.GetPublicKey().Serialize(), spectypes.DomainProposer)
		require.EqualError(t, err, "slashable proposal (HighestProposalVote), not signing")
		require.Equal(t, [32]byte{}, sig)
	})
}

func TestSignRoot(t *testing.T) {
	require.NoError(t, bls.Init(bls.BLS12_381))

	km := testKeyManager(t)

	t.Run("pk 1", func(t *testing.T) {
		pk := &bls.PublicKey{}
		require.NoError(t, pk.Deserialize(_byteArray(pk1Str)))

		msg := specqbft.Message{
			MsgType:    specqbft.CommitMsgType,
			Height:     specqbft.Height(3),
			Round:      specqbft.Round(2),
			Identifier: []byte("identifier1"),
			Root:       [32]byte{1, 2, 3},
		}

		// sign
		sig, err := km.SignRoot(&msg, spectypes.QBFTSignatureType, pk.Serialize())
		require.NoError(t, err)

		// verify
		signed := &specqbft.SignedMessage{
			Signature: sig,
			Signers:   []spectypes.OperatorID{1},
			Message:   msg,
		}

		err = types.VerifyByOperators(signed.GetSignature(), signed, networkconfig.TestNetwork.Domain, spectypes.QBFTSignatureType, []*spectypes.Operator{{OperatorID: spectypes.OperatorID(1), PubKey: pk.Serialize()}})
		// res, err := signed.VerifySig(pk)
		require.NoError(t, err)
		// require.True(t, res)
	})

	t.Run("pk 2", func(t *testing.T) {
		pk := &bls.PublicKey{}
		require.NoError(t, pk.Deserialize(_byteArray(pk2Str)))

		msg := specqbft.Message{
			MsgType:    specqbft.CommitMsgType,
			Height:     specqbft.Height(1),
			Round:      specqbft.Round(3),
			Identifier: []byte("identifier2"),
			Root:       [32]byte{4, 5, 6},
		}

		// sign
		sig, err := km.SignRoot(&msg, spectypes.QBFTSignatureType, pk.Serialize())
		require.NoError(t, err)

		// verify
		signed := &specqbft.SignedMessage{
			Signature: sig,
			Signers:   []spectypes.OperatorID{1},
			Message:   msg,
		}

		start := time.Now()
		for i := 0; i < 1000; i++ {
			err = types.VerifyByOperators(signed.GetSignature(), signed, networkconfig.TestNetwork.Domain, spectypes.QBFTSignatureType, []*spectypes.Operator{{OperatorID: spectypes.OperatorID(1), PubKey: pk.Serialize()}})
			// res, err := signed.VerifySig(pk)
			require.NoError(t, err)
			// require.True(t, res)
		}
		fmt.Println(time.Since(start))
	})

	t.Run("multiple signers", func(t *testing.T) {
		t.Skip()

		pk1 := &bls.PublicKey{}
		require.NoError(t, pk1.Deserialize(_byteArray(pk1Str)))
		pk2 := &bls.PublicKey{}
		require.NoError(t, pk2.Deserialize(_byteArray(pk2Str)))

		go func() {
			msg := specqbft.Message{
				MsgType:    specqbft.CommitMsgType,
				Height:     specqbft.Height(1),
				Round:      specqbft.Round(3),
				Identifier: []byte("identifier2"),
				Root:       [32]byte{4, 5, 6},
			}

			// sign
			sig1, err := km.SignRoot(&msg, spectypes.QBFTSignatureType, pk1.Serialize())
			require.NoError(t, err)
			sign1 := &bls.Sign{}
			err = sign1.Deserialize(sig1)
			require.NoError(t, err)

			sig2, err := km.SignRoot(&msg, spectypes.QBFTSignatureType, pk2.Serialize())
			require.NoError(t, err)
			sign2 := &bls.Sign{}
			err = sign2.Deserialize(sig2)
			require.NoError(t, err)

			sign := sign1
			sign.Add(sign2)

			// verify
			signed := &specqbft.SignedMessage{
				Signature: sign.Serialize(),
				Signers:   []spectypes.OperatorID{1, 2},
				Message:   msg,
			}

			err = types.VerifyByOperators(signed.GetSignature(), signed, networkconfig.TestNetwork.Domain, spectypes.QBFTSignatureType,
				[]*spectypes.Operator{
					{OperatorID: spectypes.OperatorID(1), PubKey: pk1.Serialize()},
					{OperatorID: spectypes.OperatorID(2), PubKey: pk2.Serialize()},
				})
			// res, err := signed.VerifySig(pk)
			require.NoError(t, err)
			// require.True(t, res)
		}()

		msg := specqbft.Message{
			MsgType:    specqbft.CommitMsgType,
			Height:     specqbft.Height(1),
			Round:      specqbft.Round(3),
			Identifier: []byte("identifier2"),
			Root:       [32]byte{4, 5, 6},
		}

		// sign
		sig1, err := km.SignRoot(&msg, spectypes.QBFTSignatureType, pk1.Serialize())
		require.NoError(t, err)
		sign1 := &bls.Sign{}
		err = sign1.Deserialize(sig1)
		require.NoError(t, err)

		sig2, err := km.SignRoot(&msg, spectypes.QBFTSignatureType, pk2.Serialize())
		require.NoError(t, err)
		sign2 := &bls.Sign{}
		err = sign2.Deserialize(sig2)
		require.NoError(t, err)

		sign := *sign1
		sign.Add(sign2)

		// verify
		signed := &specqbft.SignedMessage{
			Signature: sign.Serialize(),
			Signers:   []spectypes.OperatorID{1, 2},
			Message:   msg,
		}

		start := time.Now()
		for i := 0; i < 1000; i++ {
			err = types.VerifyByOperators(signed.GetSignature(), signed, networkconfig.TestNetwork.Domain, spectypes.QBFTSignatureType,
				[]*spectypes.Operator{
					{OperatorID: spectypes.OperatorID(1), PubKey: pk1.Serialize()},
					{OperatorID: spectypes.OperatorID(2), PubKey: pk2.Serialize()},
				})
			require.NoError(t, err)
		}
		fmt.Println(time.Since(start))
	})

	t.Run("multiple signers 2", func(t *testing.T) {
		// t.Skip()

		pk1 := &bls.PublicKey{}
		require.NoError(t, pk1.Deserialize(_byteArray(pk1Str)))
		pk2 := &bls.PublicKey{}
		require.NoError(t, pk2.Deserialize(_byteArray(pk2Str)))

		var N = 25_000
		var msgs = make([]*specqbft.SignedMessage, N)

		p := pool.New()
		for i := 0; i < N; i++ {
			i := i
			p.Go(func() {
				msg := specqbft.Message{
					MsgType:    specqbft.CommitMsgType,
					Height:     specqbft.Height(rand.Uint64()),
					Round:      specqbft.Round(3),
					Identifier: []byte("identifier2"),
					Root:       [32]byte{4, 5, 6},
				}

				// sign
				sig1, err := km.SignRoot(&msg, spectypes.QBFTSignatureType, pk1.Serialize())
				require.NoError(t, err)
				sign1 := &bls.Sign{}
				err = sign1.Deserialize(sig1)
				require.NoError(t, err)

				sig2, err := km.SignRoot(&msg, spectypes.QBFTSignatureType, pk2.Serialize())
				require.NoError(t, err)
				sign2 := &bls.Sign{}
				err = sign2.Deserialize(sig2)
				require.NoError(t, err)

				sign := *sign1
				sign.Add(sign2)

				// verify
				msgs[i] = &specqbft.SignedMessage{
					Signature: sign.Serialize(),
					Signers:   []spectypes.OperatorID{1, 2},
					Message:   msg,
				}
			})
		}
		p.Wait()

		// var total time.Duration
		// for j := 0; j < 10; j++ {
		// 	start := time.Now()
		// 	p = pool.New()
		// 	for i := range msgs {
		// 		signed := msgs[i]
		// 		p.Go(func() {
		// 			err := types.VerifyByOperators(signed.GetSignature(), signed, networkconfig.TestNetwork.Domain, spectypes.QBFTSignatureType,
		// 				[]*spectypes.Operator{
		// 					{OperatorID: spectypes.OperatorID(1), PubKey: pk1.Serialize()},
		// 					{OperatorID: spectypes.OperatorID(2), PubKey: pk2.Serialize()},
		// 				})
		// 			require.NoError(t, err)
		// 		})
		// 	}
		// 	p.Wait()
		// 	fmt.Println(time.Since(start))
		// 	total += time.Since(start)
		// }
		// fmt.Println("avg", total/10)

		f, err := os.Create("/Users/mosherevah/Downloads/cpu.pprof")
		require.NoError(t, err)
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()

		var wg sync.WaitGroup
		duration := time.Millisecond * 12000
		var total time.Duration
		start := time.Now()
		tm := tachymeter.New(&tachymeter.Config{Size: N})
		for i := 0; i < N; i++ {
			i := i
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Sleep random value between 0 and 12 seconds
				start2 := time.Now()
				sleep := time.Duration(rand.Intn(int(duration)))
				time.Sleep(sleep)

				signed := msgs[i]
				err := types.VerifyByOperators(signed.GetSignature(), signed, networkconfig.TestNetwork.Domain, spectypes.QBFTSignatureType,
					[]*spectypes.Operator{
						{OperatorID: spectypes.OperatorID(1), PubKey: pk1.Serialize()},
						{OperatorID: spectypes.OperatorID(2), PubKey: pk2.Serialize()},
					})
				require.NoError(t, err)
				dur := time.Since(start2) - sleep
				total += dur
				tm.AddTime(dur)
			}()
		}
		wg.Wait()
		fmt.Println("Total Run Time:", time.Since(start))
		fmt.Println("Average Latency:", total/time.Duration(N))
		b, _ := json.MarshalIndent(types.Verifier.Stats(), "", "  ")
		fmt.Println("Stats:", string(b))
		fmt.Println(tm.Calc())
	})
}
