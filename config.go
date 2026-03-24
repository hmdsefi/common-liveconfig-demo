package main

import (
	"log"
	"os"
	"strings"
	"time"
)

type AppConfig struct {
	EtcdEndpoints   []string
	EtcdConfigKey   string
	EtcdDialTimeout time.Duration
	PrintEvery      time.Duration
}

func loadAppConfig() AppConfig {
	return AppConfig{
		EtcdEndpoints:   parseListEnv("ETCD_ENDPOINTS", []string{"localhost:2379"}),
		EtcdConfigKey:   parseStringEnv("ETCD_CONFIG_KEY", "/demo/live-config"),
		EtcdDialTimeout: parseDurationEnv("ETCD_DIAL_TIMEOUT", 3*time.Second),
		PrintEvery:      parseDurationEnv("PRINT_EVERY", 5*time.Second),
	}
}

func defaultDemoRuntimeConfig() DemoRuntimeConfig {
	return DemoRuntimeConfig{
		ExecutionStreamRatio: 0,
		FeatureEnabled:       false,
		MaxWorkers:           2,
	}
}

func newReadTicker(cfg AppConfig) *time.Ticker {
	return time.NewTicker(cfg.PrintEvery)
}

func parseListEnv(key string, fallback []string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	parts := strings.Split(raw, ",")
	res := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			res = append(res, v)
		}
	}
	if len(res) == 0 {
		return fallback
	}
	return res
}

func parseStringEnv(key, fallback string) string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	return raw
}

func parseDurationEnv(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		log.Printf("invalid duration for %s=%q, using fallback=%s", key, raw, fallback)
		return fallback
	}
	return d
}
