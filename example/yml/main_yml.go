package main

import (
	"fmt"
	"log"
	"os"

	cfgldrlib "github.com/GAFIKART/cfg-ldr/lib"
)

// Config структура конфигурации приложения
// Для YAML провайдера используются стандартные теги yaml:
// Для Vault провайдера используются теги cfgldr:"val=..."
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	} `yaml:"database"`

	Logging struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"logging"`

	Features struct {
		CacheEnabled bool `yaml:"cache_enabled"`
		DebugMode    bool `yaml:"debug_mode"`
	} `yaml:"features"`

	Lists struct {
		AllowedHosts []string `yaml:"allowed_hosts"`
		Ports        []int    `yaml:"ports"`
		Features     []string `yaml:"features"`
	} `yaml:"lists"`
}

func main() {
	// Читаем содержимое файла конфигурации
	configYaml, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Ошибка чтения файла конфигурации: %v", err)
	}

	// Создаем параметры для загрузчика конфигурации
	configYamlStr := string(configYaml)
	params := &cfgldrlib.ConfigLoaders{
		ConfigYml: &configYamlStr,
	}

	// Загружаем конфигурацию
	config, err := cfgldrlib.LoadConfig[Config](params)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Выводим загруженную конфигурацию
	fmt.Println("=== Загруженная конфигурация (YAML) ===")
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
