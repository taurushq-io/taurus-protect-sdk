package mapper

import (
	"encoding/base64"
	"fmt"

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
func RulesContainerFromBytes(data []byte) (*model.DecodedRulesContainer, error) {
	var pbContainer pb.RulesContainer
	if err := proto.Unmarshal(data, &pbContainer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return rulesContainerFromProto(&pbContainer)
}

// rulesContainerFromProto converts a protobuf RulesContainer to the model.
func rulesContainerFromProto(pb *pb.RulesContainer) (*model.DecodedRulesContainer, error) {
	container := &model.DecodedRulesContainer{
		MinimumDistinctUserSignatures:  int(pb.GetMinimumDistinctUserSignatures()),
		MinimumDistinctGroupSignatures: int(pb.GetMinimumDistinctGroupSignatures()),
		EnforcedRulesHash:              pb.GetEnforcedRulesHash(),
		Timestamp:                      int64(pb.GetTimestamp()),
		MinimumCommitmentSignatures:    int(pb.GetMinimumCommitmentSignatures()),
		EngineIdentities:               pb.GetEngineIdentities(),
		HsmSlotId:                      pb.GetHsmSlotId(),
	}

	// Convert users
	for _, u := range pb.GetUsers() {
		user := userFromProto(u)
		container.Users = append(container.Users, user)
	}

	// Convert groups
	for _, g := range pb.GetGroups() {
		group := groupFromProto(g)
		container.Groups = append(container.Groups, group)
	}

	// Convert address whitelisting rules
	for _, r := range pb.GetAddressWhitelistingRules() {
		rule := addressWhitelistingRulesFromProto(r)
		container.AddressWhitelistingRules = append(container.AddressWhitelistingRules, rule)
	}

	// Convert contract address whitelisting rules
	for _, r := range pb.GetContractAddressWhitelistingRules() {
		rule := contractAddressWhitelistingRulesFromProto(r)
		container.ContractAddressWhitelistingRules = append(container.ContractAddressWhitelistingRules, rule)
	}

	// Convert transaction rules
	for _, r := range pb.GetTransactionRules() {
		rule := transactionRulesFromProto(r)
		container.TransactionRules = append(container.TransactionRules, rule)
	}

	return container, nil
}

// userFromProto converts a protobuf User to the model.
func userFromProto(pb *pb.User) *model.RuleUser {
	user := &model.RuleUser{
		ID:           pb.GetId(),
		PublicKeyPEM: pb.GetPublicKey(),
	}

	// Convert roles
	for _, role := range pb.GetRoles() {
		user.Roles = append(user.Roles, role.String())
	}

	// Parse public key if present
	if user.PublicKeyPEM != "" {
		publicKey, err := crypto.DecodePublicKeyPEM(user.PublicKeyPEM)
		if err == nil {
			user.PublicKey = publicKey
		}
	}

	return user
}

// groupFromProto converts a protobuf Group to the model.
func groupFromProto(pb *pb.Group) *model.RuleGroup {
	return &model.RuleGroup{
		ID:      pb.GetId(),
		UserIDs: pb.GetUserIds(),
	}
}

// addressWhitelistingRulesFromProto converts a protobuf AddressWhitelistingRules to the model.
func addressWhitelistingRulesFromProto(pb *pb.RulesContainer_AddressWhitelistingRules) *model.AddressWhitelistingRules {
	rules := &model.AddressWhitelistingRules{
		Currency: pb.GetCurrency(),
		Network:  pb.GetNetwork(),
	}

	// Convert parallel thresholds
	for _, pt := range pb.GetParallelThresholds() {
		rules.ParallelThresholds = append(rules.ParallelThresholds, sequentialThresholdsFromProto(pt))
	}

	// Convert lines
	for _, line := range pb.GetLines() {
		rules.Lines = append(rules.Lines, addressWhitelistingLineFromProto(line))
	}

	return rules
}

// addressWhitelistingLineFromProto converts a protobuf AddressWhitelistingRules.Line to the model.
func addressWhitelistingLineFromProto(pb *pb.RulesContainer_AddressWhitelistingRules_Line) *model.AddressWhitelistingLine {
	line := &model.AddressWhitelistingLine{}

	// Convert cells (each cell is a serialized RuleSource)
	for _, cellBytes := range pb.GetCells() {
		source := ruleSourceFromBytes(cellBytes)
		if source != nil {
			line.Cells = append(line.Cells, source)
		}
	}

	// Convert parallel thresholds
	for _, pt := range pb.GetParallelThresholds() {
		line.ParallelThresholds = append(line.ParallelThresholds, sequentialThresholdsFromProto(pt))
	}

	return line
}

// ruleSourceFromBytes decodes a RuleSource from serialized protobuf bytes.
func ruleSourceFromBytes(data []byte) *model.RuleSource {
	var pbSource pb.RuleSource
	if err := proto.Unmarshal(data, &pbSource); err != nil {
		return nil
	}

	source := &model.RuleSource{
		Type: model.RuleSourceType(pbSource.GetType()),
	}

	// Decode payload based on type
	if source.Type == model.RuleSourceTypeInternalWallet && len(pbSource.GetPayload()) > 0 {
		var pbWallet pb.RuleSourceInternalWallet
		if err := proto.Unmarshal(pbSource.GetPayload(), &pbWallet); err == nil {
			source.InternalWallet = &model.RuleSourceInternalWallet{
				Path: pbWallet.GetPath(),
			}
		}
	}

	return source
}

// contractAddressWhitelistingRulesFromProto converts a protobuf ContractAddressWhitelistingRules to the model.
func contractAddressWhitelistingRulesFromProto(pb *pb.RulesContainer_ContractAddressWhitelistingRules) *model.ContractAddressWhitelistingRules {
	rules := &model.ContractAddressWhitelistingRules{
		Blockchain: pb.GetBlockchain().String(),
		Network:    pb.GetNetwork(),
	}

	// Convert parallel thresholds
	for _, pt := range pb.GetParallelThresholds() {
		rules.ParallelThresholds = append(rules.ParallelThresholds, sequentialThresholdsFromProto(pt))
	}

	return rules
}

// sequentialThresholdsFromProto converts a protobuf SequentialThresholds to the model.
func sequentialThresholdsFromProto(pb *pb.SequentialThresholds) *model.SequentialThresholds {
	st := &model.SequentialThresholds{}

	for _, t := range pb.GetThresholds() {
		st.Thresholds = append(st.Thresholds, groupThresholdFromProto(t))
	}

	return st
}

// groupThresholdFromProto converts a protobuf GroupThreshold to the model.
func groupThresholdFromProto(pb *pb.GroupThreshold) *model.GroupThreshold {
	return &model.GroupThreshold{
		GroupID:           pb.GetGroupId(),
		MinimumSignatures: int(pb.GetMinimumSignatures()),
	}
}

// transactionRulesFromProto converts a protobuf TransactionRules to the model.
func transactionRulesFromProto(pb *pb.RulesContainer_TransactionRules) *model.TransactionRules {
	rules := &model.TransactionRules{
		Key: pb.GetKey(),
	}

	// Convert columns
	for _, c := range pb.GetColumns() {
		rules.Columns = append(rules.Columns, &model.RuleColumn{
			Type: c.GetType().String(),
		})
	}

	// Convert lines
	for _, l := range pb.GetLines() {
		line := &model.RuleLine{}
		for _, cell := range l.GetCells() {
			line.Cells = append(line.Cells, string(cell))
		}
		for _, pt := range l.GetParallelThresholds() {
			line.ParallelThresholds = append(line.ParallelThresholds, sequentialThresholdsFromProto(pt))
		}
		rules.Lines = append(rules.Lines, line)
	}

	// Convert details if present
	if pb.GetDetails() != nil {
		rules.Details = &model.TransactionRuleDetails{
			Domain:    pb.GetDetails().GetDomain().String(),
			SubDomain: pb.GetDetails().GetSubDomain().String(),
		}
	}

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
