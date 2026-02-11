/**
 * Lending models for Taurus Network.
 *
 * This module provides domain models for lending operations,
 * including lending offers and lending agreements between participants.
 */

/**
 * Lending agreement status enum.
 */
export enum LendingAgreementStatus {
  PENDING = "PENDING",
  ACTIVE = "ACTIVE",
  COMPLETED = "COMPLETED",
  CANCELED = "CANCELED",
  DEFAULTED = "DEFAULTED",
}

/**
 * Collateral requirement for a lending offer.
 */
export interface LendingCollateralRequirement {
  /** Currency ID. */
  readonly currencyId: string | undefined;
  /** Collateral percentage. */
  readonly percentage: string | undefined;
}

/**
 * Currency information for lending.
 */
export interface LendingCurrencyInfo {
  /** Currency ID. */
  readonly id: string | undefined;
  /** Currency symbol. */
  readonly symbol: string | undefined;
  /** Currency name. */
  readonly name: string | undefined;
  /** Number of decimals. */
  readonly decimals: number | undefined;
}

/**
 * Lending offer from a participant.
 *
 * Represents an offer to lend assets at a specified yield.
 */
export interface LendingOffer {
  /** Unique offer identifier. */
  readonly id: string | undefined;
  /** APY in basis points (525000 = 5.25%). */
  readonly annualPercentageYield: string | undefined;
  /** Loan duration (e.g., "3M", "1Y"). */
  readonly duration: string | undefined;
  /** Required collateral. */
  readonly collateralRequirement: LendingCollateralRequirement | undefined;
  /** Lender participant ID. */
  readonly participantId: string | undefined;
  /** Blockchain name. */
  readonly blockchain: string | undefined;
  /** Network name. */
  readonly network: string | undefined;
  /** Currency argument 1. */
  readonly arg1: string | undefined;
  /** Currency argument 2. */
  readonly arg2: string | undefined;
  /** Currency details. */
  readonly currencyInfo: LendingCurrencyInfo | undefined;
  /** APY as human-readable string. */
  readonly annualPercentageYieldMainUnit: string | undefined;
  /** Original creation on network. */
  readonly originCreatedAt: Date | undefined;
  /** Local creation timestamp. */
  readonly createdAt: Date | undefined;
  /** Local update timestamp. */
  readonly updatedAt: Date | undefined;
  /** Available loan amount. */
  readonly amount: string | undefined;
  /** Amount in main unit. */
  readonly amountMainUnit: string | undefined;
}

/**
 * Collateral provided for a lending agreement.
 */
export interface LendingAgreementCollateral {
  /** Currency ID. */
  readonly currencyId: string | undefined;
  /** Collateral amount. */
  readonly amount: string | undefined;
  /** Associated pledge ID. */
  readonly pledgeId: string | undefined;
}

/**
 * Transaction related to a lending agreement.
 */
export interface LendingAgreementTransaction {
  /** Transaction ID. */
  readonly id: string | undefined;
  /** Transaction type. */
  readonly transactionType: string | undefined;
  /** Transaction amount. */
  readonly amount: string | undefined;
  /** Transaction status. */
  readonly status: string | undefined;
  /** Blockchain transaction hash. */
  readonly txHash: string | undefined;
}

/**
 * Lending agreement between participants.
 *
 * Represents an active loan between a lender and borrower.
 */
export interface LendingAgreement {
  /** Unique agreement identifier. */
  readonly id: string | undefined;
  /** Lender participant ID. */
  readonly lenderParticipantId: string | undefined;
  /** Borrower participant ID. */
  readonly borrowerParticipantId: string | undefined;
  /** Associated lending offer ID. */
  readonly lendingOfferId: string | undefined;
  /** Loan amount. */
  readonly amount: string | undefined;
  /** Currency ID. */
  readonly currencyId: string | undefined;
  /** Annual yield in basis points. */
  readonly annualYield: string | undefined;
  /** Agreement status. */
  readonly status: string | undefined;
  /** Loan duration. */
  readonly duration: string | undefined;
  /** When loan started. */
  readonly startLoanDate: Date | undefined;
  /** Associated workflow ID. */
  readonly workflowId: string | undefined;
  /** Borrower's shared address ID. */
  readonly borrowerSharedAddressId: string | undefined;
  /** Lender's shared address ID. */
  readonly lenderSharedAddressId: string | undefined;
  /** Provided collateral. */
  readonly collaterals: LendingAgreementCollateral[];
  /** Related transactions. */
  readonly transactions: LendingAgreementTransaction[];
  /** Creation timestamp. */
  readonly createdAt: Date | undefined;
  /** Update timestamp. */
  readonly updatedAt: Date | undefined;
  /** Yield as percentage. */
  readonly annualYieldMainUnit: string | undefined;
  /** Currency details. */
  readonly currencyInfo: LendingCurrencyInfo | undefined;
  /** Amount in main unit. */
  readonly amountMainUnit: string | undefined;
  /** When repayment is due. */
  readonly repaymentDueDate: Date | undefined;
}

/**
 * Attachment on a lending agreement.
 *
 * Can be embedded content (base64) or external link.
 */
export interface LendingAgreementAttachment {
  /** Attachment identifier. */
  readonly id: string | undefined;
  /** Associated agreement ID. */
  readonly lendingAgreementId: string | undefined;
  /** Participant who uploaded. */
  readonly uploaderParticipantId: string | undefined;
  /** Attachment name/filename. */
  readonly name: string | undefined;
  /** Attachment type (EMBEDDED, EXTERNAL_LINK). */
  readonly type: string | undefined;
  /** MIME type. */
  readonly contentType: string | undefined;
  /** Content (base64) or URL. */
  readonly value: string | undefined;
  /** Size in bytes (for embedded). */
  readonly fileSize: string | undefined;
  /** Creation timestamp. */
  readonly createdAt: Date | undefined;
  /** Update timestamp. */
  readonly updatedAt: Date | undefined;
}

// Request models

/**
 * Request to create a lending offer.
 */
export interface CreateLendingOfferRequest {
  /** Currency identifier. */
  readonly currencyId: string;
  /** Loan amount in smallest unit. */
  readonly amount: string;
  /** APY in basis points. */
  readonly annualPercentageYield: string;
  /** Loan duration (e.g., "3M", "1Y"). */
  readonly duration: string;
  /** Collateral requirements. */
  readonly collateralRequirement?: LendingCollateralRequirement;
}

/**
 * Request to update a lending offer.
 */
export interface UpdateLendingOfferRequest {
  /** New loan amount. */
  readonly amount?: string;
  /** New APY. */
  readonly annualPercentageYield?: string;
}

/**
 * Request to create a lending agreement.
 */
export interface CreateLendingAgreementRequest {
  /** Lending offer ID. */
  readonly lendingOfferId: string;
  /** Loan amount. */
  readonly amount: string;
  /** Borrower's shared address ID. */
  readonly borrowerSharedAddressId: string;
  /** Lender's shared address ID. */
  readonly lenderSharedAddressId: string;
  /** Collateral pledge IDs. */
  readonly collateralPledgeIds?: string[];
}

/**
 * Request to create a lending agreement attachment.
 */
export interface CreateLendingAgreementAttachmentRequest {
  /** Attachment name. */
  readonly name: string;
  /** Type (EMBEDDED, EXTERNAL_LINK). */
  readonly type: string;
  /** Content (base64) or URL. */
  readonly value: string;
  /** MIME type. */
  readonly contentType?: string;
}

// Filter options

/**
 * Options for listing lending offers.
 */
export interface ListLendingOffersOptions {
  /** Maximum items to return (default: 50, max: 1000). */
  readonly limit?: number;
  /** Number of items to skip. */
  readonly offset?: number;
  /** Filter by participant ID. */
  readonly participantId?: string;
  /** Filter by currency. */
  readonly currencyId?: string;
}

/**
 * Options for listing lending agreements.
 */
export interface ListLendingAgreementsOptions {
  /** Maximum items to return (default: 50, max: 1000). */
  readonly limit?: number;
  /** Number of items to skip. */
  readonly offset?: number;
  /** Filter by statuses. */
  readonly statuses?: string[];
  /** Filter by lender. */
  readonly lenderParticipantId?: string;
  /** Filter by borrower. */
  readonly borrowerParticipantId?: string;
}
