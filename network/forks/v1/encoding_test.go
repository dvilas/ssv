package v1

import (
	"bytes"
	"github.com/bloxapp/ssv/protocol/v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestForkV1_Encoding(t *testing.T) {
	msg := &v1.SSVMessage{
		MsgType: v1.SSVConsensusMsgType,
		ID:      []byte("xxxxxxxxxxx_ATTESTER"),
		Data:    []byte("data"),
	}
	f := &ForkV1{}

	b, err := f.EncodeNetworkMsg(msg)
	require.NoError(t, err)
	require.Greater(t, len(b), 0)

	res, err := f.DecodeNetworkMsg(b)
	require.NoError(t, err)
	require.Equal(t, msg.MsgType, res.(*v1.SSVMessage).MsgType)
	require.True(t, bytes.Equal(msg.Data, res.(*v1.SSVMessage).Data))
}
