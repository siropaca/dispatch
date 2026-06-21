---
status: accepted
date: 2026-06-21
---

# ADR-0003: API は公開 REST + 内部 Connect-RPC

## Context and Problem Statement

Go ↔ React の型共有と、サービス間(内部)通信の両方が必要。どの API スタイルを採るか。

## Decision Drivers

* 公開 API の開きやすさ(外部 / webhook)
* 内部通信の型安全
* codegen による契約の機械検証

## Considered Options

* 公開 REST + 内部 Connect-RPC
* 全部 Connect-RPC
* 全部 REST
* GraphQL

## Decision Outcome

採用: "公開 REST + 内部 Connect-RPC"。公開エッジは REST(OpenAPI spec-first、oapi-codegen で Go、openapi-typescript で TS を生成)。サービス間は Connect-RPC(buf)で、Cloud Tasks のジョブ契約も proto で型付けする。生成物はコミットし CI でドリフト検出。

### Consequences

* Good: 公開は無難で外部にも開きやすく、内部は型安全で契約を機械検証できる。
* Bad: proto と openapi の二系統を保守する(codegen で吸収)。

## Pros and Cons of the Options

### 公開 REST + 内部 Connect-RPC
* Good: 適材適所。公開は標準的、内部は型安全。
* Bad: 生成系が二つ。

### 全部 Connect-RPC
* Bad: ブラウザ向けに過剰で、公開 API として開きにくい。

### 全部 REST
* Bad: 内部の型安全が弱い。

### GraphQL
* Bad: MVP に過剰。
