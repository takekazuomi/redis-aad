package main

import (
	"context"
	"crypto/sha1"
	"crypto/tls"
	"fmt"
	"log/slog"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/pflag"
)

func run(address string, interval time.Duration) {
	options := &redis.Options{
		Addr:                address,
		CredentialsProvider: CredentialsProvider,
		TLSConfig:           &tls.Config{MinVersion: tls.VersionTLS12},
	}

	slog.Info("options", "options", options)

	cli := redis.NewClient(options)

	for {

		s, err := cli.Ping(context.Background()).Result()
		if err != nil {
			slog.Error("could not establish connection", "error", err)
			return
		}

		slog.Info("ping", "response", s)
		time.Sleep(interval)
	}
}

var (

	// CredentialsProvider に渡すために、グローバル変数に保持
	username string
)

// https://learn.microsoft.com/en-us/azure/azure-cache-for-redis/cache-azure-active-directory-for-authentication
func CredentialsProvider() (string, string) {
	slog.Info("getting credentials")

	// https://learn.microsoft.com/en-us/azure/azure-cache-for-redis/cache-azure-active-directory-for-authentication#microsoft-entra-client-workflow
	// https://redis.azure.com/.default or acca5fbb-b7e4-4009-81f1-37e38fd66d78/.default
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		slog.Error("could not get credentials", "error", err)
	}

	token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: []string{"https://redis.azure.com/.default"},
	})
	if err != nil {
		slog.Error("could not get token", "error", err)

	}

	// tokenがexpiresした時の動作確認のため、tokenのhashとexpireを表示
	hash := fmt.Sprintf("%x", sha1.Sum([]byte(token.Token)))
	slog.Info("token", "username", username, "token", hash, "expires", token.ExpiresOn.Local())

	return username, token.Token
}

func main() {
	var address string
	var interval time.Duration

	pflag.StringVarP(&address, "redis-address", "h", "", "redis address")
	pflag.StringVarP(&username, "redis-username", "u", "", "redis username")
	pflag.DurationVarP(&interval, "interval", "i", time.Second*60*15, "interval to ping redis server")
	help := pflag.BoolP("help", "", false, "show help")

	pflag.Parse()
	if help != nil && *help {
		pflag.Usage()
		return
	}

	run(address, interval)
}
