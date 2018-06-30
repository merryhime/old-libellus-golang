package objstore

import (
	"encoding/hex"
	"fmt"
)

type Oid struct {
	Bytes [20]byte
}

func (o Oid) String() string {
	return fmt.Sprintf("%x", o.Bytes)
}

func (o Oid) Equals(o2 Oid) bool {
	return o.Bytes == o2.Bytes
}

func OidFromString(s string) (Oid, error) {
	if len(s) != 40 {
		return Oid{}, fmt.Errorf("bad oid length %d want 40", len(s))
	}

	var o Oid
	_, err := hex.Decode(o.Bytes[:], []byte(s))
	if err != nil {
		return Oid{}, err
	}
	return o, err
}