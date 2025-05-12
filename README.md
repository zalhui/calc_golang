# Веб-сервис для параллельного вычисления арифметических выражений

*ЕСЛИ ВОЗНИКЛИ ПРОБЛЕМЫ С МОИМ ПРОЕКТОМ, ВЫ МОЖЕТЕ СВЯЗАТЬСЯ СО МНОЙ В TELEGRAM: @unbiunii*

## Описание
Проект `calc_golang` — это веб-сервис, который вычисляет арифметические выражения с положительными числами, скобками и операциями `+`, `-`, `*`, `/`. Выражения обрабатываются асинхронно с использованием системы задач, что позволяет параллельно вычислять части сложных выражений. Сервис предоставляет REST API для отправки выражений, получения их статуса и результатов, а также просмотра всех сохраненных выражений.

В последней версии проекта реализована регистрация и аутентификация пользователей, а также персистентность, что позволяет хранить данные о пользователях и выражениях после завершения работы сервиса.

В качестве СУБД в данном проекте используется SQLlite.

Сервис разделен на две части:

- Сервер, который принимает арифметическое выражение, переводит его в набор последовательных задач и обеспечивает порядок их выполнения. Далее будем называть его оркестратором.
- Вычислитель, который может получить от оркестратора задачу, выполнить его и вернуть серверу результат. Далее будем называть его агентом.

Основные возможности сервиса:
- Регистрация и аутентификация пользователя через POST-запросы к `/api/v1/register` и `/api/v1/login`.
- Отправка выражения на вычисление через POST-запрос к `/api/v1/calculate`.
- Получение статуса и результата конкретного выражения по его ID через GET-запрос к `/api/v1/expressions/{id}`.
- Получение списка всех выражений через GET-запрос к `/api/v1/expressions`.

Сервис возвращает ответы в формате JSON. 


Ответы сопровождаются следующими кодами:

| Код ответа | Описание |
| --- | --- |
| `200` | Выражение успешно получено / список выражений получен успешно |
| `201` | Выражение принято для вычисления |
| `422` | Невалидные данные |
| `404` | Выражение не найдено |
| `400` | Неверный формат запроса |
| `405` | Неверный метод запроса |
| `500` | Что-то пошло не так |

Схема работы сервиса представлена на изображении ниже:

![Image alt](https://github.com/zalhui/calc_golang/blob/main/image.png)

Описание схемы:
- Пользователь → Оркестратор: Отправляет POST-запрос на /api/v1/calculate с выражением, например {"expression": "2+2*2"}.
- Оркестратор: Парсит выражение и генерирует задачи для выполнения.
- Оркестратор → Пользователь: Возвращает уникальный идентификатор выражения, например {"id": "expr-id"}.
- Цикл (loop):
- Агент → Оркестратор: Запрашивает задачу через GET /internal/task.
- Оркестратор → Агент: Возвращает задачу (например, часть выражения для вычисления).
- Агент: Выполняет задачу (например, считает 2*2=4).
- Агент → Оркестратор: Отправляет результат через POST /internal/task/result, например {"id": "task-id", "result": 4}.
- Оркестратор: Обновляет статус задачи на "completed".
- Пользователь → Оркестратор: Запрашивает статус выражения через GET /api/v1/expressions/{id}.
- Оркестратор → Пользователь: Возвращает данные выражения, включая статус и итоговый результат, например {"expression": { "status": "completed", "result": 6, ... }}.
## Структура проекта

- `cmd/` - директория с файлами `orchestrator/main.go` и `agent/main.go` для запуска оркестратора и агента.
- `config/` - конфигурация сервиса.
- `internal/agent/worker/` - код агента, выполняющего вычисления задач.
- `internal/auth/` - логика регистации и аутентификации.
- `internal/common/models/` - структуры данных для выражений и задач.
- `internal/db` - описание схем базы данных.
- `internal/middleware` - middleware для проверки аутентификации.
- `internal/orchestrator/application/` - логика и хэндлеры оркестратора (сервер), который принимает запросы, распределяет задачи и возвращает результаты.
- `internal/orchestrator/repository` - логика работы с бд.
- `pkg/calculation/` - логика перевода выражений в обратную польскую запись, парсинга выражений и преобразования их в задачи.
- `.env` - файл с переменными среды(время операций и вычислительная мощность).

## Запуск

### Общие шаги
1. Установите Golang: [https://go.dev/dl/](https://go.dev/dl/).
2. Установите Git: [https://git-scm.com/downloads](https://git-scm.com/downloads).
3. Склонируйте проект с GitHub:
```
git clone https://github.com/zalhui/calc_golang
```
4. Перейдите в директорию проекта:
```
cd calc_golang
```
### Сценарий 1: Запуск на Unix-системах (Linux, macOS) с помощью Makefile
1. Убедитесь, что `make` установлен (обычно предустановлен):
- Если нет, установите:
  - Linux: `sudo apt install make` (Ubuntu/Debian) или `sudo yum install make` (CentOS).
  - macOS: `xcode-select --install`.
2. Запустите проект одной командой:
```
make
```
Это запустит оркестратор и агента в фоновом режиме. Логи будут выводиться в консоль.
3. Для остановки процессов:
```
make clean
```

Сервер будет доступен по адресу `http://localhost:8080`.
### Сценарий 2: Запуск на Windows с помощью двух терминалов
1. Откройте два терминала (например, Git Bash, CMD или PowerShell).
2. В первом терминале запустите оркестратор:
```
go run ./cmd/orchestrator/main.go
```
3. Во втором терминале запустите агента:
```
go run ./cmd/agent/main.go
```
4. Для остановки нажмите `Ctrl+C` в каждом терминале.

Сервер будет доступен по адресу `http://localhost:8080`.
## Работа с сервисом

Для взаимодействия с сервисом используйте командную строку (Git Bash на Windows) или инструменты вроде Postman.

### Основные эндпоинты

1. **Регистрация пользователя**  
URL: `http://localhost:8080/api/v1/register`  
Метод: `POST`  
Тело запроса:
```
{
    "login":"ваш_логин",
    "password":"ваш_пароль"
}
```
Ответ: статус операции.

1. **Регистрация пользователя**  
URL: `http://localhost:8080/api/v1/login`  
Метод: `POST`  
Тело запроса:
```
{
    "login":"ваш_логин",
    "password":"ваш_пароль"
}
```
Ответ: персональный JWT токен.

3. **Отправка выражения на вычисление**  
URL: `http://localhost:8080/api/v1/calculate`  
Метод: `POST`  
Тело запроса:
```
{
"expression": "арифметическое выражение"
}
```
Ответ: ID выражения для последующего отслеживания.

4. **Получение статуса и результата выражения**  
URL: `http://localhost:8080/api/v1/expressions/{id}`  
Метод: `GET`  
Ответ: полная информация о выражении, включая статус, результат (если вычислено).

5. **Получение списка всех выражений**  
URL: `http://localhost:8080/api/v1/expressions`  
Метод: `GET`  
Ответ: список всех выражений с их статусами и результатами.

## Примеры работы с сервисом

*Примеры приведены для командной строки Git Bash.*

### 1. Корректный запрос на регистрацию
Запрос:

```
curl --location 'http://localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTQyMDEsInVzZXJfaWQiOiI2NDRkZTRhYi0yNzgwLTRmZTAtODczMC0xY2VlYWUxNTgyOTUifQ.5qW_h9TrjlroWnP9nKfkx47AKvUAsPOnDYbGmGfq4lM' \
--data '{
    "login":"user1",
    "password":"1"
}'
```
Ответ:
```
{"status":"success"}
```
Код: `[200]`

### 2. Корректный запрос на аутентификацию
Запрос:

```
curl --location 'http://localhost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTQyMDEsInVzZXJfaWQiOiI2NDRkZTRhYi0yNzgwLTRmZTAtODczMC0xY2VlYWUxNTgyOTUifQ.5qW_h9TrjlroWnP9nKfkx47AKvUAsPOnDYbGmGfq4lM' \
--data '{
    "login":"user1",
    "password":"1"
}'
```
Ответ:
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4"
}
```
Код: `[200]`

*В ПОСЛЕДУЮЩИХ ЗАПРОСАХ НЕОБХОДИМО ИСПОЛЬЗОВАТЬ ДАННЫЙ ТОКЕН В ЗАГОЛОВКЕ ЗАПРОСА ПО КЛЮЧУ Authorization со значением "Bearer {ваш токен}"*

### 3. Корректный запрос на вычисление
Запрос:

```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4' \
--data '{
    "expression":"2+2*2"
}'
```
Ответ:
```
{"id":"6a451ccb-36ba-4045-a3dd-a21f7beb45dd","message":"Expression accepted for processing","status":"pending"}

```
Код: `[201]`

### 4. Получение результата выражения
Запрос (через несколько секунд, чтобы вычисления завершились):
```
curl --location 'http://localhost:8080/api/v1/expressions/6a451ccb-36ba-4045-a3dd-a21f7beb45dd' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4'
```
*Результат выражения лежит в поле "result":*

Ответ:
```
{
    "created": "2025-05-12T09:15:49.1191212+03:00",
    "expression": "2+2*2",
    "id": "6a451ccb-36ba-4045-a3dd-a21f7beb45dd",
    "result": 6,
    "status": "completed"
}
```

Код: `[200]`

### 5. Получение списка всех выражений
Запрос:
```curl --location 'http://localhost:8080/api/v1/expressions' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4'
```

```
Ответ (пример с одним выражением):
{
    "expressions": [
        {
            "created": "2025-05-12T09:15:49.1191212+03:00",
            "expression": "2+2*2",
            "id": "6a451ccb-36ba-4045-a3dd-a21f7beb45dd",
            "result": {
                "Float64": 6,
                "Valid": true
            },
            "status": "completed"
        }
    ]
}
```

Код: `[200]`

### 6. Запрос с неправильным методом (не POST)
Запрос:
```
curl --location --request GET 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4' \
--data '{
    "expression":"2+2*2"
}'
```
Ответ:
```
404 page not found
```
Код: `[404]`
### 7. Запрос с некорректным телом
Запрос:
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4' \
--data '{
    "expression":"2+2*2
}'
```
Ответ:
```
Invalid request format
```
Код: `[400]`

### 8. Запрос с делением на ноль
Запрос:
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4' \
--data '{
    "expression":"2/0"
}'
```
Ответ:
```
{"id":"48c2f8cb-a5ed-4ef6-8129-6d60e374b8f6","message":"Expression accepted for processing","status":"pending"}

```
Код: `[201]`

При запросе данного выражения по айди получим ответ:
```
{
    "created": "2025-05-12T09:24:10.6219323+03:00",
    "expression": "2/0",
    "id": "48c2f8cb-a5ed-4ef6-8129-6d60e374b8f6",
    "result": 0,
    "status": "error"
}
```
В поле status находится ошибка, что свидетельствует о том, что невалидное выражение не считается сервером.
### 9. Запрос с несовпадающими скобками
Запрос:
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4' \
--data '{
    "expression":"2+(3"
}'
```
Ответ:
```
error converting expression to RPN : expression is not valid. number of brackets doesn't match
```

Код: `[422]`

### 10. Запрос с недопустимыми символами
Запрос:
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4' \
--data '{
    "expression":"2+x"
}'
```
Ответ:
```
error converting expression to RPN : expression is not valid. only numbers and ( ) + - * / allowed
```
Код: `[422]`

### 11. Запрос с недостаточным количеством значений
Запрос:
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4' \
--data '{
    "expression":"2++2"
}'
```
Ответ:
```
expression is not valid. not enough values
```
Код: `[422]`

### 12. Запрос выражения по несуществующему ID
Запрос:
```
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDcxMTY3NjUsInVzZXJfaWQiOiJiZjcxYzk3Mi05ZDg3LTRlYWItODg1NS04ZTRhNjY0NDM2ZmUifQ.XpIiMNs6Z-8i0Ps7yQXBAQYlx92A5iWvV7b-Zj8-Xw4' \
--data '{
    "expression":"2++2"
}'
```
Ответ:
```
Expression not found
```
Код: `[404]`
**В случае иной ошибки на стороне сервера будет получен ответ:**
```
{"error":"Internal server error"}
```
с кодом `[500]`.
## Замечания
- Время выполнения операций (сложение, вычитание, умножение, деление) задается в `.env` и по умолчанию составляет 1 секунду на операцию. Вы можете изменять эти значения для более наглядной демонстрации функций сервиса.
- Статус выражения может быть `"pending"` (в процессе), `"completed"` (завершено) или `"error"` (ошибка).
- Для полного завершения вычисления сложных выражений может потребоваться несколько секунд в зависимости от количества задач и настроек `.env`.
