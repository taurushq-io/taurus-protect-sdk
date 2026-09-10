package mapper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math"
	"strconv"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// deterministicMarshal keeps map ordering stable so encoding the same
// container twice yields identical bytes.
var deterministicMarshal = proto.MarshalOptions{Deterministic: true}

// RulesContainerToBase64 encodes a DecodedRulesContainer to the base64
// protobuf wire format the governance API expects. Inverse of
// RulesContainerFromBase64.
//
// The server-controlled EnforcedRulesHash and Timestamp fields are stripped:
// the server recomputes them and rejects submissions asserting stale values.
// Unknown protobuf fields, raw cells, and raw whitelisting sources captured at
// decode are re-emitted verbatim, so containers produced by newer schema
// versions round-trip without data loss.
func RulesContainerToBase64(container *model.DecodedRulesContainer) (string, error) {
	data, err := RulesContainerToBytes(container)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// RulesContainerToBytes encodes a DecodedRulesContainer to raw protobuf bytes.
// Inverse of RulesContainerFromBytes; see RulesContainerToBase64 for the
// stripping and unknown-field semantics.
func RulesContainerToBytes(container *model.DecodedRulesContainer) ([]byte, error) {
	pbContainer, err := rulesContainerToProto(container)
	if err != nil {
		return nil, err
	}
	data, err := deterministicMarshal.Marshal(pbContainer)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal rules container: %w", err)
	}
	return data, nil
}

// rulesContainerToProto converts the model to a protobuf RulesContainer.
// Mirror of rulesContainerFromProto.
func rulesContainerToProto(c *model.DecodedRulesContainer) (*pb.RulesContainer, error) {
	if c == nil {
		return nil, fmt.Errorf("rules container cannot be nil")
	}

	minUserSigs, err := uint32FromInt(c.MinimumDistinctUserSignatures, "MinimumDistinctUserSignatures")
	if err != nil {
		return nil, err
	}
	minGroupSigs, err := uint32FromInt(c.MinimumDistinctGroupSignatures, "MinimumDistinctGroupSignatures")
	if err != nil {
		return nil, err
	}
	minCommitSigs, err := uint32FromInt(c.MinimumCommitmentSignatures, "MinimumCommitmentSignatures")
	if err != nil {
		return nil, err
	}

	out := &pb.RulesContainer{
		MinimumDistinctUserSignatures:  minUserSigs,
		MinimumDistinctGroupSignatures: minGroupSigs,
		MinimumCommitmentSignatures:    minCommitSigs,
		EngineIdentities:               c.EngineIdentities,
		HsmSlotId:                      c.HsmSlotId,
		Properties:                     c.Properties,
		// EnforcedRulesHash and Timestamp are deliberately not set: they are
		// server-controlled and must be stripped from proposal submissions.
	}

	for _, u := range c.Users {
		pbUser, err := userToProto(u)
		if err != nil {
			return nil, err
		}
		out.Users = append(out.Users, pbUser)
	}

	for _, g := range c.Groups {
		out.Groups = append(out.Groups, groupToProto(g))
	}

	for _, r := range c.TransactionRules {
		pbRule, err := transactionRulesToProto(r)
		if err != nil {
			return nil, err
		}
		out.TransactionRules = append(out.TransactionRules, pbRule)
	}

	for _, r := range c.AddressWhitelistingRules {
		pbRule, err := addressWhitelistingRulesToProto(r)
		if err != nil {
			return nil, err
		}
		out.AddressWhitelistingRules = append(out.AddressWhitelistingRules, pbRule)
	}

	for _, r := range c.ContractAddressWhitelistingRules {
		pbRule, err := contractAddressWhitelistingRulesToProto(r)
		if err != nil {
			return nil, err
		}
		out.ContractAddressWhitelistingRules = append(out.ContractAddressWhitelistingRules, pbRule)
	}

	attachUnknown(out, c.UnknownFields)
	return out, nil
}

// userToProto converts a model RuleUser to protobuf. Mirror of userFromProto.
func userToProto(u *model.RuleUser) (*pb.User, error) {
	if u == nil {
		return nil, fmt.Errorf("rule user cannot be nil")
	}
	out := &pb.User{
		Id:         u.ID,
		PublicKey:  u.PublicKeyPEM,
		Properties: u.Properties,
	}
	for _, role := range u.Roles {
		n, err := enumNumber(pb.Role_value, role, "user role")
		if err != nil {
			return nil, fmt.Errorf("user %q: %w", u.ID, err)
		}
		out.Roles = append(out.Roles, pb.Role(n))
	}
	attachUnknown(out, u.UnknownFields)
	return out, nil
}

// groupToProto converts a model RuleGroup to protobuf. Mirror of groupFromProto.
func groupToProto(g *model.RuleGroup) *pb.Group {
	if g == nil {
		return nil
	}
	out := &pb.Group{
		Id:         g.ID,
		UserIds:    g.UserIDs,
		Properties: g.Properties,
	}
	attachUnknown(out, g.UnknownFields)
	return out
}

// transactionRulesToProto converts a model TransactionRules to protobuf.
// Mirror of transactionRulesFromProto.
func transactionRulesToProto(r *model.TransactionRules) (*pb.RulesContainer_TransactionRules, error) {
	if r == nil {
		return nil, fmt.Errorf("transaction rule cannot be nil")
	}
	out := &pb.RulesContainer_TransactionRules{Key: r.Key}

	for _, col := range r.Columns {
		if col == nil {
			return nil, fmt.Errorf("transaction rule %q: column cannot be nil", r.Key)
		}
		n, err := enumNumber(pb.RulesContainer_ColumnType_value, col.Type, "column type")
		if err != nil {
			return nil, fmt.Errorf("transaction rule %q: %w", r.Key, err)
		}
		pbCol := &pb.RulesContainer_Column{
			Type:        pb.RulesContainer_ColumnType(n),
			Name:        col.Name,
			MetadataKey: col.MetadataKey,
		}
		attachUnknown(pbCol, col.UnknownFields)
		out.Columns = append(out.Columns, pbCol)
	}

	for lineIndex, line := range r.Lines {
		if line == nil {
			return nil, fmt.Errorf("transaction rule %q: line %d cannot be nil", r.Key, lineIndex)
		}
		pbLine := &pb.RulesContainer_Line{
			Priority:   line.Priority,
			Properties: line.Properties,
		}
		for cellIndex, cell := range line.Cells {
			colType := ""
			if cellIndex < len(r.Columns) && r.Columns[cellIndex] != nil {
				colType = r.Columns[cellIndex].Type
			}
			cellBytes, err := ruleCellToBytes(colType, cell)
			if err != nil {
				return nil, fmt.Errorf("transaction rule %q line %d cell %d: %w", r.Key, lineIndex, cellIndex, err)
			}
			pbLine.Cells = append(pbLine.Cells, cellBytes)
		}
		thresholds, err := sequentialThresholdsListToProto(line.ParallelThresholds)
		if err != nil {
			return nil, fmt.Errorf("transaction rule %q line %d: %w", r.Key, lineIndex, err)
		}
		pbLine.ParallelThresholds = thresholds
		attachUnknown(pbLine, line.UnknownFields)
		out.Lines = append(out.Lines, pbLine)
	}

	if r.Details != nil {
		pbDetails, err := transactionRuleDetailsToProto(r.Details)
		if err != nil {
			return nil, fmt.Errorf("transaction rule %q: %w", r.Key, err)
		}
		out.Details = pbDetails
	}

	attachUnknown(out, r.UnknownFields)
	return out, nil
}

// transactionRuleDetailsToProto converts model rule details to protobuf.
func transactionRuleDetailsToProto(d *model.TransactionRuleDetails) (*pb.RulesContainer_TransactionRules_TransactionRuleDetails, error) {
	domain, err := enumNumber(pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleDomain_value, d.Domain, "rule domain")
	if err != nil {
		return nil, err
	}
	subDomain, err := enumNumber(pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleSubDomain_value, d.SubDomain, "rule sub-domain")
	if err != nil {
		return nil, err
	}
	out := &pb.RulesContainer_TransactionRules_TransactionRuleDetails{
		Domain:     pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleDomain(domain),
		SubDomain:  pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleSubDomain(subDomain),
		Blockchain: d.Blockchain,
		Network:    d.Network,
	}
	if d.EvmCallContract != nil {
		out.EvmCallContract = &pb.RulesContainer_TransactionRules_TransactionRuleDetails_EvmCallContract{
			ContractType:    d.EvmCallContract.ContractType,
			MethodSignature: d.EvmCallContract.MethodSignature,
		}
		attachUnknown(out.EvmCallContract, d.EvmCallContract.UnknownFields)
	}
	if d.XtzCallContract != nil {
		out.XtzCallContract = &pb.RulesContainer_TransactionRules_TransactionRuleDetails_XtzCallContract{
			ContractType:    d.XtzCallContract.ContractType,
			MethodSignature: d.XtzCallContract.MethodSignature,
		}
		attachUnknown(out.XtzCallContract, d.XtzCallContract.UnknownFields)
	}
	if d.CashSettlement != nil {
		out.CashSettlement = &pb.RulesContainer_TransactionRules_TransactionRuleDetails_CashSettlement{
			Provider:    d.CashSettlement.Provider,
			RequestType: d.CashSettlement.RequestType,
		}
		attachUnknown(out.CashSettlement, d.CashSettlement.UnknownFields)
	}
	if d.CosmosDetails != nil {
		out.CosmosDetails = &pb.RulesContainer_TransactionRules_TransactionRuleDetails_CosmosDetails{
			MethodSignatures: d.CosmosDetails.MethodSignatures,
		}
		attachUnknown(out.CosmosDetails, d.CosmosDetails.UnknownFields)
	}
	attachUnknown(out, d.UnknownFields)
	return out, nil
}

// addressWhitelistingRulesToProto converts model address whitelisting rules to
// protobuf. Mirror of addressWhitelistingRulesFromProto.
func addressWhitelistingRulesToProto(r *model.AddressWhitelistingRules) (*pb.RulesContainer_AddressWhitelistingRules, error) {
	if r == nil {
		return nil, fmt.Errorf("address whitelisting rule cannot be nil")
	}
	out := &pb.RulesContainer_AddressWhitelistingRules{
		Currency:   r.Currency,
		Network:    r.Network,
		Properties: r.Properties,
	}
	thresholds, err := sequentialThresholdsListToProto(r.ParallelThresholds)
	if err != nil {
		return nil, fmt.Errorf("address whitelisting rule %q: %w", r.Currency, err)
	}
	out.ParallelThresholds = thresholds

	for lineIndex, line := range r.Lines {
		if line == nil {
			return nil, fmt.Errorf("address whitelisting rule %q: line %d cannot be nil", r.Currency, lineIndex)
		}
		pbLine := &pb.RulesContainer_AddressWhitelistingRules_Line{
			Properties: line.Properties,
		}
		for cellIndex, source := range line.Cells {
			cellBytes, err := ruleSourceToBytes(source)
			if err != nil {
				return nil, fmt.Errorf("address whitelisting rule %q line %d cell %d: %w", r.Currency, lineIndex, cellIndex, err)
			}
			pbLine.Cells = append(pbLine.Cells, cellBytes)
		}
		lineThresholds, err := sequentialThresholdsListToProto(line.ParallelThresholds)
		if err != nil {
			return nil, fmt.Errorf("address whitelisting rule %q line %d: %w", r.Currency, lineIndex, err)
		}
		pbLine.ParallelThresholds = lineThresholds
		attachUnknown(pbLine, line.UnknownFields)
		out.Lines = append(out.Lines, pbLine)
	}

	attachUnknown(out, r.UnknownFields)
	return out, nil
}

// ruleSourceToBytes encodes a whitelisting rule source cell. Mirror of
// ruleSourceFromBytes; a Raw source is emitted verbatim.
func ruleSourceToBytes(s *model.RuleSource) ([]byte, error) {
	if s == nil {
		return nil, nil
	}
	if len(s.Raw) > 0 {
		return s.Raw, nil
	}

	out := &pb.RuleSource{Type: pb.RuleSource_RuleSourceType(s.Type)}
	var inner proto.Message
	switch s.Type {
	case model.RuleSourceTypeAny, model.RuleSourceTypeAnyExchange:
		// No payload.
	case model.RuleSourceTypeInternalWallet:
		if s.InternalWallet != nil {
			inner = &pb.RuleSourceInternalWallet{Path: s.InternalWallet.Path}
		}
	case model.RuleSourceTypeInternalAddress:
		if s.InternalAddress != nil {
			inner = &pb.RuleSourceInternalAddress{Address: s.InternalAddress.Address, Path: s.InternalAddress.Path}
		}
	case model.RuleSourceTypeExchange:
		if s.Exchange != nil {
			inner = &pb.RuleSourceExchange{Label: s.Exchange.Label}
		}
	case model.RuleSourceTypeExternalAddress:
		if s.ExternalAddress != nil {
			inner = &pb.RuleSourceExternalAddress{Address: s.ExternalAddress.Address, Memo: s.ExternalAddress.Memo}
		}
	default:
		return nil, fmt.Errorf("unknown rule source type %d", s.Type)
	}
	if inner != nil {
		payload, err := marshalCell(inner)
		if err != nil {
			return nil, err
		}
		out.Payload = payload
	}
	return marshalCell(out)
}

// contractAddressWhitelistingRulesToProto converts model contract whitelisting
// rules to protobuf. Mirror of contractAddressWhitelistingRulesFromProto.
func contractAddressWhitelistingRulesToProto(r *model.ContractAddressWhitelistingRules) (*pb.RulesContainer_ContractAddressWhitelistingRules, error) {
	if r == nil {
		return nil, fmt.Errorf("contract address whitelisting rule cannot be nil")
	}
	blockchain, err := blockchainToProto(r.Blockchain)
	if err != nil {
		return nil, fmt.Errorf("contract address whitelisting rule: %w", err)
	}
	out := &pb.RulesContainer_ContractAddressWhitelistingRules{
		Blockchain: blockchain,
		Network:    r.Network,
		Properties: r.Properties,
	}
	thresholds, err := sequentialThresholdsListToProto(r.ParallelThresholds)
	if err != nil {
		return nil, fmt.Errorf("contract address whitelisting rule %q: %w", r.Blockchain, err)
	}
	out.ParallelThresholds = thresholds
	attachUnknown(out, r.UnknownFields)
	return out, nil
}

// sequentialThresholdsListToProto converts model thresholds to protobuf.
func sequentialThresholdsListToProto(thresholds []*model.SequentialThresholds) ([]*pb.SequentialThresholds, error) {
	var out []*pb.SequentialThresholds
	for _, st := range thresholds {
		if st == nil {
			return nil, fmt.Errorf("sequential thresholds cannot be nil")
		}
		pbST := &pb.SequentialThresholds{}
		for _, t := range st.Thresholds {
			if t == nil {
				return nil, fmt.Errorf("group threshold cannot be nil")
			}
			minSigs, err := uint32FromInt(t.MinimumSignatures, "MinimumSignatures")
			if err != nil {
				return nil, fmt.Errorf("group %q: %w", t.GroupID, err)
			}
			pbT := &pb.GroupThreshold{
				GroupId:           t.GroupID,
				MinimumSignatures: minSigs,
			}
			attachUnknown(pbT, t.UnknownFields)
			pbST.Thresholds = append(pbST.Thresholds, pbT)
		}
		attachUnknown(pbST, st.UnknownFields)
		out = append(out, pbST)
	}
	return out, nil
}

// blockchainToProto reverse-maps a blockchain enum value name.
func blockchainToProto(name string) (pb.Blockchain, error) {
	n, err := enumNumber(pb.Blockchain_value, name, "blockchain")
	if err != nil {
		return 0, err
	}
	return pb.Blockchain(n), nil
}

// enumNumber reverse-maps an enum value name to its number. An empty name maps
// to the enum zero value. Decimal strings pass through numerically, so enum
// values decoded from a newer schema than this SDK re-encode losslessly; any
// other unrecognized name is an error (a caller-authored typo, not skew).
func enumNumber(valueMap map[string]int32, name, what string) (int32, error) {
	if name == "" {
		return 0, nil
	}
	if v, ok := valueMap[name]; ok {
		return v, nil
	}
	if n, err := strconv.ParseInt(name, 10, 32); err == nil {
		return int32(n), nil
	}
	return 0, fmt.Errorf("unknown %s %q", what, name)
}

// uint32FromInt guards int → uint32 conversions in the encode path. The upper bound is
// compared in uint64 because `n > math.MaxUint32` does not compile where int is 32 bits.
func uint32FromInt(n int, what string) (uint32, error) {
	if n < 0 || uint64(n) > math.MaxUint32 {
		return 0, fmt.Errorf("%s out of range: %d", what, n)
	}
	return uint32(n), nil
}

// attachUnknown re-attaches unknown protobuf fields captured at decode so
// proto.Marshal re-emits them.
func attachUnknown(m proto.Message, raw []byte) {
	if len(raw) > 0 {
		m.ProtoReflect().SetUnknown(protoreflect.RawFields(bytes.Clone(raw)))
	}
}
