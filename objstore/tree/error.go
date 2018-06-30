package tree

import (
	"fmt"
)

type NameAlreadyExistsError string

func (e NameAlreadyExistsError) Error() string {
	return fmt.Sprintf("tree: %#v already exists", e)
}

type NotFoundError string

func (e NotFoundError) Error() string {
	return fmt.Sprintf("tree: could not find %#v", e)
}
