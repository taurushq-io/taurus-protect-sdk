"""Lending service for Taurus Network operations."""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from typing import TYPE_CHECKING, Any, Dict, List, Optional, Tuple

from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


# =============================================================================
# Data Classes
# =============================================================================


@dataclass
class CurrencyCollateralRequirement:
    """
    Currency collateral requirement for a lending offer.

    Attributes:
        currency_id: Currency ID for the collateral.
        percentage: Required collateral percentage.
    """

    currency_id: str = ""
    percentage: str = ""


@dataclass
class LendingCollateralRequirement:
    """
    Collateral requirements for a lending offer.

    Attributes:
        accepted_currencies: List of accepted currencies and their requirements.
    """

    accepted_currencies: List[CurrencyCollateralRequirement] = field(default_factory=list)


@dataclass
class CurrencyInfo:
    """
    Currency information.

    Attributes:
        id: Currency ID.
        name: Currency name.
        symbol: Currency symbol.
        decimals: Number of decimals.
    """

    id: str = ""
    name: str = ""
    symbol: str = ""
    decimals: int = 0


@dataclass
class LendingOffer:
    """
    A lending offer in the Taurus Network.

    Represents an offer to lend cryptocurrency at specified terms.

    Attributes:
        id: Unique offer identifier.
        participant_id: Participant who created the offer.
        currency_id: Currency being offered for lending.
        currency_info: Currency details.
        amount: Amount available for lending.
        amount_main_unit: Amount in human-readable main unit.
        annual_percentage_yield: Interest rate (5 decimal precision, e.g., 525000 = 5.25%).
        annual_percentage_yield_main_unit: Human-readable APY.
        duration: Loan duration.
        blockchain: Blockchain name.
        network: Network name.
        collateral_requirement: Required collateral configuration.
        created_at: Creation timestamp.
        updated_at: Last update timestamp.
        origin_created_at: Original creation timestamp.
    """

    id: str = ""
    participant_id: str = ""
    currency_id: str = ""
    currency_info: Optional[CurrencyInfo] = None
    amount: str = ""
    amount_main_unit: str = ""
    annual_percentage_yield: str = ""
    annual_percentage_yield_main_unit: str = ""
    duration: str = ""
    blockchain: str = ""
    network: str = ""
    collateral_requirement: Optional[LendingCollateralRequirement] = None
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None
    origin_created_at: Optional[datetime] = None


@dataclass
class LendingAgreementCollateral:
    """
    Collateral for a lending agreement.

    Attributes:
        id: Unique collateral identifier.
        lending_agreement_id: Associated agreement ID.
        lender_participant_id: Lender participant ID.
        borrower_participant_id: Borrower participant ID.
        currency_id: Collateral currency ID.
        currency_info: Currency details.
        amount: Collateral amount.
        amount_main_unit: Amount in human-readable main unit.
        status: Collateral status.
        pledge_id: Associated pledge ID.
        pledge_action_id: Associated pledge action ID.
        shared_address_id: Shared address holding collateral.
        created_at: Creation timestamp.
        updated_at: Last update timestamp.
    """

    id: str = ""
    lending_agreement_id: str = ""
    lender_participant_id: str = ""
    borrower_participant_id: str = ""
    currency_id: str = ""
    currency_info: Optional[CurrencyInfo] = None
    amount: str = ""
    amount_main_unit: str = ""
    status: str = ""
    pledge_id: str = ""
    pledge_action_id: str = ""
    shared_address_id: str = ""
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


@dataclass
class LendingAgreementTransaction:
    """
    Transaction associated with a lending agreement.

    Attributes:
        id: Transaction identifier.
        type: Transaction type (e.g., "DISBURSEMENT", "REPAYMENT").
        amount: Transaction amount.
        status: Transaction status.
        created_at: Creation timestamp.
    """

    id: str = ""
    type: str = ""
    amount: str = ""
    status: str = ""
    created_at: Optional[datetime] = None


@dataclass
class LendingAgreement:
    """
    A lending agreement in the Taurus Network.

    Represents an agreement between a lender and borrower for cryptocurrency lending.

    Attributes:
        id: Unique agreement identifier.
        lender_participant_id: Lender participant ID.
        borrower_participant_id: Borrower participant ID.
        lending_offer_id: Associated lending offer ID.
        currency_id: Loan currency ID.
        currency_info: Currency details.
        amount: Loan amount.
        amount_main_unit: Amount in human-readable main unit.
        annual_yield: Interest rate.
        annual_yield_main_unit: Human-readable interest rate.
        duration: Loan duration.
        status: Agreement status.
        workflow_id: Workflow ID managing the agreement.
        lender_shared_address_id: Lender's shared address.
        borrower_shared_address_id: Borrower's shared address.
        start_loan_date: Loan start date.
        repayment_due_date: Repayment due date.
        collaterals: List of collaterals.
        transactions: List of related transactions.
        created_at: Creation timestamp.
        updated_at: Last update timestamp.
    """

    id: str = ""
    lender_participant_id: str = ""
    borrower_participant_id: str = ""
    lending_offer_id: str = ""
    currency_id: str = ""
    currency_info: Optional[CurrencyInfo] = None
    amount: str = ""
    amount_main_unit: str = ""
    annual_yield: str = ""
    annual_yield_main_unit: str = ""
    duration: str = ""
    status: str = ""
    workflow_id: str = ""
    lender_shared_address_id: str = ""
    borrower_shared_address_id: str = ""
    start_loan_date: Optional[datetime] = None
    repayment_due_date: Optional[datetime] = None
    collaterals: List[LendingAgreementCollateral] = field(default_factory=list)
    transactions: List[LendingAgreementTransaction] = field(default_factory=list)
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


@dataclass
class LendingAgreementAttachment:
    """
    Attachment on a lending agreement.

    Attributes:
        id: Unique attachment identifier.
        lending_agreement_id: Associated agreement ID.
        uploader_participant_id: Participant who uploaded the attachment.
        name: Attachment name.
        type: Attachment type (EMBEDDED or EXTERNAL_LINK).
        content_type: MIME type (e.g., "application/pdf").
        value: Content (base64 for EMBEDDED, URL for EXTERNAL_LINK).
        file_size: File size in bytes (for EMBEDDED type).
        created_at: Creation timestamp.
        updated_at: Last update timestamp.
    """

    id: str = ""
    lending_agreement_id: str = ""
    uploader_participant_id: str = ""
    name: str = ""
    type: str = ""
    content_type: str = ""
    value: str = ""
    file_size: str = ""
    created_at: Optional[datetime] = None
    updated_at: Optional[datetime] = None


# =============================================================================
# Request/Options Data Classes
# =============================================================================


@dataclass
class ListLendingOffersOptions:
    """
    Options for listing lending offers.

    Attributes:
        currency_ids: Filter by currency IDs.
        participant_id: Filter by participant ID.
        duration: Filter by duration.
        sort_order: Sort order (ASC or DESC).
        page_size: Number of results per page.
        current_page: Current page token.
        page_request: Page navigation (FIRST, PREVIOUS, NEXT, LAST).
    """

    currency_ids: Optional[List[str]] = None
    participant_id: Optional[str] = None
    duration: Optional[str] = None
    sort_order: Optional[str] = None
    page_size: int = 50
    current_page: Optional[str] = None
    page_request: Optional[str] = None


@dataclass
class ListLendingAgreementsOptions:
    """
    Options for listing lending agreements.

    Attributes:
        sort_order: Sort order (ASC or DESC).
        page_size: Number of results per page.
        current_page: Current page token.
        page_request: Page navigation (FIRST, PREVIOUS, NEXT, LAST).
    """

    sort_order: Optional[str] = None
    page_size: int = 50
    current_page: Optional[str] = None
    page_request: Optional[str] = None


@dataclass
class CollateralRequest:
    """
    Collateral configuration for creating a lending agreement.

    Attributes:
        source_shared_address_id: Shared address holding the collateral.
        currency_id: Collateral currency ID.
        amount: Collateral amount.
    """

    source_shared_address_id: str
    currency_id: str
    amount: str


@dataclass
class CreateLendingOfferRequest:
    """
    Request to create a lending offer.

    Attributes:
        currency_id: Currency ID for the loan.
        amount: Amount to offer for lending.
        annual_percentage_yield: Interest rate (5 decimal precision).
        duration: Loan duration.
        collateral_requirements: Required collateral configuration.
    """

    currency_id: str
    amount: str
    annual_percentage_yield: str
    duration: str
    collateral_requirements: Optional[List[CurrencyCollateralRequirement]] = None


@dataclass
class CreateLendingAgreementRequest:
    """
    Request to create a lending agreement.

    Attributes:
        borrower_shared_address_id: Borrower's shared address for receiving funds.
        lending_offer_id: ID of the lending offer to accept.
        lender_participant_id: Lender participant ID.
        currency_id: Loan currency ID.
        amount: Loan amount.
        annual_percentage_yield: Interest rate (5 decimal precision).
        duration: Loan duration.
        collaterals: Collateral configuration.
    """

    borrower_shared_address_id: str
    lending_offer_id: Optional[str] = None
    lender_participant_id: Optional[str] = None
    currency_id: Optional[str] = None
    amount: Optional[str] = None
    annual_percentage_yield: Optional[str] = None
    duration: Optional[str] = None
    collaterals: Optional[List[CollateralRequest]] = None


@dataclass
class UpdateLendingAgreementRequest:
    """
    Request to update a lending agreement.

    Attributes:
        lender_shared_address_id: Updated lender shared address.
    """

    lender_shared_address_id: str


@dataclass
class RepayLendingAgreementRequest:
    """
    Request to repay a lending agreement.

    Attributes:
        repayer_shared_address_id: Shared address to repay from.
    """

    repayer_shared_address_id: str


@dataclass
class CreateLendingAgreementAttachmentRequest:
    """
    Request to create an attachment on a lending agreement.

    Attributes:
        name: Attachment name.
        value: Content (base64 for EMBEDDED, URL for EXTERNAL_LINK).
        content_type: MIME type (e.g., "application/pdf").
        type: Attachment type (EMBEDDED or EXTERNAL_LINK).
    """

    name: str
    value: str
    content_type: str
    type: str = "EMBEDDED"


@dataclass
class CursorPagination:
    """
    Cursor-based pagination information.

    Attributes:
        current_page: The current page cursor.
        has_next: Whether there is a next page.
        has_previous: Whether there is a previous page.
    """

    current_page: Optional[str] = None
    has_next: bool = False
    has_previous: bool = False


# =============================================================================
# Mapper Functions
# =============================================================================


def _currency_info_from_dto(dto: Any) -> Optional[CurrencyInfo]:
    """Convert OpenAPI currency DTO to CurrencyInfo."""
    if dto is None:
        return None

    return CurrencyInfo(
        id=getattr(dto, "id", "") or "",
        name=getattr(dto, "name", "") or "",
        symbol=getattr(dto, "symbol", "") or "",
        decimals=int(getattr(dto, "decimals", 0) or 0),
    )


def _collateral_requirement_from_dto(dto: Any) -> Optional[LendingCollateralRequirement]:
    """Convert OpenAPI collateral requirement DTO to domain model."""
    if dto is None:
        return None

    accepted_currencies = []
    for currency in getattr(dto, "accepted_currencies", None) or []:
        accepted_currencies.append(
            CurrencyCollateralRequirement(
                currency_id=getattr(currency, "currency_id", "") or "",
                percentage=getattr(currency, "percentage", "")
                or getattr(currency, "ratio", "")
                or "",
            )
        )

    return LendingCollateralRequirement(accepted_currencies=accepted_currencies)


def _lending_offer_from_dto(dto: Any) -> Optional[LendingOffer]:
    """Convert OpenAPI lending offer DTO to domain model."""
    if dto is None:
        return None

    return LendingOffer(
        id=getattr(dto, "id", "") or "",
        participant_id=getattr(dto, "participant_id", "") or "",
        currency_id=getattr(dto, "currency_id", "") or "",
        currency_info=_currency_info_from_dto(getattr(dto, "currency_info", None)),
        amount=getattr(dto, "amount", "") or "",
        amount_main_unit=getattr(dto, "amount_main_unit", "") or "",
        annual_percentage_yield=getattr(dto, "annual_percentage_yield", "") or "",
        annual_percentage_yield_main_unit=getattr(dto, "annual_percentage_yield_main_unit", "")
        or "",
        duration=getattr(dto, "duration", "") or "",
        blockchain=getattr(dto, "blockchain", "") or "",
        network=getattr(dto, "network", "") or "",
        collateral_requirement=_collateral_requirement_from_dto(
            getattr(dto, "collateral_requirement", None)
        ),
        created_at=getattr(dto, "created_at", None),
        updated_at=getattr(dto, "updated_at", None),
        origin_created_at=getattr(dto, "origin_created_at", None),
    )


def _lending_agreement_collateral_from_dto(dto: Any) -> Optional[LendingAgreementCollateral]:
    """Convert OpenAPI collateral DTO to domain model."""
    if dto is None:
        return None

    return LendingAgreementCollateral(
        id=getattr(dto, "id", "") or "",
        lending_agreement_id=getattr(dto, "lending_agreement_id", "") or "",
        lender_participant_id=getattr(dto, "lender_participant_id", "") or "",
        borrower_participant_id=getattr(dto, "borrower_participant_id", "") or "",
        currency_id=getattr(dto, "currency_id", "") or "",
        currency_info=_currency_info_from_dto(getattr(dto, "currency_info", None)),
        amount=getattr(dto, "amount", "") or "",
        amount_main_unit=getattr(dto, "amount_main_unit", "") or "",
        status=getattr(dto, "status", "") or "",
        pledge_id=getattr(dto, "pledge_id", "") or "",
        pledge_action_id=getattr(dto, "pledge_action_id", "") or "",
        shared_address_id=getattr(dto, "shared_address_id", "") or "",
        created_at=getattr(dto, "created_at", None),
        updated_at=getattr(dto, "updated_at", None),
    )


def _lending_agreement_transaction_from_dto(dto: Any) -> Optional[LendingAgreementTransaction]:
    """Convert OpenAPI transaction DTO to domain model."""
    if dto is None:
        return None

    return LendingAgreementTransaction(
        id=getattr(dto, "id", "") or "",
        type=getattr(dto, "type", "") or getattr(dto, "transaction_type", "") or "",
        amount=getattr(dto, "amount", "") or "",
        status=getattr(dto, "status", "") or "",
        created_at=getattr(dto, "created_at", None),
    )


def _lending_agreement_from_dto(dto: Any) -> Optional[LendingAgreement]:
    """Convert OpenAPI lending agreement DTO to domain model."""
    if dto is None:
        return None

    collaterals = []
    for collateral in getattr(dto, "lending_agreement_collaterals", None) or []:
        c = _lending_agreement_collateral_from_dto(collateral)
        if c:
            collaterals.append(c)

    transactions = []
    for tx in getattr(dto, "lending_agreement_transactions", None) or []:
        t = _lending_agreement_transaction_from_dto(tx)
        if t:
            transactions.append(t)

    return LendingAgreement(
        id=getattr(dto, "id", "") or "",
        lender_participant_id=getattr(dto, "lender_participant_id", "") or "",
        borrower_participant_id=getattr(dto, "borrower_participant_id", "") or "",
        lending_offer_id=getattr(dto, "lending_offer_id", "") or "",
        currency_id=getattr(dto, "currency_id", "") or "",
        currency_info=_currency_info_from_dto(getattr(dto, "currency_info", None)),
        amount=getattr(dto, "amount", "") or "",
        amount_main_unit=getattr(dto, "amount_main_unit", "") or "",
        annual_yield=getattr(dto, "annual_yield", "") or "",
        annual_yield_main_unit=getattr(dto, "annual_yield_main_unit", "") or "",
        duration=getattr(dto, "duration", "") or "",
        status=getattr(dto, "status", "") or "",
        workflow_id=getattr(dto, "workflow_id", "") or "",
        lender_shared_address_id=getattr(dto, "lender_shared_address_id", "") or "",
        borrower_shared_address_id=getattr(dto, "borrower_shared_address_id", "") or "",
        start_loan_date=getattr(dto, "start_loan_date", None),
        repayment_due_date=getattr(dto, "repayment_due_date", None),
        collaterals=collaterals,
        transactions=transactions,
        created_at=getattr(dto, "created_at", None),
        updated_at=getattr(dto, "updated_at", None),
    )


def _lending_agreement_attachment_from_dto(dto: Any) -> Optional[LendingAgreementAttachment]:
    """Convert OpenAPI attachment DTO to domain model."""
    if dto is None:
        return None

    return LendingAgreementAttachment(
        id=getattr(dto, "id", "") or "",
        lending_agreement_id=getattr(dto, "lending_agreement_id", "") or "",
        uploader_participant_id=getattr(dto, "uploader_participant_id", "") or "",
        name=getattr(dto, "name", "") or "",
        type=getattr(dto, "type", "") or "",
        content_type=getattr(dto, "content_type", "") or "",
        value=getattr(dto, "value", "") or "",
        file_size=getattr(dto, "file_size", "") or "",
        created_at=getattr(dto, "created_at", None),
        updated_at=getattr(dto, "updated_at", None),
    )


# =============================================================================
# Service Class
# =============================================================================


class LendingService(BaseService):
    """
    Service for Taurus Network lending operations.

    Provides methods to manage lending offers and agreements between
    Taurus Network participants.

    Example:
        >>> # List lending offers
        >>> offers, pagination = client.taurus_network.lending.list_lending_offers()
        >>> for offer in offers:
        ...     print(f"{offer.id}: {offer.annual_percentage_yield_main_unit} APY")
        >>>
        >>> # Get single agreement
        >>> agreement = client.taurus_network.lending.get_lending_agreement("123")
        >>> print(f"Status: {agreement.status}")
    """

    def __init__(self, api_client: Any, lending_api: Any) -> None:
        """
        Initialize lending service.

        Args:
            api_client: The OpenAPI client instance.
            lending_api: The TaurusNetworkLendingApi service.
        """
        super().__init__(api_client)
        self._lending_api = lending_api

    # =========================================================================
    # Lending Agreement Methods
    # =========================================================================

    def get_lending_agreement(self, lending_agreement_id: str) -> LendingAgreement:
        """
        Get a lending agreement by ID.

        Args:
            lending_agreement_id: The lending agreement ID to retrieve.

        Returns:
            The lending agreement.

        Raises:
            ValueError: If lending_agreement_id is empty.
            NotFoundError: If agreement not found.
            APIError: If API request fails.
        """
        self._validate_required(lending_agreement_id, "lending_agreement_id")

        try:
            resp = self._lending_api.taurus_network_service_get_lending_agreement(
                lending_agreement_id
            )

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Lending agreement {lending_agreement_id} not found")

            agreement = _lending_agreement_from_dto(result)
            if agreement is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Lending agreement {lending_agreement_id} not found")

            return agreement
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_lending_agreements(
        self,
        options: Optional[ListLendingAgreementsOptions] = None,
    ) -> Tuple[List[LendingAgreement], Optional[CursorPagination]]:
        """
        List lending agreements.

        Args:
            options: Optional filtering and pagination options.

        Returns:
            Tuple of (agreements list, cursor pagination info).

        Raises:
            APIError: If API request fails.
        """
        opts = options or ListLendingAgreementsOptions()

        try:
            resp = self._lending_api.taurus_network_service_get_lending_agreements(
                sort_order=opts.sort_order,
                cursor_current_page=opts.current_page,
                cursor_page_request=opts.page_request,
                cursor_page_size=str(opts.page_size) if opts.page_size > 0 else None,
            )

            result = getattr(resp, "result", None)
            agreements = []
            if result:
                for dto in result:
                    agreement = _lending_agreement_from_dto(dto)
                    if agreement:
                        agreements.append(agreement)

            # Extract cursor pagination
            cursor = getattr(resp, "cursor", None)
            pagination = None
            if cursor:
                pagination = CursorPagination(
                    current_page=getattr(cursor, "current_page", None),
                    has_next=getattr(cursor, "has_next", False) or False,
                    has_previous=getattr(cursor, "has_previous", False) or False,
                )

            return agreements, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def list_lending_agreements_for_approval(
        self,
        options: Optional[ListLendingAgreementsOptions] = None,
    ) -> Tuple[List[LendingAgreement], Optional[CursorPagination]]:
        """
        List lending agreements pending approval.

        Args:
            options: Optional filtering and pagination options.

        Returns:
            Tuple of (agreements list, cursor pagination info).

        Raises:
            APIError: If API request fails.
        """
        opts = options or ListLendingAgreementsOptions()

        try:
            resp = self._lending_api.taurus_network_service_get_lending_agreements_for_approval(
                sort_order=opts.sort_order,
                cursor_current_page=opts.current_page,
                cursor_page_request=opts.page_request,
                cursor_page_size=str(opts.page_size) if opts.page_size > 0 else None,
            )

            result = getattr(resp, "result", None)
            agreements = []
            if result:
                for dto in result:
                    agreement = _lending_agreement_from_dto(dto)
                    if agreement:
                        agreements.append(agreement)

            # Extract cursor pagination
            cursor = getattr(resp, "cursor", None)
            pagination = None
            if cursor:
                pagination = CursorPagination(
                    current_page=getattr(cursor, "current_page", None),
                    has_next=getattr(cursor, "has_next", False) or False,
                    has_previous=getattr(cursor, "has_previous", False) or False,
                )

            return agreements, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def create_lending_agreement(self, request: CreateLendingAgreementRequest) -> str:
        """
        Create a new lending agreement.

        Creates a lending agreement. Once approved by the lender and borrower,
        the agreement will automatically execute the loan lifecycle.
        This should be called by the borrower.

        Args:
            request: Lending agreement creation parameters.

        Returns:
            The created lending agreement ID.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.borrower_shared_address_id, "borrower_shared_address_id")

        try:
            # Build the request body
            body: Dict[str, Any] = {
                "borrowerSharedAddressID": request.borrower_shared_address_id,
            }

            if request.lending_offer_id:
                body["lendingOfferID"] = request.lending_offer_id
            if request.lender_participant_id:
                body["lenderParticipantID"] = request.lender_participant_id
            if request.currency_id:
                body["currencyID"] = request.currency_id
            if request.amount:
                body["amount"] = request.amount
            if request.annual_percentage_yield:
                body["annualPercentageYield"] = request.annual_percentage_yield
            if request.duration:
                body["duration"] = request.duration

            if request.collaterals:
                collaterals = [
                    {
                        "sourceSharedAddressID": c.source_shared_address_id,
                        "currencyID": c.currency_id,
                        "amount": c.amount,
                    }
                    for c in request.collaterals
                ]
                body["collaterals"] = collaterals

            resp = self._lending_api.taurus_network_service_create_lending_agreement(body=body)

            result = getattr(resp, "result", None)
            if result:
                return getattr(result, "id", "") or ""

            # Check for id directly on response
            return getattr(resp, "id", "") or ""
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def update_lending_agreement(
        self,
        lending_agreement_id: str,
        request: UpdateLendingAgreementRequest,
    ) -> None:
        """
        Update a lending agreement.

        Updates the lender's shared address for a lending agreement.
        This should be called by the lender before approving.

        Args:
            lending_agreement_id: The lending agreement ID to update.
            request: Update parameters.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        self._validate_required(lending_agreement_id, "lending_agreement_id")
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.lender_shared_address_id, "lender_shared_address_id")

        try:
            body = {
                "lenderSharedAddressID": request.lender_shared_address_id,
            }

            self._lending_api.taurus_network_service_update_lending_agreement(
                lending_agreement_id, body=body
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def repay_lending_agreement(
        self,
        lending_agreement_id: str,
        request: RepayLendingAgreementRequest,
    ) -> None:
        """
        Record repayment for a lending agreement.

        Records a repayment of the loan from the borrower to the lender.

        Args:
            lending_agreement_id: The lending agreement ID.
            request: Repayment parameters.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        self._validate_required(lending_agreement_id, "lending_agreement_id")
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.repayer_shared_address_id, "repayer_shared_address_id")

        try:
            body = {
                "repayerSharedAddressID": request.repayer_shared_address_id,
            }

            self._lending_api.taurus_network_service_repay_lending_agreement(
                lending_agreement_id, body=body
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def cancel_lending_agreement(self, lending_agreement_id: str) -> None:
        """
        Cancel a lending agreement.

        Cancels a lending agreement when it is not yet approved by the lender.

        Args:
            lending_agreement_id: The lending agreement ID to cancel.

        Raises:
            ValueError: If lending_agreement_id is empty.
            APIError: If API request fails.
        """
        self._validate_required(lending_agreement_id, "lending_agreement_id")

        try:
            self._lending_api.taurus_network_service_cancel_lending_agreement(
                lending_agreement_id, body={}
            )
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    # =========================================================================
    # Lending Agreement Attachment Methods
    # =========================================================================

    def create_lending_agreement_attachment(
        self,
        lending_agreement_id: str,
        request: CreateLendingAgreementAttachmentRequest,
    ) -> str:
        """
        Add an attachment to a lending agreement.

        Args:
            lending_agreement_id: The lending agreement ID.
            request: Attachment creation parameters.

        Returns:
            The created attachment ID.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        self._validate_required(lending_agreement_id, "lending_agreement_id")
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.name, "name")
        self._validate_required(request.value, "value")
        self._validate_required(request.content_type, "content_type")

        try:
            body = {
                "name": request.name,
                "value": request.value,
                "contentType": request.content_type,
                "type": request.type,
            }

            resp = self._lending_api.taurus_network_service_create_lending_agreement_attachment(
                lending_agreement_id, body=body
            )

            result = getattr(resp, "result", None)
            if result:
                return getattr(result, "id", "") or ""

            # Check for id directly on response
            return getattr(resp, "id", "") or ""
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_lending_agreement_attachments(
        self, lending_agreement_id: str
    ) -> List[LendingAgreementAttachment]:
        """
        List attachments for a lending agreement.

        Args:
            lending_agreement_id: The lending agreement ID.

        Returns:
            List of attachments.

        Raises:
            ValueError: If lending_agreement_id is empty.
            APIError: If API request fails.
        """
        self._validate_required(lending_agreement_id, "lending_agreement_id")

        try:
            resp = self._lending_api.taurus_network_service_get_lending_agreement_attachments(
                lending_agreement_id
            )

            result = getattr(resp, "result", None)
            attachments = []
            if result:
                for dto in result:
                    attachment = _lending_agreement_attachment_from_dto(dto)
                    if attachment:
                        attachments.append(attachment)

            return attachments
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    # =========================================================================
    # Lending Offer Methods
    # =========================================================================

    def get_lending_offer(self, offer_id: str) -> LendingOffer:
        """
        Get a lending offer by ID.

        Args:
            offer_id: The lending offer ID to retrieve.

        Returns:
            The lending offer.

        Raises:
            ValueError: If offer_id is empty.
            NotFoundError: If offer not found.
            APIError: If API request fails.
        """
        self._validate_required(offer_id, "offer_id")

        try:
            resp = self._lending_api.taurus_network_service_get_lending_offer(offer_id)

            # Offer is directly in lending_offer field
            result = getattr(resp, "lending_offer", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Lending offer {offer_id} not found")

            offer = _lending_offer_from_dto(result)
            if offer is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Lending offer {offer_id} not found")

            return offer
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_lending_offers(
        self,
        options: Optional[ListLendingOffersOptions] = None,
    ) -> Tuple[List[LendingOffer], Optional[CursorPagination]]:
        """
        List lending offers.

        Args:
            options: Optional filtering and pagination options.

        Returns:
            Tuple of (offers list, cursor pagination info).

        Raises:
            APIError: If API request fails.
        """
        opts = options or ListLendingOffersOptions()

        try:
            resp = self._lending_api.taurus_network_service_get_lending_offers(
                sort_order=opts.sort_order,
                cursor_current_page=opts.current_page,
                cursor_page_request=opts.page_request,
                cursor_page_size=str(opts.page_size) if opts.page_size > 0 else None,
                currency_ids_currency_ids=opts.currency_ids,
                participant_id=opts.participant_id,
                duration=opts.duration,
            )

            result = getattr(resp, "result", None)
            offers = []
            if result:
                for dto in result:
                    offer = _lending_offer_from_dto(dto)
                    if offer:
                        offers.append(offer)

            # Extract cursor pagination
            cursor = getattr(resp, "cursor", None)
            pagination = None
            if cursor:
                pagination = CursorPagination(
                    current_page=getattr(cursor, "current_page", None),
                    has_next=getattr(cursor, "has_next", False) or False,
                    has_previous=getattr(cursor, "has_previous", False) or False,
                )

            return offers, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def create_lending_offer(self, request: CreateLendingOfferRequest) -> str:
        """
        Create a new lending offer.

        Creates a lending offer that other participants can accept
        to create a lending agreement.

        Args:
            request: Lending offer creation parameters.

        Returns:
            The created lending offer ID.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        if request is None:
            raise ValueError("request cannot be None")
        self._validate_required(request.currency_id, "currency_id")
        self._validate_required(request.amount, "amount")
        self._validate_required(request.annual_percentage_yield, "annual_percentage_yield")
        self._validate_required(request.duration, "duration")

        try:
            body: Dict[str, Any] = {
                "currencyID": request.currency_id,
                "amount": request.amount,
                "annualPercentageYield": request.annual_percentage_yield,
                "duration": request.duration,
            }

            if request.collateral_requirements:
                collateral_requirements = [
                    {
                        "currencyID": c.currency_id,
                        "percentage": c.percentage,
                    }
                    for c in request.collateral_requirements
                ]
                body["collateralRequirement"] = collateral_requirements

            resp = self._lending_api.taurus_network_service_create_lending_offer(body=body)

            result = getattr(resp, "result", None)
            if result:
                return getattr(result, "id", "") or ""

            # Check for id directly on response
            return getattr(resp, "id", "") or ""
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def delete_lending_offer(self, offer_id: str) -> None:
        """
        Delete a specific lending offer.

        Args:
            offer_id: The lending offer ID to delete.

        Raises:
            ValueError: If offer_id is empty.
            APIError: If API request fails.
        """
        self._validate_required(offer_id, "offer_id")

        try:
            self._lending_api.taurus_network_service_delete_lending_offer(offer_id)
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def delete_lending_offers(self) -> None:
        """
        Delete all lending offers for the current participant.

        Raises:
            APIError: If API request fails.
        """
        try:
            self._lending_api.taurus_network_service_delete_lending_offers()
        except Exception as e:
            from taurus_protect.errors import APIError

            if type(e).__name__ == "ApiException":
                raise self._handle_error(e) from e
            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    # =========================================================================
    # Convenience Alias Methods
    # =========================================================================

    def create_attachment(
        self,
        agreement_id: str,
        request: CreateLendingAgreementAttachmentRequest,
    ) -> str:
        """
        Add an attachment to a lending agreement.

        Convenience alias for create_lending_agreement_attachment().

        Args:
            agreement_id: The lending agreement ID.
            request: Attachment creation parameters.

        Returns:
            The created attachment ID.

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        return self.create_lending_agreement_attachment(agreement_id, request)

    def list_attachments(self, agreement_id: str) -> List[LendingAgreementAttachment]:
        """
        List attachments for a lending agreement.

        Convenience alias for list_lending_agreement_attachments().

        Args:
            agreement_id: The lending agreement ID.

        Returns:
            List of attachments.

        Raises:
            ValueError: If agreement_id is empty.
            APIError: If API request fails.
        """
        return self.list_lending_agreement_attachments(agreement_id)
