/**
 * Unit tests for ReservationService.
 */

import { ReservationService } from '../../../src/services/reservation-service';
import { NotFoundError, ValidationError } from '../../../src/errors';
import type { ReservationsApi } from '../../../src/internal/openapi/apis/ReservationsApi';

function createMockApi(): jest.Mocked<ReservationsApi> {
  return {
    walletServiceGetReservations: jest.fn(),
    walletServiceGetReservation: jest.fn(),
    walletServiceGetReservationUTXO: jest.fn(),
  } as unknown as jest.Mocked<ReservationsApi>;
}

describe('ReservationService', () => {
  let mockApi: jest.Mocked<ReservationsApi>;
  let service: ReservationService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new ReservationService(mockApi);
  });

  describe('list', () => {
    it('should return reservations', async () => {
      mockApi.walletServiceGetReservations.mockResolvedValue({
        result: [
          { id: 'res-1', amount: '0.5', kind: 'OUTGOING' },
          { id: 'res-2', amount: '1.0', kind: 'CONSOLIDATION' },
        ],
      } as never);

      const reservations = await service.list();
      expect(reservations.items.map((r) => r.id)).toEqual(['res-1', 'res-2']);
      expect(reservations.pagination).toEqual({ pageSize: 20, nextCursor: '', hasMore: false });
    });

    it('should return the reply cursor and send it back as NEXT', async () => {
      mockApi.walletServiceGetReservations.mockResolvedValue({
        result: [{ id: 'res-1' }],
        cursor: { currentPage: 'next-page', hasNext: true },
      } as never);

      const first = await service.list({ pageSize: 1 });
      expect(first.pagination).toEqual({ pageSize: 1, nextCursor: 'next-page', hasMore: true });

      await service.list({ pageSize: 1, cursor: first.pagination.nextCursor });
      expect(mockApi.walletServiceGetReservations).toHaveBeenLastCalledWith(
        expect.objectContaining({
          cursorCurrentPage: 'next-page',
          cursorPageRequest: 'NEXT',
          cursorPageSize: '1',
        })
      );
    });

    it('should pass filter options to API', async () => {
      mockApi.walletServiceGetReservations.mockResolvedValue({
        result: [],
      } as never);

      await service.list({
        addressId: 'addr-123',
        kinds: ['OUTGOING', 'CONSOLIDATION'],
      });

      expect(mockApi.walletServiceGetReservations).toHaveBeenCalledWith(
        expect.objectContaining({
          addressId: 'addr-123',
          kinds: ['OUTGOING', 'CONSOLIDATION'],
        })
      );
    });

    it('should handle empty results', async () => {
      mockApi.walletServiceGetReservations.mockResolvedValue({
        result: [],
      } as never);

      const reservations = await service.list();
      expect(reservations.items).toHaveLength(0);
    });

    it('should work without options', async () => {
      mockApi.walletServiceGetReservations.mockResolvedValue({
        result: [],
      } as never);

      const reservations = await service.list();
      expect(reservations).toBeDefined();
    });
  });

  describe('get', () => {
    it('should throw ValidationError when id is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('id is required');
    });

    it('should throw ValidationError when id is whitespace', async () => {
      await expect(service.get('  ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when reservation is not found', async () => {
      mockApi.walletServiceGetReservation.mockResolvedValue({
        result: undefined,
      } as never);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });

    it('should return reservation for valid id', async () => {
      mockApi.walletServiceGetReservation.mockResolvedValue({
        result: {
          id: 'res-123',
          amount: '1.5',
          kind: 'OUTGOING',
          address: 'bc1qtest',
        },
      } as never);

      const reservation = await service.get('res-123');
      expect(reservation).toBeDefined();
      expect(mockApi.walletServiceGetReservation).toHaveBeenCalledWith({ id: 'res-123' });
    });
  });

  describe('getUtxo', () => {
    it('should throw ValidationError when id is empty', async () => {
      await expect(service.getUtxo('')).rejects.toThrow(ValidationError);
      await expect(service.getUtxo('')).rejects.toThrow('id is required');
    });

    it('should throw ValidationError when id is whitespace', async () => {
      await expect(service.getUtxo('  ')).rejects.toThrow(ValidationError);
    });

    it('should throw NotFoundError when UTXO is not found', async () => {
      mockApi.walletServiceGetReservationUTXO.mockResolvedValue({
        result: undefined,
      } as never);

      await expect(service.getUtxo('nonexistent')).rejects.toThrow(NotFoundError);
    });

    it('should return UTXO for valid id', async () => {
      mockApi.walletServiceGetReservationUTXO.mockResolvedValue({
        result: {
          hash: 'txhash123',
          value: '50000',
          blockHeight: '800000',
        },
      } as never);

      const utxo = await service.getUtxo('res-123');
      expect(utxo).toBeDefined();
      expect(mockApi.walletServiceGetReservationUTXO).toHaveBeenCalledWith({ id: 'res-123' });
    });
  });
});
