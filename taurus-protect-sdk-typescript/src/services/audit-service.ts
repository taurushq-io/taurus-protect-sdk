/**
 * Audit service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving audit trail information.
 */

import type { AuditApi } from '../internal/openapi/apis/AuditApi';
import { auditTrailsFromDto } from '../mappers/audit';
import type { ListAuditTrailsOptions, ListAuditTrailsResult } from '../models/audit';
import { buildCursorPage, cursorRequest } from '../models/pagination';
import { BaseService } from './base';
import { cursorQuery } from './paging';

/**
 * Service for audit trail operations.
 *
 * Provides methods to list audit trails for tracking actions
 * performed in Taurus-PROTECT.
 *
 * @example
 * ```typescript
 * // List audit trails, page by page
 * let cursor: string | undefined;
 * do {
 *   const page = await auditService.list({ pageSize: 50, cursor });
 *   for (const audit of page.items) {
 *     console.log(`${audit.action} on ${audit.entity} by ${audit.userEmail}`);
 *   }
 *   cursor = page.pagination.hasMore ? page.pagination.nextCursor : undefined;
 * } while (cursor);
 *
 * // Filter by entity type
 * const walletAudits = await auditService.list({
 *   entities: ['WALLET'],
 *   actions: ['CREATE', 'UPDATE'],
 * });
 *
 * // Filter by date range
 * const recentAudits = await auditService.list({
 *   creationDateFrom: new Date('2024-01-01'),
 *   creationDateTo: new Date(),
 * });
 * ```
 */
export class AuditService extends BaseService {
  private readonly auditApi: AuditApi;

  /**
   * Creates a new AuditService instance.
   *
   * @param auditApi - The AuditApi instance from the OpenAPI client
   */
  constructor(auditApi: AuditApi) {
    super();
    this.auditApi = auditApi;
  }

  /**
   * Lists a page of audit trails with optional filtering.
   *
   * @param options - Filters, `pageSize` (1-100, default 20) and `cursor`
   * @returns The page of audit trails and its cursor pagination
   * @throws {@link ValidationError} If the page size is out of bounds
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List recent audit trails
   * const { items, pagination } = await auditService.list({ pageSize: 100 });
   *
   * // Filter by entity and action
   * const walletCreations = await auditService.list({
   *   entities: ['WALLET'],
   *   actions: ['CREATE'],
   * });
   *
   * // Filter by date range
   * const lastWeek = await auditService.list({
   *   creationDateFrom: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000),
   *   creationDateTo: new Date(),
   * });
   *
   * // Filter by user
   * const userActions = await auditService.list({
   *   externalUserId: 'user-123',
   * });
   * ```
   */
  async list(options?: ListAuditTrailsOptions): Promise<ListAuditTrailsResult> {
    const page = cursorRequest(options);

    return this.execute(async () => {
      const response = await this.auditApi.auditServiceGetAuditTrails({
        externalUserId: options?.externalUserId,
        entities: options?.entities,
        actions: options?.actions,
        creationDateFrom: options?.creationDateFrom,
        creationDateTo: options?.creationDateTo,
        ...cursorQuery(page),
        sortingSortOrder: options?.sortOrder,
      });

      return {
        items: auditTrailsFromDto(response.result),
        pagination: buildCursorPage(page.pageSize, response.cursor),
      };
    });
  }

  /**
   * Export audit trails to a formatted string (CSV or JSON).
   *
   * Note that only a maximum of 10000 trails are exportable at any one time.
   *
   * @param options - Optional filtering options
   * @returns The exported data as a string
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // Export all audit trails as CSV (default)
   * const csvData = await auditService.exportAuditTrails();
   *
   * // Export with filters
   * const filtered = await auditService.exportAuditTrails({
   *   externalUserId: 'user-123',
   *   creationDateFrom: new Date('2024-01-01'),
   *   format: 'json',
   * });
   *
   * // Export specific entity and action types
   * const walletCreations = await auditService.exportAuditTrails({
   *   entities: ['WALLET'],
   *   actions: ['CREATE'],
   * });
   * ```
   */
  async exportAuditTrails(options?: {
    externalUserId?: string;
    entities?: string[];
    actions?: string[];
    creationDateFrom?: Date;
    creationDateTo?: Date;
    format?: string;
  }): Promise<string> {
    return this.execute(async () => {
      const response = await this.auditApi.auditServiceExportAuditTrails({
        externalUserId: options?.externalUserId,
        entities: options?.entities,
        actions: options?.actions,
        creationDateFrom: options?.creationDateFrom,
        creationDateTo: options?.creationDateTo,
        format: options?.format,
      });
      return response.result ?? '';
    });
  }
}
