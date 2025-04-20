# go/xlog
[![main.yml](https://github.com/anticrew/log/actions/workflows/main.yml/badge.svg)](https://github.com/anticrew/log/actions/workflows/main.yml)  
_Высокопроизводительный и минималистичный Go logger_

## Почему xlog?
Единый интерфейс без привязки к реальному фреймворку. Используйте фиксированный универсальный API, а реализации 
переключайте с помощью тегов build tags:
- `anticrew_log_zap` ➜ `uber/zap`
- `anticrew_log_slog` ➜ `go/slog`

## Design
- `Arg` - именованные типизированные аргументы   
Аргументы позволяют передавать с логом дополнительные данные. Для максимальной производительности рекомендуется 
использовать типизированные атрибуты: `xlog.String`, `xlog.Bool` и другие.
- `Context-Logger` - хранение `Logger` в `context.Context`  
Методы `GetContextLogger` и `SetContextLogger` позволяют получать и устанавливать `Logger` в `context.Context`. Каждый 
экземпляр `context.Context` может содержать только один экземпляр `Logger`, при создании дочерних `context.Context` они 
продолжают поддержку получения `Logger`, как и любых других значений, добавленных через `context.WithValue`
- `Context-Arg` - хранение набора `Arg` в `context.Context`
Методы `GetContextArg` и `SetContextArg` позволяют получать и устанавливать `Arg` в `context.Context` по указанному ключу.
Методы `GetContextArgs`, `AddContextArgs` и `SetContextArgs` позволяют получать и устанавливать целый набор `Arg` в 
`context.Context` по фиксированному внутреннему ключу (используется ключ-структура `argsKey{}`).


# Install
```bash
go get github.com/anticrew/log
```

# Development
В командах ниже `{{ driver }}` необходимо заменить на корректный тег для выбора драйвера из указанных выше.

## Lint
```bash
golangci-lint run --build-tags {{ driver }}
```

## Test
Целевое покрытие тестами - 80%
```bash
go test -tags {{ driver }} -v ./...
```

# Perfomance
Бенчмарк доступен в папке benchmark, запускается 3 группы бенчмарков по 3 логгерам:
- `zap`
- `slog`
- `xlog`

Для каждого логгера запускается 4 бенчмарка:
- `logfmt` - запись логов в LogFmt 
- `text` - запись логов в стандартном для логгера текстовом формате
- `json` - запись логов в JSON
- `json-any` - запись логов в JSON с маршалингом кастомной структуры среди атрибутов

Результаты бенчмарков ниже достигнуты на следующем сетапе:
```
goos: windows
goarch: amd64
pkg: github.com/anticrew/log/benchmark
cpu: AMD Ryzen 7 5700X 8-Core Processor
```

## LogFmt
| Logger | Driver | N      | ns/op | B/op | allocs/op |
|--------|--------|--------|-------|------|-----------|
| slog   | slog   | 263697 | 4517  | 1322 | 19        |
| zap    | zap    | 298497 | 4013  | 947  | 8         |
| xlog   | slog   | 322170 | 3674  | 1117 | 13        |
| xlog   | zap    | 717826 | 1686  | 216  | 8         |

## Text
| Logger | Driver | N      | ns/op | B/op | allocs/op |
|--------|--------|--------|-------|------|-----------|
| slog   | slog   | 263697 | 4517  | 1322 | 19        |
| zap    | zap    | 303142 | 4000  | 1069 | 17        |
| xlog   | slog   | 330854 | 3596  | 1116 | 13        |
| xlog   | zap    | 295057 | 4020  | 1377 | 23        |

## JSON
| Logger | Driver | N      | ns/op | B/op | allocs/op |
|--------|--------|--------|-------|------|-----------|
| slog   | slog   | 283357 | 4231  | 1849 | 20        |
| zap    | zap    | 339984 | 3536  | 947  | 8         |
| xlog   | slog   | 372331 | 3255  | 1117 | 13        |
| xlog   | zap    | 302076 | 3855  | 1037 | 14        |

## JSON-Any
| Logger | Driver | N      | ns/op | B/op | allocs/op |
|--------|--------|--------|-------|------|-----------|
| slog   | slog   | 206641 | 5753  | 2668 | 36        |
| zap    | zap    | 244671 | 4895  | 1735 | 19        |
| xlog   | slog   | 253510 | 4861  | 1712 | 25        |
| xlog   | zap    | 202628 | 5948  | 1891 | 27        |