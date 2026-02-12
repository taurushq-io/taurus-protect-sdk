package protect

import "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/service"

// TaurusNetworkClient provides namespaced access to Taurus Network services.
// Use client.TaurusNetwork() to access this client.
//
// Example usage:
//
//	client.TaurusNetwork().Participants().GetMyParticipant(ctx)
//	client.TaurusNetwork().Pledges().Create(ctx, req)
//	client.TaurusNetwork().Lending().GetAgreement(ctx, agreementID)
//	client.TaurusNetwork().Settlements().Get(ctx, settlementID)
//	client.TaurusNetwork().Sharing().ShareAddress(ctx, req)
type TaurusNetworkClient struct {
	participants *service.TaurusNetworkParticipantService
	pledges      *service.TaurusNetworkPledgeService
	lending      *service.TaurusNetworkLendingService
	settlements  *service.TaurusNetworkSettlementService
	sharing      *service.TaurusNetworkSharingService
}

// Participants returns the participant management service.
// Provides operations for managing Taurus Network participants and their attributes.
func (tn *TaurusNetworkClient) Participants() *service.TaurusNetworkParticipantService {
	return tn.participants
}

// Pledges returns the pledge management service.
// Provides operations for creating, managing, and withdrawing pledges.
func (tn *TaurusNetworkClient) Pledges() *service.TaurusNetworkPledgeService {
	return tn.pledges
}

// Lending returns the lending offers and agreements service.
// Provides operations for managing lending offers and agreements.
func (tn *TaurusNetworkClient) Lending() *service.TaurusNetworkLendingService {
	return tn.lending
}

// Settlements returns the settlement management service.
// Provides operations for creating and managing settlements between participants.
func (tn *TaurusNetworkClient) Settlements() *service.TaurusNetworkSettlementService {
	return tn.settlements
}

// Sharing returns the address/asset sharing service.
// Provides operations for sharing and unsharing addresses and assets with other participants.
func (tn *TaurusNetworkClient) Sharing() *service.TaurusNetworkSharingService {
	return tn.sharing
}
