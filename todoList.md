# TODO List

## コマンド一覧

### new

clientのシステムファイルを用意

### generate

- model
  - migration fileに追加
- view?
  - 仕様決めから
- migration?
- scaffold?
  - 多分不要
- job?

#### view

- tsconfigが欲しい
  - tsconfig.json
- axiosが入ってない？

```json
{
  "compilerOptions": {
    "jsx": "react-jsx"
  },
  "include": [
    "src"
  ]
}
```

```bash
vi src/pages/topics/index/index.tsx
vi src/index.tsx
```

```js
// src/pages/topics/index/index.tsx
const TopicsIndex = () => {
  return (
    <div>
      Topic Index
    </div>
  )
}

export default TopicsIndex

// src/index.tsx
import TopicsIndex from './pages/topics/index/index'
<Route path="/topics" element={<TopicsIndex />} />
```

### server

dev用環境変数のファイルの読み込み

### run server/client

server, clientコマンドのエイリアス
→ 動作確認

### build

本番用コンテナの作成
all-in-one(in postgres?)とserver/client
all-in-oneはdocker composeで？
