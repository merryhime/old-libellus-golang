package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func generateRandomString() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), err
}

type config struct {
	Username string
	Password string
}

type rateLimitEntry struct {
	Count  int
	Expiry time.Time
}

type Auth struct {
	config

	httpOnly  bool
	cookies   map[string]bool
	ratelimit map[string]rateLimitEntry
}

func NewAuth(configFile string, httpOnly bool) *Auth {
	auth := &Auth{
		httpOnly:  httpOnly,
		cookies:   make(map[string]bool),
		ratelimit: make(map[string]rateLimitEntry),
	}
	raw, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Println(err)
		return auth
	}
	err = json.Unmarshal(raw, &auth.config)
	if err != nil {
		log.Println(err)
		auth.Username = ""
		auth.Password = ""
	}
	return auth
}

func (auth *Auth) getCookie(r *http.Request) *string {
	cookie, err := r.Cookie("libellus")
	if err != nil {
		return nil
	}
	c := cookie.Value
	return &c
}

func (auth *Auth) IsAuthenticated(r *http.Request) bool {
	cookie := auth.getCookie(r)
	if cookie == nil {
		return false
	}
	if value, ok := auth.cookies[*cookie]; ok {
		return value
	}
	return false
}

func (auth *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/_auth/login":
		auth.login(w, r)
	case "/_auth/logout":
		auth.logout(w, r)
	case "/_auth/clear":
		auth.clear(w, r)
	default:
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("bad authentication url"))
	}
}

func (auth *Auth) login(w http.ResponseWriter, r *http.Request) {
	if auth.IsAuthenticated(r) {
		w.Write([]byte("you're already authenticated, silly"))
		return
	}

	if err := r.ParseForm(); err != nil {
		w.Write([]byte("ParseForm failure"))
		return
	}

	outputForm := func(msg string) {
		w.Write([]byte(`
<html>
    <div>` + msg + `</div>
    <form action="/_auth/login" method="POST">
        <input type="hidden" name="redirect" value="` + html.EscapeString(r.Form.Get("redirect")) + `" />
        <div>Username: <input type="text" name="username" /></div>
        <div>Passsword: <input type="password" name="password" /></div>
        <div><input type="submit"></div>
    </form>
</html>
`))
	}

	if r.Method == http.MethodGet {
		outputForm("")
		return
	}

	if r.Method == http.MethodPost {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("SplitHostPort"))
			return
		}

		if auth.ratelimit[ip].Expiry.After(time.Now()) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("rate-limited"))
			return
		}

		rle := auth.ratelimit[ip]
		rle.Count += 1
		auth.ratelimit[ip] = rle
		if rle.Count >= 5 {
			auth.ratelimit[ip] = rateLimitEntry{
				Expiry: time.Now().Add(time.Hour),
			}
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("rate-limited"))
			return
		}

		if r.Form.Get("username") != auth.Username {
			outputForm("invalid login")
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(auth.Password), []byte(r.Form.Get("password"))) != nil {
			outputForm("invalid login")
			return
		}

		newcookie, err := generateRandomString()
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("bad luck"))
			return
		}

		auth.cookies[newcookie] = true
		http.SetCookie(w, &http.Cookie{
			Name:    "libellus",
			Value:   newcookie,
			Expires: time.Now().Add(14 * 24 * time.Hour),
			Secure:  !auth.httpOnly,
		})
		http.Redirect(w, r, "/"+r.Form.Get("redirect"), 303)

		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("invalid method"))
}

func (auth *Auth) logout(w http.ResponseWriter, r *http.Request) {
	if !auth.IsAuthenticated(r) {
		w.Write([]byte("you're not even logged-in, silly"))
		return
	}

	if r.Method == http.MethodGet {
		w.Write([]byte(`
<html>
    <form action="/_auth/logout" method="post">
        <button type="submit">Logout</button>
    </form>
</html>
`))
		return
	}

	if r.Method == http.MethodPost {
		cookie := auth.getCookie(r)
		if cookie == nil {
			w.Write([]byte("cookie == nil. Internal error?"))
			return
		}
		delete(auth.cookies, *cookie)
		w.Write([]byte("Logged out."))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("invalid method"))
}

func (auth *Auth) clear(w http.ResponseWriter, r *http.Request) {
	if !auth.IsAuthenticated(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
		return
	}

	if r.Method == http.MethodGet {
		w.Write([]byte(`
<html>
    <form action="/_auth/clear" method="post">
        <button type="submit">Clear Auth</button>
    </form>
</html>
`))
		return
	}

	if r.Method == http.MethodPost {
		num := len(auth.cookies)
		auth.cookies = make(map[string]bool)
		w.Write([]byte(fmt.Sprintf("successful auth clear of %d session(s)", num)))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("invalid method"))
}
