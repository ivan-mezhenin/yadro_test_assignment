package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"dns-manager/proto/dns"
	"dns-manager/server/config"
	"dns-manager/server/internal/adapter/resolver"
	"dns-manager/server/internal/controller"
	"dns-manager/server/internal/usecase"

	"google.golang.org/grpc"
)

func main() {
	var (
		configPath string
		resolvPath string
		backupPath string
	)
	flag.StringVar(&configPath, "config", "server/config/config.yaml", "config file path")
	flag.StringVar(&resolvPath, "resolv", "", "resolv.conf path (overrides config)")
	flag.StringVar(&backupPath, "backup", "", "backup path (overrides config)")
	flag.Parse()

	cfg := config.MustLoad(configPath)
	if resolvPath != "" {
		cfg.Server.ResolvConf = resolvPath
	}
	if backupPath != "" {
		cfg.Server.BackupPath = backupPath
	}

	logLevel := slog.LevelInfo
	if cfg.Logger.Debug {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})))

	if err := backupFile(cfg.Server.ResolvConf, cfg.Server.BackupPath); err != nil {
		slog.Warn("failed to create backup", "error", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
	if err != nil {
		slog.Error("failed to listen", "port", cfg.GRPC.Port, "error", err)
		os.Exit(1)
	}

	res := resolver.New(cfg.Server.ResolvConf)
	uc := usecase.NewDnsUseCase(res)
	h := controller.NewHandler(uc)

	s := grpc.NewServer()
	dns.RegisterDnsServiceServer(s, h)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		slog.Info("shutting down gracefully...")
		s.GracefulStop()
	}()

	slog.Info("starting gRPC server", "port", cfg.GRPC.Port)
	if err := s.Serve(lis); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func backupFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
