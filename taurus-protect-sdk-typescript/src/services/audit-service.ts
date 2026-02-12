/**
 * Audit service for Taurus-PROTECT SDK.
 *
 * Provides methods for retrieving audit trail information.
 */

import { ValidationError } from '../errors';
import type { AuditApi } from '../internal/openapi/apis/AuditApi';
import { auditTrailsFromDto } from '../mappers/audit';
import type { AuditTrail, ListAuditTrailsOptions } from '../models/audit';
import { BaseService } from './base';

/**
 * Service for audit trail operations.
 *
 * Provides methods to list audit trails for tracking actions
 * performed in Taurus-PROTECT.
 *
 * @example
 * ```typescript
 * // List audit trails
 * const audits = await auditService.list({ limit: 50 });
 * for (const audit of audits) {
 *   console.log(`${audit.action} on ${audit.entity} by ${audit.userEmail}`);
 * }
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
   * Lists audit trails with optional filtering.
   *
   * @param options - Optional filtering options
   * @returns Array of audit trails
   * @throws {@link ValidationError} If limit is invalid
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List recent audit trails
   * const audits = await auditService.list({ limit: 100 });
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
  async list(options?: ListAuditTrailsOptions): Promise<AuditTrail[]> {
    const limit = options?.limit ?? 50;

    if (limit <= 0) {
      throw new ValidationError('limit must be positive');
    }

    return this.execute(async () => {
      const response = await this.auditApi.auditServiceGetAuditTrails({
        externalUserId: options?.externalUserId,
        entities: options?.entities,
        actions: options?.actions,
        creationDateFrom: options?.creationDateFrom,
        creationDateTo: options?.creationDateTo,
        cursorPageSize: String(limit),
        sortingSortOrder: options?.sortOrder,
      });

      const resp = response as Record<string, unknown>;
      const result = resp.auditTrails ?? resp.audit_trails ?? resp.result ?? resp.trails;
      return auditTrailsFromDto(result as unknown[]);
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
