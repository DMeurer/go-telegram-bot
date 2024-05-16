cd ../

LOCK_FILE="$(pwd)/deployment/cronjob.lock"

flock -n "$LOCK_FILE" ./deployment/check-updates.sh >> ./deployment/logs/cronjob.log 2>&1
