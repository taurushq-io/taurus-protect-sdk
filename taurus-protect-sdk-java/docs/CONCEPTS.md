# Key Concepts - Java SDK

This document covers Java SDK-specific aspects of Taurus PROTECT concepts. For the full domain model documentation, see the [Common Concepts](../../docs/CONCEPTS.md) document.

## Common Documentation

The following documentation applies to all Taurus PROTECT SDKs:

- [Key Concepts](../../docs/CONCEPTS.md) - Domain model, entities, relationships
- [Authentication & TPV1](../../docs/AUTHENTICATION.md) - API authentication protocol
- [Integrity Verification](../../docs/INTEGRITY_VERIFICATION.md) - Cryptographic verification flows

---

## Java Model Classes

The Java SDK represents domain entities as immutable model classes in `com.taurushq.sdk.protect.client.model`.

### Wallet

```java
public class Wallet {
    private Long id;
    private String name;
    private String blockchain;
    private String network;
    private Boolean isOmnibus;
    private String customerId;
    private Balance balance;
    private List<Attribute> attributes;
    // getters...
}
```

### Address

```java
public class Address {
    private Long id;
    private Long walletId;
    private String address;
    private String label;
    private String customerId;
    private Balance balance;
    private List<AddressAttribute> attributes;
    // getters...
}
```

### Request

```java
public class Request {
    private Long id;
    private RequestStatus status;
    private String currency;
    private RequestMetadata metadata;
    private List<RequestApprover> approvers;
    private List<SignedRequest> signedRequests;
    private List<RequestTrail> trails;
    // getters...
}
```

### RequestMetadata

```java
public class RequestMetadata {
    // Stored fields:
    private String hash;              // Cryptographic hash for integrity verification
    private String payloadAsString;   // JSON string (cryptographically verified source)
    private JsonElement payloadAsJson; // Parsed from payloadAsString for data extraction

    // Convenience extraction methods (parse from the verified payload):
    public long getRequestId() throws RequestMetadataException { ... }
    public String getSourceAddress() throws RequestMetadataException { ... }
    public String getDestinationAddress() throws RequestMetadataException { ... }
    public RequestMetadataAmount getAmount() throws RequestMetadataException { ... }
    public String getCurrency() throws RequestMetadataException { ... }
    public String getRulesKey() throws RequestMetadataException { ... }
    // getters...
}
```

**Security Note:** The extraction methods (`getSourceAddress()`, `getDestinationAddress()`, `getAmount()`, etc.) parse from `payloadAsJson`, which is derived from `payloadAsString` -- the cryptographically verified source. The raw `payload` field is intentionally omitted to prevent extraction from an unverified source.

### Balance

```java
public class Balance {
    private BigInteger totalConfirmed;
    private BigInteger totalUnconfirmed;
    private BigInteger availableConfirmed;
    private BigInteger availableUnconfirmed;
    private BigInteger reservedConfirmed;
    private BigInteger reservedUnconfirmed;
    // getters...
}
```

---

## Java Enums

### RequestStatus

```java
public enum RequestStatus {
    CREATED,
    PENDING,
    APPROVING,
    APPROVED,
    APPROVED_2,
    REJECTED,
    HSM_READY,
    HSM_READY_2,
    HSM_SIGNED,
    HSM_SIGNED_2,
    HSM_FAILED,
    HSM_FAILED_2,
    BROADCASTING,
    BROADCASTING_2,
    BROADCASTED,
    MINED,
    CONFIRMED,
    PARTIALLY_CONFIRMED,
    PERMANENT_FAILURE,
    CANCELED,
    EXPIRED,
    READY,
    SENT,
    MANUAL_BROADCAST,
    AUTO_PREPARED,
    AUTO_PREPARED_2,
    FAST_APPROVED_2,
    BUNDLE_READY,
    BUNDLE_APPROVED,
    BUNDLE_BROADCASTING,
    UNKNOWN
    // ... plus DIEM_* and SIGNET_* statuses
}
```

---

## Java Exception Types

The SDK uses specific exceptions for different error scenarios:

| Exception | Type | When Thrown |
|-----------|------|-------------|
| `ApiException` | Checked | General API errors (network, auth, validation) |
| `AuthenticationException` | Checked (extends `ApiException`) | Authentication failures (HTTP 401) |
| `AuthorizationException` | Checked (extends `ApiException`) | Authorization failures (HTTP 403) |
| `NotFoundException` | Checked (extends `ApiException`) | Resource not found (HTTP 404) |
| `RateLimitException` | Checked (extends `ApiException`) | Rate limit exceeded (HTTP 429) |
| `ServerException` | Checked (extends `ApiException`) | Server errors (HTTP 5xx) |
| `ValidationException` | Checked (extends `ApiException`) | Input validation errors (HTTP 400) |
| `IntegrityException` | Unchecked (extends `SecurityException`) | Hash/signature verification failures |
| `WhitelistException` | Checked | Whitelist-specific verification errors |
| `ConfigurationException` | Checked | Client configuration errors |
| `RequestMetadataException` | Checked | Metadata payload parsing errors |

### Example

```java
try {
    WhitelistedAddress addr = client.getWhitelistedAddressService()
        .getWhitelistedAddress(id);
} catch (IntegrityException e) {
    // Hash mismatch or insufficient valid signatures
    System.err.println("Verification failed: " + e.getMessage());
} catch (WhitelistException e) {
    // Whitelist-specific verification failure
    System.err.println("Whitelist error: " + e.getMessage());
} catch (ApiException e) {
    // General API error
    System.err.println("API error: " + e.getCode() + " - " + e.getMessage());
}
```

---

## Pagination

The Java SDK uses cursor-based pagination with a factory pattern.

### Core Types

```java
// Enum for page navigation direction
public enum PageRequest { FIRST, PREVIOUS, NEXT, LAST }

// Request cursor (sent to API)
public class ApiRequestCursor {
    private String currentPage;    // Page token (null for initial request)
    private PageRequest pageRequest; // Navigation direction
    private long pageSize;          // Items per page
}

// Response cursor (returned from API)
public class ApiResponseCursor {
    private String currentPage;    // Current page token
    private Boolean hasPrevious;   // Has previous page
    private Boolean hasNext;       // Has next page

    public ApiRequestCursor nextPage(int pageSize);
    public ApiRequestCursor previousPage(int pageSize);
}
```

### Pagination Factory

The `Pagination` class provides static factory methods for creating cursors:

```java
import com.taurushq.sdk.protect.client.model.Pagination;

// First page (default size: 50, max: 1000)
ApiRequestCursor cursor = Pagination.first();
ApiRequestCursor cursor = Pagination.first(100);

// Navigate from response cursor
ApiRequestCursor next = Pagination.next(responseCursor, 50);
ApiRequestCursor prev = Pagination.previous(responseCursor, 50);

// Last page
ApiRequestCursor last = Pagination.last(50);
```

### Result Wrapper Pattern

All paginated operations return `*Result` classes wrapping the data with cursor info:

```java
// Fetch first page
ApiRequestCursor cursor = Pagination.first(50);
BalanceResult result = client.getBalanceService().getBalances("ETH", cursor);

// Iterate through pages
while (result.hasNext()) {
    cursor = result.nextCursor(50);
    result = client.getBalanceService().getBalances("ETH", cursor);
    for (AssetBalance balance : result.getBalances()) {
        // process each balance
    }
}
```

---

## MapStruct Mappers

The Java SDK uses [MapStruct](https://mapstruct.org/) for compile-time DTO mapping:

```java
@Mapper
public interface WalletMapper {
    WalletMapper INSTANCE = Mappers.getMapper(WalletMapper.class);

    Wallet fromDTO(TgvalidatordWallet dto);
    Wallet fromDTO(TgvalidatordWalletInfo dto);
}
```

Generated implementations are in `target/generated-sources/`.

---

## Related Documentation

### Common (applies to all SDKs)
- [Key Concepts](../../docs/CONCEPTS.md) - Full domain model
- [Authentication](../../docs/AUTHENTICATION.md) - TPV1 protocol
- [Integrity Verification](../../docs/INTEGRITY_VERIFICATION.md) - Verification flows

### Java SDK Specific
- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Authentication](AUTHENTICATION.md) - Java authentication implementation
- [Services Reference](SERVICES.md) - Complete API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples and patterns
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Java verification implementation

### Other SDKs
- [Go SDK Documentation](../../taurus-protect-sdk-go/docs/) - Go SDK reference
- [Python SDK Documentation](../../taurus-protect-sdk-python/docs/) - Python SDK reference
- [TypeScript SDK Documentation](../../taurus-protect-sdk-typescript/docs/) - TypeScript SDK reference