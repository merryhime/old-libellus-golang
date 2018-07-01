package objid

import (
	"encoding/hex"
	"fmt"
	"io"
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

func (o Oid) Write(w io.Writer) error {
	_, err := w.Write(o.Bytes[:])
	return err
}

func FromString(s string) (Oid, error) {
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

func Read(r io.Reader) (Oid, error) {
	var o Oid
	_, err := io.ReadFull(r, o.Bytes[:])
	return o, err
}
