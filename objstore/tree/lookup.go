package tree

import (
	"strings"

	"github.com/MerryMage/libellus/objstore"
	"github.com/MerryMage/libellus/objstore/objtype"
)

func Lookup(store objstore.ObjStore, oid objstore.Oid, path string) (*Entry, error) {
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

func LookupInTree(store objstore.ObjStore, t Tree, path string) (*Entry, error) {
	splitpath := strings.SplitN(path, "/", 2)

	e := t.Find(splitpath[0])
	if e == nil {
		return nil, nil
	}

	if len(splitpath) == 1 {
		return e, nil
	}

	return Lookup(store, e.Oid, splitpath[1])
}
