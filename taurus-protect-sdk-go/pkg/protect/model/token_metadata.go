package model

// CryptoPunkMetadata represents metadata for a CryptoPunk NFT.
type CryptoPunkMetadata struct {
	// PunkID is the ID of the CryptoPunk (0-9999).
	PunkID string `json:"punk_id"`
	// Attributes contains the attributes of the CryptoPunk as a JSON string.
	Attributes string `json:"attributes"`
	// Image contains the base64-encoded image data of the CryptoPunk.
	Image string `json:"image"`
}

// ERCTokenMetadata represents metadata for an ERC-721 or ERC-1155 token.
type ERCTokenMetadata struct {
	// Name is the name of the token.
	Name string `json:"name"`
	// Description is the description of the token.
	Description string `json:"description"`
	// Decimals is the number of decimals for the token.
	Decimals string `json:"decimals,omitempty"`
	// DataType is the MIME type of the base64 data (e.g., "image/png").
	DataType string `json:"data_type,omitempty"`
	// Base64Data contains the base64-encoded data (e.g., image) if requested.
	Base64Data string `json:"base64_data,omitempty"`
	// URI is the token URI pointing to the metadata.
	URI string `json:"uri,omitempty"`
}

// GetERCTokenMetadataOptions contains options for retrieving ERC token metadata.
type GetERCTokenMetadataOptions struct {
	// WithData requests base64-encoded data to be included in the response.
	WithData bool
	// Blockchain specifies the blockchain to use.
	Blockchain string
}
