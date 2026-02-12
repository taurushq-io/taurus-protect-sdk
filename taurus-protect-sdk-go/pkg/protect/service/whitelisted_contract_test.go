package service

import (
	"context"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestWhitelistedContractService_GetWhitelistedContract_EmptyID(t *testing.T) {
	// Create a service with a nil API (we're testing validation, not API calls)
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetWhitelistedContract(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty ID")
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("expected 'id cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_CreateWhitelistedContract_NilRequest(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateWhitelistedContract(context.Background(), nil)
	if err == nil {
		t.Error("expected error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("expected 'request cannot be nil', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_CreateWhitelistedContract_MissingBlockchain(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &model.CreateWhitelistedContractRequest{
		Symbol: "USDC",
	}

	_, err := svc.CreateWhitelistedContract(context.Background(), req)
	if err == nil {
		t.Error("expected error for missing blockchain")
	}
	if err.Error() != "blockchain is required" {
		t.Errorf("expected 'blockchain is required', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_CreateWhitelistedContract_MissingSymbol(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &model.CreateWhitelistedContractRequest{
		Blockchain: "ETH",
	}

	_, err := svc.CreateWhitelistedContract(context.Background(), req)
	if err == nil {
		t.Error("expected error for missing symbol")
	}
	if err.Error() != "symbol is required" {
		t.Errorf("expected 'symbol is required', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_UpdateWhitelistedContract_EmptyID(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &model.UpdateWhitelistedContractRequest{
		Name: "Test",
	}

	_, err := svc.UpdateWhitelistedContract(context.Background(), "", req)
	if err == nil {
		t.Error("expected error for empty ID")
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("expected 'id cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_UpdateWhitelistedContract_NilRequest(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	tempID, err := svc.UpdateWhitelistedContract(context.Background(), "123", nil)
	if err == nil {
		t.Error("expected error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("expected 'request cannot be nil', got '%s'", err.Error())
	}
	if tempID != "" {
		t.Errorf("expected empty tempID, got '%s'", tempID)
	}
}

func TestWhitelistedContractService_DeleteWhitelistedContract_EmptyID(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.DeleteWhitelistedContract(context.Background(), "", "test comment")
	if err == nil {
		t.Error("expected error for empty ID")
	}
	if err.Error() != "id cannot be empty" {
		t.Errorf("expected 'id cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_ApproveWhitelistedContract_EmptyIDs(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.ApproveWhitelistedContract(context.Background(), []string{}, "sig", "comment")
	if err == nil {
		t.Error("expected error for empty IDs")
	}
	if err.Error() != "ids cannot be empty" {
		t.Errorf("expected 'ids cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_ApproveWhitelistedContract_EmptySignature(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.ApproveWhitelistedContract(context.Background(), []string{"123"}, "", "comment")
	if err == nil {
		t.Error("expected error for empty signature")
	}
	if err.Error() != "signature is required" {
		t.Errorf("expected 'signature is required', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_ApproveWhitelistedContract_EmptyComment(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.ApproveWhitelistedContract(context.Background(), []string{"123"}, "sig", "")
	if err == nil {
		t.Error("expected error for empty comment")
	}
	if err.Error() != "comment is required" {
		t.Errorf("expected 'comment is required', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_RejectWhitelistedContract_EmptyIDs(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RejectWhitelistedContract(context.Background(), []string{}, "comment")
	if err == nil {
		t.Error("expected error for empty IDs")
	}
	if err.Error() != "ids cannot be empty" {
		t.Errorf("expected 'ids cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_RejectWhitelistedContract_EmptyComment(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.RejectWhitelistedContract(context.Background(), []string{"123"}, "")
	if err == nil {
		t.Error("expected error for empty comment")
	}
	if err.Error() != "comment is required" {
		t.Errorf("expected 'comment is required', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_CreateWhitelistedContractAttribute_EmptyContractID(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &model.CreateWhitelistedContractAttributeRequest{
		Key:   "test",
		Value: "value",
	}

	_, err := svc.CreateWhitelistedContractAttribute(context.Background(), "", req)
	if err == nil {
		t.Error("expected error for empty contractID")
	}
	if err.Error() != "contractID cannot be empty" {
		t.Errorf("expected 'contractID cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_CreateWhitelistedContractAttribute_NilRequest(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateWhitelistedContractAttribute(context.Background(), "123", nil)
	if err == nil {
		t.Error("expected error for nil request")
	}
	if err.Error() != "request cannot be nil" {
		t.Errorf("expected 'request cannot be nil', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_CreateWhitelistedContractAttribute_EmptyKey(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	req := &model.CreateWhitelistedContractAttributeRequest{
		Value: "value",
	}

	_, err := svc.CreateWhitelistedContractAttribute(context.Background(), "123", req)
	if err == nil {
		t.Error("expected error for empty key")
	}
	if err.Error() != "key is required" {
		t.Errorf("expected 'key is required', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_CreateWhitelistedContractAttributes_EmptyContractID(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	reqs := []model.CreateWhitelistedContractAttributeRequest{
		{Key: "test", Value: "value"},
	}

	_, err := svc.CreateWhitelistedContractAttributes(context.Background(), "", reqs)
	if err == nil {
		t.Error("expected error for empty contractID")
	}
	if err.Error() != "contractID cannot be empty" {
		t.Errorf("expected 'contractID cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_CreateWhitelistedContractAttributes_EmptyRequests(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.CreateWhitelistedContractAttributes(context.Background(), "123", []model.CreateWhitelistedContractAttributeRequest{})
	if err == nil {
		t.Error("expected error for empty requests")
	}
	if err.Error() != "at least one attribute request is required" {
		t.Errorf("expected 'at least one attribute request is required', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_CreateWhitelistedContractAttributes_MissingKey(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	reqs := []model.CreateWhitelistedContractAttributeRequest{
		{Key: "valid", Value: "value"},
		{Value: "missing key"},
	}

	_, err := svc.CreateWhitelistedContractAttributes(context.Background(), "123", reqs)
	if err == nil {
		t.Error("expected error for missing key")
	}
	expectedErr := "key is required for attribute at index 1"
	if err.Error() != expectedErr {
		t.Errorf("expected '%s', got '%s'", expectedErr, err.Error())
	}
}

func TestWhitelistedContractService_GetWhitelistedContractAttribute_EmptyContractID(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetWhitelistedContractAttribute(context.Background(), "", "attr-1")
	if err == nil {
		t.Error("expected error for empty contractID")
	}
	if err.Error() != "contractID cannot be empty" {
		t.Errorf("expected 'contractID cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_GetWhitelistedContractAttribute_EmptyAttributeID(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	_, err := svc.GetWhitelistedContractAttribute(context.Background(), "123", "")
	if err == nil {
		t.Error("expected error for empty attributeID")
	}
	if err.Error() != "attributeID cannot be empty" {
		t.Errorf("expected 'attributeID cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_DeleteWhitelistedContractAttribute_EmptyContractID(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.DeleteWhitelistedContractAttribute(context.Background(), "", "attr-1")
	if err == nil {
		t.Error("expected error for empty contractID")
	}
	if err.Error() != "contractID cannot be empty" {
		t.Errorf("expected 'contractID cannot be empty', got '%s'", err.Error())
	}
}

func TestWhitelistedContractService_DeleteWhitelistedContractAttribute_EmptyAttributeID(t *testing.T) {
	svc := &WhitelistedContractService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	err := svc.DeleteWhitelistedContractAttribute(context.Background(), "123", "")
	if err == nil {
		t.Error("expected error for empty attributeID")
	}
	if err.Error() != "attributeID cannot be empty" {
		t.Errorf("expected 'attributeID cannot be empty', got '%s'", err.Error())
	}
}
