echo "Starting" >> "deployment/logs/cronjob.log" 2>&1

cd /home/raspi/go-telegram-bot

pwd  >> "deployment/logs/cronjob.log" 2>&1

echo "Switched" >> "deployment/logs/cronjob.log" 2>&1

echo "$(date --utc): Starting Cron job..." >> "deployment/logs/cronjob.log" 2>&1

echo "written" >> "deployment/logs/cronjob.log" 2>&1

mkdir -p deployment/logs

echo "made dir" >> "deployment/logs/cronjob.log" 2>&1

LOCK_FILE="$(pwd)/deployment/cronjob.lock" >> "deployment/logs/cronjob.log" 2>&1

echo "define lockfile" >> "deployment/logs/cronjob.log" 2>&1

sudo flock -n "$LOCK_FILE" "deployment/check-updates.sh" >> "deployment/logs/cronjob.log" 2>&1

echo "done" >> "deployment/logs/cronjob.log" 2>&1

ls deployment/logs/
