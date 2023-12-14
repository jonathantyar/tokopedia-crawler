package migrations

import (
	"context"
	"database/sql"
	"jonathantyar/tokopedia-crawler/src/database"
	"jonathantyar/tokopedia-crawler/src/model"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInit, downInit)
}

func upInit(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	db, err := database.InitDB()
	if err != nil {
		return err
	}
	err = db.AutoMigrate(
		&model.Product{},
	)
	if err != nil {
		db.Rollback()
		return err
	}

	return nil
}

func downInit(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	db, err := database.InitDB()
	if err != nil {
		return err
	}
	err = db.Migrator().DropTable(
		&model.Product{},
	)
	if err != nil {
		return err
	}
	return nil

}
