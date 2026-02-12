package com.taurushq.sdk.protect.client.service;

import com.google.common.base.Strings;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.StakingMapper;
import com.taurushq.sdk.protect.client.model.ADAStakePoolInfo;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.ETHValidatorInfo;
import com.taurushq.sdk.protect.client.model.FTMValidatorInfo;
import com.taurushq.sdk.protect.client.model.ICPNeuronInfo;
import com.taurushq.sdk.protect.client.model.NEARValidatorInfo;
import com.taurushq.sdk.protect.client.model.StakeAccountResult;
import com.taurushq.sdk.protect.client.model.XTZStakingRewards;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.StakingApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetADAStakePoolInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetETHValidatorsInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFTMValidatorInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetICPNeuronInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetNEARValidatorInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetStakeAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetXTZAddressStakingRewardsReply;

import java.time.OffsetDateTime;
import java.util.List;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving staking information across multiple blockchain networks.
 * <p>
 * This service provides access to validator information, stake accounts, and staking
 * rewards for various proof-of-stake blockchains including Cardano (ADA), Ethereum (ETH),
 * Fantom (FTM), Internet Computer (ICP), NEAR Protocol, Solana, and Tezos (XTZ).
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get Cardano stake pool info
 * ADAStakePoolInfo poolInfo = client.getStakingService()
 *     .getADAStakePoolInfo("mainnet", "pool1abc123...");
 *
 * // Get Ethereum validator info
 * List<ETHValidatorInfo> validators = client.getStakingService()
 *     .getETHValidatorsInfo("mainnet", Arrays.asList("validator1", "validator2"));
 *
 * // Get stake accounts with pagination
 * StakeAccountResult result = client.getStakingService()
 *     .getStakeAccounts("address-123", null, null, null);
 * }</pre>
 *
 * @see ADAStakePoolInfo
 * @see ETHValidatorInfo
 * @see StakeAccountResult
 */
public class StakingService {

    /**
     * The underlying OpenAPI client for staking operations.
     */
    private final StakingApi stakingApi;

    /**
     * Mapper for converting OpenAPI exceptions to SDK exceptions.
     */
    private final ApiExceptionMapper apiExceptionMapper;

    /**
     * Mapper for converting staking DTOs to domain models.
     */
    private final StakingMapper stakingMapper;

    /**
     * Instantiates a new Staking service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public StakingService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.stakingApi = new StakingApi(openApiClient);
        this.stakingMapper = StakingMapper.INSTANCE;
    }

    /**
     * Retrieves information about a Cardano stake pool.
     * <p>
     * Returns details including the pool's pledge, margin, fixed costs, and active stake.
     *
     * @param network     the network (e.g., "mainnet", "preprod")
     * @param stakePoolId the stake pool ID (Bech32 format, starting with "pool1")
     * @return the stake pool information
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if network or stakePoolId is null or empty
     */
    public ADAStakePoolInfo getADAStakePoolInfo(final String network, final String stakePoolId)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(stakePoolId), "stakePoolId cannot be null or empty");

        try {
            TgvalidatordGetADAStakePoolInfoReply reply = stakingApi.stakingServiceGetADAStakePoolInfo(
                    network, stakePoolId);
            return stakingMapper.fromADAStakePoolInfoReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves information about Ethereum validators.
     * <p>
     * Returns details including validator public keys, balances, and status.
     *
     * @param network the network (e.g., "mainnet", "goerli")
     * @param ids     the list of validator IDs to query
     * @return the list of validator information
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if network is null/empty or ids is null/empty
     */
    public List<ETHValidatorInfo> getETHValidatorsInfo(final String network, final List<String> ids)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkNotNull(ids, "ids cannot be null");
        checkArgument(!ids.isEmpty(), "ids cannot be empty");

        try {
            TgvalidatordGetETHValidatorsInfoReply reply = stakingApi.stakingServiceGetETHValidatorsInfo(
                    network, ids);
            return stakingMapper.fromETHValidatorInfoList(reply.getValidators());
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves information about a Fantom validator.
     * <p>
     * Returns details including the validator's stake amounts and status.
     *
     * @param network          the network (e.g., "mainnet")
     * @param validatorAddress the validator's address
     * @return the validator information
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if network or validatorAddress is null or empty
     */
    public FTMValidatorInfo getFTMValidatorInfo(final String network, final String validatorAddress)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(validatorAddress), "validatorAddress cannot be null or empty");

        try {
            TgvalidatordGetFTMValidatorInfoReply reply = stakingApi.stakingServiceGetFTMValidatorInfo(
                    network, validatorAddress);
            return stakingMapper.fromFTMValidatorInfoReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves information about an Internet Computer Protocol neuron.
     * <p>
     * Returns details including the neuron's stake, voting power, and dissolve delay.
     *
     * @param network  the network (e.g., "mainnet")
     * @param neuronId the neuron ID
     * @return the neuron information
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if network or neuronId is null or empty
     */
    public ICPNeuronInfo getICPNeuronInfo(final String network, final String neuronId)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(neuronId), "neuronId cannot be null or empty");

        try {
            TgvalidatordGetICPNeuronInfoReply reply = stakingApi.stakingServiceGetICPNeuronInfo(
                    network, neuronId);
            return stakingMapper.fromICPNeuronInfoReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves information about a NEAR Protocol validator.
     * <p>
     * Returns details including the validator's total staked balance and fee structure.
     *
     * @param network          the network (e.g., "mainnet", "testnet")
     * @param validatorAddress the validator's contract address
     * @return the validator information
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if network or validatorAddress is null or empty
     */
    public NEARValidatorInfo getNEARValidatorInfo(final String network, final String validatorAddress)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(validatorAddress), "validatorAddress cannot be null or empty");

        try {
            TgvalidatordGetNEARValidatorInfoReply reply = stakingApi.stakingServiceGetNEARValidatorInfo(
                    network, validatorAddress);
            return stakingMapper.fromNEARValidatorInfoReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves stake accounts with optional filtering.
     * <p>
     * Returns a paginated list of stake accounts that can be filtered by address,
     * account type, or account address.
     *
     * @param addressId      filter by associated address ID (optional)
     * @param accountType    filter by account type (optional)
     * @param accountAddress filter by on-chain account address (optional)
     * @param cursor         pagination cursor (optional, null for first page)
     * @return a paginated result containing stake accounts
     * @throws ApiException if the API call fails
     */
    public StakeAccountResult getStakeAccounts(final String addressId, final String accountType,
                                               final String accountAddress, final ApiRequestCursor cursor)
            throws ApiException {

        String cursorCurrentPage = null;
        String cursorPageRequest = null;
        String cursorPageSize = null;

        if (cursor != null) {
            cursorCurrentPage = cursor.getCurrentPage();
            cursorPageRequest = cursor.getPageRequest() != null ? cursor.getPageRequest().name() : null;
            cursorPageSize = String.valueOf(cursor.getPageSize());
        }

        try {
            TgvalidatordGetStakeAccountsReply reply = stakingApi.stakingServiceGetStakeAccounts(
                    addressId,
                    accountType,
                    accountAddress,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize
            );
            return stakingMapper.fromStakeAccountsReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Retrieves Tezos staking rewards for an address over a time period.
     * <p>
     * Returns the total rewards received by the specified address within the given
     * date range.
     *
     * @param network   the network (e.g., "mainnet")
     * @param addressId the address ID in the Taurus Protect system
     * @param from      the start date (optional)
     * @param to        the end date (optional)
     * @return the staking rewards information
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if network or addressId is null or empty
     */
    public XTZStakingRewards getXTZStakingRewards(final String network, final String addressId,
                                                   final OffsetDateTime from, final OffsetDateTime to)
            throws ApiException {
        checkArgument(!Strings.isNullOrEmpty(network), "network cannot be null or empty");
        checkArgument(!Strings.isNullOrEmpty(addressId), "addressId cannot be null or empty");

        try {
            TgvalidatordGetXTZAddressStakingRewardsReply reply = stakingApi.stakingServiceGetXTZAddressStakingRewards(
                    network, addressId, from, to);
            return stakingMapper.fromXTZStakingRewardsReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
