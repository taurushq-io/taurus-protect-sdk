package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.UserDeviceMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.UserDevicePairing;
import com.taurushq.sdk.protect.client.model.UserDevicePairingInfo;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.UserDeviceApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateUserDevicePairingReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUserDevicePairingInfo;
import com.taurushq.sdk.protect.openapi.model.UserDeviceServiceApproveUserDevicePairingBody;
import com.taurushq.sdk.protect.openapi.model.UserDeviceServiceStartUserDevicePairingBody;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for managing user device pairing in the Taurus Protect system.
 * <p>
 * Device pairing is used to register and authenticate user devices
 * for secure access to the Taurus Protect platform.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Create a new device pairing
 * UserDevicePairing pairing = client.getUserDeviceService().createPairing();
 *
 * // Get pairing status
 * UserDevicePairingInfo info = client.getUserDeviceService()
 *     .getPairingStatus(pairing.getPairingId(), "nonce-123");
 *
 * // Start the pairing process
 * client.getUserDeviceService().startPairing(pairing.getPairingId(), "nonce", "publicKey");
 *
 * // Approve the pairing
 * client.getUserDeviceService().approvePairing(pairing.getPairingId(), "nonce");
 * }</pre>
 *
 * @see UserDevicePairing
 * @see UserDevicePairingInfo
 */
public class UserDeviceService {

    private final UserDeviceApi userDeviceApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final UserDeviceMapper userDeviceMapper;

    /**
     * Creates a new UserDeviceService.
     *
     * @param apiClient          the API client for making HTTP requests
     * @param apiExceptionMapper the mapper for converting API exceptions
     * @throws NullPointerException if any parameter is null
     */
    public UserDeviceService(final ApiClient apiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(apiClient, "apiClient must not be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper must not be null");
        this.userDeviceApi = new UserDeviceApi(apiClient);
        this.apiExceptionMapper = apiExceptionMapper;
        this.userDeviceMapper = UserDeviceMapper.INSTANCE;
    }

    /**
     * Creates a new device pairing.
     *
     * @return the created pairing with its ID
     * @throws ApiException if the API call fails
     */
    public UserDevicePairing createPairing() throws ApiException {
        try {
            TgvalidatordCreateUserDevicePairingReply reply = userDeviceApi.userDeviceServiceCreateUserDevicePairing();
            return userDeviceMapper.fromCreateDTO(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Gets the status of a device pairing.
     *
     * @param pairingId the pairing ID
     * @param nonce     the nonce for verification
     * @return the pairing status and info
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if pairingId or nonce is null or empty
     */
    public UserDevicePairingInfo getPairingStatus(final String pairingId, final String nonce) throws ApiException {
        checkArgument(pairingId != null && !pairingId.isEmpty(), "pairingId must not be null or empty");
        checkArgument(nonce != null && !nonce.isEmpty(), "nonce must not be null or empty");
        try {
            TgvalidatordUserDevicePairingInfo reply = userDeviceApi.userDeviceServiceGetUserDevicePairingStatus(
                    pairingId, nonce);
            return userDeviceMapper.fromInfoDTO(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Starts the device pairing process.
     *
     * @param pairingId the pairing ID
     * @param nonce     the nonce for verification
     * @param publicKey the public key for the device
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if pairingId is null or empty
     */
    public void startPairing(final String pairingId, final String nonce, final String publicKey) throws ApiException {
        checkArgument(pairingId != null && !pairingId.isEmpty(), "pairingId must not be null or empty");
        try {
            UserDeviceServiceStartUserDevicePairingBody body = new UserDeviceServiceStartUserDevicePairingBody();
            body.setNonce(nonce);
            body.setPublicKey(publicKey);
            userDeviceApi.userDeviceServiceStartUserDevicePairing(pairingId, body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Approves a device pairing.
     *
     * @param pairingId the pairing ID
     * @param nonce     the nonce for verification
     * @throws ApiException             if the API call fails
     * @throws IllegalArgumentException if pairingId or nonce is null or empty
     */
    public void approvePairing(final String pairingId, final String nonce) throws ApiException {
        checkArgument(pairingId != null && !pairingId.isEmpty(), "pairingId must not be null or empty");
        checkArgument(nonce != null && !nonce.isEmpty(), "nonce must not be null or empty");
        try {
            UserDeviceServiceApproveUserDevicePairingBody body = new UserDeviceServiceApproveUserDevicePairingBody();
            body.setNonce(nonce);
            userDeviceApi.userDeviceServiceApproveUserDevicePairing(pairingId, body);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
