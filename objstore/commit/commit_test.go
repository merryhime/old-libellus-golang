package commit

import (
	"bytes"
	"testing"

	"encoding/hex"

	"github.com/MerryMage/libellus/objstore"
)

var testCommit []byte = func() []byte {
	ret, _ := hex.DecodeString("7472656520613839636662653134393635356265343835313264386537343861306230623737396538356536390a706172656e7420383666376334636233333038323839633739626135303037366139313630643939373535616436360a617574686f72204d657272794d616765203c4d657272794d6167654075736572732e6e6f7265706c792e6769746875622e636f6d3e2031353330333537363637202b303130300a636f6d6d6974746572204d657272794d616765203c4d657272794d6167654075736572732e6e6f7265706c792e6769746875622e636f6d3e2031353330333631323334202b303130300a0a4136343a20496d706c656d656e7420464356544d5520287363616c6172290a")
	return ret
}()

func oid(s string) objstore.Oid {
	oid, err := objstore.OidFromString(s)
	if err != nil {
		panic("oid failed")
	}
	return oid
}

func TestRead(t *testing.T) {
	b := bytes.NewBuffer(testCommit)

	commit, err := Read(b)
	if err != nil {
		t.Error(err)
	}

	if len(commit.Parents) != 1 || commit.Parents[0] != oid("86f7c4cb3308289c79ba50076a9160d99755ad66") {
		t.Errorf("commit.Parents = %#v", commit.Parents)
	}
	if commit.Tree != oid("a89cfbe149655be48512d8e748a0b0b779e85e69") {
		t.Errorf("commit.Tree = %#v", commit.Tree)
	}
	if commit.Author.Name != "MerryMage" || commit.Author.Email != "MerryMage@users.noreply.github.com" || commit.Author.Timestamp != 1530357667 || commit.Author.Timezone != "+0100" {
		t.Errorf("commit.Author = %#v", commit.Author)
	}
	if commit.Committer.Name != "MerryMage" || commit.Committer.Email != "MerryMage@users.noreply.github.com" || commit.Committer.Timestamp != 1530361234 || commit.Committer.Timezone != "+0100" {
		t.Errorf("commit.Committer = %#v", commit.Committer)
	}
	if commit.Message != "A64: Implement FCVTMU (scalar)\n" {
		t.Errorf("commit.Message = %#v", commit.Message)
	}
}

func TestWrite(t *testing.T) {
	var b bytes.Buffer

	commit := Commit{
		Author: Signature{
			Name:      "MerryMage",
			Email:     "MerryMage@users.noreply.github.com",
			Timestamp: 1530357667,
			Timezone:  "+0100",
		},
		Committer: Signature{
			Name:      "MerryMage",
			Email:     "MerryMage@users.noreply.github.com",
			Timestamp: 1530361234,
			Timezone:  "+0100",
		},
		Tree: oid("a89cfbe149655be48512d8e748a0b0b779e85e69"),
		Parents: []objstore.Oid{
			oid("86f7c4cb3308289c79ba50076a9160d99755ad66"),
		},
		Message: "A64: Implement FCVTMU (scalar)\n",
	}

	err := commit.Write(&b)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(b.Bytes(), testCommit) {
		t.Errorf("b.Bytes() = %#v", b.Bytes())
	}
}
