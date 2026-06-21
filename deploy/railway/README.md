# Railway デプロイ手順

Railway に **3 サービス**を作る(同一リポジトリを参照)。詳細設計は [`../../docs/architecture/infrastructure.md`](../../docs/architecture/infrastructure.md)。

| サービス | 種別 | ビルド | 主な env |
|---|---|---|---|
| `postgres` | Railway Postgres(プラグイン) | — | 自動で接続情報を発行 |
| `api` | Dockerfile | `deploy/docker/Dockerfile.api` | `APP_ENV=production`、`DATABASE_URL`(postgres から参照)、`PORT`(Railway 自動) |
| `worker` | Dockerfile | `deploy/docker/Dockerfile.worker` | `APP_ENV=production`、`PORT`(Railway 自動) |

> 外部サービス(Clerk / OpenAI / Tavily)と GCP(Cloud Tasks / Scheduler / GCS)の env は、該当機能を実装する Phase で追加する。

## 手順

1. Railway でプロジェクトを作成し、GitHub リポジトリ `siropaca/dispatch` を連携する。
2. **postgres**: 「New → Database → PostgreSQL」を追加。発行された接続文字列を控える。
3. **api**: 「New → GitHub Repo」→ サービス設定で
   - Build: Dockerfile、Path = `deploy/docker/Dockerfile.api`、Root = リポジトリ root
   - Variables: `APP_ENV=production`、`DATABASE_URL=${{ Postgres.DATABASE_URL }}`(参照変数)
   - Healthcheck Path: `/healthz`
4. **worker**: 同様に Dockerfile = `deploy/docker/Dockerfile.worker`、`APP_ENV=production`、Healthcheck `/healthz`。
5. デプロイ後、api の公開 URL で `/healthz` が `{"status":"ok"}` を返すことを確認する。

## ビルド前提

- Dockerfile のビルドコンテキストはリポジトリ **root**(`backend/` を参照するため)。
- `PORT` は Railway が注入し、app(api / worker)はそれを読む。

## 将来(全 GCP)

本格運用時は Cloud Run(api / worker)+ Cloud SQL + Secret Manager + Artifact Registry に移行する([ADR-0001](../../docs/adr/0001-hybrid-railway-gcp.md))。Dockerfile はそのまま Cloud Run でも使える。
