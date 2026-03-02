package database

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

// DBConfig struct for Postgresql connection
type DBConfig struct {
	URL string
}

// Datastore struct that acts as a receiver for the DB update methods
type Datastore struct {
	pool *pgxpool.Pool
}

// DependencyInfo holds package dependency data from the npm Registry API.
type DependencyInfo struct {
	PackageName  string            `json:"packageName"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
	UpdatedAt    string            `json:"updatedAt"`
}

func NewDatabase(ctx context.Context, config DBConfig) (*Datastore, error) {

	dbs := Datastore{}

	if config.URL == "" {
		log.Fatal("Database URL is empty")
	}

	poolConfig, err := pgxpool.ParseConfig(config.URL)
	if err != nil {
		log.Fatalln("Unable to parse database URL:", err)
	}

	poolConfig.MinConns = 10
	poolConfig.MaxConns = 10

	dbs.pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	log.Debugf("Connected to Database %s", config.URL)

	return &dbs, nil
}

// InsertData stores package dependency info from the cronjob (npm Registry API).
func InsertData(db *Datastore, packageName, version, depsJSON, updatedAt string) error {

	if _, err := db.pool.Exec(
		context.Background(),
		`INSERT INTO package_dep(packageName, version, dependencies, updatedAt) VALUES ($1, $2, $3, $4)
		 ON CONFLICT (packageName) DO UPDATE SET version=$2, dependencies=$3, updatedAt=$4`,
		packageName, version, depsJSON, updatedAt); err != nil {
		log.Errorf("Exec failed: %v\n", err)

		return err
	}

	return nil

}

// SelectData returns dependency info for a package (REST API).
func SelectData(db *Datastore, packageName string) (*DependencyInfo, error) {

	info := DependencyInfo{}
	var depsJSON string
	err := db.pool.QueryRow(
		context.Background(),
		"SELECT packageName, version, dependencies, updatedAt FROM package_dep WHERE packageName=$1",
		packageName).Scan(&info.PackageName, &info.Version, &depsJSON, &info.UpdatedAt)
	if err != nil {
		log.Errorf("QueryRow failed: %v\n", err)

		return nil, err
	}

	if depsJSON != "" {
		_ = json.Unmarshal([]byte(depsJSON), &info.Dependencies)
	}

	return &info, nil
}

// Close closes the underying PostgreSQL database.
func (d *Datastore) Close() error {
	if d.pool != nil {
		d.pool.Close()
	}

	return nil
}
