---
status: accepted
date: 2026-06-21
---

# ADR-0007: AI は LLM / Search / Image の port 抽象

## Context and Problem Statement

「AI モデルを差し替え可能に」が要件。取材は LLM・web 検索・画像生成の 3 種を使う。どう抽象化するか。

## Decision Drivers

* ベンダーロックの回避
* TDD で決定的にテストできること
* 出典の制御しやすさ

## Considered Options

* `LLMProvider` / `SearchProvider` / `ImageProvider` の port 抽象
* 単一ベンダー直叩き
* LLM 内蔵 web 検索のみ
* langchaingo 等の重い抽象

## Decision Outcome

採用: "LLM / Search / Image の port 抽象"。`platform/ai` に 3 つの port を定義する。LLM は OpenAI / Anthropic を config で差し替え、Search は Tavily を既定(出典 URL が取れる)、Image は OpenAI gpt-image-1。すべて fake 実装を持つ。

### Consequences

* Good: ベンダーロックを避け、TDD でも決定的にテストできる。
* Bad: 薄い自前インターフェース + 公式 SDK ラップを保守する。

## Pros and Cons of the Options

### port 抽象
* Good: 差し替えとテストが容易。
* Bad: 保守対象のインターフェースが増える。

### 単一ベンダー直叩き
* Bad: 差し替え不可。

### LLM 内蔵 web 検索のみ
* Bad: 出典の制御が弱い。

### langchaingo 等の重い抽象
* Bad: MVP に過剰。
