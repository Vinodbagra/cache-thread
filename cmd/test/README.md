# Cache API Test Runner

This directory contains a comprehensive test suite for the Cache API endpoints.

## Prerequisites

1. Make sure the cache server is running on `http://localhost:8080`
2. Ensure you have a `.env` file with the required configuration

## Running the Tests

### Option 1: Run directly with Go
```bash
cd cmd/test
go run test_apis.go
```

### Option 2: Build and run
```bash
cd cmd/test
go build -o test_apis test_apis.go
./test_apis
```

## What the Tests Cover

The test suite includes **13 comprehensive tests**:

1. **Health Check** - Verifies the API is running
2. **Get Configuration** - Tests configuration endpoint
3. **Put Single** - Tests storing a single key-value pair
4. **Get Single** - Tests retrieving a stored value
5. **Put with TTL** - Tests storing with expiration time
6. **Get Non-Existent** - Tests handling of missing keys
7. **Bulk Put** - Tests storing multiple key-value pairs
8. **Bulk Get** - Tests retrieving multiple values
9. **Get Stats** - Tests cache statistics endpoint
10. **List Keys** - Tests listing all cache keys
11. **Delete Key** - Tests removing a specific key
12. **Clear Cache** - Tests clearing the entire cache
13. **Get After Clear** - Verifies cache is empty after clearing

## Expected Output

When all tests pass, you should see output like:

```
🚀 Starting Cache API Tests...
Make sure the server is running on http://localhost:8080
============================================================

📋 Test 1: Health Check
✅ Health Check Passed - Status: 200
   Response: {"status":"healthy","timestamp":"2024-01-15T10:30:00Z","version":"1.0.0","uptime":"2h30m15s"}

[... more test results ...]

============================================================
📊 TEST RESULTS SUMMARY
============================================================
Total Tests: 13
Passed: 13 ✅
Failed: 0 ❌

Success Rate: 100.0%

🎉 All tests passed! Cache API is working correctly.
```

## Troubleshooting

### Server Not Running
If you get connection errors, make sure:
1. The server is started with `go run cmd/api/main.go`
2. The server is running on port 8080
3. No firewall is blocking the connection

### Test Failures
If tests fail:
1. Check the server logs for errors
2. Verify your `.env` configuration
3. Make sure the cache service is properly initialized

## Customizing Tests

You can modify the test file to:
- Change the base URL if your server runs on a different port
- Add more test cases
- Modify test data
- Add performance tests

## Test Data

The tests use various types of data:
- Simple strings
- JSON objects
- Arrays
- Different TTL values
- Edge cases (empty keys, non-existent keys)

This ensures comprehensive coverage of the API functionality. 