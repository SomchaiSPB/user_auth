## README.md

# User Authentication and Product Management Service

This project provides a RESTful API for user authentication and product management, implemented in Go. The application supports multiple storage options and is containerized using Docker for easy deployment.

## Features

- **User Registration**: Allows the creation of new users with unique usernames and passwords.
- **User Authentication**: Authenticates users and provides JWT tokens for session management.
- **Product Retrieval**: Fetches product details by name or lists all products, with support for pagination.

## API Endpoints

### User Endpoints

- **Create User**
  - **URL**: `/users`
  - **Method**: `POST`
  - **Request Body**: `CreateUserDTO`
  - **Response**: `User`
  - **Description**: Creates a new user. Returns an error if the username already exists.

- **Authenticate User**
  - **URL**: `/auth`
  - **Method**: `POST`
  - **Request Body**: `AuthUserRequestDTO`
  - **Response**: `AuthUserResponseDTO` (JWT Token)
  - **Description**: Authenticates a user and returns a JWT token.

### Product Endpoints

- **Get Product**
  - **URL**: `/products/{name}`
  - **Method**: `GET`
  - **Response**: `Product`
  - **Description**: Retrieves product details by name, with case-insensitive search.

- **Get Products**
  - **URL**: `/api/v1/products`
  - **Method**: `GET`
  - **Response**: `[]Product`
  - **Description**: Retrieves a list of all products.

## Default Data

The application comes with default data for testing purposes:

- **Default User**:
  - **Username**: `admin@admin.com`
  - **Password**: `12345678`

- **Default Product**:
  - **Name**: `test`

You can use these credentials and product to test the API endpoints.

## Storage Options

The application supports two storage options:

1. **PostgreSQL**:
  - Used when the application is run in Docker.
  - Ideal for production environments.
2. **SQLite**:
  - Available for both Docker and binary deployments.
  - Suitable for development, testing, or small-scale deployments.

You can select the storage type using the `APP_STORAGE` environment variable. The default storage type is PostgreSQL, but it can be changed to SQLite. If SQLite is used, running the PostgreSQL database via Docker Compose is not mandatory.

## Configuration

An example configuration is provided in the `.env.example` file. Copy this file to `.env` and adjust the settings as needed:

```env
APP_HTTP_PORT=6543
APP_ENV=local
APP_STORAGE=postgres  # storage types: postgres, sqlite
APP_WITH_FAKE_DATA=true
APP_WITH_TABLE_TRUNCATE=true

DB_SQLITE_FILE=db.sqlite3

AUTH_JWT_SECRET=test123

DB_HOST=db  # use 'db' for Docker, otherwise configure as needed
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=user_auth
```

### Key Environment Variables

- **APP_HTTP_PORT**: The port on which the application will listen.
- **APP_ENV**: The application environment (`local`, `production`, etc.).
- **APP_STORAGE**: The storage type (`postgres` or `sqlite`).
- **APP_WITH_FAKE_DATA**: Whether to populate the database with fake data.
- **APP_WITH_TABLE_TRUNCATE**: Whether to truncate tables on startup.
- **DB_SQLITE_FILE**: The filename for SQLite storage.
- **AUTH_JWT_SECRET**: The secret key for signing JWT tokens.
- **DB_* Variables**: Configuration for PostgreSQL connection.

## Running the Application

### Using Docker Compose

To start the application and the PostgreSQL database, run:

```bash
docker-compose up -d --build
```

This command builds and starts the containers, making the application accessible on the configured port. PostgreSQL is initialized as specified in the `.env` file.

### Running Locally

For local development with SQLite:

1. Set `APP_STORAGE=sqlite` in your `.env` file.
2. Run the application binary.

## API Documentation

The HTTP API documentation can be viewed at:

```
http://localhost:6543/swagger/index.html#/
```

Replace `6543` with the value of `APP_HTTP_PORT` if it has been changed.

## TODO

### Testing

1. **Unit Tests**: Implement unit tests for each service and handler.
2. **Integration Tests**: Verify interactions between the application and the database.
3. **End-to-End Tests**: Test real user scenarios, including user registration, authentication, and product retrieval.
4. **Load Testing**: Ensure the application can handle expected traffic levels.

### Future Enhancements

- Implement refresh token functionality for JWT tokens.
- Add roles and permissions for more granular access control.
- Expand product management features to include CRUD operations.
- Enhance logging and monitoring for better observability.