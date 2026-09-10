/**
 * Whitelisted asset (contract) service for Taurus-PROTECT SDK.
 *
 * Every read path verifies. There is no unverified accessor: the only way to
 * obtain an asset is through the 5-step verification below, so a caller cannot
 * hold asset data that was never checked.
 *
 *	rows ──▶ hash ──▶ SuperAdmin sigs ──▶ decode ──▶ coverage ──▶ thresholds ──┬─ ok ───▶ returned
 *	                                                                           └─ fail ─▶ excluded
 *
 * On list, one unverifiable row is excluded rather than failing the whole call —
 * listing is how an operator finds the bad row, so aborting would hide its own
 * cause. Excluding stays fail-closed: an omitted asset cannot be selected. If
 * rows were returned but none survived, that is a systemic failure and errors.
 */

import type { KeyObject } from "crypto";

import type {
  ContractWhitelistingApi,
  TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply,
} from "../internal/openapi";
import { signData } from "../crypto";
import { IntegrityError, NotFoundError, ValidationError } from "../errors";
import {
  WhitelistedAssetVerifier,
  type WhitelistedAssetVerifierConfig,
  type RulesContainerDecoder,
  type UserSignaturesDecoder,
} from "../helpers/whitelisted-asset-verifier";
import type {
  WhitelistedAsset,
  SignedWhitelistedAssetEnvelope,
  WhitelistedAssetVerificationResult,
} from "../models/whitelisted-asset";
import type { Pagination } from "../models/pagination";
import { BaseService } from "./base";
import { rethrowIfNotRowLevel } from "./row-level-error";
import type { Verified } from "../helpers/verified";

/**
 * Options for listing whitelisted assets.
 */
export interface ListWhitelistedAssetsOptions {
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
  /** Include assets pending approval. */
  includeForApproval?: boolean;
  /** Filter by contract kind types (e.g., 'nft', 'token'). */
  kindTypes?: string[];
  /** Filter by specific whitelisted asset IDs. */
  ids?: string[];
}

/**
 * Options for listing whitelisted assets awaiting approval. The endpoint accepts only
 * these three.
 */
export interface ListWhitelistedAssetsForApprovalOptions {
  /** Maximum number of items to return (max 100). */
  limit?: number;
  /** Offset for pagination. */
  offset?: number;
  /** Filter by specific whitelisted asset IDs. */
  ids?: string[];
}

/**
 * A row the server returned that failed verification and was left out of `items`.
 */
export interface ExcludedWhitelistedAsset {
  /** The asset ID, or 0 when the row carried no usable ID. */
  readonly id: number;
  /** Why the row failed verification. */
  readonly reason: string;
}

/**
 * Result of listing whitelisted assets.
 */
export interface ListWhitelistedAssetsResult {
  /** List of verified whitelisted assets. */
  items: WhitelistedAsset[];
  /** Pagination information. */
  pagination: Pagination | undefined;
  /**
   * Rows that failed verification and are absent from `items`.
   *
   * A short page is not the same as a filtered one, so callers that care about
   * completeness must check this rather than inferring from `items.length`.
   */
  excludedUnverified: ExcludedWhitelistedAsset[];
}

/**
 * Configuration for WhitelistedAssetService with verification.
 */
export interface WhitelistedAssetServiceConfig extends WhitelistedAssetVerifierConfig {
  /** Function to decode rules container from base64. */
  rulesContainerDecoder: RulesContainerDecoder;
  /** Function to decode user signatures from base64. */
  userSignaturesDecoder: UserSignaturesDecoder;
}

/**
 * Service for managing whitelisted assets (contracts).
 *
 * Provides operations for retrieving and verifying whitelisted assets.
 * When verification is enabled, all get operations will verify the
 * cryptographic signatures according to governance rules.
 *
 * @example
 * ```typescript
 * const service = WhitelistedAssetService.withVerification(api, {
 *   superAdminKeys: [key1, key2],
 *   minValidSignatures: 1,
 *   rulesContainerDecoder,
 *   userSignaturesDecoder,
 * });
 * const result = await service.getWithVerification(123);
 * ```
 */
export class WhitelistedAssetService extends BaseService {
  private readonly verifier: WhitelistedAssetVerifier;
  private readonly rulesContainerDecoder: RulesContainerDecoder;
  private readonly userSignaturesDecoder: UserSignaturesDecoder;

  /**
   * Creates a new WhitelistedAssetService with mandatory verification configuration.
   *
   * Verification is always enforced. All get/list operations use the provided
   * SuperAdmin keys to verify cryptographic integrity.
   *
   * @param api - The ContractWhitelistingApi instance from OpenAPI client
   * @param config - Verification configuration (required)
   */
  constructor(
    private readonly api: ContractWhitelistingApi,
    config: WhitelistedAssetServiceConfig
  ) {
    super();
    this.verifier = new WhitelistedAssetVerifier(config);
    this.rulesContainerDecoder = config.rulesContainerDecoder;
    this.userSignaturesDecoder = config.userSignaturesDecoder;
  }

  /**
   * Creates a WhitelistedAssetService with verification enabled.
   *
   * @param api - The ContractWhitelistingApi instance
   * @param config - Verification configuration
   * @returns WhitelistedAssetService with verification enabled
   */
  static withVerification(
    api: ContractWhitelistingApi,
    config: WhitelistedAssetServiceConfig
  ): WhitelistedAssetService {
    return new WhitelistedAssetService(api, config);
  }

  /**
   * Gets a whitelisted asset by ID, running the full 5-step verification.
   *
   * Every field on the returned asset comes from the cryptographically verified
   * payload. Use {@link getWithVerification} when the verified hash is also needed.
   *
   * @param assetId - The asset ID to retrieve
   * @returns The verified whitelisted asset
   * @throws ValidationError if assetId is invalid
   * @throws NotFoundError if asset not found
   * @throws IntegrityError if verification fails
   * @throws WhitelistError if governance thresholds are not met
   * @throws APIError if API request fails
   */
  async get(assetId: number): Promise<WhitelistedAsset> {
    const result = await this.getWithVerification(assetId);
    return result.verifiedAsset;
  }

  /**
   * Gets a whitelisted asset by ID with full verification.
   *
   * Performs the complete 5-step verification of the asset.
   *
   * @param assetId - The asset ID to retrieve
   * @returns Verification result with verified asset
   * @throws ValidationError if assetId is invalid or verification not configured
   * @throws NotFoundError if asset not found
   * @throws IntegrityError if verification fails
   * @throws WhitelistError if governance thresholds not met
   * @throws APIError if API request fails
   */
  async getWithVerification(
    assetId: number
  ): Promise<WhitelistedAssetVerificationResult> {
    if (assetId <= 0) {
      throw new ValidationError("assetId must be positive");
    }

    return this.execute(async () => {
      const envelope = await this.fetchEnvelope(assetId);
      return this.verifyEnvelope(envelope);
    });
  }

  /**
   * Gets the signed envelope for a whitelisted asset, after verifying it.
   *
   * The envelope is returned only once all five steps pass, so a caller reading
   * its raw fields is reading data that was checked.
   *
   * @param assetId - The asset ID to retrieve
   * @returns The verified signed envelope
   * @throws ValidationError if assetId is invalid
   * @throws NotFoundError if asset not found
   * @throws IntegrityError if verification fails
   * @throws WhitelistError if governance thresholds are not met
   * @throws APIError if API request fails
   */
  async getEnvelope(
    assetId: number
  ): Promise<Verified<SignedWhitelistedAssetEnvelope>> {
    if (assetId <= 0) {
      throw new ValidationError("assetId must be positive");
    }

    return this.execute(async () => {
      const envelope = await this.fetchEnvelope(assetId);
      // Returns what verification produced, not the input: the branded type is the
      // caller's proof that the fields they are about to read were checked.
      return this.verifyEnvelope(envelope).verifiedEnvelope;
    });
  }

  /**
   * Lists whitelisted assets.
   *
   * @param options - Optional filtering and pagination options
   * @returns Object with assets array and pagination info
   * @throws APIError if API request fails
   */
  async list(
    options?: ListWhitelistedAssetsOptions
  ): Promise<ListWhitelistedAssetsResult> {
    const limit = options?.limit ?? 50;
    const offset = options?.offset ?? 0;

    if (limit <= 0) {
      throw new ValidationError("limit must be positive");
    }
    if (offset < 0) {
      throw new ValidationError("offset cannot be negative");
    }

    return this.execute(async () => {
      const response = await this.api.whitelistServiceGetWhitelistedContracts({
        limit: String(limit),
        offset: String(offset),
        query: options?.query,
        blockchain: options?.blockchain,
        network: options?.network,
        includeForApproval: options?.includeForApproval,
        kindTypes: options?.kindTypes,
        whitelistedContractAddressIds: options?.ids,
      });

      return this.verifyPage(response, limit, offset);
    });
  }

  /**
   * Lists whitelisted assets awaiting approval, verified exactly as {@link list}.
   *
   * Without this, the only reader of the for-approval endpoint was the unverified
   * contract-whitelisting service, so the rows an approver inspects before whitelisting
   * a contract address were never checked against governance.
   *
   * @param options - Optional filtering and pagination options
   * @throws {@link IntegrityError} If rows came back and none verified
   * @throws {@link APIError} If the API request fails
   */
  async listForApproval(
    options?: ListWhitelistedAssetsForApprovalOptions
  ): Promise<ListWhitelistedAssetsResult> {
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
        await this.api.whitelistServiceGetWhitelistedContractsForApproval({
          limit: String(limit),
          offset: String(offset),
          ids: options?.ids,
        });

      return this.verifyPage(response, limit, offset);
    });
  }

  /**
   * Signs and submits an approval for the given whitelisted assets, all-or-nothing.
   *
   * Each asset is re-read through the verified path and the hashes THOSE rows carry are
   * what gets signed, so the approver's signature covers metadata this SDK checked rather
   * than whatever a caller was handed. `ContractWhitelistingService.approve` takes an
   * opaque signature over hashes nothing verified.
   *
   * Any asset that is missing or fails verification aborts the whole call and nothing is
   * signed: one signature covers every hash in the batch, so a partial approval would mean
   * the caller believes they approved more than they did.
   *
   * @param ids - The whitelisted asset IDs to approve
   * @param privateKey - The approver's P-256 private key
   * @param comment - The approval comment
   * @throws {@link ValidationError} If parameters are invalid
   * @throws {@link IntegrityError} If any asset is missing or has no metadata hash
   * @throws {@link APIError} If the API request fails
   */
  async approve(
    ids: number[],
    privateKey: KeyObject,
    comment: string
  ): Promise<void> {
    if (!ids || ids.length === 0) {
      throw new ValidationError("ids cannot be empty");
    }
    if (!privateKey) {
      throw new ValidationError("privateKey is required");
    }
    if (!comment || comment.trim() === "") {
      throw new ValidationError("comment is required");
    }
    for (const id of ids) {
      if (!Number.isInteger(id) || id <= 0) {
        throw new ValidationError(`whitelisted asset ID ${id} must be a positive integer`);
      }
    }

    // Sorted numerically, as the request-approval path does, so the signed order is
    // independent of the order the caller passed.
    const sortedIds = [...ids].sort((a, b) => a - b);

    // ONE id-filtered page through the verifying path, not one GET per id. The list
    // endpoint verifies every row and fetches the rules container once per call, so a
    // 50-id approval costs one round trip instead of fifty. includeForApproval is
    // required: the rows being approved are pending, so the default list omits them.
    const verifiedHashes = await this.hashesToSignByID(sortedIds);

    const hashes: string[] = [];
    for (const id of sortedIds) {
      const verifiedHash = verifiedHashes.get(id);
      if (verifiedHash === undefined) {
        // A page that silently omits a row must not become an approval of fewer rows
        // than the caller asked for.
        throw new IntegrityError(
          `refusing to sign: asset ${id} was not returned by the verified read`
        );
      }
      if (!verifiedHash) {
        throw new IntegrityError(
          `refusing to sign: asset ${id} has no metadata hash`
        );
      }
      hashes.push(verifiedHash);
    }

    const signature = signData(
      privateKey,
      Buffer.from(JSON.stringify(hashes), "utf-8")
    );

    return this.execute(async () => {
      await this.api.whitelistServiceApproveWhitelistedContract({
        body: {
          ids: sortedIds.map(String),
          signature,
          comment,
        },
      });
    });
  }

  /**
   * Verifies every row of a page leniently: an unverifiable row is excluded and named
   * rather than failing the call, but a page where NO row survived throws — a filtered
   * page must never read as an empty whitelist.
   */
  /**
   * Re-reads a batch of assets through the verifying path, filtered by id, and returns
   * the hash to sign for each, keyed by id.
   *
   *   ids -> ONE filtered page -> verify every row -> map id -> metadata.hash
   *
   * The row must be VERIFIED before its hash is signed — that is what the re-read is
   * for. But the value signed is the row's CURRENT `metadata.hash`, not the
   * `verifiedHash` the verifier returned, for exactly the reason spelled out on
   * WhitelistedAddressService.hashesToSignByID: `ApproveWhitelistedContractAddressReq`
   * carries no `hashes` field, so validatord rebuilds the array itself — sorted by id
   * ascending, from `helper.ToWLCAMetadata`, which is the current schema version — and
   * verifies the submitted signature against those bytes
   * (`pkg/whitelist/service/whitelist-service.go`, ApproveWhitelistedContractAddress).
   * Signing an asset legacy variant (isNFT / kindType stripped) would be rejected.
   */
  private async hashesToSignByID(ids: number[]): Promise<Map<number, string>> {
    const response = await this.api.whitelistServiceGetWhitelistedContracts({
      limit: String(ids.length),
      offset: "0",
      includeForApproval: true,
      whitelistedContractAddressIds: ids.map(String),
    });

    const byID = new Map<number, string>();
    for (const dto of response.result ?? []) {
      const rowId = dto.id ? parseInt(dto.id, 10) : 0;
      try {
        // Verification must clear the row before its hash is signed. It also proves
        // sha256(payloadAsString) === metadata.hash (step 1), which is what makes the
        // DTO hash safe to use here rather than a value the response merely asserts.
        this.verifyEnvelope(this.mapDtoToEnvelope(dto, rowId));
        byID.set(rowId, dto.metadata?.hash ?? "");
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

  private verifyPage(
    response: TgvalidatordGetSignedWhitelistedContractAddressEnvelopesReply,
    limit: number,
    offset: number
  ): ListWhitelistedAssetsResult {
    const rows = response.result ?? [];
    const items: WhitelistedAsset[] = [];
    const excludedUnverified: ExcludedWhitelistedAsset[] = [];

    for (const dto of rows) {
      const rowId = dto.id ? parseInt(dto.id, 10) : 0;
      try {
        const result = this.verifyEnvelope(this.mapDtoToEnvelope(dto, rowId));
        items.push(result.verifiedAsset);
      } catch (error: unknown) {
        rethrowIfNotRowLevel(error);
        excludedUnverified.push({
          id: rowId,
          reason: error.message,
        });
      }
    }

    if (rows.length > 0 && items.length === 0) {
      throw new IntegrityError(
        `all ${rows.length} whitelisted asset(s) failed verification; ` +
          `first failure: ${excludedUnverified[0]?.reason ?? "unknown"}`
      );
    }

    // Reduced by the rows the caller never receives, as on the address side.
    //
    // The server counts rows it returned; excluded rows are not among `items`, so
    // reporting the server's total lets a filtered page pass for a complete one and
    // makes pagination promise rows that can never be read. The isNaN guard matters
    // because Math.max(0, NaN - n) is NaN, which is !== undefined and would reach the
    // caller as a NaN totalItems.
    const reportedTotal = response.totalItems
      ? parseInt(response.totalItems, 10)
      : undefined;
    const totalItems =
      reportedTotal === undefined || isNaN(reportedTotal)
        ? undefined
        : Math.max(0, reportedTotal - excludedUnverified.length);
    const pagination: Pagination | undefined = totalItems !== undefined
      ? { totalItems, offset, limit }
      : undefined;

    return { items, pagination, excludedUnverified };
  }

  /**
   * Fetches a single asset envelope. Performs no verification, and returns the BARE
   * envelope type: nothing that reads envelope fields will accept it until
   * {@link verifyEnvelope} has branded it.
   */
  private async fetchEnvelope(
    assetId: number
  ): Promise<SignedWhitelistedAssetEnvelope> {
    const response = await this.api.whitelistServiceGetWhitelistedContract({
      id: String(assetId),
    });

    const dto = response.result;
    if (dto == null) {
      throw new NotFoundError(`Whitelisted asset ${assetId} not found`);
    }

    return this.mapDtoToEnvelope(dto, assetId);
  }

  /**
   * Runs the 5-step verification. The single place this service checks an envelope,
   * and the only source of a `Verified<SignedWhitelistedAssetEnvelope>` — so a read
   * path that skips it cannot produce one, and will not compile.
   */
  private verifyEnvelope(
    envelope: SignedWhitelistedAssetEnvelope
  ): WhitelistedAssetVerificationResult {
    return this.verifier.verify(
      envelope,
      this.rulesContainerDecoder,
      this.userSignaturesDecoder
    );
  }

  /**
   * Maps a DTO to a SignedWhitelistedAssetEnvelope.
   *
   * Carries no security meaning on its own — the envelope it returns is
   * unverified until {@link verifyEnvelope} has run over it.
   */
  private mapDtoToEnvelope(
    dto: {
      id?: string;
      metadata?: { hash?: string; payloadAsString?: string };
      rulesContainer?: string;
      rulesSignatures?: string;
      signedContractAddress?: {
        payload?: string;
        signatures?: Array<{
          signature?: { userId?: string; signature?: string; comment?: string };
          hashes?: string[];
        }>;
      };
      blockchain?: string;
      network?: string;
    },
    assetId: number
  ): SignedWhitelistedAssetEnvelope {
    return {
      id: assetId,
      metadata: {
        hash: dto.metadata?.hash ?? "",
        payloadAsString: dto.metadata?.payloadAsString ?? "",
      },
      rulesContainerBase64: dto.rulesContainer ?? "",
      rulesSignaturesBase64: dto.rulesSignatures ?? "",
      signedContractAddress: {
        payload: dto.signedContractAddress?.payload,
        signatures: (dto.signedContractAddress?.signatures ?? []).map((sig) => ({
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
    };
  }
}
