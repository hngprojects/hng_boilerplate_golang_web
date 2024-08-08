# Build stage
FROM golang:1.20.1-alpine3.17 as build

# Set the Current Working Directory inside the container
WORKDIR /usr/src/app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download && go mod verify

# Copy the source code into the container
COPY . .

# Build the Go app
RUN if test -e app.env; then echo 'found app.env'; else mv app-sample.env app.env; fi; \
    go build -v -o /dist/app-name

# Wait-for-it stage
FROM alpine:3.17 as wait
RUN apk add --no-cache bash
ADD https://github.com/vishnubob/wait-for-it/raw/master/wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

# Deployment stage
FROM alpine:3.17
WORKDIR /usr/src/app
COPY --from=build /usr/src/app ./
COPY --from=build /dist/app-name /usr/local/bin/app-name
COPY --from=wait /wait-for-it.sh /wait-for-it.sh

# Install bash (required for wait-for-it script)
RUN apk add --no-cache bash

# Wait for DB and Redis, then start the application
CMD /wait-for-it.sh $DB_HOST:$DB_PORT -t 10 -- /wait-for-it.sh $REDIS_HOST:$REDIS_PORT -t 10 -- app-name
