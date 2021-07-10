FROM golang:alpine AS builder

#Build for R2S router
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=arm64

WORKDIR /build

# 将代码复制到容器中
COPY . .

RUN go build -o push .

FROM scratch

COPY --from=builder /build/push /

ENTRYPOINT ["/push"]