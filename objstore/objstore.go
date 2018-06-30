package objstore

import (
	"io"

	"github.com/MerryMage/libellus/objstore/objtype"
)

type ObjectNotFoundError struct{}

type Obj interface {
	io.ReadCloser

	Size() uint64
	ObjType() objtype.ObjType
}

type ObjStore interface {
	Get(oid Oid) (Obj, error)
	Exists(oid Oid) (bool, error)
}
