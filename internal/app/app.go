package app

import (
	"context"
	"github.com/NikitaVi/microservices_kafka/internal/config"
	"github.com/NikitaVi/platform_shared/pkg/closer"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := a.serviceProvider.NoteSaverConsumer(ctx).RunConsumer(ctx)
		if err != nil {
			log.Printf("failed to run consumer: %v", err)
		}
	}()

	gracefulShutdown(ctx, cancel, wg)
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(ctx context.Context) error{
		a.initConfig,
		a.initServiceProvider,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.serviceProvider = NewServiceProvider()
	return nil
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		log.Printf("graceful shutdown cancelled")
	case <-waitSignal():
		log.Printf("terminate via signal")
	}
	cancel()
	if wg != nil {
		wg.Wait()
	}
}

func waitSignal() <-chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	return sigterm
}
