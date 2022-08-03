FROM golang:1.19.0 as builder

#ENV GOPROXY=https://goproxy.cn

COPY . /s3-proxy

RUN cd /s3-proxy && go build -o s3-proxy

FROM debian:buster-slim

ENV GIN_MODE=release

COPY --from=builder /s3-proxy/s3-proxy /s3-proxy

CMD ["/s3-proxy"]