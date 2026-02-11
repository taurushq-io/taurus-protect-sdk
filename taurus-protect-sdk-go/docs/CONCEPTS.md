# Key Concepts - Go SDK

This document covers Go SDK-specific aspects of Taurus-PROTECT concepts. For the full domain model documentation, see the [Common Concepts](../../docs/CONCEPTS.md) document.

## Common Documentation

The following documentation applies to all Taurus-PROTECT SDKs:

- [Key Concepts](../../docs/CONCEPTS.md) - Domain model, entities, relationships
- [Authentication & TPV1](../../docs/AUTHENTICATION.md) - API authentication protocol
- [Integrity Verification](../../docs/INTEGRITY_VERIFICATION.md) - Cryptographic verification flows

---

## Go Model Types

The Go SDK represents domain entities as struct types in `pkg/protect/model`.

### Wallet

```go
type Wallet struct {
    ID                 string
    Name               string
    Currency           string
    Blockchain         string
    Network            string
    Balance            *Balance
    IsOmnibus          bool
    Disabled           bool
    CustomerID         string
    ExternalWalletID   string
    VisibilityGroupID  string
    AddressesCount     int64
    CreatedAt          time.Time
    UpdatedAt          time.Time
    Attributes         []WalletAttribute
}
```

### Address

```go
type Address struct {
    ID                           string
    WalletID                     string
    Address                      string
    AlternateAddress             string
    Label                        string
    Comment                      string
    Currency                     string
    CustomerID                   string
    ExternalAddressID            string
    AddressPath                  string
    AddressIndex                 string
    Nonce                        string
    Status                       string
    Balance                      *Balance
    Disabled                     bool
    CanUseAllFunds               bool
    LinkedWhitelistedAddressIDs  []string
    Attributes                   []AddressAttribute
    CreatedAt                    time.Time
    UpdatedAt                    time.Time
}
```

### Request

```go
type Request struct {
    ID                string
    TenantID          string
    Currency          string
    Status            string
    Type              string
    Rule              string
    RequestBundleID   string
    ExternalRequestID string
    NeedsApprovalFrom []string
    Metadata          *RequestMetadata
    Attributes        []RequestAttribute
    CreatedAt         time.Time
    UpdatedAt         time.Time
}
```

### RequestMetadata

```go
type RequestMetadata struct {
    Hash            string
    PayloadAsString string
}
```

**Note:** The `Payload` map field is intentionally omitted for security. Use `ParsePayloadEntries()` to extract structured data from `PayloadAsString`, which is the cryptographically verified source. Helper methods `GetSourceAddress()`, `GetDestinationAddress()`, `GetMetadataCurrency()`, and `GetAmount()` provide convenient typed access to common fields.

### Balance

```go
type Balance struct {
    TotalConfirmed       string
    TotalUnconfirmed     string
    AvailableConfirmed   string
    AvailableUnconfirmed string
    ReservedConfirmed    string
    ReservedUnconfirmed  string
}
```

### WhitelistedAddress

```go
type WhitelistedAddress struct {
    ID                 string
    TenantID           string
    Address            string
    Name               string
    Blockchain         string
    Network            string
    Status             string
    Action             string
    Rule               string
    RulesContainer     string
    RulesContainerHash string
    RulesSignatures    string
    VisibilityGroupID  string
    TnParticipantID    string
    Metadata           *WhitelistedAssetMetadata
    Scores             []Score
    Trails             []Trail
    Approvers          *Approvers
    SignedAddress      *SignedWhitelistedAddress
    Attributes         []WhitelistedAddressAttribute
}
```

### Pagination

The Go SDK uses two pagination patterns depending on the API endpoint.

#### Cursor-Based Pagination (Preferred)

Used by services like `RequestService`:

```go
// Cursor for cursor-based pagination
type Cursor struct {
    CurrentPage string
    PageSize    int64
    PageRequest string  // "FIRST", "NEXT", "PREVIOUS", "LAST"
}

// ResponseCursor returned with paginated results
type ResponseCursor struct {
    CurrentPage string
    HasPrevious bool
    HasNext     bool
}

// Result wrappers include cursor information
type RequestResult struct {
    Requests []*Request
    Cursor   *ResponseCursor
}
```

Usage:

```go
// Fetch first page
cursor := &model.Cursor{PageRequest: "FIRST", PageSize: 50}
result, err := client.Requests().ListRequests(ctx, nil, nil, "", nil, cursor)

// Iterate through pages
for result.Cursor != nil && result.Cursor.HasNext {
    cursor = &model.Cursor{
        CurrentPage: result.Cursor.CurrentPage,
        PageRequest: "NEXT",
        PageSize:    50,
    }
    result, err = client.Requests().ListRequests(ctx, nil, nil, "", nil, cursor)
}
```

#### Offset-Based Pagination (Legacy)

Used by some list endpoints that return `Pagination` metadata:

```go
type Pagination struct {
    Limit      int64
    Offset     int64
    TotalItems int64
    HasMore    bool
}
```

---

## Go Error Types

The SDK uses Go-style error types with wrapping support:

### APIError

```go
type APIError struct {
    Message     string
    Code        int
    Description string
    ErrorCode   string
    Err         error
    RetryAfter  time.Duration
}
```

**Methods:**
- `Error() string` - Error message
- `Unwrap() error` - Underlying error
- `IsRetryable() bool` - True for 429, 5xx
- `IsClientError() bool` - True for 4xx
- `IsServerError() bool` - True for 5xx

### Sentinel Errors

```go
var (
    ErrValidation     = &APIError{Code: 400, ...}
    ErrAuthentication = &APIError{Code: 401, ...}
    ErrAuthorization  = &APIError{Code: 403, ...}
    ErrNotFound       = &APIError{Code: 404, ...}
    ErrRateLimit      = &APIError{Code: 429, ...}
    ErrServer         = &APIError{Code: 500, ...}
)
```

### Specialized Errors

```go
type IntegrityError struct {
    Message string
    Err     error
}

type WhitelistError struct {
    Message string
    Err     error
}

type RequestMetadataError struct {
    Message string
    Err     error
}
```

### Error Checking

```go
import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"

// Check for specific error types
if apiErr, ok := protect.IsAPIError(err); ok {
    fmt.Printf("API Error: %d - %s\n", apiErr.Code, apiErr.Message)
}

if protect.IsIntegrityError(err) {
    // Cryptographic verification failed
}

if protect.IsWhitelistError(err) {
    // Whitelist-specific error
}

// Use errors.Is for sentinel errors
if errors.Is(err, protect.ErrNotFound) {
    // Resource not found
}
```

---

## Mapper Functions

The Go SDK uses mapper functions for DTO conversion in `pkg/protect/mapper`:

```go
// Wallet mapper
func WalletFromDTO(dto *openapi.TgvalidatordWalletInfo) *model.Wallet

// Address mapper
func AddressFromDTO(dto *openapi.TgvalidatordAddress) *model.Address

// Request mapper
func RequestFromDTO(dto *openapi.TgvalidatordRequest) *model.Request
```

### Safe Value Helpers

```go
func safeString(s *string) string     // Nil-safe string dereference
func safeBool(b *bool) bool           // Nil-safe bool dereference
func safeInt64(i *int64) int64        // Nil-safe int64 dereference
func stringPtr(s string) *string      // Create string pointer
```

---

## Context Usage

All service methods accept `context.Context` as the first parameter:

```go
ctx := context.Background()

// With timeout
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()

wallet, err := client.Wallets().GetWallet(ctx, walletID)
```

---

## Related Documentation

### Common (applies to all SDKs)
- [Key Concepts](../../docs/CONCEPTS.md) - Full domain model
- [Authentication](../../docs/AUTHENTICATION.md) - TPV1 protocol
- [Integrity Verification](../../docs/INTEGRITY_VERIFICATION.md) - Verification flows

### Go SDK Specific
- [SDK Overview](SDK_OVERVIEW.md) - Architecture and packages
- [Authentication](AUTHENTICATION.md) - Go authentication implementation
- [Services Reference](SERVICES.md) - Complete API documentation
- [Usage Examples](USAGE_EXAMPLES.md) - Code examples and patterns
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Go verification implementation

### Other SDKs
- [Java SDK Documentation](../../taurus-protect-sdk-java/docs/) - Java SDK reference
- [Python SDK Documentation](../../taurus-protect-sdk-python/docs/) - Python SDK reference
- [TypeScript SDK Documentation](../../taurus-protect-sdk-typescript/docs/) - TypeScript SDK reference