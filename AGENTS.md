# AGENTS.md — Dispatch

Dispatch は専門特化 AI 記者による情報収集 SNS。詳細な設計はすべて `docs/` にある。このファイルは入口と最低限のルールだけを示す。進行状況やフェーズは [`docs/architecture/overview.md`](docs/architecture/overview.md) を参照する。

## ドキュメントマップ

| 知りたいこと | 場所 |
|---|---|
| プロダクト仕様(何を作るか) | [`docs/concept.md`](docs/concept.md) |
| 全体設計(アーキテクチャ・スタック・トポロジ) | [`docs/architecture/overview.md`](docs/architecture/overview.md) |
| データモデル・スキーマ | [`docs/architecture/data-model.md`](docs/architecture/data-model.md) |
| 意思決定の理由(ADR) | [`docs/adr/`](docs/adr/) |
| 各フェーズの仕様 | [`docs/specs/`](docs/specs/) |

迷ったら [`docs/architecture/overview.md`](docs/architecture/overview.md) を最初に読む。

## 最低限のルール

- **ファイルは 1 つ 200 行以内。** 超えたら分割する。
- 大きめの変更(複数ファイル・設計判断を伴う)は着手前に方針を提示して確認。小さい修正は即実行。スコープは依頼の最小限。
- **TDD(red → green → refactor)。** 変更後は関連テストを実行する。
- 回答・説明・コメントは日本語。コード内の識別子・コメントは英語。
- DDD 依存ルール: `domain ← app ← adapters`。`domain` は infra を import しない。context 越境 import は禁止(depguard で強制)。詳細は overview.md。
- エラーは wrap して握りつぶさない(空 catch・安易な fallback 禁止)。Go ログは `log/slog`、TypeScript は strict で `any` を避ける。
- 新規依存パッケージの追加は事前に確認する。
- codegen(sqlc / buf / openapi)の生成物はコミット必須。
- commit / push は明示的に依頼されたときのみ。メッセージはプレフィックス形式(`feat:` / `fix:` / `refactor:` / `chore:`)+ 日本語ボディ。
