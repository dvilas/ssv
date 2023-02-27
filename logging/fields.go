package logging

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/libp2p/go-libp2p/core/peer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	spec "github.com/attestantio/go-eth2-client/spec/phase0"
	specqbft "github.com/bloxapp/ssv-spec/qbft"
	spectypes "github.com/bloxapp/ssv-spec/types"
	"github.com/bloxapp/ssv/network/records"
	"github.com/bloxapp/ssv/protocol/v2/blockchain/beacon"
	protocolp2p "github.com/bloxapp/ssv/protocol/v2/p2p"
)

const (
	FieldAddress             = "address"
	FieldBindIP              = "bindIP"
	FieldBlockCacheMetrics   = "blockCacheMetricsField"
	FieldCurrentSlot         = "currentSlot"
	FieldDurationMilliS      = "durationMilliS"
	FieldEnr                 = "enr"
	FieldEventID             = "eventID"
	FieldFromBlock           = "fromBlock"
	FieldIdentifier          = "identifier"
	FieldIndexCacheMetrics   = "indexCacheMetrics"
	FieldMessageID           = "messageID"
	FieldPeerID              = "peerID"
	FieldPrivateKey          = "privKey"
	FieldPubKey              = "pubKey"
	FieldSyncResults         = "syncResults"
	FieldStartTimeUnixNano   = "startTimeUnixNano"
	FieldSubnets             = "subnets"
	FieldSyncOffset          = "syncOffset"
	FieldTargetNodeEnr       = "targetNodeEnr"
	FieldTxHash              = "txHash"
	FieldUpdatedEnrLocalNode = "updated_enr"
	FieldValidator           = "validator"
)

func FromBlock(val fmt.Stringer) zapcore.Field {
	return zap.Stringer(FieldFromBlock, val)
}

func SyncOffset(val fmt.Stringer) zapcore.Field {
	return zap.Stringer(FieldSyncOffset, val)
}

func TxHash(val fmt.Stringer) zapcore.Field {
	return zap.Stringer(FieldTxHash, val)
}

func EventID(val fmt.Stringer) zapcore.Field {
	return zap.Stringer(FieldEventID, val)
}

func Identifier(val fmt.Stringer) zapcore.Field {
	return zap.Stringer(FieldIdentifier, val)
}

func IdentifierBytes(val []byte) zapcore.Field {
	return zap.Stringer(FieldIdentifier, hexStringer{val})
}

func PubKey(val []byte) zapcore.Field {
	return zap.Stringer(FieldPubKey, hexStringer{val})
}

func PrivKey(val []byte) zapcore.Field {
	return zap.Stringer(FieldPrivateKey, hexStringer{val})
}

func Validator(val []byte) zapcore.Field {
	return zap.Stringer(FieldValidator, hexStringer{val})
}

func Address(val url.URL) zapcore.Field {
	return zap.Stringer(FieldAddress, &val)
}

func Enr(val *enode.Node) zapcore.Field {
	return zap.Stringer(FieldEnr, val)
}

func TargetNodeEnr(val *enode.Node) zapcore.Field {
	return zap.Stringer(FieldTargetNodeEnr, val)
}

func EnrLocalNode(val *enode.LocalNode) zapcore.Field {
	return zap.Stringer(FieldEnr, val.Node())
}

func UpdatedEnrLocalNode(val *enode.LocalNode) zapcore.Field {
	return zap.Stringer(FieldUpdatedEnrLocalNode, val.Node())
}

func Subnets(val records.Subnets) zapcore.Field {
	return zap.Stringer(FieldSubnets, val)
}

func PeerID(val peer.ID) zapcore.Field {
	return zap.Stringer(FieldPeerID, val)
}

func BindIP(val net.IP) zapcore.Field {
	return zap.Stringer(FieldBindIP, val)
}

func MessageID(val spectypes.MessageID) zapcore.Field {
	return zap.Stringer(FieldMessageID, val)
}

func DurationMilliS(val time.Time) zapcore.Field {
	return zap.Stringer(FieldDurationMilliS, int64Stringer{time.Since(val).Milliseconds()})
}

func CurrentSlot(network beacon.Network) zapcore.Field {
	return zap.Stringer(FieldCurrentSlot, uint64Stringer{uint64(network.EstimatedCurrentSlot())})
}

type funcStringer struct {
	fn func() string
}

func (s funcStringer) String() string {
	return s.fn()
}

func StartTimeUnixNano(network beacon.Network, slot spec.Slot) zapcore.Field {
	return zap.Stringer(FieldStartTimeUnixNano, funcStringer{
		fn: func() string {
			return strconv.Itoa(int(network.GetSlotStartTime(slot).UnixNano()))
		},
	})
}

func BlockCacheMetrics(metrics *ristretto.Metrics) zapcore.Field {
	return zap.Stringer(FieldBlockCacheMetrics, metrics)
}

func IndexCacheMetrics(metrics *ristretto.Metrics) zapcore.Field {
	return zap.Stringer(FieldIndexCacheMetrics, metrics)
}

func SyncResults(msgs protocolp2p.SyncResults) zapcore.Field {
	return zap.Stringer(FieldSyncResults, msgs)
}

func OperatorID(operatorId spectypes.OperatorID) zap.Field {
	return zap.Uint64("operator-id", uint64(operatorId))
}

func Height(height specqbft.Height) zap.Field {
	return zap.Uint64("height", uint64(height))
}

func Round(round specqbft.Round) zap.Field {
	return zap.Uint64("round", uint64(round))
}