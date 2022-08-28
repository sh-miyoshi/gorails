# TODO List

## 全体

- templatesを埋め込む
- migration, routesの自動追加
  - markerを用意しその下に追加(markerが削除されていればWarningだけ出して終わり)

## コマンド一覧

### generate

- model
  - migration fileに追加
    - NG: templateを持っておいてそこに追加 → ユーザーが独自に追加したやつが消える
    - OK?: 特定文字を検索し、その次の行に挿入
      - 変にファイルを修正されてたら？(知らない、エラーにはしない)
- view?
  - 仕様決めから

### server

dev用環境変数のファイルの読み込み

### client

ログの表示

### build

本番用コンテナの作成
all-in-one(in postgres?)とserver/client
all-in-oneはdocker composeで？
