package com.taurushq.sdk.protect.client.model;

/**
 * Represents a comparator for action triggers.
 * <p>
 * Comparators define the condition operator (e.g., less than, greater than).
 */
public class ActionComparator {

    private String kind;

    public String getKind() {
        return kind;
    }

    public void setKind(final String kind) {
        this.kind = kind;
    }
}
