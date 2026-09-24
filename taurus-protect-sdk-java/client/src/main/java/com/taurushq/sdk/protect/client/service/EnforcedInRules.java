package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.model.Group;
import com.taurushq.sdk.protect.client.model.GroupUser;
import com.taurushq.sdk.protect.client.model.User;
import com.taurushq.sdk.protect.client.model.UserGroup;

import java.util.List;

/**
 * Fills in the {@code enforcedInRules} flags of a reply from an endpoint that computes them.
 * <p>
 * validatord omits a false bool from its JSON, so there an absent flag means false; on any other
 * endpoint it means "not computed" and stays null. {@code publicKeyEnforcedInRules} is a wrapper
 * on the wire, so its absence always means "not computed": it is never filled in.
 */
final class EnforcedInRules {

    private EnforcedInRules() {
    }

    /**
     * Fills in the flags of a user returned by an endpoint that computes them.
     *
     * @param user the mapped user, may be null
     * @return the same user
     */
    static User computed(final User user) {
        if (user == null) {
            return null;
        }
        user.setEnforcedInRules(orFalse(user.getEnforcedInRules()));
        if (user.getGroups() != null) {
            for (UserGroup group : user.getGroups()) {
                if (group != null) {
                    group.setEnforcedInRules(orFalse(group.getEnforcedInRules()));
                }
            }
        }
        return user;
    }

    /**
     * Fills in the flags of users returned by an endpoint that computes them.
     *
     * @param users the mapped users
     * @return the same list
     */
    static List<User> computedUsers(final List<User> users) {
        for (User user : users) {
            computed(user);
        }
        return users;
    }

    /**
     * Fills in the flags of groups, and of their users, returned by an endpoint that computes them.
     *
     * @param groups the mapped groups
     * @return the same list
     */
    static List<Group> computedGroups(final List<Group> groups) {
        for (Group group : groups) {
            if (group == null) {
                continue;
            }
            group.setEnforcedInRules(orFalse(group.getEnforcedInRules()));
            if (group.getUsers() != null) {
                for (GroupUser user : group.getUsers()) {
                    if (user != null) {
                        user.setEnforcedInRules(orFalse(user.getEnforcedInRules()));
                    }
                }
            }
        }
        return groups;
    }

    private static Boolean orFalse(final Boolean flag) {
        return flag == null ? Boolean.FALSE : flag;
    }
}
