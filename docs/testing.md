# Тестирование

## 1. Цели и стратегия

- **Unit-тесты (Go):** быстрые проверки домена, CLI, презентации, конвертации proto↔domain, сессии на подставном транспорте, gRPC-транспорта на локальном `net.Listener`.
- **Интеграционные тесты (Python + pytest):** сборка бинарника `bin/chat`, проверка флагов и сквозного сценария «слушатель + клиент» в отдельных процессах.

## 2. Как запускать

| Команда | Назначение |
|---------|------------|
| `go test -race ./...` | Все Go-тесты с детектором гонок |
| `go test -race ./internal/grpcchat -count=5` | Повторить пакет с gRPC несколько раз (флейки) |
| `python3 -m pytest tests/ -v` | CLI и интеграция (нужен `go` в PATH для сборки в `conftest`) |

## 3. Реализованные тесты

### 3.1 Go

| Пакет / файл | Что проверяется |
|--------------|-----------------|
| `internal/domain` — `message_test.go` | `Message.Validate`: валидное сообщение, пустой отправитель, тело длиннее лимита рун |
| `internal/cli` — `config_test.go` | `Parse`: флаги `-name`, `-listen`, `-connect`; ошибка без имени; невалидный `-connect`; невалидный `-listen`; значения по умолчанию (`:50051`, режим listen) |
| `internal/ui` — `presenter_test.go` | `Presenter.Show`: вывод содержит имя, тело, метку времени в ожидаемом формате |
| `internal/grpcchat` — `convert_test.go` | `domainToProto` / `fromProto`: round-trip; `nil`; невалидный домен; слишком длинное тело |
| `internal/grpcchat` — `listen_test.go` | `listenOn` + `Dial`: сообщение клиент→сервер и ответ сервер→клиент по одному стриму |
| `internal/session` — `session_test.go` | `Run` с подставным `pipeTransport`: `/quit`, `/exit`, отправка строки и отображение в выводе (эхо через общий канал) |

### 3.2 Python (`tests/`)

| Файл | Что проверяется |
|------|-----------------|
| `test_cli.py` | `-h` завершается с кодом 0; без `-name` — ненулевой код и упоминание имени; невалидный `-connect` |
| `test_integration.py` | Слушатель на свободном порту, клиент подключается; сообщение из stdin Alice доходит до stdout Bob |

Фикстура `chat_binary` в `conftest.py` один раз собирает `./cmd/chat` в `bin/chat`.
