## Hasher
Итоговый проект курсв REBRAIN Golang Advanced.

Данный сервис, взаимодействуя с сервисом считающим хэши (SHA-3), получает из входящих строк их хэши, сохраняет их в БД (PostgresQl) с присвоем id, по которым далее можно будет запрашивать хэши.

### Запуск проекта локально:
- Добавить .env файл с переменными окружения
```shell
cd hash-service
docker-compose up
```

### Environment
- DB_HOST=pgdb
- DB_PORT=5432
- DB_USER
- DB_PASSWORD
- DB_NAME

- GRPC_SERVER=grpc-server:50051
- HASHER_PORT
- ENVIRONMENT // prod or dev


### Swagger spec
```yaml
swagger: "2.0"
info:
  version: "1.0.0"
  title: "Итоговое задание. Хэши."
  description: "Данный сервис должен, взаимодействуя с сервисом считающим хэши (по выбранному вами протоколу), получать из входящих строк их хэши, сохранять их в свою БД (выбор так же за вами) с присвоем id, по которым далее можно будет запрашивать хэши."
schemes:
  - http
produces:
  - application/json
paths:
  /send:
    post:
      summary: "Получает на вход список строк, хэши от которых нужно посчитать и сохранить"
      parameters:
        - in: body
          name: params
          description: "Strings for hash"
          schema:
            $ref: '#/definitions/ArrayOfStrings'
      responses:
        "200":
          description: "Success"
          schema:
            $ref: '#/definitions/ArrayOfHash'
        "400":
          description: "Bad request"
        "500":
          description: "Internal Server Error"
  /check:
    get:
      summary: "Получает по id хэш из хранилища (если есть)"
      parameters:
        - in: query
          name: ids
          description: "Get hash by this id"
          required: true
          type: array
          items:
            type: string
      responses:
        "200":
          description: "Success"
          schema:
            $ref: '#/definitions/ArrayOfHash'
        "204":
          description: "No Content"
        "400":
          description: "Bad request"
        "500":
          description: "Internal Server Error"
definitions:
  ArrayOfStrings:
    type: array
    items:
      type: string
  ArrayOfHash:
    type: array
    items:
      $ref: '#/definitions/Hash'
  Hash:
    type: object
    properties:
      id:
        type: integer
        example: 38
      hash:
        type: string
        example: a7ffc6f8bf1ed76651c14756a061d662f580ff4de43b49fa82d80a4b80f8434a
    required:
      - id
      - hash
```
