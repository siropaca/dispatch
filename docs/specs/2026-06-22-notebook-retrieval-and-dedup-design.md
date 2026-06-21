# Notebook 検索 & ニュース重複排除 設計(Phase 1)

> 日付: 2026-06-22 / ステータス: **proposed**(設計方針の記録。詳細は Phase 1 着手時に確定)
> 前提: [`../architecture/overview.md`](../architecture/overview.md)、取材パイプライン [`../architecture/infrastructure.md`](../architecture/infrastructure.md)、[ADR-0007](../adr/0007-ai-provider-abstraction.md) / [ADR-0011](../adr/0011-reporting-idempotency.md) / [ADR-0012](../adr/0012-ai-cost-guardrails.md)。
> 本書は **why と shape** を残す。実装詳細(DDL・シグネチャ)はコードを正典とする。

## 解決する課題

1. **Notebook の肥大化**: `notebook_entries` は追記のみで記者ごとに無限に増える([data-model](../architecture/data-model.md))。全件を LLM に渡すとコンテキスト窓を超え、コストが線形に膨張する。「進化」を支えつつ有限のプロンプトに収める検索戦略が要る。
2. **重複発信**: 狭い専門の記者が毎日同じ分野を取材すると、同じソース・話題に何度も当たる。同じニュースを二度つぶやかせない。
3. **継続性(要件)**: 記者は **自分の最近の発信(Post)も踏まえて** 次を書く。過去発信との一貫性を保ち、自己重複を避ける。

## 方針(決定)

**段階導入(ハイブリッド)** を採る。検索を **port 化** し、Phase 1 は recency + 構造化タグ + 要約圧縮で実装する。embeddings / pgvector は **採用せず**、Notebook が育った段階で drop-in 追加できる形に予約する。

- 理由: 要件 3(最近の自分の発信を考慮)は recency クエリで満たせ embeddings は不要。pgvector は新規依存(Postgres 拡張 + 埋め込みモデル + 取材ごとの埋め込みコスト)で、Phase 1 の小さな Notebook には過剰(YAGNI)。port で隔離すれば後の移行が安い。
- 却下: 「pgvector を最初から」= 今は過剰・新規依存。「要約のみ(将来も最小)」= 意味的な近重複検知が将来も弱い。

## アーキテクチャ

取材 use case(`reporting`)が **ReportingContext** を組み立て、LLM 草稿生成に渡す。組み立ては port の裏に隠す:

```
ReportingContextBuilder  (port: app/domain に定義、adapters が実装)
  Build(ctx, correspondentID) -> ReportingContext
```

- Phase 1 実装: recency + tags + summary entry + recent posts(下記)。
- Phase 2 実装: 上記 + ベクトル近傍検索(pgvector)。port 差し替えのみで載る([ADR-0007](../adr/0007-ai-provider-abstraction.md) の port 思想と一貫)。

### ReportingContext の中身

1. **Field**(記者の狭い担当分野)— 検索クエリの素。
2. **最近の Notebook entries**(recency 上位 N + `summary` entry)。
3. **記者自身の最近の Post**(`published_at` 降順 M 件)— 継続性と自己重複回避(**要件 3**)。
4. (予約)embeddings による過去 entry の意味検索 — Phase 2。

## Notebook 圧縮

`notebook_entries.kind` を使い分ける:

- `observation`: 取材で得た生の知見(従来どおり追記)。
- `summary`: 古い `observation` 群を定期的に 1 件へ要約圧縮した entry。

圧縮は **追記のみを維持**(既存行は不変、summary は新規 entry)。これで「最近の窓」クエリが有界になり、古い文脈は summary 経由で参照できる。圧縮の起動は「entry が K 件たまるごと」または定期ジョブ([ADR-0006](../adr/0006-async-cloud-tasks-scheduler.md) の Scheduler)。

## 重複排除(2 層)

1. **ソース単位(完全一致・安価)**: 候補ソース URL を正規化(`utm_*`・fragment 除去、scheme/host 小文字化、末尾 `/` 正規化)し、**この記者が最近引用済みの URL**(sources → posts → correspondent、時間窓)と突き合わせて既出を除外。決定的で純粋関数としてテスト可能。
2. **話題単位(最近の自分との照合)**: ReportingContext の「最近の Post + summary」を草稿プロンプトに含め、「既出の話題は繰り返さない。新規性のある知見が無ければ投稿を生成しない」を指示。意味的近重複は当面 LLM 判断に委ね、ベクトル近傍は Phase 2。

### 新規性ゼロの扱い

新規の知見が無い run は **`succeeded` + `produced_post_id = null`**(正当な no-op)。

- [ADR-0011](../adr/0011-reporting-idempotency.md) の冪等設計と整合(二重投稿しない)。
- concept の「動きの少ない分野は無理に頻繁につぶやかない」= 発信頻度の自己調整の素地になる。

## スキーマへの影響

- `notebook_entries`: 追加なし。`embedding` 列は **pgvector 採用時に追加**(拡張が要るため今は作らない・本書で予約のみ)。
- `sources`: 既出 URL 照合を索引で引くため、**正規化 URL の保持を検討**(`canonical_url` 列 + index、または既存 `url` を保存し照合は窓内取得 + アプリ正規化)。← 実装時に確定(下記「未決」)。

## エラー処理

- ReportingContext 構築失敗 → run を `failed`(理由付き)。盲目生成しない。エラーは wrap。
- 要約圧縮失敗 → ログして recency のみの窓に degrade(ブロックしない)。ただし窓は件数上限でキャップしコスト暴走を防ぐ([ADR-0012](../adr/0012-ai-cost-guardrails.md))。
- 新規性ゼロ → エラーではなく no-op(上記)。

## 検証(TDD)

- **URL 正規化 / 既出判定**: table-driven の純粋関数テスト(utm 除去・fragment・末尾スラッシュ等)。
- **ReportingContextBuilder**: fake repo で「最近の Post + 最近の entry + summary」が含まれ、古い observation は要約に畳まれることを決定的に検証。
- **新規性ゼロ経路**: fake LLM が「新規なし」を返すと run = `succeeded` / post 無し。
- 結合は testcontainers(本物 Postgres)で時間窓クエリを検証。

## コスト

recency / tags は安価な SQL。要約は **定期 LLM コスト(有界)**。Phase 1 は **取材ごとの埋め込みコストが無い**。([ADR-0012](../adr/0012-ai-cost-guardrails.md))

## 未決(Phase 1 実装時に確定する)

1. **`sources.canonical_url` 列を Phase 1 で足すか**、照合時にアプリ正規化で済ますか(索引性能 vs スキーマ最小)。
2. **要約圧縮の起動条件**(件数 K / 定期 / 併用)と「最近の窓」N・M の初期値。
3. pgvector 移行の **発火条件**(Notebook 件数・想定コスト・体感品質のしきい)。
