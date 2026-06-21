---
status: accepted
date: 2026-06-21
---

# ADR-0010: テストとハーネス

## Context and Problem Statement

TDD を守りつつ、「AI エージェント開発」と「一般開発」両方のハーネス(ガードレール)を整えたい。何を採るか。

## Decision Drivers

* 決定的なテスト
* 機械的なガードレール(口頭ルールに頼らない)
* AI エージェントが安全に開発できること

## Considered Options

* testcontainers + fake + CI + depguard + AGENTS.md / skills
* 手動 QA 依存
* E2E のみ
* 境界は口頭ルール

## Decision Outcome

採用: "testcontainers + fake + CI + depguard + AGENTS.md / skills"。Go は testing + testcontainers-go(LLM / Search / Image / Queue / BlobStore は fake)、Front は Vitest + Testing Library + MSW。CI(GitHub Actions)で lint / typecheck / test / codegen ドリフト / migration 検証、pre-commit は lefthook。depguard で context 越境を禁止し、root + 各 context に AGENTS.md と project skills を置く。

### Consequences

* Good: 決定的なテストと機械的ガードレールで、品質と AI エージェント開発の安全性を同時に担保。
* Bad: 初期セットアップにコストがかかる(以降の開発速度と安全性で回収)。

## Pros and Cons of the Options

### testcontainers + fake + CI + depguard + AGENTS / skills
* Good: 機械的に品質と境界を保証。
* Bad: 初期コスト。

### 手動 QA 依存 / E2E のみ
* Bad: 遅く脆い。

### 境界は口頭ルール
* Bad: 守られず崩壊する。
