Инструкция по запуску geo-service
1. Убедись, что ты находишься в корневой директории проекта: /go-kata
2. Запускай всю систему одной командой: docker compose up --build


Почему запускать из корня:
1. Docker Compose автоматически ищет файл .env
2. Только в корне проекта определены все сервисы (proxy, geo-service, hugo) и их связи
3. Если запускать из подкаталогов (geo-service/ или hugoProxy/proxy/), запускается только один сервис
4. Только запуск из корня обеспечивает правильную работу всей системы

P.S. если требуется проверить работу только geo-service, такая возможность тоже есть:
docker compose up --build из папки go-kata/geo-service


Для проверки только geo-servie
curl -XPOST http://localhost:8080/address/search -d '{"query":"Москва"}' -H "Content-Type: application/json"

Для проверки всей системы
curl -XPOST http://localhost:8081/address/search -d '{"query":"Москва"}' -H "Content-Type: application/json"


# регистрация
curl -X POST http://localhost:8080/api/register \
-H 'Content-Type: application/json' \
-d '{"username":"alex","password":"pass"}'

# авторизация → получаем JWT
TOKEN=$(curl -s -X POST http://localhost:8080/api/login \
-H 'Content-Type: application/json' \
-d '{"username":"alex","password":"pass"}' | jq -r .token)

# защищённая ручка с токеном
curl -X POST http://localhost:8080/api/address/search \
-H "Authorization: Bearer $TOKEN" \
-H 'Content-Type: application/json' \
-d '{"query":"Москва, Тверская 7"}' | jq .

# тот же запрос без токена → 401 / 403
curl -i -X POST http://localhost:8080/api/address/search \
-H 'Content-Type: application/json' \
-d '{"query":"Москва"}'