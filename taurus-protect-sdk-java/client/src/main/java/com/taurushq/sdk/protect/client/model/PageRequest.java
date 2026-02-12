package com.taurushq.sdk.protect.client.model;

/**
 * Page request type for cursor-based pagination.
 */
public enum PageRequest {

    /**
     * Request the first page.
     */
    FIRST,

    /**
     * Request the previous page.
     */
    PREVIOUS,

    /**
     * Request the next page.
     */
    NEXT,

    /**
     * Request the last page.
     */
    LAST
}
