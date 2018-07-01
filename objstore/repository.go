package objstore

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/MerryMage/libellus/objstore/commit"
	"github.com/MerryMage/libellus/objstore/obj"
	"github.com/MerryMage/libellus/objstore/objfile"
	"github.com/MerryMage/libellus/objstore/objid"
	"github.com/MerryMage/libellus/objstore/objtype"
	"github.com/MerryMage/libellus/objstore/tree"
)

var (
	NotACommitError error = errors.New("repository: requested object was not a commit")
	NotATreeError   error = errors.New("repository: requested object was not a tree")
	NotABlobError   error = errors.New("repository: requested object was not a blob")
)

type Repository struct {
	lock     sync.RWMutex
	path     string
	objStore objfile.Store
}

func NewRepository(path string) Repository {
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		path = filepath.Join(path, ".git")
	}

	return Repository{
		path:     path,
		objStore: objfile.NewStore(path),
	}
}

func (repo Repository) Get(oid objid.Oid) (obj.Obj, error) {
	return repo.objStore.Get(oid)
}

func (repo Repository) Exists(oid objid.Oid) (bool, error) {
	return repo.objStore.Exists(oid)
}

func (repo Repository) Store(ot objtype.ObjType, payload []byte) (objid.Oid, error) {
	return repo.objStore.Store(ot, payload)
}

func (repo Repository) RefOid(ref string) (objid.Oid, error) {
	refpath := filepath.Join(repo.path, "refs", "heads", ref)

	rawoid, err := ioutil.ReadFile(refpath)
	if err != nil {
		return objid.Oid{}, err
	}

	rawoid = bytes.Trim(rawoid, " \n\r")
	return objid.FromString(string(rawoid))
}

func (repo Repository) Ref(ref string) (commit.Commit, error) {
	oid, err := repo.RefOid(ref)
	if err != nil {
		return commit.Commit{}, err
	}

	o, err := repo.Get(oid)
	if err != nil {
		return commit.Commit{}, err
	} else if o.ObjType() != objtype.Commit {
		return commit.Commit{}, NotACommitError
	}

	return commit.Read(o)
}

func (repo Repository) Tree(oid objid.Oid) (tree.Tree, error) {
	o, err := repo.Get(oid)
	if err != nil {
		return tree.Tree{}, err
	} else if o.ObjType() != objtype.Tree {
		return tree.Tree{}, NotATreeError
	}

	return tree.Read(o)
}

func (repo Repository) LookupEntryByPath(ref string, path string) (*tree.Entry, error) {
	commit, err := repo.Ref(ref)
	if err != nil {
		return nil, err
	}

	return tree.Lookup(repo, commit.Tree, path)
}

func (repo Repository) Blob(oid objid.Oid) (io.ReadCloser, error) {
	o, err := repo.Get(oid)
	if err != nil {
		return nil, err
	} else if o.ObjType() != objtype.Blob {
		return nil, NotABlobError
	}
	return o, nil
}

func (repo Repository) LookupBlobByPath(ref string, path string) (io.ReadCloser, error) {
	e, err := repo.LookupEntryByPath(ref, path)
	if err != nil {
		return nil, err
	}
	return repo.Blob(e.Oid)
}
