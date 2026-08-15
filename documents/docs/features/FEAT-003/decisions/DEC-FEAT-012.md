# DEC-FEAT-012: 公開記事だけを取得するURL安全境界

- **Status:** decided
- **Level:** L2
- **Decision:** Article Analysisは、利用者入力時と各リダイレクト前にURLのscheme、authority、解決先を検査し、公開インターネット上のHTTPまたはHTTPS記事だけを取得する。検査不能または公開外と判定した対象へは接続せず、`article_unavailable`を返す。

## 判定規則

- URLはuserinfoを持たず、HTTPは既定port 80、HTTPSは既定port 443だけを用いる。
- hostnameは、名前解決後のすべての接続候補がグローバル到達可能なIPアドレスである場合だけ公開とみなす。loopback、private、link-local、unique-local、multicast、unspecified、予約済みのIPアドレス、およびローカルホスト名は公開外である。
- 初期URLとリダイレクト先のすべてに同じ検査を適用する。検査を通らないURLへはHTTP接続しない。接続は検査済みの接続候補だけへ固定し、接続時にhostnameを再名前解決しない。同等の取得先固定を保証できない取得能力は使用しない。
- 取得は読み取りだけとし、利用者の認証情報・セッション・Cookieを送らず、フォーム送信その他の状態変更操作を行わない。

## 根拠

DEC-FEAT-010は利用者が指定した公開技術記事をCodex会話で評価することを承認した。HTTP(S)だけでは内部資源への到達を防げないため、URLとリダイレクトの両方について論理的な公開性検査が必要である。具体的なHTTPクライアントやネットワーク実装は固定しない。
