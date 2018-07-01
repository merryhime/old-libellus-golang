package objfile

import (
	"os"
	"path/filepath"

	"github.com/MerryMage/libellus/objstore"
)

type Store struct {
	path string
}

func NewStore(path string) Store {
	return Store{
		path: path,
	}
}

func (store Store) pathToObjectFile(oid objstore.Oid) string {
	oidstr := oid.String()
	return filepath.Join(store.path, "objects", oidstr[:2], oidstr[2:])
}

func (store Store) Get(oid objstore.Oid) (objstore.Obj, error) {
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

func (store Store) Exists(oid objstore.Oid) (bool, error) {
	objpath := store.pathToObjectFile(oid)

	_, err := os.Stat(objpath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}
