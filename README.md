# University Library Management System

This project is a University Library Management System built using the Gin framework in Go, designed to facilitate user registration, book management, and loan tracking.

## Project Structure

```
proyecto-bd-final
├── cmd
│   └── server
│       └── main.go          # Entry point of the application
├── internal
│   ├── config
│   │   └── database.go      # Database configuration
│   ├── controllers
│   │   ├── auth.go          # User authentication functions
│   │   ├── book.go          # Book management functions
│   │   └── loan.go          # Loan management functions
│   ├── middleware
│   │   └── auth.go          # Authentication middleware
│   ├── models
│   │   ├── user.go          # User model definition
│   │   ├── book.go          # Book model definition
│   │   └── loan.go          # Loan model definition
│   ├── repository
│   │   ├── user.go          # User repository functions
│   │   ├── book.go          # Book repository functions
│   │   └── loan.go          # Loan repository functions
│   ├── routes
│   │   └── routes.go        # Application routes
│   └── services
│       ├── auth.go          # User authentication logic
│       ├── book.go          # Book management logic
│       └── loan.go          # Loan management logic
├── pkg
│   └── utils
│       └── response.go      # Utility functions for API responses
├── migrations
│   ├── 001_create_users_table.sql  # SQL for creating Users table
│   ├── 002_create_books_table.sql  # SQL for creating Books table
│   └── 003_create_loans_table.sql  # SQL for creating Loans table
├── go.mod                     # Go module definition
├── go.sum                     # Module dependency checksums
└── README.md                  # Project documentation
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd proyecto-bd-final
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Configure the database:**
   Update the `internal/config/database.go` file with your Oracle database connection settings.

4. **Run migrations:**
   Execute the SQL scripts in the `migrations` folder to set up the database schema.

5. **Start the server:**
   ```
   go run cmd/server/main.go
   ```

## Functionality Overview

- **User Registration and Authentication:** Users can register and log in to the system.
- **Book Management:** Admins can add, update, and retrieve book information.
- **Loan Tracking:** Users can request loans for books and return them when done.

## License

This project is licensed under the MIT License.