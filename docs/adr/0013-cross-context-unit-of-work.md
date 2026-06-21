---
status: accepted
date: 2026-06-22
---

# ADR-0013: コンテキスト跨ぎの書き込みは Unit of Work で 1 トランザクション

## Context and Problem Statement

取材 1 回で `publishing`(Post / Source / Image)・`newsroom`(Notebook)・`reporting`(run)の **3 context に書き込む**([`../architecture/infrastructure.md`](../architecture/infrastructure.md) のパイプライン)。[`overview.md`](../architecture/overview.md) は「context 間は app interface か ID / イベント経由」、[`backend.md`](../architecture/backend.md) は「app = トランザクション境界」とする。複数 context の書き込みを原子化しないと **部分コミット**(Post はできたが Notebook 未追記等)が起きる。モノリス + 単一 DB なので 1 tx で束ねられる。どう原子性と境界を両立するか。

## Decision Drivers

* 部分失敗を作らない(原子性)
* DDD 境界(context は互いの `internal` を直接触らない)を壊さない
* リトライ安全([ADR-0011](./0011-reporting-idempotency.md))と整合する

## Considered Options

* Unit of Work: app の use case が 1 tx を開始し、各 context の Repository port に tx ハンドルを渡して同一 tx で書く
* tx を `context.Context` に埋めて暗黙伝播
* context ごとに別 tx(結果整合・outbox / saga で補正)

## Decision Outcome

採用: "Unit of Work(明示 tx 伝播)"。`reporting` の app use case が Transactor で tx を開き、`publishing` / `newsroom` の app interface(または Repository port)に tx を **明示的に渡して** 同一 tx 内で Post・Source・Image・Notebook・run を確定する。context 越境は app interface 経由を維持し、`sqlcgen` は `pgx.Tx` 対応のクエリを使う。

### Consequences

* Good: 1 つの原子的コミットで部分失敗を排除し、境界を保ったまま整合を取る。
* Bad: tx ハンドルを配線する記述が増える。将来 context を別 DB / サービスへ割る時は outbox / saga へ移行が必要(モノリス前提の割り切り)。

### Confirmation

パイプライン途中で画像保存を失敗させ、Post も Notebook も run(succeeded)も書かれない(全ロールバック)ことを結合テスト(testcontainers)で検証する。

## Pros and Cons of the Options

### Unit of Work(明示 tx 伝播)
* Good: 明示的で、tx の有無が型に出る。テストしやすい。
* Bad: 配線コスト。

### context.Context 埋め込み
* Good: 呼び出しが軽い。
* Bad: 暗黙で追いづらく、tx の有無が型に出ない。

### 別 tx + saga / outbox
* Bad: モノリスには過剰で複雑。

## More Information

関連: [ADR-0002](./0002-backend-modular-monolith-ddd-hexagonal.md) / [ADR-0005](./0005-persistence-postgres-sqlc-goose.md) / [ADR-0011](./0011-reporting-idempotency.md)。
