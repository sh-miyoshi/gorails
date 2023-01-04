# gorails

gorailsはRuby on Railsを参考にした、Go言語とReactでWebアプリケーションを作成するためのコマンドラインツールです。
ソースコードのディレクトリ構成など、Ruby on Railsと似たような感覚でMVCモデルのWebアプリケーションを作成できます。

## 作成するアプリケーションアーキテクチャ

gorailsで作成するアプリケーションは基本的にRuby on Railsと同じなのですが、フロントエンドとしてReactを使うため以下のような構成になります。
※Go言語のtemplate packageを使うことでRailsのほぼ同じ構成をとることもできるのですが、2023年の現実問題としてフロントエンドの開発はReact一強なため、最初からこの構成にしています。

TODO: 画像

## インストール手順

まずはGo言語とNode.js(npm)をインストールします

- Go言語インストール手順: https://go.dev/doc/install
- Node.jsインストール手順: https://nodejs.org/ja/

次にgorailsコマンドをインストールします

```bash
go install github.com/sh-miyoshi/gorails@latest
gorails --help
```

## クイックスタート

[examples/quickstart](./examples/quickstart)にサンプルWebアプリケーションの作り方をまとめたので参考にしてください。

## Notes

gorailsではModelの実装として[GORM](https://gorm.io/ja_JP/docs/index.html)を使用しています。
そのため、gorailsコマンドでModelを生成後はGORMのフォーマットに則って実装してください。
