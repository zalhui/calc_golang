# Простой веб-сервис для вычисления арифметических выражений

# Описание
проект calc_golang это веб-сервис, который вычисляет простые арифметические выражения с положительными числами, скобками и знаками + - / *.
выражения передаются через http запрос

## Структура проекта

- `cmd/` - директория с файлом main.go
- `internal/application/` - директория с кодом сервера и тестами для проверки работы сервера
- `pkg/calculator/` - директория с кодом калькулятора и тестами для проверки работы калькулятора

## Запуск

1. Установите Golang https://go.dev/dl/
2. Установите Git https://git-scm.com/downloads
3. C помощью командной строки клонируйте проект с GitHub
   
    ```
    git clone https://github.com/zalhui/calc_golang
    ```
5. Перейдите в директорию с проектом и запустите сервер
    ```
    go run ./...
    ```

# Работа с сервисом

Для работы с данным сервисом используйте командную строку

**Для корректной работы на Windows необходимо использовать *Git Bash***(устанавливается вместе с Git)

Также работа с сервисом вмозожна через Postman. Для работы вставьте в строку для URL адрес:
```
http://127.0.0.1:8080/api/v1/calculate
```

для отправки запроса используйте команду(вместо '...' введите выражение для калькулятора):

```
curl --location 'http://127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"..."
}'
```

## Примеры работы с сервисом
### Корректный запрос:
Введя данный запрос:
```
curl --location 'http://127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"2+2*2/(2+2)"
}'
```
вы получите ответ    
```
{"result":3}
```
с кодом [200]

### Запрос с методом не POST:
Введя данный запрос:
```
curl --location --request GET 'http://127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"2++2"
}'
```
вы получите ответ    
```
{"error": "only POST method allowed"}
```
с кодом [405]

### Запрос с неправильным телом:
Введя данный запрос:
```
curl --location 'http://127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"2+2
}'
```
вы получите ответ    
```
{"error": "Bad request"}
```
с кодом [400]

### Запрос с делением на 0(ноль)
Введя данный запрос:
```
curl --location 'http://127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"2+2*2/0"
}'
```
вы получите ответ    
```
{"error":"Expression is not valid. Division by zero"}
```
с кодом [422]

### Запрос с не закрытой скобкой
Введя данный запрос:
```
curl --location 'http://127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"2+(9+7"
}'
```
вы получите ответ    
```
{"error": "Expression is not valid. Number of brackets doesn't match"}
```
с кодом [422]

### Запрос с выражением с буквами
Введя данный запрос:
```
curl --location 'http://127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"2+(9+x)"
}'
```
вы получите ответ    
```
{"error": "Expression is not valid. Only numbers and ( ) + - * / allowed"}
```
с кодом [422]

### Запрос с выражением c лишними знаками действия
Введя данный запрос:
```
curl --location 'http://127.0.0.1:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"2++2"
}'
```
вы получите ответ    
```
{"error": "Expression is not valid. Not enough values"}
```
с кодом [422]
Для остальных ошибок ответ будет

    {"error":"Internal server error"}
    
с кодом [500]
