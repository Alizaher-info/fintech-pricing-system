# Cron Scheduled Tasks Setup

## Overview

Automated cleanup of old security events runs daily at 2:00 AM inside the PHP Docker container.

## Files Created

1. **`src/Command/CleanupSecurityEventsCommand.php`** - Symfony command that deletes old security events
2. **`cron/cleanup-security-events`** - Cron schedule file
3. **Updated `deploy/docker/php/Dockerfile`** - Installs cron and configures startup
4. **Updated `deploy/compose/docker-compose.yml`** - Added cron logs volume

## Configuration

### Retention Policy
Default: **365 days (1 year)**

To change retention period, edit `cron/cleanup-security-events`:
```bash
# 90 days (3 months)
0 2 * * * cd /var/www/html && php bin/console app:cleanup-security-events --days=90 >> /var/log/cron.log 2>&1

# 730 days (2 years)
0 2 * * * cd /var/www/html && php bin/console app:cleanup-security-events --days=730 >> /var/log/cron.log 2>&1
```

### Cron Schedule
Default: **Daily at 2:00 AM**

To change schedule, edit `cron/cleanup-security-events`:
```bash
0 3 * * 0    # Weekly on Sunday at 3:00 AM
0 1 1 * *    # Monthly on 1st day at 1:00 AM
*/5 * * * *  # Every 5 minutes (testing only)
```

## Rebuild Container

After any changes to cron files or Dockerfile:

```bash
cd deploy/compose
docker-compose down
docker-compose build php
docker-compose up -d
```

## Testing

### 1. Test Command Manually

```bash
# Dry run (shows what would be deleted)
docker exec fintech-pricing-system-php-1 php bin/console app:cleanup-security-events --dry-run

# Actual cleanup with custom retention
docker exec fintech-pricing-system-php-1 php bin/console app:cleanup-security-events --days=180

# Default cleanup (365 days)
docker exec fintech-pricing-system-php-1 php bin/console app:cleanup-security-events
```

### 2. Verify Cron Installation

```bash
# Check if cron is running
docker exec fintech-pricing-system-php-1 ps aux | grep cron

# View crontab
docker exec fintech-pricing-system-php-1 cat /etc/crontabs/www-data
```

### 3. Test with Short Interval

Temporarily edit `cron/cleanup-security-events`:
```bash
*/5 * * * * cd /var/www/html && php bin/console app:cleanup-security-events >> /var/log/cron.log 2>&1
```

Rebuild container, wait 5 minutes, then check logs.

### 4. Monitor Cron Logs

```bash
# View full log
docker exec fintech-pricing-system-php-1 cat /var/log/cron.log

# Tail logs in real-time
docker exec fintech-pricing-system-php-1 tail -f /var/log/cron.log

# From host machine (using volume)
docker volume inspect fintech-pricing-system_cron_logs
```

## Command Options

```bash
# Help
php bin/console app:cleanup-security-events --help

# Custom retention period
php bin/console app:cleanup-security-events --days=90

# Dry run (no deletion)
php bin/console app:cleanup-security-events --dry-run

# Combine options
php bin/console app:cleanup-security-events --days=180 --dry-run
```

## Troubleshooting

### Cron Not Running

```bash
# Check cron daemon status
docker exec fintech-pricing-system-php-1 ps aux | grep cron

# Check for errors
docker logs fintech-pricing-system-php-1

# Restart container
docker-compose restart php
```

### Permission Issues

```bash
# Check crontab permissions
docker exec fintech-pricing-system-php-1 ls -la /etc/crontabs/www-data
# Should be: -rw-r--r-- 1 root root

# Check log file
docker exec fintech-pricing-system-php-1 ls -la /var/log/cron.log
# Should be: -rw-r--r-- 1 www-data www-data
```

### Command Not Executing

```bash
# Run manually to see errors
docker exec fintech-pricing-system-php-1 php /var/www/html/bin/console app:cleanup-security-events

# Check Symfony logs
tail -f backend/var/log/dev.log
```

## Database Verification

```bash
# Count total security events
docker exec fintech-pricing-system-php-1 php bin/console doctrine:query:sql "SELECT COUNT(*) FROM security_events"

# View oldest events
docker exec fintech-pricing-system-php-1 php bin/console doctrine:query:sql "SELECT id, event_type, created_at FROM security_events ORDER BY created_at ASC LIMIT 10"

# View newest events
docker exec fintech-pricing-system-php-1 php bin/console doctrine:query:sql "SELECT id, event_type, created_at FROM security_events ORDER BY created_at DESC LIMIT 10"

# Count events older than 1 year
docker exec fintech-pricing-system-php-1 php bin/console doctrine:query:sql "SELECT COUNT(*) FROM security_events WHERE created_at < DATE_SUB(NOW(), INTERVAL 365 DAY)"
```

## Production Considerations

1. **Backup Before Cleanup**: Schedule database backups before cron runs
2. **Monitor Disk Space**: Track database size to adjust retention policy
3. **Compliance**: Verify retention period meets regulatory requirements
4. **Alerting**: Add notifications for cleanup failures
5. **Audit Trail**: Keep cron logs for at least 90 days

## Disable Cron

To disable scheduled cleanup without removing code:

1. Comment out the line in `cron/cleanup-security-events`:
   ```bash
   # 0 2 * * * cd /var/www/html && php bin/console app:cleanup-security-events >> /var/log/cron.log 2>&1
   ```

2. Rebuild container:
   ```bash
   docker-compose build php && docker-compose up -d
   ```

## Alternative: Manual Cleanup Only

If you prefer manual cleanup instead of scheduled:

1. Remove cron setup from Dockerfile
2. Run command manually when needed:
   ```bash
   docker exec fintech-pricing-system-php-1 php bin/console app:cleanup-security-events
   ```
