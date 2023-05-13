package db

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type DB interface {
	UsersDB
	Close()
}

func Setup(syncTables *bool, debug *bool) (DB, error) {
	sqlite, err := sql.Open(sqliteshim.ShimName, "file:data/hyfee.db")
	if err != nil {
		panic(err)
	}
	sqlite.SetMaxOpenConns(1)

	db := bun.NewDB(sqlite, sqlitedialect.New())

	if *debug {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	if *syncTables {
		db.NewCreateTable().Model((*User)(nil)).Exec(context.TODO());
	}

	return &sqlDB{db: db}, nil
}

type sqlDB struct {
	db *bun.DB
}

func (s *sqlDB) Close() {
	s.db.Close()
}