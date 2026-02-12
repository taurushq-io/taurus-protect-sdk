package com.taurushq.sdk.protect.client.model.taurusnetwork;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;

import java.util.List;

/**
 * Result containing a list of pledges with pagination.
 */
public class PledgeResult {

    private List<Pledge> pledges;
    private ApiResponseCursor cursor;

    public List<Pledge> getPledges() {
        return pledges;
    }

    public void setPledges(final List<Pledge> pledges) {
        this.pledges = pledges;
    }

    public ApiResponseCursor getCursor() {
        return cursor;
    }

    public void setCursor(final ApiResponseCursor cursor) {
        this.cursor = cursor;
    }
}
