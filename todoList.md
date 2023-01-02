# TODO List

## 全体

- 認証・ログイン処理をどうするか
  - ログイン処理する部分を書く
- https化
  - initializerで証明書を読み込めるようにする
- Deprication Warning
  - Warning: ReactDOM.render is no longer supported in React 18. Use createRoot instead. Until you switch to the new API, your app will behave as if it's running React 17. Learn more: https://reactjs.org/link/switch-to-createroot
- .gorails.confみたいなファイルを用意してもいいかも
  - プロジェクト独自のgorailsに関する設定をかける
  - apiのidをgo側でIDと表現するようにするなど

## コマンド一覧

### generate

バグ: application.tsで改行されない
(Nice to have)viewのmethodは複数設定できるようにしたい

### server

(Nice to have)dev用環境変数のファイルの読み込み

### client

current directoryがclientなら移動しない
ServerAddressを環境変数などから設定できるようにする
