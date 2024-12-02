FROM gomicro/goose:latest

# Устанавливаем рабочую директорию для миграций
WORKDIR /migrations

# Копируем файлы миграций в контейнер
COPY . .

# Копируем entrypoint.sh
COPY ./entrypoint.sh /entrypoint.sh
# Убедитесь, что путь к entrypoint.sh правильный

# Делаем entrypoint.sh исполняемым
RUN chmod +x /entrypoint.sh

# Указываем команду по умолчанию для запуска контейнера
ENTRYPOINT ["/entrypoint.sh"]
