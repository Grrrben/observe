# Observe

Citizen Science project that continuously documents pictures of a given subject.

## Stack

go1.16  
postgres 13.1   
html5  

docker && docker-compose  
kubernetes

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


