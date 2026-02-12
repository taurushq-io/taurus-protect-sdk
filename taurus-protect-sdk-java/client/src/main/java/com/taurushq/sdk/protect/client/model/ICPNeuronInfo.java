package com.taurushq.sdk.protect.client.model;

/**
 * Represents Internet Computer Protocol (ICP) neuron information.
 * <p>
 * This model contains information about an ICP neuron, which represents
 * staked ICP tokens used for governance and earning rewards.
 */
public class ICPNeuronInfo {

    private String neuronId;
    private String retrieveAtTimestampSeconds;
    private String neuronState;
    private String ageSeconds;
    private String dissolveDelaySeconds;
    private String votingPower;
    private String createdTimestampSeconds;
    private String stakeE8s;
    private String joinedCommunityFundTimestampSeconds;

    /**
     * Gets the neuron ID.
     *
     * @return the neuron ID
     */
    public String getNeuronId() {
        return neuronId;
    }

    /**
     * Sets the neuron ID.
     *
     * @param neuronId the neuron ID
     */
    public void setNeuronId(String neuronId) {
        this.neuronId = neuronId;
    }

    /**
     * Gets the timestamp when this info was retrieved.
     *
     * @return the retrieval timestamp in seconds
     */
    public String getRetrieveAtTimestampSeconds() {
        return retrieveAtTimestampSeconds;
    }

    /**
     * Sets the retrieval timestamp.
     *
     * @param retrieveAtTimestampSeconds the retrieval timestamp in seconds
     */
    public void setRetrieveAtTimestampSeconds(String retrieveAtTimestampSeconds) {
        this.retrieveAtTimestampSeconds = retrieveAtTimestampSeconds;
    }

    /**
     * Gets the neuron state (e.g., "LOCKED", "DISSOLVING", "DISSOLVED").
     *
     * @return the neuron state
     */
    public String getNeuronState() {
        return neuronState;
    }

    /**
     * Sets the neuron state.
     *
     * @param neuronState the neuron state
     */
    public void setNeuronState(String neuronState) {
        this.neuronState = neuronState;
    }

    /**
     * Gets the neuron's age in seconds.
     *
     * @return the age in seconds
     */
    public String getAgeSeconds() {
        return ageSeconds;
    }

    /**
     * Sets the neuron's age.
     *
     * @param ageSeconds the age in seconds
     */
    public void setAgeSeconds(String ageSeconds) {
        this.ageSeconds = ageSeconds;
    }

    /**
     * Gets the dissolve delay in seconds.
     *
     * @return the dissolve delay in seconds
     */
    public String getDissolveDelaySeconds() {
        return dissolveDelaySeconds;
    }

    /**
     * Sets the dissolve delay.
     *
     * @param dissolveDelaySeconds the dissolve delay in seconds
     */
    public void setDissolveDelaySeconds(String dissolveDelaySeconds) {
        this.dissolveDelaySeconds = dissolveDelaySeconds;
    }

    /**
     * Gets the neuron's voting power.
     *
     * @return the voting power
     */
    public String getVotingPower() {
        return votingPower;
    }

    /**
     * Sets the voting power.
     *
     * @param votingPower the voting power
     */
    public void setVotingPower(String votingPower) {
        this.votingPower = votingPower;
    }

    /**
     * Gets the timestamp when the neuron was created.
     *
     * @return the creation timestamp in seconds
     */
    public String getCreatedTimestampSeconds() {
        return createdTimestampSeconds;
    }

    /**
     * Sets the creation timestamp.
     *
     * @param createdTimestampSeconds the creation timestamp in seconds
     */
    public void setCreatedTimestampSeconds(String createdTimestampSeconds) {
        this.createdTimestampSeconds = createdTimestampSeconds;
    }

    /**
     * Gets the staked amount in e8s (1 ICP = 10^8 e8s).
     *
     * @return the stake in e8s
     */
    public String getStakeE8s() {
        return stakeE8s;
    }

    /**
     * Sets the stake.
     *
     * @param stakeE8s the stake in e8s
     */
    public void setStakeE8s(String stakeE8s) {
        this.stakeE8s = stakeE8s;
    }

    /**
     * Gets the timestamp when the neuron joined the community fund.
     *
     * @return the join timestamp in seconds, or null if not joined
     */
    public String getJoinedCommunityFundTimestampSeconds() {
        return joinedCommunityFundTimestampSeconds;
    }

    /**
     * Sets the community fund join timestamp.
     *
     * @param joinedCommunityFundTimestampSeconds the join timestamp in seconds
     */
    public void setJoinedCommunityFundTimestampSeconds(String joinedCommunityFundTimestampSeconds) {
        this.joinedCommunityFundTimestampSeconds = joinedCommunityFundTimestampSeconds;
    }

    @Override
    public String toString() {
        return "ICPNeuronInfo{"
                + "neuronId='" + neuronId + '\''
                + ", neuronState='" + neuronState + '\''
                + ", stakeE8s='" + stakeE8s + '\''
                + ", votingPower='" + votingPower + '\''
                + '}';
    }
}
