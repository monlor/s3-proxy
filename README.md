## Why

For Apple Shotcuts or MWeb APP or Bash scripts

### Apple Shotcuts
https://www.icloud.com/shortcuts/936ea0c7a9404a0ebf0cbabf9d522149

## How to use

### Deploy

```bash
docker run -d \
  --restart unless-stopped \
  --name s3-proxy \
  -e S3_ACCESS_ID=xxx \
  -e S3_SECRET_KEY=xxx \
  -e S3_REGION=xxx \
  -e S3_BUCKET_NAME=xxx \
  -e HTTP_USERNAME=xxx \
  -e HTTP_PASSWORD=xxx \
  -p 8080:8080 \
  ghcr.io/monlor/s3-proxy:main
```

### Upload by curl

```bash
curl -X POST -u foo:bar -F file=@test.txt localhost:8080/upload
```

### Command line

#### add environment variables

```bash
export S3_PROXY_URL=https://s3.xxx.com
export S3_PROXY_USERNAME=xxx
export S3_PROXY_PASSWORD=xxx
# replace url with cdn (optional)
export S3_CDN_URL=xxx
```

#### install s3-proxy

```bash
brew tap monlor/taps
brew install monlor/taps/s3-proxy
```

#### Config 

```bash
s3-proxy
```

### Typora

add custom command

```bash
s3-proxy
```