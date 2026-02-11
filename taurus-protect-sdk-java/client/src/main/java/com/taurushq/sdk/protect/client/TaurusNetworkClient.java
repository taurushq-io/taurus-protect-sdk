package com.taurushq.sdk.protect.client;

import com.taurushq.sdk.protect.client.service.TaurusNetworkLendingService;
import com.taurushq.sdk.protect.client.service.TaurusNetworkParticipantService;
import com.taurushq.sdk.protect.client.service.TaurusNetworkPledgeService;
import com.taurushq.sdk.protect.client.service.TaurusNetworkSettlementService;
import com.taurushq.sdk.protect.client.service.TaurusNetworkSharingService;

/**
 * Provides namespaced access to all Taurus Network services.
 * <p>
 * Access via: {@code client.taurusNetwork()}
 * <p>
 * This client provides a clean hierarchical API for Taurus Network operations,
 * grouping related functionality into logical sub-services:
 * <ul>
 *   <li>{@link #participants()} - Participant management</li>
 *   <li>{@link #pledges()} - Pledge lifecycle operations</li>
 *   <li>{@link #lending()} - Lending offers and agreements</li>
 *   <li>{@link #settlements()} - Settlement operations</li>
 *   <li>{@link #sharing()} - Address/asset sharing</li>
 * </ul>
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get current participant
 * Participant me = client.taurusNetwork().participants().getMyParticipant();
 *
 * // List pledges
 * PledgeResult pledges = client.taurusNetwork().pledges().list(null, null, null, null, null, null);
 *
 * // Get lending agreement
 * LendingAgreement agreement = client.taurusNetwork().lending().getAgreement("agreement-id");
 *
 * // Get settlement
 * Settlement settlement = client.taurusNetwork().settlements().get("settlement-id");
 *
 * // List shared addresses
 * SharedAddressResult addresses = client.taurusNetwork().sharing()
 *     .listSharedAddresses(null, null, null, "ETH", "mainnet", null, null, null);
 * }</pre>
 */
public class TaurusNetworkClient {

    private final TaurusNetworkParticipantService participantService;
    private final TaurusNetworkPledgeService pledgeService;
    private final TaurusNetworkLendingService lendingService;
    private final TaurusNetworkSettlementService settlementService;
    private final TaurusNetworkSharingService sharingService;

    /**
     * Instantiates a new Taurus Network client.
     *
     * @param participantService the participant service
     * @param pledgeService      the pledge service
     * @param lendingService     the lending service
     * @param settlementService  the settlement service
     * @param sharingService     the sharing service
     */
    public TaurusNetworkClient(
            final TaurusNetworkParticipantService participantService,
            final TaurusNetworkPledgeService pledgeService,
            final TaurusNetworkLendingService lendingService,
            final TaurusNetworkSettlementService settlementService,
            final TaurusNetworkSharingService sharingService) {
        this.participantService = participantService;
        this.pledgeService = pledgeService;
        this.lendingService = lendingService;
        this.settlementService = settlementService;
        this.sharingService = sharingService;
    }

    /**
     * Returns the participant management service.
     * <p>
     * Provides operations for managing Taurus Network participants and their attributes.
     *
     * @return the participant service
     */
    public TaurusNetworkParticipantService participants() {
        return participantService;
    }

    /**
     * Returns the pledge management service.
     * <p>
     * Provides operations for creating, managing, and withdrawing pledges.
     *
     * @return the pledge service
     */
    public TaurusNetworkPledgeService pledges() {
        return pledgeService;
    }

    /**
     * Returns the lending offers and agreements service.
     * <p>
     * Provides operations for managing lending offers and agreements.
     *
     * @return the lending service
     */
    public TaurusNetworkLendingService lending() {
        return lendingService;
    }

    /**
     * Returns the settlement management service.
     * <p>
     * Provides operations for creating and managing settlements between participants.
     *
     * @return the settlement service
     */
    public TaurusNetworkSettlementService settlements() {
        return settlementService;
    }

    /**
     * Returns the address/asset sharing service.
     * <p>
     * Provides operations for sharing and unsharing addresses and assets with other participants.
     *
     * @return the sharing service
     */
    public TaurusNetworkSharingService sharing() {
        return sharingService;
    }
}
