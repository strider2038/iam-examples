# build stage
FROM golang:alpine AS build

ADD . /project

RUN set -ex \
    && adduser -D -g '' appuser \
    && cd /project \
    && go mod download \
    && cd /project/cmd/app \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o app \
    && ls -la | grep "app"

# final stage
FROM alpine

WORKDIR "/app"

COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /project/cmd/app /app/

USER appuser

EXPOSE 4040

ENTRYPOINT [ "/app/app" ]
