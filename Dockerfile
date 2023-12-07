FROM golang AS build-env
LABEL authors="QunBo"

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /workspace
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-s -w" -o /bin/server ./cmd/main.go

# 创建配置文件目录并复制配置文件
RUN mkdir -p /workspace/config
COPY config/config.yaml /workspace/config/

FROM alpine
RUN ln -s /var/cache/apk /etc/apk/cache
RUN --mount=type=cache,target=/var/cache/apk --mount=type=cache,target=/etc/apk/cache \
    sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update --no-cache \
    && apk add --no-cache ca-certificates tzdata bash curl

# 从构建环境复制server二进制文件
COPY --from=build-env /bin/server /server
# 从构建环境复制配置文件
COPY --from=build-env /workspace/config/config.yaml /workspace/config/

ENV GIN_MODE=release
EXPOSE 8080
ENTRYPOINT ["/server"]
