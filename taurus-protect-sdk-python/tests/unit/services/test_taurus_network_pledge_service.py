"""Unit tests for TaurusNetwork PledgeService."""

from __future__ import annotations

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.services.taurus_network.pledge_service import PledgeService


class TestGetPledge:
    """Tests for PledgeService.get_pledge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="pledge_id"):
            service.get_pledge(pledge_id="")

    def test_raises_not_found_when_none(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        api.taurus_network_service_get_pledge.return_value = resp

        from taurus_protect.errors import NotFoundError

        with pytest.raises(NotFoundError):
            service.get_pledge(pledge_id="p-missing")


class TestListPledges:
    """Tests for PledgeService.list_pledges()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.taurus_network_service_get_pledges.return_value = resp

        pledges, pagination = service.list_pledges()

        assert pledges == []


class TestCreatePledge:
    """Tests for PledgeService.create_pledge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.create_pledge(req=None)

    def test_raises_on_missing_shared_address_id(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import CreatePledgeRequest

        req = CreatePledgeRequest(
            shared_address_id="",
            currency_id="ETH",
            amount="1000",
        )
        with pytest.raises(ValueError, match="shared_address_id"):
            service.create_pledge(req=req)


class TestApprovePledgeActions:
    """Tests for PledgeService.approve_pledge_actions()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_empty_actions(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="actions list cannot be empty"):
            service.approve_pledge_actions(actions=[], private_key=MagicMock())

    def test_raises_on_none_private_key(self) -> None:
        service, _ = self._make_service()
        action = MagicMock()
        action.id = "a-1"
        action.metadata = MagicMock()
        action.metadata.hash = "abc123"
        with pytest.raises(ValueError, match="private_key cannot be None"):
            service.approve_pledge_actions(actions=[action], private_key=None)

    def test_raises_on_missing_metadata(self) -> None:
        service, _ = self._make_service()
        action = MagicMock()
        action.id = "a-1"
        action.metadata = None
        with pytest.raises(ValueError, match="action metadata cannot be None"):
            service.approve_pledge_actions(
                actions=[action], private_key=MagicMock()
            )

    def test_raises_on_empty_hash(self) -> None:
        service, _ = self._make_service()
        action = MagicMock()
        action.id = "a-1"
        action.metadata = MagicMock()
        action.metadata.hash = ""
        with pytest.raises(ValueError, match="action metadata hash cannot be empty"):
            service.approve_pledge_actions(
                actions=[action], private_key=MagicMock()
            )


class TestApprovePledgeActionsVerification:
    """
    approve_pledge_actions was the last signing site in any of the four SDKs with no
    verification: it checked that a hash was non-empty and then signed it, so the
    approver attested to a hash nothing had checked. Python is the only SDK where the
    SDK signs pledge actions at all (Go and TypeScript take a caller-supplied signature,
    Java has no such method), which is why it drifted alone.
    """

    def _make_service(self) -> tuple:
        pledge_api = MagicMock()
        reply = MagicMock()
        reply.approved_count = 1
        pledge_api.taurus_network_service_approve_pledge_actions.return_value = reply
        service = PledgeService(api_client=MagicMock(), pledge_api=pledge_api)
        return service, pledge_api

    @staticmethod
    def _action(action_id: str, payload: str, hash_value: str) -> MagicMock:
        action = MagicMock()
        action.id = action_id
        action.metadata = MagicMock()
        action.metadata.payload = payload
        action.metadata.hash = hash_value
        return action

    def test_signs_an_action_whose_hash_covers_its_payload(self) -> None:
        from taurus_protect.crypto.hashing import calculate_hex_hash

        service, api = self._make_service()
        payload = '{"amount":"1"}'
        action = self._action("a-1", payload, calculate_hex_hash(payload))

        with patch(
            "taurus_protect.services.taurus_network.pledge_service.sign_data",
            return_value="signature_base64",
        ):
            assert service.approve_pledge_actions(actions=[action], private_key=MagicMock()) == 1
        api.taurus_network_service_approve_pledge_actions.assert_called_once()

    def test_refuses_a_hash_that_does_not_cover_its_payload(self) -> None:
        from taurus_protect.crypto.hashing import calculate_hex_hash
        from taurus_protect.errors import IntegrityError

        service, api = self._make_service()
        # The documented attack: alter the payload, leave the hash alone.
        action = self._action("a-1", '{"amount":"999"}', calculate_hex_hash('{"amount":"1"}'))

        with pytest.raises(IntegrityError, match="refusing to sign pledge action a-1"):
            service.approve_pledge_actions(actions=[action], private_key=MagicMock())
        api.taurus_network_service_approve_pledge_actions.assert_not_called()

    def test_refuses_a_hash_with_no_payload(self) -> None:
        from taurus_protect.errors import IntegrityError

        service, api = self._make_service()
        action = self._action("a-1", "", "abc123")

        with pytest.raises(IntegrityError, match="payload is missing"):
            service.approve_pledge_actions(actions=[action], private_key=MagicMock())
        api.taurus_network_service_approve_pledge_actions.assert_not_called()

    def test_is_all_or_nothing(self) -> None:
        """One signature covers every hash in the batch, so partial is not an option."""
        from taurus_protect.crypto.hashing import calculate_hex_hash
        from taurus_protect.errors import IntegrityError

        service, api = self._make_service()
        good_payload = '{"amount":"1"}'
        ok = self._action("a-1", good_payload, calculate_hex_hash(good_payload))
        bad = self._action("a-2", '{"amount":"999"}', calculate_hex_hash(good_payload))

        with pytest.raises(IntegrityError):
            service.approve_pledge_actions(actions=[ok, bad], private_key=MagicMock())
        api.taurus_network_service_approve_pledge_actions.assert_not_called()


    def test_signs_compact_json_matching_the_other_sdks(self) -> None:
        """The signed bytes are compact JSON, as in Go/Java/TS and every whitelist path.

        ``json.dumps`` defaults to ``", "`` between elements, so this SDK used to sign a
        DIFFERENT byte string for any batch of two or more -- the server verifies the
        signature against its own rebuilt array, so a multi-action approval signed here
        was rejected there.
        """
        from taurus_protect.crypto.hashing import calculate_hex_hash

        service, _ = self._make_service()
        p1, p2 = '{"amount":"1"}', '{"amount":"2"}'
        h1, h2 = calculate_hex_hash(p1), calculate_hex_hash(p2)
        actions = [self._action("a-1", p1, h1), self._action("a-2", p2, h2)]

        captured: dict = {}

        def capture(_key: object, data: bytes) -> str:
            captured["signed"] = data.decode("utf-8")
            return "signature_base64"

        with patch(
            "taurus_protect.services.taurus_network.pledge_service.sign_data",
            side_effect=capture,
        ):
            service.approve_pledge_actions(actions=actions, private_key=MagicMock())

        assert captured["signed"] == f'["{h1}","{h2}"]', "no spaces: compact separators"


class TestPledgeActionReadVerification:
    """The read paths verify each action's hash against its payload.

    Approval verifies again in the same call -- that is deliberate, not redundant: an
    action can reach ``approve_pledge_actions`` decoded from a queue or cache rather than
    from this SDK. What the read-path check buys is that the payload an integrator
    DISPLAYS for review is one the hash commits to.
    """

    def _make_service(self, rows: list) -> tuple:
        pledge_api = MagicMock()
        reply = MagicMock()
        reply.result = rows
        reply.total_items = str(len(rows))
        reply.offset = "0"
        pledge_api.taurus_network_service_get_pledge_actions.return_value = reply
        pledge_api.taurus_network_service_get_pledge_actions_for_approval.return_value = reply
        service = PledgeService(api_client=MagicMock(), pledge_api=pledge_api)
        return service, pledge_api

    @staticmethod
    def _mapped(action_id: str, payload: str, hash_value: str) -> MagicMock:
        action = MagicMock()
        action.id = action_id
        action.metadata = MagicMock()
        action.metadata.payload = payload
        action.metadata.hash = hash_value
        return action

    def test_a_tampered_payload_is_refused_on_both_list_paths(self) -> None:
        from taurus_protect.crypto.hashing import calculate_hex_hash
        from taurus_protect.errors import IntegrityError

        # The hash of a benign top-up, delivered with the payload of something else.
        tampered = self._mapped(
            "a-1", '{"amount":"999"}', calculate_hex_hash('{"amount":"1"}')
        )
        service, _ = self._make_service([MagicMock()])

        with patch(
            "taurus_protect.services.taurus_network.pledge_service.pledge_actions_from_dto",
            return_value=[tampered],
        ):
            with pytest.raises(IntegrityError, match="hash verification failed"):
                service.list_pledge_actions()
            with pytest.raises(IntegrityError, match="hash verification failed"):
                service.list_pledge_actions_for_approval()

    def test_a_matching_payload_passes(self) -> None:
        from taurus_protect.crypto.hashing import calculate_hex_hash

        payload = '{"amount":"1"}'
        good = self._mapped("a-1", payload, calculate_hex_hash(payload))
        service, _ = self._make_service([MagicMock()])

        with patch(
            "taurus_protect.services.taurus_network.pledge_service.pledge_actions_from_dto",
            return_value=[good],
        ):
            actions, _pagination = service.list_pledge_actions()

        assert [a.id for a in actions] == ["a-1"]

    def test_an_action_with_no_metadata_is_not_an_error(self) -> None:
        """An early-status action has nothing to read, so absence is not tampering."""
        bare = MagicMock()
        bare.id = "a-1"
        bare.metadata = None
        service, _ = self._make_service([MagicMock()])

        with patch(
            "taurus_protect.services.taurus_network.pledge_service.pledge_actions_from_dto",
            return_value=[bare],
        ):
            actions, _pagination = service.list_pledge_actions()

        assert [a.id for a in actions] == ["a-1"]

    def test_a_hash_with_no_payload_is_refused(self) -> None:
        from taurus_protect.errors import IntegrityError

        orphan = self._mapped("a-1", "", "abc123")
        service, _ = self._make_service([MagicMock()])

        with patch(
            "taurus_protect.services.taurus_network.pledge_service.pledge_actions_from_dto",
            return_value=[orphan],
        ):
            with pytest.raises(IntegrityError, match="nothing to verify it against"):
                service.list_pledge_actions()

    def test_the_integrity_error_is_not_remapped_to_a_retryable_server_error(self) -> None:
        """T7: the funnel must not turn a tampered response into "retry me".

        ``IntegrityError`` is not an ``APIError``, so the pre-existing
        ``isinstance(e, APIError)`` funnel re-wrapped it as ``ServerError(500)`` --
        whose ``is_retryable()`` is True. A caller following that advice retries a
        response an attacker controls.
        """
        from taurus_protect.crypto.hashing import calculate_hex_hash
        from taurus_protect.errors import IntegrityError, ServerError

        tampered = self._mapped(
            "a-1", '{"amount":"999"}', calculate_hex_hash('{"amount":"1"}')
        )
        service, _ = self._make_service([MagicMock()])

        with patch(
            "taurus_protect.services.taurus_network.pledge_service.pledge_actions_from_dto",
            return_value=[tampered],
        ):
            with pytest.raises(IntegrityError) as exc_info:
                service.list_pledge_actions()

        assert not isinstance(exc_info.value, ServerError)


class TestRejectPledgeActions:
    """Tests for PledgeService.reject_pledge_actions()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.reject_pledge_actions(req=None)

    def test_raises_on_empty_ids(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import (
            RejectPledgeActionsRequest,
        )

        req = RejectPledgeActionsRequest(ids=[], comment="rejected")
        with pytest.raises(ValueError, match="ids list cannot be empty"):
            service.reject_pledge_actions(req=req)

    def test_raises_on_empty_comment(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import (
            RejectPledgeActionsRequest,
        )

        req = RejectPledgeActionsRequest(ids=["a-1"], comment="")
        with pytest.raises(ValueError, match="comment"):
            service.reject_pledge_actions(req=req)


class TestUnpledge:
    """Tests for PledgeService.unpledge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="pledge_id"):
            service.unpledge(pledge_id="")


class TestRejectPledge:
    """Tests for PledgeService.reject_pledge()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_raises_on_empty_id(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import RejectPledgeRequest

        req = RejectPledgeRequest(comment="not approved")
        with pytest.raises(ValueError, match="pledge_id"):
            service.reject_pledge(pledge_id="", req=req)

    def test_raises_on_none_request(self) -> None:
        service, _ = self._make_service()
        with pytest.raises(ValueError, match="request cannot be None"):
            service.reject_pledge(pledge_id="p-1", req=None)

    def test_raises_on_empty_comment(self) -> None:
        service, _ = self._make_service()
        from taurus_protect.models.taurus_network.pledge import RejectPledgeRequest

        req = RejectPledgeRequest(comment="")
        with pytest.raises(ValueError, match="comment"):
            service.reject_pledge(pledge_id="p-1", req=req)


class TestListPledgeWithdrawals:
    """Tests for PledgeService.list_pledge_withdrawals()."""

    def _make_service(self) -> tuple:
        api_client = MagicMock()
        pledge_api = MagicMock()
        service = PledgeService(api_client=api_client, pledge_api=pledge_api)
        return service, pledge_api

    def test_returns_empty_when_no_results(self) -> None:
        service, api = self._make_service()
        resp = MagicMock()
        resp.result = None
        resp.total_items = None
        resp.offset = None
        api.taurus_network_service_get_pledges_withdrawals.return_value = resp

        withdrawals, pagination = service.list_pledge_withdrawals()

        assert withdrawals == []
