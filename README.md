# FL Weather Service

## Overview
This project is a web service that combines two existing APIs:
1. Fetches a random location from [Patch3s Locations API](https://locations.patch3s.dev/api/random).
2. Retrieves the current weather forecast for the fetched location 
3. Combines the results and returns them in a simple response.

## Implementation Stages

### Stage 1: Initial Working Version
- Implemented using Go's standard `net/http` package.
- Successfully fetched and combined data from both APIs.
- Basic logging and error handling included.
- Achieved within **less than 2 hours**.

### Stage 2: Production-Ready Enhancements
- **Switched to [Chi Router](https://github.com/go-chi/chi)** for better request handling and middleware support.
- **Reorganized project structure** following best practices for maintainability.
- **Added Config from environment variables** for better configuration management.
- **Replaced standard HTTP client with [Resty](https://github.com/go-resty/resty)** to handle retries and enhanced error logging.
- **Implemented caching** to reduce redundant API calls and handle traffic spikes.
- **Added unit tests with Gomock** for robust testing.
- **Implemented graceful shutdown** to ensure smooth service termination.
- Taken **more than 4 hours** for this stage.

## Running the Service
### Prerequisites
- Go 1.18+ 

### Running Locally
```sh
# Clone the repository
git clone <repo-url>
cd <repo-folder>

# Run the service
go run cmd/main.go
```


### Testing the API
```sh
curl http://localhost:5000
```

## Future Improvements
- Add rate limiting for API protection.
- Introduce monitoring and observability tools.
- Implement authentication and authorization if needed.

## Author
Sothyvorn

