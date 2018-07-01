package obj

import (
	"io"

	"github.com/MerryMage/libellus/objstore/objid"
	"github.com/MerryMage/libellus/objstore/objtype"
)

type ObjectNotFoundError interface {
	objectNotFoundError()
}

type Obj interface {
	io.ReadCloser

	Size() uint64
	ObjType() objtype.ObjType
}

type ObjGetter interface {
	Get(oid objid.Oid) (Obj, error)
	Exists(oid objid.Oid) (bool, error)
}

type ObjStorer interface {
	Store(oid objid.Oid, payload []byte) error
}

type ObjGetStorer interface {
	ObjGetter
	ObjStorer
}
