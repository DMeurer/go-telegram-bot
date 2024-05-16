# pull the latest changes from the repository
git pull --force

# rebuilding the docker compose, with zero downtime
docker compose up -d --no-deps --build

# remove the old images
docker image prune -f

# remove the old containers
docker image prune -f
