## Пример .env файла:
```
PORT=8080
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB_NAME=pr_service_db
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
ATTEMPTS_TO_CONNECT=3
BASE_DELAY=1
```

## Запуск приложения:
```
docker-compose up --build -d
```

## Добавлен эндпоинт статистики:

+ Группировка по user_id:
Запрос:
```
http://localhost:8080/pullRequest/stats или http://localhost:8080/pullRequest/stats?group_by=user
```
Пример ответа:
```
{
    "group_by": "user",
    "stats": [
        {
            "key": "user4",
            "count": 2
        },
        {
            "key": "user1",
            "count": 1
        },
        {
            "key": "user3",
            "count": 1
        }
    ]
}
```
+ Группировка по pull_request_id:
Запрос:
```
http://localhost:8080/pullRequest/stats?group_by=pr
```
Пример ответа:
```
{
    "group_by": "pr",
    "stats": [
        {
            "key": "pr-52",
            "count": 2
        },
        {
            "key": "pr-53",
            "count": 2
        }
    ]
}
```

## Насчёт эндпоинта /team/add
1) Проверки: ни один из участников не должен уже принадлежать другой команде и нельзя создать команду с названием, уже существующим
2) Я не понял, что значит "обновляет пользователей". Единственное, что относится к обновлению, это тот случай, в котором пользователь уже существует и не привязан к какой-то команде (если вручную удалили у него это поле в базе). В таком случае, существующий пользователь обновляется (ему присвается переданное в запросе название команды)