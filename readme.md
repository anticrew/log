# go/xlog
_Высокопроизводительный и минималистичный Go logger_

## Почему xlog?
- `zero allocation`  
  Все структуры, часто использующиеся в работе переиспользуются через `sync.Pool` и неуправляемую память.
- `zero dependency`  
  Зависимости в `go.mod` - `zap` для бенчмарка и `testify` для тестов, ни одна зависимость не попадает в итоговую сборку.
- поддержка `middleware`  
  `Middleware` может кастомизировать лог перед его записью, например, интерполировать атрибуты в строку сообщения или 
  изменять ключи, типы, значения атрибутов
- запись внутренних ошибок  
  Внутренние ошибки `xlog` не опускаются, а записываются следом за логом, при обработке которого они произошли. 
  Ошибки, связанные с работой `Middleware` или получением стека вызовов не влияют на запись лога.

## Design
- `xlog.Attr` - именованные типизированные атрибуты   
  Атрибуты позволяют передавать с логом дополнительные данные. Для максимальной производительности рекомендуется
  использовать типизированные атрибуты `xlog.String`, `xlog.Bool` и другие
- Цепочка ответственности `Log > Middleware > Engine > Marshaler`
   - `Log` обрабатывает записи и отсеивает ненужные
   - `Middleware` насыщает дополнительными данными и изменяет атрибуты
   - `Engine` буферизирует запись и превращает ее в поток байт через `Marshaler`, а затем записывает в вывод.

# Install
```bash
go get github.com/anticrew/log
```

# Development
## Lint

```bash
golangci-lint run
```

## Test
Целевое покрытие тестами - 80%
```bash
go test -v ./...
```

# Perfomance
Бенчмарк доступен в папке benchmark, запускается 3 группы бенчмарков по 3 логгерам:
- `zap`
- `slog`
- `xlog`

Для каждого логгера запускается 3 бенчмарка:
- `logfmt` - запись логов в LogFmt 
- `json` - запись логов в JSON
- `json-any` - запись логов в JSON с маршалингом кастомной структуры среди атрибутов

Для `xlog` также запускается бенчмарк `pretty` - форматированный лог в стиле `npm`

Результаты бенчмарков ниже достигнуты на следующем сетапе:
```
goos: windows
goarch: amd64
pkg: github.com/anticrew/log/benchmark
cpu: 13th Gen Intel(R) Core(TM) i9-13900H
```

## LogFmt
```
Benchmark_Slog
Benchmark_Slog/logfmt
Benchmark_Slog/logfmt-20         	  186484	      6393 ns/op	    1371 B/op	      19 allocs/op

Benchmark_Xlog
Benchmark_Xlog/logfmt
Benchmark_Xlog/logfmt-20         	  272448	      4308 ns/op	    1293 B/op	      12 allocs/op

Benchmark_Zap
Benchmark_Zap/logfmt
Benchmark_Zap/logfmt-20          	  216344	      5457 ns/op	     924 B/op	       8 allocs/op
```

## Json
```
Benchmark_Slog
Benchmark_Slog/json
Benchmark_Slog/json-20           	  205005	      6119 ns/op	    1849 B/op	      20 allocs/op

Benchmark_Xlog
Benchmark_Xlog/json
Benchmark_Xlog/json-20           	  220084	      5342 ns/op	    1295 B/op	      12 allocs/op

Benchmark_Zap
Benchmark_Zap/json
Benchmark_Zap/json-20            	  256888	      5888 ns/op	     923 B/op	       8 allocs/op
```

## Json Any
```
Benchmark_Slog
Benchmark_Slog/json-any
Benchmark_Slog/json-any-20       	  126555	      9697 ns/op	    2968 B/op	      42 allocs/op

Benchmark_Xlog
Benchmark_Xlog/json-any
Benchmark_Xlog/json-any-20       	  121662	      9736 ns/op	    2087 B/op	      29 allocs/op

Benchmark_Zap
Benchmark_Zap/json-any
Benchmark_Zap/json-any-20        	   89996	     12906 ns/op	    1963 B/op	      25 allocs/op
```

## Pretty (xlog only)
```
Benchmark_Xlog
Benchmark_Xlog/pretty
Benchmark_Xlog/pretty-20         	  241713	      4939 ns/op	    1294 B/op	      12 allocs/op
```

### Пример pretty-лога
```
[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:19 main.main
 TRACE  trace message
  instance: default

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:20 main.main
 TRACE  trace message
  instance: local

[ 2025-06-04 19:45:33 ]/tmp/log/example/main.go:22 main.main
 DEBUG  debug message
  instance: default

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:23 main.main
 DEBUG  debug message
  instance: local

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:25 main.main
 INFO  info message
  instance: default

[ 2025-06-04 19:45:33 /tmp/log/example/main.go:26 main.main
 DEBUG  info message
  instance: local

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:28 main.main
 WARN  warning message
  instance: default

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:29 main.main
 WARN  warning message
  instance: local

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:31 main.main
 DEBUG  debug message with "interpolation" and attr
  instance: default
  key: interpolation

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:32 main.main
 DEBUG  debug message with "interpolation" and attr
  instance: local
  key: interpolation

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:34 main.main
 DEBUG  debug message with attr only
  instance: default
  key: interpolation

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:35 main.main
 DEBUG  debug message with attr only
  instance: local
  key: interpolation

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:37 main.main
 WARN  warn message with interpolation from "default" instance
  instance: default

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:38 main.main
 WARN  warn message with interpolation from "local" instance
  instance: local

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:41 main.main
 ERROR  oh, error?
  error: not found
  instance: default

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:41 main.main
 ERROR  not found
  error: not found
  instance: default

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:42 main.main
 ERROR  oh, error?
  error: not found
  instance: local

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:42 main.main
 ERROR  not found
  error: not found
  instance: local

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:44 main.main
 FATAL  oh, fail!
  error: not found
  instance: default

[ 2025-06-04 19:45:33 ]	/tmp/log/example/main.go:44 main.main
 ERROR  not found
  error: not found
  instance: default
```