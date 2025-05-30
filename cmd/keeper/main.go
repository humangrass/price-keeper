// @title Price Keeper API
// @version 1.0
// @description API for cryptocurrency tokens
// @host localhost:8888
// @BasePath /api/
package main

import (
	"context"
	"log"
	"os"

	"github.com/humangrass/price-keeper/config"
	"github.com/humangrass/price-keeper/domain/repository"
	"github.com/humangrass/price-keeper/internal/instance"
	"github.com/humangrass/price-keeper/internal/usecases/keeper"
	"github.com/humangrass/price-keeper/internal/usecases/plodder"

	"github.com/humangrass/gommon/signal"
	"github.com/urfave/cli/v2"
)

func main() {
	application := cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config-file",
				Required: false,
				Value:    "keeper.example.yaml",
				Usage:    "Keeper config yaml file",
				EnvVars:  []string{"KEEPER_CONFIG_FILE"},
			},
		},
		Action: Main,
		After: func(c *cli.Context) error {
			log.Println("stopped")
			return nil
		},
	}

	if err := application.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func Main(ctx *cli.Context) error {
	appContext, cancel := context.WithCancel(ctx.Context)
	defer func() {
		cancel()
	}()

	cfg, err := config.NewKeeperConfig(ctx.String("config-file"))
	if err != nil {
		return err
	}

	inst, err := instance.NewInstance(appContext, &instance.Opt{
		Database:     &cfg.Database,
		IsProduction: cfg.IsProduction,
		XHttpOpt:     cfg.Server,
	})
	if err != nil {
		return err
	}

	await, _ := signal.Notifier(func() {
		inst.Logger.Sugar().Info("recieve stop signal, start shutdown process..")
	})

	priceRepo := repository.NewPricesRepository(inst.Pool)
	tokenRepo := repository.NewTokensRepository(inst.Pool)
	pairRepo := repository.NewPairsRepository(inst.Pool)
	ucKeeper := keeper.NewKeeperUseCase(priceRepo, tokenRepo, pairRepo, inst.Logger)
	ucPlodder := plodder.NewPlodderUseCase(pairRepo, priceRepo, inst.Logger, cfg.RefreshInterval)

	ucKeeper.RegisterRoutes(inst.Server.Mux)
	inst.Server.Start()

	go ucPlodder.Run(appContext)

	return await()
}
