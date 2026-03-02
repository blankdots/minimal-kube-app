package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/blankdots/minimal-kube-app/internal/config"
	"github.com/blankdots/minimal-kube-app/internal/database"
	log "github.com/sirupsen/logrus"
)

var AppConfig *config.Config

var db *database.Datastore

type harvester struct {
	apiBase  string
	packages []string
}

// NPMRegistryRecord matches npm Registry API /package/latest response (no auth).
type NPMRegistryRecord struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Dependencies map[string]string `json:"dependencies"`
	Time         map[string]string `json:"time"` // optional, has "modified" etc
}

func newHarvester() *harvester {

	AppConfig, err := config.App("cronjob")
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err = database.NewDatabase(ctx, AppConfig.Database)
	if err != nil {
		log.Panicf("database connection failed, reason: %v", err)
	}

	rv := &harvester{
		apiBase:  AppConfig.CronJob.APIBase,
		packages: AppConfig.CronJob.Packages,
	}

	harvest(rv)

	return rv
}

func harvest(h *harvester) {
	for _, pkg := range h.packages {
		url := fmt.Sprintf("%s/%s/latest", h.apiBase, pkg)
		log.Debugf("Fetching %s", url)

		response, err := http.Get(url) //nolint:gosec // configurable URL
		if err != nil {
			log.Errorf("error on API call to %s: %v", url, err)
			continue
		}

		responseData, err := io.ReadAll(response.Body)
		response.Body.Close() //nolint:errcheck // defer to io.ReadAll error handling
		if err != nil {
			log.Errorf("error reading response data: %v", err)
			continue
		}

		var rec NPMRegistryRecord
		if err := json.Unmarshal(responseData, &rec); err != nil {
			log.Errorf("error unmarshal for %s: %v", pkg, err)
			continue
		}

		depsJSON, _ := json.Marshal(rec.Dependencies)
		if rec.Dependencies == nil {
			depsJSON = []byte("{}")
		}

		updatedAt := time.Now().UTC().Format(time.RFC3339)
		if t, ok := rec.Time["modified"]; ok && t != "" {
			updatedAt = t
		}

		log.Debugf("Inserting package: %s@%s (%d deps)", rec.Name, rec.Version, len(rec.Dependencies))
		if err := database.InsertData(db, rec.Name, rec.Version, string(depsJSON), updatedAt); err != nil {
			log.Errorf("error inserting %s: %v", pkg, err)
		}
	}
}

func main() {
	newHarvester()
}
