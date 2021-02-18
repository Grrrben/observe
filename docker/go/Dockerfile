# Compile stage
FROM golang:1.15.7 AS compile_env

ADD . /gosource
WORKDIR /gosource/cmd

RUN go test -v ../test/.
RUN go build -o /goapp

# Final stage
FROM debian:buster

RUN mkdir app
WORKDIR /app

COPY --from=compile_env /goapp ./
COPY --from=compile_env /gosource/.env ./

CMD ["./goapp"]
