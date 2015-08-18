// (c) Chi Vinh Le <cvl@chinet.info> – 13.06.2015

package main

import (
	"flag"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"log"
	"net/http"
	"text/template"
)

type Session struct {
	tempCred *oauth.Credentials
	client   *oauth.Client
}

var session = map[string]Session{}

// respond responds to a request by executing the html template t with data.
func respond(w http.ResponseWriter, t *template.Template, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		log.Print(err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	respond(w, homeTmpl, nil)
}
func serveAuthenticate(w http.ResponseWriter, r *http.Request) {
	p := r.FormValue("provider")
	k := r.FormValue("key")
	s := r.FormValue("secret")

	if p != "bitbucket" {
		w.WriteHeader(400)
		w.Write([]byte("unsupported"))
		return
	}

	var client = &oauth.Client{
		Credentials:                   oauth.Credentials{k, s},
		TemporaryCredentialRequestURI: "https://bitbucket.org/api/1.0/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://bitbucket.org/api/1.0/oauth/authenticate",
		TokenRequestURI:               "https://bitbucket.org/api/1.0/oauth/access_token",
	}

	callback := "http://" + r.Host + "/callback"
	tempCred, err := client.RequestTemporaryCredentials(nil, callback, nil)
	if err != nil {
		http.Error(w, "Error getting temp cred, "+err.Error(), 500)
		return
	}
	session[tempCred.Token] = Session{tempCred, client}
	http.Redirect(w, r, client.AuthorizationURL(tempCred, nil), 302)
}

// serveOAuthCallback handles callbacks from the OAuth server.
func serveOAuthCallback(w http.ResponseWriter, r *http.Request) {
	v := r.FormValue("oauth_verifier")
	t := r.FormValue("oauth_token")

	if _, ok := session[t]; !ok {
		http.Error(w, "Unknown oauth_token.", 500)
		return
	}
	client := session[t].client
	tempCred := session[t].tempCred
	tokenCred, _, err := client.RequestToken(nil, tempCred, v)
	if err != nil {
		http.Error(w, "Error getting request token, "+err.Error(), 500)
		return
	}
	delete(session, t)
	respond(w, tokenTmpl, tokenCred)
}

func main() {
	var httpAddr = flag.String("addr", "0.0.0.0:8080", "HTTP server address")
	flag.Parse()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/authenticate", serveAuthenticate)
	http.HandleFunc("/callback", serveOAuthCallback)
	fmt.Printf("Listening on %s\n", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatalf("Error listening, %v", err)
	}
}

var (
	homeTmpl = template.Must(template.New("home").Parse(
		`<html>
<body>
<h1>OAuth1 Authenticator</h1>
<form action="/authenticate" method="post">
  Provider: <select name="provider"><option>bitbucket</option><option>github</option></select><br/>
  Consumer key: <input type="text" name="key" /><br/>
  Consumer secret: <input type="password" name="secret" /><br/>
  <input type="submit" />
</form>
</body>
</html>`))
	tokenTmpl = template.Must(template.New("token").Parse(
		`<html>
<body>
<h1>OAuth1 Authenticator</h1>
  <p>
    Access Token: {{.Token}}<br/>
    Access Token Secret: {{.Secret}}<br/>
</body>
</html>`))
)
