package objfile

import (
	"testing"

	"bytes"
	"encoding/hex"
	"io/ioutil"

	"github.com/MerryMage/libellus/objstore"
	"github.com/MerryMage/libellus/objstore/objtype"
)

func TestWrite(t *testing.T) {
	var b bytes.Buffer

	w, err := NewWriter(&b, objtype.Blob, 15)
	if err != nil {
		t.Error(err)
	}

	n, err := w.Write([]byte("This is a test\n"))
	if err != nil || n != 15 {
		t.Error(n, err)
	}

	w.Close()

	expectedOid, _ := objstore.OidFromString("0527e6bd2d76b45e2933183f1b506c7ac49f5872")
	if !w.Oid().Equals(expectedOid) {
		t.Errorf("w.Oid() = %q", w.Oid())
	}

	readHelper(t, b.Bytes())
}

func readHelper(t *testing.T, input []byte) {
	b := bytes.NewBuffer(input)

	r, err := NewReader(b)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	if r.Size() != 15 {
		t.Errorf("r.Size() = %#v", r.Size())
	}

	if r.ObjType() != objtype.Blob {
		t.Errorf("r.ObjType() = %#v", r.ObjType().String())
	}

	contents, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	expectedContents := []byte("This is a test\n")
	if !bytes.Equal(contents, expectedContents) {
		t.Errorf("contents = %#v", contents)
	}
}

func TestRead(t *testing.T) {
	input, _ := hex.DecodeString("78014bcac94f5230346508c9c82c5600a2448592d4e2122e0055ab0725")
	readHelper(t, input)
}
