---
status: accepted
date: 2026-06-22
---

# ADR-0011: 取材ジョブは reporting_run_id で冪等にする

## Context and Problem Statement

Cloud Tasks([ADR-0006](./0006-async-cloud-tasks-scheduler.md))は **at-least-once** 配送で、リトライ・重複配送が起こり得る。取材は web 検索 + LLM + 画像生成で高コストかつ副作用(Post 作成・課金)を伴うため、重複実行で **二重課金・重複つぶやき** が発生する。現状の `RunReportingRequest` は `correspondent_id` のみで冪等キーが無い。どう冪等性を担保するか。

## Decision Drivers

* 二重課金・重複 Post を防ぐ
* Cloud Tasks の再送に耐える(リトライ安全)
* 取材の来歴・コストを記録する `reporting_runs` と整合する

## Considered Options

* enqueue 時に `reporting_runs` 行を作り `reporting_run_id` をタスクに載せ、worker はその id で冪等化
* worker 側で `(correspondent_id, 期間)` のユニーク制約で重複検知
* 重複を許容し事後で dedup

## Decision Outcome

採用: "enqueue 時に run_id 採番"。api(enqueuer)が `reporting_runs` を `status=queued` で作成し、`RunReportingRequest` に `reporting_run_id` を追加して載せる。worker は受信時に run を `running` へ条件付き遷移(既に `running`/`succeeded` なら no-op)し、`produced_post_id` の有無で Post 作成済みかを判定して再課金しない。

### Consequences

* Good: 再送に強く、二重課金・重複 Post を防ぐ。来歴(run)が常に先に存在し、コスト記録の器になる。
* Bad: enqueue と worker の 2 箇所で run の状態遷移を扱う。proto に `reporting_run_id` を足す(契約変更)。

### Confirmation

同一 `reporting_run_id` を 2 回 push する結合テスト(testcontainers)で、Post が 1 件・課金 1 回に収まることを検証する。

## Pros and Cons of the Options

### enqueue 時 run_id 採番
* Good: 明示的で、来歴・コスト記録と一体。部分失敗の再開判定も run 状態で行える。
* Bad: run の状態機械を持つ必要がある。

### ユニーク制約のみ
* Good: 実装が単純。
* Bad: 「期間」の境界定義が曖昧で、部分失敗からの再開が難しい。

### 重複許容 + 事後 dedup
* Bad: コストと UX(重複つぶやき)が許容できない。

## More Information

関連: [ADR-0006](./0006-async-cloud-tasks-scheduler.md) / [ADR-0013](./0013-cross-context-unit-of-work.md)(部分失敗の原子性)。データモデルは [`../architecture/data-model.md`](../architecture/data-model.md) の `reporting_runs`。
