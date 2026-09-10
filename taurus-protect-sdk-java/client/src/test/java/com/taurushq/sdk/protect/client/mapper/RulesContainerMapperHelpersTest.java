package com.taurushq.sdk.protect.client.mapper;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.util.Arrays;
import java.util.Base64;
import java.util.Collections;
import java.util.List;

import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.proto.v1.RequestReply;
import org.junit.jupiter.api.Test;

/**
 * Coverage for the {@link RulesContainerMapper} encoder helper methods: enum-name
 * resolution (which tolerates unknown/null names rather than throwing so a newer
 * schema does not break encoding), role-name filtering, and the base64 wrapper.
 */
class RulesContainerMapperHelpersTest {

    private static final RulesContainerMapper M = RulesContainerMapper.INSTANCE;

    @Test
    void enumValuePassesUnknownNumbersThroughAndRejectsBogusNames() {
        RequestReply.RulesContainer.ColumnType known = RequestReply.RulesContainer.ColumnType.RuleFiatAmount;
        assertEquals(known.getNumber(), M.enumValue(RequestReply.RulesContainer.ColumnType.class, "RuleFiatAmount"));
        // A value newer than this SDK arrives as its decimal string and keeps its number,
        // so re-encoding does not rewrite the column to the zero value.
        assertEquals(77, M.enumValue(RequestReply.RulesContainer.ColumnType.class, "77"));
        // An absent name is the proto default; a name that is neither known nor numeric is
        // a caller error and must fail loudly rather than silently becoming 0.
        assertEquals(0, M.enumValue(RequestReply.RulesContainer.ColumnType.class, null));
        assertEquals(0, M.enumValue(RequestReply.RulesContainer.ColumnType.class, ""));
        assertThrows(IllegalArgumentException.class,
                () -> M.enumValue(RequestReply.RulesContainer.ColumnType.class, "TotallyBogus"));
    }

    @Test
    void rolesFromStringsPassesUnknownRoleNumbersThrough() {
        // Dropping a role would strip a privilege from data SuperAdmins are about to sign,
        // so a role newer than this SDK keeps its number.
        List<Integer> roles = M.rolesFromStrings(Arrays.asList("SUPERADMIN", "9", "HSMSLOT"));
        assertEquals(Arrays.asList(
                RequestReply.Role.SUPERADMIN.getNumber(), 9, RequestReply.Role.HSMSLOT.getNumber()), roles);
        assertTrue(M.rolesFromStrings(null).isEmpty());
        assertThrows(IllegalArgumentException.class,
                () -> M.rolesFromStrings(Arrays.asList("BOGUS_ROLE")));
    }

    @Test
    void toBase64StringMatchesBase64OfBytes() {
        DecodedRulesContainer c = new DecodedRulesContainer();
        RuleUser u = new RuleUser();
        u.setId("u1");
        u.setRoles(Collections.singletonList("SUPERADMIN"));
        c.setUsers(Collections.singletonList(u));
        c.setMinimumDistinctUserSignatures(1);
        assertEquals(Base64.getEncoder().encodeToString(M.toBytes(c)), M.toBase64String(c));
    }
}
