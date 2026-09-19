package database

import (
	"fmt"
	"path"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DBSettings interface {
	GetName() string
}

func Open(dir string, dbs DBSettings) gorm.Dialector {
	return sqlite.Open(path.Join(dir, fmt.Sprintf("%s.db", dbs.GetName())))
}
