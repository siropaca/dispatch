---
status: accepted
date: 2026-06-21
---

# ADR-0002: バックエンドはモジュラモノリス + DDD + ヘキサゴナル

## Context and Problem Statement

Go バックエンドで DDD / TDD を守り、AI モデルを差し替え可能にし、長期保守に耐える構造が必要。どのアーキテクチャを採るか。

## Decision Drivers

* 配信機構(HTTP 等)を domain に漏らさない
* 将来の分割に備えた明確な境界
* 個人開発で運用が軽いこと

## Considered Options

* モジュラモノリス + DDD + ヘキサゴナル(薄い chi / connect-go / std)
* マイクロサービス
* Gin / Echo 等の FW 中心設計
* レイヤなしの素朴構成

## Decision Outcome

採用: "モジュラモノリス + DDD + ヘキサゴナル"。重量級 FW は不採用で chi + connect-go + std net/http を薄く使う。依存は `domain ← app ← adapters`。context 越境 import は depguard で禁止する。

### Consequences

* Good: 単一デプロイ単位で運用が軽く、境界は明確で将来分割に備えられる。
* Good: 配信機構を差し替え可能にでき、domain を純粋に保てる。
* Bad: ボイラープレートが増える(codegen と project skills で軽減)。

### Confirmation

depguard で context 越境・domain への infra import を CI で検出する。

## Pros and Cons of the Options

### モジュラモノリス + DDD + ヘキサゴナル
* Good: 境界の明確さと軽い運用を両立。
* Bad: 初期の記述量が増える。

### マイクロサービス
* Bad: 個人開発には過剰で運用コストが高い。

### Gin / Echo 等の FW 中心設計
* Bad: フレームワークが domain に漏れる。

### レイヤなしの素朴構成
* Bad: 境界が崩壊し保守不能になる。
