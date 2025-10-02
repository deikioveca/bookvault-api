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