# Структура проекта
- `\cmd\app`: инициализация приложения
- `\internal`: основные пакеты приложения
    - `\adapter`: адаптеры, адаптируют сущности приложения и внешних источников
        - `\db\postgres`: репозиторий PostgreSQL
    - `\composites`: композиты
    - `\config`: конфигурация приложения
    - `\controller\http`: HTTP контроллер, реализует логику ответов на запросы пользователя
        - `\dto`: DTO объекты контроллера для транспортировки данных в бизнес логику
        - `\v1`: handlers v1
    - `\domain`: ядро приложения
        - `\dto`: DTO объекты доменной области для транспортировки данных между слоями домена и репозитория
        - `\entity`: доменные сущности приложения (модели)
        - `\service`: сервисы, оперируют доменными сущностями, репозиториями (технический слой)
        - `\usecase`: реализация бизнес логики приложения
    - `\location`: определение местоположения беспроводного клиента
    - `\router`: инициализация роутера, маршрутизация
    - `\middleware`: промежуточный слой
    - `\migrations`: миграции
- `\pkg`: внешние зависимости
    - `\client\postgres`: подключение к базе данных
    - `\httperrors`: пользовательские ошибки
    - `\logger`: пользовательский логгер
    - `\utils`: вспомогательный код
- `\plugins`: плагины
- `\public`: публичная папка сервера
- `example.env`: шаблон файла переменных окружения
- `main.go`: точка входа в приложение

# Project location-backend

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## Before build and run
Create .env file and copy example.env content into it. `jwt secret` in `.env` must be equal `token_hmac_secret_key` in `centrifugo_config.json`

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

run backend containers
```bash
make docker-run
```

shutdown backend containers
```bash
make docker-down
```

build and run backend containers
```bash
make docker-build
```

rebuild the app container
```bash
make docker-rebuild-app
```

cleanup all docker data
```bash
make docker-clean
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```