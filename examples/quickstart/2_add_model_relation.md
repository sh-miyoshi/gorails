# modelにリレーションをつける

※このドキュメントは[1_quick_start](./1_quick_start.md)のアプリが前提となっています。
まだやってない方は先にそちらを実施してください。

## 概要

このドキュメントではModelにリレーションをつける方法の解説をします。
具体的にはquick startの掲示板アプリのTopicにメッセージを保存できるようにしたいと思います。

### 出来上がりのイメージ

WIP

## 1. Modelにリレーションをつける

メッセージを保存するためのPost modelを作成し、Topicと関連付けてやります
関連性としては下の図のようにhas manyの関係です

![image](./assets/relations.png)

まずはPostモデルをコマンドで作成します

```bash
gorails generate model post --columns name:string --columns body:string --columns topic_id:string 
```

次にmodelファイルを修正して、TopicとPostを関連付けます

## 2. Controllerの作成

## 3. 画面の修正
