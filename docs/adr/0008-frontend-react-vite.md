---
status: accepted
date: 2026-06-21
---

# ADR-0008: フロントは Vite + React + TanStack

## Context and Problem Statement

ピュア React + TypeScript の SPA。ルーター等は推奨に委ねられた。何を採るか。

## Decision Drivers

* 型安全なルーティング / データ取得
* 軽量
* 「ピュア React」方針

## Considered Options

* Vite + React + TanStack(Router / Query)+ Tailwind
* Next.js 等のメタフレームワーク
* React Router
* Redux 中心の状態管理

## Decision Outcome

採用: "Vite + React + TS(strict) + TanStack Router + TanStack Query + Tailwind CSS"。REST クライアントは openapi-typescript で生成。MVP では api が静的配信し single origin(CORS 回避)。

### Consequences

* Good: 型安全なルーティングとサーバ状態管理、軽量、生成クライアントで公開 REST と型が一致。
* Bad: クライアント側ルーティング。SEO が要る面が出たら将来 SSR を別途検討する。

## Pros and Cons of the Options

### Vite + React + TanStack
* Good: 型安全で軽量。
* Bad: TanStack の学習コスト。

### Next.js 等のメタフレームワーク
* Bad: 「ピュア React」方針と不一致。

### React Router
* Bad: 可だが型安全性は TanStack Router が上。

### Redux 中心
* Bad: server state は TanStack Query で十分。
