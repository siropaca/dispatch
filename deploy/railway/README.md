# Railway デプロイ手順

Railway に **4 サービス**を作る(同一リポジトリを参照)。公開は **web のみ**で、`/` は SPA・`/api/*` は api へ振り分ける単一オリジン構成([ADR-0015](../../docs/adr/0015-public-topology-edge-proxy.md))。詳細設計は [`../../docs/architecture/infrastructure.md`](../../docs/architecture/infrastructure.md)。

| サービス | 種別 | ビルド | 主な env |
|---|---|---|---|
| `postgres` | Railway Postgres(プラグイン) | — | 自動で接続情報を発行 |
| `web` | Dockerfile | `deploy/docker/Dockerfile.web` | `API_UPSTREAM`(api のプライベート宛先)、`PORT`(Railway 自動)。**公開ドメインはこれに付与** |
| `api` | Dockerfile | `deploy/docker/Dockerfile.api` | `APP_ENV=production`、`DATABASE_URL`(postgres から参照)、`PORT`(Railway 自動)。**プライベート(公開ドメインなし)** |
| `worker` | Dockerfile | `deploy/docker/Dockerfile.worker` | `APP_ENV=production`、`PORT`(Railway 自動) |

> 外部サービス(Clerk / OpenAI / Tavily)と GCP(Cloud Tasks / Scheduler / GCS)の env は、該当機能を実装する Phase で追加する。

## 手順

1. Railway でプロジェクトを作成し、GitHub リポジトリ `siropaca/dispatch` を連携する。
2. **postgres**: 「New → Database → PostgreSQL」を追加。発行された接続文字列を控える。
3. **api(プライベート)**: 「New → GitHub Repo」→ サービス設定で
   - Build: Dockerfile、Path = `deploy/docker/Dockerfile.api`、Root = リポジトリ root
   - Variables: `APP_ENV=production`、`DATABASE_URL=${{ Postgres.DATABASE_URL }}`(参照変数)
   - Deploy → Healthcheck Path: `/api/healthz`
   - **公開ドメインは発行しない**(web からのみ private network で到達する)
4. **web(公開エントリ)**: Dockerfile = `deploy/docker/Dockerfile.web`、Root = リポジトリ root
   - Variables: `API_UPSTREAM` = api のプライベート宛先。Railway の参照変数で `${{ api.RAILWAY_PRIVATE_DOMAIN }}:${{ api.PORT }}` のように指定(実変数名はダッシュボードで確認)
   - Deploy → Healthcheck Path: `/api/healthz`(Caddy 経由で api を確認)
   - **Networking → Custom Domain** で独自ドメインを付与する(公開 URL をこれで固定。[ADR-0015](../../docs/adr/0015-public-topology-edge-proxy.md))
5. **worker**: Dockerfile = `deploy/docker/Dockerfile.worker`、`APP_ENV=production`、Healthcheck `/healthz`(worker はルート維持)。公開ドメイン不要(Cloud Tasks 配線 = Phase 1 で OIDC 付きの受け口にする)。
6. デプロイ後、独自ドメインで `/`(SPA)と `/api/healthz`(`{"status":"ok"}`)が応答することを確認する。

## ビルド前提

- Dockerfile のビルドコンテキストはリポジトリ **root**(api/worker は `backend/`、web は JS workspace を参照するため)。`.dockerignore` で `node_modules` 等を除外。
- `PORT` は Railway が注入する(api/worker は app が読み、web は Caddy が `:{$PORT}` で待ち受ける)。

## マイグレーション

- 現状は **自動実行されない**(api は起動時に DB へ Ping するだけ。空 DB でも `/healthz` は通る)。
- **Phase 1 で自動化する**: migration を `//go:embed` でバイナリに同梱し、Railway の **Pre-Deploy Command** で `goose up` を実行する(アプリ起動とは分離、多レプリカでも安全・冪等)。将来の Cloud Run では Cloud Run Job として同じバイナリを実行。

## 将来(全 GCP)

本格運用時は Cloud Run(api / worker)+ Cloud SQL + Secret Manager + Artifact Registry に移行する([ADR-0001](../../docs/adr/0001-hybrid-railway-gcp.md))。Dockerfile はそのまま Cloud Run でも使える。
