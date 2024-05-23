package credentialsprovider

type CredentialsProvider_ interface {
	CredentialsProvider() (username string, password string)
}

type ProviderName string

const (
	ProviderNameAzure ProviderName = "AzureAD"
	ProtocolNameKey   ProviderName = "Key"
)

var provider CredentialsProvider_

func Init(p CredentialsProvider_) {
	provider = p
}

var hosts = map[string]CredentialsProvider_{}

func New(address string, p CredentialsProvider_) {
	hosts[address] = p
}

func Provider() (string, string) {
	return provider.CredentialsProvider()
}
