package com.taurushq.sdk.protect.client.service;

import com.google.protobuf.ByteString;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.rulescontainer.AddressWhitelistingLine;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSource;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceInternalWallet;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleSourceType;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.Test;

import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.spec.ECGenParameterSpec;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Guards the fail-closed check on untypeable whitelisting source cells.
 *
 * <p>A container signed by a validatord newer than this SDK can carry a source cell the
 * decoder could only preserve verbatim. Such a line simply failed to match the wallet
 * path, so verification fell through to the container default thresholds — approving the
 * address against a weaker quorum than its governance line demands, silently and with no
 * error. {@code getApplicableThresholds} now refuses instead.
 *
 * <p>The predicate is asserted directly rather than through the enclosing private method:
 * this project forbids Mockito, so a service test cannot stub the API client needed to
 * drive a full verification. The three other SDKs cover the abort behaviourally.
 */
class WhitelistedAddressServiceUntypedSourceTest {

    private static WhitelistedAddressService newService() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        KeyPair pair = generator.generateKeyPair();
        return new WhitelistedAddressService(
                new ApiClient(),
                new ApiExceptionMapper(),
                Collections.singletonList(pair.getPublic()),
                1);
    }

    private static AddressWhitelistingLine lineWithCell(RuleSource cell) {
        AddressWhitelistingLine line = new AddressWhitelistingLine();
        List<RuleSource> cells = new ArrayList<>();
        cells.add(cell);
        line.setCells(cells);
        return line;
    }

    @Test
    void rawSourceCellIsReportedAsUntyped() throws Exception {
        RuleSource raw = new RuleSource();
        raw.setRaw(ByteString.copyFrom(new byte[]{0x08, 0x63}));

        assertTrue(newService().lineHasUntypedSource(lineWithCell(raw)),
                "a verbatim-preserved source cell must be reported as untyped");
    }

    @Test
    void typedSourceCellIsNotReportedAsUntyped() throws Exception {
        RuleSource typed = new RuleSource();
        typed.setType(RuleSourceType.RuleSourceInternalWallet);
        RuleSourceInternalWallet wallet = new RuleSourceInternalWallet();
        wallet.setPath("m/44/60/0");
        typed.setInternalWallet(wallet);

        assertFalse(newService().lineHasUntypedSource(lineWithCell(typed)),
                "a fully typed source cell must not trip the fail-closed check");
    }

    @Test
    void emptyRawIsNotReportedAsUntyped() throws Exception {
        RuleSource emptyRaw = new RuleSource();
        emptyRaw.setType(RuleSourceType.RuleSourceAny);
        emptyRaw.setRaw(ByteString.EMPTY);

        assertFalse(newService().lineHasUntypedSource(lineWithCell(emptyRaw)),
                "an empty raw buffer carries nothing this SDK failed to interpret");
    }

    @Test
    void lineWithNoCellsIsNotReportedAsUntyped() throws Exception {
        AddressWhitelistingLine empty = new AddressWhitelistingLine();
        empty.setCells(new ArrayList<>());

        assertFalse(newService().lineHasUntypedSource(empty));
        assertFalse(newService().lineHasUntypedSource(null));
    }
}
