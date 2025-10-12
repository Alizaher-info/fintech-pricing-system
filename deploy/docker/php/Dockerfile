FROM php:8.2-fpm-alpine

# Use Alpine packages for faster installation where possible
RUN apk add --no-cache \
    $PHPIZE_DEPS \
    linux-headers \
    git \
    curl \
    unzip \
    # Try to install gRPC from Alpine packages first (much faster)
    grpc \
    grpc-dev \
    protobuf \
    protobuf-dev

# Installer für PHP-Extensions (nimmt dir viel Arbeit ab)  
COPY --from=mlocati/php-extension-installer:2 /usr/bin/install-php-extensions /usr/local/bin/

# Install other extensions (skip gRPC if Alpine version worked)
RUN install-php-extensions pdo_mysql intl zip opcache

# Only install gRPC from source if Alpine packages didn't work
RUN php -m | grep -q grpc || install-php-extensions grpc

WORKDIR /var/www/html
