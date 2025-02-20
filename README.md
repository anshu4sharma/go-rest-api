# Go REST API with Authentication

A scalable and secure REST API built with Go, featuring JWT authentication, Google OAuth integration, and role-based authorization.

## Features

- Clean Architecture pattern
- JWT Authentication
- Google OAuth2 Integration
- Role-based Authorization
- MySQL Database with GORM
- Middleware for Authentication and Authorization
- Secure Password Hashing
- Environment Configuration
- Structured Error Handling

## Project Structure

- `cmd/`: Main entry point for the application
- `internal/`: Contains all the core logic of the application
- `pkg/`: External packages
- `tests/`: Test files
- `Dockerfile`: Docker configuration
- `Makefile`: Build and run commands

## Installation

1. Clone the repository

```bash
git clone https://github.com/yourusername/go-rest-api.git
```

2. Install dependencies

```bash
go mod tidy
```

3. Create a `.env` file in the root directory with the following variables:

```bash

```

4. Build the application

```bash
go build -o go-rest-api
```

5. Run the application

```bash
./go-rest-api
```
# go-rest-api
