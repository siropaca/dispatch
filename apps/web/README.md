# web — Dispatch フロントエンド(SPA)

ピュア React の SPA。全体設計は [`../../docs/architecture/overview.md`](../../docs/architecture/overview.md)、技術選定は [ADR-0008](../../docs/adr/0008-frontend-react-vite.md)。

## スタック

- Vite + React 19 + TypeScript(strict)
- ルーティング: TanStack Router / サーバ状態: TanStack Query
- スタイル: Tailwind CSS v4(`@tailwindcss/vite`)
- API: 公開 REST を生成クライアント [`@dispatch/api-client`](../../packages/api-client)(openapi-typescript 生成)経由で呼ぶ
- lint / format: Biome(リポジトリ root の `biome.json`)。型チェックは `tsc`

## 構成

```
src/
  main.tsx     # エントリ(QueryClientProvider + RouterProvider)
  router.tsx   # TanStack Router 定義(ルートツリー)
  health.tsx   # /healthz を取得して表示する最小ページ
  index.css    # @import "tailwindcss";
```

## 開発

リポジトリ root の Taskfile 経由を推奨(toolchain は mise 管理)。

```sh
task web:dev     # 開発サーバ(Vite)
task web:build   # 型チェック + 本番ビルド(tsc -b && vite build)
```

直接実行する場合は `pnpm dev` / `pnpm build` / `pnpm preview`。

dev サーバは `/healthz`・`/api` を api(`http://localhost:8080`)へプロキシする(`vite.config.ts`)。本番は api が SPA を同一オリジンで配信する想定。
