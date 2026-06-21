---
status: accepted
date: 2026-06-21
---

# ADR-0005: 永続化は Postgres + pgx + sqlc + goose

## Context and Problem Statement

Postgres を採用。DDD と両立し、型安全で、ORM の抽象漏れを避けたい。データアクセス層に何を使うか。

## Decision Drivers

* SQL を明示し ORM マジックを避ける
* 型安全
* DDD 境界の維持

## Considered Options

* pgx + sqlc + goose
* GORM
* ent
* 生 SQL(database/sql 手書き)

## Decision Outcome

採用: "pgx v5 + sqlc + goose"。ID は UUID v7、一覧はカーソルページング。sqlc は `platform/db` に共有生成し、各 context の `adapters/postgres` がそれを使って domain の Repository port を実装する。

### Consequences

* Good: SQL を明示でき、生成コードで型安全。
* Good: 共有生成層 + context 別 adapter で sqlc 設定を単純に保ちつつ境界を維持。
* Bad: SQL を書く必要がある。テストは testcontainers-go で本物の Postgres を使う。

## Pros and Cons of the Options

### pgx + sqlc + goose
* Good: 型安全かつ明示的。
* Bad: 記述量が増える。

### GORM / ent
* Bad: ORM の抽象漏れ・マジックが DDD と相性が悪い。

### 生 SQL(database/sql)
* Bad: 型安全性がなく保守が辛い。
