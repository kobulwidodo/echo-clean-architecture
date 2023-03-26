package config

import (
	"go-clean/src/lib/sql"
)

type Application struct {
	Meta ApplicationMeta
	SQL  sql.Config
}

type ApplicationMeta struct {
	Title       string
	Description string
	Host        string
	BasePath    string
	Version     string
}

type CORSConfig struct {
	Mode string
}

func Init() Application {
	return Application{}
}
