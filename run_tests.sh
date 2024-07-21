#!/bin/bash

# Starting Docker Compose for the test PostgreSQL instance
echo "Starting Docker Compose for the test PostgreSQL instance..."
docker-compose -f docker-compose-test.yaml up -d

# Check if Docker Compose started successfully
if [ $? -ne 0 ]; then
    echo "Failed to start Docker Compose."
    exit $?
fi

# Run the Go tests
echo "Running Go tests..."
cd tests
go test -v ./...

# Check if tests ran successfully
if [ $? -ne 0 ]; then
    echo "Go tests failed."
    # Stop Docker Compose and clean up
    echo "Stopping Docker Compose..."
    docker-compose -f ../docker-compose-test.yaml down
    exit $?
fi

# Change back to the original directory
cd ..

# Stop Docker Compose and clean up
echo "Stopping Docker Compose..."
docker-compose -f docker-compose-test.yaml down

# Check if Docker Compose stopped successfully
if [ $? -ne 0 ]; then
    echo "Failed to stop Docker Compose."
    exit $?
fi

echo "All done!"
