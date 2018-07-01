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

type ObjGetter interface {
	Get(oid Oid) (Obj, error)
	Exists(oid Oid) (bool, error)
}

type ObjStorer interface {
	Store(ot objtype.ObjType, payload []byte) (Oid, error)
}

type ObjGetStorer interface {
	ObjGetter
	ObjStorer
}
