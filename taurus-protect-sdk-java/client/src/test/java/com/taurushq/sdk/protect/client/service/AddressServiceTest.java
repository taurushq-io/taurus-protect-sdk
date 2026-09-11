package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.Security;
import java.security.spec.ECGenParameterSpec;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

class AddressServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private RulesContainerCache rulesContainerCache;
    private AddressService addressService;

    @BeforeAll
    static void setUpProvider() {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    @BeforeEach
    void setUp() throws Exception {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();

        // Generate a real EC P-256 key to satisfy GovernanceRuleService constructor
        KeyPairGenerator kpg = KeyPairGenerator.getInstance("EC");
        kpg.initialize(new ECGenParameterSpec("secp256r1"));
        KeyPair kp = kpg.generateKeyPair();
        PublicKey publicKey = kp.getPublic();

        GovernanceRuleService governanceRuleService = new GovernanceRuleService(
                apiClient, apiExceptionMapper, Collections.singletonList(publicKey), 1);
        rulesContainerCache = new RulesContainerCache(governanceRuleService);
        addressService = new AddressService(apiClient, apiExceptionMapper, rulesContainerCache);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new AddressService(null, apiExceptionMapper, rulesContainerCache));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new AddressService(apiClient, null, rulesContainerCache));
    }

    @Test
    void constructor_throwsOnNullRulesContainerCache() {
        assertThrows(NullPointerException.class, () ->
                new AddressService(apiClient, apiExceptionMapper, null));
    }

    @Test
    void createAddress_throwsOnNullRequest() {
        assertThrows(NullPointerException.class, () ->
                addressService.createAddress(null));
    }

    @Test
    void getAddress_throwsOnZeroId() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddress(0));
    }

    @Test
    void getAddress_throwsOnNegativeId() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddress(-1));
    }

    @Test
    void getAddresses_throwsOnZeroWalletId() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddresses(0, 10, 0));
    }

    @Test
    void getAddresses_throwsOnZeroLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddresses(1, 0, 0));
    }

    @Test
    void getAddresses_throwsOnNegativeOffset() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.getAddresses(1, 10, -1));
    }

    @Test
    void createAddressAttribute_throwsOnZeroAddressId() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.createAddressAttribute(0, "key", "value"));
    }

    @Test
    void createAddressAttribute_throwsOnEmptyKey() {
        assertThrows(IllegalArgumentException.class, () ->
                addressService.createAddressAttribute(1, "", "value"));
    }


    // ------------------------------------------------------------------------
    // The address verification seam.
    //
    // These drive AddressService.verifiedAddress / verifiedAddresses directly with
    // hand-built rows. That is deliberate and it is why the seam is package-private:
    // test scope here is JUnit only (no Mockito, no WireMock, no HTTP stub), so the
    // full createAddress/getAddress/getAssetAddresses round trip is NOT drivable and
    // these tests do not claim to cover it. What they do cover is the invariant every
    // one of those paths now routes through — including AssetService, which calls the
    // same static seam rather than carrying its own copy.
    // ------------------------------------------------------------------------

    /**
     * A container source that fails the test if it is consulted at all. Used for the
     * arms that must reach a decision without a governance fetch.
     */
    private static final AddressService.RulesContainerSource NO_FETCH_EXPECTED = () -> {
        throw new AssertionError("the rules container must not be fetched here");
    };

    private static Address address(long id, String value, String signature, String status) {
        Address address = new Address();
        address.setId(id);
        address.setAddress(value);
        address.setSignature(signature);
        address.setStatus(status);
        return address;
    }

    private static DecodedRulesContainer containerWithHsmKey(PublicKey hsmKey) {
        DecodedRulesContainer container = new DecodedRulesContainer();
        RuleUser hsm = new RuleUser();
        hsm.setId("hsm-slot-1");
        // setPublicKey, not setPublicKeyPem: findHsmPublicKey reads getPublicKey(), so a
        // PEM-only user leaves the container with no HSM key and every one of these
        // tests would pass on "no HSMSLOT user found" instead of on the real check.
        hsm.setPublicKey(hsmKey);
        hsm.setRoles(Collections.singletonList("HSMSLOT"));
        container.setUsers(new ArrayList<>(Arrays.asList(hsm)));
        return container;
    }

    private static KeyPair hsmKeyPair() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        return generator.generateKeyPair();
    }

    /**
     * Positive control for the two tests below: a genuinely signed address verifies.
     * <p>
     * Without it a tamper test proves nothing — a broken fixture (no HSM key resolved,
     * wrong signing algorithm) also throws, and the test would stay green against an
     * implementation that rejects everything.
     */
    @Test
    void verifiedAddress_signatureMatchesTheAddress_isReturned() throws Exception {
        KeyPair hsm = hsmKeyPair();
        String value = "0x00000000000000000000000000000000000000aa";
        String signature = CryptoTPV1.calculateBase64Signature(
                hsm.getPrivate(), value.getBytes(StandardCharsets.UTF_8));
        Address signed = address(7L, value, signature, "confirmed");
        DecodedRulesContainer container = containerWithHsmKey(hsm.getPublic());

        Address returned = AddressService.verifiedAddress(signed, () -> container);

        assertSame(signed, returned);
    }

    /**
     * The signature covers the address STRING, so swapping the string after signing is
     * how a response-controlling server redirects a deposit. It must not verify.
     */
    @Test
    void verifiedAddress_signaturePresentButAddressTampered_throwsIntegrityException() throws Exception {
        KeyPair hsm = hsmKeyPair();
        String signature = CryptoTPV1.calculateBase64Signature(
                hsm.getPrivate(),
                "0x00000000000000000000000000000000000000aa".getBytes(StandardCharsets.UTF_8));
        // Same signature, attacker's address.
        Address tampered = address(7L, "0x00000000000000000000000000000000000000bb",
                signature, "confirmed");
        DecodedRulesContainer container = containerWithHsmKey(hsm.getPublic());

        IntegrityException e = assertThrows(IntegrityException.class,
                () -> AddressService.verifiedAddress(tampered, () -> container));
        assertTrue(e.getMessage().contains("7"), e.getMessage());
    }

    /**
     * The createAddress hole. A creation reply carrying an address string and NO
     * signature used to be handed straight to the caller, in the same type the getters
     * verify. There is nothing to check, so the only safe answer is to refuse — not to
     * pass the server's string through because no check was possible.
     * <p>
     * Uses {@link #NO_FETCH_EXPECTED}: the refusal must be reached without a governance
     * round trip, since an unsigned address is not a transient condition to retry.
     */
    @Test
    void verifiedAddress_signatureAbsentAndAddressPresent_isRefused() {
        Address unsigned = address(99L, "0x00000000000000000000000000000000000000bb",
                null, "created");

        IntegrityException e = assertThrows(IntegrityException.class,
                () -> AddressService.verifiedAddress(unsigned, NO_FETCH_EXPECTED));

        // The caller has to be able to act on this: which row, what state it was in,
        // and what to do next.
        assertTrue(e.getMessage().contains("99"), e.getMessage());
        assertTrue(e.getMessage().contains("created"), e.getMessage());
        assertTrue(e.getMessage().contains("getAddress"), e.getMessage());

        // An empty signature is the same case as an absent one.
        assertThrows(IntegrityException.class, () -> AddressService.verifiedAddress(
                address(99L, "0x00000000000000000000000000000000000000bb", "", "created"),
                NO_FETCH_EXPECTED));
    }

    /**
     * Asynchronous creation: the row exists with status {@code creating} and no address
     * string yet. There is no destination to misuse and nothing has been signed, so
     * this must come back rather than being refused — and without a container fetch.
     */
    @Test
    void verifiedAddress_emptyAddressWithCreatingStatus_isReturned() {
        Address creating = address(123L, "", null, "creating");

        Address returned = assertDoesNotThrow(
                () -> AddressService.verifiedAddress(creating, NO_FETCH_EXPECTED));

        assertSame(creating, returned);
        // Java could not express this case at all before: Address had no status field
        // and AddressMapper silently dropped the DTO's, so "not signed yet" and
        // "signature stripped" were indistinguishable.
        assertEquals("creating", returned.getStatus());
    }

    /**
     * A null address string is the same case as an empty one — some replies omit the
     * field entirely rather than sending "".
     */
    @Test
    void verifiedAddress_nullAddressString_isReturned() {
        Address creating = address(124L, null, null, "creating");

        assertSame(creating, assertDoesNotThrow(
                () -> AddressService.verifiedAddress(creating, NO_FETCH_EXPECTED)));
    }

    /**
     * The list seam is the same invariant applied per row: one bad row aborts the page.
     * Addresses are destinations, so excluding-and-reporting (what the whitelist list
     * paths do) would leave the caller choosing from a silently filtered set.
     */
    @Test
    void verifiedAddresses_oneUnsignedRowAbortsThePage() throws Exception {
        KeyPair hsm = hsmKeyPair();
        String value = "0x00000000000000000000000000000000000000aa";
        String signature = CryptoTPV1.calculateBase64Signature(
                hsm.getPrivate(), value.getBytes(StandardCharsets.UTF_8));
        List<Address> page = Arrays.asList(
                address(1L, value, signature, "confirmed"),
                address(2L, "0x00000000000000000000000000000000000000bb", null, "created"));
        DecodedRulesContainer container = containerWithHsmKey(hsm.getPublic());

        IntegrityException e = assertThrows(IntegrityException.class,
                () -> AddressService.verifiedAddresses(page, () -> container));
        assertTrue(e.getMessage().contains("2"), e.getMessage());
    }

    /**
     * A page of in-flight creations needs no rules container, so none is fetched. This
     * is what keeps createAddress usable against an asynchronous wallet even when the
     * governance read would fail.
     */
    @Test
    void verifiedAddresses_pageOfCreatingRowsFetchesNoContainer() {
        List<Address> page = Arrays.asList(
                address(1L, "", null, "creating"),
                address(2L, null, null, "creating"));

        assertSame(page, assertDoesNotThrow(
                () -> AddressService.verifiedAddresses(page, NO_FETCH_EXPECTED)));
    }

    @Test
    void verifiedAddresses_nullAndEmptyPagesFetchNoContainer() {
        assertDoesNotThrow(() -> AddressService.verifiedAddresses(null, NO_FETCH_EXPECTED));
        assertDoesNotThrow(() -> AddressService.verifiedAddresses(
                Collections.<Address>emptyList(), NO_FETCH_EXPECTED));
    }
}
