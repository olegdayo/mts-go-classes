# Go stuff hapening here
![Build](https://github.com/OFFLUCK/mts-go-classes/actions/workflows/build.yml/badge.svg)

## This is a CRUD app

Examples of requests:
- GET localhost:8080/ping
- GET localhost:8080/users
- GET localhost:8080/user/{username}
- PUT localhost:8080/user/{old_username}?username={new_username}&password={new_password}
- DELETE localhost:8080/user/{username}
- POST localhost:8080/registration?username={username}&password={password}
- GET localhost:8080/login
    - Needs basic auth header
- GET localhost:8080/verify
    - Needs basic auth header
