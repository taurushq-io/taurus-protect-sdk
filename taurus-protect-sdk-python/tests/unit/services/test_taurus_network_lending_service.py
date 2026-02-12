"""Unit tests for TaurusNetwork LendingService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.taurus_network.lending_service import LendingService


class TestGetLendingAgreement:
    """Tests for LendingService.get_lending_agreement()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        lending_api = MagicMock()
        service = LendingService(api_client=api_client, lending_api=lending_api)
        return service, lending_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="lending_agreement_id"):
            service.get_lending_agreement(lending_agreement_id="")

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.taurus_network_service_get_lending_agreement.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get_lending_agreement(lending_agreement_id="la-missing")


class TestListLendingAgreements:
    """Tests for LendingService.list_lending_agreements()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        lending_api = MagicMock()
        service = LendingService(api_client=api_client, lending_api=lending_api)
        return service, lending_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.taurus_network_service_get_lending_agreements.return_value = resp

        agreements, pagination = service.list_lending_agreements()

        assert agreements == []


class TestGetLendingOffer:
    """Tests for LendingService.get_lending_offer()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        lending_api = MagicMock()
        service = LendingService(api_client=api_client, lending_api=lending_api)
        return service, lending_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="offer_id"):
            service.get_lending_offer(offer_id="")

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.lending_offer = None
        api.taurus_network_service_get_lending_offer.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get_lending_offer(offer_id="lo-missing")


class TestListLendingOffers:
    """Tests for LendingService.list_lending_offers()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        lending_api = MagicMock()
        service = LendingService(api_client=api_client, lending_api=lending_api)
        return service, lending_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.taurus_network_service_get_lending_offers.return_value = resp

        offers, pagination = service.list_lending_offers()

        assert offers == []


class TestCancelLendingAgreement:
    """Tests for LendingService.cancel_lending_agreement()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        lending_api = MagicMock()
        service = LendingService(api_client=api_client, lending_api=lending_api)
        return service, lending_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="lending_agreement_id"):
            service.cancel_lending_agreement(lending_agreement_id="")
