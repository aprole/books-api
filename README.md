# Books API

This is an HTTP API built with Go and Gorilla mux for managing a collection of books.

## Installation

1. Install Go (https://golang.org/doc/install)
2. Clone the repository: `git clone https://github.com/aprole/books-api.git`
3. Build the binary: `go build`

## Usage

1. Start the server: `./books-api`
2. Use a tool like cURL or Postman to make requests to the API endpoints.

## Endpoints

### Books

- `GET /api/books` - Get all books
- `GET /api/books/{id}` - Get a book by ID
- `POST /api/books` - Create a new book
- `PUT /api/books/{id}` - Update a book by ID
- `DELETE /api/books/{id}` - Delete a book by ID

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](https://choosealicense.com/licenses/mit/)