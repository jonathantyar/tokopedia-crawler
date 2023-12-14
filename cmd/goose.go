package cmd

import (
	"errors"
	"fmt"
	"io/fs"
	"jonathantyar/tokopedia-crawler/src/constant"
	"jonathantyar/tokopedia-crawler/src/database"

	_ "jonathantyar/tokopedia-crawler/src/migration"
	"log"
	"os"
	"regexp"
	"text/template"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var dir, table, schema string
var verbose, help, version, sequential, allowMissing, noVersioning bool

var GooseCmd = &cobra.Command{
	Use:   "goose",
	Short: "Use to migrate table to database",
	Run:   gooseCommand,
}

func init() {
	GooseCmd.Flags().StringVarP(&dir, "dir", "d", defaultMigrationDir, "migrations directory")
	GooseCmd.Flags().StringVarP(&table, "table", "t", "goose_db_version", "migrations table name")
	GooseCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
	GooseCmd.Flags().BoolVarP(&help, "help", "h", false, "print help")
	GooseCmd.Flags().BoolVarP(&version, "version", "", false, "print version")
	GooseCmd.Flags().BoolVarP(&sequential, "sequential", "s", false, "use sequential numbering for new migrations")
	GooseCmd.Flags().BoolVarP(&allowMissing, "allow-missing", "", false, "applies missing (out-of-order) migrations")
	GooseCmd.Flags().BoolVarP(&noVersioning, "no-versioning", "", false, "apply migration commands with no versioning, in file order, from directory pointed to")
	GooseCmd.Flags().StringVarP(&schema, "schema", "", "migration", "for type (migration) or (seeder)")
}

type stdLogger struct{}

func (*stdLogger) Fatal(v ...interface{})                 { fmt.Println(v...) }
func (*stdLogger) Fatalf(format string, v ...interface{}) { fmt.Printf(format, v...) }
func (*stdLogger) Print(v ...interface{})                 { fmt.Print(v...) }
func (*stdLogger) Println(v ...interface{})               { fmt.Println(v...) }
func (*stdLogger) Printf(format string, v ...interface{}) { fmt.Printf(format, v...) }

func gooseCommand(cmd *cobra.Command, args []string) {
	if verbose {
		goose.SetVerbose(true)
	}
	if sequential {
		goose.SetSequential(true)
	}
	if schema != "" {
		dir = "./src/migration"
	}
	goose.SetTableName(table)

	switch args[0] {
	case "init":
		if err := gooseInit(dir); err != nil {
			fmt.Printf("goose run: %v", err)
		}
		return
	case "create":
		if err := goose.Run("create", nil, dir, args[1:]...); err != nil {
			fmt.Printf("goose run: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, dir); err != nil {
			fmt.Printf("goose run: %v", err)
		}
		return
	}

	args = mergeArgs(args)
	if len(args) < 1 {
		return
	}

	goose.SetLogger(&stdLogger{})

	command := args[0]
	config := database.Config{
		Host:     viper.GetString(constant.DB_HOST),
		Port:     viper.GetInt(constant.DB_PORT),
		User:     viper.GetString(constant.DB_USER),
		Password: viper.GetString(constant.DB_PASS),
		Database: viper.GetString(constant.DB_NAME),
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := goose.OpenDBWithDriver("postgres", connStr)
	if err != nil {
		fmt.Printf("-dbstring %v, %v\n", connStr, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	options := []goose.OptionsFunc{}
	if allowMissing {
		options = append(options, goose.WithAllowMissing())
	}
	if noVersioning {
		options = append(options, goose.WithNoVersioning())
	}
	log.Println(dir)
	if err := goose.RunWithOptions(
		command,
		db,
		dir,
		arguments,
		options...,
	); err != nil {
		fmt.Printf("goose run: %v", err)
	}
}

func normalizeDBString(driver string, str string, certfile string, sslcert string, sslkey string) string {
	return str
}

var tlsReg = regexp.MustCompile(`(\?|&)tls=[^&]*(?:&|$)`)

const (
	envGooseDriver       = "GOOSE_DRIVER"
	envGooseDBString     = "GOOSE_DBSTRING"
	envGooseMigrationDir = "GOOSE_MIGRATION_DIR"
)

const (
	defaultMigrationDir = "src/migration"
)

func mergeArgs(args []string) []string {
	if len(args) < 1 {
		return args
	}
	if d := os.Getenv(envGooseDriver); d != "" {
		args = append([]string{d}, args...)
	}
	if d := os.Getenv(envGooseDBString); d != "" {
		args = append([]string{args[0], d}, args[1:]...)
	}
	return args
}

func usage() {
	fmt.Println(usagePrefix)
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND
or
Set environment key
GOOSE_DRIVER=DRIVER
GOOSE_DBSTRING=DBSTRING
Usage: goose [OPTIONS] COMMAND
Examples:
    goose mysql "user:password@/dbname?parseTime=true" status
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
`
)

var sqlMigrationTemplate = template.Must(template.New("goose.sql-migration").Parse(`-- Thank you for giving goose a try!
--
-- This file was automatically created running goose init. If you're familiar with goose
-- feel free to remove/rename this file, write some SQL and goose up. Briefly,
--
-- Documentation can be found here: https://pressly.github.io/goose
--
-- A single goose .sql file holds both Up and Down migrations.
--
-- All goose .sql files are expected to have a -- +goose Up directive.
-- The -- +goose Down directive is optional, but recommended, and must come after the Up directive.
--
-- The -- +goose NO TRANSACTION directive may be added to the top of the file to run statements
-- outside a transaction. Both Up and Down migrations within this file will be run without a transaction.
--
-- More complex statements that have semicolons within them must be annotated with
-- the -- +goose StatementBegin and -- +goose StatementEnd directives to be properly recognized.
--
-- Use GitHub issues for reporting bugs and requesting features, enjoy!
-- +goose Up
SELECT 'up SQL query';
-- +goose Down
SELECT 'down SQL query';
`))

// initDir will create a directory with an empty SQL migration file.
func gooseInit(dir string) error {
	if dir == "" || dir == defaultMigrationDir {
		dir = "migrations"
	}
	_, err := os.Stat(dir)
	switch {
	case errors.Is(err, fs.ErrNotExist):
	case err == nil, errors.Is(err, fs.ErrExist):
		return fmt.Errorf("directory already exists: %s", dir)
	default:
		return err
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return goose.CreateWithTemplate(nil, dir, sqlMigrationTemplate, "initial", "sql")
}
