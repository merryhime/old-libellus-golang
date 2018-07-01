package tree

import (
	"fmt"

	"github.com/MerryMage/libellus/objstore"
)

type NameAlreadyExistsError string

func (e NameAlreadyExistsError) Error() string {
	return fmt.Sprintf("tree: %#v already exists", e)
}

type NotFoundError string

func (e NotFoundError) Error() string {
	return fmt.Sprintf("tree: could not find %#v", e)
}

type NotATreeError objstore.Oid

func (e NotATreeError) Error() string {
	return fmt.Sprintf("tree: oid %s not a tree", e)
}
