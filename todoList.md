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

(Nice to have)viewのmethodは複数設定できるようにしたい
generate modelで複数columnsを指定じに不要な改行が入る
generate apiで`[]Item`はyamlパースができないので`"`で囲って文字列にする必要がある
generate apiで配列指定の場合client側の生成がおかしくなる
(Nice to have)--with-testオプションの追加

### server

(Nice to have)dev用環境変数のファイルの読み込み

### client

ServerAddressを環境変数などから設定できるようにする
