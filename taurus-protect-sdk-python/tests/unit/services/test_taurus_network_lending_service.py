"""Unit tests for TaurusNetwork LendingService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect._internal.openapi import TaurusNetworkLendingApi
from taurus_protect.models.pagination import CursorPage
from taurus_protect.models.taurus_network.lending import (
    ListLendingAgreementsOptions,
    ListLendingOffersOptions,
)
from taurus_protect.services.taurus_network.lending_service import LendingService
from tests.unit.transport_stub import StubTransport, api_client


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
    """list_lending_agreements reads ``lendingAgreements`` (it read ``result``, always empty)."""

    def _service(self) -> LendingService:
        ac = api_client()
        return LendingService(ac, TaurusNetworkLendingApi(ac))

    def test_rows_page_and_sort(self) -> None:
        reply = {
            "lendingAgreements": [{"id": "la1"}],
            "cursor": {"currentPage": "n", "hasNext": True},
        }
        with StubTransport(reply) as transport:
            agreements, page = self._service().list_lending_agreements(
                ListLendingAgreementsOptions(sort_order="ASC", page_size=5)
            )

        assert [a.id for a in agreements] == ["la1"]
        assert page == CursorPage(page_size=5, next_cursor="n", has_more=True)
        assert transport.last.query == [("cursor.pageSize", "5"), ("sortOrder", "ASC")]

    def test_ids_only_filter_the_approval_queue(self) -> None:
        with StubTransport() as transport:
            with pytest.raises(ValueError, match="ids"):
                self._service().list_lending_agreements(ListLendingAgreementsOptions(ids=["1"]))
            self._service().list_lending_agreements_for_approval(
                ListLendingAgreementsOptions(ids=["1"], cursor="n")
            )

        assert len(transport.requests) == 1
        assert transport.last.path == "/api/rest/v1/tn/lending/agreements/for-approval"
        assert transport.last.query == sorted(
            [
                ("ids", "1"),
                ("cursor.currentPage", "n"),
                ("cursor.pageRequest", "NEXT"),
                ("cursor.pageSize", "20"),
            ]
        )


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
    """list_lending_offers reads ``lendingOffers`` (it read ``result``, always empty)."""

    def _service(self) -> LendingService:
        ac = api_client()
        return LendingService(ac, TaurusNetworkLendingApi(ac))

    def test_rows_page_and_filters(self) -> None:
        reply = {"lendingOffers": [{"id": "lo1"}], "cursor": {"currentPage": "n"}}
        with StubTransport(reply) as transport:
            offers, page = self._service().list_lending_offers(
                ListLendingOffersOptions(
                    currency_ids=["c1", "c2"], participant_id="p", duration="30d"
                )
            )

        assert [o.id for o in offers] == ["lo1"]
        assert page == CursorPage(page_size=20)
        assert transport.last.query == sorted(
            [
                ("currencyIDs.currencyIDs", "c1"),
                ("currencyIDs.currencyIDs", "c2"),
                ("participantID", "p"),
                ("duration", "30d"),
                ("cursor.pageSize", "20"),
            ]
        )


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
