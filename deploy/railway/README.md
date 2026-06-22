# Railway デプロイ手順

Railway に **4 サービス**を作る(同一リポジトリを参照)。公開は **web のみ**で、`/` は SPA・`/api/*` は api へ振り分ける単一オリジン構成([ADR-0015](../../docs/adr/0015-public-topology-edge-proxy.md))。詳細設計は [`../../docs/architecture/infrastructure.md`](../../docs/architecture/infrastructure.md)。

| サービス | 種別 | ビルド | 主な env |
|---|---|---|---|
| `postgres` | Railway Postgres(プラグイン) | — | 自動で接続情報を発行 |
| `web` | Dockerfile | `deploy/docker/Dockerfile.web` | `API_UPSTREAM`(api のプライベート宛先)、`PORT`(Railway 自動)。**公開ドメインはこれに付与** |
| `api` | Dockerfile | `deploy/docker/Dockerfile.api` | `APP_ENV=production`、`DATABASE_URL`(postgres から参照)、`PORT=8080`(**明示**。private 宛先のポート確定)。**プライベート(公開ドメインなし)** |
| `worker` | Dockerfile | `deploy/docker/Dockerfile.worker` | `APP_ENV=production`、`PORT`(Railway 自動) |

> 外部サービス(Clerk / OpenAI / Tavily)と GCP(Cloud Tasks / Scheduler / GCS)の env は、該当機能を実装する Phase で追加する。

## 手順

1. Railway でプロジェクトを作成し、GitHub リポジトリ `siropaca/dispatch` を連携する。
2. **postgres**: 「New → Database → PostgreSQL」を追加。発行された接続文字列を控える。
3. **api(プライベート)**: 「New → GitHub Repo」→ サービス設定で
   - Build: Dockerfile、Path = `deploy/docker/Dockerfile.api`、Root = リポジトリ root
   - Variables: `APP_ENV=production`、`DATABASE_URL=${{ Postgres.DATABASE_URL }}`(参照変数)、`PORT=8080`(**明示**。下記「private networking」参照)
   - Deploy → Healthcheck Path: `/api/healthz`
   - **公開ドメインは発行しない**(web からのみ private network で到達する)
4. **web(公開エントリ)**: Dockerfile = `deploy/docker/Dockerfile.web`、Root = リポジトリ root
   - Variables: `API_UPSTREAM = ${{ api.RAILWAY_PRIVATE_DOMAIN }}:8080`(api に `PORT=8080` を設定済みなら `${{ api.PORT }}` でも可)
   - Deploy → Healthcheck Path: `/api/healthz`(Caddy 経由で api を確認)
   - **Networking**: 独自ドメインがあれば **Custom Domain** で付与し公開 URL を固定する([ADR-0015](../../docs/adr/0015-public-topology-edge-proxy.md))。未取得の間は **Generate Domain** の Railway ドメインで運用し、取得後に差し替える(api/worker は非公開のまま)。
5. **worker**: Dockerfile = `deploy/docker/Dockerfile.worker`、`APP_ENV=production`、Healthcheck `/healthz`(worker はルート維持)。公開ドメイン不要(Cloud Tasks 配線 = Phase 1 で OIDC 付きの受け口にする)。
6. デプロイ後、web の公開 URL で `/`(SPA)と `/api/healthz`(`{"status":"ok"}`)が応答することを確認する。

## private networking(api 宛先)

- web(Caddy)→ api は Railway の private network 経由。**公開側のようなエッジプロキシが無いため、api が実際に listen するポートを直接指定する**。
- よって api には `PORT=8080` を**明示**し、`API_UPSTREAM` は `<api>.railway.internal:8080`(= `${{ api.RAILWAY_PRIVATE_DOMAIN }}:8080`)を指す。PORT を明示しないと `${{ api.PORT }}` が参照候補に出ない。
- Go の `:8080` は IPv6(デュアルスタック)で待ち受けるため private(IPv6)から到達できる。繋がらず Caddy が 502 を返す場合は、ポート一致・private domain のサービス名・api の listen を確認する。

## ビルド前提

- Dockerfile のビルドコンテキストはリポジトリ **root**(api/worker は `backend/`、web は JS workspace を参照するため)。`.dockerignore` で `node_modules` 等を除外。
- `PORT`: api は `8080` を明示(上記)。web / worker は Railway 注入(web は Caddy が `:{$PORT}` で待ち受け)。

## マイグレーション

- 現状は **自動実行されない**(api は起動時に DB へ Ping するだけ。空 DB でも `/api/healthz` は通る)。
- **Phase 1 で自動化する**: migration を `//go:embed` でバイナリに同梱し、Railway の **Pre-Deploy Command** で `goose up` を実行する(アプリ起動とは分離、多レプリカでも安全・冪等)。将来の Cloud Run では Cloud Run Job として同じバイナリを実行。

## 将来(全 GCP)

本格運用時は Cloud Run(web / api / worker)+ Cloud SQL + Secret Manager + Artifact Registry に移行する([ADR-0001](../../docs/adr/0001-hybrid-railway-gcp.md))。公開エッジは外部 HTTPS LB のパスルールで `/` と `/api` を振り分け、URL は不変に保つ。Dockerfile はそのまま Cloud Run でも使える。
