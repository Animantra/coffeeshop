package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/thangchung/go-coffeeshop/cmd/counter/config"
	"github.com/thangchung/go-coffeeshop/internal/counter/app"
	"github.com/thangchung/go-coffeeshop/pkg/logger"
	"github.com/thangchung/go-coffeeshop/pkg/postgres"
	"github.com/thangchung/go-coffeeshop/pkg/rabbitmq"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc"

	pkgConsumer "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/consumer"
	pkgPublisher "github.com/thangchung/go-coffeeshop/pkg/rabbitmq/publisher"

	_ "github.com/lib/pq"
)

func main() {
	// Сеттим GOMAXPROCS для корректной работы в контейнерах
	_, err := maxprocs.Set()
	if err != nil {
		slog.Error("failed set max procs", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed get config", err)
	}

	slog.Info("⚡ init app", "name", cfg.Name, "version", cfg.Version)

	// Настройка logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logger.ConvertLogLevel(cfg.Log.Level))

	// Интеграция Logrus с slog
	slog.New(logger.NewLogrusHandler(logrus.StandardLogger()))

	server := grpc.NewServer()

	// Асинхронный перехват отмены контекста для мягкой остановки gRPC
	go func() {
		<-ctx.Done()
		slog.Info("🔄 Shutting down gRPC server gracefully...")
		server.GracefulStop()
	}()

	cleanup := prepareApp(ctx, cancel, cfg, server)

	// Настройка сетевого лисенера для gRPC
	address := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	network := "tcp"

	l, err := net.Listen(network, address)
	if err != nil {
		slog.Error("failed to listen to address", err, "network", network, "address", address)
		return
	}

	slog.Info("🌏 start server...", "address", address)

	defer func() {
		if err1 := l.Close(); err1 != nil {
			slog.Error("failed to close listener", err1, "network", network, "address", address)
		}
	}()

	// Сервер метрик и хелсчеков
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ok","service":"counter-service"}`))
		})
		if err := http.ListenAndServe(":9091", nil); err != nil {
			slog.Error("metrics server failed", err)
		}
	}()

	// Запускаем gRPC-сервер в горутине, чтобы он не блокировал main
	go func() {
		if err := server.Serve(l); err != nil && err != grpc.ErrServerStopped {
			slog.Error("failed start gRPC server", err, "network", network, "address", address)
			cancel()
		}
	}()

	// Канал для системных сигналов (Ctrl+C, SIGTERM от Kubernetes)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Точка блокировки: ждем либо сигнал от ОС, либо критическую отмену контекста
	select {
	case v := <-quit:
		slog.Info("🔴 Signal received, exiting...", "signal", v)
	case <-ctx.Done():
		slog.Info("⚠️ Context canceled, app is stopping...")
	}

	// Финальная очистка ресурсов (закрытие коннектов к БД и брокеру)
	cleanup()
	slog.Info("👋 Counter-service successfully stopped")
}

func prepareApp(ctx context.Context, cancel context.CancelFunc, cfg *config.Config, server *grpc.Server) func() {
	a, cleanup, err := app.InitApp(cfg, postgres.DBConnString(cfg.PG.DsnURL), rabbitmq.RabbitMQConnStr(cfg.RabbitMQ.URL), server)
	if err != nil {
		slog.Error("failed init app", err)
		cancel()
		return func() {}
	}

	a.BaristaOrderPub.Configure(
		pkgPublisher.ExchangeName("barista-order-exchange"),
		pkgPublisher.BindingKey("barista-order-routing-key"),
		pkgPublisher.MessageTypeName("barista-order-created"),
	)

	a.KitchenOrderPub.Configure(
		pkgPublisher.ExchangeName("kitchen-order-exchange"),
		pkgPublisher.BindingKey("kitchen-order-routing-key"),
		pkgPublisher.MessageTypeName("kitchen-order-created"),
	)

	a.Consumer.Configure(
		pkgConsumer.ExchangeName("counter-order-exchange"),
		pkgConsumer.QueueName("counter-order-queue"),
		pkgConsumer.BindingKey("counter-order-routing-key"),
		pkgConsumer.ConsumerTag("counter-order-consumer"),
	)

	// Безопасный запуск консьюмера без дедлока главного потока
	go func() {
		err1 := a.Consumer.StartConsumer(a.Worker)
		if err1 != nil {
			slog.Error("failed to start Consumer", err1)
			cancel()
		}
	}()

	return cleanup
}
