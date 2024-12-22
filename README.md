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
    git clone https://github.com/nikitakutergin59/RPN
    ```
5. Перейдите в директорию с проектом и запустите сервер
    ```
    go run ./...
    ```

## Примеры запросов

Создайте новую командную строку и введите команду

-**Например:**

    curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"52+52\"}" http://localhost:8080
    
-**вы увидете ответ:**
 
    {"result":104}
  
с кодом [200]

-**Ещё один пример работы, если вы введёте команду, например:**

    curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"10/0 \"}" http://localhost:8080
    
-**вы получите ответ:**

    {"error":"Division by zero"}
    
с кодом [422]

-**а так же если вы введёте команду например:**

    curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"ткцотикзтц9отищзтзтщшкцьтещшцть \"}" http://localhost:8080
    
вы получите ответ

    {"error":"Expression is not valid"}

так-же с кодом [422]

Для остальных ошибок ответ будет

    {"error":"Internal server error"}
    
с кодом [500]
