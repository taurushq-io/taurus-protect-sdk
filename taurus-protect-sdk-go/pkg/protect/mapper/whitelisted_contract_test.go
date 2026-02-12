package mapper

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestWhitelistedContractFromDTO_Nil(t *testing.T) {
	result := WhitelistedContractFromDTO(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestWhitelistedContractFromDTO_Basic(t *testing.T) {
	id := "123"
	tenantID := "tenant-1"
	status := "approved"
	action := "create"
	blockchain := "ETH"
	network := "mainnet"
	rule := "rule-1"
	rulesContainer := "container-1"
	rulesSig := "signatures-1"
	businessRuleEnabled := true

	dto := &openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope{
		Id:                  &id,
		TenantId:            &tenantID,
		Status:              &status,
		Action:              &action,
		Blockchain:          &blockchain,
		Network:             &network,
		Rule:                &rule,
		RulesContainer:      &rulesContainer,
		RulesSignatures:     &rulesSig,
		BusinessRuleEnabled: &businessRuleEnabled,
	}

	result := WhitelistedContractFromDTO(dto)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.ID != id {
		t.Errorf("expected ID %s, got %s", id, result.ID)
	}
	if result.TenantID != tenantID {
		t.Errorf("expected TenantID %s, got %s", tenantID, result.TenantID)
	}
	if result.Status != status {
		t.Errorf("expected Status %s, got %s", status, result.Status)
	}
	if result.Action != action {
		t.Errorf("expected Action %s, got %s", action, result.Action)
	}
	if result.Blockchain != blockchain {
		t.Errorf("expected Blockchain %s, got %s", blockchain, result.Blockchain)
	}
	if result.Network != network {
		t.Errorf("expected Network %s, got %s", network, result.Network)
	}
	if result.Rule != rule {
		t.Errorf("expected Rule %s, got %s", rule, result.Rule)
	}
	if result.RulesContainer != rulesContainer {
		t.Errorf("expected RulesContainer %s, got %s", rulesContainer, result.RulesContainer)
	}
	if result.RulesSignatures != rulesSig {
		t.Errorf("expected RulesSignatures %s, got %s", rulesSig, result.RulesSignatures)
	}
	if result.BusinessRuleEnabled != businessRuleEnabled {
		t.Errorf("expected BusinessRuleEnabled %v, got %v", businessRuleEnabled, result.BusinessRuleEnabled)
	}
}

func TestWhitelistedContractFromDTO_WithAttributes(t *testing.T) {
	id := "123"
	attrID := "attr-1"
	attrKey := "key1"
	attrValue := "value1"
	attrContentType := "text/plain"
	isFile := true

	dto := &openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope{
		Id: &id,
		Attributes: []openapi.TgvalidatordWhitelistedContractAddressAttribute{
			{
				Id:          &attrID,
				Key:         &attrKey,
				Value:       &attrValue,
				ContentType: &attrContentType,
				Isfile:      &isFile,
			},
		},
	}

	result := WhitelistedContractFromDTO(dto)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.Attributes) != 1 {
		t.Fatalf("expected 1 attribute, got %d", len(result.Attributes))
	}
	attr := result.Attributes[0]
	if attr.ID != attrID {
		t.Errorf("expected attribute ID %s, got %s", attrID, attr.ID)
	}
	if attr.Key != attrKey {
		t.Errorf("expected attribute Key %s, got %s", attrKey, attr.Key)
	}
	if attr.Value != attrValue {
		t.Errorf("expected attribute Value %s, got %s", attrValue, attr.Value)
	}
	if attr.ContentType != attrContentType {
		t.Errorf("expected attribute ContentType %s, got %s", attrContentType, attr.ContentType)
	}
	if attr.IsFile != isFile {
		t.Errorf("expected attribute IsFile %v, got %v", isFile, attr.IsFile)
	}
}

func TestWhitelistedContractsFromDTO_Nil(t *testing.T) {
	result := WhitelistedContractsFromDTO(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestWhitelistedContractsFromDTO_Multiple(t *testing.T) {
	id1 := "123"
	id2 := "456"

	dtos := []openapi.TgvalidatordSignedWhitelistedContractAddressEnvelope{
		{Id: &id1},
		{Id: &id2},
	}

	result := WhitelistedContractsFromDTO(dtos)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 contracts, got %d", len(result))
	}
	if result[0].ID != id1 {
		t.Errorf("expected first contract ID %s, got %s", id1, result[0].ID)
	}
	if result[1].ID != id2 {
		t.Errorf("expected second contract ID %s, got %s", id2, result[1].ID)
	}
}

func TestWhitelistedContractAttributeFromDTO_Nil(t *testing.T) {
	result := WhitelistedContractAttributeFromDTO(nil)
	if result.ID != "" || result.Key != "" || result.Value != "" {
		t.Errorf("expected empty attribute, got %+v", result)
	}
}

func TestWhitelistedContractAttributeFromDTO_Full(t *testing.T) {
	id := "attr-1"
	key := "test-key"
	value := "test-value"
	contentType := "application/json"
	owner := "user-1"
	attrType := "custom"
	subtype := "metadata"
	isFile := false

	dto := &openapi.TgvalidatordWhitelistedContractAddressAttribute{
		Id:          &id,
		Key:         &key,
		Value:       &value,
		ContentType: &contentType,
		Owner:       &owner,
		Type:        &attrType,
		Subtype:     &subtype,
		Isfile:      &isFile,
	}

	result := WhitelistedContractAttributeFromDTO(dto)

	if result.ID != id {
		t.Errorf("expected ID %s, got %s", id, result.ID)
	}
	if result.Key != key {
		t.Errorf("expected Key %s, got %s", key, result.Key)
	}
	if result.Value != value {
		t.Errorf("expected Value %s, got %s", value, result.Value)
	}
	if result.ContentType != contentType {
		t.Errorf("expected ContentType %s, got %s", contentType, result.ContentType)
	}
	if result.Owner != owner {
		t.Errorf("expected Owner %s, got %s", owner, result.Owner)
	}
	if result.Type != attrType {
		t.Errorf("expected Type %s, got %s", attrType, result.Type)
	}
	if result.Subtype != subtype {
		t.Errorf("expected Subtype %s, got %s", subtype, result.Subtype)
	}
	if result.IsFile != isFile {
		t.Errorf("expected IsFile %v, got %v", isFile, result.IsFile)
	}
}

func TestWhitelistedContractAttributesFromDTO_Nil(t *testing.T) {
	result := WhitelistedContractAttributesFromDTO(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestWhitelistedContractAttributesFromDTO_Multiple(t *testing.T) {
	key1 := "key1"
	key2 := "key2"

	dtos := []openapi.TgvalidatordWhitelistedContractAddressAttribute{
		{Key: &key1},
		{Key: &key2},
	}

	result := WhitelistedContractAttributesFromDTO(dtos)

	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 attributes, got %d", len(result))
	}
	if result[0].Key != key1 {
		t.Errorf("expected first attribute Key %s, got %s", key1, result[0].Key)
	}
	if result[1].Key != key2 {
		t.Errorf("expected second attribute Key %s, got %s", key2, result[1].Key)
	}
}
