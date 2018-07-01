package tree

import (
	"testing"

	"bytes"
	"encoding/hex"

	"github.com/MerryMage/libellus/objstore/filemode"
	"github.com/MerryMage/libellus/objstore/objid"
)

func oid(s string) objid.Oid {
	oid, err := objid.FromString(s)
	if err != nil {
		panic("oid failed")
	}
	return oid
}

func TestWrite(t *testing.T) {
	var tree Tree
	tree.Add(Entry{
		Mode: filemode.Regular,
		Name: ".gitignore",
		Oid:  oid("f1c181ec9c5c921245027c6b452ecfc1d3626364"),
	})
	tree.Add(Entry{
		Mode: filemode.Dir,
		Name: "objstore",
		Oid:  oid("80eafdf2e2c6c9dd03ae2fbab2bf8d72bd58c688"),
	})
	tree.Add(Entry{
		Mode: filemode.Regular,
		Name: "README.md",
		Oid:  oid("1c8392b7465cdf257da626d81dfc0641b34310de"),
	})
	tree.Add(Entry{
		Mode: filemode.Regular,
		Name: "main.go",
		Oid:  oid("f6f20359ce44134b101df38857f131530a3fe1ab"),
	})
	tree.Add(Entry{
		Mode: filemode.Regular,
		Name: "LICENSE",
		Oid:  oid("ea06ac37261238e918f74d9ade554e0c5cb2e107"),
	})

	var b bytes.Buffer
	tree.Write(&b)

	expectedResult, _ := hex.DecodeString("313030363434202e67697469676e6f726500f1c181ec9c5c921245027c6b452ecfc1d3626364313030363434204c4943454e534500ea06ac37261238e918f74d9ade554e0c5cb2e10731303036343420524541444d452e6d64001c8392b7465cdf257da626d81dfc0641b34310de313030363434206d61696e2e676f00f6f20359ce44134b101df38857f131530a3fe1ab3430303030206f626a73746f72650080eafdf2e2c6c9dd03ae2fbab2bf8d72bd58c688")

	if !bytes.Equal(b.Bytes(), expectedResult) {
		t.Errorf("b.Bytes() = %#v", b.Bytes())
	}
}

func TestRead(t *testing.T) {
	input, _ := hex.DecodeString("313030363434202e67697469676e6f726500f1c181ec9c5c921245027c6b452ecfc1d3626364313030363434204c4943454e534500ea06ac37261238e918f74d9ade554e0c5cb2e10731303036343420524541444d452e6d64001c8392b7465cdf257da626d81dfc0641b34310de313030363434206d61696e2e676f00f6f20359ce44134b101df38857f131530a3fe1ab3430303030206f626a73746f72650080eafdf2e2c6c9dd03ae2fbab2bf8d72bd58c688")
	b := bytes.NewBuffer(input)

	tree, err := Read(b)
	if err != nil {
		t.Error(err)
	}

	if tree.Entries[0].Mode != filemode.Regular {
		t.Errorf("tree.Entries[0].Mode = %#v", tree.Entries[0].Mode)
	}

	if tree.Entries[0].Name != ".gitignore" {
		t.Errorf("tree.Entries[0].Name = %#v", tree.Entries[0].Name)
	}

	if tree.Entries[0].Oid.String() != "f1c181ec9c5c921245027c6b452ecfc1d3626364" {
		t.Errorf("tree.Entries[0].Oid.Stirng() = %#v", tree.Entries[0].Oid.String())
	}

	if tree.Entries[1].Name != "LICENSE" {
		t.Errorf("tree.Entries[1].Name = %#v", tree.Entries[1].Name)
	}

	if tree.Entries[2].Name != "README.md" {
		t.Errorf("tree.Entries[2].Name = %#v", tree.Entries[2].Name)
	}

	if tree.Entries[3].Name != "main.go" {
		t.Errorf("tree.Entries[3].Name = %#v", tree.Entries[3].Name)
	}

	if tree.Entries[4].Name != "objstore" {
		t.Errorf("tree.Entries[4].Name = %#v", tree.Entries[4].Name)
	}
}
