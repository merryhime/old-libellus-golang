package commit

import (
	"github.com/MerryMage/libellus/objstore/objid"
)

type Commit struct {
	Author    Signature
	Committer Signature
	Message   string
	Tree      objid.Oid
	Parents   []objid.Oid
}
