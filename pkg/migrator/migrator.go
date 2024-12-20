package migrator

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/backend/bff-cognito/pkg/logger"

	"github.com/caarlos0/env"
	packr "github.com/gobuffalo/packr/v2"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"

	gotdotenv "github.com/joho/godotenv"
)

const (
	migrationDriver = "postgres"
)

type Migrator struct {
	Db         *sql.DB
	Opts       MigratorOpts
	Migrations *migrate.PackrMigrationSource
	l          *logger.Logger
}

type MigratorOpts struct {
	PostgresDns   string `env:"POSTGRES_DNS"`
	MigrationPath string
}

func NewMigrator(l *logger.Logger, opts MigratorOpts) (*Migrator, error) {
	db, err := sql.Open("postgres", opts.PostgresDns)
	if err != nil {
		return nil, err
	}

	f, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	migrationPath := fmt.Sprintf("%s/%s", filepath.Dir(f), filepath.Base(f))
	if opts.MigrationPath != "" {
		migrationPath = opts.MigrationPath
	}
	opts.MigrationPath = migrationPath

	migrations := &migrate.PackrMigrationSource{
		Box: packr.New(migrationDriver, migrationPath),
	}

	return &Migrator{
		Opts:       opts,
		Migrations: migrations,
		Db:         db,
		l:          l,
	}, nil
}

func (m *Migrator) Up() (int, error) {
	n, err := migrate.Exec(m.Db, migrationDriver, m.Migrations, migrate.Up)

	if err != nil {
		return 0, err
	}

	return n, nil
}

func (m *Migrator) Down() (int, error) {
	n, err := migrate.Exec(m.Db, migrationDriver, m.Migrations, migrate.Down)

	if err != nil {
		return 0, err
	}

	return n, nil
}

func (m *Migrator) CreateFile(name string) error {
	var box = packr.New("api-migrations", m.Opts.MigrationPath)

	fmt.Println(m.Opts.MigrationPath)
	if _, err := os.Stat(box.Path); os.IsNotExist(err) {
		return err
	}

	var templateContent = `-- +migrate Up
-- +migrate Down
`
	var tpl *template.Template = template.Must(template.New("new_migration").Parse(templateContent))
	fileName := fmt.Sprintf("%s-%s.sql", time.Now().Format("20060102150405"), strings.TrimSpace(name))
	pathName := path.Join(box.Path, fileName)
	f, err := os.Create(pathName)

	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if err := tpl.Execute(f, nil); err != nil {
		return err
	}

	return nil
}

var migrationDir string

func LoadConfigFromEnv(c *MigratorOpts) error {
	if os.Getenv("ENV") != "production" {
		if err := gotdotenv.Load(); err != nil {
			return err
		}
	}

	if err := env.Parse(c); err != nil {
		return err
	}

	return nil
}
