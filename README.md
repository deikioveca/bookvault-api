BookVault API
- 
bookvault-api is a RESTful backend service for managing an online bookstore.
It provides functionality for user management, book catalog management, shopping cart operations, order processing, and book reviews.
The API is designed to be modular, secure, and scalable, making it suitable for building a complete e-commerce platform for books.

Features
- 
* User Management:
  * Register and authenticate users

  * Role-based access (admin and user)

  * Manage user details
* Book Management:

  * Add, view, and update books

  * Search books by title or author

  * Manage stock
* Shopping Cart:

  * Add and remove books from the cart

  * Update book quantities

  * Clear cart
* Orders:

  * Create and cancel orders

  * View orders by user or status

  * Update order status
* Reviews:

  * Add, update, and delete reviews for books

  * Fetch reviews by book or by user
 

Technologies Used
-
* Go
* GORM
* PostgreSQL
* Docker

Running Guide
- 
* Prerequisites

  * Go 1.25+ installed (for local testing)

  * Docker & Docker Compose (for containerized setup)

  * PostgreSQL (for local setup)

Running locally
-
* clone the repo
* go mod download
* start postgresql and ensure databases bookvault and bookvault_test exist.
* change .env file according to your setup
* go run main.go

Running with Docker
-
* docker compose up -> this command will start two services: app(bookvault-api) and db(PostgreSQL)
* to run tests inside the container you must open a shell inside the app container with the command -> docker compose exec app sh
* to test services: cd tests/services -> go test -v
* to test handlers: cd tests/handlers -> go test -v
