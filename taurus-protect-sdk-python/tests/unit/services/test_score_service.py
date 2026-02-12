"""Unit tests for ScoreService."""

from __future__ import annotations

from unittest.mock import MagicMock

import pytest

from taurus_protect.services.score_service import ScoreService


class TestGetAddressScore:
    """Tests for ScoreService.get_address_score()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        scores_api = MagicMock()
        service = ScoreService(api_client=api_client, scores_api=scores_api)
        return service, scores_api

    def test_raises_on_empty_address_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="address_id"):
            service.get_address_score(address_id="")

    def test_returns_empty_when_no_scores(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.scores = None
        api.score_service_refresh_address_score.return_value = resp

        result = service.get_address_score(address_id="123")

        assert result == []

    def test_passes_provider_in_body(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.scores = None
        api.score_service_refresh_address_score.return_value = resp

        service.get_address_score(address_id="123", provider="chainalysis")

        api.score_service_refresh_address_score.assert_called_once_with(
            address_id="123",
            body={"provider": "chainalysis"},
        )


class TestGetTransactionScore:
    """Tests for ScoreService.get_transaction_score()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        scores_api = MagicMock()
        service = ScoreService(api_client=api_client, scores_api=scores_api)
        return service, scores_api

    def test_raises_on_empty_tx_hash(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="tx_hash"):
            service.get_transaction_score(tx_hash="")

    def test_raises_not_implemented(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(NotImplementedError):
            service.get_transaction_score(tx_hash="0xabc")


class TestRefreshWhitelistedAddressScore:
    """Tests for ScoreService.refresh_whitelisted_address_score()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        scores_api = MagicMock()
        service = ScoreService(api_client=api_client, scores_api=scores_api)
        return service, scores_api

    def test_raises_on_empty_address_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="address_id"):
            service.refresh_whitelisted_address_score(address_id="")

    def test_returns_empty_when_no_scores(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.scores = None
        api.score_service_refresh_wla_score.return_value = resp

        result = service.refresh_whitelisted_address_score(address_id="456")

        assert result == []
