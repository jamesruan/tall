package tall

import (
	"encoding/hex"
)

// hex representation of bytes
type HexBytes string

func (h HexBytes) String() string {
	return string(h)
}

func (h *HexBytes) FromString(s string) error {
	if _, err := hex.DecodeString(s); err != nil {
		return err
	} else {
		*h = HexBytes(s)
		return nil
	}
}

func (h *HexBytes) FromBytes(b []byte) {
	*h = HexBytes(hex.EncodeToString(b))
}

func (h HexBytes) Bytes() []byte {
	if b, err := hex.DecodeString(string(h)); err != nil {
		panic(err)
	} else {
		return b
	}
}
