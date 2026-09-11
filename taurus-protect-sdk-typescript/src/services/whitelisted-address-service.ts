/**
 * Whitelisted address service for Taurus-PROTECT SDK.
 *
 * Provides operations for retrieving and verifying whitelisted addresses.
 * All addresses retrieved through this service can be automatically verified
 * for cryptographic integrity using governance rules when verification is enabled.
 */

import type { KeyObject } from "crypto";
import type { AddressWhitelistingApi } from "../internal/openapi";
import { signData } from "../crypto/signing";
import {
  ContainerIntegrityError,
  IntegrityError,
  NotFoundError,
  ValidationError,
} from "../errors";
import { constantTimeCompare } from "../helpers/constant-time";
import { containerHashLabel } from "../helpers/whitelist-hash-helper";
import {
  WhitelistedAddressVerifier,
  type WhitelistedAddressVerifierConfig,
  type RulesContainerDecoder,
  type UserSignaturesDecoder,
} from "../helpers/whitelisted-address-verifier";
import type { Verified } from "../helpers/verified";
import type { DecodedRulesContainer } from "../models/governance-rules";
import { WhitelistedAddressApproval } from "../models/whitelisted-address";
import type {
  WhitelistedAddress,
  SignedWhitelistedAddressEnvelope,
  WhitelistedAddressVerificationResult,
  InternalAddress,
  InternalWallet,
} from "../models/whitelisted-address";
import type { Pagination } from "../models/pagination";
import { BaseService } from "./base";
import { rethrowIfNotRowLevel } from "./row-level-error";

/**
 * Options for listing whitelisted addresses.
 */
/**
 * Filters for the whitelisted-address approval queue.
 *
 * The approval queue is a different ENDPOINT, not a status filter on the general list:
 * it is scoped to what the calling user may act on, which no filter reproduces.
 */
export interface ListWhitelistedAddressesForApprovalOptions {
  /** Maximum number of rows to return. */
  readonly limit?: number;
  /** Number of rows to skip. */
  readonly offset?: number;
  /** Filter by specific whitelisted address IDs. */
  readonly ids?: string[];
  /** Filter by blockchain symbol. */
  readonly blockchain?: string;
  /** Filter by network. */
  readonly network?: string;
  /** Filter by address type. */
  readonly addressType?: string;
  /** Matches customer id, address, blockchain, label, memo and address type. */
  readonly query?: string;
  /**
   * Include rows the calling user has already signed, which is how an approver tells
   * "waiting for me" from "waiting for someone else".
   */
  readonly includeAlreadySignedByUser?: boolean;
}

export interface ListWhitelistedAddressesOptions {
  /** Filter by specific whitelisted address IDs. */
  readonly ids?: string[];
  /** Maximum number of items to return (max 100). */
  limit?: number;
  /** Offset for pagination. */
  offset?: number;
  /** Search query. */
  query?: string;
  /** Filter by blockchain. */
  blockchain?: string;
  /** Filter by network. */
  network?: string;
  /** Filter by address type. */
  addressType?: string;
  /** Include addresses pending approval. */
  includeForApproval?: boolean;
  /** Filter by specific addresses. */
  addresses?: string[];
  /** Filter by tag IDs. */
  tagIds?: string[];
  /** Filter by contract types. */
  contractTypes?: string[];
}

/**
 * Result of listing whitelisted addresses.
 */
/** A row the server returned that failed verification and was left out of `items`. */
export interface ExcludedWhitelistedAddress {
  /** The address ID, or "" when the row carried no usable ID. */
  readonly id: string;
  /** Why the row failed verification. */
  readonly reason: string;
}

export interface ListWhitelistedAddressesResult {
  /** List of verified whitelisted addresses. */
  items: WhitelistedAddress[];
  /** Pagination information. */
  pagination: Pagination | undefined;
  /**
   * Rows that failed verification and are absent from `items`.
   *
   * A short page is not the same as a filtered one, so callers that care about
   * completeness must check this rather than inferring from `items.length`.
   */
  excludedUnverified: ExcludedWhitelistedAddress[];
  /**
   * Pins the named rows for {@link WhitelistedAddressService.approve}, recording the
   * metadata hash each one carried on THIS read.
   *
   * An id that is not in this result is an error rather than a silent omission: it means
   * the caller is trying to approve something this read did not return — either it was
   * excluded as unverifiable, or it was never on the page — and approving fewer rows
   * than asked for would tell the approver they approved more than they did.
   *
   * @param ids - the row ids to pin
   * @throws {@link ValidationError} If `ids` is empty or names a row not in this result
   * @throws {@link IntegrityError} If a named row carries no metadata hash to pin to
   */
  select(...ids: string[]): WhitelistedAddressApproval;
  /**
   * Pins every row this read returned.
   *
   * Note it pins what SURVIVED verification, not what the server sent: rows in
   * `excludedUnverified` are not included, so a caller who wants to know about them must
   * read that field. Approving is all-or-nothing over what is pinned here.
   *
   * @throws {@link IntegrityError} If this read returned no verified rows
   */
  selectAll(): WhitelistedAddressApproval;
}

/**
 * Assembles a list result and attaches the two selection methods over the hashes the
 * read captured.
 *
 * A module-level function rather than a method so both `list` and `listForApproval`
 * mint the same shape — the two readers that would otherwise drift, exactly as
 * `rethrowIfNotRowLevel` was extracted for.
 *
 * @param items - the rows that survived verification
 * @param pagination - the page window, already reduced by the exclusions
 * @param excludedUnverified - the rows withheld, with reasons
 * @param pinnedHashes - row id -> the metadata hash that row carried on this read
 */
function buildAddressListResult(
  items: WhitelistedAddress[],
  pagination: Pagination | undefined,
  excludedUnverified: ExcludedWhitelistedAddress[],
  pinnedHashes: ReadonlyMap<string, string>
): ListWhitelistedAddressesResult {
  const select = (...ids: string[]): WhitelistedAddressApproval => {
    if (ids.length === 0) {
      throw new ValidationError("cannot select an empty set of ids");
    }
    const pinned = new Map<string, string>();
    for (const id of ids) {
      const hash = pinnedHashes.get(id);
      if (hash === undefined) {
        throw new ValidationError(
          `whitelisted address ${id} is not in this verified read: it was either ` +
            `excluded as unverifiable or not on this page`
        );
      }
      if (!hash) {
        throw new IntegrityError(
          `whitelisted address ${id} carries no metadata hash, so there is nothing ` +
            `to pin the approval to`
        );
      }
      pinned.set(id, hash);
    }
    return new WhitelistedAddressApproval(pinned);
  };

  return {
    items,
    pagination,
    excludedUnverified,
    select,
    selectAll: (): WhitelistedAddressApproval => {
      const ids = [...pinnedHashes.keys()];
      if (ids.length === 0) {
        throw new IntegrityError(
          "this read returned no verified addresses to approve"
        );
      }
      return select(...ids);
    },
  };
}

/**
 * Configuration for WhitelistedAddressService with verification.
 */
export interface WhitelistedAddressServiceConfig extends WhitelistedAddressVerifierConfig {
  /** Function to decode rules container from base64. */
  rulesContainerDecoder: RulesContainerDecoder;
  /** Function to decode user signatures from base64. */
  userSignaturesDecoder: UserSignaturesDecoder;
}

/**
 * Service for managing whitelisted addresses.
 *
 * Provides operations for retrieving and verifying whitelisted addresses.
 * When verification is enabled, all get operations will verify the
 * cryptographic signatures according to governance rules.
 *
 * @example
 * ```typescript
 * const service = WhitelistedAddressService.withVerification(api, {
 *   superAdminKeys: [key1, key2],
 *   minValidSignatures: 1,
 *   rulesContainerDecoder,
 *   userSignaturesDecoder,
 * });
 * const result = await service.getWithVerification('123');
 * ```
 */
export class WhitelistedAddressService extends BaseService {
  private readonly verifier: WhitelistedAddressVerifier;
  private readonly rulesContainerDecoder: RulesContainerDecoder;
  private readonly userSignaturesDecoder: UserSignaturesDecoder;

  /**
   * Creates a new WhitelistedAddressService with mandatory verification configuration.
   *
   * Verification is always enforced. All get/list operations use the provided
   * SuperAdmin keys to verify cryptographic integrity.
   *
   * @param api - The AddressWhitelistingApi instance from OpenAPI client
   * @param config - Verification configuration (required)
   */
  constructor(
    private readonly api: AddressWhitelistingApi,
    config: WhitelistedAddressServiceConfig
  ) {
    super();
    this.verifier = new WhitelistedAddressVerifier(config);
    this.rulesContainerDecoder = config.rulesContainerDecoder;
    this.userSignaturesDecoder = config.userSignaturesDecoder;
  }

  /**
   * Creates a WhitelistedAddressService with verification enabled.
   *
   * @param api - The AddressWhitelistingApi instance
   * @param config - Verification configuration
   * @returns WhitelistedAddressService with verification enabled
   */
  static withVerification(
    api: AddressWhitelistingApi,
    config: WhitelistedAddressServiceConfig
  ): WhitelistedAddressService {
    return new WhitelistedAddressService(api, config);
  }

  /**
   * Gets a whitelisted address by ID with mandatory verification.
   *
   * Performs the full 6-step verification of the address, matching
   * Java/Python/Go behavior where verification is always enforced.
   *
   * @param addressId - The address ID to retrieve
   * @returns The verified whitelisted address
   * @throws ValidationError if addressId is invalid
   * @throws NotFoundError if address not found
   * @throws IntegrityError if verification fails
   * @throws APIError if API request fails
   */
  async get(addressId: string): Promise<WhitelistedAddress> {
    if (!addressId) {
      throw new ValidationError("addressId is required");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedAddress({
        id: addressId,
      });

      const dto = response.result;
      if (dto == null) {
        throw new NotFoundError(`Whitelisted address ${addressId} not found`);
      }

      // Map to envelope and verify (6-step verification)
      const envelope = this.mapDtoToEnvelope(dto, addressId);
      const result = this.verifier.verify(
        envelope,
        this.rulesContainerDecoder,
        this.userSignaturesDecoder
      );

      return this.withEnvelopeFields(result.verifiedWhitelistedAddress, dto);
    });
  }

  /**
   * Gets a whitelisted address by ID with full verification.
   *
   * Performs the complete 6-step verification of the address.
   *
   * @param addressId - The address ID to retrieve
   * @returns Verification result with verified address
   * @throws ValidationError if addressId is invalid or verification not configured
   * @throws NotFoundError if address not found
   * @throws IntegrityError if verification fails
   * @throws WhitelistError if governance thresholds not met
   * @throws APIError if API request fails
   */
  async getWithVerification(
    addressId: string
  ): Promise<WhitelistedAddressVerificationResult> {
    if (!addressId) {
      throw new ValidationError("addressId is required");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedAddress({
        id: addressId,
      });

      const dto = response.result;
      if (dto == null) {
        throw new NotFoundError(`Whitelisted address ${addressId} not found`);
      }

      // Map DTO to envelope
      const envelope = this.mapDtoToEnvelope(dto, addressId);

      // Perform verification
      return this.verifier.verify(
        envelope,
        this.rulesContainerDecoder,
        this.userSignaturesDecoder
      );
    });
  }

  /**
   * Gets the signed envelope for a whitelisted address, after verifying it.
   *
   * Returns what verification PRODUCED, not the input: the branded type is the caller's
   * proof that the fields they are about to read were checked, and a read path that
   * skipped `verify()` cannot produce one — it will not compile.
   *
   * This used to run the verification and then return the unverified INPUT envelope,
   * discarding the result. Structurally that is the same defect the `Verified<>` brand
   * was introduced for: the type could not tell an envelope off the wire from one the
   * verifier had cleared, so the two were interchangeable and only the asset side got
   * it right.
   *
   * @param addressId - The address ID to retrieve
   * @returns The verified signed envelope
   * @throws ValidationError if addressId is invalid
   * @throws NotFoundError if address not found
   * @throws IntegrityError if verification fails
   * @throws WhitelistError if governance thresholds are not met
   * @throws APIError if API request fails
   */
  async getEnvelope(
    addressId: string
  ): Promise<Verified<SignedWhitelistedAddressEnvelope>> {
    if (!addressId) {
      throw new ValidationError("addressId is required");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedAddress({
        id: addressId,
      });

      const dto = response.result;
      if (dto == null) {
        throw new NotFoundError(`Whitelisted address ${addressId} not found`);
      }

      return this.verifier.verify(
        this.mapDtoToEnvelope(dto, addressId),
        this.rulesContainerDecoder,
        this.userSignaturesDecoder
      ).verifiedEnvelope;
    });
  }

  /**
   * Lists whitelisted addresses with mandatory verification.
   *
   * Each address undergoes the full 6-step verification, matching
   * Java/Python/Go behavior where verification is always enforced.
   * Fails on the first verification error (strict mode).
   *
   * @param options - Optional filtering and pagination options
   * @returns Object with verified addresses array and pagination info
   * @throws IntegrityError if verification fails for any address
   * @throws WhitelistError if governance thresholds not met
   * @throws APIError if API request fails
   */
  async list(
    options?: ListWhitelistedAddressesOptions
  ): Promise<ListWhitelistedAddressesResult> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError("limit must be positive");
    }
    if (offset < 0) {
      throw new ValidationError("offset cannot be negative");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedAddresses({
        limit: String(limit),
        offset: String(offset),
        query: options?.query,
        blockchain: options?.blockchain,
        network: options?.network,
        addressType: options?.addressType,
        includeForApproval: options?.includeForApproval,
        addresses: options?.addresses,
        tagIDs: options?.tagIds,
        contractTypes: options?.contractTypes,
        rulesContainerNormalized: true,
        ids: options?.ids,
      });

      // Build rules container cache from normalized response
      const rulesContainerCache = this.buildRulesContainerCache(response);

      //	rows ──▶ verify ──┬─ ok ───▶ returned
      //	                  └─ fail ─▶ excluded + reported, list survives
      //
      // One unverifiable row used to fail the whole call, which took down the
      // whitelist for every consumer — and listing is how an operator would find
      // the bad row, so the failure hid its own cause. Excluding stays fail-closed:
      // an omitted destination cannot be selected, so a transfer to it is refused.
      const rows = response.result ?? [];
      const items: WhitelistedAddress[] = [];
      const excludedUnverified: ExcludedWhitelistedAddress[] = [];
      // The content pin, captured while the DTO is still in hand: this SDK's
      // WhitelistedAddress carries no metadata, so the hash an approver would sign is
      // not recoverable from `items` afterwards.
      const pinnedHashes = new Map<string, string>();

      for (const dto of rows) {
        const dtoId = dto.id ?? "";
        try {
          const envelope = this.mapDtoToEnvelope(dto, dtoId);

          // Look up cached rules container by hash
          const cached = envelope.rulesContainerHash
            ? rulesContainerCache.get(envelope.rulesContainerHash)
            : undefined;

          const result = this.verifier.verify(
            envelope,
            this.rulesContainerDecoder,
            this.userSignaturesDecoder,
            cached
          );
          items.push(this.withEnvelopeFields(result.verifiedWhitelistedAddress, dto));
          pinnedHashes.set(dtoId, dto.metadata?.hash ?? "");
        } catch (error: unknown) {
          rethrowIfNotRowLevel(error);
          excludedUnverified.push({
            id: dtoId,
            reason: (error as Error).message,
          });
        }
      }

      // Rows came back but none survived: a systemic failure, not an empty
      // whitelist, and returning [] would be indistinguishable from one.
      if (rows.length > 0 && items.length === 0) {
        throw new IntegrityError(
          `all ${rows.length} whitelisted address(es) failed verification; ` +
            `first failure: ${excludedUnverified[0]?.reason ?? "unknown"}`
        );
      }

      // Extract pagination, minus the rows the caller never receives.
      //
      // The server counts rows it returned; excluded rows are not among the items.
      // Reporting the server's total lets a filtered page pass for a complete one.
      const reportedTotal = response.totalItems
        ? parseInt(response.totalItems, 10)
        : undefined;
      // isNaN guard as in listForApproval: without it a non-numeric totalItems yields
      // Math.max(0, NaN - n) === NaN, which is !== undefined, so NaN reaches the caller.
      const totalItems =
        reportedTotal === undefined || isNaN(reportedTotal)
          ? undefined
          : Math.max(0, reportedTotal - excludedUnverified.length);
      const pagination: Pagination | undefined = totalItems !== undefined
        ? { totalItems, offset, limit }
        : undefined;

      return buildAddressListResult(
        items,
        pagination,
        excludedUnverified,
        pinnedHashes
      );
    });
  }

  /**
   * Adds the non-security fields the SIGNED payload does not carry.
   *
   * `createdAt` comes from the trails array (the "created" action) and `attributes`
   * from the envelope's attribute list — both DTO-sourced, and both safe to take
   * from there because neither participates in the governance decision. Go, Java and
   * Python all populate them; this SDK left them `undefined` and `{}`.
   *
   * Security-critical fields are NOT touched here: they stay exactly as the verified
   * payload produced them.
   */
  /**
   * Lists whitelisted addresses awaiting approval, verified exactly as {@link list} is.
   *
   * Both this endpoint and {@link approve} are generated in all four SDK clients and
   * were wrapped by none of them, so the rows an approver reads before whitelisting a
   * destination were not reachable through the SDK at all — verified or not.
   */
  async listForApproval(
    options?: ListWhitelistedAddressesForApprovalOptions
  ): Promise<ListWhitelistedAddressesResult> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError("limit must be positive");
    }
    if (offset < 0) {
      throw new ValidationError("offset cannot be negative");
    }

    return this.execute(async () => {
      const response =
        await this.api.whitelistServiceGetWhitelistedAddressesForApproval({
          limit: String(limit),
          offset: String(offset),
          ids: options?.ids,
          blockchain: options?.blockchain,
          network: options?.network,
          addressType: options?.addressType,
          query: options?.query,
          includeAlreadySignedByUser: options?.includeAlreadySignedByUser,
        });

      // This endpoint has no normalized-container mode, so the per-row containers are
      // used and the cache starts empty.
      const rows = response.result ?? [];
      const items: WhitelistedAddress[] = [];
      const excludedUnverified: ExcludedWhitelistedAddress[] = [];
      const pinnedHashes = new Map<string, string>();

      for (const dto of rows) {
        const rowId = dto.id ?? "";
        try {
          const result = this.verifier.verify(
            this.mapDtoToEnvelope(dto, rowId),
            this.rulesContainerDecoder,
            this.userSignaturesDecoder
          );
          items.push(
            this.withEnvelopeFields(result.verifiedWhitelistedAddress, dto)
          );
          pinnedHashes.set(rowId, dto.metadata?.hash ?? "");
        } catch (error: unknown) {
          rethrowIfNotRowLevel(error);
          excludedUnverified.push({
            id: rowId,
            reason: (error as Error).message,
          });
        }
      }

      if (rows.length > 0 && items.length === 0) {
        throw new IntegrityError(
          `all ${rows.length} whitelisted address(es) failed verification; ` +
            `first failure: ${excludedUnverified[0]?.reason ?? "unknown"}`
        );
      }

      const reportedTotal = response.totalItems
        ? parseInt(response.totalItems, 10)
        : undefined;
      const totalItems =
        reportedTotal === undefined || isNaN(reportedTotal)
          ? undefined
          : Math.max(0, reportedTotal - excludedUnverified.length);

      return buildAddressListResult(
        items,
        // TypeScript's Pagination is {totalItems, offset, limit} — it has no `hasMore`,
        // unlike Go's. Keep the shape per SDK rather than inventing a field here.
        totalItems !== undefined ? { limit, offset, totalItems } : undefined,
        excludedUnverified,
        pinnedHashes
      );
    });
  }

  /**
   * Signs and submits an approval for the whitelisted addresses an approver REVIEWED,
   * all-or-nothing.
   *
   * The selection comes from a preceding verified read — `result.select(...)` or
   * `result.selectAll()` on {@link list} / {@link listForApproval} — so it carries the
   * metadata hash each row had at review time. The batch is then re-read through the
   * verifying path and every re-read hash must equal its pin, or the call aborts and
   * NOTHING is signed.
   *
   *   selection -> sort numerically -> ONE filtered verified page -> completeness check
   *                                                                        |
   *                                                        pin == re-read hash ?
   *                                                                        |
   *                                            sign(JSON(hashes)) -> POST once
   *
   * The pin is what defeats the substitution the re-read alone cannot see: the API takes
   * only ids, so without it a response-controlling server can answer the review call
   * with one row and the re-read with another whose existing signatures already satisfy
   * the container it presents — harvesting a genuine approver signature over content the
   * approver never saw. An abort here is CORRECT, not a false positive; the remedy is
   * re-read, re-review, re-approve. Note `metadata.hash` is recomputed by validatord on
   * every read from live linked-address/wallet rows, so renaming a linked wallet moves
   * it legitimately — which is precisely the change an approver most needs to be told
   * about.
   *
   * Taking a `string[]` of ids was the previous shape and is deliberately gone: an
   * optional pin that defaults to "no pin" restores the unpinned behaviour silently,
   * which is the mistake `approveRulesProposal` already records.
   *
   * @param selection - rows pinned by a verified read
   * @param privateKey - the approver's P-256 private key
   * @param comment - the approval comment
   * @throws {@link ValidationError} If the selection is empty or an argument is missing
   * @throws {@link IntegrityError} If a row is missing, unverifiable, or has moved
   *   since it was reviewed
   */
  async approve(
    selection: WhitelistedAddressApproval,
    privateKey: KeyObject,
    comment: string
  ): Promise<void> {
    if (!selection || selection.isEmpty()) {
      throw new ValidationError(
        "selection cannot be empty: pin the rows with result.select(...ids) or " +
          "result.selectAll() from a verified read"
      );
    }
    if (!privateKey) {
      throw new ValidationError("privateKey is required");
    }
    if (!comment || comment.trim() === "") {
      throw new ValidationError("comment is required");
    }
    const ids = selection.ids();
    for (const id of ids) {
      if (!/^\d+$/.test(id)) {
        throw new ValidationError(
          `whitelisted address ID ${id} must be a positive integer`
        );
      }
    }

    // The endpoint requires ascending order, and sorting here also makes the signed
    // order independent of the order the caller passed.
    const sortedIds = [...ids].sort((a, b) => Number(a) - Number(b));

    const verifiedHashes = await this.hashesToSignByID(sortedIds);

    const hashes: string[] = [];
    for (const id of sortedIds) {
      const hash = verifiedHashes.get(id);
      if (hash === undefined) {
        // A page that silently omits a row must not become an approval of fewer rows
        // than the caller asked for.
        throw new IntegrityError(
          `refusing to sign: address ${id} was not returned by the verified read`
        );
      }
      if (!hash) {
        throw new IntegrityError(
          `refusing to sign: address ${id} has no metadata hash`
        );
      }

      // The pin. Constant-time because this compares hash material, matching how
      // approveRulesProposal handles its container pin.
      const pinned = selection.pinnedHash(id);
      if (pinned === undefined) {
        throw new IntegrityError(
          `refusing to sign: address ${id} is not in the reviewed selection`
        );
      }
      if (!constantTimeCompare(pinned, hash)) {
        throw new IntegrityError(
          `refusing to sign: whitelisted address ${id} changed since it was reviewed: ` +
            `reviewed hash ${pinned}, re-read hash ${hash}. Re-read, re-review and ` +
            `re-approve`
        );
      }

      hashes.push(hash);
    }

    const signature = signData(
      privateKey,
      Buffer.from(JSON.stringify(hashes), "utf-8")
    );

    return this.execute(async () => {
      await this.api.whitelistServiceApproveWhitelistedAddress({
        body: { ids: sortedIds, signature, comment },
      });
    });
  }

  /**
   * Re-reads a batch of addresses through the verifying path, filtered by id, and
   * returns the hash to sign for each, keyed by id.
   *
   * The row must be VERIFIED before its hash is signed — that is what the re-read is
   * for. But the value signed is the row's CURRENT `metadata.hash`, not the
   * `verifiedHash` the verifier returned.
   *
   * Those differ only for a legacy row, where step 4 matched a hash computed from a
   * stripped payload. Signing the legacy variant there would be rejected by the server:
   * `ApproveWhitelistedAddressRequest` carries no `hashes` field, so validatord rebuilds
   * the array itself from `helper.ToWLAMetadata`, which is hard-wired to the current
   * schema version, and verifies the submitted signature against those bytes
   * (`pkg/whitelist/service/whitelist-service.go` and `wl_address_signatures.go`). The
   * three legacy variants are accepted only by the per-signature tolerance check in
   * `VerifyWLAddressesIntegrity`, which keeps OLD signatures valid — it does not make a
   * legacy hash acceptable for a NEW one. A legacy row therefore ends up carrying a mix,
   * and both kinds count toward the group threshold.
   *
   * This SDK signed `verifiedHash` and was the only one of the four that did; Go, Java
   * and Python all sign `metadata.hash`. Do not "align" the other three to this one.
   *
   * Ordering matters for the same reason: the server sorts by whitelisted-address id
   * ascending before building the array, which is what the caller's sort reproduces.
   */
  private async hashesToSignByID(
    ids: string[]
  ): Promise<Map<string, string>> {
    const response = await this.api.whitelistServiceGetWhitelistedAddresses({
      limit: String(ids.length),
      offset: "0",
      rulesContainerNormalized: true,
      includeForApproval: true,
      ids,
    });

    const cache = this.buildRulesContainerCache(response);
    const byID = new Map<string, string>();
    for (const dto of response.result ?? []) {
      try {
        // Verification must clear the row before its hash is signed. It also proves
        // sha256(payloadAsString) === metadata.hash (step 1), which is what makes the
        // DTO hash safe to use here rather than a value the response merely asserts.
        this.verifier.verify(
          this.mapDtoToEnvelope(dto, dto.id ?? ""),
          this.rulesContainerDecoder,
          this.userSignaturesDecoder,
          dto.rulesContainerHash
            ? cache.get(dto.rulesContainerHash)
            : undefined
        );
        byID.set(dto.id ?? "", dto.metadata?.hash ?? "");
      } catch (error: unknown) {
        // A container this SDK cannot interpret is not a property of this row, and a
        // defect is not a verification failure — both must surface. Only a row-level
        // integrity failure is left out of the map, where the caller's completeness
        // check turns it into a refusal naming the id, which is more use than a bare
        // verification error.
        rethrowIfNotRowLevel(error);
        continue;
      }
    }
    return byID;
  }

  private withEnvelopeFields(
    address: WhitelistedAddress,
    dto: {
      trails?: Array<{ action?: string; date?: Date }>;
      attributes?: Array<{ key?: string; value?: string }>;
    }
  ): WhitelistedAddress {
    let createdAt: Date | undefined;
    for (const trail of dto.trails ?? []) {
      if (trail.action === "created") {
        createdAt = trail.date;
        break;
      }
    }

    const attributes: Record<string, unknown> = {};
    for (const attr of dto.attributes ?? []) {
      if (attr.key != null) {
        attributes[attr.key] = attr.value;
      }
    }

    return { ...address, createdAt, attributes };
  }

  /**
   * Builds a cache of verified rules containers from the normalized response.
   * When rulesContainerNormalized=true, the API returns deduplicated rules containers.
   * Each is verified once and cached by hash.
   */
  private buildRulesContainerCache(
    response: { rulesContainers?: Array<{ hash?: string; rulesContainer?: string; rulesSignatures?: string }> }
  ): Map<string, DecodedRulesContainer> {
    const cache = new Map<string, DecodedRulesContainer>();
    if (!response.rulesContainers || response.rulesContainers.length === 0) {
      return cache;
    }

    // Deduplicate by base64 container string to avoid re-verifying identical containers
    const verifiedContainers = new Map<string, DecodedRulesContainer>();

    for (const hashContainer of response.rulesContainers) {
      const containerHash = hashContainer.hash;
      const containerBase64 = hashContainer.rulesContainer;
      const signaturesBase64 = hashContainer.rulesSignatures;

      if (!containerHash || !containerBase64) {
        continue;
      }

      // The hash is the LABEL a row uses to pick its container, and it arrives in the
      // same response as the container. Recompute it: otherwise a server can file
      // container A under container B's label and steer any row to any other
      // validly-signed container — an older ruleset with a weaker group threshold, say.
      // Both pass the SuperAdmin check, so signature verification alone misses this.
      const computed = containerHashLabel(containerBase64);
      if (computed !== containerHash) {
        throw new ContainerIntegrityError(
          `rules container hash mismatch: response labelled it ${containerHash} ` +
            `but its bytes hash to ${computed}`
        );
      }

      // Check if already verified (dedup by content)
      let decoded = verifiedContainers.get(containerBase64);
      if (!decoded) {
        decoded = this.verifier.verifyAndDecodeRulesContainer(
          containerBase64,
          signaturesBase64 ?? "",
          this.rulesContainerDecoder,
          this.userSignaturesDecoder
        );
        verifiedContainers.set(containerBase64, decoded);
      }

      cache.set(containerHash, decoded);
    }

    return cache;
  }

  /**
   * Maps a DTO to a SignedWhitelistedAddressEnvelope.
   */
  private mapDtoToEnvelope(
    dto: {
      id?: string;
      metadata?: { hash?: string; payloadAsString?: string; payload?: object };
      rulesContainer?: string;
      rulesContainerHash?: string;
      rulesSignatures?: string;
      signedAddress?: {
        payload?: string;
        signatures?: Array<{
          signature?: { userId?: string; signature?: string; comment?: string };
          hashes?: string[];
        }>;
      };
      blockchain?: string;
      network?: string;
    },
    addressId: string
  ): SignedWhitelistedAddressEnvelope {
    // SECURITY: Extract linked addresses/wallets from verified payloadAsString
    //
    // The metadata contains two representations:
    //   - payload: Raw object from API response (UNVERIFIED - could be tampered)
    //   - payloadAsString: JSON string that is cryptographically hashed (VERIFIED)
    //
    // The hash verification chain is:
    //   1. Server: metadata.hash = SHA256(payloadAsString)
    //   2. Hash is signed by governance rules (SuperAdmin keys)
    //   3. Client verifies: computed_hash(payloadAsString) === metadata.hash
    //
    // ATTACK VECTOR (if using raw payload):
    // An attacker intercepting the API response could:
    //   1. Modify payload.linkedInternalAddresses to include malicious addresses
    //   2. Leave payloadAsString unchanged (hash still verifies)
    //   3. Client trusts the modified linked addresses → SECURITY BYPASS
    //
    // By parsing from payloadAsString, we ensure the integrity chain is preserved.
    let linkedInternalAddresses: InternalAddress[] = [];
    let linkedWallets: InternalWallet[] = [];

    if (dto.metadata?.payloadAsString) {
      try {
        const verifiedPayload = JSON.parse(dto.metadata.payloadAsString) as Record<string, unknown>;
        if (Array.isArray(verifiedPayload.linkedInternalAddresses)) {
          linkedInternalAddresses = verifiedPayload.linkedInternalAddresses.map(
            (addr: { id?: number; label?: string }) => ({
              id: addr.id ?? 0,
              label: addr.label,
            })
          );
        }
        if (Array.isArray(verifiedPayload.linkedWallets)) {
          linkedWallets = verifiedPayload.linkedWallets.map(
            (wallet: { id?: number; path?: string; name?: string }) => ({
              id: wallet.id ?? 0,
              path: wallet.path,
              label: wallet.name,
            })
          );
        }
      } catch (error: unknown) {
        // JSON parsing failure is expected for malformed payloads — leave arrays empty.
        // The main address data comes from parseWhitelistedAddressFromJson in the verifier.
        if (!(error instanceof SyntaxError)) {
          throw error;
        }
      }
    }

    return {
      id: addressId,
      metadata: {
        hash: dto.metadata?.hash ?? "",
        payloadAsString: dto.metadata?.payloadAsString ?? "",
      },
      rulesContainerBase64: dto.rulesContainer ?? "",
      rulesSignaturesBase64: dto.rulesSignatures ?? "",
      rulesContainerHash: dto.rulesContainerHash,
      signedAddress: {
        payload: dto.signedAddress?.payload,
        signatures: (dto.signedAddress?.signatures ?? []).map((sig) => ({
          userSignature: sig.signature
            ? {
                userId: sig.signature.userId,
                signature: sig.signature.signature,
                comment: sig.signature.comment,
              }
            : undefined,
          hashes: sig.hashes ?? [],
        })),
      },
      blockchain: dto.blockchain ?? "",
      network: dto.network ?? "",
      linkedInternalAddresses,
      linkedWallets,
    };
  }
}
