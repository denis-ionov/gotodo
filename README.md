# gotodo
## Инструкция по запуску
1. Установите зависимости:
```
go mod tidy
```
2. Запустите сервер:

```
go run cmd/main.go
```

## Тестирование API
Используйте Postman для тестирования
1. Регистрация пользователя
```   
POST http://localhost:8080/api/register
Content-Type: application/json
{
    "name": "Денис Денисов",
    "email": "denis@example.com",
    "password": "123456"
}
```
2. Вход в систему
```   
POST http://localhost:8080/api/login
Content-Type: application/json
{
    "email": "denis@example.com",
    "password": "123456"
}
```
3. Создание задачи с токеном
```
POST http://localhost:8080/api/tasks
Authorization: Bearer <токен полученный на шаге 2>
Content-Type: application/json
{
    "title": "Первая задача",
    "description": "Описание задачи"
}
```
4. Получение всех задач
```
GET http://localhost:8080/api/tasks
Authorization: Bearer <токен полученный на шаге 2>
```
## Особенности реализации
1. Разделение на слои (handlers, service, repository)
2. JWT аутентификация
3. Роуты, защищенные middleware
4. Использование Gin binding для валидации входных данных
5. Хеширование паролей с помощью bcrypt
6. Использование UUID для уникальных идентификаторов
## TODO
1. Перевести на docker контейнеры
2. Заменить in-memory репозитории на БД PostgreSQL
