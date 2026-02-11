package com.taurushq.sdk.protect.client.model;

/**
 * Represents the tenant configuration in the Taurus Protect system.
 * <p>
 * Contains various configuration settings for the tenant, including
 * security requirements, feature flags, and system parameters.
 *
 * @see ConfigService
 */
public class TenantConfig {

    private String tenantId;
    private String superAdminMinimumSignatures;
    private String baseCurrency;
    private Boolean isMFAMandatory;
    private Boolean excludeContainer;
    private Float feeLimitFactor;
    private String protectEngineVersion;
    private Boolean restrictSourcesForWhitelistedAddresses;
    private Boolean isProtectEngineCold;
    private Boolean isColdProtectEngineOffline;
    private Boolean isPhysicalAirGapEnabled;

    public String getTenantId() {
        return tenantId;
    }

    public void setTenantId(final String tenantId) {
        this.tenantId = tenantId;
    }

    public String getSuperAdminMinimumSignatures() {
        return superAdminMinimumSignatures;
    }

    public void setSuperAdminMinimumSignatures(final String superAdminMinimumSignatures) {
        this.superAdminMinimumSignatures = superAdminMinimumSignatures;
    }

    public String getBaseCurrency() {
        return baseCurrency;
    }

    public void setBaseCurrency(final String baseCurrency) {
        this.baseCurrency = baseCurrency;
    }

    public Boolean getMFAMandatory() {
        return isMFAMandatory;
    }

    public void setMFAMandatory(final Boolean mfaMandatory) {
        isMFAMandatory = mfaMandatory;
    }

    public Boolean getExcludeContainer() {
        return excludeContainer;
    }

    public void setExcludeContainer(final Boolean excludeContainer) {
        this.excludeContainer = excludeContainer;
    }

    public Float getFeeLimitFactor() {
        return feeLimitFactor;
    }

    public void setFeeLimitFactor(final Float feeLimitFactor) {
        this.feeLimitFactor = feeLimitFactor;
    }

    public String getProtectEngineVersion() {
        return protectEngineVersion;
    }

    public void setProtectEngineVersion(final String protectEngineVersion) {
        this.protectEngineVersion = protectEngineVersion;
    }

    public Boolean getRestrictSourcesForWhitelistedAddresses() {
        return restrictSourcesForWhitelistedAddresses;
    }

    public void setRestrictSourcesForWhitelistedAddresses(final Boolean restrictSourcesForWhitelistedAddresses) {
        this.restrictSourcesForWhitelistedAddresses = restrictSourcesForWhitelistedAddresses;
    }

    public Boolean getProtectEngineCold() {
        return isProtectEngineCold;
    }

    public void setProtectEngineCold(final Boolean protectEngineCold) {
        isProtectEngineCold = protectEngineCold;
    }

    public Boolean getColdProtectEngineOffline() {
        return isColdProtectEngineOffline;
    }

    public void setColdProtectEngineOffline(final Boolean coldProtectEngineOffline) {
        isColdProtectEngineOffline = coldProtectEngineOffline;
    }

    public Boolean getPhysicalAirGapEnabled() {
        return isPhysicalAirGapEnabled;
    }

    public void setPhysicalAirGapEnabled(final Boolean physicalAirGapEnabled) {
        isPhysicalAirGapEnabled = physicalAirGapEnabled;
    }
}
