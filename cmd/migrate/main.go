// nolint:cyclop // it's ok
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	_ "github.com/lib/pq"

	"github.com/humangrass/price-keeper/config"

	"github.com/pressly/goose/v3"
	"github.com/urfave/cli/v2"
)

type migrateOptions struct {
	db  *sql.DB
	dir string
}

func initOptions(configPath string) (*migrateOptions, error) {
	cfg, err := config.NewKeeperConfig(configPath)
	if err != nil {
		return nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	db, err := sql.Open("postgres", cfg.Database.ConnectionString())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &migrateOptions{
		db:  db,
		dir: path.Join(wd, "migrations"),
	}, nil
}

// nolint:lll,funlen,cyclop,gocognit // it's ok
func main() {
	application := cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config-file",
				Required: false,
				Value:    "config.yaml",
				Usage:    "YAML config filepath",
				EnvVars:  []string{"ASIS_CONFIG_FILE"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "up",
				Usage: "applies all available migrations",
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.UpContext(ctx.Context, opt.db, opt.dir); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
			},
			{
				Name:  "down",
				Usage: "rolls back a single migration from the current version",
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.DownContext(ctx.Context, opt.db, opt.dir); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
				After: func(c *cli.Context) error {
					log.Println("successfully")
					return nil
				},
			},
			{
				Name:  "up-to",
				Usage: "migrates up to a specific version",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:    "version",
						Aliases: []string{"v"},
						Value:   0,
						Usage:   "version to up",
					},
				},
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.UpToContext(ctx.Context, opt.db, opt.dir, ctx.Int64("version")); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
				After: func(c *cli.Context) error {
					log.Println("successfully")
					return nil
				},
			},
			{
				Name:  "down-to",
				Usage: "rolls back migrations to a specific version",
				Flags: []cli.Flag{
					&cli.Int64Flag{
						Name:    "version",
						Aliases: []string{"v"},
						Value:   0,
						Usage:   "version to down",
					},
				},
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.DownToContext(ctx.Context, opt.db, opt.dir, ctx.Int64("version")); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
				After: func(c *cli.Context) error {
					log.Println("successfully")
					return nil
				},
			},
			{
				Name:  "up-by-one",
				Usage: "migrates up by a single version",
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.UpByOneContext(ctx.Context, opt.db, opt.dir); err != nil {
						if errors.Is(err, goose.ErrNoNextVersion) {
							return nil
						}
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
			},
			{
				Name:  "create",
				Usage: "creates new migration file with the current timestamp",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "",
						Usage:   "name",
					},
				},
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.Create(opt.db, opt.dir, ctx.String("name"), "sql"); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "dump the migration status for the current DB",
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.StatusContext(ctx.Context, opt.db, opt.dir); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
			},
			{
				Name:  "version",
				Usage: "print the current version of the database",
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.VersionContext(ctx.Context, opt.db, opt.dir); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
			},
			{
				Name:  "redo",
				Usage: "re-run the latest migration",
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.RedoContext(ctx.Context, opt.db, opt.dir); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
			},
			{
				Name:  "reset",
				Usage: "roll back all migrations",
				Action: func(ctx *cli.Context) error {
					opt, err := initOptions(ctx.String("config-file"))
					if err != nil {
						return err
					}

					if err = goose.ResetContext(ctx.Context, opt.db, opt.dir); err != nil {
						return fmt.Errorf("failed to run migrations: %w", err)
					}
					return nil
				},
			},
		},
	}

	if err := application.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
