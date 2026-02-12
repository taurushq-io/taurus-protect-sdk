package integration

import (
	"context"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestIntegration_ListRequests(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()
	result, err := client.Requests().ListRequests(ctx, &model.ListRequestsOptions{
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("ListRequests() error = %v", err)
	}

	t.Logf("Found %d requests, HasNext: %v", len(result.Requests), result.HasNext)

	for _, r := range result.Requests {
		t.Logf("Request: ID=%s, Status=%s, Currency=%s", r.ID, r.Status, r.Currency)
	}
}

func TestIntegration_GetRequest(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	// First, get a list to find a valid request ID
	ctx := context.Background()
	result, err := client.Requests().ListRequests(ctx, &model.ListRequestsOptions{
		PageSize: 1,
	})
	if err != nil {
		t.Fatalf("ListRequests() error = %v", err)
	}
	if len(result.Requests) == 0 {
		t.Skip("No requests available for testing")
	}

	requestID := result.Requests[0].ID
	request, err := client.Requests().GetRequest(ctx, requestID)
	if err != nil {
		t.Fatalf("GetRequest(%s) error = %v", requestID, err)
	}

	t.Logf("Request details:")
	t.Logf("  ID: %s", request.ID)
	t.Logf("  Status: %s", request.Status)
	t.Logf("  Currency: %s", request.Currency)
	t.Logf("  Type: %s", request.Type)
	if request.Metadata != nil {
		t.Logf("  Metadata hash: %s", request.Metadata.Hash)
	}
}

func TestIntegration_RequestMetadataVerification(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// Get a request
	result, err := client.Requests().ListRequests(ctx, &model.ListRequestsOptions{
		PageSize: 1,
	})
	if err != nil {
		t.Fatalf("ListRequests() error = %v", err)
	}
	if len(result.Requests) == 0 {
		t.Skip("No requests available for testing")
	}

	request, err := client.Requests().GetRequest(ctx, result.Requests[0].ID)
	if err != nil {
		t.Fatalf("GetRequest() error = %v", err)
	}

	// Verify we can access metadata
	if request.Metadata == nil {
		t.Skip("Request has no metadata")
	}

	t.Logf("Request ID: %s", request.ID)
	t.Logf("Metadata hash: %s", request.Metadata.Hash)
	t.Logf("Metadata PayloadAsString length: %d", len(request.Metadata.PayloadAsString))

	// Verify metadata payload is present and non-empty
	if request.Metadata.PayloadAsString == "" {
		t.Fatal("Metadata PayloadAsString should not be empty")
	}

	// Log metadata payload (first 200 chars for readability)
	payload := request.Metadata.PayloadAsString
	if len(payload) > 200 {
		t.Logf("Metadata payload (truncated): %s...", payload[:200])
	} else {
		t.Logf("Metadata payload: %s", payload)
	}

	// Verify basic metadata fields
	t.Logf("Request Currency: %s", request.Currency)
	t.Logf("Request Type: %s", request.Type)
	t.Logf("Request Status: %s", request.Status)

	// Test the convenience methods for extracting payload data
	entries, err := request.Metadata.ParsePayloadEntries()
	if err != nil {
		t.Fatalf("ParsePayloadEntries() error = %v", err)
	}
	t.Logf("Parsed %d payload entries", len(entries))
	for _, entry := range entries {
		t.Logf("  Entry: key=%s", entry.Key)
	}

	// Test GetMetadataCurrency
	metadataCurrency := request.Metadata.GetMetadataCurrency()
	t.Logf("Metadata currency: %s", metadataCurrency)

	// Test GetMetadataRequestID
	metadataRequestID := request.Metadata.GetMetadataRequestID()
	t.Logf("Metadata request ID: %d", metadataRequestID)

	// Test GetSourceAddress (may be empty for some request types)
	sourceAddr := request.Metadata.GetSourceAddress()
	if sourceAddr != "" {
		t.Logf("Source address: %s", sourceAddr)
	}

	// Test GetDestinationAddress (may be empty for some request types)
	destAddr := request.Metadata.GetDestinationAddress()
	if destAddr != "" {
		t.Logf("Destination address: %s", destAddr)
	}

	// Test GetAmount (may be nil for some request types)
	amount := request.Metadata.GetAmount()
	if amount != nil {
		t.Logf("Amount: valueFrom=%s, valueTo=%s, rate=%s, currencyFrom=%s, currencyTo=%s",
			amount.ValueFrom, amount.ValueTo, amount.Rate, amount.CurrencyFrom, amount.CurrencyTo)
	}
}

func TestIntegration_ListRequestsByStatus(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()

	ctx := context.Background()

	// List requests with CONFIRMED status (the final successful state)
	// Note: API uses "CONFIRMED" not "COMPLETED" for completed requests
	confirmedResult, err := client.Requests().ListRequests(ctx, &model.ListRequestsOptions{
		PageSize: 10,
		Status:   "CONFIRMED",
	})
	if err != nil {
		t.Fatalf("ListRequests(CONFIRMED) error = %v", err)
	}

	t.Logf("Found %d CONFIRMED requests, HasNext: %v", len(confirmedResult.Requests), confirmedResult.HasNext)

	for _, r := range confirmedResult.Requests {
		t.Logf("Request: ID=%s, Status=%s, Type=%s", r.ID, r.Status, r.Type)
		// Verify all returned requests are CONFIRMED
		if r.Status != "CONFIRMED" && r.Status != "confirmed" {
			t.Errorf("Expected CONFIRMED status, got %s", r.Status)
		}
	}

	// Also try PENDING status
	pendingResult, err := client.Requests().ListRequests(ctx, &model.ListRequestsOptions{
		PageSize: 5,
		Status:   "PENDING",
	})
	if err != nil {
		t.Fatalf("ListRequests(PENDING) error = %v", err)
	}

	t.Logf("Found %d PENDING requests, HasNext: %v", len(pendingResult.Requests), pendingResult.HasNext)
}
