FROM golang:alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
# RUN apk add curl
# RUN curl -L  https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.20
WORKDIR /app
COPY start.sh .
RUN chmod +x /app/start.sh
COPY --from=builder /app/main .
# COPY --from=builder /app/migrate.linux-amd64 /app/migrate
#Copy env
COPY app.env .
COPY db/migration /app/db/migration
EXPOSE 8080
CMD [ "/app/main" ]


# entrypoint becomes the main entry and CMD becomes arguments passed to the entrypoint
ENTRYPOINT [ "/app/start.sh" ]

