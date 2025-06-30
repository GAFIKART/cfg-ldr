# cfg-ldr

Библиотека для загрузки конфигурации в Go приложениях с поддержкой YAML файлов и HashiCorp Vault.

## Описание

`cfg-ldr` - это универсальная библиотека для загрузки конфигурации в Go приложениях. Библиотека поддерживает два основных провайдера конфигурации:

- **YAML файлы** - для локальной разработки и простых конфигураций
- **HashiCorp Vault** - для безопасного хранения секретов и конфигурации в продакшене

## Возможности

- ✅ Загрузка конфигурации из YAML файлов
- ✅ Загрузка конфигурации из HashiCorp Vault
- ✅ Автоматическое определение типа провайдера
- ✅ Поддержка вложенных структур
- ✅ Автоматическое преобразование типов данных
- ✅ Валидация параметров
- ✅ Гибкая система тегов для настройки маппинга
- ✅ Обработка ошибок с подробными сообщениями

## Установка

```bash
go get github.com/GAFIKART/cfg-ldr
```

## Быстрый старт

### Загрузка из YAML файла

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    cfgldrlib "github.com/GAFIKART/cfg-ldr/lib"
)

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
}

func main() {
    // Читаем YAML файл
    configYaml, err := os.ReadFile("config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    configYamlStr := string(configYaml)
    params := &cfgldrlib.ConfigLoaders{
        ConfigYml: &configYamlStr,
    }
    
    // Загружаем конфигурацию
    config, err := cfgldrlib.LoadConfig[Config](params)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Сервер: %s:%d\n", config.Server.Host, config.Server.Port)
}
```

### Загрузка из Vault

```go
package main

import (
    "fmt"
    "log"
    
    cfgldrlib "github.com/GAFIKART/cfg-ldr/lib"
    "github.com/hashicorp/vault/api"
)

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
}

func main() {
    // Создаем клиент Vault
    vaultClient, err := api.NewClient(&api.Config{
        Address: "http://localhost:8200",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    vaultClient.SetToken("your-vault-token")
    
    kvName := "app-config"
    params := &cfgldrlib.ConfigLoaders{
        VaultParams: &cfgldrlib.VaultParamsT{
            KvName:      &kvName,
            VaultClient: vaultClient,
        },
    }
    
    // Загружаем конфигурацию
    config, err := cfgldrlib.LoadConfig[Config](params)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Сервер: %s:%d\n", config.Server.Host, config.Server.Port)
}
```

## API

### Основные типы

#### ConfigLoaders

Структура для настройки загрузчика конфигурации:

```go
type ConfigLoaders struct {
    ConfigYml      *string        // Содержимое YAML файла
    ConfigProvider *ConfigProvider // Явное указание провайдера
    VaultParams    *VaultParamsT  // Параметры для Vault
}
```

#### VaultParamsT

Параметры для работы с Vault:

```go
type VaultParamsT struct {
    KvName      *string     // Имя KV секрета в Vault
    VaultClient *api.Client // Клиент Vault
}
```

#### ConfigProvider

Типы провайдеров конфигурации:

```go
const (
    ConfigProviderYml   ConfigProvider = "yml"
    ConfigProviderVault ConfigProvider = "vault"
)
```

### Основная функция

#### LoadConfig[T]

Загружает конфигурацию в структуру указанного типа:

```go
func LoadConfig[T any](params *ConfigLoaders) (*T, error)
```

## Система тегов

### YAML провайдер

Для YAML провайдера используются стандартные теги `yaml`:

```go
type Config struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}
```

### Vault провайдер

Для Vault провайдера используются теги `cfgldr` с параметром `val=`:

```go
type Config struct {
    Host string `cfgldr:"val=host"`
    Port int    `cfgldr:"val=port"`
}
```

#### Специальные теги

- `cfgldr:"-"` - исключить поле из обработки
- `cfgldr:"val=key_name"` - указать имя ключа в Vault

## Поддерживаемые типы данных

Библиотека поддерживает следующие типы данных:

- `string` - строки
- `int`, `int8`, `int16`, `int32`, `int64` - целые числа
- `float32`, `float64` - числа с плавающей точкой
- `bool` - логические значения

Автоматическое преобразование типов:
- Числа могут быть загружены из строк
- Булевы значения поддерживают строки "true"/"false", "1"/"0"

## Структура Vault секретов

Для корректной работы с Vault секреты должны быть организованы следующим образом:

```
app-config/
├── server/
│   ├── host
│   └── port
├── database/
│   ├── host
│   ├── port
│   ├── name
│   ├── user
│   └── password
└── logging/
    ├── level
    └── file
```

Каждый секрет должен содержать данные в формате KV v2.

## Примеры

Полные примеры использования находятся в директории `example/`:

- `example/yml/` - пример работы с YAML файлами
- `example/vault/` - пример работы с HashiCorp Vault

### Пример YAML конфигурации

```yaml
server:
  host: "0.0.0.0"
  port: 9090

database:
  host: "db.example.com"
  port: 5432
  name: "production_db"
  user: "app_user"
  password: "secret_password_123"

logging:
  level: "debug"
  file: "production.log"

features:
  cache_enabled: true
  debug_mode: false
```

## Обработка ошибок

Библиотека предоставляет подробные сообщения об ошибках:

- Валидация параметров
- Ошибки чтения файлов
- Ошибки подключения к Vault
- Ошибки парсинга данных
- Ошибки преобразования типов

## Требования

- Go 1.24.0 или выше
- HashiCorp Vault API (для работы с Vault)

## Зависимости

- `github.com/hashicorp/vault/api` - для работы с HashiCorp Vault
- `gopkg.in/yaml.v3` - для парсинга YAML файлов

