# Effective Mobile Test Task

**REST-сервис для агрегации данных об онлайн подписках пользователей**

### Предварительные требования
- Docker и Docker Compose
- Git

### Установка и запуск
git clone https://github.com/CAATHARSIS/Effective-Mobile-Test

cd Effective-Mobile-Test

Создать .env файл, примерное содержание:

`DB_PORT=5432`

`DB_USER=postgres`

`DB_PASSWORD=postgres`

`DB_NAME=Effective-Mobile-Test`

docker-compose up --build

### Проверка работы
Приложение: http://localhost:8080

API документация: http://localhost:8080/swagger