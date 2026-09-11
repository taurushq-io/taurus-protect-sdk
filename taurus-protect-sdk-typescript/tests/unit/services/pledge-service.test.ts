/**
 * Unit tests for PledgeService (Taurus Network) — metadata verification and approval.
 *
 * The approval path signs a JSON array of pledge-action metadata hashes. Signing a
 * hash that has not been checked against its payload means the approver's key attests
 * to content nobody read: a response-controlling server can return an action whose
 * `payloadAsString` describes a benign top-up while its `hash` is that of a withdrawal
 * of the whole collateral to an address of the attacker's choosing.
 */

import crypto from "crypto";
import { IntegrityError, ValidationError } from "../../../src/errors";
import { calculateHexHash, verifySignature } from "../../../src/crypto";
import type { TaurusNetworkPledgeApi } from "../../../src/internal/openapi/apis/TaurusNetworkPledgeApi";
import { PledgeService } from "../../../src/services/taurus-network/pledge-service";
import type { PledgeAction } from "../../../src/services/taurus-network/pledge-service";

function createMockApi(): jest.Mocked<TaurusNetworkPledgeApi> {
  return {
    taurusNetworkServiceGetPledge: jest.fn(),
    taurusNetworkServiceGetPledges: jest.fn(),
    taurusNetworkServiceCreatePledge: jest.fn(),
    taurusNetworkServiceUpdatePledge: jest.fn(),
    taurusNetworkServiceAddPledgeCollateral: jest.fn(),
    taurusNetworkServiceWithdrawPledge: jest.fn(),
    taurusNetworkServiceInitiateWithdrawPledge: jest.fn(),
    taurusNetworkServiceUnpledge: jest.fn(),
    taurusNetworkServiceRejectPledge: jest.fn(),
    taurusNetworkServiceGetPledgeActions: jest.fn(),
    taurusNetworkServiceGetPledgeActionsForApproval: jest.fn(),
    taurusNetworkServiceApprovePledgeActions: jest.fn(),
    taurusNetworkServiceRejectPledgeActions: jest.fn(),
    taurusNetworkServiceGetPledgesWithdrawals: jest.fn(),
  } as unknown as jest.Mocked<TaurusNetworkPledgeApi>;
}

/** A payload string and the hash that genuinely commits to it. */
function honestAction(id: string, payload: string): PledgeAction {
  return {
    id,
    pledgeId: "p-1",
    actionType: "PLEDGE_CREATE",
    status: "PENDING",
    metadata: { hash: calculateHexHash(payload), payloadAsString: payload },
  };
}

describe("PledgeService", () => {
  let mockApi: jest.Mocked<TaurusNetworkPledgeApi>;
  let service: PledgeService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new PledgeService(mockApi);
  });

  // ===========================================================================
  // Read paths: every action's hash must cover its payload
  // ===========================================================================

  describe("listPledgeActions", () => {
    it("should return actions whose hashes cover their payloads", async () => {
      const payload = '{"amount":"100","destination":"0xhonest"}';
      mockApi.taurusNetworkServiceGetPledgeActions.mockResolvedValue({
        result: [
          {
            id: "a-1",
            pledgeID: "p-1",
            metadata: { hash: calculateHexHash(payload), payloadAsString: payload },
          },
        ],
      } as never);

      const { actions } = await service.listPledgeActions();

      expect(actions).toHaveLength(1);
      expect(actions[0].metadata?.payloadAsString).toBe(payload);
    });

    it("should fail the whole page when one action's hash does not cover its payload", async () => {
      // The documented attack: alter the payload, leave the hash alone.
      mockApi.taurusNetworkServiceGetPledgeActions.mockResolvedValue({
        result: [
          {
            id: "a-1",
            pledgeID: "p-1",
            metadata: {
              hash: calculateHexHash('{"amount":"100"}'),
              payloadAsString: '{"amount":"1000000"}',
            },
          },
        ],
      } as never);

      await expect(service.listPledgeActions()).rejects.toThrow(IntegrityError);
      await expect(service.listPledgeActions()).rejects.toThrow(/a-1/);
    });

    it("should accept an action with no metadata at all", async () => {
      // An early-status action has nothing to read, so nothing to verify.
      mockApi.taurusNetworkServiceGetPledgeActions.mockResolvedValue({
        result: [{ id: "a-1", pledgeID: "p-1", status: "CREATING" }],
      } as never);

      const { actions } = await service.listPledgeActions();

      expect(actions).toHaveLength(1);
      expect(actions[0].metadata).toBeUndefined();
    });

    it("should reject a hash that arrives with no payload to check it against", async () => {
      mockApi.taurusNetworkServiceGetPledgeActions.mockResolvedValue({
        result: [
          { id: "a-1", pledgeID: "p-1", metadata: { hash: calculateHexHash("x") } },
        ],
      } as never);

      await expect(service.listPledgeActions()).rejects.toThrow(IntegrityError);
    });
  });

  describe("listPledgeActionsForApproval", () => {
    it("should fail the whole page when one action's hash does not cover its payload", async () => {
      mockApi.taurusNetworkServiceGetPledgeActionsForApproval.mockResolvedValue({
        result: [
          {
            id: "a-9",
            pledgeID: "p-1",
            metadata: {
              hash: calculateHexHash('{"amount":"1"}'),
              payloadAsString: '{"amount":"999"}',
            },
          },
        ],
      } as never);

      await expect(service.listPledgeActionsForApproval()).rejects.toThrow(
        IntegrityError
      );
    });

    it("should return verified actions", async () => {
      const payload = '{"amount":"5"}';
      mockApi.taurusNetworkServiceGetPledgeActionsForApproval.mockResolvedValue({
        result: [
          {
            id: "a-9",
            pledgeID: "p-1",
            metadata: { hash: calculateHexHash(payload), payloadAsString: payload },
          },
        ],
      } as never);

      const { actions } = await service.listPledgeActionsForApproval();
      expect(actions).toHaveLength(1);
    });
  });

  // ===========================================================================
  // approvePledgeActions: signs inside the SDK, over hashes it verified itself
  // ===========================================================================

  describe("approvePledgeActions", () => {
    const { privateKey, publicKey } = crypto.generateKeyPairSync("ec", {
      namedCurve: "P-256",
    });

    it("should sign the JSON array of hashes, ids sorted ascending", async () => {
      mockApi.taurusNetworkServiceApprovePledgeActions.mockResolvedValue({} as never);

      // Deliberately out of order: the signed array must not depend on the
      // caller's ordering, it must reproduce the server's own.
      const b = honestAction("a-2", '{"amount":"2"}');
      const a = honestAction("a-1", '{"amount":"1"}');

      const count = await service.approvePledgeActions([b, a], privateKey);

      expect(count).toBe(2);
      const call = mockApi.taurusNetworkServiceApprovePledgeActions.mock.calls[0][0];
      expect(call.body.ids).toEqual(["a-1", "a-2"]);

      const expectedHashes = JSON.stringify([
        a.metadata!.hash,
        b.metadata!.hash,
      ]);
      expect(
        verifySignature(
          publicKey,
          Buffer.from(expectedHashes, "utf-8"),
          call.body.signature as string
        )
      ).toBe(true);
    });

    it("should refuse to sign a hash that does not cover its payload", async () => {
      const action = honestAction("a-1", '{"amount":"1"}');
      const tampered: PledgeAction = {
        ...action,
        metadata: { ...action.metadata!, payloadAsString: '{"amount":"999999"}' },
      };

      await expect(
        service.approvePledgeActions([tampered], privateKey)
      ).rejects.toThrow(IntegrityError);
      await expect(
        service.approvePledgeActions([tampered], privateKey)
      ).rejects.toThrow(/a-1/);
      expect(
        mockApi.taurusNetworkServiceApprovePledgeActions
      ).not.toHaveBeenCalled();
    });

    it("should refuse the whole batch when any one action fails verification", async () => {
      // One signature covers every hash in the batch, so there is no partial
      // submission: signing the survivors would tell the approver they approved
      // less than they did.
      const good = honestAction("a-1", '{"amount":"1"}');
      const bad = honestAction("a-2", '{"amount":"2"}');
      const tampered: PledgeAction = {
        ...bad,
        metadata: { ...bad.metadata!, payloadAsString: '{"amount":"666"}' },
      };

      await expect(
        service.approvePledgeActions([good, tampered], privateKey)
      ).rejects.toThrow(IntegrityError);
      expect(
        mockApi.taurusNetworkServiceApprovePledgeActions
      ).not.toHaveBeenCalled();
    });

    it("should refuse an action carrying a hash but no payload", async () => {
      const action: PledgeAction = {
        id: "a-1",
        metadata: { hash: calculateHexHash("something") },
      };

      await expect(
        service.approvePledgeActions([action], privateKey)
      ).rejects.toThrow(IntegrityError);
    });

    it("should refuse an action with no metadata", async () => {
      await expect(
        service.approvePledgeActions([{ id: "a-1" }], privateKey)
      ).rejects.toThrow(ValidationError);
    });

    it("should refuse an action with no metadata hash", async () => {
      await expect(
        service.approvePledgeActions(
          [{ id: "a-1", metadata: { payloadAsString: "{}" } }],
          privateKey
        )
      ).rejects.toThrow(ValidationError);
    });

    it("should reject an empty batch", async () => {
      await expect(service.approvePledgeActions([], privateKey)).rejects.toThrow(
        ValidationError
      );
    });

    it("should reject a missing private key", async () => {
      const action = honestAction("a-1", '{"amount":"1"}');
      await expect(
        service.approvePledgeActions(
          [action],
          undefined as unknown as crypto.KeyObject
        )
      ).rejects.toThrow(ValidationError);
    });

    it("should pass the comment through", async () => {
      mockApi.taurusNetworkServiceApprovePledgeActions.mockResolvedValue({} as never);
      const action = honestAction("a-1", '{"amount":"1"}');

      await service.approvePledgeActions([action], privateKey, "Q1 collateral");

      const call = mockApi.taurusNetworkServiceApprovePledgeActions.mock.calls[0][0];
      expect(call.body.comment).toBe("Q1 collateral");
    });
  });
});
