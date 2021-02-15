# Observe

Citizen Science project dag foto's van 1 object of landschap door de tijd vastlegt.

## Stack

go1.15.7  
postgres 13.1   
html5

docker && docker-compose

## useful cmd's

DB from scratch

```
$ docker container ls -a
$ docker container rm [id of DB container]
$ docker-compose up --build 
```

Postgres

```
psql -U postgres observe
```

```
docker-compose exec db psql -U postgres observe -c "SELECT * FROM project;"
docker-compose exec db psql -U postgres observe -c "SELECT * FROM observation;"
```


