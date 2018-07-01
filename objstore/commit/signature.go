package commit

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

var SignatureError error = errors.New("commit: signature invalid format")

type Signature struct {
	Name      string
	Email     string
	Timestamp int64
	Timezone  string
}

func NewSignature(b []byte) (Signature, error) {
	var sig Signature

	emailstart := bytes.IndexByte(b, byte('<'))
	emailend := bytes.IndexByte(b, byte('>'))
	if emailstart == -1 || emailend == -1 || emailstart >= emailend {
		return sig, SignatureError
	}

	sig.Name = string(bytes.Trim(b[:emailstart], " "))
	sig.Email = string(b[emailstart+1 : emailend])

	if len(b[emailend+2:]) > 1 {
		timestr := bytes.Split(b[emailend+2:], []byte(" "))
		ts, err := strconv.ParseInt(string(timestr[0]), 10, 64)
		if err != nil {
			return sig, err
		}
		sig.Timestamp = ts
		if len(timestr) == 2 {
			sig.Timezone = string(timestr[1])
		} else if len(timestr) > 2 {
			return sig, SignatureError
		}
	}

	return sig, nil
}

func (sig Signature) String() string {
	if sig.Timezone != "" {
		return fmt.Sprintf("%s <%s> %d %s", sig.Name, sig.Email, sig.Timestamp, sig.Timezone)
	}
	return fmt.Sprintf("%s <%s> %d", sig.Name, sig.Email, sig.Timestamp)
}
