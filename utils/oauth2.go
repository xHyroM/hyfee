package utils

import (
	"hyros_coffee/db"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/oauth2"
	"github.com/disgoorg/snowflake/v2"
)

var cache Cache = *NewCache()

func GetGuilds(client oauth2.Client, database db.DB, userId snowflake.ID, originalSession oauth2.Session) []discord.OAuth2Guild {
	cached, found := cache.Get(userId.String())
	if found {
		return cached.([]discord.OAuth2Guild)
	}
	
	session, err := client.VerifySession(originalSession)
	if err != nil {
		return []discord.OAuth2Guild{}
	}

	if session.AccessToken != originalSession.AccessToken {
		database.UpdateUser(userId, session)
	}

	guilds, err := client.GetGuilds(session)

	if err != nil {
		return []discord.OAuth2Guild{}
	}

	cache.Set(userId.String(), guilds, 180000 * 1000000)

	return guilds
}