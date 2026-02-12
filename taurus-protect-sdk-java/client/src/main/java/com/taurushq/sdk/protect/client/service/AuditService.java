package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.mapper.AuditMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.AuditTrailResult;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.api.AuditApi;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordExportAuditTrailsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetAuditTrailsReply;

import java.time.OffsetDateTime;
import java.util.List;

import static com.google.common.base.Preconditions.checkNotNull;

/**
 * Service for retrieving audit trail information in the Taurus Protect system.
 * <p>
 * This service provides access to the audit log, which records all significant
 * actions performed in the system for compliance and security purposes.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Get recent audit trails
 * AuditTrailResult result = client.getAuditService().getAuditTrails(
 *     null,                                    // externalUserId
 *     Arrays.asList("request"),                // entities
 *     Arrays.asList("approve"),                // actions
 *     OffsetDateTime.now().minusDays(7),       // from
 *     OffsetDateTime.now(),                    // to
 *     null                                     // cursor
 * );
 *
 * for (AuditTrail trail : result.getAuditTrails()) {
 *     System.out.println(trail.getEntity() + ": " + trail.getAction());
 * }
 * }</pre>
 *
 * @see AuditTrail
 * @see AuditTrailResult
 */
public class AuditService {

    private final AuditApi auditApi;
    private final ApiExceptionMapper apiExceptionMapper;
    private final AuditMapper auditMapper;

    /**
     * Instantiates a new Audit service.
     *
     * @param openApiClient      the OpenAPI client
     * @param apiExceptionMapper the API exception mapper
     */
    public AuditService(final ApiClient openApiClient, final ApiExceptionMapper apiExceptionMapper) {
        checkNotNull(openApiClient, "openApiClient cannot be null");
        checkNotNull(apiExceptionMapper, "apiExceptionMapper cannot be null");

        this.apiExceptionMapper = apiExceptionMapper;
        this.auditApi = new AuditApi(openApiClient);
        this.auditMapper = AuditMapper.INSTANCE;
    }

    /**
     * Retrieves audit trails with optional filtering.
     *
     * @param externalUserId filter by external user ID (optional)
     * @param entities       filter by entity types (optional)
     * @param actions        filter by action types (optional)
     * @param from           filter from date (optional)
     * @param to             filter to date (optional)
     * @param cursor         pagination cursor (optional)
     * @return the audit trail result with pagination
     * @throws ApiException if the API call fails
     */
    public AuditTrailResult getAuditTrails(final String externalUserId,
                                            final List<String> entities,
                                            final List<String> actions,
                                            final OffsetDateTime from,
                                            final OffsetDateTime to,
                                            final ApiRequestCursor cursor) throws ApiException {

        String cursorCurrentPage = null;
        String cursorPageRequest = null;
        String cursorPageSize = null;

        if (cursor != null) {
            cursorCurrentPage = cursor.getCurrentPage();
            cursorPageRequest = cursor.getPageRequest() != null ? cursor.getPageRequest().name() : null;
            cursorPageSize = String.valueOf(cursor.getPageSize());
        }

        try {
            TgvalidatordGetAuditTrailsReply reply = auditApi.auditServiceGetAuditTrails(
                    externalUserId,
                    entities,
                    actions,
                    from,
                    to,
                    cursorCurrentPage,
                    cursorPageRequest,
                    cursorPageSize,
                    null,  // sortingSortBy
                    null   // sortingSortOrder
            );
            return auditMapper.fromReply(reply);
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }

    /**
     * Exports audit trails to a specified format.
     *
     * @param externalUserId filter by external user ID (optional)
     * @param entities       filter by entity types (optional)
     * @param actions        filter by action types (optional)
     * @param from           filter from date (optional)
     * @param to             filter to date (optional)
     * @param format         the export format (e.g., "csv", "json")
     * @return the exported data as a string
     * @throws ApiException if the API call fails
     */
    public String exportAuditTrails(final String externalUserId,
                                     final List<String> entities,
                                     final List<String> actions,
                                     final OffsetDateTime from,
                                     final OffsetDateTime to,
                                     final String format) throws ApiException {
        try {
            TgvalidatordExportAuditTrailsReply reply = auditApi.auditServiceExportAuditTrails(
                    externalUserId,
                    entities,
                    actions,
                    from,
                    to,
                    format
            );
            return reply.getResult();
        } catch (com.taurushq.sdk.protect.openapi.ApiException e) {
            throw apiExceptionMapper.toApiException(e);
        }
    }
}
