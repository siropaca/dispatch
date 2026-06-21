---
status: accepted
date: 2026-06-21
---

# ADR-0009: モノリポは mise + Taskfile + pnpm

## Context and Problem Statement

ポリグロット(Go + TS)のモノリポ。ローカルのランタイム管理は mise を使う。モノリポのツールチェーンに何を採るか。

## Decision Drivers

* JS 寄りのモノリポツールを Go に強制しない
* 軽量さ
* Go / TS 双方のタスクを一元化

## Considered Options

* mise + Taskfile(go-task)+ pnpm workspace
* Nx / Turborepo
* 複数リポジトリ

## Decision Outcome

採用: "mise + Taskfile + pnpm workspace"。mise で toolchain を固定、Taskfile でタスクを統括、JS 側は pnpm workspace。Nx / Turborepo は不採用。

### Consequences

* Good: 素直で軽く、Go・TS 双方のタスクを Taskfile に一元化できる。
* Bad: タスクは Taskfile に集約する運用ルールが必要(`task dev` / `test` / `lint` / `gen` / `migrate` 等)。

## Pros and Cons of the Options

### mise + Taskfile + pnpm
* Good: ポリグロットに素直に適合。
* Bad: JS 専用モノリポ機能(影響分析等)は薄い。

### Nx / Turborepo
* Bad: JS 中心で Go と相性が悪い。

### 複数リポジトリ
* Bad: 連携・バージョン整合のコストが高い。
