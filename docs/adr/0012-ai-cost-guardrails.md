---
status: accepted
date: 2026-06-22
---

# ADR-0012: AI コストは事前ガード(per-run 上限 + 予算上限)

## Context and Problem Statement

取材は LLM・web 検索・画像生成で課金が発生する。`reporting_runs.cost`([`../architecture/cross-cutting.md`](../architecture/cross-cutting.md))は **結果を記録するだけ** で暴走を止めない。バグ・ループ・プロンプト膨張で予算を焼く事故は AI パイプラインで頻出する。どうコストの上限(ブラスト半径)を画すか。

## Decision Drivers

* 想定外コストの上限を事前に画す(被害の最小化)
* 個人開発の予算を守る
* fake provider で決定的にテストできること

## Considered Options

* パイプラインに pre-flight 上限(per-run の検索回数 / トークン / 画像枚数)+ 記者ごと・全体の日次予算
* 事後集計のみ(現状)+ アラート
* プロバイダ側のハード制限のみに依存

## Decision Outcome

採用: "事前ガード"。`platform/ai` の port(LLM/Search/Image)呼び出しをコストメータでラップし、per-run 上限(検索 N 回・入出力トークン・画像 ≤1 など)を超えたら run を `failed`(理由付き)で打ち切る。記者ごと / 全体の日次予算超過は、スケジューラの取材対象抽出段でスキップする。実コストは従来どおり `reporting_runs.cost` に記録。閾値は config 化する。

### Consequences

* Good: 事故時の被害を上限で抑える。閾値を env で運用調整できる。
* Bad: メータリングの実装と、上限到達時の打ち切り処理が要る。

### Confirmation

上限超過で run が `failed` になり、それ以降のプロバイダ呼び出しが行われないことを fake provider で検証する。

## Pros and Cons of the Options

### 事前ガード
* Good: 予防的でブラスト半径を画せる。記者単位の制御ができる。
* Bad: 実装コスト。

### 事後集計のみ
* Bad: 焼けてから気づく。

### プロバイダ制限依存
* Bad: 粒度が粗く、記者単位の制御ができない。

## More Information

関連: [ADR-0007](./0007-ai-provider-abstraction.md) / [ADR-0011](./0011-reporting-idempotency.md)。コスト記録は [`../architecture/cross-cutting.md`](../architecture/cross-cutting.md)。
