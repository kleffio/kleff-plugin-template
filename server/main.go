// Command plugin is the gRPC server for this Kleff plugin.
// It declares a nav item and settings page via GetUIManifest so the
// panel can surface them, and reports healthy to the platform health check.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	pluginsv1 "github.com/kleffio/plugin-sdk/v1"
	"google.golang.org/grpc"
)

// ── Plugin server ─────────────────────────────────────────────────────────────

type pluginServer struct {
	pluginsv1.UnimplementedPluginHealthServer
	pluginsv1.UnimplementedPluginUIServer
}

func (s *pluginServer) Health(_ context.Context, _ *pluginsv1.HealthRequest) (*pluginsv1.HealthResponse, error) {
	return &pluginsv1.HealthResponse{Status: pluginsv1.HealthStatusHealthy}, nil
}

func (s *pluginServer) GetCapabilities(_ context.Context, _ *pluginsv1.GetCapabilitiesRequest) (*pluginsv1.GetCapabilitiesResponse, error) {
	return &pluginsv1.GetCapabilitiesResponse{
		Capabilities: []string{pluginsv1.CapabilityUIManifest},
	}, nil
}

// GetUIManifest tells the panel which nav items and settings pages this plugin
// contributes. Keep this in sync with what you declare in src/index.tsx.
func (s *pluginServer) GetUIManifest(_ context.Context, _ *pluginsv1.GetUIManifestRequest) (*pluginsv1.GetUIManifestResponse, error) {
	pluginID := env("PLUGIN_ID", "my-plugin")
	return &pluginsv1.GetUIManifestResponse{
		Manifest: &pluginsv1.UIManifest{
			PluginID: pluginID,
			NavItems: []*pluginsv1.NavItem{
				{
					Label: "My Plugin",
					Icon:  "Puzzle",
					Path:  fmt.Sprintf("/p/%s/example", pluginID),
				},
			},
		},
	}, nil
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	srv := &pluginServer{}

	gs := grpc.NewServer()
	pluginsv1.RegisterPluginHealthServer(gs, srv)
	pluginsv1.RegisterPluginUIServer(gs, srv)

	port := env("PLUGIN_PORT", "50051")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("listen failed", "error", err)
		os.Exit(1)
	}

	go func() {
		logger.Info("plugin listening", "port", port)
		if err := gs.Serve(lis); err != nil {
			logger.Error("gRPC server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	logger.Info("shutting down")
	gs.GracefulStop()
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
