# specs

各フェーズ・機能の確定した設計(spec)の目次。`/brainstorm` で壁打ちして `YYYY-MM-DD-<topic>-design.md` 形式で追加し(1 ファイル 200 行以内)、この表に 1 行足す。

| spec | 内容 | ステータス |
|---|---|---|
| [2026-06-21 Phase 0 基盤構築](./2026-06-21-phase-0-foundation.md) | walking skeleton の build plan | M1–M8 完了 |
| [2026-06-22 frontend/api 分離トポロジ](./2026-06-22-frontend-api-split-topology-design.md) | 単一オリジン + frontend/api コンテナ分離 | proposed |
| [2026-06-22 Notebook 検索 & 重複排除](./2026-06-22-notebook-retrieval-and-dedup-design.md) | 取材メモの検索戦略・重複排除 | proposed |

- 全体設計は [`../architecture/`](../architecture/)、意思決定は [`../adr/index.md`](../adr/index.md)
