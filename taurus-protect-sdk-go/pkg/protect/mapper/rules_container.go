package mapper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log/slog"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"google.golang.org/protobuf/proto"
)

// RulesContainerFromBase64 decodes a base64-encoded protobuf RulesContainer into a model.
func RulesContainerFromBase64(base64Data string) (*model.DecodedRulesContainer, error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return RulesContainerFromBytes(data)
}

// RulesContainerFromBytes decodes raw protobuf bytes into a DecodedRulesContainer.
//
// Data the current SDK version does not model — unknown protobuf fields,
// unknown enum values, unknown cell types — is never dropped: it is captured
// on the decoded model (UnknownFields buffers, numeric enum names, RawCell /
// raw sources) and re-emitted verbatim by RulesContainerToBytes. When any such
// data is present a warning is logged and
// DecodedRulesContainer.HasUnknownFields reports true.
func RulesContainerFromBytes(data []byte) (*model.DecodedRulesContainer, error) {
	var pbContainer pb.RulesContainer
	if err := proto.Unmarshal(data, &pbContainer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return rulesContainerFromProto(&pbContainer)
}

// decodeStats counts data preserved-but-not-modeled during a decode, to warn once.
type decodeStats struct {
	unknownFieldNodes int
	rawCells          int
}

func (s *decodeStats) total() int {
	if s == nil {
		return 0
	}
	return s.unknownFieldNodes + s.rawCells
}

// messageHasUnknown reports whether a decoded protobuf message retains fields
// unknown to this SDK's schema.
func messageHasUnknown(m proto.Message) bool {
	return len(m.ProtoReflect().GetUnknown()) > 0
}

// captureUnknown copies a decoded message's unknown protobuf fields onto the
// model node so the encode path can re-emit them.
func captureUnknown(m proto.Message, dst *[]byte, stats *decodeStats) {
	raw := m.ProtoReflect().GetUnknown()
	if len(raw) == 0 {
		return
	}
	*dst = bytes.Clone(raw)
	if stats != nil {
		stats.unknownFieldNodes++
	}
}

// rulesContainerFromProto converts a protobuf RulesContainer to the model.
func rulesContainerFromProto(pbContainer *pb.RulesContainer) (*model.DecodedRulesContainer, error) {
	stats := &decodeStats{}
	container := &model.DecodedRulesContainer{
		MinimumDistinctUserSignatures:  int(pbContainer.GetMinimumDistinctUserSignatures()),
		MinimumDistinctGroupSignatures: int(pbContainer.GetMinimumDistinctGroupSignatures()),
		EnforcedRulesHash:              pbContainer.GetEnforcedRulesHash(),
		Timestamp:                      int64(pbContainer.GetTimestamp()),
		MinimumCommitmentSignatures:    int(pbContainer.GetMinimumCommitmentSignatures()),
		EngineIdentities:               pbContainer.GetEngineIdentities(),
		HsmSlotId:                      pbContainer.GetHsmSlotId(),
		Properties:                     pbContainer.GetProperties(),
	}

	// Convert users
	for _, u := range pbContainer.GetUsers() {
		user := userFromProto(u, stats)
		container.Users = append(container.Users, user)
	}

	// Convert groups
	for _, g := range pbContainer.GetGroups() {
		group := groupFromProto(g, stats)
		container.Groups = append(container.Groups, group)
	}

	// Convert address whitelisting rules
	for _, r := range pbContainer.GetAddressWhitelistingRules() {
		rule := addressWhitelistingRulesFromProto(r, stats)
		container.AddressWhitelistingRules = append(container.AddressWhitelistingRules, rule)
	}

	// Convert contract address whitelisting rules
	for _, r := range pbContainer.GetContractAddressWhitelistingRules() {
		rule := contractAddressWhitelistingRulesFromProto(r, stats)
		container.ContractAddressWhitelistingRules = append(container.ContractAddressWhitelistingRules, rule)
	}

	// Convert transaction rules
	for _, r := range pbContainer.GetTransactionRules() {
		rule := transactionRulesFromProto(r, stats)
		container.TransactionRules = append(container.TransactionRules, rule)
	}

	captureUnknown(pbContainer, &container.UnknownFields, stats)

	if stats.total() > 0 {
		slog.Warn("rules container carries data unknown to this SDK version; it is preserved and re-emitted verbatim on encode — consider upgrading the SDK before editing the container",
			"unknown_field_nodes", stats.unknownFieldNodes,
			"raw_cells", stats.rawCells,
		)
	}

	return container, nil
}

// userFromProto converts a protobuf User to the model.
func userFromProto(pbUser *pb.User, stats *decodeStats) *model.RuleUser {
	user := &model.RuleUser{
		ID:           pbUser.GetId(),
		PublicKeyPEM: pbUser.GetPublicKey(),
		Properties:   pbUser.GetProperties(),
	}

	// Convert roles
	for _, role := range pbUser.GetRoles() {
		user.Roles = append(user.Roles, role.String())
	}

	// Parse public key if present
	if user.PublicKeyPEM != "" {
		publicKey, err := crypto.DecodePublicKeyPEM(user.PublicKeyPEM)
		if err == nil {
			user.PublicKey = publicKey
		}
	}

	captureUnknown(pbUser, &user.UnknownFields, stats)
	return user
}

// groupFromProto converts a protobuf Group to the model.
func groupFromProto(pbGroup *pb.Group, stats *decodeStats) *model.RuleGroup {
	group := &model.RuleGroup{
		ID:         pbGroup.GetId(),
		UserIDs:    pbGroup.GetUserIds(),
		Properties: pbGroup.GetProperties(),
	}
	captureUnknown(pbGroup, &group.UnknownFields, stats)
	return group
}

// addressWhitelistingRulesFromProto converts a protobuf AddressWhitelistingRules to the model.
func addressWhitelistingRulesFromProto(pbRules *pb.RulesContainer_AddressWhitelistingRules, stats *decodeStats) *model.AddressWhitelistingRules {
	rules := &model.AddressWhitelistingRules{
		Currency:   pbRules.GetCurrency(),
		Network:    pbRules.GetNetwork(),
		Properties: pbRules.GetProperties(),
	}

	// Convert parallel thresholds
	for _, pt := range pbRules.GetParallelThresholds() {
		rules.ParallelThresholds = append(rules.ParallelThresholds, sequentialThresholdsFromProto(pt, stats))
	}

	// Convert lines
	for _, line := range pbRules.GetLines() {
		rules.Lines = append(rules.Lines, addressWhitelistingLineFromProto(line, stats))
	}

	captureUnknown(pbRules, &rules.UnknownFields, stats)
	return rules
}

// addressWhitelistingLineFromProto converts a protobuf AddressWhitelistingRules.Line to the model.
func addressWhitelistingLineFromProto(pbLine *pb.RulesContainer_AddressWhitelistingRules_Line, stats *decodeStats) *model.AddressWhitelistingLine {
	line := &model.AddressWhitelistingLine{
		Properties: pbLine.GetProperties(),
	}

	// Convert cells (each cell is a serialized RuleSource)
	for _, cellBytes := range pbLine.GetCells() {
		line.Cells = append(line.Cells, ruleSourceFromBytes(cellBytes, stats))
	}

	// Convert parallel thresholds
	for _, pt := range pbLine.GetParallelThresholds() {
		line.ParallelThresholds = append(line.ParallelThresholds, sequentialThresholdsFromProto(pt, stats))
	}

	captureUnknown(pbLine, &line.UnknownFields, stats)
	return line
}

// ruleSourceFromBytes decodes a RuleSource from serialized protobuf bytes.
// Sources this SDK version cannot fully represent (unknown type, unknown
// protobuf fields, undecodable payload) keep their exact wire bytes in Raw and
// are re-emitted verbatim on encode.
func ruleSourceFromBytes(data []byte, stats *decodeStats) *model.RuleSource {
	rawSource := func() *model.RuleSource {
		if stats != nil {
			stats.rawCells++
		}
		return &model.RuleSource{Raw: bytes.Clone(data)}
	}

	typed := ruleSourceFromBytesTyped(data)
	if typed == nil {
		return rawSource()
	}
	// Same lossless guard the cell codec uses, and for the same reason: an
	// unknown-field check does not catch a non-canonical encoding. Java already
	// byte-compares on this path (RulesContainerMapper.ruleSourceToBytes); Go,
	// Python and TypeScript did not.
	return losslessOrRaw(data, typed, ruleSourceToBytes, rawSource)
}

// ruleSourceFromBytesTyped decodes a RuleSource into its typed form, or returns nil
// when this SDK cannot represent it (unknown type, unknown protobuf field, or an
// undecodable payload). The caller turns nil into a verbatim raw source.
func ruleSourceFromBytesTyped(data []byte) *model.RuleSource {
	var pbSource pb.RuleSource
	if proto.Unmarshal(data, &pbSource) != nil || messageHasUnknown(&pbSource) {
		return nil
	}

	source := &model.RuleSource{
		Type: model.RuleSourceType(pbSource.GetType()),
	}
	payload := pbSource.GetPayload()

	// Decode payload based on type. Payload-less types tolerate a present-but-
	// empty payload; per the historical decode behavior an empty payload on a
	// payload-carrying type leaves its typed field nil.
	switch pbSource.GetType() {
	case pb.RuleSource_RuleSourceAny, pb.RuleSource_RuleSourceAnyExchange:
		if len(payload) > 0 {
			return nil
		}
	case pb.RuleSource_RuleSourceInternalWallet:
		if len(payload) > 0 {
			var m pb.RuleSourceInternalWallet
			if proto.Unmarshal(payload, &m) != nil || messageHasUnknown(&m) {
				return nil
			}
			source.InternalWallet = &model.RuleSourceInternalWallet{Path: m.GetPath()}
		}
	case pb.RuleSource_RuleSourceInternalAddress:
		if len(payload) > 0 {
			var m pb.RuleSourceInternalAddress
			if proto.Unmarshal(payload, &m) != nil || messageHasUnknown(&m) {
				return nil
			}
			source.InternalAddress = &model.RuleSourceInternalAddress{Address: m.GetAddress(), Path: m.GetPath()}
		}
	case pb.RuleSource_RuleSourceExchange:
		if len(payload) > 0 {
			var m pb.RuleSourceExchange
			if proto.Unmarshal(payload, &m) != nil || messageHasUnknown(&m) {
				return nil
			}
			source.Exchange = &model.RuleSourceExchange{Label: m.GetLabel()}
		}
	case pb.RuleSource_RuleSourceExternalAddress:
		if len(payload) > 0 {
			var m pb.RuleSourceExternalAddress
			if proto.Unmarshal(payload, &m) != nil || messageHasUnknown(&m) {
				return nil
			}
			source.ExternalAddress = &model.RuleSourceExternalAddress{Address: m.GetAddress(), Memo: m.GetMemo()}
		}
	default:
		// Source type newer than this SDK: preserve verbatim.
		return nil
	}

	return source
}

// contractAddressWhitelistingRulesFromProto converts a protobuf ContractAddressWhitelistingRules to the model.
func contractAddressWhitelistingRulesFromProto(pbRules *pb.RulesContainer_ContractAddressWhitelistingRules, stats *decodeStats) *model.ContractAddressWhitelistingRules {
	rules := &model.ContractAddressWhitelistingRules{
		Blockchain: pbRules.GetBlockchain().String(),
		Network:    pbRules.GetNetwork(),
		Properties: pbRules.GetProperties(),
	}

	// Convert parallel thresholds
	for _, pt := range pbRules.GetParallelThresholds() {
		rules.ParallelThresholds = append(rules.ParallelThresholds, sequentialThresholdsFromProto(pt, stats))
	}

	captureUnknown(pbRules, &rules.UnknownFields, stats)
	return rules
}

// sequentialThresholdsFromProto converts a protobuf SequentialThresholds to the model.
func sequentialThresholdsFromProto(pbST *pb.SequentialThresholds, stats *decodeStats) *model.SequentialThresholds {
	st := &model.SequentialThresholds{}

	for _, t := range pbST.GetThresholds() {
		st.Thresholds = append(st.Thresholds, groupThresholdFromProto(t, stats))
	}

	captureUnknown(pbST, &st.UnknownFields, stats)
	return st
}

// groupThresholdFromProto converts a protobuf GroupThreshold to the model.
func groupThresholdFromProto(pbT *pb.GroupThreshold, stats *decodeStats) *model.GroupThreshold {
	threshold := &model.GroupThreshold{
		GroupID:           pbT.GetGroupId(),
		MinimumSignatures: int(pbT.GetMinimumSignatures()),
	}
	captureUnknown(pbT, &threshold.UnknownFields, stats)
	return threshold
}

// transactionRulesFromProto converts a protobuf TransactionRules to the model.
func transactionRulesFromProto(pbRules *pb.RulesContainer_TransactionRules, stats *decodeStats) *model.TransactionRules {
	rules := &model.TransactionRules{
		Key: pbRules.GetKey(),
	}

	// Convert columns
	for _, c := range pbRules.GetColumns() {
		column := &model.RuleColumn{
			Type:        c.GetType().String(),
			Name:        c.GetName(),
			MetadataKey: c.GetMetadataKey(),
		}
		captureUnknown(c, &column.UnknownFields, stats)
		rules.Columns = append(rules.Columns, column)
	}

	// Convert lines; cells decode against the column they align with.
	for _, l := range pbRules.GetLines() {
		line := &model.RuleLine{
			Priority:   l.GetPriority(),
			Properties: l.GetProperties(),
		}
		for i, cell := range l.GetCells() {
			colType := ""
			if i < len(rules.Columns) {
				colType = rules.Columns[i].Type
			}
			line.Cells = append(line.Cells, ruleCellFromBytes(colType, cell, stats))
		}
		for _, pt := range l.GetParallelThresholds() {
			line.ParallelThresholds = append(line.ParallelThresholds, sequentialThresholdsFromProto(pt, stats))
		}
		captureUnknown(l, &line.UnknownFields, stats)
		rules.Lines = append(rules.Lines, line)
	}

	// Convert details if present
	if pbDetails := pbRules.GetDetails(); pbDetails != nil {
		details := &model.TransactionRuleDetails{
			Domain:     pbDetails.GetDomain().String(),
			SubDomain:  pbDetails.GetSubDomain().String(),
			Blockchain: pbDetails.GetBlockchain(),
			Network:    pbDetails.GetNetwork(),
		}
		if evm := pbDetails.GetEvmCallContract(); evm != nil {
			details.EvmCallContract = &model.EvmCallContract{
				ContractType:    evm.GetContractType(),
				MethodSignature: evm.GetMethodSignature(),
			}
			captureUnknown(evm, &details.EvmCallContract.UnknownFields, stats)
		}
		if xtz := pbDetails.GetXtzCallContract(); xtz != nil {
			details.XtzCallContract = &model.XtzCallContract{
				ContractType:    xtz.GetContractType(),
				MethodSignature: xtz.GetMethodSignature(),
			}
			captureUnknown(xtz, &details.XtzCallContract.UnknownFields, stats)
		}
		if cash := pbDetails.GetCashSettlement(); cash != nil {
			details.CashSettlement = &model.CashSettlement{
				Provider:    cash.GetProvider(),
				RequestType: cash.GetRequestType(),
			}
			captureUnknown(cash, &details.CashSettlement.UnknownFields, stats)
		}
		if cosmos := pbDetails.GetCosmosDetails(); cosmos != nil {
			details.CosmosDetails = &model.CosmosDetails{
				MethodSignatures: cosmos.GetMethodSignatures(),
			}
			captureUnknown(cosmos, &details.CosmosDetails.UnknownFields, stats)
		}
		captureUnknown(pbDetails, &details.UnknownFields, stats)
		rules.Details = details
	}

	captureUnknown(pbRules, &rules.UnknownFields, stats)
	return rules
}

// UserSignaturesFromBase64 decodes base64-encoded protobuf UserSignatures into model.
func UserSignaturesFromBase64(base64Data string) ([]*model.RuleUserSignature, error) {
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	return UserSignaturesFromBytes(data)
}

// UserSignaturesFromBytes decodes raw protobuf bytes into RuleUserSignature slice.
func UserSignaturesFromBytes(data []byte) ([]*model.RuleUserSignature, error) {
	var pbSigs pb.UserSignatures
	if err := proto.Unmarshal(data, &pbSigs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	var signatures []*model.RuleUserSignature
	for _, sig := range pbSigs.GetSignatures() {
		signatures = append(signatures, &model.RuleUserSignature{
			UserID:    sig.GetUserId(),
			Signature: base64.StdEncoding.EncodeToString(sig.GetSignature()),
		})
	}

	return signatures, nil
}
