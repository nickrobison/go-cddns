FROM alpine:latest
LABEL version="0.1.0"
LABEL description="Docker image for running go-cddns application"
LABEL maintainer="Nick Robison nick@nickrobison.com"

# Setup the necessary SSL certs
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true

# Add the go binary and config files
COPY bin/go-cddns*_amd64 /app/go-cddns
WORKDIR /app
CMD ["./go-cddns", "--config=/etc/config.json"]
