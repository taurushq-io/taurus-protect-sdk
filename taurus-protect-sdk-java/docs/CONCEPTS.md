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
| `AuthorizationException` | Checked (extends `ApiException`) | Authorization failures (HTTP 403); `getRequiredRoles()` names the roles that would satisfy the check |
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

Every list method follows the cross-SDK pagination contract (repo-root `CLAUDE.md` →
"Pagination (cross-SDK)"): a page size is always sent, `0`/`null` means
`Pagination.DEFAULT_PAGE_SIZE` (**20**), and a negative size or one above
`Pagination.MAX_PAGE_SIZE` (**100**) throws `IllegalArgumentException` naming the option before
any request is made. A negative offset is rejected the same way.

Every paged result carries a pagination value that is never null on a successful call; an empty
reply (`{}`) is a valid last page.

### Offset lists — `OffsetPagination`

Wallets, addresses, transactions, users, groups, fee payers, actions and the two whitelists page
by offset (validatord has no cursor for them). Their result types (`WalletResult`,
`AddressResult`, `TransactionResult`, `UserResult`, `GroupResult`, `FeePayerResult`,
`ActionResult`, `WhitelistedAddressListResult`, `WhitelistedAssetResult`) expose
`getPagination()`:

```java
public final class OffsetPagination {
    int getLimit();          // the page size sent
    long getOffset();        // the offset sent (never the reply's)
    long getTotalItems();    // server total, reduced by rows the SDK withheld
    long getNextOffset();    // pass this back as the next offset
    boolean hasMore();       // offset < nextOffset < server total
}
```

```java
long offset = 0;
WalletResult page;
do {
    page = client.getWalletService().getWallets(20, offset);
    for (Wallet w : page.getWallets()) {
        // process
    }
    offset = page.getPagination().getNextOffset();
} while (page.getPagination().hasMore());
```

Always continue from `getNextOffset()`, never from `offset + rows`: each endpoint has its own
rule (the wallet and address replies return the next offset; users and groups can append a
synthetic row beyond the limit; the contract whitelist keeps skipped rows' slots), and
`OffsetPagination.of` applies it.

### Cursor lists — `CursorPage`

Every other list pages by cursor. Each cursor method takes a page size and a cursor — the
previous page's `getPage().getNextCursor()`, or `null` for the first page — and its result
(`RequestResult`, `ChangeResult`, `BalanceResult`, `PriceResult`, ...) extends
`CursorPagedResult`:

```java
public final class CursorPage {
    int getPageSize();       // the page size sent
    String getNextCursor();  // "" when there is no next page
    boolean hasMore();       // the reply cursor's hasNext
    Long getTotalItems();    // the server total where the endpoint reports one, else null
}
```

```java
String cursor = null;
ChangeResult page;
do {
    page = client.getChangeService().getChanges(null, null, 20, cursor);
    // process page.getChanges()
    cursor = page.getPage().getNextCursor();
} while (page.getPage().hasMore());
```

With a cursor the SDK sends `currentPage=<cursor>`, `pageRequest=NEXT` and the page size, in
whatever form the endpoint takes (`cursor.*` or `requestCursor.*` query parameters, or a body
cursor). A cursor is opaque base64 text and may contain `+ / =`; it is URL-encoded on the wire.
Wallet tokens and the governance rules history page by a bare token instead; the result's
`getPage()` is the same `CursorPage`, and their `getTotalItems()` is set.

### Low-level request cursors

Each cursor method that existed before the contract also keeps an `ApiRequestCursor` overload,
which can express every `PageRequest` (`FIRST`, `PREVIOUS`, `NEXT`, `LAST`). A `null` cursor is the
first page with the default size; an explicit size must be between 1 and 100.

```java
// Enum for page navigation direction
public enum PageRequest { FIRST, PREVIOUS, NEXT, LAST }

ApiRequestCursor first = Pagination.first(20);                    // pageRequest=FIRST
ApiRequestCursor contract = Pagination.page(20, nextCursor);      // the contract form
ApiRequestCursor back = Pagination.previous(result.getCursor(), 20);
```

`TaurusNetworkSharingService.listSharedAddresses` has seven filters, which leaves no room for
separate page-size and cursor parameters, so it takes the page only as
`Pagination.page(pageSize, cursor)`.

### Exempt: price history and the transaction export

`PriceService.getPriceHistory` returns at most `Pagination.MAX_PRICE_HISTORY_LIMIT` (365) daily
points (default 20). `TransactionService.exportTransactions` cannot page at all: the server ignores
any offset and always exports from the first matching transaction, so it takes only a limit
(default 20, no SDK maximum) and returns a `TransactionExportResult` with the exported text and the
server's `getTotalItems()`.

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

## AuthorizationException.getRequiredRoles()

A 403 is mapped to `AuthorizationException`, whose `getRequiredRoles()` returns an
unmodifiable list.

```java
try {
    client.getWalletService().createWallet(/* ... */);
} catch (AuthorizationException e) {
    List<String> roles = e.getRequiredRoles();
    System.out.println(roles.isEmpty()
            ? "Permission denied"
            : "Permission denied; requires one of: " + String.join(", ", roles));
}
```

The error carries the **roles that would satisfy the failed check**, so a caller can say
which role to request instead of only "forbidden". The list is empty when the denial was
not role-based (a disabled endpoint, a visibility restriction). One entry means that role
is required; several mean any one of them suffices.

> Surface the role list, never the server's message. The role names are safe to display;
> the surrounding text is server-controlled. A consumer that forwards errors into another
> system should render roles through an allowlist (the roles are lowercase alphanumeric).
