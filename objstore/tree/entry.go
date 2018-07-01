package tree

import (
	"github.com/MerryMage/libellus/objstore/filemode"
	"github.com/MerryMage/libellus/objstore/objid"
)

type Entry struct {
	Name string
	Mode filemode.FileMode
	Oid  objid.Oid
}

func (e Entry) sortName() string {
	if e.Mode == filemode.Dir {
		return e.Name + "/"
	}
	return e.Name
}
