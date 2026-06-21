---
status: accepted
date: 2026-06-21
---

# ADR-0001: インフラは Railway × GCP ハイブリッド(将来 全 GCP)

## Context and Problem Statement

個人開発で本番公開を目指す。小規模スタートだが、本格運用時はスケール・可用性を見据えて全 GCP に寄せたい。どこに何をホストするか。

## Decision Drivers

* 個人開発の DX とコスト
* 本番(全 GCP)への移行リスクを下げる
* キュー・スケジューラ・オブジェクト保存など周辺インフラの質

## Considered Options

* Railway × GCP ハイブリッド(compute / DB は Railway、周辺は GCP)
* 全部 Railway
* 最初から全 GCP

## Decision Outcome

採用: "Railway × GCP ハイブリッド"。compute(api / worker)と DB は Railway、キュー・スケジューラ・オブジェクト保存は GCP(Cloud Tasks / Cloud Scheduler / GCS)。本格運用時は compute を Cloud Run、DB を Cloud SQL、シークレットを Secret Manager に移し全 GCP 化する。

### Consequences

* Good: GCP ネイティブ部品(Tasks / Scheduler / GCS)を最初から本番と同一にでき、移行リスクが小さい。
* Good: Railway の DX とコストで個人開発が軽い。
* Bad: 一時的にクラウドを跨ぐ越境レイテンシが出る(全 GCP 化で解消)。

## Pros and Cons of the Options

### Railway × GCP ハイブリッド
* Good: 移行が config 中心(`QueueProvider` / `BlobStore` 抽象で隔離)。
* Bad: 二クラウド運用の複雑さ。

### 全部 Railway
* Bad: キュー等のマネージド部品が弱く、将来の GCP 移行で作り直しになる。

### 最初から全 GCP
* Bad: 個人開発の初期には運用・コストが過剰。
