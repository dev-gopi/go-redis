# go-redis (Redis clone)

A small Redis-like server written in Go for learning and experimentation.

## What it is

- Implements a subset of Redis commands (strings, hashes, JSON root values, basic scanning, TTL, memory estimates).
- Supports multiple logical databases (SELECT).
- Persistence via snapshot (`data/dump.rdb`), AOF (`data/appendonly.aof`), and a WAL (`data/wal.log`).
- Simple auth support (disabled by default).

This is an educational project — not production-ready. The AOF/WAL formats are line-based and have limitations.

## Build

Requires Go 1.20+.

```bash
go build ./...
```

## Run

Start the server (default listens on `:6379`):

```bash
go run cmd/server/main.go
```

The server automatically:
- Loads `data/dump.rdb` on startup (snapshot)
- Replays `data/appendonly.aof`
- Initializes `data/wal.log`
- Starts a background snapshot saver (every minute)
- Starts AOF auto-rotate at ~10MB

## Client

Use a Redis-compatible client (recommended):

```bash
redis-cli -p 6379
```

Examples:

```text
PING
SET mykey hello
GET mykey
HSET myhash field1 value1
HGET myhash field1
JSON.SET myjson . {"test":"value"}
JSON.GET myjson .
SCAN 0 MATCH * COUNT 100
SELECT 1
```

## Persistence files

- `data/dump.rdb` — JSON snapshot of all DBs
- `data/appendonly.aof` — textual command append log
- `data/wal.log` — simple WAL entries

Notes:
- `JSON.SET` currently supports only the root path (`.`).
- The AOF/WAL format is simple text and may not properly escape all argument characters; prefer `redis-cli` to avoid framing issues.

## Authentication

Auth is available via the `AUTH` command and guarded by an internal manager (disabled by default). To enable, modify `internal/auth/auth.go` and set `Manager.Enabled = true` and appropriate credentials.

## Supported commands (selected)

- `PING`, `AUTH`, `INFO`, `CLIENT`
- Strings: `SET`, `GET`, `DEL`, `EXISTS`, `STRLEN`, `GETRANGE`
- Hashes: `HSET`, `HGET`, `HEXISTS`, `HDEL`, `HMGET`, `HLEN`, `HSCAN`
- JSON: `JSON.SET`, `JSON.GET`, `JSON.DEL` (root path only)
- Keys: `SCAN`, `TYPE`, `TTL`, `DBSIZE`

## Development notes

- RESP parsing is implemented in `internal/protocol/resp_parser.go`. Use `redis-cli` to avoid framing problems.
- Snapshot load now normalizes persisted types so hash values and JSON values are usable after restart.
- If you want to add commands or persistence formats, look under `internal/commands` and `internal/persistence`.

## License

This project is for learning; no license specified.

---
