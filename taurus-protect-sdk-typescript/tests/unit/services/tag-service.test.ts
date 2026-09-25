/**
 * Unit tests for TagService.
 */

import { TagService } from '../../../src/services/tag-service';
import { ValidationError, NotFoundError } from '../../../src/errors';
import { TagsApi } from '../../../src/internal/openapi/apis/TagsApi';
import { RecordingTransport } from '../pagination/harness';

function createMockApi(): jest.Mocked<TagsApi> {
  return {
    tagServiceGetTags: jest.fn(),
    tagServiceCreateTag: jest.fn(),
    tagServiceDeleteTag: jest.fn(),
  } as unknown as jest.Mocked<TagsApi>;
}

describe('TagService', () => {
  let mockApi: jest.Mocked<TagsApi>;
  let service: TagService;

  beforeEach(() => {
    mockApi = createMockApi();
    service = new TagService(mockApi);
  });

  describe('get', () => {
    it('should throw ValidationError when tagId is empty', async () => {
      await expect(service.get('')).rejects.toThrow(ValidationError);
      await expect(service.get('')).rejects.toThrow('tagId is required');
    });

    it('should throw NotFoundError when tag not found', async () => {
      mockApi.tagServiceGetTags.mockResolvedValue({
        result: [],
      } as never);

      await expect(service.get('nonexistent')).rejects.toThrow(NotFoundError);
    });

    it('should return tag when found', async () => {
      mockApi.tagServiceGetTags.mockResolvedValue({
        result: [{ id: 'tag-123', value: 'production' }],
      } as never);

      const tag = await service.get('tag-123');
      expect(tag).toBeDefined();
    });
  });

  describe('list', () => {
    it('should return tags', async () => {
      mockApi.tagServiceGetTags.mockResolvedValue({
        result: [
          { id: 'tag-1', value: 'production' },
          { id: 'tag-2', value: 'staging' },
        ],
      } as never);

      const tags = await service.list();
      expect(tags).toHaveLength(2);
    });

    it('should handle empty results', async () => {
      mockApi.tagServiceGetTags.mockResolvedValue({
        result: [],
      } as never);

      const tags = await service.list();
      expect(tags).toHaveLength(0);
    });

    // The endpoint is not paged: list used to cut the reply at 20 (50 before) with no
    // pagination, so a longer list read as complete.
    it('should return every tag the reply carries, sending no page parameters', async () => {
      const rows = Array.from({ length: 25 }, (_v, i) => ({ id: `t${i + 1}`, value: `tag ${i + 1}` }));
      const transport = new RecordingTransport({ result: rows });

      const tags = await new TagService(new TagsApi(transport.configuration)).list();

      expect(tags.map((t) => t.id)).toEqual(rows.map((r) => r.id));
      const request = transport.single();
      expect(request.path).toBe('/api/rest/v1/tags');
      expect(request.query).toEqual([]);
    });

    it('should send the ids and query filters', async () => {
      const transport = new RecordingTransport({ result: [] });

      await new TagService(new TagsApi(transport.configuration)).list({
        ids: ['t1', 't2'],
        query: 'prod',
      });

      expect(transport.single().query).toEqual([
        ['ids', 't1'],
        ['ids', 't2'],
        ['query', 'prod'],
      ]);
    });

    it.each([{ limit: 10 }, { offset: 20 }])(
      'should reject the removed paging option %p by name, sending nothing',
      async (paging) => {
        await expect(service.list(paging as never)).rejects.toThrow(
          /tags.list is not paged: it returns every tag, so limit and offset are not accepted/
        );
        expect(mockApi.tagServiceGetTags).not.toHaveBeenCalled();
      }
    );
  });

  describe('create', () => {
    it('should throw ValidationError when name is empty', async () => {
      await expect(service.create({ name: '', color: '#FF0000' })).rejects.toThrow(ValidationError);
      await expect(service.create({ name: '', color: '#FF0000' })).rejects.toThrow('name is required');
    });

    it('should throw ValidationError when color is empty', async () => {
      await expect(service.create({ name: 'new-tag', color: '' })).rejects.toThrow(ValidationError);
      await expect(service.create({ name: 'new-tag', color: '' })).rejects.toThrow('color is required');
    });

    it('should create tag with valid request', async () => {
      mockApi.tagServiceCreateTag.mockResolvedValue({
        result: { id: 'tag-new', value: 'new-tag' },
      } as never);

      const tag = await service.create({ name: 'new-tag', color: '#FF0000' });
      expect(tag).toBeDefined();
    });
  });

  describe('delete', () => {
    it('should throw ValidationError when tagId is empty', async () => {
      await expect(service.delete('')).rejects.toThrow(ValidationError);
    });

    it('should delete tag with valid ID', async () => {
      mockApi.tagServiceDeleteTag.mockResolvedValue(undefined as never);
      await service.delete('tag-123');
      expect(mockApi.tagServiceDeleteTag).toHaveBeenCalledWith({ id: 'tag-123' });
    });
  });
});
