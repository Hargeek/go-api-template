FROM golang:1.22.6-alpine AS build-env

ENV GOSUMDB=off \
    GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"

WORKDIR /workspace
COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN ln -s /var/cache/apk /etc/apk/cache
RUN --mount=type=cache,target=/var/cache/apk --mount=type=cache,target=/etc/apk/cache \
    --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update --no-cache \
    && apk add --no-cache git make \
    && make buildx

FROM alpine:3.14.0
RUN ln -s /var/cache/apk /etc/apk/cache
RUN --mount=type=cache,target=/var/cache/apk --mount=type=cache,target=/etc/apk/cache \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update --no-cache \
    && apk add --no-cache ca-certificates tzdata bash curl git

COPY --from=build-env /bin/server /server

ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT [ "/server" ]