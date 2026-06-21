---
status: accepted
date: 2026-06-22
---

# ADR-0015: 公開トポロジはエッジプロキシで単一オリジン、frontend/api はコンテナ分離

## Context and Problem Statement

[ADR-0008](./0008-frontend-react-vite.md) では MVP で api が SPA を静的配信し single origin としていた。しかし将来フロントを別配信に分けると **公開 URL が変わり**、ブックマーク・SEO・Clerk の許可オリジン・OAuth リダイレクト URI 等を巻き込む。URL を最初から固定したい。frontend と api を分離しつつ公開 URL を不変にするトポロジは何か。

## Decision Drivers

* 公開 URL を将来にわたり不変にする
* frontend と api を独立にビルド / デプロイ / スケールできる
* 将来の全 GCP(Cloud Run + 外部 HTTPS LB、[ADR-0001](./0001-hybrid-railway-gcp.md))へ素直に対応
* 個人開発の運用が軽い

## Considered Options

* 単一ドメイン + エッジプロキシ(Caddy)でパスルーティング(`/` → frontend、`/api/*` → api)
* サブドメイン分割(app と `api.*` を別オリジン、CORS)
* api が SPA を配信(ADR-0008 現状、単一コンテナ)

## Decision Outcome

採用: "単一ドメイン + エッジプロキシ"。独自ドメインを公開エントリにし、Caddy が `/` を frontend コンテナ、`/api/*` を api コンテナへ振り分ける。api は公開ドメインを持たず **プライベート**(web からのみ到達)。同一オリジンなので CORS 不要。Railway は web(Caddy)/ api / worker / postgres の構成。将来は外部 HTTPS LB のパスルールへ 1:1 で移行し、**公開 URL は不変**。

### Consequences

* Good: 公開 URL が 1 つで永続。frontend / api を独立運用でき、CORS が無い。GCP LB に直マップできる。
* Bad: Railway に web(Caddy)サービスが 1 つ増える。private networking の配線が要る。

### Confirmation

デプロイ後、公開ドメインで `/`(SPA)と `/api/healthz`(api)が到達することを確認する。

## Pros and Cons of the Options

### 単一ドメイン + エッジプロキシ
* Good: 単一オリジン・URL 不変・GCP 親和。
* Bad: プロキシ 1 つ分の運用。

### サブドメイン分割
* Good: プロキシ不要で単純。
* Bad: CORS が要る。オリジンが 2 つに増える。

### api が SPA 配信(現状)
* Bad: 後で分けると公開 URL が変わる(本 ADR の動機)。

## More Information

[ADR-0008](./0008-frontend-react-vite.md)(フロント技術選定。配信トポロジは本 ADR で更新)/ [ADR-0001](./0001-hybrid-railway-gcp.md)。実装設計は [`../specs/2026-06-22-frontend-api-split-topology-design.md`](../specs/2026-06-22-frontend-api-split-topology-design.md)。
