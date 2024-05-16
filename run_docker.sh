docker build -t go_tel_bot .
docker run -dp 8080:8080 go_tel_bot
echo "Container is running on port 8080"
