package ioutil

import (
	"crypto/sha1"
	"hash"

	"github.com/MerryMage/libellus/objstore/objid"
)

type Hasher struct {
	hash.Hash
}

func NewHasher() Hasher {
	return Hasher{sha1.New()}
}

func (h Hasher) Oid() (oid objid.Oid) {
	copy(oid.Bytes[:], h.Hash.Sum(nil))
	return oid
}
