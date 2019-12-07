FROM alpine

RUN apk add --update ca-certificates

WORKDIR /src/users-api

COPY bin/noken-users-api /usr/bin/noken-users-api

EXPOSE 3050

CMD ["/bin/sh", "-l", "-c", "noken-users-api"]