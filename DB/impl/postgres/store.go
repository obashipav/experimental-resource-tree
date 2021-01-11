package postgres

import (
	"context"
	"github.com/OBASHITechnology/resourceList/DB/impl/core"
	"github.com/OBASHITechnology/resourceList/DB/impl/postgres/common"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type store struct{ pool *pgxpool.Pool }

func New() core.IStore {
	var (
		ctx, cancel = context.WithTimeout(context.Background(), common.DEFAULT_REQUEST_TTL)
		instance    = &store{}
		err         error
		poolConf    *pgxpool.Config
	)
	defer cancel()

	poolConf, err = pgxpool.ParseConfig("postgres://obashi:obashi@localhost:5432/resource?sslmode=disable&pool_max_conns=5")
	if err != nil {
		log.Fatal("failed to parse the db config: ", err)
	}

	instance.pool, err = pgxpool.ConnectConfig(ctx, poolConf)
	if err != nil {
		log.Fatal("failed to connect to db ", err)
	}

	return instance
}
