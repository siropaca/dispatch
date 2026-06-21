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
   - Deploy → Healthcheck Path: `/healthz`
   - **Networking → Generate Domain** で公開 URL を発行する
4. **worker**: 同様に Dockerfile = `deploy/docker/Dockerfile.worker`、`APP_ENV=production`、Healthcheck `/healthz`。Phase 0 では公開ドメイン不要(Cloud Tasks 配線 = Phase 1 で、OIDC 付きの公開エンドポイントにする)。
5. デプロイ後、api の公開 URL + `/healthz` が `{"status":"ok"}` を返すことを確認する。

## ビルド前提

- Dockerfile のビルドコンテキストはリポジトリ **root**(`backend/` を参照するため)。
- `PORT` は Railway が注入し、app(api / worker)はそれを読む。

## マイグレーション

- 現状は **自動実行されない**(api は起動時に DB へ Ping するだけ。空 DB でも `/healthz` は通る)。
- **Phase 1 で自動化する**: migration を `//go:embed` でバイナリに同梱し、Railway の **Pre-Deploy Command** で `goose up` を実行する(アプリ起動とは分離、多レプリカでも安全・冪等)。将来の Cloud Run では Cloud Run Job として同じバイナリを実行。

## 将来(全 GCP)

本格運用時は Cloud Run(api / worker)+ Cloud SQL + Secret Manager + Artifact Registry に移行する([ADR-0001](../../docs/adr/0001-hybrid-railway-gcp.md))。Dockerfile はそのまま Cloud Run でも使える。
