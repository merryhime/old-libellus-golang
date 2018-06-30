package tree

import (
	"github.com/MerryMage/libellus/objstore"
	"github.com/MerryMage/libellus/objstore/filemode"
)

type Entry struct {
	Name string
	Mode filemode.FileMode
	Oid  objstore.Oid
}

func (e Entry) sortName() string {
	if e.Mode == filemode.Dir {
		return e.Name + "/"
	}
	return e.Name
}
