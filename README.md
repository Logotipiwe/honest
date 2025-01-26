<h1>How Are You Really?</h1>

# Локальный запуск

- Поднять конфиг сервер
- Запустить go конфиг

# Swagger

Генерируется командой ```swag```

Установка команды swag https://github.com/swaggo/swag?tab=readme-ov-file#getting-started

```shell
go install github.com/swaggo/swag/cmd/swag@latest
```

для macOS также добавление в PATH:
```shell
echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.zshrc
source ~/.zshrc
```

Перед коммитом - обновить доку
```shell
swag i -d ./src/cmd,./src/internal
```

Прод - https://logotipiwe.ru/haur/swagger/index.html

# TODO
- сваггер
- скрытые колоды