# backend — Dispatch バックエンド(Go)

モジュラモノリス + DDD + ヘキサゴナル(Ports & Adapters)の Go モジュール。詳細設計は [`../docs/architecture/backend.md`](../docs/architecture/backend.md)、全体は [`../docs/architecture/overview.md`](../docs/architecture/overview.md)。

## エントリポイント

- `cmd/api` — 公開 REST エッジ(chi、`/healthz`、Clerk JWT 検証は今後)。MVP では SPA も配信。
- `cmd/worker` — Cloud Tasks の push を受ける HTTP ハンドラ(内部 Connect)。取材パイプラインは Phase 1。

両者は `internal/` を共有する。

## 構成

```
cmd/{api,worker}/
internal/
  {identity,newsroom,reporting,publishing,timeline,interaction}/  # context は各 Phase で追加
    domain/ app/ adapters/
  platform/   # config, db(+sqlcgen), httpserver, httpapi(oapi 生成), telemetry
  proto/      # buf 生成の内部 Connect(dispatch.reporting.v1 …)
db/{migrations(goose), queries(sqlc)}/  sqlc.yaml
.golangci.yml
```

**依存ルール**: `domain ← app ← adapters`。`domain` は infra を import しない。context 越境 import は禁止(depguard で強制)。

## 開発

リポジトリ root の Taskfile 経由を推奨(toolchain は mise 管理)。

```sh
task dev          # api / worker / web を同時起動
task api:dev      # api のみ
task worker:dev   # worker のみ
task test         # go test ./...
task lint         # golangci-lint(+ biome)
task migrate      # goose up(要 DATABASE_URL)
task gen          # codegen(sqlc / oapi-codegen / openapi-typescript / buf)
```

直接実行は `mise exec -- go test ./...` のように mise 経由で。

## 環境変数

`backend/.env.example` 参照(`task` 実行時に `backend/.env` が自動ロードされる)。`api` は `DATABASE_URL` 必須、`PORT`(既定 8080)/ `APP_ENV` を読む。`worker` は `PORT`(既定 8081)/ `APP_ENV`。

## codegen(生成物はコミット必須)

| 生成物 | 入力 | 出力 |
|---|---|---|
| sqlc | `db/queries` + `db/migrations` | `internal/platform/db/sqlcgen` |
| oapi-codegen | `../api/openapi.yaml` | `internal/platform/httpapi` |
| buf(proto / connect) | `../proto` | `internal/proto` |

CI は `task gen` 後に差分が出れば fail(ドリフト検出)。

## テスト

ドメインは table-driven の単体、結合は testcontainers-go(本物 Postgres、Phase 1 で導入)。AI / Queue / BlobStore 等の port は fake で決定的に。
