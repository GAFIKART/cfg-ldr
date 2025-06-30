package cfgldrlib

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func loadConfigFromYml[T any](configYML *string) (*T, error) {
	readerConfig := strings.NewReader(*configYML)
	decoder := yaml.NewDecoder(readerConfig)

	var config T
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func loadConfigFromVault[T any](vaultParams *VaultParamsT) (*T, error) {
	var config T
	err := fillStructFromVault(&config, vaultParams.VaultClient, *vaultParams.KvName)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
