# Compile stage
FROM golang:1.15.7 AS build-env

ADD . /dockerdev
WORKDIR /dockerdev/cmd

RUN go build -o /goapp

# Final stage
FROM debian:buster

EXPOSE 8000

RUN mkdir app
WORKDIR /app

COPY --from=build-env /goapp ./
COPY --from=build-env /dockerdev/.env ./

CMD ["./goapp"]