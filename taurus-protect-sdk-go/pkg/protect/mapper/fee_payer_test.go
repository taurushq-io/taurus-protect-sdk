package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestFeePayerFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordFeePayerEnvelope
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns fee payer with zero values",
			dto:  &openapi.TgvalidatordFeePayerEnvelope{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordFeePayerEnvelope {
				id := "feepayer-123"
				tenantId := "tenant-456"
				blockchain := "ETH"
				network := "mainnet"
				name := "Main Fee Payer"
				creationDate := time.Now()
				return &openapi.TgvalidatordFeePayerEnvelope{
					Id:           &id,
					TenantId:     &tenantId,
					Blockchain:   &blockchain,
					Network:      &network,
					Name:         &name,
					CreationDate: &creationDate,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FeePayerFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("FeePayerFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("FeePayerFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.CreationDate != nil && !got.CreationDate.Equal(*tt.dto.CreationDate) {
				t.Errorf("CreationDate = %v, want %v", got.CreationDate, *tt.dto.CreationDate)
			}
		})
	}
}

func TestFeePayersFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordFeePayerEnvelope
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordFeePayerEnvelope{},
			want: 0,
		},
		{
			name: "converts multiple fee payers",
			dtos: func() []openapi.TgvalidatordFeePayerEnvelope {
				id1 := "feepayer-1"
				id2 := "feepayer-2"
				return []openapi.TgvalidatordFeePayerEnvelope{
					{Id: &id1},
					{Id: &id2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FeePayersFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("FeePayersFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("FeePayersFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestFeePayerETHFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordFeePayer
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns ETH config with zero values",
			dto:  &openapi.TgvalidatordFeePayer{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordFeePayer {
				blockchain := "ETH"
				kind := "local"
				remoteEncrypted := "YmFzZTY0ZW5jb2RlZA=="
				addressId := "addr-123"
				autoApprove := true
				return &openapi.TgvalidatordFeePayer{
					Blockchain: &blockchain,
					Eth: &openapi.FeePayerETH{
						Kind:            &kind,
						RemoteEncrypted: &remoteEncrypted,
						Local: &openapi.ETHLocal{
							AddressId:   &addressId,
							AutoApprove: &autoApprove,
						},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FeePayerETHFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("FeePayerETHFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("FeePayerETHFromDTO() returned nil for non-nil input")
			}
			// Verify blockchain field
			if tt.dto.Blockchain != nil && got.Blockchain != *tt.dto.Blockchain {
				t.Errorf("Blockchain = %v, want %v", got.Blockchain, *tt.dto.Blockchain)
			}
			// Verify ETH-specific fields
			if tt.dto.Eth != nil {
				if tt.dto.Eth.Kind != nil && got.Kind != *tt.dto.Eth.Kind {
					t.Errorf("Kind = %v, want %v", got.Kind, *tt.dto.Eth.Kind)
				}
				if tt.dto.Eth.RemoteEncrypted != nil && got.RemoteEncrypted != *tt.dto.Eth.RemoteEncrypted {
					t.Errorf("RemoteEncrypted = %v, want %v", got.RemoteEncrypted, *tt.dto.Eth.RemoteEncrypted)
				}
				if tt.dto.Eth.Local != nil && got.Local == nil {
					t.Error("Local should not be nil when DTO has local")
				}
			}
		})
	}
}

func TestETHLocalConfigFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.ETHLocal
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns config with zero values",
			dto:  &openapi.ETHLocal{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.ETHLocal {
				addressId := "addr-123"
				forwarderAddressId := "forwarder-456"
				autoApprove := true
				creatorAddressId := "creator-789"
				domainSeparator := "ZG9tYWluU2VwYXJhdG9y"
				forwarderKind := openapi.TGVALIDATORDFEEPAYERFORWARDERKIND_OPEN_ZEPPELIN_FORWARDER
				return &openapi.ETHLocal{
					AddressId:          &addressId,
					ForwarderAddressId: &forwarderAddressId,
					AutoApprove:        &autoApprove,
					CreatorAddressId:   &creatorAddressId,
					DomainSeparator:    &domainSeparator,
					ForwarderKind:      &forwarderKind,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ETHLocalConfigFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ETHLocalConfigFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ETHLocalConfigFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.AddressId != nil && got.AddressID != *tt.dto.AddressId {
				t.Errorf("AddressID = %v, want %v", got.AddressID, *tt.dto.AddressId)
			}
			if tt.dto.ForwarderAddressId != nil && got.ForwarderAddressID != *tt.dto.ForwarderAddressId {
				t.Errorf("ForwarderAddressID = %v, want %v", got.ForwarderAddressID, *tt.dto.ForwarderAddressId)
			}
			if tt.dto.AutoApprove != nil && got.AutoApprove != *tt.dto.AutoApprove {
				t.Errorf("AutoApprove = %v, want %v", got.AutoApprove, *tt.dto.AutoApprove)
			}
			if tt.dto.CreatorAddressId != nil && got.CreatorAddressID != *tt.dto.CreatorAddressId {
				t.Errorf("CreatorAddressID = %v, want %v", got.CreatorAddressID, *tt.dto.CreatorAddressId)
			}
			if tt.dto.DomainSeparator != nil && got.DomainSeparator != *tt.dto.DomainSeparator {
				t.Errorf("DomainSeparator = %v, want %v", got.DomainSeparator, *tt.dto.DomainSeparator)
			}
			if tt.dto.ForwarderKind != nil && got.ForwarderKind != string(*tt.dto.ForwarderKind) {
				t.Errorf("ForwarderKind = %v, want %v", got.ForwarderKind, string(*tt.dto.ForwarderKind))
			}
		})
	}
}

func TestETHRemoteConfigFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.ETHRemote
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns config with zero values",
			dto:  &openapi.ETHRemote{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.ETHRemote {
				url := "https://feepayer.example.com"
				username := "user"
				password := "pass"
				privateKey := "0x123"
				fromAddressId := "from-123"
				forwarderAddress := "0xforwarder"
				forwarderAddressId := "forwarder-456"
				creatorAddress := "0xcreator"
				creatorAddressId := "creator-789"
				domainSeparator := "ZG9tYWluU2VwYXJhdG9y"
				forwarderKind := openapi.TGVALIDATORDFEEPAYERFORWARDERKIND_OPEN_ZEPPELIN_FORWARDER
				return &openapi.ETHRemote{
					Url:                &url,
					Username:           &username,
					Password:           &password,
					PrivateKey:         &privateKey,
					FromAddressId:      &fromAddressId,
					ForwarderAddress:   &forwarderAddress,
					ForwarderAddressId: &forwarderAddressId,
					CreatorAddress:     &creatorAddress,
					CreatorAddressId:   &creatorAddressId,
					DomainSeparator:    &domainSeparator,
					ForwarderKind:      &forwarderKind,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ETHRemoteConfigFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("ETHRemoteConfigFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("ETHRemoteConfigFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Url != nil && got.URL != *tt.dto.Url {
				t.Errorf("URL = %v, want %v", got.URL, *tt.dto.Url)
			}
			if tt.dto.Username != nil && got.Username != *tt.dto.Username {
				t.Errorf("Username = %v, want %v", got.Username, *tt.dto.Username)
			}
			if tt.dto.Password != nil && got.Password != *tt.dto.Password {
				t.Errorf("Password = %v, want %v", got.Password, *tt.dto.Password)
			}
			if tt.dto.PrivateKey != nil && got.PrivateKey != *tt.dto.PrivateKey {
				t.Errorf("PrivateKey = %v, want %v", got.PrivateKey, *tt.dto.PrivateKey)
			}
			if tt.dto.FromAddressId != nil && got.FromAddressID != *tt.dto.FromAddressId {
				t.Errorf("FromAddressID = %v, want %v", got.FromAddressID, *tt.dto.FromAddressId)
			}
			if tt.dto.ForwarderAddress != nil && got.ForwarderAddress != *tt.dto.ForwarderAddress {
				t.Errorf("ForwarderAddress = %v, want %v", got.ForwarderAddress, *tt.dto.ForwarderAddress)
			}
			if tt.dto.ForwarderAddressId != nil && got.ForwarderAddressID != *tt.dto.ForwarderAddressId {
				t.Errorf("ForwarderAddressID = %v, want %v", got.ForwarderAddressID, *tt.dto.ForwarderAddressId)
			}
			if tt.dto.CreatorAddress != nil && got.CreatorAddress != *tt.dto.CreatorAddress {
				t.Errorf("CreatorAddress = %v, want %v", got.CreatorAddress, *tt.dto.CreatorAddress)
			}
			if tt.dto.CreatorAddressId != nil && got.CreatorAddressID != *tt.dto.CreatorAddressId {
				t.Errorf("CreatorAddressID = %v, want %v", got.CreatorAddressID, *tt.dto.CreatorAddressId)
			}
			if tt.dto.DomainSeparator != nil && got.DomainSeparator != *tt.dto.DomainSeparator {
				t.Errorf("DomainSeparator = %v, want %v", got.DomainSeparator, *tt.dto.DomainSeparator)
			}
			if tt.dto.ForwarderKind != nil && got.ForwarderKind != string(*tt.dto.ForwarderKind) {
				t.Errorf("ForwarderKind = %v, want %v", got.ForwarderKind, string(*tt.dto.ForwarderKind))
			}
		})
	}
}

func TestChecksumRequestToDTO(t *testing.T) {
	tests := []struct {
		name string
		req  *model.ChecksumRequest
		want string
	}{
		{
			name: "nil request returns empty DTO",
			req:  nil,
			want: "",
		},
		{
			name: "empty request returns empty data",
			req:  &model.ChecksumRequest{},
			want: "",
		},
		{
			name: "request with data maps correctly",
			req:  &model.ChecksumRequest{Data: "SGVsbG8gV29ybGQ="},
			want: "SGVsbG8gV29ybGQ=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ChecksumRequestToDTO(tt.req)
			gotData := ""
			if got.Data != nil {
				gotData = *got.Data
			}
			if gotData != tt.want {
				t.Errorf("ChecksumRequestToDTO().Data = %v, want %v", gotData, tt.want)
			}
		})
	}
}

func TestFeePayerFromDTO_WithNestedFeePayer(t *testing.T) {
	id := "feepayer-123"
	blockchain := "ETH"
	kind := "local"
	addressId := "addr-456"
	autoApprove := true

	dto := &openapi.TgvalidatordFeePayerEnvelope{
		Id:         &id,
		Blockchain: &blockchain,
		FeePayer: &openapi.TgvalidatordFeePayer{
			Blockchain: &blockchain,
			Eth: &openapi.FeePayerETH{
				Kind: &kind,
				Local: &openapi.ETHLocal{
					AddressId:   &addressId,
					AutoApprove: &autoApprove,
				},
			},
		},
	}

	got := FeePayerFromDTO(dto)
	if got == nil {
		t.Fatal("FeePayerFromDTO() returned nil for non-nil input")
	}
	if got.ETH == nil {
		t.Fatal("ETH should not be nil when DTO has FeePayer")
	}
	if got.ETH.Kind != kind {
		t.Errorf("ETH.Kind = %v, want %v", got.ETH.Kind, kind)
	}
	if got.ETH.Local == nil {
		t.Fatal("ETH.Local should not be nil when DTO has Local config")
	}
	if got.ETH.Local.AddressID != addressId {
		t.Errorf("ETH.Local.AddressID = %v, want %v", got.ETH.Local.AddressID, addressId)
	}
	if got.ETH.Local.AutoApprove != autoApprove {
		t.Errorf("ETH.Local.AutoApprove = %v, want %v", got.ETH.Local.AutoApprove, autoApprove)
	}
}

func TestFeePayerFromDTO_NilCreationDate(t *testing.T) {
	id := "feepayer-123"
	dto := &openapi.TgvalidatordFeePayerEnvelope{
		Id:           &id,
		CreationDate: nil,
	}

	got := FeePayerFromDTO(dto)
	if got == nil {
		t.Fatal("FeePayerFromDTO() returned nil for non-nil input")
	}
	// When creation date is nil, it should be the zero time value
	if !got.CreationDate.IsZero() {
		t.Errorf("CreationDate should be zero time when nil, got %v", got.CreationDate)
	}
}

func TestFeePayerFromDTO_NilFeePayer(t *testing.T) {
	id := "feepayer-123"
	dto := &openapi.TgvalidatordFeePayerEnvelope{
		Id:       &id,
		FeePayer: nil,
	}

	got := FeePayerFromDTO(dto)
	if got == nil {
		t.Fatal("FeePayerFromDTO() returned nil for non-nil input")
	}
	if got.ETH != nil {
		t.Errorf("ETH should be nil when DTO FeePayer is nil, got %v", got.ETH)
	}
	if got.ID != "feepayer-123" {
		t.Errorf("ID = %v, want feepayer-123", got.ID)
	}
}

func TestETHLocalConfigFromDTO_AutoApproveValues(t *testing.T) {
	tests := []struct {
		name            string
		autoApprove     *bool
		wantAutoApprove bool
	}{
		{
			name:            "nil auto approve defaults to false",
			autoApprove:     nil,
			wantAutoApprove: false,
		},
		{
			name:            "true auto approve",
			autoApprove:     boolPtr(true),
			wantAutoApprove: true,
		},
		{
			name:            "false auto approve",
			autoApprove:     boolPtr(false),
			wantAutoApprove: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.ETHLocal{
				AutoApprove: tt.autoApprove,
			}
			got := ETHLocalConfigFromDTO(dto)
			if got.AutoApprove != tt.wantAutoApprove {
				t.Errorf("AutoApprove = %v, want %v", got.AutoApprove, tt.wantAutoApprove)
			}
		})
	}
}
