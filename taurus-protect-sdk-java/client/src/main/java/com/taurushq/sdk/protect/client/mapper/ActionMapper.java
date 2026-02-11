package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Action;
import com.taurushq.sdk.protect.client.model.ActionAmount;
import com.taurushq.sdk.protect.client.model.ActionAttribute;
import com.taurushq.sdk.protect.client.model.ActionComparator;
import com.taurushq.sdk.protect.client.model.ActionDestination;
import com.taurushq.sdk.protect.client.model.ActionEnvelope;
import com.taurushq.sdk.protect.client.model.ActionSource;
import com.taurushq.sdk.protect.client.model.ActionTarget;
import com.taurushq.sdk.protect.client.model.ActionTask;
import com.taurushq.sdk.protect.client.model.ActionTrail;
import com.taurushq.sdk.protect.client.model.ActionTrigger;
import com.taurushq.sdk.protect.client.model.TargetAddress;
import com.taurushq.sdk.protect.client.model.TargetWallet;
import com.taurushq.sdk.protect.client.model.TaskNotification;
import com.taurushq.sdk.protect.client.model.TaskTransfer;
import com.taurushq.sdk.protect.client.model.TriggerBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAction;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionAmount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionAttribute;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionDestination;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionEnvelopeTrail;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionSource;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordActionEnvelope;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting Action-related DTOs to client model objects.
 */
@Mapper
public interface ActionMapper {

    ActionMapper INSTANCE = Mappers.getMapper(ActionMapper.class);

    @Mapping(source = "lastcheckeddate", target = "lastCheckedDate")
    ActionEnvelope fromDTO(TgvalidatordActionEnvelope dto);

    List<ActionEnvelope> fromDTOList(List<TgvalidatordActionEnvelope> dtos);

    Action fromActionDTO(TgvalidatordAction dto);

    ActionTrigger fromTriggerDTO(com.taurushq.sdk.protect.openapi.model.ActionTrigger dto);

    TriggerBalance fromTriggerBalanceDTO(com.taurushq.sdk.protect.openapi.model.TriggerBalance dto);

    ActionTarget fromTargetDTO(com.taurushq.sdk.protect.openapi.model.ActionTarget dto);

    TargetAddress fromTargetAddressDTO(com.taurushq.sdk.protect.openapi.model.TargetAddress dto);

    TargetWallet fromTargetWalletDTO(com.taurushq.sdk.protect.openapi.model.TargetWallet dto);

    ActionComparator fromComparatorDTO(com.taurushq.sdk.protect.openapi.model.ActionComparator dto);

    ActionAmount fromAmountDTO(TgvalidatordActionAmount dto);

    ActionTask fromTaskDTO(com.taurushq.sdk.protect.openapi.model.ActionTask dto);

    List<ActionTask> fromTaskDTOList(List<com.taurushq.sdk.protect.openapi.model.ActionTask> dtos);

    TaskTransfer fromTransferDTO(com.taurushq.sdk.protect.openapi.model.TaskTransfer dto);

    TaskNotification fromNotificationDTO(com.taurushq.sdk.protect.openapi.model.TaskNotification dto);

    ActionSource fromSourceDTO(TgvalidatordActionSource dto);

    ActionDestination fromDestinationDTO(TgvalidatordActionDestination dto);

    ActionAttribute fromAttributeDTO(TgvalidatordActionAttribute dto);

    List<ActionAttribute> fromAttributeDTOList(List<TgvalidatordActionAttribute> dtos);

    ActionTrail fromTrailDTO(TgvalidatordActionEnvelopeTrail dto);

    List<ActionTrail> fromTrailDTOList(List<TgvalidatordActionEnvelopeTrail> dtos);
}
