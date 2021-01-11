package DB

import (
	"github.com/OBASHITechnology/resourceList/DB/impl/core"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres"
)

var Store core.IStore

func init() {
	Store = postgres.New()
}
