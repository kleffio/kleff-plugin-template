// Command plugin is the gRPC server for this Kleff plugin.
// It declares capabilities via GetCapabilities and exposes a UI manifest so the
// panel can surface nav items and pages contributed by this plugin.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	pluginsv1 "github.com/kleffio/plugin-sdk-go/v1"
	"google.golang.org/grpc"
)

type pluginServer struct {
	pluginsv1.UnimplementedPluginHealthServer
	pluginsv1.UnimplementedPluginUIServer
}

func (s *pluginServer) Health(_ context.Context, _ *pluginsv1.HealthRequest) (*pluginsv1.HealthResponse, error) {
	return &pluginsv1.HealthResponse{Status: pluginsv1.HealthStatusHealthy, Message: "my-plugin 0.1.0"}, nil
}

func (s *pluginServer) GetCapabilities(_ context.Context, _ *pluginsv1.GetCapabilitiesRequest) (*pluginsv1.GetCapabilitiesResponse, error) {
	return &pluginsv1.GetCapabilitiesResponse{
		Capabilities: []string{pluginsv1.CapabilityUIManifest},
	}, nil
}

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

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	port := env("PORT", "50051")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Error("listen failed", "error", err)
		os.Exit(1)
	}

	srv := grpc.NewServer()
	p := &pluginServer{}
	pluginsv1.RegisterPluginHealthServer(srv, p)
	pluginsv1.RegisterPluginUIServer(srv, p)

	go func() {
		logger.Info("plugin listening", "port", port)
		if serveErr := srv.Serve(lis); serveErr != nil {
			logger.Error("gRPC server error", "error", serveErr)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	logger.Info("shutting down")
	srv.GracefulStop()
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
