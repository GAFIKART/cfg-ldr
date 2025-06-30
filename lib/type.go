package cfgldrlib

import (
	"github.com/hashicorp/vault/api"
)

const TagName = "cfgldr"
const ValTagParam = "val="

type ConfigProvider string

const (
	ConfigProviderYml   ConfigProvider = "yml"
	ConfigProviderVault ConfigProvider = "vault"
)

type ConfigLoaders struct {
	ConfigYml      *string
	ConfigProvider *ConfigProvider
	VaultParams    *VaultParamsT
}

type VaultParamsT struct {
	KvName      *string
	VaultClient *api.Client
}
