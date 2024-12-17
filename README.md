## Development

`docker-compose up --build`

backend is open on port `8080` in a docker container and exposed to `localhost:8080`

frontend is open on port `80` in a docker container and exposed to `localhost:3000`

database (mysql) is open on port `3306` in a docker container and exposed to `localhost:3306`

### Debug

`docker exec -it <backend-container>`

`docker exec -it <frontend-container>`

`docker exec -it <mysql-container> mysql -u user -p`

## Production




