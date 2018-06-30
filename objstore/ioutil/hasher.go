package ioutil

import (
	"crypto/sha1"
	"hash"

	"github.com/MerryMage/libellus/objstore"
)

type Hasher struct {
	hash.Hash
}

func NewHasher() Hasher {
	return Hasher{sha1.New()}
}

func (h Hasher) Oid() (oid objstore.Oid) {
	copy(oid.Bytes[:], h.Hash.Sum(nil))
	return oid
}
