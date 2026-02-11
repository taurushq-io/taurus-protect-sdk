"""Unit tests for GovernanceRuleService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.errors import IntegrityError
from taurus_protect.services.governance_rule_service import GovernanceRuleService


class TestGetRules:
    """Tests for GovernanceRuleService.get_rules()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        governance_rules_api = MagicMock()
        super_admin_keys = [MagicMock()]
        service = GovernanceRuleService(
            api_client=api_client,
            governance_rules_api=governance_rules_api,
            super_admin_keys=super_admin_keys,
            min_valid_signatures=1,
        )
        return service, governance_rules_api

    def test_get_rules_returns_none_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.rule_service_get_rules.return_value = reply

        result = service.get_rules()
        assert result is None

    def test_get_rules_maps_and_verifies(self) -> None:
        service, api = self._make_service()

        dto = MagicMock()
        dto.rules_container = "base64data"
        dto.rules_signatures = []
        dto.locked = False
        dto.creation_date = None
        dto.update_date = None
        reply = MagicMock()
        reply.result = dto
        api.rule_service_get_rules.return_value = reply

        with patch(
            "taurus_protect.services.governance_rule_service.verify_governance_rules"
        ):
            result = service.get_rules()

        assert result is not None
        assert result.rules_container == "base64data"


class TestGetRulesById:
    """Tests for GovernanceRuleService.get_rules_by_id()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        governance_rules_api = MagicMock()
        service = GovernanceRuleService(
            api_client=api_client,
            governance_rules_api=governance_rules_api,
            super_admin_keys=[MagicMock()],
            min_valid_signatures=1,
        )
        return service, governance_rules_api

    def test_get_rules_by_id_raises_for_empty_id(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="rules_id"):
            service.get_rules_by_id("")

    def test_get_rules_by_id_returns_none_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.rule_service_get_rules_by_id.return_value = reply

        result = service.get_rules_by_id("123")
        assert result is None


class TestGetRulesProposal:
    """Tests for GovernanceRuleService.get_rules_proposal()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        governance_rules_api = MagicMock()
        service = GovernanceRuleService(
            api_client=api_client,
            governance_rules_api=governance_rules_api,
            super_admin_keys=[MagicMock()],
            min_valid_signatures=1,
        )
        return service, governance_rules_api

    def test_get_rules_proposal_returns_none_when_no_result(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.result = None
        api.rule_service_get_rules_proposal.return_value = reply

        result = service.get_rules_proposal()
        assert result is None

    def test_get_rules_proposal_does_not_verify(self) -> None:
        service, api = self._make_service()

        dto = MagicMock()
        dto.rules_container = "base64data"
        dto.rules_signatures = []
        dto.locked = False
        dto.creation_date = None
        dto.update_date = None
        reply = MagicMock()
        reply.result = dto
        api.rule_service_get_rules_proposal.return_value = reply

        # Should not call verify_governance_rules for proposals
        with patch(
            "taurus_protect.services.governance_rule_service.verify_governance_rules"
        ) as mock_verify:
            service.get_rules_proposal()
            mock_verify.assert_not_called()


class TestGetPublicKeys:
    """Tests for GovernanceRuleService.get_public_keys()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        governance_rules_api = MagicMock()
        service = GovernanceRuleService(
            api_client=api_client,
            governance_rules_api=governance_rules_api,
            super_admin_keys=[MagicMock()],
            min_valid_signatures=1,
        )
        return service, governance_rules_api

    def test_get_public_keys_returns_empty_list(self) -> None:
        service, api = self._make_service()

        reply = MagicMock()
        reply.public_keys = None
        api.rule_service_get_public_keys.return_value = reply

        result = service.get_public_keys()
        assert result == []

    def test_get_public_keys_maps_dtos(self) -> None:
        service, api = self._make_service()

        pk_dto = MagicMock()
        pk_dto.user_id = "user-1"
        pk_dto.public_key = "pem-key"
        reply = MagicMock()
        reply.public_keys = [pk_dto]
        api.rule_service_get_public_keys.return_value = reply

        result = service.get_public_keys()
        assert len(result) == 1
        assert result[0].user_id == "user-1"
        assert result[0].public_key == "pem-key"


class TestVerifyGovernanceRules:
    """Tests for GovernanceRuleService.verify_governance_rules()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        governance_rules_api = MagicMock()
        service = GovernanceRuleService(
            api_client=api_client,
            governance_rules_api=governance_rules_api,
            super_admin_keys=[MagicMock()],
            min_valid_signatures=1,
        )
        return service, governance_rules_api

    def test_verify_returns_rules_on_success(self) -> None:
        service, _ = self._make_service()

        mock_rules = MagicMock()
        with patch(
            "taurus_protect.services.governance_rule_service.verify_governance_rules"
        ):
            result = service.verify_governance_rules(mock_rules)

        assert result is mock_rules

    def test_verify_raises_integrity_error_on_failure(self) -> None:
        service, _ = self._make_service()

        mock_rules = MagicMock()
        with patch(
            "taurus_protect.services.governance_rule_service.verify_governance_rules",
            side_effect=IntegrityError("insufficient signatures"),
        ):
            with pytest.raises(IntegrityError, match="insufficient signatures"):
                service.verify_governance_rules(mock_rules)


class TestProperties:
    """Tests for GovernanceRuleService properties."""

    def test_super_admin_keys_returns_copy(self) -> None:
        keys = [MagicMock(), MagicMock()]
        service = GovernanceRuleService(
            api_client=MagicMock(),
            governance_rules_api=MagicMock(),
            super_admin_keys=keys,
            min_valid_signatures=2,
        )

        result = service.super_admin_keys
        assert len(result) == 2
        # Must be a copy
        assert result is not keys

    def test_min_valid_signatures(self) -> None:
        service = GovernanceRuleService(
            api_client=MagicMock(),
            governance_rules_api=MagicMock(),
            super_admin_keys=[MagicMock()],
            min_valid_signatures=3,
        )

        assert service.min_valid_signatures == 3


class TestGetDecodedRulesContainer:
    """Tests for GovernanceRuleService.get_decoded_rules_container()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        governance_rules_api = MagicMock()
        service = GovernanceRuleService(
            api_client=api_client,
            governance_rules_api=governance_rules_api,
            super_admin_keys=[MagicMock()],
            min_valid_signatures=1,
        )
        return service, governance_rules_api

    def test_raises_for_none_rules(self) -> None:
        service, _ = self._make_service()

        with pytest.raises(ValueError, match="rules cannot be None"):
            service.get_decoded_rules_container(None)

    def test_raises_integrity_error_for_none_container(self) -> None:
        service, _ = self._make_service()

        mock_rules = MagicMock()
        mock_rules.rules_container = None

        with patch(
            "taurus_protect.services.governance_rule_service.verify_governance_rules"
        ):
            with pytest.raises(IntegrityError, match="Rules container is None"):
                service.get_decoded_rules_container(mock_rules)
