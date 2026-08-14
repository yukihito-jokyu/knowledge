# FEAT-001 Command Catalog（v1）

全操作は [../design.md](../design.md) の名前付きoption入力、JSON出力、共通envelope／error／exit codeに従う。初期提供に `search-semantic` は含めない。

| Operation | 分類 | 資料 | 状態 |
| --- | --- | --- | --- |
| `search-text` | lexical search | [search-text.md](commands/search-text.md) | v1 |
| `search-concept` | concept search | [search-concept.md](commands/search-concept.md) | v1 |
| `search-related` | relation search | [search-related.md](commands/search-related.md) | v1 |
| `get` | retrieval | [get.md](commands/get.md) | v1 |
| `get-evidence` | retrieval | [get-evidence.md](commands/get-evidence.md) | v1 |
| `search-contradictions` | relation search | [search-contradictions.md](commands/search-contradictions.md) | v1 |
| `search-temporal` | temporal search | [search-temporal.md](commands/search-temporal.md) | v1 |
| `create` | mutation | [create.md](commands/create.md) | v1 |
| `attach-evidence` | mutation | [attach-evidence.md](commands/attach-evidence.md) | v1 |
| `revise` | mutation | [revise.md](commands/revise.md) | v1 |
| `supersede` | mutation | [supersede.md](commands/supersede.md) | v1 |

Index management is a CLI responsibility. v1 creates and maintains lexical derived data within the same mutation or migration transaction and exposes no separate public operation (DEC-FEAT-002).
