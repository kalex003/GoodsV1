#!/bin/bash

# Если что-то идет не так, выводим сообщение об ошибке и останавливаем скрипт
set -e

# Проверяем, заданы ли необходимые переменные окружения
if [ -z "$DATABASE_URL" ]; then
  echo "Error: DATABASE_URL is not set."
  exit 1
fi

# Логируем информацию о миграции
echo "Running Goose migrations..."

# Запуск миграций с Goose
goose -dir /migrations postgres "$DATABASE_URL" up

# Если миграции выполнены успешно, выводим сообщение
echo "Migrations completed successfully."
