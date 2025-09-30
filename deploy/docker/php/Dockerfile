FROM php:8.2-fpm-alpine

# Installer für PHP-Extensions (nimmt dir viel Arbeit ab)
COPY --from=mlocati/php-extension-installer:2 /usr/bin/install-php-extensions /usr/local/bin/

# Basis-Buildtools (für einzelne Extensions nützlich)
RUN apk add --no-cache $PHPIZE_DEPS linux-headers git curl unzip

# gRPC + weitere Extensions
RUN install-php-extensions grpc pdo_mysql intl zip opcache

# (optional) Composer
# COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/html
