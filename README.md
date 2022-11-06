# Line_Bot

***

## Needs to install

1. docker
2. go 1.18

***

## Setup

1. set channel_secret and channel_token in `config/token.json`
2. set user ud in `config/user.json`

***

## APIs

| Description             | Method  | Token       | Path                       |
| ----------------------- | ------- | ----------- | -------------------------- |
| Get User's Message List | GET     | *necessary* | /api/querymessage/:user_id |
| Send Message To User    | POST    | *necessary* | /api/sendmessage           |
---
### > GET  /api/querymessage/:user_id
#### Request Header

* Content-Type: application/json
* Authorization: Bearer `{channel access token}`

#### Response Body

* id: `string`
* message: `string`

``` JSON
[
    {
        "id": "Ua42f1d02f01d55c94b8a45a665a4fbbd",
        "message": "Hello World",
    }, {
        ...
    }
]
```
---
### > POST /api/sendmessage

#### Request Header

* Content-Type: application/json
* Authorization: Bearer `{channel access token}`

#### Request Body

* user: `string`
* type: `string`
* text: `string`

``` JSON
    {
        "user": "Ua42f1d02f01d55c94b8a45a665a4fbbd",
        "type": "text",
        "text": "Hello World",
    }
```

#### Response Body

`[]`
