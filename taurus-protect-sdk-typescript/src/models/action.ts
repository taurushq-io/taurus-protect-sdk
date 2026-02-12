/**
 * Action models for Taurus-PROTECT SDK.
 *
 * Actions allow automated workflows to be triggered based on specific conditions
 * such as balance thresholds. When conditions are met, tasks like transfers or
 * notifications can be executed automatically.
 */

/**
 * Represents an action envelope containing an automated action configuration
 * with its metadata and execution history.
 */
export interface ActionEnvelope {
  /** Unique action identifier */
  readonly id?: string;
  /** Tenant identifier */
  readonly tenantId?: string;
  /** Action label/name */
  readonly label?: string;
  /** Action configuration with trigger and tasks */
  readonly action?: Action;
  /** Current action status */
  readonly status?: string;
  /** When the action was created */
  readonly creationDate?: Date;
  /** When the action was last updated */
  readonly updateDate?: Date;
  /** When the action was last checked */
  readonly lastCheckedDate?: Date;
  /** Whether the action auto-approves requests */
  readonly autoApprove?: boolean;
  /** Custom attributes associated with the action */
  readonly attributes?: ActionAttribute[];
  /** Execution history trails */
  readonly trails?: ActionTrail[];
}

/**
 * Action configuration containing trigger conditions and tasks to execute.
 */
export interface Action {
  /** Trigger configuration that determines when the action fires */
  readonly trigger?: ActionTrigger;
  /** Tasks to execute when the trigger conditions are met */
  readonly tasks?: ActionTask[];
}

/**
 * Trigger configuration for an action.
 */
export interface ActionTrigger {
  /** Type of trigger (e.g., "balance") */
  readonly kind?: string;
  /** Balance-based trigger configuration */
  readonly balance?: TriggerBalance;
}

/**
 * Balance-based trigger configuration.
 */
export interface TriggerBalance {
  /** Target address or wallet to monitor */
  readonly target?: ActionTarget;
  /** Comparison operator (e.g., "less_than", "greater_than") */
  readonly comparator?: ActionComparator;
  /** Amount threshold */
  readonly amount?: ActionAmount;
}

/**
 * Target for balance monitoring.
 */
export interface ActionTarget {
  /** Type of target (e.g., "address", "wallet") */
  readonly kind?: string;
  /** Address target configuration */
  readonly address?: TargetAddress;
  /** Wallet target configuration */
  readonly wallet?: TargetWallet;
}

/**
 * Address target for balance monitoring.
 */
export interface TargetAddress {
  /** Type identifier */
  readonly kind?: string;
  /** Address ID to monitor */
  readonly addressId?: string;
}

/**
 * Wallet target for balance monitoring.
 */
export interface TargetWallet {
  /** Type identifier */
  readonly kind?: string;
  /** Wallet ID to monitor */
  readonly walletId?: string;
}

/**
 * Comparator for balance conditions.
 */
export interface ActionComparator {
  /** Comparison kind (e.g., "less_than", "greater_than") */
  readonly kind?: string;
}

/**
 * Amount configuration for triggers and transfers.
 */
export interface ActionAmount {
  /** Type of amount (e.g., "fixed", "percentage") */
  readonly kind?: string;
  /** Crypto amount value */
  readonly cryptoAmount?: string;
}

/**
 * Task to execute when action triggers.
 */
export interface ActionTask {
  /** Type of task (e.g., "transfer", "notification") */
  readonly kind?: string;
  /** Transfer task configuration */
  readonly transfer?: TaskTransfer;
  /** Notification task configuration */
  readonly notification?: TaskNotification;
}

/**
 * Transfer task configuration.
 */
export interface TaskTransfer {
  /** Source of the transfer */
  readonly from?: ActionSource;
  /** Destination of the transfer */
  readonly to?: ActionDestination;
  /** Amount to transfer */
  readonly amount?: ActionAmount;
  /** Whether this is a top-up transfer */
  readonly topUp?: boolean;
  /** Whether to use all available funds */
  readonly useAllFunds?: boolean;
}

/**
 * Source for a transfer task.
 */
export interface ActionSource {
  /** Type of source (e.g., "address", "wallet") */
  readonly kind?: string;
  /** Source address ID */
  readonly addressId?: string;
  /** Source wallet ID */
  readonly walletId?: string;
}

/**
 * Destination for a transfer task.
 */
export interface ActionDestination {
  /** Type of destination (e.g., "address", "wallet", "whitelisted_address") */
  readonly kind?: string;
  /** Destination address ID */
  readonly addressId?: string;
  /** Destination whitelisted address ID */
  readonly whitelistedAddressId?: string;
  /** Destination wallet ID */
  readonly walletId?: string;
}

/**
 * Notification task configuration.
 */
export interface TaskNotification {
  /** Email addresses to notify */
  readonly emailAddresses?: string[];
  /** Notification message content */
  readonly notificationMessage?: string;
  /** Number of reminder notifications to send */
  readonly numberOfReminders?: string;
}

/**
 * Custom attribute associated with an action.
 */
export interface ActionAttribute {
  /** Attribute identifier */
  readonly id?: string;
  /** Tenant identifier */
  readonly tenantId?: string;
  /** Attribute key */
  readonly key?: string;
  /** Attribute value */
  readonly value?: string;
  /** Content type of the value */
  readonly contentType?: string;
}

/**
 * Execution trail record for an action.
 */
export interface ActionTrail {
  /** Trail record identifier */
  readonly id?: string;
  /** Action taken */
  readonly action?: string;
  /** Comment or description */
  readonly comment?: string;
  /** When the action occurred */
  readonly date?: Date;
  /** Status after the action */
  readonly actionStatus?: string;
}

/**
 * Options for listing actions.
 */
export interface ListActionsOptions {
  /** Maximum number of actions to return */
  limit?: string;
  /** Offset for pagination */
  offset?: string;
  /** Filter by specific action IDs */
  ids?: string[];
}
