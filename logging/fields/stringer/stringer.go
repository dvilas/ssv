package stringer

import (
	"encoding/hex"
	"strconv"
)

type HexStringer struct {
	Val []byte
}

func (h HexStringer) String() string {
	return hex.EncodeToString(h.Val)
}

type Int64Stringer struct {
	Val int64
}

func (h Int64Stringer) String() string {
	return strconv.Itoa(int(h.Val))
}

type Uint64Stringer struct {
	Val uint64
}

func (h Uint64Stringer) String() string {
	return strconv.FormatUint(h.Val, 10)
}

type FuncStringer struct {
	Fn func() string
}

func (s FuncStringer) String() string {
	return s.Fn()
}