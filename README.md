製品出荷管理CRUD APIサーバー

* 梱包
* 出荷
* 納期

の日付等のデータをJSONに保存し、操作するAPIを提供します。


## Quick Start
1. schd を実行する
1. ブラウザアドレス欄に http://localhost:8080/view/list を打ち込み、テーブル表示を確認する。
1. ブラウザアドレス欄に http://localhost:8080/data/list を打ち込み、JSONを確認する。
1. ブラウザアドレス欄に http://localhost:8080/data/cal を打ち込み、JSONを確認する。
1. ブラウザアドレス欄に http://localhost:8080/data を打ち込み、JSONを確認する。 このJSONはtest/sample.json の中身と同一。
1. ブラウザアドレス欄に http://localhost:8080/data/741744 を打ち込み、JSONを確認する。
このJSONはtest/sample.json の中身のキー=741744 オブジェクトと同一。
また、http://localhost:8080/view/list に表示されるリンクの遷移先。


## API

| 説明 | メソッド | URI | パラメータ | パラメータ例 |
|----|------|-----|-------|-------|
| 保存されているJSONの中身を出力 | GET | /data |  数字6桁 | http://localhost:8080/data/000000 |
| 日付をキーに、項目ごとに製番リストを保持するCal構造体をJSONで返す | GET | /data/cal |  なし | http://localhost:8080/data/cal |
| Cal構造体から日付をプライマリキーとするテーブル形式のRowsをJSONで返す | GET | /data/list |  なし |  http://localhost:8080/data/list |
| Cal構造体から日付をプライマリキーとするテーブル形式のRowsをHTMLで返す | GET | /view/list |  なし |  http://localhost:8080/view/list |

