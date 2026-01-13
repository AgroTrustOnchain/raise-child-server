package db

import (
	"errors"
	"log"

	action_type "raise-child/constants/action_type"
	"raise-child/constants/noti"

	"github.com/golang-migrate/migrate/v4"
)

// Database migration
func MigrateDB(action, scriptFolderDir string, version int, server ISQLServer, logger *log.Logger) error {
	// Identify action
	var refixDir string
	switch action {
	case action_type.MIGRATION_TYPE:
		refixDir = "migration"
	case action_type.ROLLBACK_TYPE:
		refixDir = "rollback"
	default: // No method found
		return errors.New("Invalid migration command.")
	}

	// Initialize migration
	migration, err := migrate.New(
		scriptFolderDir+refixDir,
		server.GetCnnStr(),
	)

	if err != nil {
		logger.Println(noti.DB_MIGRATION_ERR_MSG + err.Error())
		return errors.New(noti.DB_MIGRATION_INFORM_MSG)
	}

	if err := migration.Migrate(uint(version)); err != nil && err != migrate.ErrNoChange {
		logger.Println(noti.DB_CONNECTION_ERR_MSG + err.Error())
		return errors.New(noti.DB_MIGRATION_INFORM_MSG)
	}

	log.Println(action + "successful.")
	return nil
}
