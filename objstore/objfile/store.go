package objfile

import (
	"os"
	"path/filepath"

	"github.com/MerryMage/libellus/objstore/obj"
	"github.com/MerryMage/libellus/objstore/objid"
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

func (store Store) Store(oid objid.Oid, payload []byte) error {
	objpath := store.pathToObjectFile(oid)

	err := os.MkdirAll(store.dirContainingObjectFile(oid), 0777)
	if err != nil {
		return err
	}

	f, err := os.Create(objpath)
	if err != nil {
		return err
	}

	_, err = f.Write(payload)
	if err != nil {
		f.Close()
		os.Remove(objpath)
		return err
	}

	err = f.Sync()
	if err != nil {
		f.Close()
		os.Remove(objpath)
		return err
	}

	err = f.Close()
	if err != nil {
		os.Remove(objpath)
		return err
	}

	return nil
}
