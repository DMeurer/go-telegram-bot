cd /home/raspi/go-telegram-bot || echo "$(date --utc): Failed to cd into  /home/raspi/go-telegram-bot aborting..." >> "deployment/logs/cronjob.log" 2>&1

echo "$(date --utc): Starting Cron job..." >> "deployment/logs/cronjob.log" 2>&1

mkdir -p deployment/logs

LOCK_FILE="/tmp/go-telegram-bot.lockfile"

flock -n "$LOCK_FILE" deployment/check-updates.sh >> "deployment/logs/cronjob.log" 2>&1
