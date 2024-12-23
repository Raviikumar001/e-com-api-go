# E-commerce RBAC API Documentation

## Table of Contents

1. [Authentication](#authentication)
2. [Roles & Permissions](#roles--permissions)

```
project-root/
├── main.go                 # Application entry point
│
│
├── internal/                   # Private application code
│   ├── config/                # Configuration code
│   │   └── config.go
│   │
│   ├── database/              # Database related code
│   │   ├── db.go             # Database connection
│   │   └── seed.go           # Database seeding
│   │
│   ├── handlers/              # HTTP request handlers
│   │   ├── auth_handler.go
│   │   ├── product_handler.go
│   │   └── storefront_handler.go
│   │
│   ├── middleware/            # Middleware components
│   │   ├── auth_middleware.go
│   │   ├── error_handler.go
│   │   └── rbac_middleware.go
│   │
│   ├── models/                # Data models
│   │   ├── base.go
│   │   ├── constants.go
│   │   ├── order.go
│   │   ├── permission.go
│   │   ├── product.go
│   │   ├── role.go
│   │   ├── storefront.go
│   │   └── user.go
│   │
│   ├── pkg/                   # Internal packages
│   │   └── response.go
│   │
│   ├── routes/               # Route definitions
│   │   └── routes.go
│   │
│   └── utils/               # Utility functions
│       ├── jwt.go
│       ├── password.go
│       └── token.go
│
├── .env                     # Environment variables
├── .env.example            # Example environment variables
├── .gitignore
├── go.mod                  # Go module file
├── go.sum                  # Go module checksum
└── README.md              # Project
```

## documentation

## Authentication

All API calls require authentication token except login and register.
Add token to request header:

```
Authorization: Bearer your_jwt_token
```

**Login**

POST /auth/login
Content-Type: application/json

```
Request:
{
    "email": "user@example.com",
    "password": "password123"
}

Response:
{
    "message": "Login successful",
    "token": "eyJhbGciOiJIUzI1..."
}
```

Register

POST /auth/register
Content-Type: application/json

```
Request:
{
    "email": "newuser@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe",
    "role_id": 3  // Role ID for seller
}

Response:
{
    "message": "Registration successful",
    "user": {
        "id": 1,
        "email": "newuser@example.com",
        "first_name": "John",
        "last_name": "Doe",
        "role_id": 3
    }
}
```

Roles & Permissions
Role IDs

```
1 - Admin
2 - Wholesaler
3 - Seller
4 - Customer
```

Role Permissions

```
Admin (ID: 1)

Full access to all endpoints
Can view all user details
Can manage all products
Can see cost prices and wholesale information

Wholesaler (ID: 2)

Can manage wholesale products
Can set wholesale prices
Can view wholesale-specific information
Cannot view retail customer data

Seller (ID: 3)

Can manage own products
Can view own product details
Can update/delete own products
Cannot view other sellers' data

Customer (ID: 4)

Can view published products
Limited product details access
Cannot view wholesale prices
```
