package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"strings"
	"syscall"

	"github.com/polymarket/common/pkg/liveconfig"
)

type DemoRuntimeConfig struct {
	ExecutionStreamRatio int  `json:"execution_stream_ratio"`
	FeatureEnabled       bool `json:"feature_enabled"`
	MaxWorkers           int  `json:"max_workers"`
}

func (c DemoRuntimeConfig) String() string {
	return fmt.Sprintf(
		"execution_stream_ratio=%d feature_enabled=%t max_workers=%d",
		c.ExecutionStreamRatio,
		c.FeatureEnabled,
		c.MaxWorkers,
	)
}

func main() {
	cfg := loadAppConfig()
	store := liveconfig.NewAtomicValueStore(defaultDemoRuntimeConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	syncer, err := liveconfig.NewJSONWatchSyncer(
		liveconfig.EtcdConfig{
			Endpoints:   cfg.EtcdEndpoints,
			DialTimeout: cfg.EtcdDialTimeout,
			Key:         cfg.EtcdConfigKey,
		},
		store,
		validateDemoRuntimeConfig,
		func(next DemoRuntimeConfig) {
			log.Printf("[updated] %s", next.String())
		},
	)
	if err != nil {
		log.Fatalf("unable to initialize live config sync: %v", err)
	}
	defer func() {
		if err := syncer.Close(); err != nil {
			log.Printf("close syncer error: %v", err)
		}
	}()

	syncer.Start(ctx)
	log.Printf("watching etcd key=%s endpoints=%s", cfg.EtcdConfigKey, strings.Join(cfg.EtcdEndpoints, ","))
	log.Printf("initial config: %s", store.Get().String())

	ticker := newReadTicker(cfg)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down")
			return
		case <-ticker.C:
			cfg := store.Get()
			log.Printf("[read] %s", cfg.String())
		}
	}
}

func validateDemoRuntimeConfig(cfg DemoRuntimeConfig) error {
	if cfg.ExecutionStreamRatio < 0 || cfg.ExecutionStreamRatio > 100 {
		return fmt.Errorf("execution_stream_ratio must be between 0 and 100")
	}
	if cfg.MaxWorkers <= 0 {
		return fmt.Errorf("max_workers must be > 0")
	}
	return nil
}
