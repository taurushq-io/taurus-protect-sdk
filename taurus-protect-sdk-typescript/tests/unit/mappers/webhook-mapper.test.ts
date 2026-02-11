/**
 * Unit tests for webhook mapper functions.
 */

import { webhookFromDto, webhooksFromDto } from '../../../src/mappers/webhook';
import { WebhookStatus } from '../../../src/models/webhook';

describe('webhookFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'wh-1',
      url: 'https://example.com/webhook',
      events: ['TRANSACTION', 'REQUEST'],
      status: 'ACTIVE',
      secret: 'secret123',
      createdAt: new Date('2024-01-01'),
      updatedAt: new Date('2024-06-01'),
      failureCount: 0,
      lastFailureMessage: null,
    };

    const result = webhookFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('wh-1');
    expect(result!.url).toBe('https://example.com/webhook');
    expect(result!.events).toEqual(['TRANSACTION', 'REQUEST']);
    expect(result!.status).toBe(WebhookStatus.ACTIVE);
  });

  it('should handle status mapping', () => {
    expect(webhookFromDto({ id: '1', status: 'ACTIVE' })!.status).toBe(WebhookStatus.ACTIVE);
    expect(webhookFromDto({ id: '1', status: 'ENABLED' })!.status).toBe(WebhookStatus.ACTIVE);
    expect(webhookFromDto({ id: '1', status: 'INACTIVE' })!.status).toBe(WebhookStatus.INACTIVE);
    expect(webhookFromDto({ id: '1', status: 'DISABLED' })!.status).toBe(WebhookStatus.INACTIVE);
    expect(webhookFromDto({ id: '1', status: 'FAILED' })!.status).toBe(WebhookStatus.FAILED);
    expect(webhookFromDto({ id: '1', status: 'UNKNOWN' })!.status).toBeUndefined();
  });

  it('should handle comma-separated events string', () => {
    const dto = {
      id: 'wh-2',
      type: 'TRANSACTION,REQUEST',
    };
    const result = webhookFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.events).toEqual(['TRANSACTION', 'REQUEST']);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      webhook_id: 'wh-3',
      created_at: new Date('2024-01-01'),
      failure_count: 3,
      last_failure_message: 'timeout',
    };
    const result = webhookFromDto(dto);
    expect(result).toBeDefined();
    expect(result!.id).toBe('wh-3');
    expect(result!.failureCount).toBe(3);
    expect(result!.lastFailureMessage).toBe('timeout');
  });

  it('should return undefined for null input', () => {
    expect(webhookFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(webhookFromDto(undefined)).toBeUndefined();
  });
});

describe('webhooksFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: '1', url: 'https://example.com/1' },
      { id: '2', url: 'https://example.com/2' },
    ];
    const result = webhooksFromDto(dtos);
    expect(result).toHaveLength(2);
  });

  it('should return empty array for null input', () => {
    expect(webhooksFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(webhooksFromDto(undefined)).toEqual([]);
  });
});
