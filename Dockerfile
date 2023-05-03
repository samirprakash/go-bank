# build stage
FROM golang:1.20.3-alpine3.17 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

# run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY app.env .
COPY startup.sh .
COPY wait-for.sh .
COPY db/migration ./migration

EXPOSE 8080
ENTRYPOINT [ "/app/startup.sh" ]
CMD [ "/app/main" ]