package main

import (
	"fmt"
	"log"

	cfgldrlib "github.com/GAFIKART/cfg-ldr/lib"
	"github.com/hashicorp/vault/api"
)

// Config структура конфигурации приложения
// Для Vault провайдера используются теги cfgldr:"val=..." для указания ключей в Vault
// Теги val= определяют имена ключей в секретах Vault
type Config struct {
	Server struct {
		Host string `cfgldr:"val=host"`
		Port int    `cfgldr:"val=port"`
	}

	Database struct {
		Host     string `cfgldr:"val=host"`
		Port     int    `cfgldr:"val=port"`
		Name     string `cfgldr:"val=name"`
		User     string `cfgldr:"val=user"`
		Password string `cfgldr:"val=password"`
	}

	Logging struct {
		Level string `cfgldr:"val=level"`
		File  string `cfgldr:"val=file"`
	}

	Features struct {
		CacheEnabled bool `cfgldr:"val=cache_enabled"`
		DebugMode    bool `cfgldr:"val=debug_mode"`
	}

	Lists struct {
		AllowedHosts []string `cfgldr:"val=allowed_hosts"`
		Ports        []int    `cfgldr:"val=ports"`
		Features     []string `cfgldr:"val=features"`
	}
}

func main() {
	// Создаем клиент Vault
	vaultClient, err := api.NewClient(&api.Config{
		Address: "http://localhost:8200", // Адрес Vault сервера
	})
	if err != nil {
		log.Fatalf("Ошибка создания клиента Vault: %v", err)
	}

	// Устанавливаем токен для аутентификации (в продакшене используйте более безопасные методы)
	vaultClient.SetToken("your-vault-token")

	// Создаем параметры для загрузчика конфигурации
	kvName := "app-config" // Имя KV секрета в Vault
	params := &cfgldrlib.ConfigLoaders{
		VaultParams: &cfgldrlib.VaultParamsT{
			KvName:      &kvName,
			VaultClient: vaultClient,
		},
	}

	// Загружаем конфигурацию из Vault
	config, err := cfgldrlib.LoadConfig[Config](params)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации из Vault: %v", err)
	}

	// Выводим загруженную конфигурацию
	fmt.Println("=== Загруженная конфигурация (Vault) ===")
	fmt.Printf("Сервер: %s:%d\n", config.Server.Host, config.Server.Port)
	fmt.Printf("База данных: %s:%d/%s (пользователь: %s)\n",
		config.Database.Host, config.Database.Port, config.Database.Name, config.Database.User)
	fmt.Printf("Логирование: уровень=%s, файл=%s\n", config.Logging.Level, config.Logging.File)
	fmt.Printf("Функции: кэш=%v, отладка=%v\n", config.Features.CacheEnabled, config.Features.DebugMode)

	// Выводим слайсы
	fmt.Printf("Разрешенные хосты: %v\n", config.Lists.AllowedHosts)
	fmt.Printf("Порты: %v\n", config.Lists.Ports)
	fmt.Printf("Функции: %v\n", config.Lists.Features)
}
