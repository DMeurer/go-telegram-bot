cd /home/raspi/go-telegram-bot || exit

mkdir -p deployment/logs

LOCK_FILE="$(pwd)/deployment/cronjob.lock"

flock -n "$LOCK_FILE" "$(pwd)/deployment/check-updates.sh" >> "$(pwd)/deployment/logs/cronjob.log" 2>&1
