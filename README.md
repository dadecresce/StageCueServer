# StageCueServer v0.2

This release introduces:

* **TOML configuration** via `internal/config`
* **Structured logging** using `zap`
* Graceful shutdown with SIGINT/SIGTERM

## Quickstart

```bash
go mod tidy
go build -o bin/StageCueServer ./cmd/server
./bin/StageCueServer -config config.sample.toml
```

Visit [http://localhost:8080/healthz](http://localhost:8080/healthz) ⇒ `ok`
