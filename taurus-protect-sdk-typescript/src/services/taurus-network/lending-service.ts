/**
 * Lending service for Taurus Network in Taurus-PROTECT SDK.
 *
 * Provides methods for managing Taurus Network lending offers and agreements.
 */

import { NotFoundError, ValidationError } from '../../errors';
import type { TaurusNetworkLendingApi } from '../../internal/openapi/apis/TaurusNetworkLendingApi';
import type {
  TgvalidatordLendingAgreement,
  TgvalidatordTnLendingOffer,
  TgvalidatordLendingAgreementAttachment,
} from '../../internal/openapi/models/index';
import { BaseService } from '../base';

/**
 * Cursor-based pagination information.
 */
export interface CursorPagination {
  currentPage?: string;
  hasNext: boolean;
  hasPrevious: boolean;
}

/**
 * Currency information.
 */
export interface CurrencyInfo {
  id?: string;
  name?: string;
  symbol?: string;
  decimals?: number;
}

/**
 * Currency collateral requirement.
 */
export interface CurrencyCollateralRequirement {
  blockchain?: string;
  network?: string;
  ratio?: string;
  currencyInfo?: CurrencyInfo;
}

/**
 * Lending collateral requirement.
 */
export interface LendingCollateralRequirement {
  acceptedCurrencies?: CurrencyCollateralRequirement[];
}

/**
 * A lending offer in the Taurus Network.
 */
export interface LendingOffer {
  id: string;
  participantId?: string;
  currencyInfo?: CurrencyInfo;
  amount?: string;
  amountMainUnit?: string;
  annualPercentageYield?: string;
  annualPercentageYieldMainUnit?: string;
  duration?: string;
  blockchain?: string;
  network?: string;
  collateralRequirement?: LendingCollateralRequirement;
  originCreatedAt?: Date;
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Lending agreement collateral.
 */
export interface LendingAgreementCollateral {
  id?: string;
  lendingAgreementId?: string;
  lenderParticipantId?: string;
  borrowerParticipantId?: string;
  currencyId?: string;
  currencyInfo?: CurrencyInfo;
  amount?: string;
  amountMainUnit?: string;
  status?: string;
  pledgeId?: string;
  pledgeActionId?: string;
  sharedAddressId?: string;
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Lending agreement transaction.
 */
export interface LendingAgreementTransaction {
  id?: string;
  lendingAgreementId?: string;
  amount?: string;
  currencyId?: string;
  requestId?: string;
  transactionId?: string;
  transactionHash?: string;
  transactionBlockNumber?: string;
  type?: string;
  amountMainUnit?: string;
  currencyInfo?: CurrencyInfo;
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * A lending agreement in the Taurus Network.
 */
export interface LendingAgreement {
  id: string;
  lenderParticipantId?: string;
  borrowerParticipantId?: string;
  lendingOfferId?: string;
  currencyId?: string;
  currencyInfo?: CurrencyInfo;
  amount?: string;
  amountMainUnit?: string;
  annualYield?: string;
  annualYieldMainUnit?: string;
  duration?: string;
  status?: string;
  workflowId?: string;
  lenderSharedAddressId?: string;
  borrowerSharedAddressId?: string;
  startLoanDate?: Date;
  repaymentDueDate?: Date;
  collaterals?: LendingAgreementCollateral[];
  transactions?: LendingAgreementTransaction[];
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Lending agreement attachment.
 */
export interface LendingAgreementAttachment {
  id?: string;
  lendingAgreementId?: string;
  uploaderParticipantId?: string;
  name?: string;
  type?: string;
  contentType?: string;
  value?: string;
  fileSize?: string;
  createdAt?: Date;
  updatedAt?: Date;
}

/**
 * Options for listing lending offers.
 */
export interface ListLendingOffersOptions {
  currencyIds?: string[];
  participantId?: string;
  duration?: string;
  sortOrder?: string;
  pageSize?: number;
  currentPage?: string;
  pageRequest?: string;
}

/**
 * Options for listing lending agreements.
 */
export interface ListLendingAgreementsOptions {
  ids?: string[];
  sortOrder?: string;
  pageSize?: number;
  currentPage?: string;
  pageRequest?: string;
}

/**
 * Collateral configuration for creating a lending agreement.
 */
export interface CollateralRequest {
  sourceSharedAddressId: string;
  currencyId: string;
  amount: string;
}

/**
 * Request to create a lending offer.
 */
export interface CreateLendingOfferRequest {
  currencyId: string;
  amount: string;
  annualPercentageYield: string;
  duration: string;
  collateralRequirements?: Array<{
    ratio?: string;
    currencyId?: string;
  }>;
}

/**
 * Request to create a lending agreement.
 */
export interface CreateLendingAgreementRequest {
  borrowerSharedAddressId: string;
  lendingOfferId?: string;
  lenderParticipantId?: string;
  currencyId?: string;
  amount?: string;
  annualPercentageYield?: string;
  duration?: string;
  collaterals?: CollateralRequest[];
}

/**
 * Request to update a lending agreement.
 */
export interface UpdateLendingAgreementRequest {
  lenderSharedAddressId: string;
}

/**
 * Request to repay a lending agreement.
 */
export interface RepayLendingAgreementRequest {
  repayerSharedAddressId: string;
}

/**
 * Request to create a lending agreement attachment.
 */
export interface CreateLendingAgreementAttachmentRequest {
  name: string;
  value: string;
  contentType: string;
  type?: string;
}

/**
 * Maps currency info from DTO.
 */
function currencyInfoFromDto(dto?: { id?: string; name?: string; symbol?: string; decimals?: string }): CurrencyInfo | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    id: dto.id,
    name: dto.name,
    symbol: dto.symbol,
    decimals: dto.decimals ? parseInt(String(dto.decimals), 10) : undefined,
  };
}

/**
 * Maps a DTO to a LendingOffer.
 */
function lendingOfferFromDto(dto?: TgvalidatordTnLendingOffer): LendingOffer | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id ?? '',
    participantId: dto.participantID,
    currencyInfo: currencyInfoFromDto(dto.currencyInfo),
    amount: dto.amount,
    amountMainUnit: dto.amountMainUnit,
    annualPercentageYield: dto.annualPercentageYield,
    annualPercentageYieldMainUnit: dto.annualPercentageYieldMainUnit,
    duration: dto.duration,
    blockchain: dto.blockchain,
    network: dto.network,
    collateralRequirement: dto.collateralRequirement
      ? {
          acceptedCurrencies: dto.collateralRequirement.acceptedCurrencies?.map((c) => ({
            blockchain: c.blockchain,
            network: c.network,
            ratio: c.ratio,
            currencyInfo: currencyInfoFromDto(c.currencyInfo),
          })),
        }
      : undefined,
    originCreatedAt: dto.originCreatedAt,
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
  };
}

/**
 * Maps a DTO to a LendingAgreement.
 */
function lendingAgreementFromDto(dto?: TgvalidatordLendingAgreement): LendingAgreement | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id ?? '',
    lenderParticipantId: dto.lenderParticipantID,
    borrowerParticipantId: dto.borrowerParticipantID,
    lendingOfferId: dto.lendingOfferID,
    currencyId: dto.currencyID,
    currencyInfo: currencyInfoFromDto(dto.currencyInfo),
    amount: dto.amount,
    amountMainUnit: dto.amountMainUnit,
    annualYield: dto.annualYield,
    annualYieldMainUnit: dto.annualYieldMainUnit,
    duration: dto.duration,
    status: dto.status,
    workflowId: dto.workflowID,
    lenderSharedAddressId: dto.lenderSharedAddressID,
    borrowerSharedAddressId: dto.borrowerSharedAddressID,
    startLoanDate: dto.startLoanDate,
    repaymentDueDate: dto.repaymentDueDate,
    collaterals: dto.lendingAgreementCollaterals?.map((c) => ({
      id: c.id,
      lendingAgreementId: c.lendingAgreementID,
      lenderParticipantId: c.lenderParticipantID,
      borrowerParticipantId: c.borrowerParticipantID,
      currencyId: c.currencyID,
      currencyInfo: currencyInfoFromDto(c.currencyInfo),
      amount: c.amount,
      amountMainUnit: c.amountMainUnit,
      status: c.status,
      pledgeId: c.pledgeID,
      pledgeActionId: c.pledgeActionID,
      sharedAddressId: c.sharedAddressID,
      createdAt: c.createdAt,
      updatedAt: c.updatedAt,
    })),
    transactions: dto.lendingAgreementTransactions?.map((t) => ({
      id: t.id,
      lendingAgreementId: t.lendingAgreementID,
      amount: t.amount,
      currencyId: t.currencyID,
      requestId: t.requestID,
      transactionId: t.transactionID,
      transactionHash: t.transactionHash,
      transactionBlockNumber: t.transactionBlockNumber,
      type: t.type,
      amountMainUnit: t.amountMainUnit,
      currencyInfo: currencyInfoFromDto(t.currencyInfo),
      createdAt: t.createdAt,
      updatedAt: t.updatedAt,
    })),
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
  };
}

/**
 * Maps a DTO to a LendingAgreementAttachment.
 */
function lendingAgreementAttachmentFromDto(
  dto?: TgvalidatordLendingAgreementAttachment
): LendingAgreementAttachment | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: dto.id,
    lendingAgreementId: dto.lendingAgreementID,
    uploaderParticipantId: dto.uploaderParticipantID,
    name: dto.name,
    type: dto.type,
    contentType: dto.contentType,
    value: dto.value,
    fileSize: dto.fileSize,
    createdAt: dto.createdAt,
    updatedAt: dto.updatedAt,
  };
}

/**
 * Extracts cursor pagination from response.
 */
function extractCursorPagination(cursor?: {
  currentPage?: string;
  hasNext?: boolean;
  hasPrevious?: boolean;
}): CursorPagination | undefined {
  if (!cursor) {
    return undefined;
  }

  return {
    currentPage: cursor.currentPage,
    hasNext: cursor.hasNext ?? false,
    hasPrevious: cursor.hasPrevious ?? false,
  };
}

/**
 * Service for Taurus Network lending operations.
 *
 * Provides methods to manage lending offers and agreements between
 * Taurus Network participants.
 *
 * @example
 * ```typescript
 * // List lending offers
 * const { offers, pagination } = await lendingService.listLendingOffers();
 * for (const offer of offers) {
 *   console.log(`${offer.id}: ${offer.annualPercentageYieldMainUnit} APY`);
 * }
 *
 * // Get single agreement
 * const agreement = await lendingService.getLendingAgreement('123');
 * console.log(`Status: ${agreement.status}`);
 * ```
 */
export class LendingService extends BaseService {
  private readonly lendingApi: TaurusNetworkLendingApi;

  /**
   * Creates a new LendingService instance.
   *
   * @param lendingApi - The TaurusNetworkLendingApi instance from the OpenAPI client
   */
  constructor(lendingApi: TaurusNetworkLendingApi) {
    super();
    this.lendingApi = lendingApi;
  }

  // =========================================================================
  // Lending Offer Methods
  // =========================================================================

  /**
   * Gets a lending offer by ID.
   *
   * @param offerId - The lending offer ID to retrieve
   * @returns The lending offer
   * @throws {@link ValidationError} If offerId is empty
   * @throws {@link NotFoundError} If offer not found
   * @throws {@link APIError} If API request fails
   */
  async getLendingOffer(offerId: string): Promise<LendingOffer> {
    if (!offerId || offerId.trim() === '') {
      throw new ValidationError('offerId is required');
    }

    return this.execute(async () => {
      const response = await this.lendingApi.taurusNetworkServiceGetLendingOffer({
        offerID: offerId,
      });

      const offer = lendingOfferFromDto(response.lendingOffer);
      if (!offer) {
        throw new NotFoundError(`Lending offer ${offerId} not found`);
      }

      return offer;
    });
  }

  /**
   * Lists lending offers.
   *
   * @param options - Optional filtering and pagination options
   * @returns Offers list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async listLendingOffers(
    options?: ListLendingOffersOptions
  ): Promise<{ offers: LendingOffer[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.lendingApi.taurusNetworkServiceGetLendingOffers({
        currencyIDsCurrencyIDs: options?.currencyIds,
        participantID: options?.participantId,
        duration: options?.duration,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const offers: LendingOffer[] = [];
      if (response.lendingOffers) {
        for (const dto of response.lendingOffers) {
          const offer = lendingOfferFromDto(dto);
          if (offer) {
            offers.push(offer);
          }
        }
      }

      return {
        offers,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Creates a new lending offer.
   *
   * Creates a lending offer that other participants can accept
   * to create a lending agreement.
   *
   * @param request - Lending offer creation parameters
   * @returns The created lending offer ID
   * @throws {@link ValidationError} If required fields are missing
   * @throws {@link APIError} If API request fails
   */
  async createLendingOffer(request: CreateLendingOfferRequest): Promise<string> {
    if (!request.currencyId || request.currencyId.trim() === '') {
      throw new ValidationError('currencyId is required');
    }
    if (!request.amount || request.amount.trim() === '') {
      throw new ValidationError('amount is required');
    }
    if (!request.annualPercentageYield || request.annualPercentageYield.trim() === '') {
      throw new ValidationError('annualPercentageYield is required');
    }
    if (!request.duration || request.duration.trim() === '') {
      throw new ValidationError('duration is required');
    }

    return this.execute(async () => {
      const response = await this.lendingApi.taurusNetworkServiceCreateLendingOffer({
        body: {
          currencyID: request.currencyId,
          amount: request.amount,
          annualPercentageYield: request.annualPercentageYield,
          duration: request.duration,
          collateralRequirement: request.collateralRequirements?.map((c) => ({
            ratio: c.ratio,
            currencyID: c.currencyId,
          })),
        },
      });

      return response.offerID ?? '';
    });
  }

  /**
   * Deletes a specific lending offer.
   *
   * @param offerId - The lending offer ID to delete
   * @throws {@link ValidationError} If offerId is empty
   * @throws {@link APIError} If API request fails
   */
  async deleteLendingOffer(offerId: string): Promise<void> {
    if (!offerId || offerId.trim() === '') {
      throw new ValidationError('offerId is required');
    }

    return this.execute(async () => {
      await this.lendingApi.taurusNetworkServiceDeleteLendingOffer({
        offerID: offerId,
      });
    });
  }

  /**
   * Deletes all lending offers for the current participant.
   *
   * @throws {@link APIError} If API request fails
   */
  async deleteLendingOffers(): Promise<void> {
    return this.execute(async () => {
      await this.lendingApi.taurusNetworkServiceDeleteLendingOffers();
    });
  }

  // =========================================================================
  // Lending Agreement Methods
  // =========================================================================

  /**
   * Gets a lending agreement by ID.
   *
   * @param lendingAgreementId - The lending agreement ID to retrieve
   * @returns The lending agreement
   * @throws {@link ValidationError} If lendingAgreementId is empty
   * @throws {@link NotFoundError} If agreement not found
   * @throws {@link APIError} If API request fails
   */
  async getLendingAgreement(lendingAgreementId: string): Promise<LendingAgreement> {
    if (!lendingAgreementId || lendingAgreementId.trim() === '') {
      throw new ValidationError('lendingAgreementId is required');
    }

    return this.execute(async () => {
      const response = await this.lendingApi.taurusNetworkServiceGetLendingAgreement({
        lendingAgreementID: lendingAgreementId,
      });

      const agreement = lendingAgreementFromDto(response.result);
      if (!agreement) {
        throw new NotFoundError(`Lending agreement ${lendingAgreementId} not found`);
      }

      return agreement;
    });
  }

  /**
   * Lists lending agreements.
   *
   * @param options - Optional filtering and pagination options
   * @returns Agreements list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async listLendingAgreements(
    options?: ListLendingAgreementsOptions
  ): Promise<{ agreements: LendingAgreement[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.lendingApi.taurusNetworkServiceGetLendingAgreements({
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const agreements: LendingAgreement[] = [];
      if (response.lendingAgreements) {
        for (const dto of response.lendingAgreements) {
          const agreement = lendingAgreementFromDto(dto);
          if (agreement) {
            agreements.push(agreement);
          }
        }
      }

      return {
        agreements,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Lists lending agreements pending approval.
   *
   * @param options - Optional filtering and pagination options
   * @returns Agreements list and pagination info
   * @throws {@link APIError} If API request fails
   */
  async listLendingAgreementsForApproval(
    options?: ListLendingAgreementsOptions
  ): Promise<{ agreements: LendingAgreement[]; pagination?: CursorPagination }> {
    return this.execute(async () => {
      const response = await this.lendingApi.taurusNetworkServiceGetLendingAgreementsForApproval({
        ids: options?.ids,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize ? String(options.pageSize) : undefined,
      });

      const agreements: LendingAgreement[] = [];
      if (response.result) {
        for (const dto of response.result) {
          const agreement = lendingAgreementFromDto(dto);
          if (agreement) {
            agreements.push(agreement);
          }
        }
      }

      return {
        agreements,
        pagination: extractCursorPagination(response.cursor),
      };
    });
  }

  /**
   * Creates a new lending agreement.
   *
   * Creates a lending agreement. Once approved by the lender and borrower,
   * the agreement will automatically execute the loan lifecycle.
   * This should be called by the borrower.
   *
   * @param request - Lending agreement creation parameters
   * @returns The created lending agreement ID
   * @throws {@link ValidationError} If required fields are missing
   * @throws {@link APIError} If API request fails
   */
  async createLendingAgreement(request: CreateLendingAgreementRequest): Promise<string> {
    if (!request.borrowerSharedAddressId || request.borrowerSharedAddressId.trim() === '') {
      throw new ValidationError('borrowerSharedAddressId is required');
    }

    return this.execute(async () => {
      const response = await this.lendingApi.taurusNetworkServiceCreateLendingAgreement({
        body: {
          borrowerSharedAddressID: request.borrowerSharedAddressId,
          lendingOfferID: request.lendingOfferId,
          lenderParticipantID: request.lenderParticipantId,
          currencyID: request.currencyId,
          amount: request.amount,
          annualPercentageYield: request.annualPercentageYield,
          duration: request.duration,
          collaterals: request.collaterals?.map((c) => ({
            sourceSharedAddressID: c.sourceSharedAddressId,
            currencyID: c.currencyId,
            amount: c.amount,
          })),
        },
      });

      return response.lendingAgreementID ?? '';
    });
  }

  /**
   * Updates a lending agreement.
   *
   * Updates the lender's shared address for a lending agreement.
   * This should be called by the lender before approving.
   *
   * @param lendingAgreementId - The lending agreement ID to update
   * @param request - Update parameters
   * @throws {@link ValidationError} If required fields are missing
   * @throws {@link APIError} If API request fails
   */
  async updateLendingAgreement(
    lendingAgreementId: string,
    request: UpdateLendingAgreementRequest
  ): Promise<void> {
    if (!lendingAgreementId || lendingAgreementId.trim() === '') {
      throw new ValidationError('lendingAgreementId is required');
    }
    if (!request.lenderSharedAddressId || request.lenderSharedAddressId.trim() === '') {
      throw new ValidationError('lenderSharedAddressId is required');
    }

    return this.execute(async () => {
      await this.lendingApi.taurusNetworkServiceUpdateLendingAgreement({
        lendingAgreementID: lendingAgreementId,
        body: {
          lenderSharedAddressID: request.lenderSharedAddressId,
        },
      });
    });
  }

  /**
   * Records repayment for a lending agreement.
   *
   * Records a repayment of the loan from the borrower to the lender.
   *
   * @param lendingAgreementId - The lending agreement ID
   * @param request - Repayment parameters
   * @throws {@link ValidationError} If required fields are missing
   * @throws {@link APIError} If API request fails
   */
  async repayLendingAgreement(
    lendingAgreementId: string,
    request: RepayLendingAgreementRequest
  ): Promise<void> {
    if (!lendingAgreementId || lendingAgreementId.trim() === '') {
      throw new ValidationError('lendingAgreementId is required');
    }
    if (!request.repayerSharedAddressId || request.repayerSharedAddressId.trim() === '') {
      throw new ValidationError('repayerSharedAddressId is required');
    }

    return this.execute(async () => {
      await this.lendingApi.taurusNetworkServiceRepayLendingAgreement({
        lendingAgreementID: lendingAgreementId,
        body: {
          repayerSharedAddressID: request.repayerSharedAddressId,
        },
      });
    });
  }

  /**
   * Cancels a lending agreement.
   *
   * Cancels a lending agreement when it is not yet approved by the lender.
   *
   * @param lendingAgreementId - The lending agreement ID to cancel
   * @throws {@link ValidationError} If lendingAgreementId is empty
   * @throws {@link APIError} If API request fails
   */
  async cancelLendingAgreement(lendingAgreementId: string): Promise<void> {
    if (!lendingAgreementId || lendingAgreementId.trim() === '') {
      throw new ValidationError('lendingAgreementId is required');
    }

    return this.execute(async () => {
      await this.lendingApi.taurusNetworkServiceCancelLendingAgreement({
        lendingAgreementID: lendingAgreementId,
        body: {},
      });
    });
  }

  // =========================================================================
  // Lending Agreement Attachment Methods
  // =========================================================================

  /**
   * Adds an attachment to a lending agreement.
   *
   * @param lendingAgreementId - The lending agreement ID
   * @param request - Attachment creation parameters
   * @returns The created attachment ID
   * @throws {@link ValidationError} If required fields are missing
   * @throws {@link APIError} If API request fails
   */
  async createLendingAgreementAttachment(
    lendingAgreementId: string,
    request: CreateLendingAgreementAttachmentRequest
  ): Promise<string> {
    if (!lendingAgreementId || lendingAgreementId.trim() === '') {
      throw new ValidationError('lendingAgreementId is required');
    }
    if (!request.name || request.name.trim() === '') {
      throw new ValidationError('name is required');
    }
    if (!request.value || request.value.trim() === '') {
      throw new ValidationError('value is required');
    }
    if (!request.contentType || request.contentType.trim() === '') {
      throw new ValidationError('contentType is required');
    }

    return this.execute(async () => {
      const response = (await this.lendingApi.taurusNetworkServiceCreateLendingAgreementAttachment({
        lendingAgreementID: lendingAgreementId,
        body: {
          name: request.name,
          value: request.value,
          contentType: request.contentType,
          type: request.type,
        },
      })) as { attachmentID?: string };

      return response.attachmentID ?? '';
    });
  }

  /**
   * Lists attachments for a lending agreement.
   *
   * @param lendingAgreementId - The lending agreement ID
   * @returns List of attachments
   * @throws {@link ValidationError} If lendingAgreementId is empty
   * @throws {@link APIError} If API request fails
   */
  async listLendingAgreementAttachments(
    lendingAgreementId: string
  ): Promise<LendingAgreementAttachment[]> {
    if (!lendingAgreementId || lendingAgreementId.trim() === '') {
      throw new ValidationError('lendingAgreementId is required');
    }

    return this.execute(async () => {
      const response = await this.lendingApi.taurusNetworkServiceGetLendingAgreementAttachments({
        lendingAgreementID: lendingAgreementId,
      });

      const attachments: LendingAgreementAttachment[] = [];
      if (response.result) {
        for (const dto of response.result) {
          const attachment = lendingAgreementAttachmentFromDto(dto);
          if (attachment) {
            attachments.push(attachment);
          }
        }
      }

      return attachments;
    });
  }
}
