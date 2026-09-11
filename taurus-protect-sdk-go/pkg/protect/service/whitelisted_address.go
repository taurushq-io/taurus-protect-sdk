package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedAddressService provides whitelisted address management operations.
type WhitelistedAddressService struct {
	api       *openapi.AddressWhitelistingAPIService
	errMapper *ErrorMapper
	verifier  *helper.WhitelistedAddressVerifier
	logger    Logger
}

// WhitelistedAddressServiceConfig holds configuration for the WhitelistedAddressService.
type WhitelistedAddressServiceConfig struct {
	// SuperAdminKeys are the public keys used to verify governance rules signatures.
	SuperAdminKeys []*ecdsa.PublicKey
	// MinValidSignatures is the minimum number of valid SuperAdmin signatures required.
	MinValidSignatures int
}

// NewWhitelistedAddressServiceWithVerification creates a new WhitelistedAddressService with
// integrity verification enabled. When SuperAdmin keys are provided, all retrieved addresses
// will be cryptographically verified before being returned.
func NewWhitelistedAddressServiceWithVerification(
	client *openapi.APIClient,
	config *WhitelistedAddressServiceConfig,
	opts ...ServiceOption,
) *WhitelistedAddressService {
	svc := &WhitelistedAddressService{
		api:       client.AddressWhitelistingAPI,
		errMapper: NewErrorMapper(),
		logger:    applyServiceOptions(opts).logger,
	}

	// Create verifier if SuperAdmin keys are configured
	if config != nil && len(config.SuperAdminKeys) > 0 {
		svc.verifier = helper.NewWhitelistedAddressVerifier(
			config.SuperAdminKeys,
			config.MinValidSignatures,
		)
	}

	return svc
}

// GetWhitelistedAddress retrieves a whitelisted address by ID.
// If verification is enabled (SuperAdmin keys configured), the address integrity
// will be cryptographically verified before being returned.
func (s *WhitelistedAddressService) GetWhitelistedAddress(ctx context.Context, id string) (*model.WhitelistedAddress, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	resp, httpResp, err := s.api.WhitelistServiceGetWhitelistedAddress(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("whitelisted address not found")
	}

	// SECURITY: Verify hash BEFORE calling mapper to prevent extraction of unverified data.
	// The mapper extracts security-critical fields from PayloadAsString, so we must verify
	// that PayloadAsString has not been tampered with before the mapper runs.
	if err := s.verifyMetadataHashFromDTO(resp.Result); err != nil {
		return nil, err
	}

	// Now safe to call mapper (extracts from verified PayloadAsString)
	addr, err := mapper.WhitelistedAddressFromDTO(resp.Result)
	if err != nil {
		return nil, err
	}

	// Full verification (rules container signatures, whitelist signatures) — always enforced.
	// verifyAddress also replaces the identity fields with the ones the verified payload
	// carries, so what is returned here is payload-derived, not DTO-derived.
	if addr != nil {
		if err := s.verifyAddress(addr); err != nil {
			return nil, err
		}
	}

	return addr, nil
}

// ListWhitelistedAddresses retrieves a list of whitelisted addresses, verifying every
// envelope leniently: a row that fails its own integrity check is excluded and named in
// the result rather than failing the call.
//
//	rows ──▶ hash ──▶ signature ──┬─ ok ─────▶ Addresses
//	                              └─ fail ───▶ ExcludedUnverified (+ log)
//
// Three failures are NOT per-row and abort the whole call instead:
//   - no verifier configured — configuration, identical for every row
//   - the rules container cannot be interpreted by this SDK version
//   - rows came back but none survived — a filtered page must never read as an
//     empty whitelist
func (s *WhitelistedAddressService) ListWhitelistedAddresses(ctx context.Context, opts *model.ListWhitelistedAddressesOptions) (*model.WhitelistedAddressResult, error) {
	req := s.api.WhitelistServiceGetWhitelistedAddresses(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
		if opts.Network != "" {
			req = req.Network(opts.Network)
		}
		// Query matches customer id, address, blockchain, label, memo and address type.
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
		// Which destinations a given wallet or address may send to.
		if opts.AllowedForWalletID != "" {
			req = req.AllowedForWalletId(opts.AllowedForWalletID)
		}
		if opts.AllowedForAddressID != "" {
			req = req.AllowedForAddressId(opts.AllowedForAddressID)
		}
		if opts.TNParticipantID != "" {
			req = req.TnParticipantID(opts.TNParticipantID)
		}
		if len(opts.TagIDs) > 0 {
			req = req.TagIDs(opts.TagIDs)
		}
		if len(opts.ContractTypes) > 0 {
			req = req.ContractTypes(opts.ContractTypes)
		}
		if len(opts.ExchangeAccountIDs) > 0 {
			req = req.ExchangeAccountIds(opts.ExchangeAccountIDs)
		}
		if opts.AddressType != "" {
			req = req.AddressType(opts.AddressType)
		}
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
		if len(opts.Addresses) > 0 {
			req = req.Addresses(opts.Addresses)
		}
		if opts.IncludeForApproval {
			req = req.IncludeForApproval(true)
		}
	}

	// Request normalized rules containers for caching optimization
	req = req.RulesContainerNormalized(true)

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	addresses, excluded, err := s.verifiedAddresses(ctx, resp.Result, resp.RulesContainers)
	if err != nil {
		return nil, err
	}

	return &model.WhitelistedAddressResult{
		Addresses:          addresses,
		Pagination:         adjustedPagination(resp.TotalItems, len(excluded), opts),
		ExcludedUnverified: excluded,
	}, nil
}

// ListWhitelistedAddressesForApproval retrieves whitelisted addresses awaiting approval,
// verified exactly as ListWhitelistedAddresses is.
//
// Both this endpoint and the approval endpoint below are generated in all four SDK
// clients and were wrapped by none of them, so the rows an approver reads before
// whitelisting a destination were not reachable through the SDK at all — verified or not.
//
// IncludeAlreadySignedByUser is exposed because it is what lets an approver see whether
// their own signature is already on a row.
func (s *WhitelistedAddressService) ListWhitelistedAddressesForApproval(
	ctx context.Context,
	opts *model.ListWhitelistedAddressesForApprovalOptions,
) (*model.WhitelistedAddressResult, error) {
	req := s.api.WhitelistServiceGetWhitelistedAddressesForApproval(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Offset > 0 {
			req = req.Offset(fmt.Sprintf("%d", opts.Offset))
		}
		if len(opts.IDs) > 0 {
			req = req.Ids(opts.IDs)
		}
		if opts.Blockchain != "" {
			req = req.Blockchain(opts.Blockchain)
		}
		if opts.Network != "" {
			req = req.Network(opts.Network)
		}
		if opts.AddressType != "" {
			req = req.AddressType(opts.AddressType)
		}
		if opts.Query != "" {
			req = req.Query(opts.Query)
		}
		if opts.IncludeAlreadySignedByUser {
			req = req.IncludeAlreadySignedByUser(true)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	// This endpoint has no normalized-container mode, so the per-row containers are used.
	addresses, excluded, err := s.verifiedAddresses(ctx, resp.Result, nil)
	if err != nil {
		return nil, err
	}

	var limit, offset int64
	if opts != nil {
		limit, offset = opts.Limit, opts.Offset
	}
	return &model.WhitelistedAddressResult{
		Addresses: addresses,
		Pagination: adjustedPagination(resp.TotalItems, len(excluded),
			&model.ListWhitelistedAddressesOptions{Limit: limit, Offset: offset}),
		ExcludedUnverified: excluded,
	}, nil
}

// ApproveWhitelistedAddresses signs and submits an approval for the reviewed whitelisted
// addresses, all-or-nothing.
//
// `selection` comes from a verified read — `result.Select(ids...)` or `result.SelectAll()` — and
// carries the metadata hash each row had AT REVIEW TIME. That pin is the point:
//
//	verified read ─▶ Select(ids) ─▶ pinned hashes
//	                                     │
//	  sort ─▶ ONE filtered verified re-read ─▶ completeness ─▶ PIN MATCH ─┐
//	                                                                      ▼
//	                                                  sign(JSON(hashes)) ─▶ POST once
//
// Without it the approval accepted bare ids and signed whatever the server returned under them,
// so a response-controlling server could substitute a row whose existing signatures already
// satisfy the container it presents and harvest a genuine approver signature over content the
// approver never saw. GovernanceRuleService.ApproveRulesProposal was hardened the same way with
// its mandatory expectedContainerHash; an empty pin must not silently restore the old behaviour,
// which is why an empty selection is an error rather than "approve nothing".
//
// A pin mismatch is a REAL signal, not a false positive. The metadata hash is recomputed by the
// server on every read from the immutable envelope plus the row's live linked-address and
// linked-wallet rows, so it moves when a linked address is renamed. Normally that also breaks
// signature coverage and the row is excluded anyway; for a legacy-signed row it can move while
// the row still verifies. Either way the content changed since review, so re-read, re-review and
// re-approve — exactly as for a changed rules proposal.
//
// Any address that is missing, fails verification, or whose hash no longer matches the pin aborts
// the whole call and nothing is signed: one signature covers every hash in the batch, so a
// partial approval would mean the caller believes they approved more than they did.
func (s *WhitelistedAddressService) ApproveWhitelistedAddresses(
	ctx context.Context,
	selection *model.WhitelistedAddressApproval,
	privateKey *ecdsa.PrivateKey,
	comment string,
) error {
	if selection.IsEmpty() {
		return fmt.Errorf("selection cannot be empty: pin the rows with " +
			"result.Select(ids...) or result.SelectAll() from a verified read")
	}
	if privateKey == nil {
		return fmt.Errorf("privateKey cannot be nil")
	}
	if comment == "" {
		return fmt.Errorf("comment is required")
	}

	ids := selection.IDs()

	// The endpoint requires ascending order, and sorting here also makes the signed
	// order independent of the order the caller passed.
	sortedIDs := make([]string, len(ids))
	copy(sortedIDs, ids)
	for _, id := range sortedIDs {
		if _, err := strconv.ParseInt(id, 10, 64); err != nil {
			return fmt.Errorf("whitelisted address ID %q is not a valid numeric ID: %w", id, err)
		}
	}
	sort.Slice(sortedIDs, func(i, j int) bool {
		a, _ := strconv.ParseInt(sortedIDs[i], 10, 64)
		b, _ := strconv.ParseInt(sortedIDs[j], 10, 64)
		return a < b
	})

	// ONE id-filtered page through the verifying list path, not one GET per id.
	result, err := s.ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{
		IDs:                sortedIDs,
		Limit:              int64(len(sortedIDs)),
		IncludeForApproval: true,
	})
	if err != nil {
		return fmt.Errorf("refusing to sign: the verified read failed: %w", err)
	}

	byID := make(map[string]*model.WhitelistedAddress, len(result.Addresses))
	for _, addr := range result.Addresses {
		if addr != nil {
			byID[addr.ID] = addr
		}
	}

	hashes := make([]string, 0, len(sortedIDs))
	for _, id := range sortedIDs {
		addr, ok := byID[id]
		if !ok {
			// Excluded, or simply absent. Either way a page that omits a row must not
			// become an approval of fewer rows than the caller asked for.
			return &model.IntegrityError{
				Message: fmt.Sprintf("refusing to sign: address %s was not returned by the verified read", id),
			}
		}
		if addr.Metadata == nil || addr.Metadata.Hash == "" {
			return &model.IntegrityError{
				Message: fmt.Sprintf("refusing to sign: address %s has no metadata hash", id),
			}
		}

		// The pin. Constant-time because this compares hash material, matching
		// ApproveRulesProposal's use of subtle.ConstantTimeCompare on its container pin.
		pinnedHash, pinned := selection.PinnedHash(id)
		if !pinned {
			return &model.IntegrityError{
				Message: fmt.Sprintf("refusing to sign: address %s is not in the reviewed selection", id),
			}
		}
		if subtle.ConstantTimeCompare([]byte(pinnedHash), []byte(addr.Metadata.Hash)) != 1 {
			return &model.IntegrityError{
				Message: fmt.Sprintf("refusing to sign: whitelisted address %s changed since it was "+
					"reviewed: reviewed hash %s, re-read hash %s. Re-read, re-review and re-approve",
					id, pinnedHash, addr.Metadata.Hash),
			}
		}

		hashes = append(hashes, addr.Metadata.Hash)
	}

	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		return fmt.Errorf("failed to serialize hashes: %w", err)
	}

	signature, err := crypto.SignData(privateKey, hashesJSON)
	if err != nil {
		return fmt.Errorf("failed to sign address hashes: %w", err)
	}

	approveReq := openapi.TgvalidatordApproveWhitelistedAddressRequest{
		Ids:       sortedIDs,
		Signature: signature,
		Comment:   comment,
	}

	if _, httpResp, err := s.api.WhitelistServiceApproveWhitelistedAddress(ctx).
		Body(approveReq).
		Execute(); err != nil {
		return s.errMapper.MapError(err, httpResp)
	}

	return nil
}

// verifiedAddresses verifies a page of address envelopes and returns the survivors plus
// the rows it withheld.
//
//	rows ──▶ hash ──▶ signature ──┬─ ok ─────▶ returned
//	                              └─ fail ───▶ excluded + logged, list survives
//
// LENIENT per row, but three things abort the whole call rather than being excluded:
// an unconfigured verifier (identical for every row, so excluding would return an empty
// whitelist presented as complete), a ContainerIntegrityError (it invalidates every row
// judged against that container), and rows-returned-but-none-surviving.
//
// Shared by ListWhitelistedAddresses and ListWhitelistedAddressesForApproval. The
// for-approval endpoint had no verified reader at all, and duplicating ~80 lines of
// verification into a second one is precisely how the list paths drifted from get()
// in the first place.
func (s *WhitelistedAddressService) verifiedAddresses(
	ctx context.Context,
	rows []openapi.TgvalidatordSignedWhitelistedAddressEnvelope,
	normalizedContainers []openapi.TgvalidatordHashRulesContainer,
) ([]*model.WhitelistedAddress, []model.ExcludedWhitelistedAddress, error) {
	containers := buildRulesContainerCache(s.verifier, normalizedContainers)

	// A missing verifier is CONFIGURATION, not one bad row: it fails identically for
	// every address, so excluding on it would return an empty whitelist presented as
	// complete — strictly worse than refusing. Checked once, here, before the per-row
	// loop, so everything inside that loop is genuinely per-row.
	if s.verifier == nil {
		return nil, nil, &model.IntegrityError{
			Message: "verification is required but no verifier is configured",
		}
	}

	// SECURITY: verify the hash BEFORE the mapper reads PayloadAsString.
	//
	//	rows ──▶ hash ──▶ signature ──┬─ ok ─────▶ returned
	//	                              └─ fail ───▶ excluded + logged, list survives
	//
	// One unverifiable row previously failed the whole call, which took down the
	// whitelist for every consumer — and listing is how an operator would find the
	// bad row, so the failure hid its own cause. Excluding is equally fail-closed: an
	// omitted destination cannot be selected, so a transfer to it is refused.
	excluded := make([]model.ExcludedWhitelistedAddress, 0)

	verified := make([]openapi.TgvalidatordSignedWhitelistedAddressEnvelope, 0, len(rows))
	for i, dto := range rows {
		if err := s.verifyMetadataHashFromDTO(&dto); err != nil {
			addrID := ""
			if dto.Id != nil {
				addrID = *dto.Id
			}
			s.logExcluded(ctx, addrID, i, err)
			excluded = append(excluded, model.ExcludedWhitelistedAddress{ID: addrID, Reason: err.Error()})
			continue
		}
		verified = append(verified, dto)
	}

	// Now safe to call mapper (extracts from verified PayloadAsString).
	//
	// Mapped ROW BY ROW rather than in bulk: a signed payload this SDK cannot parse is one
	// bad row, so it belongs in ExcludedUnverified like every other row-level failure. The
	// bulk form would fail the whole page.
	addresses := make([]*model.WhitelistedAddress, 0, len(verified))
	for i := range verified {
		addr, err := mapper.WhitelistedAddressFromDTO(&verified[i])
		if err != nil {
			addrID := ""
			if verified[i].Id != nil {
				addrID = *verified[i].Id
			}
			s.logExcluded(ctx, addrID, i, err)
			excluded = append(excluded, model.ExcludedWhitelistedAddress{ID: addrID, Reason: err.Error()})
			continue
		}
		if addr == nil {
			continue
		}
		// Look up cached rules container by hash
		var cached *model.DecodedRulesContainer
		if addr.RulesContainerHash != "" {
			cached = containers.verified[addr.RulesContainerHash]
		}
		// Name the container that failed. Falling through would report "required data
		// missing" — in normalized mode the per-row container is empty by design, so
		// that reads as a malformed response rather than a container that did not verify.
		if cached == nil && addr.RulesContainerHash != "" {
			if reason, failed := containers.reasonFor(addr.RulesContainerHash); failed {
				err := &model.IntegrityError{
					Message: fmt.Sprintf("rules container %s failed verification: %s",
						addr.RulesContainerHash, reason),
				}
				s.logExcluded(ctx, addr.ID, -1, err)
				excluded = append(excluded, model.ExcludedWhitelistedAddress{ID: addr.ID, Reason: err.Error()})
				continue
			}
		}
		if err := s.verifyAddressWithCache(addr, cached); err != nil {
			// A container this SDK cannot interpret is not a property of this row —
			// it invalidates every row judged against it. Excluding them one by one
			// would empty the whitelist and report success.
			var containerErr *model.ContainerIntegrityError
			if errors.As(err, &containerErr) {
				return nil, nil, err
			}
			s.logExcluded(ctx, addr.ID, -1, err)
			excluded = append(excluded, model.ExcludedWhitelistedAddress{ID: addr.ID, Reason: err.Error()})
			continue
		}
		addresses = append(addresses, addr)
	}

	// Rows came back but none survived: systemic, not an empty whitelist, and returning
	// an empty list would be indistinguishable from one.
	if len(rows) > 0 && len(addresses) == 0 {
		firstReason := "unknown"
		if len(excluded) > 0 {
			firstReason = excluded[0].Reason
		}
		return nil, nil, &model.IntegrityError{
			Message: fmt.Sprintf("all %d whitelisted address(es) failed verification; first failure: %s",
				len(rows), firstReason),
		}
	}

	return addresses, excluded, nil
}

// adjustedPagination builds the page window with TotalItems reduced by the number of
// rows excluded for failing verification.
//
// The server counts rows it returned; the caller receives only the ones that verified.
// Reporting the server's total would make TotalItems promise rows that can never be
// read, and would let a filtered page pass for a complete one.
//
// HasMore is derived from the server's UNREDUCED total, because excludedCount covers
// this page only: subtracting it from a global total ended pagination early, so a page
// with many exclusions mid-result-set reported HasMore=false while whole further pages
// existed (total 250, limit 100, offset 100, 60 excluded => 200 < 190 is false, yet 50
// rows remain at offset 200).
func adjustedPagination(totalItems *string, excludedCount int, opts *model.ListWhitelistedAddressesOptions) *model.Pagination {
	if totalItems == nil {
		return nil
	}
	pagination := &model.Pagination{}
	var serverTotal int64
	if total, parseErr := strconv.ParseInt(*totalItems, 10, 64); parseErr == nil {
		serverTotal = total
		pagination.TotalItems = total - int64(excludedCount)
		if pagination.TotalItems < 0 {
			pagination.TotalItems = 0
		}
	}
	if opts != nil {
		pagination.Limit = opts.Limit
		pagination.Offset = opts.Offset
		// Overflow-safe form; offset+limit can wrap on caller-supplied values.
		pagination.HasMore = serverTotal > opts.Offset && serverTotal-opts.Offset > opts.Limit
	}
	return pagination
}

// containerVerifier is what the container caches need: steps 2-3 for one container.
// Both whitelist verifiers satisfy it, so the caches serve both services.
type containerVerifier interface {
	VerifyAndDecodeRulesContainer(
		rulesContainerBase64 string,
		rulesSignaturesBase64 string,
		rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
		userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
	) (*model.DecodedRulesContainer, error)
}

// rulesContainerCache holds the containers that verified, keyed by hash, and why each one
// that failed was rejected. Keeping the failures is what lets a row name the container
// responsible instead of reporting the generic "required data missing".
type rulesContainerCache struct {
	verified map[string]*model.DecodedRulesContainer
	failed   map[string]string
}

// reasonFor returns why the container behind this hash was rejected, if it was.
func (c rulesContainerCache) reasonFor(containerHash string) (string, bool) {
	reason, failed := c.failed[containerHash]
	return reason, failed
}

// inlineContainerCache memoizes steps 2-3 per page for endpoints that carry the container
// on every row rather than in a deduplicated array — the contracts list has no
// rulesContainerNormalized parameter, so N rows sharing one container otherwise pay N
// SuperAdmin verifications and N protobuf decodes. Outcomes are memoized either way, so a
// container that fails is not re-verified per row.
type inlineContainerCache struct {
	verifier containerVerifier
	decoded  map[string]*model.DecodedRulesContainer
	failures map[string]error
}

func newInlineContainerCache(verifier containerVerifier) *inlineContainerCache {
	return &inlineContainerCache{
		verifier: verifier,
		decoded:  make(map[string]*model.DecodedRulesContainer),
		failures: make(map[string]error),
	}
}

// get returns the verified container for these bytes, verifying on first sight.
func (c *inlineContainerCache) get(containerBase64, signaturesBase64 string) (*model.DecodedRulesContainer, error) {
	// Both halves key the entry: the same container under different signatures is a
	// different claim and must be verified again.
	key := containerBase64 + "\x00" + signaturesBase64

	if container, ok := c.decoded[key]; ok {
		return container, nil
	}
	if err, ok := c.failures[key]; ok {
		return nil, err
	}

	container, err := c.verifier.VerifyAndDecodeRulesContainer(
		containerBase64,
		signaturesBase64,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
	)
	if err != nil {
		c.failures[key] = err
		return nil, err
	}
	c.decoded[key] = container
	return container, nil
}

// containerHashLabel recomputes the label validatord files a normalized rules container
// under, so a row's rulesContainerHash resolves only to bytes that really hash to it.
//
// The convention is validatord's, at internal/api/v1/whitelist-controller.go:
//
//	h := base64.StdEncoding.EncodeToString(crypto.Sha256([]byte(e.GetRulesContainer())))
//
// Note what is hashed: `GetRulesContainer()` is ALREADY a base64 string there, so the
// digest is over the base64 TEXT, not over the protobuf bytes, and the result is base64
// rather than hex. Decoding the container first and hashing the protobuf yields a
// different value and would reject every container. Do not "fix" it to hash the decoded
// bytes — and do not confuse it with `enforcedRulesHash`, which is base64(SHA256(raw
// protobuf)) and is a backlink to a ruleset's predecessor, not its own identity.
func containerHashLabel(containerBase64 string) string {
	sum := sha256.Sum256([]byte(containerBase64))
	return base64.StdEncoding.EncodeToString(sum[:])
}

// buildRulesContainerCache verifies each distinct container in a normalized response once
// and caches it by hash.
//
// A container that fails is RECORDED, not skipped silently, but does not fail the whole
// call: a page can legitimately span a rules-container change, so rows judged against a
// good container still return. If every row referenced the bad one, the caller's
// none-survived check makes it a hard error.
func buildRulesContainerCache(
	verifier containerVerifier,
	containers []openapi.TgvalidatordHashRulesContainer,
) rulesContainerCache {
	result := rulesContainerCache{
		verified: make(map[string]*model.DecodedRulesContainer),
		failed:   make(map[string]string),
	}
	if len(containers) == 0 || verifier == nil {
		return result
	}

	// Deduplicate by base64 container string to avoid re-verifying identical containers
	verifiedContainers := make(map[string]*model.DecodedRulesContainer)
	failedContainers := make(map[string]string)

	for _, hashContainer := range containers {
		containerHash := ""
		if hashContainer.Hash != nil {
			containerHash = *hashContainer.Hash
		}
		containerBase64 := ""
		if hashContainer.RulesContainer != nil {
			containerBase64 = *hashContainer.RulesContainer
		}
		signaturesBase64 := ""
		if hashContainer.RulesSignatures != nil {
			signaturesBase64 = *hashContainer.RulesSignatures
		}

		if containerHash == "" || containerBase64 == "" {
			continue
		}

		// The hash is the LABEL a row uses to pick its container, and it is supplied by
		// the same response as the container. Recompute it: otherwise a server can file
		// container A under the label of container B and steer any row to any other
		// validly-signed container — an older ruleset with a weaker group threshold, say.
		// Both would still pass step 2, so signature verification alone does not catch it.
		if computed := containerHashLabel(containerBase64); computed != containerHash {
			reason := fmt.Sprintf(
				"rules container hash mismatch: response labelled it %s but its bytes hash to %s",
				containerHash, computed)
			failedContainers[containerBase64] = reason
			result.failed[containerHash] = reason
			continue
		}

		// Already known bad (dedup by content): record against this hash too.
		if reason, bad := failedContainers[containerBase64]; bad {
			result.failed[containerHash] = reason
			continue
		}

		// Check if already verified (dedup by content)
		decoded, exists := verifiedContainers[containerBase64]
		if !exists {
			var err error
			decoded, err = verifier.VerifyAndDecodeRulesContainer(
				containerBase64,
				signaturesBase64,
				mapper.RulesContainerFromBase64,
				mapper.UserSignaturesFromBase64,
			)
			if err != nil {
				failedContainers[containerBase64] = err.Error()
				result.failed[containerHash] = err.Error()
				continue
			}
			verifiedContainers[containerBase64] = decoded
		}

		result.verified[containerHash] = decoded
	}

	return result
}

// verifyAddressWithCache performs verification with an optional cached rules container.
// When cached is non-nil, steps 2-3 are skipped and RulesContainer on the address
// may be empty (normalized mode moves it to the response-level array).
// logExcluded reports a whitelisted address dropped for failing integrity. index is
// the position in the response, or -1 once the rows have been mapped.
//
// Metadata only: id and reason. Never the payload — the caller's logger writes this
// to disk and the reason is an error string we control, not response content.
func (s *WhitelistedAddressService) logExcluded(ctx context.Context, addrID string, index int, err error) {
	fields := []Field{
		{Key: "resource", Value: "whitelisted_address"},
		{Key: "address_id", Value: addrID},
		{Key: "reason", Value: err.Error()},
	}
	if index >= 0 {
		fields = append(fields, Field{Key: "index", Value: index})
	}
	s.logger.Warn(ctx, "whitelisted address excluded: integrity verification failed", fields...)
}

func (s *WhitelistedAddressService) verifyAddressWithCache(addr *model.WhitelistedAddress, cached *model.DecodedRulesContainer) error {
	// The nil-verifier case is handled by the caller before the per-row loop: it is
	// row-independent, and treating it as a per-row failure would silently empty the
	// whole list.
	if s.verifier == nil {
		return &model.IntegrityError{Message: "verification is required but no verifier is configured"}
	}

	if addr.Metadata == nil || addr.SignedAddress == nil {
		return &model.IntegrityError{Message: "verification enabled but required data missing"}
	}

	// When not using cache, RulesContainer must be present on the address itself
	if cached == nil && addr.RulesContainer == "" {
		return &model.IntegrityError{Message: "verification enabled but required data missing"}
	}

	result, err := s.verifier.VerifyWhitelistedAddress(
		addr,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
		cached,
	)
	if err != nil {
		return err
	}
	applyVerifiedIdentity(addr, result.VerifiedAddress)
	return nil
}

// applyVerifiedIdentity overwrites the identity fields on addr with the values parsed from the
// payload verification actually cleared.
//
// Both service seams used to discard the verifier's result (`_, err :=`), so every read path
// except GetWhitelistedAddressEnvelope returned the mapper's object instead. That left two
// defects live:
//
//   - TnParticipantID came from the unsigned DTO although the signed payload carries it and
//     nothing cross-checked the two, so a server could attribute a genuinely whitelisted
//     address to a different Taurus-NETWORK participant and the SDK returned it as verified.
//   - When step 4 matched a LEGACY hash, the mapper had parsed the DELIVERED payload — which
//     carries the members the legacy strip removed, i.e. values no signature covered.
//
// Network is the one field that keeps its DTO value when the payload omits it: the payload
// legitimately has no `network` member when the governing rule's includeNetworkInPayload is
// off, and blanking the label would lose information the caller needs. That is safe only
// because rule selection no longer lets an unsigned network reach a weaker quorum — see
// helper.ResolveRuleKey.
func applyVerifiedIdentity(addr *model.WhitelistedAddress, verified *model.WhitelistedAddress) {
	if addr == nil || verified == nil {
		return
	}
	addr.Blockchain = verified.Blockchain
	addr.TnParticipantID = verified.TnParticipantID
	if verified.Network != "" {
		addr.Network = verified.Network
	}
	addr.Address = verified.Address
	addr.Label = verified.Label
	addr.Memo = verified.Memo
	addr.CustomerId = verified.CustomerId
	addr.ContractType = verified.ContractType
	addr.AddressType = verified.AddressType
	addr.ExchangeAccountId = verified.ExchangeAccountId
	addr.LinkedInternalAddresses = verified.LinkedInternalAddresses
	addr.LinkedWallets = verified.LinkedWallets
}

// verifyMetadataHashFromDTO verifies the metadata hash before calling the mapper.
// SECURITY: This MUST be called before the mapper to prevent extraction of unverified data.
// Returns nil if verification passes, or an error if hash doesn't match.
func (s *WhitelistedAddressService) verifyMetadataHashFromDTO(dto *openapi.TgvalidatordSignedWhitelistedAddressEnvelope) error {
	if dto == nil || dto.Metadata == nil {
		return nil // No metadata to verify
	}

	payloadAsString := dto.Metadata.PayloadAsString
	providedHash := dto.Metadata.Hash

	// Skip verification if no payload data (e.g., addresses pending approval)
	if payloadAsString == nil || *payloadAsString == "" {
		return nil
	}

	// If there's payloadAsString but no hash, that's an error
	if providedHash == nil || *providedHash == "" {
		return &model.IntegrityError{Message: "metadata hash is missing but payloadAsString is present"}
	}

	computedHash := crypto.CalculateHexHash(*payloadAsString)
	if !helper.ConstantTimeCompare(computedHash, *providedHash) {
		return &model.IntegrityError{Message: "metadata hash verification failed"}
	}

	return nil
}

// verifyAddress performs the 6-step integrity verification on a whitelisted address.
// Returns nil if verification passes, or an error describing the failure.
func (s *WhitelistedAddressService) verifyAddress(addr *model.WhitelistedAddress) error {
	if s.verifier == nil {
		return &model.IntegrityError{Message: "verification is required but no verifier is configured"}
	}

	// Verification is enabled but required data is missing — this is an error.
	// An attacker could strip verification data to bypass checks.
	if addr.Metadata == nil || addr.RulesContainer == "" || addr.SignedAddress == nil {
		return &model.IntegrityError{Message: "verification enabled but required data missing"}
	}

	result, err := s.verifier.VerifyWhitelistedAddress(
		addr,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
	)
	if err != nil {
		return err
	}
	applyVerifiedIdentity(addr, result.VerifiedAddress)
	return nil
}

// GetWhitelistedAddressEnvelope retrieves a whitelisted address envelope by ID and performs
// the complete 6-step verification flow. The envelope contains both the verified whitelisted
// address and the decoded rules container.
//
// This method requires verification to be enabled (SuperAdmin keys must be configured).
// If verification is not configured, an error is returned.
//
// The 6-step verification flow:
// 1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
// 2. Verify rules container signatures (SuperAdmin signatures)
// 3. Decode rules container (base64 -> protobuf -> model)
// 4. Verify hash coverage (metadata.hash in at least one signature.hashes)
// 5. Verify whitelist signatures (user signatures meet governance thresholds)
// 6. Parse WhitelistedAddress from verified payload
func (s *WhitelistedAddressService) GetWhitelistedAddressEnvelope(
	ctx context.Context,
	id string,
) (*model.WhitelistedAddressEnvelope, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if s.verifier == nil {
		return nil, &model.IntegrityError{
			Message: "verification is required for GetWhitelistedAddressEnvelope but no SuperAdmin keys are configured",
		}
	}

	resp, httpResp, err := s.api.WhitelistServiceGetWhitelistedAddress(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, fmt.Errorf("whitelisted address not found")
	}

	// Create envelope from DTO
	envelope := mapper.WhitelistedAddressEnvelopeFromDTO(resp.Result)

	// Initialize and verify the envelope
	if err := s.initializeEnvelope(envelope); err != nil {
		return nil, err
	}

	return envelope, nil
}

// initializeEnvelope performs the 6-step verification and populates the verified fields.
func (s *WhitelistedAddressService) initializeEnvelope(envelope *model.WhitelistedAddressEnvelope) error {
	if envelope == nil {
		return fmt.Errorf("envelope cannot be nil")
	}

	// Check required data for verification
	if envelope.Metadata == nil {
		return &model.IntegrityError{Message: "metadata is required for verification"}
	}
	if envelope.RulesContainer == "" {
		return &model.IntegrityError{Message: "rules container is required for verification"}
	}
	if envelope.SignedAddress == nil {
		return &model.IntegrityError{Message: "signed address is required for verification"}
	}

	// Create a temporary WhitelistedAddress for verification
	tempAddr := &model.WhitelistedAddress{
		ID:                      envelope.ID,
		Blockchain:              envelope.Blockchain,
		Network:                 envelope.Network,
		Metadata:                envelope.Metadata,
		SignedAddress:           envelope.SignedAddress,
		RulesContainer:          envelope.RulesContainer,
		RulesSignatures:         envelope.RulesSignatures,
		LinkedInternalAddresses: envelope.LinkedInternalAddresses,
		LinkedWallets:           envelope.LinkedWallets,
	}

	// Perform the 6-step verification
	result, err := s.verifier.VerifyWhitelistedAddress(
		tempAddr,
		mapper.RulesContainerFromBase64,
		mapper.UserSignaturesFromBase64,
	)
	if err != nil {
		return err
	}

	// Set verified data on envelope
	envelope.SetVerified(result.VerifiedAddress, result.RulesContainer)

	return nil
}
