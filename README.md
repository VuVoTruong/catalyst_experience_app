# Overview
It's a simple API for Catalyst Experience App. 

## Note
1. I used Go Echo template which can be found here: https://github.com/nixsolutions/golang-echo-boilerplate
2. You need to change DB configuration in .env file
3. Document is used with github.com/swaggo/swag/cmd/swag
    + Configure: swag init -g cmd/main.go
    + Run: http://localhost:7788/swagger/index.html
4. Call Register to create a valid user to authenticate for backend APIs
5. You now can check around with other APIs

## Features:
- Registration
- Authentication with JWT
- Invitation Token actions include: generating, validating, updating for recall and listing (admin page)
- Swagger docs
- Docker help

## Deployment
Run: docker-compose up -d --build

## TODO:
Still missing check duplicate code when generating