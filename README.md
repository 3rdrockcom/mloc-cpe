# Customer Profiling Engine
CPE allows a user to feed customer transactions to the system in order to generate a classification and other useful statisitics.

## Instructions
```bash
# Perform Database Migrations
./goose -dir=./migrations mysql "username:password@tcp(localhost:3306)/cpe?parseTime=true" up;

# Build application
go build -o cpe;

# Run application
./cpe -dsn="username:password@tcp(localhost:3306)/cpe?parseTime=true";
```

## API Documentation
API documentation is available on Postman (Collection: CPE).