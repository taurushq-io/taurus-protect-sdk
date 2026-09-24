/**
 * Tag service for Taurus-PROTECT SDK.
 *
 * Provides methods for tag management operations.
 */

import { NotFoundError, ValidationError } from '../errors';
import type { TagsApi } from '../internal/openapi/apis/TagsApi';
import { tagFromDto, tagsFromDto } from '../mappers/user';
import type { CreateTagRequest, ListTagsOptions, Tag } from '../models/user';
import { BaseService } from './base';

/**
 * Service for tag management operations.
 *
 * Provides methods to list, get, create, and delete tags.
 * Tags can be applied to wallets, addresses, and other entities
 * for organization and filtering.
 *
 * @example
 * ```typescript
 * // List every tag
 * const tags = await tagService.list();
 * for (const tag of tags) {
 *   console.log(`${tag.name}: ${tag.color}`);
 * }
 *
 * // Get single tag
 * const tag = await tagService.get('tag-123');
 * console.log(`Tag: ${tag.name}`);
 *
 * // Create tag
 * const newTag = await tagService.create({ name: 'Important', color: '#FF0000' });
 * console.log(`Created tag: ${newTag.id}`);
 *
 * // Delete tag
 * await tagService.delete('tag-123');
 * ```
 */
export class TagService extends BaseService {
  private readonly tagsApi: TagsApi;

  /**
   * Creates a new TagService instance.
   *
   * @param tagsApi - The TagsApi instance from the OpenAPI client
   */
  constructor(tagsApi: TagsApi) {
    super();
    this.tagsApi = tagsApi;
  }

  /**
   * Gets a tag by ID.
   *
   * @param tagId - The tag ID to retrieve
   * @returns The tag
   * @throws {@link ValidationError} If tagId is invalid
   * @throws {@link NotFoundError} If tag not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const tag = await tagService.get('tag-123');
   * console.log(`Tag: ${tag.name}, Color: ${tag.color}`);
   * ```
   */
  async get(tagId: string): Promise<Tag> {
    if (!tagId || tagId.trim() === '') {
      throw new ValidationError('tagId is required');
    }

    return this.execute(async () => {
      // The API doesn't have a direct get-by-id endpoint,
      // so we use the list endpoint with id filter
      const response = await this.tagsApi.tagServiceGetTags({
        ids: [tagId],
      });

      const result = response.result;

      if (!result || result.length === 0) {
        throw new NotFoundError(`Tag ${tagId} not found`);
      }

      const tag = tagFromDto(result[0]);
      if (!tag) {
        throw new NotFoundError(`Tag ${tagId} not found`);
      }

      return tag;
    });
  }

  /**
   * Lists every tag, optionally filtered by id or query.
   *
   * The tags endpoint is not paged, so the result is every tag the reply carries; there
   * is no `limit` or `offset`, and passing either is rejected by name.
   *
   * @param options - Optional `ids` and `query` filters
   * @returns Every matching tag
   * @throws {@link ValidationError} If `limit` or `offset` is passed
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * // List every tag
   * const tags = await tagService.list();
   *
   * // Search tags
   * const matching = await tagService.list({ query: 'important' });
   * ```
   */
  async list(options?: ListTagsOptions): Promise<Tag[]> {
    // A caller still paging with limit/offset would get every tag on each call and
    // could loop forever; refuse instead of ignoring them.
    const paging = options as { limit?: unknown; offset?: unknown } | undefined;
    if (paging?.limit !== undefined || paging?.offset !== undefined) {
      throw new ValidationError(
        "tags.list is not paged: it returns every tag, so limit and offset are not accepted"
      );
    }

    return this.execute(async () => {
      const response = await this.tagsApi.tagServiceGetTags({
        ids: options?.ids,
        query: options?.query,
      });
      return tagsFromDto(response.result);
    });
  }

  /**
   * Creates a new tag.
   *
   * @param request - Tag creation request
   * @returns The created tag
   * @throws {@link ValidationError} If name or color is empty
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * const tag = await tagService.create({
   *   name: 'High Priority',
   *   color: '#FF0000',
   * });
   * console.log(`Created tag: ${tag.id}`);
   * ```
   */
  async create(request: CreateTagRequest): Promise<Tag> {
    if (!request.name || request.name.trim() === '') {
      throw new ValidationError('name is required');
    }
    if (!request.color || request.color.trim() === '') {
      throw new ValidationError('color is required');
    }

    return this.execute(async () => {
      const response = await this.tagsApi.tagServiceCreateTag({
        body: {
          value: request.name,
          color: request.color,
        },
      });

      const tag = tagFromDto(response.result);

      if (!tag) {
        throw new ValidationError('Failed to create tag: no result returned');
      }

      return tag;
    });
  }

  /**
   * Deletes a tag.
   *
   * This removes the tag and all its assignments from entities.
   *
   * @param tagId - The tag ID to delete
   * @throws {@link ValidationError} If tagId is empty
   * @throws {@link NotFoundError} If tag not found
   * @throws {@link APIError} If API request fails
   *
   * @example
   * ```typescript
   * await tagService.delete('tag-123');
   * console.log('Tag deleted');
   * ```
   */
  async delete(tagId: string): Promise<void> {
    if (!tagId || tagId.trim() === '') {
      throw new ValidationError('tagId is required');
    }

    return this.execute(async () => {
      await this.tagsApi.tagServiceDeleteTag({ id: tagId });
    });
  }
}
