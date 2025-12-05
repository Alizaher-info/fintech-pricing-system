# Base Stage - Heavy operations that rarely change
FROM php:8.2-fpm-alpine as base

# Install system packages (cached unless Dockerfile changes)
# hadolint ignore=DL3018,SC2086
RUN apk add --no-cache $PHPIZE_DEPS linux-headers git curl unzip dcron busybox-suid

# Install PHP extension installer (cached unless Dockerfile changes)
COPY --from=mlocati/php-extension-installer:2 /usr/bin/install-php-extensions /usr/local/bin/

# Install PHP extensions (cached unless Dockerfile changes) 
RUN install-php-extensions grpc pdo_mysql intl zip opcache xdebug

# Copy Xdebug configuration
COPY xdebug.ini /usr/local/etc/php/conf.d/docker-php-ext-xdebug.ini

# Install Composer (cached unless Dockerfile changes)
COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

# Development Stage - Only rebuilds when composer files change
FROM base as development

WORKDIR /var/www/html

# Copy composer files (only invalidates cache when these files change)
COPY composer.json composer.lock ./

# Create src/Grpc/GPBMetadata directory to prevent autoload scan errors
# Install dependencies and setup cron (only runs when composer files change)
RUN mkdir -p src/Grpc/GPBMetadata && \
    composer install --optimize-autoloader --no-scripts --no-cache

# Setup cron for scheduled tasks
COPY cron/cleanup-security-events /etc/crontabs/www-data
RUN chmod 0644 /etc/crontabs/www-data && \
    mkdir -p /var/log && \
    touch /var/log/cron.log && \
    chown www-data:www-data /var/log/cron.log && \
    echo '#!/bin/sh' > /usr/local/bin/docker-entrypoint.sh && \
    echo 'crond -l 2 -f &' >> /usr/local/bin/docker-entrypoint.sh && \
    echo 'php-fpm' >> /usr/local/bin/docker-entrypoint.sh && \
    chmod +x /usr/local/bin/docker-entrypoint.sh

CMD ["/usr/local/bin/docker-entrypoint.sh"]

# Source code is mounted as volume, not copied
