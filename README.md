# py-server

Pay Yourself backend


`min-client` is a minimal example of a client side JS application interacting with the server

`server` is written in Go

## API defs

### Expected request headers

* `Token`: A valid Google user token string, generated using PY's client ID

### `GET` `/v1/usersave`

* 200: `json` of user save belonging to token's ID
* 403: no token was provided, or the provided token was invalid
* 404: no such save belonging to the token's ID (but the token is valid)

### `POST` `/v1/usersave`

Expects JSON body with valid usersave

* 200: save successful
* 403: no token was provided, or the provided token was invalid
* 404: no such save belonging to the token's ID (but the token is valid)

### `DELETE` `/v1/usersave`

* 200: remove successful
* 403: no token was provided, or the provided token was invalid
* 404: no such save belonging to the token's ID (but the token is valid)