# TODO List

## 全体

- modelでID, CreatedAt, UpdatedAtは自動で入れたい
- api interfaceを自動で作ってもいいかも
  - server側のinterface定義
  - client側のinterface定義

## コマンド一覧

### generate

- view routesの自動追加
  - markerを用意しその下に追加(markerが削除されていればWarningだけ出して終わり)
- viewのmethodは複数設定できるようにしたい
- 自動でviewのmethodの最初は大文字にする

### server

dev用環境変数のファイルの読み込み
DBが起動しているかチェック

### client

ログの表示
ServerAddressを環境変数などから設定できるようにする

### build

本番用コンテナの作成
all-in-one(in postgres?)とserver/client
all-in-oneはdocker composeで？
server and clientイメージとDBイメージの2つ
