# Use Microsoft Entra ID for cache authentication

Azure Redis AAD 認証の Go サンプル

https://learn.microsoft.com/en-us/azure/azure-cache-for-redis/cache-azure-active-directory-for-authentication

## 1. 事前準備

### 1.1. Azure Portal で Azure Cache for Redis を作成

1. AAD認証を有効にする
2. アクセスポリシー追加

## 2. アプリケーションの実行

以下のシナリオでは、az loginして、az cli トークンを使ってのアクセスを想定

ログインする

```
az login
```

下記で実行

```
$ go  run . -h {REDIS_ADDRESS:port} -u {objectid で}  | tee log.txt
```

## 3. 解説

ここでは、go-redis の、[Options](https://pkg.go.dev/github.com/redis/go-redis/v9@v9.5.1#Options)、CredentialsProvider をAAD認証で実装しています。
Tokenの所得は、[Azure Identity Client Module for Go](https://pkg.go.dev/github.com/Azure/azure-sdk-for-go/sdk/azidentity@v1.5.2)(以下azidentity)を使っています。
DefaultAzureCredential を使っているので、コマンドラインから実行した場合は、AzureCLICredential が実行されているはず（多くの場合）です。
azidentityは、memory 上の token cacheがデフォルトで実装されるので、最初のアクセス以降のオーバーヘッドは少ないはずです。手元で確認したところ、token のexpireは1時間でした。



## 4. Token の更新

シンプルにトークンの使いまわしと、更新を確認
`2024/05/23 09:35:23` が最初で、1時間後に、`expires=2024-05-23T10:46:30.000+09:00`。tokenのhashは、`945b57`。`10:35:29` までは同じtoken、`10:50:31`に`140bf0`、expiresが`11:58:31`に更新。
それっぽい。redis clientは同じのを使っているので、接続は使いまわしているはずで、認証(AUTHコマンド？)を毎回送る必要はないはず。AUTHコマンドを実行するときに、、CredentialsProviderを呼んでるのではないのか？と思うので、後で確認する。
VM上で、MSI使うと、expires は一日らしい。確かに、az cli 環境よりMSIの方がセキュアなのかもしれない。

```sh
$ go  run . -h $REDIS_ADDRESS -u f05421d4-465e-4766-9d70-xxxxxxxxxxxx  | tee log.txt
2024/05/23 09:35:22 INFO getting credentials
2024/05/23 09:35:23 INFO token username=f05421d4-465e-4766-9d70-xxxxxxxxxxxx token=945b579debc891f139747f50caee915fa89730a0 expires=2024-05-23T10:46:30.000+09:00
2024/05/23 09:35:23 INFO ping response=PONG
2024/05/23 09:50:23 INFO getting credentials
2024/05/23 09:50:25 INFO token username=f05421d4-465e-4766-9d70-xxxxxxxxxxxx token=945b579debc891f139747f50caee915fa89730a0 expires=2024-05-23T10:46:30.000+09:00
2024/05/23 09:50:25 INFO ping response=PONG
2024/05/23 10:05:25 INFO getting credentials
2024/05/23 10:05:26 INFO token username=f05421d4-465e-4766-9d70-xxxxxxxxxxxx token=945b579debc891f139747f50caee915fa89730a0 expires=2024-05-23T10:46:30.000+09:00
2024/05/23 10:05:26 INFO ping response=PONG
2024/05/23 10:20:26 INFO getting credentials
2024/05/23 10:20:28 INFO token username=f05421d4-465e-4766-9d70-xxxxxxxxxxxx token=945b579debc891f139747f50caee915fa89730a0 expires=2024-05-23T10:46:30.000+09:00
2024/05/23 10:20:28 INFO ping response=PONG
2024/05/23 10:35:28 INFO getting credentials
2024/05/23 10:35:29 INFO token username=f05421d4-465e-4766-9d70-xxxxxxxxxxxx token=945b579debc891f139747f50caee915fa89730a0 expires=2024-05-23T10:46:30.000+09:00
2024/05/23 10:35:29 INFO ping response=PONG
2024/05/23 10:50:29 INFO getting credentials
2024/05/23 10:50:31 INFO token username=f05421d4-465e-4766-9d70-xxxxxxxxxxxx token=140bf0d72abc585496101609daeae7f182f3791a expires=2024-05-23T11:58:31.000+09:00
2024/05/23 10:50:31 INFO ping response=PONG
```

