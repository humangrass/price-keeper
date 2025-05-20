package main

import (
	"log"
	"os"

	"github.com/humangrass/price-keeper/config"
	"github.com/humangrass/price-keeper/internal/usecases/jupiter"

	"github.com/urfave/cli/v2"
)

func main() {
	application := cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "jupiter-token-id",
				Required: false,
				Value:    "So11111111111111111111111111111111111111112",
				Usage:    "Token ID in Solana",
				EnvVars:  []string{"JUPITER_TOKEN_ID"},
			},
			&cli.StringFlag{
				Name:     "jupiter-vs-token-id",
				Required: false,
				Usage:    "VS Token ID in Solana. By default USD",
				EnvVars:  []string{"JUPITER_VS_TOKEN_ID"},
			},
			&cli.BoolFlag{
				Name:     "jupiter-extra-info",
				Required: false,
				Value:    false,
				Usage:    "Show extra info. Cannot use with VS Token ID param",
				EnvVars:  []string{"JUPITER_EXTRA_INFO"},
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
	cfg := config.NewJupiterCfg(ctx)
	err := cfg.IsValid()
	if err != nil {
		return err
	}

	client := jupiter.NewClient(*cfg)

	_, err = client.GetPrice()
	return err
}
