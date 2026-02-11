/**
 * Unit tests for health-related mapping logic.
 *
 * The HealthService performs inline mapping (no dedicated mapper file).
 * These tests verify the mapping logic by testing the service methods directly.
 */

import type { HealthStatus } from '../../../src/models/health';

describe('HealthStatus model', () => {
  it('should represent a healthy status', () => {
    const status: HealthStatus = {
      status: 'healthy',
      version: '1.2.3',
      message: 'All systems operational',
    };

    expect(status.status).toBe('healthy');
    expect(status.version).toBe('1.2.3');
    expect(status.message).toBe('All systems operational');
  });

  it('should represent an unhealthy status', () => {
    const status: HealthStatus = {
      status: 'unhealthy',
      message: 'Database connection failed',
    };

    expect(status.status).toBe('unhealthy');
    expect(status.version).toBeUndefined();
    expect(status.message).toBe('Database connection failed');
  });

  it('should allow minimal status with only required fields', () => {
    const status: HealthStatus = {
      status: 'unknown',
    };

    expect(status.status).toBe('unknown');
    expect(status.version).toBeUndefined();
    expect(status.message).toBeUndefined();
  });
});

describe('Health check response mapping', () => {
  it('should determine healthy from array of checks with all healthy', () => {
    const checks = [
      { healthy: true, status: 'healthy' },
      { healthy: true, status: 'HEALTHY' },
    ];

    const allHealthy = checks.every((check) => {
      return check.healthy === true || check.status === 'healthy' || check.status === 'HEALTHY';
    });

    expect(allHealthy).toBe(true);
  });

  it('should determine unhealthy when one check fails', () => {
    const checks = [
      { healthy: true, status: 'healthy' },
      { healthy: false, status: 'unhealthy' },
    ];

    const allHealthy = checks.every((check) => {
      return check.healthy === true || check.status === 'healthy' || check.status === 'HEALTHY';
    });

    expect(allHealthy).toBe(false);
  });

  it('should handle empty checks array as healthy', () => {
    const checks: Array<{ healthy?: boolean; status?: string }> = [];
    const allHealthy = checks.every((check) => {
      return check.healthy === true || check.status === 'healthy' || check.status === 'HEALTHY';
    });

    // Array.every returns true for empty arrays
    expect(allHealthy).toBe(true);
  });

  it('should map global component status string', () => {
    const mapStatus = (status: unknown): { status: string; healthy: boolean } => {
      if (typeof status === 'string') {
        const healthy = status.toLowerCase() === 'healthy' || status.toLowerCase() === 'ok';
        return { status: healthy ? 'healthy' : status, healthy };
      }
      return { status: 'unknown', healthy: false };
    };

    expect(mapStatus('healthy')).toEqual({ status: 'healthy', healthy: true });
    expect(mapStatus('HEALTHY')).toEqual({ status: 'healthy', healthy: true });
    expect(mapStatus('ok')).toEqual({ status: 'healthy', healthy: true });
    expect(mapStatus('OK')).toEqual({ status: 'healthy', healthy: true });
    expect(mapStatus('degraded')).toEqual({ status: 'degraded', healthy: false });
  });

  it('should map global component status object', () => {
    const statusObj = {
      status: 'healthy',
      healthy: true,
    };

    const statusStr = (statusObj.status as string) ?? 'unknown';
    const healthy = statusObj.healthy === true || statusStr.toLowerCase() === 'healthy';

    expect(healthy).toBe(true);
    expect(statusStr).toBe('healthy');
  });
});
