package main

import (
	"context"
	"log"
	"os"

	"github.com/humangrass/price-keeper/config"
	"github.com/humangrass/price-keeper/domain/repository"
	"github.com/humangrass/price-keeper/internal/instance"
	"github.com/humangrass/price-keeper/internal/usecases/keeper"

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

	baseRepo := repository.NewBaseRepository(inst.Pool)
	uc := keeper.NewKeeperUseCase(&baseRepo, inst.Logger)

	uc.RegisterRoutes(inst.Server.Mux)
	inst.Server.Start()

	return await()
}
