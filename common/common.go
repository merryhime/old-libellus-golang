package common

import (
	"github.com/gobuffalo/packr"

	"github.com/MerryMage/libellus/auth"
	"github.com/MerryMage/libellus/objstore"
	"github.com/MerryMage/libellus/wikidata"
)

type Config struct {
	HttpOnly        bool
	CanonicalDomain string
	PrivateWikiDir  string
	PrivateSrsDir   string
	Repo            *objstore.Repository
	Authentication  *auth.Auth
	StaticData      packr.Box
	WikiData        *wikidata.WikiData
}
