/**
 * Unit tests for webhook call mapper functions.
 */

import {
  webhookCallFromDto,
  webhookCallsFromDto,
  webhookCallCursorFromDto,
  webhookCallResultFromDto,
} from '../../../src/mappers/webhook-call';

describe('webhookCallFromDto', () => {
  it('should map all fields from DTO', () => {
    const dto = {
      id: 'wc-1',
      eventId: 'evt-1',
      webhookId: 'wh-1',
      payload: '{"type":"request.approved"}',
      status: 'delivered',
      statusMessage: 'OK',
      attempts: '1',
      updatedAt: new Date('2024-06-01'),
      createdAt: new Date('2024-01-01'),
    };

    const result = webhookCallFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.id).toBe('wc-1');
    expect(result!.eventId).toBe('evt-1');
    expect(result!.webhookId).toBe('wh-1');
    expect(result!.payload).toBe('{"type":"request.approved"}');
    expect(result!.status).toBe('delivered');
    expect(result!.statusMessage).toBe('OK');
    expect(result!.attempts).toBe('1');
    expect(result!.updatedAt).toBeInstanceOf(Date);
    expect(result!.createdAt).toBeInstanceOf(Date);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      id: 'wc-2',
      event_id: 'evt-2',
      webhook_id: 'wh-2',
      status_message: 'Failed',
      updated_at: '2024-07-01T00:00:00Z',
      created_at: '2024-02-01T00:00:00Z',
    };

    const result = webhookCallFromDto(dto);

    expect(result!.eventId).toBe('evt-2');
    expect(result!.webhookId).toBe('wh-2');
    expect(result!.statusMessage).toBe('Failed');
    expect(result!.updatedAt).toBeInstanceOf(Date);
    expect(result!.createdAt).toBeInstanceOf(Date);
  });

  it('should return undefined for null input', () => {
    expect(webhookCallFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(webhookCallFromDto(undefined)).toBeUndefined();
  });

  it('should handle empty object', () => {
    const result = webhookCallFromDto({});
    expect(result).toBeDefined();
    expect(result!.id).toBeUndefined();
  });
});

describe('webhookCallsFromDto', () => {
  it('should map array of DTOs', () => {
    const dtos = [
      { id: 'wc-1', status: 'delivered' },
      { id: 'wc-2', status: 'failed' },
    ];

    const result = webhookCallsFromDto(dtos);

    expect(result).toHaveLength(2);
    expect(result[0].status).toBe('delivered');
    expect(result[1].status).toBe('failed');
  });

  it('should return empty array for null input', () => {
    expect(webhookCallsFromDto(null)).toEqual([]);
  });

  it('should return empty array for undefined input', () => {
    expect(webhookCallsFromDto(undefined)).toEqual([]);
  });

  it('should return empty array for empty array', () => {
    expect(webhookCallsFromDto([])).toEqual([]);
  });
});

describe('webhookCallCursorFromDto', () => {
  it('should map cursor with next and previous pages', () => {
    const dto = {
      nextPage: 'page-2',
      previousPage: 'page-0',
    };

    const result = webhookCallCursorFromDto(dto);

    expect(result).toBeDefined();
    expect(result!.nextPage).toBe('page-2');
    expect(result!.previousPage).toBe('page-0');
    expect(result!.hasNextPage).toBe(true);
    expect(result!.hasPreviousPage).toBe(true);
  });

  it('should handle snake_case field names', () => {
    const dto = {
      next_page: 'p2',
      previous_page: 'p0',
    };

    const result = webhookCallCursorFromDto(dto);

    expect(result!.nextPage).toBe('p2');
    expect(result!.previousPage).toBe('p0');
    expect(result!.hasNextPage).toBe(true);
    expect(result!.hasPreviousPage).toBe(true);
  });

  it('should set hasNextPage to false when no nextPage', () => {
    const dto = { previousPage: 'page-0' };

    const result = webhookCallCursorFromDto(dto);

    expect(result!.hasNextPage).toBe(false);
    expect(result!.hasPreviousPage).toBe(true);
  });

  it('should set hasPreviousPage to false when no previousPage', () => {
    const dto = { nextPage: 'page-2' };

    const result = webhookCallCursorFromDto(dto);

    expect(result!.hasNextPage).toBe(true);
    expect(result!.hasPreviousPage).toBe(false);
  });

  it('should return undefined for null input', () => {
    expect(webhookCallCursorFromDto(null)).toBeUndefined();
  });

  it('should return undefined for undefined input', () => {
    expect(webhookCallCursorFromDto(undefined)).toBeUndefined();
  });
});

describe('webhookCallResultFromDto', () => {
  it('should map result with calls and cursor', () => {
    const dto = {
      calls: [
        { id: 'wc-1', status: 'delivered' },
        { id: 'wc-2', status: 'failed' },
      ],
      cursor: { nextPage: 'page-2' },
    };

    const result = webhookCallResultFromDto(dto);

    expect(result.calls).toHaveLength(2);
    expect(result.cursor).toBeDefined();
    expect(result.cursor!.hasNextPage).toBe(true);
  });

  it('should return default for null input', () => {
    const result = webhookCallResultFromDto(null);

    expect(result.calls).toEqual([]);
    expect(result.cursor).toBeUndefined();
  });

  it('should return default for undefined input', () => {
    const result = webhookCallResultFromDto(undefined);

    expect(result.calls).toEqual([]);
  });

  it('should handle missing cursor', () => {
    const dto = {
      calls: [{ id: 'wc-1' }],
    };

    const result = webhookCallResultFromDto(dto);

    expect(result.calls).toHaveLength(1);
    expect(result.cursor).toBeUndefined();
  });
});
