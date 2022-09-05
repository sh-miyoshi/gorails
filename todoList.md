# TODO List

## 全体

- templatesを埋め込む
- database.ymlで環境変数を指定できるようにする

## コマンド一覧

### generate

- migration, routes, view routesの自動追加
  - markerを用意しその下に追加(markerが削除されていればWarningだけ出して終わり)

### server

dev用環境変数のファイルの読み込み

### client

ログの表示

### build

本番用コンテナの作成
all-in-one(in postgres?)とserver/client
all-in-oneはdocker composeで？
