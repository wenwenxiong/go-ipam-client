FROM golang:1.18-alpine as builder

RUN apk add \
    binutils \
    gcc \
    git \
    libc-dev \
    make

WORKDIR /work
COPY . .
ENV GOPROXY="https://goproxy.cn,https://goproxy.io,direct"
RUN make client

FROM alpine:3.16
COPY --from=builder /work/bin/* /
ENTRYPOINT [ "/go-ipam-client" ]
