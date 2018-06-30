package tree

import (
	"io"

	"github.com/MerryMage/libellus/objstore"
	"github.com/MerryMage/libellus/objstore/filemode"
	"github.com/MerryMage/libellus/objstore/ioutil"
)

func Read(r io.Reader) (Tree, error) {
	ret := Tree{}

	for {
		raw, err := ioutil.ReadUntil(r, ' ')
		if err != nil {
			if err == io.EOF {
				break
			}
			return ret, err
		}

		mode, err := filemode.New(string(raw))
		if err != nil {
			return ret, err
		}

		raw, err = ioutil.ReadUntil(r, 0)
		if err != nil {
			return ret, err
		}
		name := string(raw)

		oid, err := objstore.ReadOid(r)
		if err != nil {
			return ret, err
		}

		ret.Add(Entry{
			Mode: mode,
			Name: name,
			Oid:  oid,
		})
	}

	return ret, nil
}
