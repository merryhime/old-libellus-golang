package common

import (
	"github.com/MerryMage/libellus/auth"
	"github.com/MerryMage/libellus/objstore"
)

type Config struct {
	HttpOnly        bool
	CanonicalDomain string
	PrivateWikiDir  string
	PrivateSrsDir   string
	Repo            *objstore.Repository
	Authentication  *auth.Auth
}
