# Phase 0 基盤構築(walking skeleton)build spec

> 日付: 2026-06-21
> ステータス: **M1–M8 完了**(CI green)。以後の正典はコード。本書は build plan の記録。
> 前提(決定済み): [`../architecture/overview.md`](../architecture/overview.md) と [`../adr/`](../adr/README.md)。本書は設計の再決定ではなく **build plan**(何を・どの順で作り、何を満たせば完了か)。

## ゴール

「実機能は無いが、Phase 1 以降を載せられる土台が動いてデプロイできる」状態(= walking skeleton)。具体的には health API + 空 SPA + マイグレーション 1 本 + 疎通 1 経路 を Railway にデプロイできること。

## 非ゴール

- 取材・タイムライン・交流などの実機能(Phase 1+)。
- Cloud Tasks / Scheduler / GCS の配線(Phase 1。Phase 0 では worker に Connect スタブを置くのみ)。
- 認証の実装(Phase 2。ここでは `AuthProvider` 抽象の枠のみ)。

## マイルストーン

各 M は小さく、TDD・1 ファイル 200 行以内・完了条件つきで進める。

### M1 toolchain & 骨格
- `mise.toml`(go / node / task / sqlc / goose / buf / golangci-lint / lefthook をピン。pnpm は Corepack で pin)
- `Taskfile.yml`(dev / test / lint / gen / migrate / db:* の枠)
- `pnpm-workspace.yaml`(apps/* packages/*)、`.editorconfig`
- 完了: `mise install` でツールが揃い `task --list` が動く

### M2 ローカル DB
- `deploy/docker/docker-compose.yml`(Postgres)
- 完了: `task db:up` で起動・接続可

### M3 backend 最小(api)
- `backend/go.mod`、`cmd/api`(chi + `/healthz`)、`platform/{config,db,telemetry}`
- TDD: healthz ハンドラ / config ロード
- 完了: `task test` green・`/healthz` が 200

### M4 migration + sqlc
- `backend/db/migrations`(goose 初回)、`backend/db/queries`、`sqlc.yaml`
- 完了: `task migrate` 適用・`sqlc generate` 通過・生成物コミット

### M5 frontend 最小(web)
- `apps/web`(Vite + React + TS + Tailwind + TanStack Router / Query)、`/healthz` を表示
- 完了: `task web:dev` 起動・`task web:build` 通過

### M6 ハーネス / CI
- golangci-lint(+ depguard 境界ルール)、biome、lefthook、`.github/workflows/ci.yml`
- 完了: CI が lint / typecheck / test / codegen ドリフト を green

### M7 API 契約 + 内部 Connect(最小)
- `api/openapi.yaml`(最小)+ oapi-codegen + openapi-typescript、`proto` + buf(最小)
- 完了: 生成通過・web が生成クライアントで `/healthz` 疎通

### M8 デプロイ(Railway)
- `deploy/docker/Dockerfile.{api,worker}`、Railway 設定。GCP 配線は最小(後続)
- 完了: Railway で `/healthz` 到達

## 順序と依存

M1 → M2 → M3 →(M4・M5 は並行可)→ M6 → M7 → M8。worker(`cmd/worker`)は M8 で作成(HTTP `/healthz` + Connect スタブ)。Cloud Tasks 配線は Phase 1。

## 横断方針(ポインタ)

TDD / 200 行 / depguard 境界 / codegen ドリフト検出 は [`../architecture/cross-cutting.md`](../architecture/cross-cutting.md)。

## リスク / 留意

- Go・各 CLI は未インストール。mise で provision(M1)。
- 今は worker(Railway)↔ GCP の越境あり(全 GCP 化で解消、[ADR-0001](../adr/0001-hybrid-railway-gcp.md))。

## Phase 0 完了条件

M1〜M8 の各完了条件を満たし、Railway に骨格がデプロイされて `/healthz` が応答する。
