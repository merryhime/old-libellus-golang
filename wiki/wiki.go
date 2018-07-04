package wiki

import (
	"net/http"
	"path"

	"github.com/MerryMage/libellus/common"
)

type Wiki struct {
	config *common.Config
}

func NewWiki(config *common.Config) *Wiki {
	return &Wiki{
		config: config,
	}
}

func (wiki *Wiki) invalidPathResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("404"))
}

func validatePath(p *string) bool {
	*p = path.Clean(*p)

	if (*p)[0] != '/' {
		return false
	}

	for _, ch := range *p {
		if (ch < 'a' || ch > 'z') && (ch < 'A' || ch > 'Z') && (ch < '0' || ch > '9') && ch != '/' && ch != '-' {
			return false
		}
	}

	return true
}

func (wiki *Wiki) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if !validatePath(&path) {
		wiki.invalidPathResponse(w, r)
		return
	}

	w.Write([]byte("we're ok: " + path))
}
