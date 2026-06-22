# バックエンド設計

> 全体像は [`overview.md`](./overview.md)。本書はレイヤ・ライブラリ・ディレクトリ・ports・API を扱う。

## 1. レイヤ(ヘキサゴナル)

```
   ┌──────────── adapters(infra / delivery)─────────────────────┐
   │ driving: REST(oapi)・Connect(task)・scheduler tick           │
   │ driven : sqlc/pgx repo・OpenAI/Tavily・Cloud Tasks・GCS・Clerk │
   └───────────────────────┬─────────────────────────────────────┘
                implements ports / calls use cases
                  ┌─────────▼─────────┐
                  │  app(use cases)   │  ← トランザクション境界
                  └─────────┬─────────┘
                  ┌─────────▼─────────┐
                  │     domain        │  純粋・infra 非依存
                  │ 集約 / VO / port(if)│
                  └───────────────────┘
```

**依存ルール**: 依存は常に内向き。`domain` は何も import しない(infra 非依存)。`app` は domain を使い、port(interface)を介して外部を呼ぶ。`adapters` が port を実装する。context 間は互いの `internal` を直接参照せず、app の interface か ID / イベント経由。違反は depguard で弾く。

## 2. フレームワーク(ライブラリ構成)

Go の重量級 FW(Gin / Echo / Fiber)は不採用。薄いライブラリをエッジに置く。

| 層 | 採用 | 役割 |
|---|---|---|
| HTTP ルータ | **chi**(+ std `net/http`) | ルーティング・ミドルウェア(`http.Handler` 互換) |
| 公開 REST | **oapi-codegen** | `openapi.yaml` → サーバ interface 生成 |
| 内部 RPC | **connect-go**(buf) | タスク契約・サービス間(型付き) |
| DB ドライバ | **pgx v5** | Postgres 接続 |
| クエリ | **sqlc** | SQL → 型安全 Go |
| マイグレーション | **goose** | スキーマ変更 |
| キュー | **Cloud Tasks SDK** | ジョブ enqueue / handle |
| 設定 | typed env loader(env / koanf) | 設定 |
| ログ | `log/slog` | 構造化ログ |
| テスト | `testing` + **testcontainers-go** | 単体 / 結合 |
| DI | 手動コンストラクタ注入(`cmd` で wiring) | 肥大化したら google/wire |

## 3. ディレクトリ構成

```
backend/                  # Go module(1 つ)
  cmd/{api,worker}/       # 2 entrypoint、internal/ を共有
  internal/
    {identity,newsroom,reporting,publishing,timeline,interaction}/  # context は各 Phase で追加
      domain/   # 集約・VO・Repository interface(port)
      app/      # use case
      adapters/ # postgres(sqlc) / http(oapi) / connect
    platform/   # config, db(+sqlcgen), httpserver, httpapi(oapi 生成), telemetry … ai/queue/blobstore は Phase 1
    proto/      # buf 生成の内部 Connect(dispatch.reporting.v1 …)
  db/{migrations(goose),queries(sqlc)}/  sqlc.yaml
```

sqlc は `platform/db/sqlcgen` に共有の型付きクエリ層を生成し、各 context の `adapters/postgres` がそれを使って自分の domain port を実装する(sqlc 設定を単純に保ちつつ境界を維持)。

## 4. ports(domain / app に定義 → infra が実装)

`CorrespondentRepository` / `PostRepository` / `NotebookRepository` / `LLMProvider` / `SearchProvider` / `ImageProvider` / `QueueProvider` / `BlobStore` / `AuthProvider`

すべてローカル・テスト用の fake 実装を持つ(GCP 非依存のローカル開発と TDD のため)。

## 5. API スタイル

- **公開エッジ(ブラウザ SPA・将来の 3rd party / webhook)= REST + OpenAPI。** spec-first(`contracts/openapi.yaml`)で oapi-codegen が Go サーバ interface を、openapi-typescript が TS クライアント(`packages/api-client`)を生成。フロントに Connect は出さない。
- **サービス間(内部)= Connect-RPC。** ジョブ(タスク)契約を proto で定義し、Cloud Tasks が JSON で worker の Connect ハンドラに配送 → 内部契約が型付きになる。
- 生成物(sqlc / buf / openapi)はコミット必須。CI が最新でなければ fail(ドリフト検出)。

## 6. CQRS ライト

timeline などの読み取りは集約を経由せず、専用 sqlc クエリ(read model)で引く。書き込みは集約 + Repository を通す。イベントソーシングは採用しない。
