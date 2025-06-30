package cfgldrlib

import "errors"

func validateParams(params *ConfigLoaders) error {
	if params == nil {
		return errors.New("params is nil")
	}
	if params.ConfigYml == nil && params.VaultParams == nil {
		return errors.New("config yml or vault client is nil")
	}

	// Проверяем VaultParams только если они предоставлены
	if params.VaultParams != nil {
		if params.VaultParams.KvName == nil {
			return errors.New("vault kv name is nil")
		}
		if params.VaultParams.VaultClient == nil {
			return errors.New("vault client is nil")
		}
	}

	return nil
}

func loadConfigProvider(params *ConfigLoaders) (ConfigProvider, error) {
	if params.ConfigProvider != nil {
		return *params.ConfigProvider, nil
	}

	if params.ConfigYml != nil {
		return ConfigProviderYml, nil
	}

	if params.VaultParams != nil && params.VaultParams.VaultClient != nil {
		return ConfigProviderVault, nil
	}

	return "", errors.New("config provider error detected")
}
