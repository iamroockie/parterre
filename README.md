# Parterre

[![Check](https://github.com/iamroockie/parterre/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/iamroockie/parterre/actions/workflows/ci.yml)

Платформа продажи билетов на события с выбором **конкретного места** в зале.

Ключевое требование к системе:

> Место продаётся **ровно один раз**, и это доказано воспроизводимым
> нагрузочным тестом.

## Статус

> **Приложение на стадии разработки**

## Локальный запуск

Нужны Docker и [Task](https://taskfile.dev); для тестов и запуска api с хоста — Go 1.27.

```sh
cp .env.example .env
task up          # postgres и api в docker
task migrate:up  # накатывает миграции
curl localhost:8080/healthz
```

`task up` пересобирает образ приложения и возвращает управление только после того,
как оба контейнера прошли проверку здоровья. Миграции накатываются отдельной
командой. Данные лежат в томе docker и переживают `task down`.

### api процессом на хосте

Второй способ, для цикла «правка — перезапуск»: в docker остаётся только postgres,
приложение работает на хосте.

```sh
task up:db
task run:api  # http://localhost:8080
```

Перезапуск обходится без пересборки образа.

## Команды

| Команда                     | Что делает                          |
| --------------------------- | ----------------------------------- |
| `task up` / `task down`     | поднять / остановить postgres и api |
| `task up:db`                | поднять только postgres             |
| `task run:api`              | запустить api процессом на хосте    |
| `task test`                 | тесты без интеграционных            |
| `task test:all`             | все тесты, включая интеграционные   |
| `task lint` / `task format` | линтинг / форматирование            |
| `task check`                | линтинг и все тесты                 |
| `task db:psql`              | консоль postgres                    |

### Миграции

| Команда                             | Что делает                     |
| ----------------------------------- | ------------------------------ |
| `task migrate:up`                   | накатить все ненакатанные      |
| `task migrate:down`                 | откатить последнюю             |
| `task migrate:status`               | что накатано                   |
| `task migrate:validate`             | проверить файлы миграций       |
| `task migrate:create -- add_events` | создать `00002_add_events.sql` |

Файлы лежат в `internal/platform/postgres/migrations/sql` и зашиты в бинарник через `embed`.

Интеграционные тесты поднимают отдельный postgres в docker и базу из `task up` не трогают.
