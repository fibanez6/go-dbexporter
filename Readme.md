# Importing data into a SQL database

Start Postgres:

 ```
$ docker-compose up
  ```

Create 'test' database:
 
 ```
 $ docker cp dbdriver/postgres.sql my_postgres:/tmp/postgres.sql
 $ docker exec -it my_postgres psql -U postgres -c "create database test"
 $ docker exec -it my_postgres psql -U postgres test -f /tmp/postgres.sql
```
 Build application 
  ```
 $ go build .
  ```
Run application
 ```
 ./go-dbexporter [-file=path/to/file] [-set-offset-beginning=true]
 ```
 
 Usage of ./go-dbexporter
  ```
  ./go-dbexporter -h
  ```
Configure PgAdmin:
 ```
 host: host.docker.internal
 database: postgres
 user: fibanez@fibanez.com
 password: fibanez
 ```
