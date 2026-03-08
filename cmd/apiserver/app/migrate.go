package app

import (
	"context"
	"fmt"

	"github.com/mgcis-cn/ibookfs/cmd/apiserver/app/options"
	"github.com/mgcis-cn/ibookfs/pkg/database"

	"github.com/spf13/cobra"
)

// NewMigrateCommand creates the migrate subcommand.
func NewMigrateCommand() *cobra.Command {
	runOpts := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long: `Run database schema migrations.

This command applies all pending database migrations to create or update
the database schema. It should be run manually when deploying new versions
that include schema changes.

Example:
  apiserver migrate --env docker
  apiserver migrate --env prod`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			if err := c.Scan(runOpts); err != nil {
				return fmt.Errorf("failed to scan config: %w", err)
			}
			return runMigrate(runOpts)
		},
	}

	fs := cmd.Flags()
	AddConfigFlags(fs)
	return cmd
}

func runMigrate(opts *options.ServerRunOptions) error {
	ctx := context.Background()

	fmt.Println("Initializing database connection...")
	databaseF, err := database.NewFactory(ctx, opts.Data.Database, database.WithNameFunc)
	if err != nil {
		return fmt.Errorf("failed to create database factory: %w", err)
	}
	defer databaseF.Close()

	db, err := databaseF.MustGet("default")
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	fmt.Println("Running database migrations...")
	if err := database.AutoMigrate(db); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Database migrations completed successfully!")
	return nil
}
