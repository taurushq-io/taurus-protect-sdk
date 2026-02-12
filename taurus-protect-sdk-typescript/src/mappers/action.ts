/**
 * Action mapper functions for converting OpenAPI DTOs to domain models.
 */

import type {
  Action,
  ActionAmount,
  ActionAttribute,
  ActionComparator,
  ActionDestination,
  ActionEnvelope,
  ActionSource,
  ActionTarget,
  ActionTask,
  ActionTrail,
  ActionTrigger,
  TargetAddress,
  TargetWallet,
  TaskNotification,
  TaskTransfer,
  TriggerBalance,
} from '../models/action';
import { safeBool, safeDate, safeMap, safeString } from './base';

/**
 * Maps a target address DTO to a TargetAddress domain model.
 */
export function targetAddressFromDto(dto: unknown): TargetAddress | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
    addressId: safeString(d.addressID ?? d.addressId ?? d.address_id),
  };
}

/**
 * Maps a target wallet DTO to a TargetWallet domain model.
 */
export function targetWalletFromDto(dto: unknown): TargetWallet | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
    walletId: safeString(d.walletID ?? d.walletId ?? d.wallet_id),
  };
}

/**
 * Maps an action target DTO to an ActionTarget domain model.
 */
export function actionTargetFromDto(dto: unknown): ActionTarget | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
    address: targetAddressFromDto(d.address),
    wallet: targetWalletFromDto(d.wallet),
  };
}

/**
 * Maps an action comparator DTO to an ActionComparator domain model.
 */
export function actionComparatorFromDto(dto: unknown): ActionComparator | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
  };
}

/**
 * Maps an action amount DTO to an ActionAmount domain model.
 */
export function actionAmountFromDto(dto: unknown): ActionAmount | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
    cryptoAmount: safeString(d.cryptoAmount ?? d.crypto_amount),
  };
}

/**
 * Maps a trigger balance DTO to a TriggerBalance domain model.
 */
export function triggerBalanceFromDto(dto: unknown): TriggerBalance | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    target: actionTargetFromDto(d.target),
    comparator: actionComparatorFromDto(d.comparator),
    amount: actionAmountFromDto(d.amount),
  };
}

/**
 * Maps an action trigger DTO to an ActionTrigger domain model.
 */
export function actionTriggerFromDto(dto: unknown): ActionTrigger | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
    balance: triggerBalanceFromDto(d.balance),
  };
}

/**
 * Maps an action source DTO to an ActionSource domain model.
 */
export function actionSourceFromDto(dto: unknown): ActionSource | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
    addressId: safeString(d.addressID ?? d.addressId ?? d.address_id),
    walletId: safeString(d.walletID ?? d.walletId ?? d.wallet_id),
  };
}

/**
 * Maps an action destination DTO to an ActionDestination domain model.
 */
export function actionDestinationFromDto(dto: unknown): ActionDestination | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
    addressId: safeString(d.addressID ?? d.addressId ?? d.address_id),
    whitelistedAddressId: safeString(
      d.whitelistedAddressID ?? d.whitelistedAddressId ?? d.whitelisted_address_id
    ),
    walletId: safeString(d.walletID ?? d.walletId ?? d.wallet_id),
  };
}

/**
 * Maps a task transfer DTO to a TaskTransfer domain model.
 */
export function taskTransferFromDto(dto: unknown): TaskTransfer | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    from: actionSourceFromDto(d.from),
    to: actionDestinationFromDto(d.to),
    amount: actionAmountFromDto(d.amount),
    topUp: safeBool(d.topUp ?? d.top_up),
    useAllFunds: safeBool(d.useAllFunds ?? d.use_all_funds),
  };
}

/**
 * Maps a task notification DTO to a TaskNotification domain model.
 */
export function taskNotificationFromDto(dto: unknown): TaskNotification | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const emailAddresses = d.emailAddresses ?? d.email_addresses;
  return {
    emailAddresses: Array.isArray(emailAddresses)
      ? emailAddresses.map((e) => String(e))
      : undefined,
    notificationMessage: safeString(d.notificationMessage ?? d.notification_message),
    numberOfReminders: safeString(d.numberOfReminders ?? d.number_of_reminders),
  };
}

/**
 * Maps an action task DTO to an ActionTask domain model.
 */
export function actionTaskFromDto(dto: unknown): ActionTask | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    kind: safeString(d.kind),
    transfer: taskTransferFromDto(d.transfer),
    notification: taskNotificationFromDto(d.notification),
  };
}

/**
 * Maps an action DTO to an Action domain model.
 */
export function actionFromDto(dto: unknown): Action | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const tasks = d.tasks;
  return {
    trigger: actionTriggerFromDto(d.trigger),
    tasks: Array.isArray(tasks) ? safeMap(tasks, actionTaskFromDto) : undefined,
  };
}

/**
 * Maps an action attribute DTO to an ActionAttribute domain model.
 */
export function actionAttributeFromDto(dto: unknown): ActionAttribute | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    tenantId: safeString(d.tenantId ?? d.tenant_id),
    key: safeString(d.key),
    value: safeString(d.value),
    contentType: safeString(d.contentType ?? d.content_type),
  };
}

/**
 * Maps an action trail DTO to an ActionTrail domain model.
 */
export function actionTrailFromDto(dto: unknown): ActionTrail | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    action: safeString(d.action),
    comment: safeString(d.comment),
    date: safeDate(d.date),
    actionStatus: safeString(d.actionStatus ?? d.action_status),
  };
}

/**
 * Maps an action envelope DTO to an ActionEnvelope domain model.
 */
export function actionEnvelopeFromDto(dto: unknown): ActionEnvelope | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  const attributes = d.attributes;
  const trails = d.trails;

  return {
    id: safeString(d.id),
    tenantId: safeString(d.tenantId ?? d.tenant_id),
    label: safeString(d.label),
    action: actionFromDto(d.action),
    status: safeString(d.status),
    creationDate: safeDate(d.creationDate ?? d.creation_date),
    updateDate: safeDate(d.updateDate ?? d.update_date),
    lastCheckedDate: safeDate(d.lastcheckeddate ?? d.lastCheckedDate ?? d.last_checked_date),
    autoApprove: safeBool(d.autoApprove ?? d.auto_approve),
    attributes: Array.isArray(attributes)
      ? safeMap(attributes, actionAttributeFromDto)
      : undefined,
    trails: Array.isArray(trails) ? safeMap(trails, actionTrailFromDto) : undefined,
  };
}

/**
 * Maps an array of action envelope DTOs to ActionEnvelope domain models.
 */
export function actionEnvelopesFromDto(
  dtos: unknown[] | null | undefined
): ActionEnvelope[] {
  return safeMap(dtos, actionEnvelopeFromDto);
}
