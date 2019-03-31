FetchTaskServer
=====================

В проекте используется [Nats](https://www.nats.io/) для создания API. И mux для frontend API
Для загрузки зависимостей используется [Go modules](https://github.com/golang/go/wiki/Modules)
Проект основывается на другом моём [проекте](https://github.com/Atluss/Go-Nats-Api-Example).

 Запуск
-----------------------------------
Запускать сначала докер для nats (см. ниже), потом запускаем /api/api.go
Настройки api хранятся в settings.json, который должен находится в той же директории где запускается api. 

 Запросы
 ----------------------------------
 
Название  | Запрос | Описание
 ----------------|----------------------|-----
 Загрузить ресурс в список элементов | **/v1/fetch** | в тело запроса передать json (стандарт [RFC7159](https://tools.ietf.org/html/rfc7159)) формат см. ниже
 Получение элемента | **/v1/get/{id}** | где id это id записи
 Получение списка элементов | **/v1/list/** |
 Удаление элемента | **/v1/delete/{id}** | где id это ид записи
 
 ##### Описание формата для запроса **/v1/fetch**
 Оба поля в запросе обязательны, реализовано для метода ***GET***, в ***Url*** не забывайте передавать адрес с указание протокала.
 ```json
{"Method":"GET", "Url":"https://yandex.ru/"}
```

Запуск докера
-----------------------------------
Как установить и запустить Docker: 
 1. [Установка Docker-CE (ubuntu)](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
 2. [Установка Docker compose](https://docs.docker.com/compose/install/)
 3. Распоковать docker/docker.zip в папку(этот контейнер содержит: Nats 1.4.1, Postgres 11.2)
 4. В распакованной папке запустить: `sudo docker-compose up`
 
Файл настроек
-----------------------------------
Файл настроек используется формат json по [RFC7159](https://tools.ietf.org/html/rfc7159)
 
Пример настроек settings.json:
 ```json
 {
   "name": "FetchTaskServer",
   "version": "1.0.0",
   "port" : "10000",
   "nats": {
     "version" : "1.4.2",
     "reconnectedWait" : 5,
     "address" : [
       {
         "host" : "localhost",
         "port" : "54222"
       },
       {
         "host" : "localhost",
         "port" : "54222"
       }
     ]
   }
 }
 ```