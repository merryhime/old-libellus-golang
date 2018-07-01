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

func NewRepository(path string) *Repository {
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		path = filepath.Join(path, ".git")
	}

	return &Repository{
		path:     path,
		objStore: objfile.NewStore(path),
	}
}

func (repo *Repository) Get(oid objid.Oid) (obj.Obj, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	return repo.objStore.Get(oid)
}

func (repo *Repository) Exists(oid objid.Oid) (bool, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()
	return repo.exists(oid)
}

func (repo *Repository) exists(oid objid.Oid) (bool, error) {
	return repo.objStore.Exists(oid)
}

func (repo *Repository) Store(ot objtype.ObjType, payload []byte) (objid.Oid, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()
	return repo.store(ot, payload)
}

func (repo *Repository) store(ot objtype.ObjType, payload []byte) (objid.Oid, error) {
	var b bytes.Buffer

	w, err := objfile.NewWriter(&b, ot, uint64(len(payload)))
	if err != nil {
		return objid.Oid{}, err
	}

	_, err = w.Write(payload)
	if err != nil {
		return objid.Oid{}, err
	}
	w.Close()

	oid := w.Oid()
	return oid, repo.objStore.Store(oid, b.Bytes())
}

func (repo *Repository) storeBlob(b []byte) (objid.Oid, error) {
	return repo.store(objtype.Blob, b)
}

func (repo *Repository) storeTree(t tree.Tree) (objid.Oid, error) {
	var b bytes.Buffer
	err := t.Write(&b)
	if err != nil {
		return objid.Oid{}, err
	}
	return repo.store(objtype.Tree, b.Bytes())
}

func (repo *Repository) storeCommit(c commit.Commit) (objid.Oid, error) {
	var b bytes.Buffer
	err := c.Write(&b)
	if err != nil {
		return objid.Oid{}, err
	}
	return repo.store(objtype.Commit, b.Bytes())
}

func (repo *Repository) RefOid(ref string) (objid.Oid, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	refpath := filepath.Join(repo.path, "refs", "heads", ref)

	rawoid, err := ioutil.ReadFile(refpath)
	if err != nil {
		return objid.Oid{}, err
	}

	rawoid = bytes.Trim(rawoid, " \n\r")
	return objid.FromString(string(rawoid))
}

func (repo *Repository) Ref(ref string) (commit.Commit, objid.Oid, error) {
	oid, err := repo.RefOid(ref)
	if err != nil {
		return commit.Commit{}, objid.Oid{}, err
	}

	o, err := repo.Get(oid)
	if err != nil {
		return commit.Commit{}, objid.Oid{}, err
	} else if o.ObjType() != objtype.Commit {
		return commit.Commit{}, objid.Oid{}, NotACommitError
	}

	c, err := commit.Read(o)
	return c, oid, err
}

func (repo *Repository) writeRef(ref string, oid objid.Oid) error {
	refpath := filepath.Join(repo.path, "refs", "heads", ref)
	f, err := os.OpenFile(refpath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	_, err = f.WriteString(oid.String() + "\n")
	if err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func (repo *Repository) Tree(oid objid.Oid) (tree.Tree, error) {
	o, err := repo.Get(oid)
	if err != nil {
		return tree.Tree{}, err
	} else if o.ObjType() != objtype.Tree {
		return tree.Tree{}, NotATreeError
	}

	return tree.Read(o)
}

func (repo *Repository) LookupEntryByPath(ref string, path string) (*tree.Entry, error) {
	commit, _, err := repo.Ref(ref)
	if err != nil {
		return nil, err
	}

	return tree.Lookup(repo, commit.Tree, path)
}

func (repo *Repository) Blob(oid objid.Oid) (io.ReadCloser, error) {
	o, err := repo.Get(oid)
	if err != nil {
		return nil, err
	} else if o.ObjType() != objtype.Blob {
		return nil, NotABlobError
	}
	return o, nil
}

func (repo *Repository) LookupBlobByPath(ref string, path string) (io.ReadCloser, error) {
	e, err := repo.LookupEntryByPath(ref, path)
	if err != nil {
		return nil, err
	}
	return repo.Blob(e.Oid)
}

func (repo *Repository) LookupTreeByPath(ref string, path string) (tree.Tree, error) {
	e, err := repo.LookupEntryByPath(ref, path)
	if err != nil {
		return tree.Tree{}, err
	}
	return repo.Tree(e.Oid)
}
