package tree

import (
	"strings"

	"github.com/MerryMage/libellus/objstore/obj"
	"github.com/MerryMage/libellus/objstore/objid"
	"github.com/MerryMage/libellus/objstore/objtype"
)

func Lookup(store obj.ObjGetter, oid objid.Oid, path string) (*Entry, error) {
	obj, err := store.Get(oid)
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	if obj.ObjType() != objtype.Tree {
		return nil, NotATreeError(oid)
	}

	t, err := Read(obj)
	if err != nil {
		return nil, err
	}
	return LookupInTree(store, t, path)
}

func LookupInTree(store obj.ObjGetter, t Tree, path string) (*Entry, error) {
	splitpath := strings.SplitN(path, "/", 2)

	e := t.Find(splitpath[0])
	if e == nil {
		return nil, NotFoundError(path)
	}

	if len(splitpath) == 1 {
		return e, nil
	}

	return Lookup(store, e.Oid, splitpath[1])
}
