# frontend / api コンテナ分離 + 単一オリジン 設計

> 日付: 2026-06-22 / ステータス: **proposed**(実装前のレビュー待ち)
> 決定: [ADR-0015](../adr/0015-public-topology-edge-proxy.md)。前提: [`../architecture/infrastructure.md`](../architecture/infrastructure.md)。
> 本書は why と shape を残す。実装詳細はコードを正典とする。

## 目的

公開 URL を将来にわたり不変にしたまま、frontend と api を別コンテナに分ける。エッジプロキシ(Caddy)で単一オリジンを維持する([ADR-0015](../adr/0015-public-topology-edge-proxy.md))。

## URL 契約(これを固定する)

```
https://<独自ドメイン>/         → SPA(web / Caddy)
https://<独自ドメイン>/api/*    → api(Go)
```

将来 全 GCP では「外部 HTTPS LB が `/` と `/api` を振り分け」に置換し、**この URL は不変**。

## サービス構成

| サービス | 公開 | 役割 |
|---|---|---|
| **web**(Caddy) | 公開(独自ドメイン) | SPA 配信 + `/api/*` を api へリバースプロキシ。**唯一の公開エントリ** |
| **api**(Go) | プライベート | `/api` 配下で REST 配信。web からのみ到達 |
| **worker**(Go) | プライベート | 現状どおり(変更なし) |
| **postgres** | プライベート | 現状どおり |

## コンポーネント

### web コンテナ(新規)
- `deploy/docker/Dockerfile.web`: stage1 で `pnpm build`(Vite)→ stage2 `caddy:alpine` に `dist` と `Caddyfile` を載せる。
- `deploy/docker/Caddyfile`(概略):
  - `handle /api/*` → `reverse_proxy {$API_UPSTREAM}`(api のプライベート宛先)。
  - `handle` → `root * /srv` + `try_files {path} /index.html` + `file_server`(**SPA フォールバック**でクライアントルーティングを成立)。
  - listen は `:{$PORT}`(Railway / Cloud Run が注入)。

### api を `/api` 配下へ
- ルート結線を `httpapi.HandlerFromMuxWithBaseURL(server, r, "/api")` に変更(oapi-codegen 生成済みの BaseURL 版を使用)。
- probes は `/api/healthz`・`/api/readyz` になる(契約と一貫)。Railway / LB のヘルスチェックパスもこれに合わせる。
- `contracts/openapi.yaml` に `servers: [{ url: /api }]` を記載(クライアント基準 URL の明示)。
- 生成物(`httpapi.gen.go` / `schema.d.ts`)を再生成。

### api-client / フロント
- `getHealth` の呼び先を `/api/healthz` に変更(`baseUrl` 既定は同一オリジン)。
- 追加エンドポイントは今後すべて `/api` 配下。

### dev ワークフロー
- `task dev`: api(:8080・`/api` 配下)/ worker / web(Vite :5173)。
- Vite proxy を `'/api' → http://localhost:8080` に変更(旧 `/healthz` proxy は廃止)。Caddy は dev では使わない(本番パリティは Docker イメージで担保)。

## Railway デプロイ

- **web**: `Dockerfile.web`、独自ドメインを付与、`API_UPSTREAM=<api のプライベート宛先>`(`*.railway.internal:PORT`、セットアップ時に確定)。Healthcheck `/api/healthz`(proxy 経由)。
- **api**: 公開ドメインを外しプライベート化。Healthcheck `/api/healthz`。
- `deploy/railway/README.md` を 4 サービス構成へ更新。Clerk 許可オリジンは単一ドメイン(Phase 2 で設定)。

## エラー処理 / 留意

- web → api の到達不可時、Caddy は 502 を返す(SPA は `/api` のエラーをハンドリング)。
- private networking の宛先(内部ホスト名・ポート・IPv6)は Railway 仕様に従いセットアップ時に確定。
- worker は Caddy の背後に置かない(独立プライベート、`/healthz` は現状維持)。

## 検証

- 既存テスト更新(TDD): `health_test` / `ready_test` のパスを `/api/...` に。`HandlerFromMuxWithBaseURL` 結線で `/api/readyz` が応答することをルートテストで確認。
- 生成ドリフト: `task gen` 後に差分が出ないこと。
- web ビルド(`pnpm build`)/ biome / golangci / go test が green。
- 任意: ローカルで Caddy イメージを起動し `/`(SPA)と `/api/healthz`(→api)へ到達するスモーク。

## 非ゴール

- 認証(Phase 2)。CORS は単一オリジンのため不要。
- SSR / CDN(将来必要なら別途)。
