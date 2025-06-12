package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const baseURL = "http://localhost:8080/api/cache"

// TestResults holds the results of API tests
type TestResults struct {
	TotalTests         int
	PassedTests        int
	FailedTests        int
	FailedTestsDetails []string
}

func main() {
	fmt.Println("🚀 Starting Cache API Tests...")
	fmt.Println("Make sure the server is running on http://localhost:8080")
	fmt.Println(strings.Repeat("=", 60))

	results := &TestResults{}

	// Test 1: Health Check
	testHealthCheck(results)

	// Test 2: Get Configuration
	testGetConfiguration(results)

	// Test 3: Put single key-value
	testPutSingle(results)

	// Test 4: Get single key-value
	testGetSingle(results)

	// Test 5: Put with TTL
	testPutWithTTL(results)

	// Test 6: Get non-existent key
	testGetNonExistent(results)

	// Test 7: Bulk Put
	testBulkPut(results)

	// Test 8: Bulk Get
	testBulkGet(results)

	// Test 9: Get Stats
	testGetStats(results)

	// Test 10: List Keys
	testListKeys(results)

	// Test 11: Delete key
	testDeleteKey(results)

	// Test 12: Clear cache
	testClearCache(results)

	// Test 13: Get after clear (should be empty)
	testGetAfterClear(results)

	// Print final results
	printResults(results)
}

func testHealthCheck(results *TestResults) {
	fmt.Println("\n📋 Test 1: Health Check")

	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		failTest(results, "Health Check", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		failTest(results, "Health Check", fmt.Sprintf("Expected 200, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Health Check Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testGetConfiguration(results *TestResults) {
	fmt.Println("\n📋 Test 2: Get Configuration")

	resp, err := http.Get(baseURL + "/config")
	if err != nil {
		failTest(results, "Get Configuration", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		failTest(results, "Get Configuration", fmt.Sprintf("Expected 200, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Get Configuration Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testPutSingle(results *TestResults) {
	fmt.Println("\n📋 Test 3: Put Single Key-Value")

	data := map[string]interface{}{
		"key":   "test:user:1",
		"value": map[string]interface{}{"name": "John Doe", "age": 30},
	}

	jsonData, _ := json.Marshal(data)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", baseURL+"/put", bytes.NewBuffer(jsonData))
	if err != nil {
		failTest(results, "Put Single", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		failTest(results, "Put Single", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		failTest(results, "Put Single", fmt.Sprintf("Expected 201, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Put Single Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testGetSingle(results *TestResults) {
	fmt.Println("\n📋 Test 4: Get Single Key-Value")

	resp, err := http.Get(baseURL + "/get/test:user:1")
	if err != nil {
		failTest(results, "Get Single", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		failTest(results, "Get Single", fmt.Sprintf("Expected 200, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Get Single Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testPutWithTTL(results *TestResults) {
	fmt.Println("\n📋 Test 5: Put with TTL")

	data := map[string]interface{}{
		"key":   "test:temp:1",
		"value": "This will expire in 5 seconds",
		"ttl":   5, // 5 seconds
	}

	jsonData, _ := json.Marshal(data)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", baseURL+"/put", bytes.NewBuffer(jsonData))
	if err != nil {
		failTest(results, "Put with TTL", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		failTest(results, "Put with TTL", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		failTest(results, "Put with TTL", fmt.Sprintf("Expected 201, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Put with TTL Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testGetNonExistent(results *TestResults) {
	fmt.Println("\n📋 Test 6: Get Non-Existent Key")

	resp, err := http.Get(baseURL + "/get/non:existent:key")
	if err != nil {
		failTest(results, "Get Non-Existent", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		failTest(results, "Get Non-Existent", fmt.Sprintf("Expected 404, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Get Non-Existent Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testBulkPut(results *TestResults) {
	fmt.Println("\n📋 Test 7: Bulk Put")

	data := map[string]interface{}{
		"items": []map[string]interface{}{
			{
				"key":   "bulk:user:1",
				"value": map[string]interface{}{"name": "Alice", "role": "admin"},
				"ttl":   3600,
			},
			{
				"key":   "bulk:user:2",
				"value": map[string]interface{}{"name": "Bob", "role": "user"},
				"ttl":   1800,
			},
			{
				"key":   "bulk:user:3",
				"value": map[string]interface{}{"name": "Charlie", "role": "moderator"},
			},
		},
	}

	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(baseURL+"/bulk/put", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		failTest(results, "Bulk Put", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		failTest(results, "Bulk Put", fmt.Sprintf("Expected 200, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Bulk Put Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testBulkGet(results *TestResults) {
	fmt.Println("\n📋 Test 8: Bulk Get")

	data := map[string]interface{}{
		"keys": []string{"bulk:user:1", "bulk:user:2", "bulk:user:3", "non:existent:key"},
	}

	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(baseURL+"/bulk/get", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		failTest(results, "Bulk Get", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		failTest(results, "Bulk Get", fmt.Sprintf("Expected 200, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Bulk Get Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testGetStats(results *TestResults) {
	fmt.Println("\n📋 Test 9: Get Stats")

	resp, err := http.Get(baseURL + "/stats")
	if err != nil {
		failTest(results, "Get Stats", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		failTest(results, "Get Stats", fmt.Sprintf("Expected 200, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Get Stats Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testListKeys(results *TestResults) {
	fmt.Println("\n📋 Test 10: List Keys")

	resp, err := http.Get(baseURL + "/keys?limit=10")
	if err != nil {
		failTest(results, "List Keys", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		failTest(results, "List Keys", fmt.Sprintf("Expected 200, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ List Keys Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testDeleteKey(results *TestResults) {
	fmt.Println("\n📋 Test 11: Delete Key")

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", baseURL+"/delete/bulk:user:1", nil)
	if err != nil {
		failTest(results, "Delete Key", err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		failTest(results, "Delete Key", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		failTest(results, "Delete Key", fmt.Sprintf("Expected 200, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Delete Key Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testClearCache(results *TestResults) {
	fmt.Println("\n📋 Test 12: Clear Cache")

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", baseURL+"/clear", nil)
	if err != nil {
		failTest(results, "Clear Cache", err.Error())
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		failTest(results, "Clear Cache", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		failTest(results, "Clear Cache", fmt.Sprintf("Expected 200, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Clear Cache Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func testGetAfterClear(results *TestResults) {
	fmt.Println("\n📋 Test 13: Get After Clear")

	resp, err := http.Get(baseURL + "/get/test:user:1")
	if err != nil {
		failTest(results, "Get After Clear", err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		failTest(results, "Get After Clear", fmt.Sprintf("Expected 404, got %d", resp.StatusCode))
		return
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✅ Get After Clear Passed - Status: %d\n", resp.StatusCode)
	fmt.Printf("   Response: %s\n", string(body))
	passTest(results)
}

func passTest(results *TestResults) {
	results.TotalTests++
	results.PassedTests++
}

func failTest(results *TestResults, testName, reason string) {
	results.TotalTests++
	results.FailedTests++
	results.FailedTestsDetails = append(results.FailedTestsDetails, fmt.Sprintf("%s: %s", testName, reason))
}

func printResults(results *TestResults) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 TEST RESULTS SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total Tests: %d\n", results.TotalTests)
	fmt.Printf("Passed: %d ✅\n", results.PassedTests)
	fmt.Printf("Failed: %d ❌\n", results.FailedTests)

	if results.FailedTests > 0 {
		fmt.Println("\n❌ Failed Tests:")
		for _, detail := range results.FailedTestsDetails {
			fmt.Printf("   - %s\n", detail)
		}
	}

	successRate := float64(results.PassedTests) / float64(results.TotalTests) * 100
	fmt.Printf("\nSuccess Rate: %.1f%%\n", successRate)

	if results.FailedTests == 0 {
		fmt.Println("\n🎉 All tests passed! Cache API is working correctly.")
	} else {
		fmt.Println("\n⚠️  Some tests failed. Please check the server logs and try again.")
	}
}
