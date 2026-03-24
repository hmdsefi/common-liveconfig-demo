# common-liveconfig-demo

Minimal demo app outside `clob` that uses `github.com/polymarket/common/pkg/liveconfig`
to pull non-sensitive runtime config from etcd and update it live.

## What it demonstrates

- Shared `common` package usage from outside `clob`
- Initial fetch from etcd key
- Live watch updates on the same key
- Atomic, race-safe reads via `AtomicValueStore`
- Validation before applying updates

## Config schema in etcd

The etcd value must be JSON:

```json
{
  "execution_stream_ratio": 25,
  "feature_enabled": true,
  "max_workers": 4
}
```

## Run

From this folder:

```bash
docker compose up -d
go mod tidy
go run .
```

This demo uses etcd **without username/password** (local non-auth setup).

Optional env vars:

- `ETCD_ENDPOINTS` (default: `localhost:2379`)
- `ETCD_CONFIG_KEY` (default: `/demo/live-config`)
- `ETCD_DIAL_TIMEOUT` (default: `3s`)
- `PRINT_EVERY` (default: `5s`)

Example:

```bash
ETCD_ENDPOINTS=localhost:2379 ETCD_CONFIG_KEY=/demo/live-config PRINT_EVERY=2s go run .
```

## Update the config live

In another terminal:

```bash
etcdctl put /demo/live-config '{"execution_stream_ratio":10,"feature_enabled":false,"max_workers":2}'
etcdctl put /demo/live-config '{"execution_stream_ratio":80,"feature_enabled":true,"max_workers":8}'
```

or
```bash
docker exec -it etcd sh -lc "ETCDCTL_API=3 etcdctl put /demo/live-config '{\"execution_stream_ratio\":10,\"feature_enabled\":false,\"max_workers\":2}'"
docker exec -it etcd sh -lc "ETCDCTL_API=3 etcdctl put /demo/live-config '{\"execution_stream_ratio\":80,\"feature_enabled\":true,\"max_workers\":8}'"
```

The running app will log `[updated]` when etcd changes are applied and `[read]` on each poll interval.
