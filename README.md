# blog-server
![Lint](https://github.com/masibw/blog-server/workflows/Lint/badge.svg)
![Test](https://github.com/masibw/blog-server/workflows/Test/badge.svg)

新 `mesimasi.com` のサーバーサイドです。

なんかやばそうなところを見つけたらこっそり教えてください

# 使い方
`.env.local`と`.env.test`,`.env.prod`を適宜作成・書き換えてください  
S3へ画像をアップロードするための`GET /api/v1/images`を使うには`.env.secret`を設定する必要があります

`make up`でNginx,App,MySQLが起動します  
`make down`でNginx,App,MySQLが終了します

## テスト
`make up-test`でテスト用のMySQLが起動します  
`make test`でテストを実行します  
`make down-test`でテスト用のMySQLが終了します  

## ログ
`make logs T=app`のように`make logs T={コンテナ名}`でコンテナのログを見れます  

## 管理用ユーザーの作成
記事を書いたりするための運営用ユーザーの作成はAppを起動した状態で`make admin`を実行すると行えます  
メールアドレスとパスワードを入力してください

