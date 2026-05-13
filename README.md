# dns-manager

Клиент-серверное приложение для управления DNS-серверами через `/etc/resolv.conf`.

## Конфигурация

Настройки читаются из YAML-файла. Любое поле можно переопределить переменной окружения.

CLI-флаги (`--resolv`, `--backup`, `--server`, `--config`) имеют приоритет над config и env.

## Запуск

### Сервер

```bash
make run-server

# или напрямую:
go run ./server/cmd/server --resolv /tmp/test_resolv.conf

```

### Клиент

```bash
make run-client CMD="add 8.8.8.8"
make run-client CMD="list"
make run-client CMD="remove 8.8.8.8"
make run-client CMD="help"

# или напрямую:
go run ./client/cmd/client add 8.8.8.8
go run ./client/cmd/client list
```

## Тесты

```bash
make test
```
