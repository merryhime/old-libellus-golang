package commit

import (
	"github.com/MerryMage/libellus/objstore"
)

type Commit struct {
	Author    Signature
	Committer Signature
	Message   string
	Tree      objstore.Oid
	Parents   []objstore.Oid
}
