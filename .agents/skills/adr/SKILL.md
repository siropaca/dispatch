---
description: 技術的な意思決定を ADR(MADR 4 形式)として docs/adr/ に記録する。アーキテクチャ・スタック・設計方針を決めた/変えたとき、トレードオフを伴う選択をしたときに使う(/adr [決定の要約])。
allowed-tools: Read, Write, Edit, Glob, Grep, Bash
---

# /adr — 意思決定を ADR(MADR 4)として docs/adr/ に残す

技術的な決定を 1 ファイルに記録し、後から「なぜそうしたか」を辿れるようにする。書式は **MADR 4.x**(見出しは英語の標準名、本文は日本語)。

## いつ使うか

- アーキテクチャ / スタック / 設計方針を決めた、または変更した
- トレードオフを伴う選択をした(代替案を却下した)
- **1 ADR = 1 決定**。軽微な実装判断は対象外(コードとコミットで足りる)

## 手順

1. **文脈を集める**: 関連する会話・[`docs/architecture/`](../../../docs/architecture/)・既存 ADR を確認し、1 決定にスコープを絞る。
2. **番号を決める**: `ls docs/adr` で既存 `NNNN-*.md` の最大番号 +1(4 桁ゼロ詰め)。
3. **作成**: `docs/adr/NNNN-<kebab-title>.md`。[`template.md`](../../../docs/adr/template.md) をコピー。日付は `date +%F`。
4. **記入(MADR)**: frontmatter(status / date)/ Context and Problem Statement / Decision Drivers / Considered Options / Decision Outcome(+ Consequences)/ Pros and Cons of the Options。任意項目(Confirmation / More Information / consulted 等)は不要なら削る。what ではなく **why** を書く。1 ファイル 200 行以内。
5. **索引を更新**: `docs/adr/index.md` の表に 1 行追加(番号・決定の要約・status)。
6. **関連付け**: 既存 ADR を置換する場合は双方にリンクし、旧側の status を `superseded by ADR-XXXX` に更新。
7. **確認**: 作成したパスをユーザーに伝える。コミット / push は明示依頼時のみ。

## status

`proposed` / `accepted` / `rejected` / `deprecated` / `superseded by ADR-XXXX`。新規の確定事項は通常 `accepted`。

## 書式

[`docs/adr/template.md`](../../../docs/adr/template.md)(MADR 4)をコピーして使う。既存の ADR(`0001-*` 以降)も参照し、トーンと粒度を合わせる。
