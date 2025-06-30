package cfgldrlib

import "errors"

func LoadConfig[T any](params *ConfigLoaders) (*T, error) {
	err := validateParams(params)
	if err != nil {
		return nil, err
	}
	configProvider, err := loadConfigProvider(params)
	if err != nil {
		return nil, err
	}

	if configProvider == ConfigProviderYml {
		return loadConfigFromYml[T](params.ConfigYml)
	} else if configProvider == ConfigProviderVault {
		return loadConfigFromVault[T](params.VaultParams)
	}

	return nil, errors.New("config error detected")
}
