---
status: accepted
date: 2026-06-21
---

# ADR-0006: 非同期処理は Cloud Tasks + Cloud Scheduler

## Context and Problem Statement

記者の定期取材は重い(web サーチ + LLM + 画像生成)。通常のリクエスト / レスポンスに収まらない。インフラは GCP 採用方針([ADR-0001](./0001-hybrid-railway-gcp.md))。非同期基盤に何を使うか。

## Decision Drivers

* GCP ネイティブで本番(Cloud Run)と同一
* Redis などの追加インフラを増やさない
* Cloud Run への移行容易性

## Considered Options

* Cloud Tasks + Cloud Scheduler
* River 等(Postgres / Redis ベースのキュー)
* 常駐ポーリング worker

## Decision Outcome

採用: "Cloud Tasks + Cloud Scheduler"。Cloud Tasks(HTTP push・リトライ・per-task スケジュール)+ Cloud Scheduler(cron)。worker は push を受ける HTTP ハンドラ(OIDC 検証)。ジョブ契約は Connect 型で定義する。

### Consequences

* Good: 本番(Cloud Run)と同一で Redis 不要、移行が容易。
* Bad: worker に公開 HTTP エンドポイントが必要(OIDC トークンで保護)。ローカル / テストは fake queue を使う。

## Pros and Cons of the Options

### Cloud Tasks + Cloud Scheduler
* Good: GCP と同一、追加インフラ不要。
* Bad: 今はクラウド越境(全 GCP 化で解消)。

### River 等
* Bad: GCP 方針と将来移行で不利。

### 常駐ポーリング worker
* Bad: Cloud Run のスケール特性と相性が悪い。

## More Information

関連: [ADR-0001](./0001-hybrid-railway-gcp.md)
