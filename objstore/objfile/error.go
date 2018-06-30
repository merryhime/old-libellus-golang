package objfile

import (
	"fmt"

	"github.com/MerryMage/libellus/objstore"
)

type ObjectNotFoundError struct {
	objstore.ObjectNotFoundError
	Path string
	Oid  objstore.Oid
}

func (e ObjectNotFoundError) Error() string {
	return fmt.Sprintf("objfile: could not find object %s in %#v", e.Oid, e.Path)
}
