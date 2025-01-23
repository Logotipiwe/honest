<h1>How Are You Really?</h1>

# Локальный запуск

- Поднять конфиг сервер
- Запустить go конфиг

# Swagger

Установка https://github.com/swaggo/swag?tab=readme-ov-file#getting-started

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

для macOS также:
```shell
echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.zshrc
source ~/.zshrc
```

Перед коммитом - обновить доку
```shell
swag i
```

Прод - https://logotipiwe.ru/haur/swagger/index.html

# TODO
- юнит тесты
- тесты адаптеров - применить там флайвей
- сваггер
- скрытые колоды