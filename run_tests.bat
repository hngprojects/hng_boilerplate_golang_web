
@echo off
REM Starting Docker Compose for the test PostgreSQL instance
echo Starting Docker Compose for the test PostgreSQL instance...
docker-compose -f docker-compose-test.yaml up -d

REM Check if Docker Compose started successfully
if %errorlevel% neq 0 (
    echo Failed to start Docker Compose.
    exit /b %errorlevel%
)

REM Run the Go tests
echo Running Go tests...
cd tests
go test -v ./...

REM Check if tests ran successfully
if %errorlevel% neq 0 (
    echo Go tests failed.
    REM Stop Docker Compose and clean up
    echo Stopping Docker Compose...
    docker-compose -f ../docker-compose-test.yaml down
    exit /b %errorlevel%
)

REM Change back to the original directory
cd ..

REM Stop Docker Compose and clean up
echo Stopping Docker Compose...
docker-compose -f docker-compose-test.yaml down

REM Check if Docker Compose stopped successfully
if %errorlevel% neq 0 (
    echo Failed to stop Docker Compose.
    exit /b %errorlevel%
)

echo All done!
