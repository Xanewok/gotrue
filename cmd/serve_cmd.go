package cmd

import (
	"context"
	"net"

	"github.com/netlify/gotrue/api"
	"github.com/netlify/gotrue/conf"
	"github.com/netlify/gotrue/observability"
	"github.com/netlify/gotrue/storage"
	"github.com/netlify/gotrue/utilities"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = cobra.Command{
	Use:  "serve",
	Long: "Start API server",
	Run: func(cmd *cobra.Command, args []string) {
		serve(cmd.Context())
	},
}

func serve(ctx context.Context) {
	config, err := conf.LoadGlobal(configFile)
	if err != nil {
		logrus.WithError(err).Fatal("unable to load config")
	}

	if err := observability.ConfigureLogging(&config.Logging); err != nil {
		logrus.WithError(err).Error("unable to configure logging")
	}

	if err := observability.ConfigureTracing(ctx, &config.Tracing); err != nil {
		logrus.WithError(err).Error("unable to configure tracing")
	}

	if err := observability.ConfigureMetrics(ctx, &config.Metrics); err != nil {
		logrus.WithError(err).Error("unable to configure metrics")
	}

	db, err := storage.Dial(config)
	if err != nil {
		logrus.WithError(err).Fatalf("error opening database: %+v", err)
	}
	defer db.Close()

	api := api.NewAPIWithVersion(ctx, config, db, utilities.Version)

	addr := net.JoinHostPort(config.API.Host, config.API.Port)
	logrus.Infof("GoTrue API started on: %s", addr)

	api.ListenAndServe(ctx, addr)
}
