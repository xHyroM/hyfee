package hyfee

import (
	"net/http"
	"os"

	"github.com/disgoorg/disgo/oauth2"
)

func (b *Bot) SetupOAuth2() {
	b.OAuth2Client = oauth2.New(b.Client.ApplicationID(), os.Getenv("CLIENT_SECRET"))
}

func (b *Bot) OAuthHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")
	state := query.Get("state")

	session, _, err := b.OAuth2Client.StartSession(code, state)
	if err != nil {
		b.Logger.Error("Failed to start session ", err)
		return
	}

	oauthUser, err := b.OAuth2Client.GetUser(session)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = b.Database.Get(oauthUser.ID.String())
	if err != nil {
		if err = b.Database.AddUser(oauthUser.ID, session); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	if err = b.Database.UpdateUser(oauthUser.ID, session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return
}