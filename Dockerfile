# syntax=docker/dockerfile:1

FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# copy project root to /app
COPY . .

# install go dependencys
RUN go mod download

# Builds your app with optional configuration
# RUN go build -o /godocker

# Tells Docker which network port your container listens on
EXPOSE 8080
