# Usage Examples

This document provides complete TypeScript code examples for common SDK operations.

## Table of Contents

1. [Client Setup](#client-setup)
2. [Wallet Management](#wallet-management)
3. [Address Management](#address-management)
4. [Request Approval Workflow](#request-approval-workflow)
5. [Transaction Queries](#transaction-queries)
6. [Balance Queries](#balance-queries)
7. [Whitelisted Addresses](#whitelisted-addresses)
8. [Governance Rules](#governance-rules)
9. [TaurusNetwork Operations](#taurusnetwork-operations)
10. [Webhook Management](#webhook-management)
11. [Pagination Patterns](#pagination-patterns)
12. [Error Handling](#error-handling)
13. [Async/Await Patterns](#asyncawait-patterns)

---

## Client Setup

### Basic Initialization

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function createClient(): Promise<ProtectClient> {
  const client = ProtectClient.create({
    host: 'https://api.protect.taurushq.com',
    apiKey: 'your-api-key-uuid',
    apiSecret: 'your-api-secret-hex',
  });

  return client;
}

// Usage
const client = await createClient();
try {
  // Use client...
} finally {
  client.close();
}
```

### Environment-Based Configuration

```typescript
import { ProtectClient, ProtectClientConfig } from '@taurushq/protect-sdk';

function createFromEnvironment(): ProtectClient {
  const host = process.env.TAURUS_API_HOST;
  const apiKey = process.env.TAURUS_API_KEY;
  const apiSecret = process.env.TAURUS_API_SECRET;

  if (!host || !apiKey || !apiSecret) {
    throw new Error('Missing required environment variables');
  }

  const config: ProtectClientConfig = {
    host,
    apiKey,
    apiSecret,
  };

  // Optional: SuperAdmin keys for verification
  if (process.env.TAURUS_SUPERADMIN_KEYS) {
    config.superAdminKeysPem = JSON.parse(process.env.TAURUS_SUPERADMIN_KEYS);
    config.minValidSignatures = parseInt(process.env.TAURUS_MIN_SIGNATURES ?? '2', 10);
  }

  return ProtectClient.create(config);
}
```

### With SuperAdmin Key Verification

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

const superAdmin1 = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----`;

const superAdmin2 = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----`;

const client = ProtectClient.create({
  host: 'https://api.protect.taurushq.com',
  apiKey: 'your-api-key-uuid',
  apiSecret: 'your-api-secret-hex',
  superAdminKeysPem: [superAdmin1, superAdmin2],
  minValidSignatures: 2,
  rulesCacheTtlMs: 600000, // 10 minutes
});
```

---

## Wallet Management

### Create a Wallet

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function createWalletExample(client: ProtectClient): Promise<void> {
  // Basic wallet creation
  const wallet = await client.wallets.create({
    blockchain: 'ETH',
    network: 'mainnet',
    name: 'Treasury',
    isOmnibus: false,
  });

  console.log(`Created wallet: ${wallet.id}`);
  console.log(`Name: ${wallet.name}`);

  // Wallet with additional metadata
  const walletWithMeta = await client.wallets.create({
    blockchain: 'ETH',
    network: 'mainnet',
    name: 'Operations Wallet',
    isOmnibus: false,
    comment: 'Used for daily operations',
    customerId: 'CUST-12345',
  });

  console.log(`Created wallet with customer ID: ${walletWithMeta.id}`);
}
```

### List and Search Wallets

```typescript
import { ProtectClient, Wallet } from '@taurushq/protect-sdk';

async function listWalletsExample(client: ProtectClient): Promise<void> {
  // List first page of wallets
  const result = await client.wallets.list({ limit: 50, offset: 0 });
  console.log(`Found ${result.pagination.totalItems} wallets`);

  for (const wallet of result.items) {
    console.log(`${wallet.id}: ${wallet.name} (${wallet.blockchain})`);
  }

  // List all wallets with pagination
  const allWallets: Wallet[] = [];
  let offset = 0;
  const limit = 50;

  while (true) {
    const page = await client.wallets.list({ limit, offset });
    allWallets.push(...page.items);

    if (page.items.length < limit) {
      break; // No more pages
    }
    offset += limit;
  }

  console.log(`Total wallets loaded: ${allWallets.length}`);

  // Filter by blockchain
  const ethWallets = await client.wallets.list({
    blockchain: 'ETH',
    excludeDisabled: true,
    onlyPositiveBalance: true,
  });

  console.log(`ETH wallets with balance: ${ethWallets.items.length}`);
}
```

### Get Wallet Balance History

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function walletBalanceHistoryExample(
  client: ProtectClient,
  walletId: number
): Promise<void> {
  // Get hourly balance history
  const history = await client.wallets.getBalanceHistory(walletId, 1);

  for (const point of history) {
    console.log(`${point.pointDate}: ${point.balance?.totalConfirmed}`);
  }

  // Get daily balance history
  const dailyHistory = await client.wallets.getBalanceHistory(walletId, 24);
  console.log(`Daily data points: ${dailyHistory.length}`);
}
```

### Manage Wallet Attributes

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function walletAttributesExample(
  client: ProtectClient,
  walletId: number
): Promise<void> {
  // Add attributes
  await client.wallets.createAttribute(walletId, 'department', 'Finance');
  await client.wallets.createAttribute(walletId, 'costCenter', 'CC-100');

  // Get wallet to see attributes
  const wallet = await client.wallets.get(walletId);
  for (const attr of wallet.attributes ?? []) {
    console.log(`${attr.key}: ${attr.value}`);
  }
}
```

---

## Address Management

### Create Addresses

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function createAddressExample(
  client: ProtectClient,
  walletId: number
): Promise<void> {
  // Create address with object-style request
  const address = await client.addresses.create({
    walletId: String(walletId),
    label: 'Customer Deposit',
    comment: 'Auto-generated for customer deposits',
    customerId: 'USER-789',
  });

  console.log(`Address: ${address.address}`);
  console.log(`ID: ${address.id}`);

  // Create address with explicit parameters
  const address2 = await client.addresses.createAddress(
    walletId,
    'Hot Wallet Address',
    'Primary hot wallet',
    'CUST-001'
  );

  console.log(`Created: ${address2.address}`);
}
```

### Get Address with Verification

When SuperAdmin keys are configured, addresses are automatically verified:

```typescript
import { ProtectClient, IntegrityError } from '@taurushq/protect-sdk';

async function getAddressExample(
  client: ProtectClient,
  addressId: number
): Promise<void> {
  try {
    // This performs signature verification automatically
    const address = await client.addresses.get(addressId);

    console.log(`Address: ${address.address}`);
    console.log(`Balance: ${address.balance?.totalConfirmed}`);
    console.log(`Signature: ${address.signature}`);
  } catch (error) {
    if (error instanceof IntegrityError) {
      console.error('Address signature verification failed!');
      console.error('Data may have been tampered with.');
    }
    throw error;
  }
}
```

### List Addresses

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function listAddressesExample(
  client: ProtectClient,
  walletId: number
): Promise<void> {
  // List addresses for a wallet
  const result = await client.addresses.list(walletId, {
    limit: 50,
    offset: 0,
  });

  console.log(`Addresses: ${result.items.length}`);
  if (result.pagination) {
    console.log(`Total: ${result.pagination.totalItems}`);
  }

  for (const addr of result.items) {
    console.log(`${addr.id}: ${addr.address} (${addr.label})`);
  }

  // List with additional filters
  const filtered = await client.addresses.listWithOptions({
    walletId: String(walletId),
    blockchain: 'ETH',
    network: 'mainnet',
    limit: 100,
  });

  console.log(`Filtered addresses: ${filtered.items.length}`);
}
```

### Add Custom Attributes

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function addAttributeExample(
  client: ProtectClient,
  addressId: number
): Promise<void> {
  await client.addresses.createAttribute(addressId, 'department', 'Finance');
  await client.addresses.createAttribute(addressId, 'purpose', 'Customer deposits');

  console.log('Attributes added successfully');
}
```

---

## Request Approval Workflow

### Complete Approval Flow

```typescript
import { ProtectClient, decodePrivateKeyPem, IntegrityError } from '@taurushq/protect-sdk';

async function completeApprovalFlow(client: ProtectClient): Promise<void> {
  // Step 1: Create an external transfer request
  const request = await client.requests.createExternalTransferRequest({
    fromAddressId: 1001,
    toWhitelistedAddressId: 2001,
    amount: '1000000000000000000', // 1 ETH in wei
    comment: 'Monthly treasury transfer',
  });

  console.log(`Created request ID: ${request.id}`);
  console.log(`Status: ${request.status}`);

  // Step 2: Retrieve and verify request (hash verification is automatic)
  const retrieved = await client.requests.get(request.id);

  // Step 3: Review metadata before signing
  const metadata = retrieved.metadata;
  if (metadata) {
    console.log('=== Review Request Metadata ===');
    console.log(`Request ID: ${retrieved.id}`);
    console.log(`Currency: ${retrieved.currency}`);
    console.log(`Hash: ${metadata.hash}`);
    console.log(`Payload: ${metadata.payloadAsString}`);
    // The payload object contains the parsed transfer details
    if (metadata.payload) {
      console.log(`Parsed payload: ${JSON.stringify(metadata.payload, null, 2)}`);
    }
  }

  // Step 4: Approve with private key (after human review)
  const privateKeyPem = await loadPrivateKey(); // From secure storage
  const privateKey = decodePrivateKeyPem(privateKeyPem);

  const signatures = await client.requests.approveRequest(
    retrieved,
    privateKey,
    'Approved after compliance review'
  );

  console.log(`Signatures performed: ${signatures}`);

  // Step 5: Check updated status
  const updated = await client.requests.get(request.id);
  console.log(`New status: ${updated.status}`);
}

async function loadPrivateKey(): Promise<string> {
  // Load from secure storage (HSM, vault, etc.)
  return process.env.USER_PRIVATE_KEY ?? '';
}
```

### Batch Approval

```typescript
import { ProtectClient, decodePrivateKeyPem, Request } from '@taurushq/protect-sdk';

async function batchApprovalExample(client: ProtectClient): Promise<void> {
  // Get all requests pending approval
  const { requests } = await client.requests.listForApproval({ limit: 100 });

  // Filter by currency
  const ethRequests = requests.filter(r => r.currency === 'ETH');

  if (ethRequests.length === 0) {
    console.log('No ETH requests pending approval');
    return;
  }

  // Load private key
  const privateKeyPem = process.env.USER_PRIVATE_KEY ?? '';
  const privateKey = decodePrivateKeyPem(privateKeyPem);

  // Batch approve
  const count = await client.requests.approveRequests(
    ethRequests,
    privateKey,
    'Batch approval for Q1 transfers'
  );

  console.log(`Approved ${ethRequests.length} requests with ${count} signatures`);
}
```

### Create Transfer Requests

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function createTransferRequests(client: ProtectClient): Promise<void> {
  // Internal transfer (between addresses in same tenant)
  const internalRequest = await client.requests.createInternalTransferRequest({
    fromAddressId: 100,
    toAddressId: 200,
    amount: '500000000000000000', // 0.5 ETH
    comment: 'Internal rebalancing',
  });
  console.log(`Internal request: ${internalRequest.id}`);

  // External transfer (to whitelisted address)
  const externalRequest = await client.requests.createExternalTransferRequest({
    fromAddressId: 100,
    toWhitelistedAddressId: 300,
    amount: '1000000000000000000', // 1 ETH
    comment: 'External payment',
    destinationAddressMemo: 'Invoice #12345',
  });
  console.log(`External request: ${externalRequest.id}`);

  // Transfer from wallet (omnibus)
  const walletRequest = await client.requests.createInternalTransferFromWalletRequest({
    fromWalletId: 1,
    toAddressId: 200,
    amount: '2000000000000000000', // 2 ETH
    gasLimit: '21000',
  });
  console.log(`Wallet request: ${walletRequest.id}`);
}
```

### Reject Requests

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function rejectRequestExample(client: ProtectClient): Promise<void> {
  // Reject single request
  await client.requests.rejectRequest(123, 'Amount exceeds daily limit');

  // Reject multiple requests
  await client.requests.rejectRequests(
    [101, 102, 103],
    'Suspicious activity detected'
  );

  console.log('Requests rejected');
}
```

---

## Transaction Queries

### Query Transactions

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function queryTransactionsExample(client: ProtectClient): Promise<void> {
  // Get transaction by ID
  const tx = await client.transactions.get(12345);
  console.log(`Hash: ${tx.hash}`);
  console.log(`Status: ${tx.status}`);
  console.log(`Amount: ${tx.amount}`);

  // Get transaction by hash
  const txByHash = await client.transactions.getByHash(
    '0xe41578e07623a4a3646cf6393a512d975adbc4d6446849148d8c742069dfb34f'
  );
  console.log(`Found: ${txByHash.id}`);

  // List transactions with filters
  const now = new Date();
  const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);

  const result = await client.transactions.list({
    fromDate: weekAgo,
    toDate: now,
    currency: 'ETH',
    limit: 100,
  });

  console.log(`Found ${result.items.length} transactions`);
  for (const t of result.items) {
    console.log(`${t.id}: ${t.hash?.substring(0, 10)}... (${t.status})`);
  }
}
```

---

## Balance Queries

### Get All Balances

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function getAllBalancesExample(client: ProtectClient): Promise<void> {
  const result = await client.balances.list({ limit: 100 });

  console.log('=== Asset Balances ===');
  for (const balance of result.items) {
    console.log(`${balance.currency}: ${balance.available} available`);
  }

  // Summarize by currency
  const byCurrency = new Map<string, string>();
  for (const balance of result.items) {
    byCurrency.set(balance.currency, balance.available);
  }

  console.log('\n=== Summary ===');
  for (const [currency, amount] of byCurrency) {
    console.log(`${currency}: ${amount}`);
  }
}
```

---

## Whitelisted Addresses

### Get Whitelisted Address

```typescript
import { ProtectClient, WhitelistError } from '@taurushq/protect-sdk';

async function getWhitelistedAddressExample(client: ProtectClient): Promise<void> {
  try {
    // Simple: get address directly
    const addr = await client.whitelistedAddresses.get('1001');

    console.log(`Blockchain: ${addr.blockchain}`);
    console.log(`Address: ${addr.address}`);
    console.log(`Label: ${addr.label}`);
  } catch (error) {
    if (error instanceof WhitelistError) {
      console.error('Whitelist verification failed:', error.message);
    }
    throw error;
  }
}
```

### Get Whitelisted Address with Full Verification

For enhanced security, use the verification methods:

```typescript
import { ProtectClient, WhitelistedAddressService, rulesContainerFromBase64, userSignaturesFromBase64 } from '@taurushq/protect-sdk';

async function verifiedWhitelistExample(): Promise<void> {
  // Create client
  const client = ProtectClient.create({
    host: 'https://api.protect.taurushq.com',
    apiKey: 'your-api-key',
    apiSecret: 'your-api-secret',
  });

  // Create verified service with decoders
  const verifiedService = WhitelistedAddressService.withVerification(
    client.addressWhitelistingApi,
    {
      superAdminKeysPem: [superAdmin1, superAdmin2],
      minValidSignatures: 2,
      rulesContainerDecoder: rulesContainerFromBase64,
      userSignaturesDecoder: userSignaturesFromBase64,
    }
  );

  // Get address with full 6-step verification
  const result = await verifiedService.getWithVerification('1001');

  console.log('Verification passed!');
  console.log(`Address: ${result.verifiedWhitelistedAddress.address}`);
  console.log(`Hash verified: ${result.hashVerified}`);
  console.log(`Signatures verified: ${result.signaturesVerified}`);
}
```

### List Whitelisted Addresses

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function listWhitelistedAddressesExample(client: ProtectClient): Promise<void> {
  // List all
  const all = await client.whitelistedAddresses.list({ limit: 100 });
  console.log(`Total: ${all.pagination?.totalItems ?? all.items.length}`);

  // Filter by blockchain
  const ethAddresses = await client.whitelistedAddresses.list({
    blockchain: 'ETH',
    limit: 100,
  });

  // Filter by blockchain and network
  const mainnetAddresses = await client.whitelistedAddresses.list({
    blockchain: 'ETH',
    network: 'mainnet',
    limit: 100,
  });

  for (const addr of mainnetAddresses.items) {
    console.log(`${addr.blockchain}: ${addr.address} (${addr.label})`);
  }
}
```

---

## Governance Rules

### Get and Verify Governance Rules

```typescript
import { ProtectClient, IntegrityError } from '@taurushq/protect-sdk';

async function governanceRulesExample(client: ProtectClient): Promise<void> {
  try {
    // Get current rules (verification happens if configured)
    const rules = await client.governanceRules.get();

    if (rules) {
      console.log(`Rules ID: ${rules.id}`);
      console.log(`Locked: ${rules.locked}`);
      console.log(`Signatures: ${rules.rulesSignatures.length}`);
      console.log(`Created: ${rules.creationDate}`);
    }
  } catch (error) {
    if (error instanceof IntegrityError) {
      console.error('Governance rules verification failed!');
    }
    throw error;
  }
}
```

### Decode Rules Container

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function decodeRulesExample(client: ProtectClient): Promise<void> {
  // Get decoded rules container
  const container = await client.governanceRules.getDecodedRulesContainer();

  console.log('=== Decoded Rules Container ===');
  console.log(`Users: ${container.users.length}`);
  console.log(`Groups: ${container.groups.length}`);

  // List users
  for (const user of container.users) {
    console.log(`User: ${user.email} (ID: ${user.id})`);
  }
}
```

### Get Rules History

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function rulesHistoryExample(client: ProtectClient): Promise<void> {
  // Get first page
  const page1 = await client.governanceRules.getHistory({ limit: 10 });
  console.log(`Found ${page1.items.length} rule sets`);

  for (const rules of page1.items) {
    console.log(`${rules.id}: created ${rules.creationDate}`);
  }

  // Get next page if available
  if (page1.nextCursor) {
    const page2 = await client.governanceRules.getHistory({
      limit: 10,
      cursor: page1.nextCursor,
    });
    console.log(`Page 2: ${page2.items.length} rule sets`);
  }
}
```

---

## TaurusNetwork Operations

TaurusNetwork provides low-level API access for Taurus Network operations. The following examples show how to use the raw OpenAPI-generated APIs.

### Participant Management

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function participantExample(client: ProtectClient): Promise<void> {
  // Get my participant info
  const myParticipantResponse = await client.taurusNetwork.participantApi.getMyParticipant();
  const myParticipant = myParticipantResponse.result;
  console.log(`My participant: ${myParticipant?.participant?.name}`);
  console.log(`Country: ${myParticipant?.participant?.country}`);

  if (myParticipant?.settings) {
    console.log(`Status: ${myParticipant.settings.status}`);
  }

  // List all visible participants
  const participantsResponse = await client.taurusNetwork.participantApi.getAllParticipants();
  const participants = participantsResponse.result?.participants ?? [];
  for (const p of participants) {
    console.log(`${p.name}: ${p.country}`);
  }

  // Get specific participant with valuation
  const participantResponse = await client.taurusNetwork.participantApi.getParticipant({
    participantId: 'participant-id',
    includeTotalPledgesValuation: true,
  });
  const participant = participantResponse.result?.participant;
  console.log(`Outgoing pledges: ${participant?.outgoingTotalPledgesValuationBaseCurrency}`);
}
```

### Pledge Operations

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function pledgeExample(client: ProtectClient): Promise<void> {
  // List all pledges
  const pledgesResponse = await client.taurusNetwork.pledgeApi.getAllPledges({ pageSize: 50 });
  const pledges = pledgesResponse.result?.pledges ?? [];
  console.log(`Found ${pledges.length} pledges`);

  for (const pledge of pledges) {
    console.log(`${pledge.id}: ${pledge.status}`);
  }

  // Get specific pledge
  const pledgeResponse = await client.taurusNetwork.pledgeApi.getPledge({ pledgeId: 'pledge-id' });
  const pledge = pledgeResponse.result?.pledge;
  console.log(`Pledge amount: ${pledge?.pledgedAmount}`);
  console.log(`Asset: ${pledge?.currencyInfo?.symbol}`);
}
```

### Lending Operations

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function lendingExample(client: ProtectClient): Promise<void> {
  // List lending offers
  const offersResponse = await client.taurusNetwork.lendingApi.getAllLendingOffers({ pageSize: 50 });
  const offers = offersResponse.result?.lendingOffers ?? [];
  console.log(`Found ${offers.length} offers`);

  // List lending agreements
  const agreementsResponse = await client.taurusNetwork.lendingApi.getAllLendingAgreements({ pageSize: 50 });
  const agreements = agreementsResponse.result?.lendingAgreements ?? [];
  console.log(`Found ${agreements.length} agreements`);
}
```

### Settlement Operations

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function settlementExample(client: ProtectClient): Promise<void> {
  // List settlements
  const settlementsResponse = await client.taurusNetwork.settlementApi.getAllSettlements({ pageSize: 50 });
  const settlements = settlementsResponse.result?.settlements ?? [];
  console.log(`Found ${settlements.length} settlements`);

  for (const settlement of settlements) {
    console.log(`${settlement.id}: ${settlement.status}`);
  }
}
```

### Sharing Operations

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function sharingExample(client: ProtectClient): Promise<void> {
  // List shared address assets
  const sharedResponse = await client.taurusNetwork.sharedAddressAssetApi.getAllSharedAddressAssets({ pageSize: 50 });
  const sharedAssets = sharedResponse.result?.sharedAddressAssets ?? [];
  console.log(`Found ${sharedAssets.length} shared assets`);

  for (const asset of sharedAssets) {
    console.log(`${asset.id}: ${asset.currencyInfo?.symbol}`);
  }
}
```

---

## Webhook Management

### Create and Configure Webhooks

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function webhookExample(client: ProtectClient): Promise<void> {
  // Create a webhook
  const webhook = await client.webhooks.create({
    url: 'https://example.com/webhook/transactions',
    eventType: 'TRANSACTION',
    secret: 'my-webhook-secret-key',
  });
  console.log(`Created webhook: ${webhook.id}`);

  // List all webhooks
  const webhooks = await client.webhooks.list();
  for (const wh of webhooks.items) {
    console.log(`${wh.id}: ${wh.eventType} - ${wh.status}`);
  }

  // Update webhook status
  await client.webhooks.updateStatus(webhook.id, 'DISABLED');
  console.log('Webhook disabled');

  await client.webhooks.updateStatus(webhook.id, 'ENABLED');
  console.log('Webhook re-enabled');

  // Delete webhook
  await client.webhooks.delete(webhook.id);
  console.log('Webhook deleted');
}
```

---

## Pagination Patterns

### Offset-Based Pagination

```typescript
import { ProtectClient, Wallet } from '@taurushq/protect-sdk';

async function offsetPaginationExample(client: ProtectClient): Promise<void> {
  const pageSize = 50;
  let offset = 0;
  const allWallets: Wallet[] = [];

  while (true) {
    const result = await client.wallets.list({ limit: pageSize, offset });
    allWallets.push(...result.items);

    console.log(`Page ${offset / pageSize + 1}: ${result.items.length} items`);

    // Check if we've reached the end
    if (result.items.length < pageSize) {
      break;
    }

    offset += pageSize;
  }

  console.log(`Total wallets: ${allWallets.length}`);
}
```

### Cursor-Based Pagination

```typescript
import { ProtectClient, Request } from '@taurushq/protect-sdk';

async function cursorPaginationExample(client: ProtectClient): Promise<void> {
  const allRequests: Request[] = [];
  let currentPage: string | undefined;
  let pageNumber = 0;

  while (true) {
    const result = await client.requests.list({
      limit: 50,
      currentPage,
    });

    pageNumber++;
    allRequests.push(...result.requests);
    console.log(`Page ${pageNumber}: ${result.requests.length} items`);

    if (!result.cursor.hasMore) {
      break;
    }

    currentPage = result.cursor.nextCursor;
  }

  console.log(`Total requests: ${allRequests.length}`);
}
```

### Generic Async Iterator Pattern

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function* paginateAll<T>(
  fetchPage: (offset: number) => Promise<{ items: T[]; pagination?: { totalItems: number } }>,
  pageSize: number = 50
): AsyncGenerator<T, void, unknown> {
  let offset = 0;

  while (true) {
    const result = await fetchPage(offset);

    for (const item of result.items) {
      yield item;
    }

    if (result.items.length < pageSize) {
      break;
    }

    offset += pageSize;
  }
}

// Usage
async function useGenericPagination(client: ProtectClient): Promise<void> {
  for await (const wallet of paginateAll(
    (offset) => client.wallets.list({ limit: 50, offset }),
    50
  )) {
    console.log(`Wallet: ${wallet.name}`);
  }
}
```

---

## Error Handling

### Comprehensive Exception Handling

```typescript
import {
  ProtectClient,
  APIError,
  ValidationError,
  AuthenticationError,
  AuthorizationError,
  NotFoundError,
  RateLimitError,
  ServerError,
  IntegrityError,
  WhitelistError,
  ConfigurationError,
} from '@taurushq/protect-sdk';

async function errorHandlingExample(client: ProtectClient): Promise<void> {
  try {
    const wallet = await client.wallets.get(999999);
  } catch (error) {
    if (error instanceof IntegrityError) {
      // Cryptographic verification failed - security issue
      console.error('SECURITY: Integrity check failed:', error.message);
      // This indicates potential tampering - investigate!
      // NEVER retry integrity errors
      return;
    }

    if (error instanceof WhitelistError) {
      // Whitelist-specific verification failed
      console.error('Whitelist verification failed:', error.message);
      return;
    }

    if (error instanceof ValidationError) {
      // 400 Bad Request - input validation failed
      console.error('Validation error:', error.message);
      return;
    }

    if (error instanceof AuthenticationError) {
      // 401 Unauthorized - check credentials
      console.error('Authentication failed:', error.message);
      return;
    }

    if (error instanceof AuthorizationError) {
      // 403 Forbidden - insufficient permissions
      console.error('Permission denied:', error.message);
      return;
    }

    if (error instanceof NotFoundError) {
      // 404 Not Found
      console.error('Resource not found:', error.message);
      return;
    }

    if (error instanceof RateLimitError) {
      // 429 Too Many Requests
      const delay = error.suggestedRetryDelayMs();
      console.error(`Rate limited. Retry after ${delay}ms`);
      return;
    }

    if (error instanceof ServerError) {
      // 5xx Server Error
      console.error('Server error:', error.message, `(HTTP ${error.statusCode})`);
      return;
    }

    if (error instanceof APIError) {
      // Generic API error
      console.error('API error:', error.message);
      console.error('Status:', error.statusCode);
      console.error('Error code:', error.errorCode);
      console.error('Retryable:', error.isRetryable());
      return;
    }

    // Unknown error
    throw error;
  }
}
```

### Retry Logic with Exponential Backoff

```typescript
import { ProtectClient, APIError } from '@taurushq/protect-sdk';

async function withRetry<T>(
  operation: () => Promise<T>,
  maxRetries: number = 3
): Promise<T> {
  let lastError: Error | undefined;

  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      return await operation();
    } catch (error) {
      lastError = error as Error;

      if (error instanceof APIError && error.isRetryable()) {
        const delay = error.suggestedRetryDelayMs() ?? Math.pow(2, attempt) * 1000;
        console.log(`Retry ${attempt + 1}/${maxRetries} after ${delay}ms`);
        await new Promise(resolve => setTimeout(resolve, delay));
        continue;
      }

      // Non-retryable error
      throw error;
    }
  }

  throw lastError;
}

// Usage
async function useWithRetry(client: ProtectClient, walletId: number): Promise<void> {
  const wallet = await withRetry(() => client.wallets.get(walletId), 3);
  console.log(`Got wallet: ${wallet.name}`);
}
```

---

## Async/Await Patterns

### Parallel Operations

```typescript
import { ProtectClient } from '@taurushq/protect-sdk';

async function parallelOperations(client: ProtectClient): Promise<void> {
  // Execute multiple independent operations in parallel
  const [wallets, requests, transactions] = await Promise.all([
    client.wallets.list({ limit: 10 }),
    client.requests.list({ limit: 10 }),
    client.transactions.list({ limit: 10 }),
  ]);

  console.log(`Wallets: ${wallets.items.length}`);
  console.log(`Requests: ${requests.requests.length}`);
  console.log(`Transactions: ${transactions.items.length}`);
}
```

### Batch Processing with Concurrency Control

```typescript
import { ProtectClient, Request } from '@taurushq/protect-sdk';

async function processInBatches<T, R>(
  items: T[],
  batchSize: number,
  processor: (item: T) => Promise<R>
): Promise<R[]> {
  const results: R[] = [];

  for (let i = 0; i < items.length; i += batchSize) {
    const batch = items.slice(i, i + batchSize);
    const batchResults = await Promise.all(batch.map(processor));
    results.push(...batchResults);

    console.log(`Processed ${Math.min(i + batchSize, items.length)}/${items.length}`);
  }

  return results;
}

// Usage
async function processManyRequests(client: ProtectClient): Promise<void> {
  const { requests } = await client.requests.list({ limit: 100 });

  // Process requests in batches of 5 to avoid overwhelming the API
  const verifiedRequests = await processInBatches(
    requests,
    5,
    (request) => client.requests.get(request.id)
  );

  console.log(`Verified ${verifiedRequests.length} requests`);
}
```

### Using try-finally for Cleanup

```typescript
import { ProtectClient, ProtectClientConfig } from '@taurushq/protect-sdk';

async function withClient<T>(
  config: ProtectClientConfig,
  operation: (client: ProtectClient) => Promise<T>
): Promise<T> {
  const client = ProtectClient.create(config);
  try {
    return await operation(client);
  } finally {
    client.close();
  }
}

// Usage
async function main(): Promise<void> {
  const result = await withClient(
    {
      host: process.env.API_HOST!,
      apiKey: process.env.API_KEY!,
      apiSecret: process.env.API_SECRET!,
    },
    async (client) => {
      const wallets = await client.wallets.list();
      return wallets.items.length;
    }
  );

  console.log(`Found ${result} wallets`);
}
```

---

## Best Practices

### 1. Always Close the Client

```typescript
const client = ProtectClient.create(config);
try {
  await doSomething(client);
} finally {
  client.close();
}
```

### 2. Use SuperAdmin Keys in Production

Always configure SuperAdmin keys in production to enable signature verification:

```typescript
const client = ProtectClient.create({
  host: config.host,
  apiKey: config.apiKey,
  apiSecret: config.apiSecret,
  superAdminKeysPem: config.superAdminKeys,
  minValidSignatures: 2,
});
```

### 3. Handle Integrity Errors Specially

Never retry integrity errors - they indicate potential security issues:

```typescript
try {
  const address = await client.addresses.get(id);
} catch (error) {
  if (error instanceof IntegrityError) {
    await logSecurityEvent('integrity_violation', { addressId: id, error });
    await alertSecurityTeam(error);
    throw error; // Do NOT continue with this data
  }
}
```

### 4. Protect Your Private Keys

```typescript
import { createPrivateKey, KeyObject } from 'crypto';

function loadPrivateKeyFromEnv(): KeyObject {
  const keyBase64 = process.env.PROTECT_PRIVATE_KEY;
  if (!keyBase64) {
    throw new Error('PROTECT_PRIVATE_KEY environment variable not set');
  }

  const keyPem = Buffer.from(keyBase64, 'base64').toString('utf-8');
  return createPrivateKey({ key: keyPem, format: 'pem' });
}

// Never log or expose private keys
const privateKey = loadPrivateKeyFromEnv();
```

---

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Authentication](AUTHENTICATION.md) - Security and signing
- [Services Reference](SERVICES.md) - Complete API documentation
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Verification details
