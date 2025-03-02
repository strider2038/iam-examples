# build stage
FROM golang:alpine AS build

ADD . /project

RUN set -ex \
    && adduser -D -g '' appuser \
    && cd /project \
    && go mod download \
    && cd /project/cmd/api_gateway \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o api_gateway \
    && ls -la | grep "api_gateway"

# final stage
FROM alpine

WORKDIR "/app"

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /project/cmd/api_gateway /app/

USER appuser

EXPOSE 8089

ENTRYPOINT [ "/app/api_gateway" ]
