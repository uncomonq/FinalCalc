# 📐 FinalCalc

**FinalCalc** — это веб-приложение для выполнения математических вычислений. Серверная часть написана на Go и использует gRPC, а также предоставляет простой веб-интерфейс и REST API для взаимодействия с пользователем или внешними сервисами.

---

## 🧰 Возможности

- Ввод математических выражений через веб-интерфейс или API
- Обработка выражений через gRPC-сервис
- Расширяемая архитектура с разделением по пакетам
- Поддержка тестирования с помощью Postman

---

## 📦 Зависимости и требования

Для корректной работы проекта на **Windows** необходима установка следующих компонентов:

- [Go 1.20+](https://golang.org/dl/)
- [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) **(обязателен для сборки на Windows, так как используется `cgo`)**
- [Postman](https://www.postman.com/) для ручного тестирования API

---

## 🚀 Установка и запуск

### 1. Клонирование проекта

```bash
git clone https://github.com/uncomonq/FinalCalc.git
cd FinalCalc
```

### 2. Установка зависимостей

```bash
go mod tidy
```
### 3. Запуск

```bash
go run cmd/main.go
```
### Веб интерфейс

http://localhost:8080/api/v1/web

### Тесты

Для запуска тестов нужно выполнить такую команду:
```bash
go test ./internal/application
```
### Примеры работы

#### Регистрация и вход

**HOST**=localhost, а **PORT**=8080

```
curl --location 'http://localhost:8080/api/v1/register' \ // 
--header 'Content-Type: application/json' \
--data '{
  "login": <логин>,
  "password": <пароль>,
}'
```
###### Как это работает?

Отправляется запрос **POST** на ```/api/v1/register``` с телом:
```
{ "login": <логин>, "password": <пароль>}
```

В ответ получает 200 OK или 422 в случае ошибки:

- **Неправильный метод**(только **POST**): ```method not allowed```.
- **Неверное тело запроса:** ```invalid body```.
- **Пользователь существует:** ```user already exists```.
- **Короткий пароль** (меньше *5 символов*): ```short password```.

```
curl --location 'http://localhost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
  "login": <логин>,
  "password": <пароль>,
}'
```
###### Как это работает?

Отправляется запрос **POST** на ```/api/v1/login``` с телом:
```
{ "login": <логин>, "password":<пароль> }
```
- Тело ответа:
```
{"access_token": <токен>}
```

В ответ получает 200 OK и **JWT-токен для последующей авторизации**, 422 в случае ошибки.

- **Неверный метод**(должен быть **только POST**): ```method not allowed```
- **Неверное тело запроса:** ```invalid body```
- **Пользователь не найден:** ```invalid login or password```

Для всех операций, кроме **login** и **register**, нужен Authorization:
```
Authorization: Bearer <полученный при login токен>
```

 - Заголовок не найден: ```invalid header```.
 - Токен не верный: ```invalid token```.

#### Вычисления

```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": <выражение>
}'
```

##### Коды ответа: 
 - 201 - выражение принято для вычисления
 - 422 - невалидные данные
 - 500 - что-то пошло не так

##### Тело ответа

```
{
    "id": <уникальный идентификатор выражения> // его ID
}
```

##### Тело ответа:

Получение всех сохранённых выражений(**ID** не нужен).

```
{
    "expressions": [
        {
            "id": 2352351,
            "status": "OK",  // готово
            "result": 3
            "error": ""
        },
        
            "id": 5372342,
            "status": "Calculation",  // в процессе
            "result": 3
            "error": ""
        },
        {
            "id": 8251431,
            "status": "error",  // статус ошибки
            "result": 0
            "error": "someting error"
        },
        {
            "id": 34942763,
            "status": "Wait", // ждёт выполнения
            "result": 0
            "error": ""
        }
    ]
}
```
##### Коды ответа:
 - 200 - успешно получен список выражений
 - 500 - что-то пошло не так

#### Получение выражения по его идентификатору

```
curl --location 'http://localhost:8080/api/v1/expressions/<id выражения>'
```

##### Тело ответа:

```
{
    "expression":
        {
            "id": <идентификатор выражения>,
            "status": <статус вычисления выражения>,
            "result": <результат выражения>
        }
}
```

##### Коды ответа:
 - 200 - успешно получено выражение
 - 404 - нет такого выражения
 - 500 - что-то пошло не так

### Наглядный пример (копируйте curl и пробуйте)

```
curl --location 'http://localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
  "login": "user0",
  "password": "user0_password"
}'
```

```
curl --location 'http://localhost:8080/api/v1/login' \ // 
--header 'Content-Type: application/json' \
--data '{
  "login": "user0",
  "password": "user0_password"
}'
```
Считываем из тела ответа токен.

```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' --header 'Authorization: Bearer <токен>' \
--data '{
  "expression": "2+2/2"
}'
```

##### Ответ
Статус 201(успешно создано);
```
{
    "id": 12345 // пример
}
```

#### Получение выражения по ID

```
curl --location 'localhost:8080/api/v1/expressions/12345'
```

##### Ответ
Статус 200(успешно получено);
```
{
    "expression":
        {
            "id": 12345,
            "status": "OK",
            "result": 321
        }
}
```

#### Получаем все выражения

```
curl --location 'http://localhost:8080/api/v1/expressions'
--header 'Authorization: Bearer <токен>'
```

##### Ответ
Статус 200(успешно получены);
```
{
    "expressions": [
        {
            "id": 12345,
            "status": "OK",
            "result": 321
        },
    ]
}
```
