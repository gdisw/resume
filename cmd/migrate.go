package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func init() {
	migrateCmd.AddCommand(migrateNewCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateForceCmd)
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateResetCmd)
	rootCmd.AddCommand(migrateCmd)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migration creation, up and down.",
}

var migrateNewCmd = &cobra.Command{
	Use:          "new",
	Short:        "Creates a new up and down migrations.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("missing migration name")
		}

		version := time.Now().UTC().Format("20060102150405")
		name := strings.Join(append([]string{version}, args...), "_")

		for _, direction := range []string{"up", "down"} {
			path := fmt.Sprintf("migrations/%s.%s.sql", name, direction)

			f, err := os.Create(path)
			if err != nil {
				return err
			}

			if err := f.Close(); err != nil {
				return err
			}
		}

		return nil
	},
}

var migrateVersionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Get the applied version.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := migrator()
		if err != nil {
			return err
		}
		defer m.Close()

		version, dirty, err := m.Version()
		if err != nil {
			return err
		}

		fmt.Println("Version:", version)
		fmt.Println("Dirty:", dirty)

		return nil
	},
}

var migrateForceCmd = &cobra.Command{
	Use:          "force",
	Short:        "Force a migration version.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := migrator()
		if err != nil {
			return err
		}
		defer m.Close()

		version, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		return m.Force(version)
	},
}

var migrateUpCmd = &cobra.Command{
	Use:          "up",
	Short:        "Applies all of the existing up migrations.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := migrator()
		if err != nil {
			return err
		}
		defer m.Close()

		if len(args) > 0 && args[0] == "step" {
			err = m.Steps(1)
		} else {
			err = m.Up()
		}

		if err == migrate.ErrNoChange {
			return nil
		}

		return err
	},
}

var migrateDownCmd = &cobra.Command{
	Use:          "down",
	Short:        "Applies the last down migration.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := migrator()
		if err != nil {
			return err
		}
		defer m.Close()

		return m.Steps(-1)
	},
}

var migrateResetCmd = &cobra.Command{
	Use:          "reset",
	Short:        "Applies all of the existing down migrations.",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		m, err := migrator()
		if err != nil {
			return err
		}
		defer m.Close()

		return m.Down()
	},
}

func migrator() (*migrate.Migrate, error) {
	src := "file://migrations"
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		return nil, errors.New("DATABASE_DSN is not set")
	}

	return migrate.New(src, dsn)
}
