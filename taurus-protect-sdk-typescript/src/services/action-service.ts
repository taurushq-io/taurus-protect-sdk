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
import {
  buildOffsetPagination,
  offsetRequest,
  type PaginatedResult,
} from '../models/pagination';
import { BaseService } from './base';
import { offsetQuery } from './paging';

/**
 * Service for managing automated actions in the Taurus-PROTECT system.
 *
 * Actions allow automated workflows to be triggered based on specific conditions
 * such as balance thresholds. When conditions are met, tasks like transfers or
 * notifications can be executed automatically.
 *
 * @example
 * ```typescript
 * // First page of actions
 * const { items, pagination } = await actionService.list();
 * for (const action of items) {
 *   console.log(`${action.id}: ${action.label} (${action.status})`);
 * }
 *
 * // Get a specific action
 * const action = await actionService.get('action-123');
 * console.log(`Action: ${action.label}, Auto-approve: ${action.autoApprove}`);
 *
 * // Next page
 * if (pagination.hasMore) {
 *   await actionService.list({ offset: pagination.nextOffset });
 * }
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
   * Lists a page of actions.
   *
   * @param options - Filters, `limit` (1-100, default 20) and `offset`
   * @returns The page of actions and its pagination, including the server's total
   * @throws {@link ValidationError} If limit or offset are out of bounds
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const page = await actionService.list({ limit: 10 });
   * for (const action of page.items) {
   *   console.log(`Action ${action.id}: ${action.label} (${action.status})`);
   * }
   * if (page.pagination.hasMore) {
   *   await actionService.list({ limit: 10, offset: page.pagination.nextOffset });
   * }
   * ```
   */
  async list(options?: ListActionsOptions): Promise<PaginatedResult<ActionEnvelope>> {
    const page = offsetRequest(options);

    return this.execute(async () => {
      const response = await this.actionsApi.actionServiceGetActions({
        ...offsetQuery(page),
        ids: options?.ids,
      });

      const rows = response.result ?? [];
      return {
        items: actionEnvelopesFromDto(rows),
        pagination: buildOffsetPagination('plus_rows', page, response, rows.length),
      };
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

      const action = actionEnvelopeFromDto(response.action);

      if (!action) {
        throw new NotFoundError(`Action with id '${actionId}' not found`);
      }

      return action;
    });
  }
}
