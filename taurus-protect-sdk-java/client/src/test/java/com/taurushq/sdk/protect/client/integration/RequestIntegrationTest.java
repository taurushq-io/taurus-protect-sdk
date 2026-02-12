package com.taurushq.sdk.protect.client.integration;

import com.google.gson.JsonElement;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestMetadataException;
import com.taurushq.sdk.protect.client.model.RequestResult;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for RequestService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class RequestIntegrationTest {

    private ProtectClient client;

    @BeforeAll
    void setup() throws Exception {
        TestHelper.skipIfNotEnabled();
        client = TestHelper.getTestClient(1);
    }

    @AfterAll
    void teardown() {
        if (client != null) {
            client.close();
        }
    }

    @Test
    void listRequests() throws ApiException {
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 10);
        RequestResult result = client.getRequestService().getRequests(null, null, null, null, cursor);

        List<Request> requests = result.getRequests();
        System.out.println("Found " + requests.size() + " requests");
        for (Request r : requests) {
            System.out.println("Request: ID=" + r.getId() + ", Status=" + r.getStatus());
        }

        assertNotNull(requests);
    }

    @Test
    void getRequest() throws ApiException {
        // First get a list to find a valid request ID
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 1);
        RequestResult result = client.getRequestService().getRequests(null, null, null, null, cursor);

        List<Request> requests = result.getRequests();
        if (requests.isEmpty()) {
            System.out.println("No requests available for testing");
            return;
        }

        long requestId = requests.get(0).getId();
        Request r = client.getRequestService().getRequest(requestId);

        System.out.println("Request details:");
        System.out.println("  ID: " + r.getId());
        System.out.println("  Status: " + r.getStatus());
        System.out.println("  METADATA PAYLOAD: " + r.getMetadata().getPayloadAsString());

        r.getSignedRequests().forEach(s -> System.out.println("  " + s.getHash() + ": " + s.getStatus()));

        assertNotNull(r);
        assertEquals(requestId, r.getId());
    }

    @Test
    void requestMetadataVerification() throws ApiException {
        // Get a request
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 1);
        RequestResult result = client.getRequestService().getRequests(null, null, null, null, cursor);

        List<Request> requests = result.getRequests();
        if (requests.isEmpty()) {
            System.out.println("No requests available for testing");
            return;
        }

        Request r = client.getRequestService().getRequest(requests.get(0).getId());

        // Verify we can access basic metadata
        try {
            System.out.println("Request ID: " + r.getMetadata().getRequestId());
        } catch (RequestMetadataException e) {
            System.out.println("Request ID: (not available)");
        }
        try {
            System.out.println("Currency: " + r.getMetadata().getCurrency());
        } catch (RequestMetadataException e) {
            System.out.println("Currency: (not available)");
        }
        try {
            System.out.println("Source Address: " + r.getMetadata().getSourceAddress());
        } catch (RequestMetadataException e) {
            System.out.println("Source Address: (not available)");
        }
        try {
            System.out.println("Destination Address: " + r.getMetadata().getDestinationAddress());
        } catch (RequestMetadataException e) {
            System.out.println("Destination Address: (not available - some request types don't have destinations)");
        }

        // Verify metadata payload can be parsed as JSON
        String payloadStr = r.getMetadata().getPayloadAsString();
        JsonElement root = JsonParser.parseString(payloadStr);
        assertTrue(root.isJsonArray(), "Metadata payload should be a JSON array");

        // Verify we can iterate through metadata fields
        int fieldCount = 0;
        for (JsonElement bm : root.getAsJsonArray()) {
            Map<String, JsonElement> map = bm.getAsJsonObject().asMap();
            String key = map.get("key").getAsString();
            System.out.println("  Metadata field: " + key);
            fieldCount++;
        }

        System.out.println("Total metadata fields: " + fieldCount);
        assertTrue(fieldCount > 0, "Request should have metadata fields");
    }
}
