package objstore

import (
	"errors"
	"strings"

	"github.com/MerryMage/libellus/objstore/commit"
	"github.com/MerryMage/libellus/objstore/filemode"
	"github.com/MerryMage/libellus/objstore/objid"
	"github.com/MerryMage/libellus/objstore/tree"
)

var (
	PathAlreadyExistsError error = errors.New("transaction: path already exists")
	PathDoesNotExistError  error = errors.New("transaction: path does not exist")
)

type transactionTreeEntry struct {
	Mode filemode.FileMode
	Oid  objid.Oid
}

type Transaction struct {
	ref      string
	repo     *Repository
	flatTree map[string]transactionTreeEntry
	parent   objid.Oid
}

func (repo *Repository) StartTransaction(ref string) (*Transaction, error) {
	prevcommit, prevcommitoid, err := repo.Ref(ref)
	if err != nil {
		return nil, err
	}

	trans := &Transaction{
		ref:      ref,
		repo:     repo,
		flatTree: make(map[string]transactionTreeEntry),
		parent:   prevcommitoid,
	}

	err = trans.flattenTree("", prevcommit.Tree)
	return trans, err
}

func catPath(parent string, next string) string {
	if parent != "" {
		return parent + "/" + next
	}
	return next
}

func (trans *Transaction) flattenTree(currentPath string, currentTree objid.Oid) error {
	tree, err := trans.repo.Tree(currentTree)
	if err != nil {
		return err
	}

	for _, e := range tree.Entries {
		path := catPath(currentPath, e.Name)

		if e.Mode == filemode.Dir {
			err = trans.flattenTree(path, e.Oid)
			if err != nil {
				return err
			}
			continue
		}

		trans.flatTree[path] = transactionTreeEntry{
			Mode: e.Mode,
			Oid:  e.Oid,
		}
	}

	return nil
}

func (trans *Transaction) Add(path string, payload []byte) error {
	if _, ok := trans.flatTree[path]; ok {
		return PathAlreadyExistsError
	}

	return trans.AddOrReplace(path, payload)
}

func (trans *Transaction) AddOrReplace(path string, payload []byte) error {
	oid, err := trans.repo.storeBlob(payload)
	if err != nil {
		return err
	}

	trans.flatTree[path] = transactionTreeEntry{
		Mode: filemode.Regular,
		Oid:  oid,
	}

	return nil
}

func (trans *Transaction) Delete(path string) error {
	delete(trans.flatTree, path)
	return nil
}

func (trans *Transaction) Move(src string, dest string) error {
	if _, ok := trans.flatTree[dest]; ok {
		return PathAlreadyExistsError
	}
	if _, ok := trans.flatTree[src]; !ok {
		return PathDoesNotExistError
	}

	tmp := trans.flatTree[src]
	delete(trans.flatTree, src)
	trans.flatTree[dest] = tmp

	return nil
}

func (trans *Transaction) unflattenTree() map[string]*tree.Tree {
	trees := make(map[string]*tree.Tree)

	for path, entry := range trans.flatTree {
		splitpath := strings.Split(path, "/")
		for i := len(splitpath) - 1; i >= 0; i-- {
			parentpath := strings.Join(splitpath[:i], "/")
			name := splitpath[i]

			var mode filemode.FileMode = filemode.Dir
			if i == len(splitpath)-1 {
				mode = entry.Mode
			}

			var oid objid.Oid
			if i == len(splitpath)-1 {
				oid = entry.Oid
			}

			if trees[parentpath] == nil {
				trees[parentpath] = &tree.Tree{}
			}

			trees[parentpath].Add(tree.Entry{
				Name: name,
				Mode: mode,
				Oid:  oid,
			})
		}
	}

	return trees
}

func (trans *Transaction) writeTree(unflattendTree map[string]*tree.Tree, path string) (objid.Oid, error) {
	currentTree := unflattendTree[path]
	for i := range currentTree.Entries {
		if currentTree.Entries[i].Mode == filemode.Dir {
			childoid, err := trans.writeTree(unflattendTree, catPath(path, currentTree.Entries[i].Name))
			if err != nil {
				return objid.Oid{}, err
			}
			currentTree.Entries[i].Oid = childoid
		}
	}
	return trans.repo.storeTree(*currentTree)
}

func (trans *Transaction) Store(c commit.Commit) error {
	unflattenedTree := trans.unflattenTree()
	tree, err := trans.writeTree(unflattenedTree, "")
	if err != nil {
		return err
	}

	c.Tree = tree
	c.Parents = []objid.Oid{trans.parent}

	coid, err := trans.repo.storeCommit(c)
	if err != nil {
		return err
	}

	return trans.repo.writeRef(trans.ref, coid)
}
