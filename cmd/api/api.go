package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/blankdots/minimal-kube-app/internal/config"
	"github.com/blankdots/minimal-kube-app/internal/database"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var AppConfig *config.Config

var db *database.Datastore

var token string

type ErrorResponse struct {
	Message interface{} `json:"message"`
}

func main() {

	AppConfig, err := config.App("api")
	if err != nil {
		log.Fatal(err)
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err = database.NewDatabase(ctx, AppConfig.Database)
	if err != nil {
		log.Panicf("database connection failed, reason: %v", err)
	}

	token = AppConfig.API.StaticToken

	router := gin.New()

	// adding this not to overcrowd logs with /health endpoint requests
	router.Use(
		gin.LoggerWithWriter(gin.DefaultWriter, "/health"),
		gin.Recovery(),
	)

	router.HandleMethodNotAllowed = true

	// Health has no auth; protected routes use auth middleware
	router.GET("/health", healthResponse)

	protected := router.Group("/")
	protected.Use(authMiddleware())
	protected.GET("/query", apiResponse)

	srv := &http.Server{
		Addr:              AppConfig.API.Host + ":" + fmt.Sprint(AppConfig.API.Port),
		Handler:           router,
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       5 * time.Minute,
		WriteTimeout:      20 * time.Second,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	shutdown()
	log.Debug("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Debug("Server forced to shutdown: ", err)
	}

	log.Info("Server exiting")
}

func shutdown() {
	defer db.Close()
}

func extractToken(c *gin.Context) (string, error) {
	// Authorization: Bearer <token>
	if auth := c.GetHeader("Authorization"); auth != "" {
		const prefix = "Bearer "
		if len(auth) >= len(prefix) && strings.EqualFold(auth[:len(prefix)], prefix) {
			return strings.TrimSpace(auth[len(prefix):]), nil
		}
		return "", errors.New("expected Authorization: Bearer <token>")
	}
	// X-API-Key: <token>
	if key := c.GetHeader("X-API-Key"); key != "" {
		return strings.TrimSpace(key), nil
	}
	return "", errors.New("missing auth: provide Authorization: Bearer <token> or X-API-Key header")
}

const authRealm = "api"

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken, err := extractToken(c)
		if err != nil {
			c.Header("WWW-Authenticate", `Bearer realm="`+authRealm+`"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: err.Error(),
			})
			return
		}
		if clientToken != token {
			log.Debugf("auth failed: token mismatch")
			c.Header("WWW-Authenticate", `Bearer realm="`+authRealm+`"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Message: "invalid token",
			})
			return
		}
		c.Next()
	}
}

func apiResponse(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	packageName := c.Query("package")
	data, _ := database.SelectData(db, packageName)
	c.JSON(200, data)
}

// HealthResponse
func healthResponse(c *gin.Context) {
	// ok response to health

	c.Writer.WriteHeader(http.StatusOK)
}
