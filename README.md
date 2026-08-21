# Parterre

[![Check](https://github.com/iamroockie/parterre/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/iamroockie/parterre/actions/workflows/ci.yml)

Платформа продажи билетов на события с выбором **конкретного места** в зале.

Ключевое обещание системы:

> Место продаётся **ровно один раз**, и это доказано воспроизводимым
> нагрузочным тестом, а не заявлено в README.

## Статус

> **Приложение на стадии разработки**

## Стек

- Go 1.27
- PostgreSQL
- [chi](https://github.com/go-chi/chi) — маршрутизация
- [Task](https://taskfile.dev) — единая точка входа для команд
- [golangci-lint](https://golangci-lint.run) — линтеры

## Быстрый старт

Нужны [Go 1.27+](https://go.dev/dl/) и [Task](https://taskfile.dev/installation/).

```sh
cp .env.example .env
task run:api
```

Проверить, что поднялось:

```sh
curl -i localhost:8080/healthz
```

Остановить — `Ctrl+C`: сервер перестаёт принимать новые запросы, дорабатывает
текущие и только потом выходит.

## Команды

| Команда              | Что делает                                    |
| -------------------- | --------------------------------------------- |
| `task run:api`       | запустить API                                  |
| `task build:api`     | собрать бинарник в `bin/`                      |
| `task lint`          | линтеры                                        |
| `task test`          | юнит-тесты (без интеграционных)                |
| `task check`         | lint + test                                    |
| `task format`        | отформатировать код                            |
| `task coverage:func` | покрытие тестами в процентах                   |
| `task coverage:html` | покрытие тестами отчётом в `coverage.html`     |

`task --list` покажет все.

## Конфигурация

Приложение читает настройки **только из переменных окружения** — файла
конфигурации нет. `.env` нужен исключительно для локальной разработки: Task
подхватывает его сам, а в git он не попадает. Полный список переменных со
значениями по умолчанию — в [`.env.example`](.env.example).

`APP_ENV` заодно переключает формат логов: `local` — читаемый текст,
`prod` — JSON.

## Структура

```
cmd/api/              точка входа
internal/
  platform/           инфраструктура, не знающая о предметной области
  config/             разбор переменных окружения
  logger/             slog + передача логгера через context
  httpx/              роутер, middleware
migrations/           миграции БД
tests/                интеграционные тесты
docs/                 документация
```

## Документация

- [`Issues`](https://github.com/iamroockie/parterre/issues) — задачи с критериями приёмки,
  сгруппированные по эпикам через milestones
- [`docs/adr/`](docs/adr/) — архитектурные решения и их обоснование
- [`docs/architecture.md`](docs/architecture.md) — что строим и почему именно так
- [`docs/git-workflow.md`](docs/git-workflow.md) — ветки, коммиты, PR
