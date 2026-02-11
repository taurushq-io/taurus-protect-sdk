/**
 * Change service for Taurus-PROTECT SDK.
 *
 * Provides methods for managing configuration changes that require approval.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { ChangesApi } from '../internal/openapi/apis/ChangesApi';
import { changeFromDto, createChangeRequestToDto, listChangesResultFromDto } from '../mappers/change';
import type {
  Change,
  CreateChangeRequest,
  ListChangesOptions,
  ListChangesForApprovalOptions,
  ListChangesResult,
} from '../models/change';
import { BaseService } from './base';

/**
 * Service for managing configuration changes in the Taurus-PROTECT system.
 *
 * Changes represent modifications to system configuration that require approval
 * before taking effect. This includes user management, role assignments, and
 * other administrative operations that follow an approval workflow.
 *
 * @example
 * ```typescript
 * // Get changes pending approval
 * const result = await changeService.listForApproval({ pageSize: 50 });
 * for (const change of result.changes) {
 *   console.log(`${change.action} on ${change.entity} by ${change.creatorId}`);
 * }
 *
 * // Approve a change
 * await changeService.approve(changeId);
 *
 * // Reject a change
 * await changeService.reject(changeId);
 *
 * // Get changes with filters
 * const userChanges = await changeService.list({
 *   entity: 'user',
 *   status: 'Created',
 * });
 * ```
 */
export class ChangeService extends BaseService {
  private readonly changesApi: ChangesApi;

  /**
   * Creates a new ChangeService instance.
   *
   * @param changesApi - The ChangesApi instance from the OpenAPI client
   */
  constructor(changesApi: ChangesApi) {
    super();
    this.changesApi = changesApi;
  }

  /**
   * Gets a change by ID.
   *
   * @param id - The change ID
   * @returns The change
   * @throws {@link ValidationError} If id is empty
   * @throws {@link NotFoundError} If change is not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const change = await changeService.get('123');
   * console.log(`${change.action} on ${change.entity}`);
   * ```
   */
  async get(id: string): Promise<Change> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      const response = await this.changesApi.changeServiceGetChange({ id });
      const resp = response as Record<string, unknown>;
      const result = resp.result ?? resp.change;
      const change = changeFromDto(result);

      if (!change) {
        throw new NotFoundError(`Change with id '${id}' not found`);
      }

      return change;
    });
  }

  /**
   * Lists changes with optional filtering.
   *
   * @param options - Optional filtering options
   * @returns The list of changes with pagination info
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List all changes
   * const result = await changeService.list();
   *
   * // Filter by entity type and status
   * const userChanges = await changeService.list({
   *   entity: 'user',
   *   status: 'Created',
   * });
   *
   * // Paginate through results
   * const firstPage = await changeService.list({ pageSize: 50 });
   * if (firstPage.hasNext) {
   *   const nextPage = await changeService.list({
   *     pageSize: 50,
   *     currentPage: firstPage.currentPage,
   *     pageRequest: 'NEXT',
   *   });
   * }
   * ```
   */
  async list(options?: ListChangesOptions): Promise<ListChangesResult> {
    return this.execute(async () => {
      const response = await this.changesApi.changeServiceGetChanges({
        entity: options?.entity,
        status: options?.status,
        creatorId: options?.creatorId,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize?.toString(),
        entityIDs: options?.entityIDs,
        entityUUIDs: options?.entityUUIDs,
      });

      return listChangesResultFromDto(response);
    });
  }

  /**
   * Lists changes pending approval.
   *
   * @param options - Optional filtering options
   * @returns The list of changes with pagination info
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Get all changes pending approval
   * const result = await changeService.listForApproval();
   *
   * // Filter by entity types
   * const userAndGroupChanges = await changeService.listForApproval({
   *   entities: ['user', 'group'],
   * });
   * ```
   */
  async listForApproval(options?: ListChangesForApprovalOptions): Promise<ListChangesResult> {
    return this.execute(async () => {
      const response = await this.changesApi.changeServiceGetChangesForApproval({
        entities: options?.entities,
        sortOrder: options?.sortOrder,
        cursorCurrentPage: options?.currentPage,
        cursorPageRequest: options?.pageRequest,
        cursorPageSize: options?.pageSize?.toString(),
        entityIDs: options?.entityIDs,
        entityUUIDs: options?.entityUUIDs,
      });

      return listChangesResultFromDto(response);
    });
  }

  /**
   * Approves a change.
   *
   * @param id - The change ID to approve
   * @throws {@link ValidationError} If id is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await changeService.approve('123');
   * ```
   */
  async approve(id: string): Promise<void> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      await this.changesApi.changeServiceApproveChange({ id, body: {} });
    });
  }

  /**
   * Approves multiple changes.
   *
   * @param ids - The list of change IDs to approve
   * @throws {@link ValidationError} If ids is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await changeService.approveMany(['123', '456', '789']);
   * ```
   */
  async approveMany(ids: string[]): Promise<void> {
    if (!ids || ids.length === 0) {
      throw new ValidationError('ids cannot be empty');
    }

    return this.execute(async () => {
      await this.changesApi.changeServiceApproveChanges({
        body: { ids },
      });
    });
  }

  /**
   * Rejects a change.
   *
   * @param id - The change ID to reject
   * @throws {@link ValidationError} If id is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await changeService.reject('123');
   * ```
   */
  async reject(id: string): Promise<void> {
    if (!id || id.trim() === '') {
      throw new ValidationError('id is required');
    }

    return this.execute(async () => {
      await this.changesApi.changeServiceRejectChange({ id, body: {} });
    });
  }

  /**
   * Rejects multiple changes.
   *
   * @param ids - The list of change IDs to reject
   * @throws {@link ValidationError} If ids is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await changeService.rejectMany(['123', '456', '789']);
   * ```
   */
  async rejectMany(ids: string[]): Promise<void> {
    if (!ids || ids.length === 0) {
      throw new ValidationError('ids cannot be empty');
    }

    return this.execute(async () => {
      await this.changesApi.changeServiceRejectChanges({
        body: { ids },
      });
    });
  }

  /**
   * Creates a change request.
   *
   * @param request - The change request details
   * @returns The created change ID
   * @throws {@link ValidationError} If request is missing or has invalid fields
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const changeId = await changeService.create({
   *   action: 'update',
   *   entity: 'businessrule',
   *   entityId: '123',
   *   changes: { rulevalue: '100' },
   *   comment: 'Update rule value',
   * });
   * ```
   */
  async create(request: CreateChangeRequest): Promise<string> {
    if (!request) {
      throw new ValidationError('request is required');
    }
    if (!request.action || request.action.trim() === '') {
      throw new ValidationError('action is required');
    }
    if (!request.entity || request.entity.trim() === '') {
      throw new ValidationError('entity is required');
    }

    return this.execute(async () => {
      const dto = createChangeRequestToDto(request);
      const response = await this.changesApi.changeServiceCreateChange({ body: dto as any });
      const resp = response as Record<string, unknown>;
      const result = resp.result as Record<string, unknown> | undefined;
      return result?.id != null ? String(result.id) : '';
    });
  }
}
