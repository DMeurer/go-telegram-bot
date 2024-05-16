cd ../

LOCK_FILE="$(pwd)/deployment/cronjob.lock"

flock -n "$LOCK_FILE" ./deployment/check-updates.sh
