package state

import (
	"gator/internal/config"
	"gator/internal/database"
)

type State struct {
	Config_ptr *config.Config
	Db_ptr     *database.Queries
}
