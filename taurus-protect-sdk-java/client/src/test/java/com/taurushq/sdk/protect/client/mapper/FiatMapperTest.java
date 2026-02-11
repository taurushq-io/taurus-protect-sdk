package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.FiatProvider;
import com.taurushq.sdk.protect.client.model.FiatProviderAccount;
import com.taurushq.sdk.protect.client.model.FiatProviderAccountResult;
import com.taurushq.sdk.protect.client.model.FiatProviderCounterpartyAccount;
import com.taurushq.sdk.protect.client.model.FiatProviderCounterpartyAccountResult;
import com.taurushq.sdk.protect.client.model.FiatProviderOperation;
import com.taurushq.sdk.protect.client.model.FiatProviderOperationResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFiatProvider;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFiatProviderAccount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFiatProviderCounterpartyAccount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFiatProviderOperation;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderCounterpartyAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderOperationsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FiatMapperTest {

    @Test
    void fromProviderDTO_mapsFields() {
        TgvalidatordFiatProvider dto = new TgvalidatordFiatProvider();
        dto.setProvider("circle");
        dto.setLabel("Test Provider");
        dto.setBaseCurrencyValuation("1000.00");

        FiatProvider result = FiatMapper.INSTANCE.fromProviderDTO(dto);

        assertNotNull(result);
        assertEquals("circle", result.getProvider());
        assertEquals("Test Provider", result.getLabel());
        assertEquals("1000.00", result.getBaseCurrencyValuation());
    }

    @Test
    void fromProviderDTO_handlesNull() {
        FiatProvider result = FiatMapper.INSTANCE.fromProviderDTO(null);
        assertNull(result);
    }

    @Test
    void fromProviderDTOList_mapsList() {
        TgvalidatordFiatProvider dto1 = new TgvalidatordFiatProvider();
        dto1.setProvider("circle");

        TgvalidatordFiatProvider dto2 = new TgvalidatordFiatProvider();
        dto2.setProvider("bank");

        List<FiatProvider> result = FiatMapper.INSTANCE.fromProviderDTOList(
                Arrays.asList(dto1, dto2));

        assertNotNull(result);
        assertEquals(2, result.size());
    }

    @Test
    void fromAccountDTO_mapsFields() {
        TgvalidatordFiatProviderAccount dto = new TgvalidatordFiatProviderAccount();
        dto.setId("account-123");
        dto.setProvider("provider-1");
        dto.setLabel("My Account");
        dto.setAccountType("CHECKING");

        FiatProviderAccount result = FiatMapper.INSTANCE.fromAccountDTO(dto);

        assertNotNull(result);
        assertEquals("account-123", result.getId());
        assertEquals("provider-1", result.getProvider());
        assertEquals("My Account", result.getLabel());
        assertEquals("CHECKING", result.getAccountType());
    }

    @Test
    void fromAccountDTO_handlesNull() {
        FiatProviderAccount result = FiatMapper.INSTANCE.fromAccountDTO(null);
        assertNull(result);
    }

    @Test
    void fromAccountsReply_mapsAccountsAndCursor() {
        TgvalidatordFiatProviderAccount account = new TgvalidatordFiatProviderAccount();
        account.setId("account-123");

        TgvalidatordResponseCursor cursor = new TgvalidatordResponseCursor();
        cursor.setCurrentPage("page-1");

        TgvalidatordGetFiatProviderAccountsReply reply = new TgvalidatordGetFiatProviderAccountsReply();
        reply.setResult(Arrays.asList(account));
        reply.setCursor(cursor);

        FiatProviderAccountResult result = FiatMapper.INSTANCE.fromAccountsReply(reply);

        assertNotNull(result);
        assertNotNull(result.getAccounts());
        assertEquals(1, result.getAccounts().size());
        assertNotNull(result.getCursor());
    }

    @Test
    void fromCounterpartyAccountDTO_mapsFields() {
        TgvalidatordFiatProviderCounterpartyAccount dto = new TgvalidatordFiatProviderCounterpartyAccount();
        dto.setId("cp-account-123");
        dto.setProvider("provider-1");
        dto.setLabel("Counterparty");
        dto.setCounterpartyID("cp-456");

        FiatProviderCounterpartyAccount result = FiatMapper.INSTANCE.fromCounterpartyAccountDTO(dto);

        assertNotNull(result);
        assertEquals("cp-account-123", result.getId());
        assertEquals("provider-1", result.getProvider());
        assertEquals("Counterparty", result.getLabel());
        assertEquals("cp-456", result.getCounterpartyId());
    }

    @Test
    void fromCounterpartyAccountDTO_handlesNull() {
        FiatProviderCounterpartyAccount result = FiatMapper.INSTANCE.fromCounterpartyAccountDTO(null);
        assertNull(result);
    }

    @Test
    void fromCounterpartyAccountsReply_mapsAccountsAndCursor() {
        TgvalidatordFiatProviderCounterpartyAccount account = new TgvalidatordFiatProviderCounterpartyAccount();
        account.setId("cp-account-123");

        TgvalidatordResponseCursor cursor = new TgvalidatordResponseCursor();
        cursor.setCurrentPage("page-1");

        TgvalidatordGetFiatProviderCounterpartyAccountsReply reply =
                new TgvalidatordGetFiatProviderCounterpartyAccountsReply();
        reply.setResult(Arrays.asList(account));
        reply.setCursor(cursor);

        FiatProviderCounterpartyAccountResult result = FiatMapper.INSTANCE.fromCounterpartyAccountsReply(reply);

        assertNotNull(result);
        assertNotNull(result.getAccounts());
        assertEquals(1, result.getAccounts().size());
        assertNotNull(result.getCursor());
    }

    @Test
    void fromOperationDTO_mapsFields() {
        TgvalidatordFiatProviderOperation dto = new TgvalidatordFiatProviderOperation();
        dto.setId("op-123");
        dto.setProvider("provider-1");
        dto.setStatus("COMPLETED");

        FiatProviderOperation result = FiatMapper.INSTANCE.fromOperationDTO(dto);

        assertNotNull(result);
        assertEquals("op-123", result.getId());
        assertEquals("provider-1", result.getProvider());
        assertEquals("COMPLETED", result.getStatus());
    }

    @Test
    void fromOperationDTO_handlesNull() {
        FiatProviderOperation result = FiatMapper.INSTANCE.fromOperationDTO(null);
        assertNull(result);
    }

    @Test
    void fromOperationsReply_mapsOperationsAndCursor() {
        TgvalidatordFiatProviderOperation operation = new TgvalidatordFiatProviderOperation();
        operation.setId("op-123");

        TgvalidatordResponseCursor cursor = new TgvalidatordResponseCursor();
        cursor.setCurrentPage("page-1");

        TgvalidatordGetFiatProviderOperationsReply reply = new TgvalidatordGetFiatProviderOperationsReply();
        reply.setResult(Arrays.asList(operation));
        reply.setCursor(cursor);

        FiatProviderOperationResult result = FiatMapper.INSTANCE.fromOperationsReply(reply);

        assertNotNull(result);
        assertNotNull(result.getOperations());
        assertEquals(1, result.getOperations().size());
        assertNotNull(result.getCursor());
    }
}
