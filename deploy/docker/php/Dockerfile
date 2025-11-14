# Base Stage - Heavy operations that rarely change
FROM php:8.2-fpm-alpine as base

# Install system packages (cached unless Dockerfile changes)
RUN apk add --no-cache $PHPIZE_DEPS linux-headers git curl unzip

# Install PHP extension installer (cached unless Dockerfile changes)
COPY --from=mlocati/php-extension-installer:2 /usr/bin/install-php-extensions /usr/local/bin/

# Install PHP extensions (cached unless Dockerfile changes) 
RUN install-php-extensions grpc pdo_mysql intl zip opcache

# Install Composer (cached unless Dockerfile changes)
COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

# Development Stage - Only rebuilds when composer files change
FROM base as development

WORKDIR /var/www/html

# Copy composer files (only invalidates cache when these files change)
COPY composer.json composer.lock ./

# Install dependencies (only runs when composer files change)
RUN composer install --optimize-autoloader --no-scripts

# Source code is mounted as volume, not copied
