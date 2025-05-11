package config

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type Jupiter struct {
	TokenID   string
	VSTokenID string
	ExtraInfo bool
}

func (cfg Jupiter) IsValid() error {
	if cfg.VSTokenID != "" && cfg.ExtraInfo {
		return fmt.Errorf("cannot use both --jupiter-vs-token-id and --jupiter-extra-info")
	}
	return nil
}

func NewJupiterCfg(ctx *cli.Context) *Jupiter {
	return &Jupiter{
		TokenID:   ctx.String("jupiter-token-id"),
		VSTokenID: ctx.String("jupiter-vs-token-id"),
		ExtraInfo: ctx.Bool("jupiter-extra-info"),
	}
}
