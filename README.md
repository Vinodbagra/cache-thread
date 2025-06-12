# cache-thread
in-memory LRU cache

To set up this project, you need to have Go installed on your machine.

## Installation

First you have to clone the repository to your local machine:

then just run the following command to install the dependencies:

```go mod tidy```


## Usage

To run the application, use the following command:

```go run cmd/main.go```


This will start the application and listen on port 8080.

## Testing

To run the tests, use the following command:

```cd cmd/test```
```./run_tests.sh```

This will run all the tests and display the results.

You can also provide the max size of the cache and the default TTL in the env file.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

This project is licensed under the MIT License.