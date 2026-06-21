---
status: accepted
date: 2026-06-21
---

# ADR-0004: 認証は Clerk(`AuthProvider` 抽象の裏)

## Context and Problem Statement

認証は自作したくない(セキュリティリスク)。React + Go、Railway 中心の構成で、何を使うか。

## Decision Drivers

* セキュリティの自作回避
* React 親和・実装速度
* 将来の差し替え可能性

## Considered Options

* Clerk(`AuthProvider` 抽象の裏)
* 自前実装(Go セッション / JWT)
* Zitadel 自前ホスト(OIDC)

## Decision Outcome

採用: "Clerk"。`identity` context の `AuthProvider` 抽象の裏に置き、Go は JWT(JWKS)を検証して user を解決する。

### Consequences

* Good: React 親和で実装が最速、かつセキュア。
* Good: 抽象化により将来 Zitadel 等の自前 OIDC へ差し替えても `identity` 内で済む。
* Bad: 外部 SaaS 依存が一つ増える(`identity` 内に隔離)。

## Pros and Cons of the Options

### Clerk
* Good: DX と安全性が高い。
* Bad: 外部 SaaS 依存。

### 自前実装
* Bad: セキュリティ負担が大きい。

### Zitadel 自前ホスト
* Good: 外部依存なし・標準準拠。
* Bad: 今は運用負荷が過剰(将来の差し替え候補として残す)。
