"""Integration tests for extended services.

Tests read-only operations (list) for services that did not previously have
integration test coverage:
  - BusinessRuleService
  - WebhookService
  - WebhookCallService
  - StakingService (list_validators)
  - FeePayerService
  - ExchangeService
  - AssetService
  - FiatService
  - JobService
  - ContractWhitelistingService
  - TaurusNetwork (participants, pledges, lending, settlements, sharing)

These tests require a live API connection. Enable by:
    export PROTECT_INTEGRATION_TEST=true
    export PROTECT_API_HOST="https://your-api-host.com"
    export PROTECT_API_KEY="your-api-key"
    export PROTECT_API_SECRET="your-hex-encoded-secret"
    pytest tests/integration/test_extended_services.py -v

Or configure defaults in conftest.py.
"""

from __future__ import annotations

import logging
from typing import TYPE_CHECKING

import pytest

if TYPE_CHECKING:
    from taurus_protect.client import ProtectClient

logger = logging.getLogger(__name__)


# =============================================================================
# BusinessRuleService
# =============================================================================


@pytest.mark.integration
def test_list_business_rules(client: ProtectClient) -> None:
    """Test listing business rules."""
    result = client.business_rules.list()

    logger.info("Found %d business rules", len(result.rules))
    assert result.rules is not None

    logger.info("Pagination: current_page=%s, has_next=%s", result.current_page, result.has_next)

    for rule in result.rules[:5]:
        logger.info("  Rule: ID=%s", rule.id)


# =============================================================================
# WebhookService
# =============================================================================


@pytest.mark.integration
def test_list_webhooks(client: ProtectClient) -> None:
    """Test listing webhooks."""
    webhooks, pagination = client.webhooks.list()

    logger.info("Found %d webhooks", len(webhooks))
    assert webhooks is not None

    if pagination is not None:
        logger.info("Pagination: total=%s, has_more=%s", pagination.total_items, pagination.has_more)

    for wh in webhooks[:5]:
        logger.info("  Webhook: ID=%s, URL=%s", wh.id, getattr(wh, "url", None))


# =============================================================================
# WebhookCallService
# =============================================================================


@pytest.mark.integration
def test_list_webhook_calls(client: ProtectClient) -> None:
    """Test listing webhook calls."""
    calls, pagination = client.webhook_calls.list()

    logger.info("Found %d webhook calls", len(calls))
    assert calls is not None

    if pagination is not None:
        logger.info("Pagination: total=%s, has_more=%s", pagination.total_items, pagination.has_more)

    for call in calls[:5]:
        logger.info("  WebhookCall: ID=%s, Status=%s", call.id, getattr(call, "status", None))


# =============================================================================
# StakingService
# =============================================================================


@pytest.mark.integration
def test_list_staking_validators(client: ProtectClient) -> None:
    """Test listing staking validators for ETH."""
    validators, pagination = client.staking.list_validators(
        blockchain="ETH",
        network="mainnet",
        limit=10,
    )

    logger.info("Found %d ETH validators", len(validators))
    assert validators is not None

    for v in validators[:5]:
        logger.info("  Validator: %s", getattr(v, "name", v))


# =============================================================================
# FeePayerService
# =============================================================================


@pytest.mark.integration
def test_list_fee_payers(client: ProtectClient) -> None:
    """Test listing fee payers."""
    fee_payers, pagination = client.fee_payers.list()

    logger.info("Found %d fee payers", len(fee_payers))
    assert fee_payers is not None

    for fp in fee_payers[:5]:
        logger.info("  FeePayer: ID=%s", fp.id)


# =============================================================================
# ExchangeService
# =============================================================================


@pytest.mark.integration
def test_list_exchanges(client: ProtectClient) -> None:
    """Test listing exchange counterparties."""
    counterparties = client.exchanges.list_counterparties()

    logger.info("Found %d exchange counterparties", len(counterparties))
    assert counterparties is not None

    for cp in counterparties[:5]:
        logger.info("  Exchange: %s", getattr(cp, "name", cp))


# =============================================================================
# AssetService
# =============================================================================


@pytest.mark.integration
def test_list_assets(client: ProtectClient) -> None:
    """Test listing assets (via asset wallets endpoint)."""
    assets, pagination = client.assets.list(currency="ETH", limit=10)

    logger.info("Found %d assets", len(assets))
    assert assets is not None

    for asset in assets[:5]:
        logger.info("  Asset: %s", getattr(asset, "name", asset))


# =============================================================================
# FiatService
# =============================================================================


@pytest.mark.integration
def test_list_fiat_providers(client: ProtectClient) -> None:
    """Test listing fiat providers."""
    providers = client.fiat.list_providers()

    logger.info("Found %d fiat providers", len(providers))
    assert providers is not None

    for provider in providers[:5]:
        logger.info("  Provider: %s", provider)


# =============================================================================
# JobService
# =============================================================================


@pytest.mark.integration
def test_list_jobs(client: ProtectClient) -> None:
    """Test listing jobs. Requires tgvalidatord role."""
    from taurus_protect.errors import AuthorizationError

    try:
        jobs, pagination = client.jobs.list()

        logger.info("Found %d jobs", len(jobs))
        assert jobs is not None

        for job in jobs[:5]:
            logger.info("  Job: ID=%s, Name=%s", getattr(job, "id", None), getattr(job, "name", None))
    except AuthorizationError:
        # Jobs endpoint requires 'tgvalidatord' role which may not be
        # available with the current test credentials (matches Java SDK behavior)
        logger.info("Jobs endpoint requires tgvalidatord role (403 Forbidden)")


# =============================================================================
# TaurusNetwork - Participants
# =============================================================================


@pytest.mark.integration
def test_list_taurus_network_participants(client: ProtectClient) -> None:
    """Test listing TaurusNetwork participants."""
    participants = client.taurus_network.participants.list()

    logger.info("Found %d participants", len(participants))
    assert participants is not None

    for p in participants[:5]:
        logger.info("  Participant: ID=%s, Name=%s", p.id, getattr(p, "name", None))


# =============================================================================
# TaurusNetwork - Pledges
# =============================================================================


@pytest.mark.integration
def test_list_taurus_network_pledges(client: ProtectClient) -> None:
    """Test listing TaurusNetwork pledges."""
    pledges, pagination = client.taurus_network.pledges.list_pledges()

    logger.info("Found %d pledges", len(pledges))
    assert pledges is not None

    for pledge in pledges[:5]:
        logger.info("  Pledge: ID=%s, Status=%s", pledge.id, getattr(pledge, "status", None))


# =============================================================================
# TaurusNetwork - Lending
# =============================================================================


@pytest.mark.integration
def test_list_taurus_network_lending_offers(client: ProtectClient) -> None:
    """Test listing TaurusNetwork lending offers."""
    offers, pagination = client.taurus_network.lending.list_lending_offers()

    logger.info("Found %d lending offers", len(offers))
    assert offers is not None

    for offer in offers[:5]:
        logger.info("  Offer: ID=%s", offer.id)


@pytest.mark.integration
def test_list_taurus_network_lending_agreements(client: ProtectClient) -> None:
    """Test listing TaurusNetwork lending agreements."""
    agreements, pagination = client.taurus_network.lending.list_lending_agreements()

    logger.info("Found %d lending agreements", len(agreements))
    assert agreements is not None

    for agreement in agreements[:5]:
        logger.info("  Agreement: ID=%s", agreement.id)


# =============================================================================
# TaurusNetwork - Settlements
# =============================================================================


@pytest.mark.integration
def test_list_taurus_network_settlements(client: ProtectClient) -> None:
    """Test listing TaurusNetwork settlements."""
    settlements, pagination = client.taurus_network.settlements.list_settlements()

    logger.info("Found %d settlements", len(settlements))
    assert settlements is not None

    for s in settlements[:5]:
        logger.info("  Settlement: ID=%s, Status=%s", s.id, getattr(s, "status", None))


# =============================================================================
# TaurusNetwork - Sharing
# =============================================================================


@pytest.mark.integration
def test_list_taurus_network_shared_addresses(client: ProtectClient) -> None:
    """Test listing TaurusNetwork shared addresses."""
    addresses, pagination = client.taurus_network.sharing.list_shared_addresses()

    logger.info("Found %d shared addresses", len(addresses))
    assert addresses is not None

    for addr in addresses[:5]:
        logger.info("  SharedAddress: ID=%s", addr.id)


@pytest.mark.integration
def test_list_taurus_network_shared_assets(client: ProtectClient) -> None:
    """Test listing TaurusNetwork shared assets."""
    assets, pagination = client.taurus_network.sharing.list_shared_assets()

    logger.info("Found %d shared assets", len(assets))
    assert assets is not None

    for asset in assets[:5]:
        logger.info("  SharedAsset: ID=%s", asset.id)


# =============================================================================
# ContractWhitelistingService
# =============================================================================


@pytest.mark.integration
def test_list_whitelisted_contracts(client: ProtectClient) -> None:
    """Test listing whitelisted contracts."""
    contracts, pagination = client.contract_whitelisting.list()

    logger.info("Found %d whitelisted contracts", len(contracts))
    assert contracts is not None

    for c in contracts[:5]:
        logger.info("  Contract: ID=%s", getattr(c, "id", c))
