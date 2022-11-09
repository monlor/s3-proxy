FROM --platform=$BUILDPLATFORM golang:alpine as build
# Redundant, current golang images already include ca-certificates
# ENV GOPROXY=https://goproxy.cn
ARG TARGETARCH
RUN apk --no-cache add ca-certificates
WORKDIR /go/src/app
COPY . .
RUN GOOS=linux GOARCH=$TARGETARCH go build -o s3-proxy

FROM alpine:latest
ENV GIN_MODE=release
EXPOSE 8081
# copy the ca-certificate.crt from the build stage
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/s3-proxy /s3-proxy
ENTRYPOINT ["/s3-proxy"]