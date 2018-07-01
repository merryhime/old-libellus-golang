package objfile

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/MerryMage/libellus/objstore/obj"
	"github.com/MerryMage/libellus/objstore/objid"
	"github.com/MerryMage/libellus/objstore/objtype"
)

type Store struct {
	path string
}

func NewStore(path string) Store {
	return Store{
		path: path,
	}
}

func (store Store) pathToObjectFile(oid objid.Oid) string {
	oidstr := oid.String()
	return filepath.Join(store.path, "objects", oidstr[:2], oidstr[2:])
}

func (store Store) dirContainingObjectFile(oid objid.Oid) string {
	oidstr := oid.String()
	return filepath.Join(store.path, "objects", oidstr[:2])
}

func (store Store) Get(oid objid.Oid) (obj.Obj, error) {
	objpath := store.pathToObjectFile(oid)

	_, err := os.Stat(objpath)
	if os.IsNotExist(err) {
		return nil, ObjectNotFoundError{Path: store.path, Oid: oid}
	} else if err != nil {
		return nil, err
	}

	f, err := os.Open(objpath)
	if err != nil {
		return nil, err
	}

	return NewReader(f, true)
}

func (store Store) Exists(oid objid.Oid) (bool, error) {
	objpath := store.pathToObjectFile(oid)

	_, err := os.Stat(objpath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (store Store) Store(ot objtype.ObjType, payload []byte) (objid.Oid, error) {
	var b bytes.Buffer

	w, err := NewWriter(&b, ot, uint64(len(payload)))
	if err != nil {
		return objid.Oid{}, err
	}

	_, err = w.Write(payload)
	if err != nil {
		return objid.Oid{}, err
	}
	w.Close()

	oid := w.Oid()
	objpath := store.pathToObjectFile(oid)

	err = os.MkdirAll(store.dirContainingObjectFile(oid), 0777)
	if err != nil {
		return oid, err
	}

	f, err := os.Create(objpath)
	if err != nil {
		return oid, err
	}

	_, err = f.Write(b.Bytes())
	if err != nil {
		f.Close()
		os.Remove(objpath)
		return oid, err
	}

	return oid, f.Close()
}
