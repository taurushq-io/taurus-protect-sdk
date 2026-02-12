/**
 * Action service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving automated action configurations.
 * Actions allow automated workflows to be triggered based on specific conditions
 * such as balance thresholds.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { ActionsApi } from '../internal/openapi/apis/ActionsApi';
import { actionEnvelopeFromDto, actionEnvelopesFromDto } from '../mappers/action';
import type { ActionEnvelope, ListActionsOptions } from '../models/action';
import { BaseService } from './base';

/**
 * Service for managing automated actions in the Taurus-PROTECT system.
 *
 * Actions allow automated workflows to be triggered based on specific conditions
 * such as balance thresholds. When conditions are met, tasks like transfers or
 * notifications can be executed automatically.
 *
 * @example
 * ```typescript
 * // List all actions
 * const actions = await actionService.list();
 * for (const action of actions) {
 *   console.log(`${action.id}: ${action.label} (${action.status})`);
 * }
 *
 * // Get a specific action
 * const action = await actionService.get('action-123');
 * console.log(`Action: ${action.label}, Auto-approve: ${action.autoApprove}`);
 *
 * // List with pagination
 * const paged = await actionService.list({ limit: '10', offset: '0' });
 * ```
 */
export class ActionService extends BaseService {
  private readonly actionsApi: ActionsApi;

  /**
   * Creates a new ActionService instance.
   *
   * @param actionsApi - The ActionsApi instance from the OpenAPI client
   */
  constructor(actionsApi: ActionsApi) {
    super();
    this.actionsApi = actionsApi;
  }

  /**
   * Lists all actions.
   *
   * @returns Array of action envelopes
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const actions = await actionService.list();
   * for (const action of actions) {
   *   console.log(`${action.id}: ${action.label}`);
   * }
   * ```
   */
  async list(): Promise<ActionEnvelope[]>;

  /**
   * Lists actions with optional filters.
   *
   * @param options - Optional filtering and pagination options
   * @returns Array of action envelopes
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List with pagination
   * const actions = await actionService.list({
   *   limit: '10',
   *   offset: '0',
   * });
   *
   * // List specific actions by ID
   * const specific = await actionService.list({
   *   ids: ['action-1', 'action-2'],
   * });
   * ```
   */
  async list(options: ListActionsOptions): Promise<ActionEnvelope[]>;

  async list(options?: ListActionsOptions): Promise<ActionEnvelope[]> {
    return this.execute(async () => {
      const response = await this.actionsApi.actionServiceGetActions({
        limit: options?.limit,
        offset: options?.offset,
        ids: options?.ids,
      });

      const result =
        (response as Record<string, unknown>).result ??
        (response as Record<string, unknown>).actions;
      return actionEnvelopesFromDto(result as unknown[]);
    });
  }

  /**
   * Retrieves a specific action by its ID.
   *
   * @param actionId - The ID of the action to retrieve
   * @returns The action envelope
   * @throws {@link ValidationError} If actionId is empty
   * @throws {@link NotFoundError} If action is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const action = await actionService.get('action-123');
   * console.log(`Action: ${action.label}`);
   * console.log(`Status: ${action.status}`);
   * console.log(`Auto-approve: ${action.autoApprove}`);
   * ```
   */
  async get(actionId: string): Promise<ActionEnvelope> {
    if (!actionId || actionId.trim() === '') {
      throw new ValidationError('actionId is required');
    }

    return this.execute(async () => {
      const response = await this.actionsApi.actionServiceGetAction({
        id: actionId,
      });

      const result =
        (response as Record<string, unknown>).action ??
        (response as Record<string, unknown>).result;
      const action = actionEnvelopeFromDto(result);

      if (!action) {
        throw new NotFoundError(`Action with id '${actionId}' not found`);
      }

      return action;
    });
  }
}
