FetchTaskServer
=====================

В проекте используется [Nats](https://www.nats.io/) для создания API. И mux для frontend API
Для загрузки зависимостей используется [Go modules](https://github.com/golang/go/wiki/Modules)
Проект основывается на другом моём [проекте](https://github.com/Atluss/Go-Nats-Api-Example).

 Запусе
-----------------------------------
Запускать сначала докер для nats (см. ниже), потом запускаем /api/api.go

 Запросы
 ----------------------------------
* Загрузить ресурс: **/v1/fetch** в тело запроса передать json в формате:
 ```json
{"Method":"GET", "Url":"https://yandex.ru/"}
```
* Получение значение по id: **/v1/get/{id}**, где id это id записи
* Получение всего списка значений: **/v1/list/**
* Удаление элемента по id: **/v1/delete/{id}**, где id это ид записи

Запуск докера
-----------------------------------
Как установить и запустить: 
 1. [Установка Docker-CE (ubuntu)](https://docs.docker.com/install/linux/docker-ce/ubuntu/)
 2. [Установка Docker compose](https://docs.docker.com/compose/install/)
 3. Распоковать docker/docker.zip в папку(это контейнер содержит: Nats 1.4.1, Postgres 11.2)
 4. В распакованной папке запустить: `sudo docker-compose up`
 
Файл настроек
-----------------------------------
Файл настроек используется формат json по [RFC7159](https://tools.ietf.org/html/rfc7159)
 
Пример настроек settings.json:
 ```json
 {
   "name": "api",
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