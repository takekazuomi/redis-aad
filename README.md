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


