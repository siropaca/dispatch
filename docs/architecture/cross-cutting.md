# 横断方針(テスト・ハーネス・安全・可観測性)

> 全体像は [`overview.md`](./overview.md)。

## 1. テスト(TDD)

- **Go**: ドメインは純粋関数として table-driven 単体。結合は **testcontainers-go**(本物の Postgres)。LLM / Search / Image / Queue / BlobStore は fake で決定的に。プロバイダ実物への contract test は opt-in。
- **Front**: Vitest + Testing Library。REST クライアントは MSW でモック。Playwright(E2E)は後。
- red → green → refactor を守る。テストが存在する機能を変えたらテストも更新する。

## 2. ハーネス(開発時)

- **AI エージェント開発**: root + 各 context の `AGENTS.md`、定型作業の project skills(context 追加 / REST エンドポイント追加 / タスクハンドラ追加 / migration + sqlc)、depguard で context 越境 import を禁止。
- **一般開発**: CI(GitHub Actions)= golangci-lint・biome・tsc・go vet・go test・vitest + codegen ドリフト検出 + goose マイグレーション検証。pre-commit は lefthook。

## 3. 製品の安全ガードレール(concept 由来・設計に内蔵)

| 項目 | 担保方法 |
|---|---|
| 出典明示 | report 型は出典 ≥1 を DB + ドメイン不変条件で強制 |
| AI 生成明示 | `images.ai_generated` 常 true + `posts.reporting_run_id` で来歴 → UI バッジ |
| 規制分野ガード(医療・法務・投資) | 取材パイプラインに差し込む pluggable policy(プロンプト制約 + 分類 → 免責付与 or ブロック) |
| 質問(Ask) | 本人のみ閲覧 + 回数上限 |

## 4. 可観測性 / 運用

- `log/slog`(JSON)、エラーは wrap して握りつぶさない、health / readiness エンドポイント。
- **コスト計測**: `reporting_runs` に LLM / 画像のトークン・費用を記録(AI コスト可視化)。
- OpenTelemetry は seam として用意(本番 = Cloud Trace / Monitoring)、Sentry は任意。

## 5. 設定 / シークレット

- 型付き env loader で読み込む。今は Railway env、本格運用では GCP Secret Manager。
- GCP 認証は今はサービスアカウント鍵、将来は Workload Identity。
