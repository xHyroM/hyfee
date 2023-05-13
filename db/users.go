package db

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/oauth2"
	"github.com/disgoorg/snowflake/v2"
)

type User struct {
	Id 		       string                `bun:"id,pk"`
	AccessToken  string                `bun:"access_token,notnull"`
	RefreshToken string                `bun:"refresh_token,notnull"`
	Scopes       []discord.OAuth2Scope `bun:"scopes,notnull"`
	TokenType    discord.TokenType     `bun:"token_type,notnull"`
	Expiration   time.Time             `bun:"expiration,notnull"`
}

type UsersDB interface {
	Get(userId string) (User, error)
	AddUser(userId snowflake.ID, session oauth2.Session) (error)
	UpdateUser(userId snowflake.ID, session oauth2.Session) (error)
}

func (s *sqlDB) Get(userId string) (user User, err error) {
	err = s.db.NewSelect().
		Model(&user).
		Where("id = ?", userId).
		Scan(context.TODO())
	return
}

func (s *sqlDB) AddUser(userId snowflake.ID, session oauth2.Session) (err error) {
	_, err = s.db.NewInsert().
		Model(&User{
			Id: 				         userId.String(),
			AccessToken:         session.AccessToken,
			RefreshToken:        session.RefreshToken,
			Scopes:              session.Scopes,
			TokenType:           session.TokenType,
			Expiration:          session.Expiration,
		}).
		Exec(context.TODO())
	return
}

func (s *sqlDB) UpdateUser(userId snowflake.ID, session oauth2.Session) (err error) {
	_, err = s.db.NewUpdate().
		Model(&User{
			Id: 				         userId.String(),
			AccessToken:         session.AccessToken,
			RefreshToken:        session.RefreshToken,
			Scopes:              session.Scopes,
			TokenType:           session.TokenType,
			Expiration:          session.Expiration,
		}).
		Where("id = ?", userId.String()).
		Exec(context.TODO())
	return
}