package instance

import (
	"context"
	"fmt"

	"github.com/humangrass/price-keeper/pgk/logger"
	"github.com/humangrass/price-keeper/pgk/xhttp"

	"github.com/humangrass/gommon/database"
	"github.com/humangrass/gommon/database/postgres"
	"github.com/humangrass/gommon/drop"
)

type Instance struct {
	*drop.Impl

	Server *xhttp.Server
	Logger *logger.Logger
	Pool   database.Pool
}

func NewInstance(ctx context.Context, opt *Opt) (*Instance, error) {
	i := &Instance{}
	i.Impl = drop.NewContext(ctx)

	var err error

	i.Logger, err = logger.New(opt.IsProduction)
	if err != nil {
		return nil, err
	}
	i.AddDropper(i.Logger)

	server, err := xhttp.NewServer(opt.XHttpOpt, i.Logger.Sugar())
	if err != nil {
		return nil, err
	}
	i.Server = server
	i.AddDropper(i.Server)

	if opt.Database.Dialect == "postgres" {
		i.Pool, err = postgres.NewPool(i.Context(), opt.Database)
		if err != nil {
			return nil, err
		}
		i.AddDropper(i.Pool.(*postgres.Pool))
	} else {
		return nil, fmt.Errorf("database dialect %s not supported", opt.Database.Dialect)
	}

	return i, nil
}
