/**
 * Unit tests for the governance proposal lifecycle:
 * updateRulesProposal / approveRulesProposal / rejectRulesProposal.
 */
import * as crypto from 'crypto';
import { GovernanceRuleService } from '../../../src/services/governance-rule-service';
import { APIError, IntegrityError, ValidationError } from '../../../src/errors';
import { verifySignature } from '../../../src/crypto';
import { rulesContainerToBase64 } from '../../../src/mappers/protobuf-rules-container-encode';
import { createEmptyRulesContainer } from '../../../src/models/governance-rules';
import type { DecodedRulesContainer } from '../../../src/models/governance-rules';
import type { GovernanceRulesApi } from '../../../src/internal/openapi/apis/GovernanceRulesApi';

function mockApi(): jest.Mocked<GovernanceRulesApi> {
  return {
    ruleServiceGetRulesProposal: jest.fn(),
    ruleServiceUpdateRulesProposal: jest.fn().mockResolvedValue({}),
    ruleServiceApproveRulesProposal: jest.fn().mockResolvedValue({}),
    ruleServiceRejectRulesProposal: jest.fn().mockResolvedValue({}),
  } as unknown as jest.Mocked<GovernanceRulesApi>;
}

/**
 * Verification is mandatory at construction now, so even the proposal tests need keys.
 * The proposal paths never verify — a proposal legitimately carries 0..N signatures —
 * so any P-256 public key is a valid fixture here.
 */
function verifyingConfig() {
  return {
    superAdminKeys: [
      crypto.generateKeyPairSync('ec', { namedCurve: 'P-256' }).publicKey,
    ],
    minValidSignatures: 1,
  };
}

function sampleContainer(): DecodedRulesContainer {
  return {
    ...createEmptyRulesContainer(),
    minimumDistinctUserSignatures: 1,
    transactionRules: [{ key: 'ETH/transfer', columns: [], lines: [] }],
  };
}

/** The pin an approver passes back: sha256 of the DECODED container bytes. */
function pinFor(containerB64: string): string {
  return crypto.createHash('sha256').update(Buffer.from(containerB64, 'base64')).digest('hex');
}

describe('GovernanceRuleService proposal lifecycle', () => {
  it('updateRulesProposal submits the encoded container', async () => {
    const api = mockApi();
    const svc = new GovernanceRuleService(api, verifyingConfig());
    const container = sampleContainer();

    await svc.updateRulesProposal(container);

    expect(api.ruleServiceUpdateRulesProposal).toHaveBeenCalledTimes(1);
    const arg = (api.ruleServiceUpdateRulesProposal as jest.Mock).mock.calls[0][0];
    expect(arg.body.rulesContainer).toBe(rulesContainerToBase64(container));
    // round-trips back to the same content
    // (decode handled by other tests; here we only assert the wire value matches)
  });

  it('approveRulesProposal signs the pending container bytes and submits', async () => {
    const api = mockApi();
    const { privateKey, publicKey } = crypto.generateKeyPairSync('ec', { namedCurve: 'P-256' });
    const pendingB64 = rulesContainerToBase64(sampleContainer());
    (api.ruleServiceGetRulesProposal as jest.Mock).mockResolvedValue({
      result: { rulesContainer: pendingB64, rulesSignatures: [], locked: false, trails: [] },
    });

    const svc = new GovernanceRuleService(api, verifyingConfig());
    await svc.approveRulesProposal(privateKey, 'looks good', pinFor(pendingB64));

    expect(api.ruleServiceApproveRulesProposal).toHaveBeenCalledTimes(1);
    const body = (api.ruleServiceApproveRulesProposal as jest.Mock).mock.calls[0][0].body;
    expect(body.comment).toBe('looks good');

    // The submitted signature must verify against the decoded pending bytes.
    const data = new Uint8Array(Buffer.from(pendingB64, 'base64'));
    expect(verifySignature(publicKey, data, body.signature)).toBe(true);
    // ECDSA is non-deterministic, so we verify the signature rather than compare bytes;
    // a signature over different data must NOT verify.
    const otherData = new Uint8Array([...data, 1]);
    expect(verifySignature(publicKey, otherData, body.signature)).toBe(false);
  });

  it('approveRulesProposal throws when no proposal is pending', async () => {
    const api = mockApi();
    (api.ruleServiceGetRulesProposal as jest.Mock).mockResolvedValue({
      result: { rulesContainer: undefined, rulesSignatures: [], locked: false, trails: [] },
    });
    const { privateKey } = crypto.generateKeyPairSync('ec', { namedCurve: 'P-256' });
    const svc = new GovernanceRuleService(api, verifyingConfig());

    await expect(
      svc.approveRulesProposal(privateKey, 'x', 'a'.repeat(64))
    ).rejects.toBeInstanceOf(ValidationError);
    expect(api.ruleServiceApproveRulesProposal).not.toHaveBeenCalled();
  });

  // The content-binding gate. A server able to shape responses can serve the benign
  // proposal to the review call and a DIFFERENT container to the re-fetch inside
  // approveRulesProposal. Without the pin that yields a genuine SuperAdmin signature over
  // attacker-chosen bytes, which then verifies clean everywhere -- including in verifiers
  // that never trusted that server. The approval must abort and submit nothing.
  it('approveRulesProposal refuses when the pending container changed since review', async () => {
    const api = mockApi();
    const reviewedB64 = rulesContainerToBase64(sampleContainer());
    const reviewedBytes = Buffer.from(reviewedB64, 'base64');
    // The re-fetch serves the reviewed bytes with a field appended -- the shape a
    // merge-on-concatenation protobuf attack produces.
    const substituted = Buffer.concat([reviewedBytes, Buffer.from([0x40, 0x01])]).toString('base64');
    (api.ruleServiceGetRulesProposal as jest.Mock).mockResolvedValue({
      result: { rulesContainer: substituted, rulesSignatures: [], locked: false, trails: [] },
    });
    const { privateKey } = crypto.generateKeyPairSync('ec', { namedCurve: 'P-256' });
    const svc = new GovernanceRuleService(api, verifyingConfig());

    await expect(
      svc.approveRulesProposal(privateKey, 'lgtm', pinFor(reviewedB64))
    ).rejects.toBeInstanceOf(IntegrityError);
    expect(api.ruleServiceApproveRulesProposal).not.toHaveBeenCalled();
  });

  // The pin is mandatory: an empty one would restore the unpinned behaviour silently.
  it('approveRulesProposal requires a pin', async () => {
    const api = mockApi();
    const { privateKey } = crypto.generateKeyPairSync('ec', { namedCurve: 'P-256' });
    const svc = new GovernanceRuleService(api, verifyingConfig());

    await expect(svc.approveRulesProposal(privateKey, 'lgtm', '')).rejects.toBeInstanceOf(
      ValidationError
    );
    expect(api.ruleServiceGetRulesProposal).not.toHaveBeenCalled();
  });

  it('rejectRulesProposal submits the comment', async () => {
    const api = mockApi();
    const svc = new GovernanceRuleService(api, verifyingConfig());

    await svc.rejectRulesProposal('not compliant');

    expect(api.ruleServiceRejectRulesProposal).toHaveBeenCalledTimes(1);
    const body = (api.ruleServiceRejectRulesProposal as jest.Mock).mock.calls[0][0].body;
    expect(body.comment).toBe('not compliant');
  });

  it('updateRulesProposal maps API failures to an APIError', async () => {
    const api = mockApi();
    (api.ruleServiceUpdateRulesProposal as jest.Mock).mockRejectedValue(new Error('boom'));
    const svc = new GovernanceRuleService(api, verifyingConfig());
    await expect(svc.updateRulesProposal(sampleContainer())).rejects.toBeInstanceOf(APIError);
  });

  it('approveRulesProposal maps approve-call failures to an APIError', async () => {
    const api = mockApi();
    const pendingB64 = rulesContainerToBase64(sampleContainer());
    (api.ruleServiceGetRulesProposal as jest.Mock).mockResolvedValue({
      result: { rulesContainer: pendingB64, rulesSignatures: [], locked: false, trails: [] },
    });
    (api.ruleServiceApproveRulesProposal as jest.Mock).mockRejectedValue(new Error('boom'));
    const { privateKey } = crypto.generateKeyPairSync('ec', { namedCurve: 'P-256' });
    const svc = new GovernanceRuleService(api, verifyingConfig());
    await expect(
      svc.approveRulesProposal(privateKey, 'x', pinFor(pendingB64))
    ).rejects.toBeInstanceOf(APIError);
  });

  it('rejectRulesProposal maps API failures to an APIError', async () => {
    const api = mockApi();
    (api.ruleServiceRejectRulesProposal as jest.Mock).mockRejectedValue(new Error('boom'));
    const svc = new GovernanceRuleService(api, verifyingConfig());
    await expect(svc.rejectRulesProposal('x')).rejects.toBeInstanceOf(APIError);
  });
});
