# ADR(Architecture Decision Records)

Dispatch の主要な技術的意思決定の記録(**MADR 4** 形式)。**1 ADR = 1 決定**。新規は [`template.md`](./template.md) をコピーして連番(`NNNN-*.md`)で追加する(`/adr` スキルでも作成可)。

| # | 決定 | ステータス |
|---|---|---|
| [0001](./0001-hybrid-railway-gcp.md) | インフラ: Railway × GCP ハイブリッド(将来 全 GCP) | 採用 |
| [0002](./0002-backend-modular-monolith-ddd-hexagonal.md) | バックエンド: モジュラモノリス + DDD + ヘキサゴナル | 採用 |
| [0003](./0003-api-rest-external-connect-internal.md) | API: 公開 REST + 内部 Connect-RPC | 採用 |
| [0004](./0004-auth-clerk.md) | 認証: Clerk(`AuthProvider` 抽象の裏) | 採用 |
| [0005](./0005-persistence-postgres-sqlc-goose.md) | 永続化: Postgres + pgx + sqlc + goose | 採用 |
| [0006](./0006-async-cloud-tasks-scheduler.md) | 非同期: Cloud Tasks + Cloud Scheduler | 採用 |
| [0007](./0007-ai-provider-abstraction.md) | AI: LLM / Search / Image の port 抽象 | 採用 |
| [0008](./0008-frontend-react-vite.md) | フロント: Vite + React + TanStack | 採用 |
| [0009](./0009-monorepo-toolchain.md) | モノリポ: mise + Taskfile + pnpm | 採用 |
| [0010](./0010-testing-harness.md) | テスト & ハーネス | 採用 |
