package helper

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/taurusgroup/tg-validatord/internal/api/v1/tgvalidatord"
	"github.com/taurusgroup/tg-validatord/internal/config"
	"github.com/taurusgroup/tg-validatord/internal/constants"
	score_model "github.com/taurusgroup/tg-validatord/pkg/score/model"
	"github.com/taurusgroup/tg-validatord/pkg/whitelist/model"
)

func Team1PrivKey(t *testing.T) *ecdsa.PrivateKey {
	const pemData = `
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIAlDUW8kMrUaxaa1ZmKGysHVCsj5AJmcN8Vc89URxdgjoAoGCCqGSM49
AwEHoUQDQgAELJhEUNLLHgI8LiWJaeJGpaBfdvgoYyKsjSFyTMxECR/E+1qpzDlN
Nug7hDPgBPpZ3Z+U8QWjaKB4Mrbj2/kImQ==
-----END EC PRIVATE KEY-----`

	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		t.Fatal("failed to parse PEM block containing the private key")
	}

	secKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		t.Fatal("failed to parse secret key: " + err.Error())
	}

	return secKey
}

func Team2PrivKey(t *testing.T) *ecdsa.PrivateKey {
	const pemData = `
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIC2Ds7ItfXaKwo12IM3avq0iua3d9+znj2f7fOlSyMQwoAoGCCqGSM49
AwEHoUQDQgAE4UAo+GNxw4tEqdmTDi6v2MZq2ug8ajW8/MnqJ/Qd4bniPyDvITrS
EZUyP9TvYpkLAbg/ACxcR/yA6lBEgtE60Q==
-----END EC PRIVATE KEY-----`

	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		t.Fatal("failed to parse PEM block containing the private key")
	}

	secKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		t.Fatal("failed to parse secret key: " + err.Error())
	}

	return secKey
}

func Team3PrivKey(t *testing.T) *ecdsa.PrivateKey {
	const pemData = `
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEILf9mDoHW+tqOskGJxbZgMNIZv2cfuMDbF9Pl7DmwIjxoAoGCCqGSM49
AwEHoUQDQgAEHjiwU8CwxPPF+SxbTVOm7Z1nd9P6mTbFC2FYUYiSlAJwfE7woM7B
AOg6b8W0y0JIJAc+t16V7zOMWzsFVmee6g==
-----END EC PRIVATE KEY-----`

	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		t.Fatal("failed to parse PEM block containing the private key")
	}

	secKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		t.Fatal("failed to parse secret key: " + err.Error())
	}

	return secKey
}

func User1PrivKey(t *testing.T) *ecdsa.PrivateKey {
	const secretPEMKey = `
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIOd7BwfDcXGDo0cTF9KczH9/jq27xIUEFk6v7iCeY5n3oAoGCCqGSM49
AwEHoUQDQgAEmtXvCSwMCarLGbVX/l6x0GTnkXMreg6fLAVtHkwKZ6H4L7J9WhRC
VtTzTOgfvOi2zt68Jm7tbhDY9OYWuITOBA==
-----END EC PRIVATE KEY-----`

	block, _ := pem.Decode([]byte(secretPEMKey))
	if block == nil {
		t.Fatal("failed To parse PEM block containing the private key")
	}

	secKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		t.Fatal("failed To parse secret key: " + err.Error())
	}

	return secKey
}

func User1PubKey(t *testing.T) *ecdsa.PublicKey {

	const pubPEM = `
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEmtXvCSwMCarLGbVX/l6x0GTnkXMr
eg6fLAVtHkwKZ6H4L7J9WhRCVtTzTOgfvOi2zt68Jm7tbhDY9OYWuITOBA==
-----END PUBLIC KEY-----`

	k, err := loadPubKey(pubPEM)
	assert.Nil(t, err)
	return k
}

func loadPubKey(pubPEM string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPEM))
	if block == nil {
		return nil, fmt.Errorf("unable To decode pubKey: %v", pubPEM)
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed To parse pubKey")
	}

	p, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not an ecdsa.PublicKey")
	}
	return p, nil
}

func TestHashWLAddress(t *testing.T) {

	blockchainCfg := config.Blockchain{Blockchain: "BTC", Network: "mainnet", IncludeNetworkInPayload: false}

	_, err := ToWLAMetadata(nil, blockchainCfg)
	assert.NotNil(t, err)

	h, err := ToWLAMetadata(&model.WLAddress{
		Id:                2,
		TenantID:          1,
		Envelope:          createEnv(t, 0, "mainnet"),
		Action:            "action",
		Address:           "address",
		Blockchain:        "BTC",
		Network:           "mainnet",
		AddressType:       "individual",
		Label:             "label",
		Memo:              "memo",
		CustomerID:        "custID",
		ExchangeAccountID: nil,
		Scores: []score_model.Score{{
			ID:         0,
			AddressID:  0,
			Provider:   "",
			Type:       "",
			Score:      "",
			UpdateDate: time.Time{},
		}},
		Trails: []*model.WLAddressTrail{{
			Id:             0,
			WlAddressId:    0,
			UserId:         0,
			ExternalUserId: "",
			Action:         "",
			Date:           time.Time{},
			Comment:        "",
		}},
	}, blockchainCfg)
	assert.Nil(t, err)
	assert.Equal(t, "4c5b0f8f7912fd0419fc0aeaa01539955714aa4f20f72fe7ca3a691308ef881a", h.Hash)
	fmt.Printf("hash: %v\n", h)

	var eid uint64 = 10

	h, err = ToWLAMetadata(&model.WLAddress{
		Id:                2,
		TenantID:          1,
		Envelope:          createEnv(t, eid, "mainnet"),
		Action:            "action",
		Address:           "address",
		Blockchain:        "BTC",
		Network:           "mainnet",
		AddressType:       "individual",
		Label:             "label",
		Memo:              "memo",
		CustomerID:        "custID",
		ExchangeAccountID: &eid,
	}, blockchainCfg)
	assert.Nil(t, err)
	assert.Equal(t, "793116dd37a39365803558417a23855324575ec44a8bcfa354d5e13cddc35ea9", h.Hash)
	fmt.Printf("hash: %v\n", h)

	eid = 0
	h, err = ToWLAMetadata(&model.WLAddress{
		Id:                2,
		TenantID:          1,
		Envelope:          createEnv(t, eid, "mainnet"),
		Action:            "action",
		Address:           "address",
		Blockchain:        "BTC",
		Network:           "mainnet",
		AddressType:       "individual",
		Label:             "label",
		Memo:              "memo",
		CustomerID:        "custID",
		ExchangeAccountID: &eid,
	}, blockchainCfg)
	assert.Nil(t, err)
	assert.Equal(t, "4c5b0f8f7912fd0419fc0aeaa01539955714aa4f20f72fe7ca3a691308ef881a", h.Hash)
	fmt.Printf("hash: %v\n", h)

	eid = 0
	h, err = ToWLAMetadata(&model.WLAddress{
		Id:                2,
		TenantID:          1,
		Envelope:          createEnv(t, eid, "mainnet"),
		Action:            "action",
		Address:           "address",
		Blockchain:        "BTC",
		Network:           "mainnet",
		AddressType:       "individual",
		Label:             "label",
		Memo:              "memo",
		CustomerID:        "custID",
		ExchangeAccountID: &eid,
		LinkedInternalAddresses: []model.InternalAddress{
			{
				ID:      "1",
				Address: "addr1",
				Label:   "test1",
			}, {
				ID:      "2",
				Address: "addr2",
				Label:   "test2",
			},
		},
	}, blockchainCfg)
	assert.Nil(t, err)
	assert.Equal(t, "86bb97105ee1625d6140c59df3c33959fb11916c9796c65b215d5a314f0b657c", h.Hash)
	fmt.Printf("hash: %v\n", h)

	eid = 0
	h, err = ToWLAMetadata(&model.WLAddress{
		Id:                2,
		TenantID:          1,
		Envelope:          createEnv(t, eid, "mainnet"),
		Action:            "action",
		Address:           "address",
		Blockchain:        "BTC",
		Network:           "mainnet",
		AddressType:       "individual",
		Label:             "label",
		Memo:              "memo",
		CustomerID:        "custID",
		ExchangeAccountID: &eid,
		LinkedWallets: []model.InternalWallet{
			{
				ID:   1,
				Name: "test1",
			},
			{
				ID:   2,
				Name: "test2",
			},
		},
	}, blockchainCfg)
	assert.Nil(t, err)
	assert.Equal(t, "0259de5ce54cda3cb08c50e7c50cf7c601f765fb0d0be050208b36ac7251cc03", h.Hash)
	fmt.Printf("hash: %v\n", h)

	blockchainCfg = config.Blockchain{Blockchain: "BTC", Network: "testnet", IncludeNetworkInPayload: false}

	h, err = ToWLAMetadata(&model.WLAddress{
		Id:          2,
		TenantID:    1,
		Envelope:    createEnv(t, 0, "testnet"),
		Action:      "action",
		Address:     "address",
		Blockchain:  "BTC",
		Network:     "testnet",
		AddressType: "individual",
		Label:       "label",
		Memo:        "memo",
		CustomerID:  "custID",
	}, blockchainCfg)
	assert.Nil(t, err)
	assert.Equal(t, "4c5b0f8f7912fd0419fc0aeaa01539955714aa4f20f72fe7ca3a691308ef881a", h.Hash)
	fmt.Printf("hash: %v\n", h)
	hashBTCNetworkNotInPayload := "4c5b0f8f7912fd0419fc0aeaa01539955714aa4f20f72fe7ca3a691308ef881a"

	blockchainCfg = config.Blockchain{Blockchain: "BTC", Network: "testnet", IncludeNetworkInPayload: true}

	h, err = ToWLAMetadata(&model.WLAddress{
		Id:          2,
		TenantID:    1,
		Envelope:    createEnv(t, 0, "testnet"),
		Action:      "action",
		Address:     "address",
		Blockchain:  "BTC",
		Network:     "testnet",
		AddressType: "individual",
		Label:       "label",
		Memo:        "memo",
		CustomerID:  "custID",
	}, blockchainCfg)
	assert.Nil(t, err)
	assert.Equal(t, "a31e2cb9da39399c07402a87d31b34cde00e154cc93a8fc7b06905cc39eaeb58", h.Hash)
	fmt.Printf("hash: %v\n", h)
	require.NotEqual(t, hashBTCNetworkNotInPayload, h.Hash)
}

func hashContractAddressPayload(payload string) (string, error) {
	wa, err := DecodeWhitelistedContractAddress(payload)
	if err != nil {
		return "", err
	}

	b := &model.BasicWLContractAddress{
		Blockchain:      wa.GetBlockchain(),
		Name:            wa.GetName(),
		Symbol:          wa.GetSymbol(),
		Decimals:        fmt.Sprintf("%v", wa.GetDecimals()),
		ContractAddress: wa.GetContractAddress(),
	}

	fmt.Printf("basic address: %v\n", b)

	j, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	fmt.Printf("json: %v\n", string(j))

	digest := sha256.Sum256(j)
	return hex.EncodeToString(digest[:]), nil

}

func hashAddressPayload(payload string) (string, error) {
	wa, err := DecodeWhitelistedAddress(payload)
	if err != nil {
		return "", err
	}

	eaid := ""
	if wa.GetExchangeAccountId() > 0 {
		eaid = fmt.Sprintf("%v", wa.GetExchangeAccountId())
	}
	b := &model.BasicWLAddress{
		Address:           wa.GetAddress(),
		Currency:          wa.GetBlockchain(),
		AddressType:       wa.GetAddressType().String(),
		Label:             wa.GetLabel(),
		Memo:              wa.GetMemo(),
		CustomerID:        wa.GetCustomerId(),
		ExchangeAccountID: eaid,
		ContractType:      wa.GetContractType(),
	}

	fmt.Printf("basic address: %v\n", b)

	j, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	fmt.Printf("json: %v\n", string(j))

	digest := sha256.Sum256(j)
	return hex.EncodeToString(digest[:]), nil

}

func TestAddressApproveSignature(t *testing.T) {
	payloads := []string{
		"CgNFVEgQARoqMHhmNjMxY2U4OTNlZGI0NDBlNDkxODhhOTkxMjUwNTFkMDc5NjgxODY0KkBteSBzZWNvbmQgRVRIIGFkZHJlc3Mgb24gYmlzdGFtcCAoYWN0dWFsbHksIGFuIGludGVybmFsIGFkZHJlc3MpOAM=",
		"CgNFVEgQARoqMHhmYTJkM2E1NDNmMDAyNDk4MWNiOWVlNzk4ZmRkZDkzMDQ1MDgzZTY4Kj9teSBmaXJzdCBFVEggYWRkcmVzcyBvbiBiaXN0YW1wIChhY3R1YWxseSwgYW4gaW50ZXJuYWwgYWRkcmVzcyk4Aw==",
		"CgNFVEgQARoqMHg0MzNjZmU0OTljNjQ1ZTA4MzdmODU0MmFlYzU3MTM4YmJkMDg3NWUxKj5teSBmaXJzdCBFVEggYWRkcmVzcyBvbiBrcmFrZW4gKGFjdHVhbGx5LCBhbiBpbnRlcm5hbCBhZGRyZXNzKTgI",
	}

	hashes := []string{}

	for _, p := range payloads {
		h, err := hashAddressPayload(p)
		assert.Nil(t, err)
		hashes = append(hashes, h)
	}

	fmt.Printf("hashes: %v\n", hashes)

	sig1, err := SignHashes(hashes, Team1PrivKey(t))
	assert.Nil(t, err)
	fmt.Printf("sig1: %v\n", base64.StdEncoding.EncodeToString(sig1))
	sig2, err := SignHashes(hashes, Team2PrivKey(t))
	assert.Nil(t, err)
	fmt.Printf("sig2: %v\n", base64.StdEncoding.EncodeToString(sig2))
	sig3, err := SignHashes(hashes, Team3PrivKey(t))
	assert.Nil(t, err)
	fmt.Printf("sig3: %v\n", base64.StdEncoding.EncodeToString(sig3))
}

func TestContractAddressApproveSignature(t *testing.T) {
	payload := "CgNFVEgSD0NoYWluTGluayBUb2tlbhoETElOSyASKioweDUxNDkxMDc3MWFmOWNhNjU2YWY4NDBkZmY4M2U4MjY0ZWNmOTg2Y2E="

	h, err := hashContractAddressPayload(payload)
	assert.Nil(t, err)

	hashes := []string{h}
	fmt.Printf("hashes: %v\n", hashes)

	sig1, err := SignHashes(hashes, Team1PrivKey(t))
	assert.Nil(t, err)
	fmt.Printf("sig1: %v\n", base64.StdEncoding.EncodeToString(sig1))
	sig2, err := SignHashes(hashes, Team2PrivKey(t))
	assert.Nil(t, err)
	fmt.Printf("sig2: %v\n", base64.StdEncoding.EncodeToString(sig2))
	sig3, err := SignHashes(hashes, Team3PrivKey(t))
	assert.Nil(t, err)
	fmt.Printf("sig3: %v\n", base64.StdEncoding.EncodeToString(sig3))
}

func TestSubmitSignature(t *testing.T) {
	blockchain := constants.ETHBlockchain
	network := "mainnet"
	address := "0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d"

	blockchainCfg := config.Blockchain{Blockchain: blockchain, Network: network, IncludeNetworkInPayload: false}

	env, err := EncodeWhitelistedAddress(&tgvalidatord.WhitelistedAddress{
		Blockchain:  blockchain.String(),
		Network:     network,
		AddressType: tgvalidatord.WhitelistedAddress_individual,
		Address:     address,
		Memo:        "memo",
		Label:       "label",
		CustomerId:  "custID",
	}, nil)
	assert.Nil(t, err)

	wla := &model.WLAddress{
		Id:          2,
		TenantID:    1,
		Envelope:    env,
		Action:      "action",
		Address:     address,
		Blockchain:  blockchain,
		Network:     network,
		AddressType: "individual",
		Label:       "label",
		Memo:        "memo",
		CustomerID:  "custID",
	}

	h, err := ToWLAMetadata(wla, blockchainCfg)
	assert.Nil(t, err)

	hashes := []string{h.Hash}
	fmt.Printf("hashes: %v\n", hashes)

	sig1, err := SignHashes(hashes, Team1PrivKey(t))
	assert.Nil(t, err)
	fmt.Printf("sig1: %v\n", base64.StdEncoding.EncodeToString(sig1))
	sig2, err := SignHashes(hashes, Team2PrivKey(t))
	assert.Nil(t, err)
	fmt.Printf("sig2: %v\n", base64.StdEncoding.EncodeToString(sig2))
	sig3, err := SignHashes(hashes, Team3PrivKey(t))
	assert.Nil(t, err)
	fmt.Printf("sig3: %v\n", base64.StdEncoding.EncodeToString(sig3))
}

func TestCheckWLAddressHashesSignature(t *testing.T) {

	blockchainCfg := config.Blockchain{Blockchain: "BTC", Network: "mainnet", IncludeNetworkInPayload: false}

	wla := &model.WLAddress{
		Id:          2,
		TenantID:    1,
		Envelope:    createEnv(t, 0, "mainnet"),
		Action:      "action",
		Address:     "address",
		Blockchain:  "BTC",
		Network:     "mainnet",
		AddressType: "individual",
		Label:       "label",
		Memo:        "memo",
		CustomerID:  "custID",
	}

	h, err := ToWLAMetadata(wla, blockchainCfg)
	assert.Nil(t, err)

	hashes := []string{h.Hash}

	sig, err := SignHashes(hashes, User1PrivKey(t))
	assert.Nil(t, err)

	err = CheckWLAddressHashesSignature(blockchainCfg, wla, hashes, sig, User1PubKey(t))
	assert.Nil(t, err)

	blockchainCfg = config.Blockchain{Blockchain: "BTC", Network: "mainnet", IncludeNetworkInPayload: true}

	h2, err := ToWLAMetadata(wla, blockchainCfg)
	assert.Nil(t, err)
	hashes2 := []string{h2.Hash}

	// We check that the sig isn't valid when we use the one from the previous WLA (network not in payload)
	err = CheckWLAddressHashesSignature(blockchainCfg, wla, hashes2, sig, User1PubKey(t))
	assert.Error(t, err)

	sig2, err := SignHashes(hashes2, User1PrivKey(t))
	assert.Nil(t, err)

	err = CheckWLAddressHashesSignature(blockchainCfg, wla, hashes2, sig2, User1PubKey(t))
	assert.NoError(t, err)

}

func createEnv(t *testing.T, exchangeID uint64, network string) string {
	env, err := EncodeWhitelistedAddress(&tgvalidatord.WhitelistedAddress{
		Blockchain:        "BTC",
		Network:           network,
		AddressType:       tgvalidatord.WhitelistedAddress_individual,
		Address:           "address",
		Memo:              "memo",
		Label:             "label",
		CustomerId:        "custID",
		ExchangeAccountId: exchangeID,
	}, nil)
	assert.Nil(t, err)
	return env
}
