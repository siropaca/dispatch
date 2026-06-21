---
status: accepted
date: 2026-06-22
---

# ADR-0014: 規制分野(医療・法務・投資)のガードレール

## Context and Problem Statement

ペルソナを持つ記者は出典付きの報道に加え「意見」も発信する([`../concept.md`](../concept.md))。医療・法務・投資の助言は **法務リスクが最も高い**。[`../architecture/cross-cutting.md`](../architecture/cross-cutting.md) では「pluggable policy」と一行で置かれるだけで具体化されていない。MVP の公式記者でも、誤った助言や無資格助言と受け取られる発信を避ける必要がある。どう守るか。

## Decision Drivers

* 法務・レピュテーションリスクの低減
* concept の「規制分野のガードレール」要件
* 取材パイプラインに差し込め、テストできること

## Considered Options

* 分類 + 免責付与 / ブロック + 人手レビュー(MVP)をパイプライン段として実装
* プロンプト制約のみ(LLM に「助言しない」と指示)
* 規制分野を全面禁止(扱わない)

## Decision Outcome

採用: "分類 + 免責 / ブロック + 人手レビュー(MVP)"。生成草稿を規制分野分類器(LLM / ルール)にかけ、(a) 該当かつ低リスク → 免責文 + opinion フラグを強制付与、(b) 高リスク(個別の診断 / 法的助言 / 投資勧誘)→ ブロックして run を `failed`(理由付き)。MVP の公式記者は公開前に人手レビューを併用する。policy はパイプラインの pluggable な 1 段として実装し、分野ごとに有効化する。

### Consequences

* Good: 高リスク発信を機械 + 人手の二段で止める。分野追加に開いている。
* Bad: 分類の誤判定(過剰ブロック / 見逃し)が残る。人手レビューは MVP のスループットを制限する。

### Confirmation

既知の高リスク文例で必ずブロック、境界例で免責が付与されることを fake LLM で検証する。判定結果は `reporting_runs` にログとして残す。

## Pros and Cons of the Options

### 分類 + 免責 / ブロック + レビュー
* Good: 多層防御。
* Bad: 運用コスト(レビュー)。

### プロンプト制約のみ
* Bad: すり抜ける。検証できない。

### 全面禁止
* Bad: プロダクト価値(専門分野)を毀損する。

## More Information

関連: [ADR-0007](./0007-ai-provider-abstraction.md) / [ADR-0011](./0011-reporting-idempotency.md)。安全方針は [`../architecture/cross-cutting.md`](../architecture/cross-cutting.md)、リスクは [`../concept.md`](../concept.md)。
