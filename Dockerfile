FROM alpine

RUN apk add --update ca-certificates

WORKDIR /src/users-api

COPY bin/microservices-demo-users-api /usr/bin/users-api
COPY bin/goose /usr/bin/goose
COPY bin/wait-db /usr/bin/wait-db

COPY database/migrations/* /src/users-api/migrations/

EXPOSE 3050

CMD ["/bin/sh", "-l", "-c", "wait-db && cd /src/users-api/migrations/ && goose postgres ${POSTGRES_DSN} up && users-api"]