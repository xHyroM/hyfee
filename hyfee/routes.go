package hyfee

import (
	"hyros_coffee/utils"
	"net/http"
	"os"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/json"
	_ "github.com/joho/godotenv/autoload"
)

var baseUrl = os.Getenv("BASE_URL")

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

	
	if utils.Contains(session.Scopes, discord.OAuth2ScopeRoleConnectionsWrite) {
		connection, err := b.OAuth2Client.GetApplicationRoleConnection(session, b.Client.ApplicationID())
		if err != nil {
			if _, err := b.OAuth2Client.UpdateApplicationRoleConnection(session, b.Client.ApplicationID(), discord.ApplicationRoleConnectionUpdate{
				PlatformName: json.Ptr("Monitored"),
				Metadata: json.Ptr(map[string]string {
					"since": time.Now().UTC().Format(time.RFC3339),
				}),
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			if _, err := b.OAuth2Client.UpdateApplicationRoleConnection(session, b.Client.ApplicationID(), discord.ApplicationRoleConnectionUpdate{
				PlatformName: json.Ptr("Monitored"),
				Metadata: json.Ptr(map[string]string {
					"since": connection.Metadata["since"],
				}),
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
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

	w.Write([]byte("You have successfully authenicated yourself with Discord OAuth2. You can now close this tab :D"))
	return
}

func (b *Bot) LinkedRolesHandler(w http.ResponseWriter, r *http.Request) {
	url := b.OAuth2Client.GenerateAuthorizationURL(baseUrl+"/callback", discord.PermissionsNone, 0, false, discord.OAuth2ScopeGuilds, discord.OAuth2ScopeRoleConnectionsWrite, discord.OAuth2ScopeIdentify)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}