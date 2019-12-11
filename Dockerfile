FROM alpine

RUN apk add --update ca-certificates

WORKDIR /src/users-api

COPY bin/noken-users-api /usr/bin/noken-users-api

COPY database/migrations/* /src/users-api/migrations/

EXPOSE 3050

CMD ["/bin/sh", "-l", "-c", "wait-db && cd /src/users-api/migrations/ && goose postgres ${POSTGRES_DSN} up && noken-users-api"]