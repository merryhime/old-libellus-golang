package objfile

import (
	"fmt"

	"github.com/MerryMage/libellus/objstore/objid"
)

type ObjectNotFoundError struct {
	Path string
	Oid  objid.Oid
}

func (ObjectNotFoundError) objectNotFoundError() {}

func (e ObjectNotFoundError) Error() string {
	return fmt.Sprintf("objfile: could not find object %s in %#v", e.Oid, e.Path)
}
