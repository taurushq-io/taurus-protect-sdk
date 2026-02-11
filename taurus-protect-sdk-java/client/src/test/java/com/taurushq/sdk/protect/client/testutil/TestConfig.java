package com.taurushq.sdk.protect.client.testutil;

import java.io.FileInputStream;
import java.io.IOException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.Collections;
import java.util.List;
import java.util.Properties;

/**
 * Configuration for tests with multi-identity support.
 * <p>
 * Loads configuration from {@code test.properties} in the project root,
 * with environment variable overrides for CI/CD pipelines.
 * <p>
 * Each identity may have API credentials (for making API calls), a private key
 * (for signing operations), and/or a public key (for SuperAdmin verification).
 *
 * <p>Environment variables:</p>
 * <ul>
 *   <li>PROTECT_INTEGRATION_TEST - Set to "true" to enable tests</li>
 *   <li>PROTECT_API_HOST - API host URL</li>
 *   <li>PROTECT_API_KEY_N - API key for identity N (1-based)</li>
 *   <li>PROTECT_API_SECRET_N - API secret for identity N</li>
 *   <li>PROTECT_PRIVATE_KEY_N - Private key (PEM) for identity N</li>
 *   <li>PROTECT_PUBLIC_KEY_N - Public key (PEM) for identity N</li>
 * </ul>
 */
public final class TestConfig {

    private static final String ENV_INTEGRATION_TEST = "PROTECT_INTEGRATION_TEST";
    private static final String ENV_API_HOST = "PROTECT_API_HOST";

    private static final Properties properties;
    private static final List<Identity> identities;

    static {
        properties = loadProperties();
        identities = loadIdentities();
    }

    private TestConfig() {
        // Utility class
    }

    // ── Public API ──────────────────────────────────────────────────────

    /**
     * Returns the API host, preferring ENV var over properties file.
     */
    public static String getHost() {
        String env = System.getenv(ENV_API_HOST);
        if (env != null && !env.isEmpty()) {
            return env;
        }
        if (properties != null) {
            String prop = properties.getProperty("host", "");
            if (!prop.isEmpty()) {
                return prop;
            }
        }
        return "";
    }

    /**
     * Returns the identity at the given 1-based index.
     *
     * @param index 1-based identity index
     * @return the identity
     * @throws IndexOutOfBoundsException if index is out of range
     */
    public static Identity getIdentity(int index) {
        if (index < 1 || index > identities.size()) {
            throw new IndexOutOfBoundsException(
                    "Identity index " + index + " out of range [1.." + identities.size() + "]");
        }
        return identities.get(index - 1);
    }

    /**
     * Returns how many identities are configured.
     */
    public static int getIdentityCount() {
        return identities.size();
    }

    /**
     * Returns all configured identities (unmodifiable).
     */
    public static List<Identity> getIdentities() {
        return identities;
    }

    /**
     * Returns PEM-encoded SuperAdmin public keys from identities that have a public key.
     */
    public static List<String> getSuperAdminKeys() {
        List<String> keys = new ArrayList<>();
        for (Identity identity : identities) {
            if (identity.hasPublicKey()) {
                keys.add(identity.getPublicKey());
            }
        }
        return Collections.unmodifiableList(keys);
    }

    /**
     * Returns the minimum number of valid signatures required.
     */
    public static int getMinValidSignatures() {
        if (properties != null) {
            String prop = properties.getProperty("minValidSignatures", "2");
            try {
                return Integer.parseInt(prop);
            } catch (NumberFormatException e) {
                return 2;
            }
        }
        return 2;
    }

    /**
     * Returns true if tests should run.
     * Enable automatically if at least one identity has API credentials.
     */
    public static boolean isEnabled() {
        String env = System.getenv(ENV_INTEGRATION_TEST);
        if ("true".equalsIgnoreCase(env)) {
            return true;
        }
        // Enable automatically if at least one identity has API credentials
        for (Identity identity : identities) {
            if (identity.hasApiCredentials()) {
                return true;
            }
        }
        return false;
    }

    // ── Properties file loading ─────────────────────────────────────────

    private static Properties loadProperties() {
        // Try classpath first (src/test/resources/), then filesystem fallbacks
        // The classpath approach works regardless of CWD when run via Maven/Gradle
        try {
            java.io.InputStream is = TestConfig.class.getClassLoader()
                    .getResourceAsStream("test.properties");
            if (is != null) {
                Properties props = new Properties();
                try {
                    props.load(is);
                    return props;
                } finally {
                    is.close();
                }
            }
        } catch (IOException e) {
            System.err.println("Warning: Failed to load from classpath: " + e.getMessage());
        }

        // Filesystem fallbacks for running outside Maven
        String[] searchPaths = {
                "client/src/test/resources/test.properties",
                "src/test/resources/test.properties",
                "taurus-protect-sdk-java/client/src/test/resources/test.properties"
        };

        for (String searchPath : searchPaths) {
            Path path = Paths.get(searchPath).toAbsolutePath().normalize();
            if (path.toFile().exists()) {
                Properties props = new Properties();
                try (FileInputStream fis = new FileInputStream(path.toFile())) {
                    props.load(fis);
                    return props;
                } catch (IOException e) {
                    System.err.println("Warning: Failed to load " + path + ": " + e.getMessage());
                }
            }
        }
        return null;
    }

    /**
     * Resolves a setting value: env var takes priority, then properties file, then empty string.
     *
     * @param envVar        environment variable name (null to skip env lookup)
     * @param propertiesKey properties file key
     * @return the resolved value, or empty string if not found
     */
    private static String resolve(String envVar, String propertiesKey) {
        if (envVar != null) {
            String env = System.getenv(envVar);
            if (env != null && !env.isEmpty()) {
                return env;
            }
        }
        if (properties != null) {
            String prop = properties.getProperty(propertiesKey, "");
            if (!prop.isEmpty()) {
                return prop;
            }
        }
        return "";
    }

    /**
     * Loads identities by scanning for identity.1, identity.2, ... until none found.
     */
    private static List<Identity> loadIdentities() {
        List<Identity> list = new ArrayList<>();
        for (int i = 1; ; i++) {
            String name = resolve(null, "identity." + i + ".name");
            String apiKey = resolve("PROTECT_API_KEY_" + i, "identity." + i + ".apiKey");
            String apiSecret = resolve("PROTECT_API_SECRET_" + i, "identity." + i + ".apiSecret");
            String privateKey = resolve("PROTECT_PRIVATE_KEY_" + i, "identity." + i + ".privateKey");
            String publicKey = resolve("PROTECT_PUBLIC_KEY_" + i, "identity." + i + ".publicKey");

            // Stop when nothing at all is defined for this index
            if (name.isEmpty() && apiKey.isEmpty() && apiSecret.isEmpty()
                    && privateKey.isEmpty() && publicKey.isEmpty()) {
                break;
            }

            if (name.isEmpty()) {
                name = "identity-" + i;
            }

            // Unescape \n in PEM keys loaded from properties files
            if (!privateKey.isEmpty()) {
                privateKey = privateKey.replace("\\n", "\n");
            }
            if (!publicKey.isEmpty()) {
                publicKey = publicKey.replace("\\n", "\n");
            }

            list.add(new Identity(i, name, apiKey, apiSecret, privateKey, publicKey));
        }
        return Collections.unmodifiableList(list);
    }

    // ── Identity inner class ────────────────────────────────────────────

    /**
     * Represents a user identity in the system.
     * <p>
     * An identity may have API credentials (for making API calls via ProtectClient),
     * a private key (for signing operations like request approval),
     * and/or a public key (for verification as a SuperAdmin).
     * <p>
     * Not all fields are required — an identity with only a publicKey is valid
     * (e.g., a SuperAdmin whose public key we need for verification but whose
     * API credentials we don't use in this test).
     */
    public static final class Identity {
        private final int index;
        private final String name;
        private final String apiKey;
        private final String apiSecret;
        private final String privateKey;
        private final String publicKey;

        Identity(int index, String name, String apiKey, String apiSecret,
                 String privateKey, String publicKey) {
            this.index = index;
            this.name = name;
            this.apiKey = apiKey;
            this.apiSecret = apiSecret;
            this.privateKey = privateKey;
            this.publicKey = publicKey;
        }

        /** Returns the 1-based index of this identity. */
        public int getIndex() {
            return index;
        }

        /** Returns the human-readable name. */
        public String getName() {
            return name;
        }

        /** Returns the API key, or empty string if not configured. */
        public String getApiKey() {
            return apiKey;
        }

        /** Returns the API secret, or empty string if not configured. */
        public String getApiSecret() {
            return apiSecret;
        }

        /** Returns the PEM-encoded EC private key, or empty string if not configured. */
        public String getPrivateKey() {
            return privateKey;
        }

        /** Returns the PEM-encoded EC public key, or empty string if not configured. */
        public String getPublicKey() {
            return publicKey;
        }

        /** Returns true if this identity has a private key for signing. */
        public boolean hasPrivateKey() {
            return privateKey != null && !privateKey.isEmpty();
        }

        /** Returns true if this identity has a public key for verification. */
        public boolean hasPublicKey() {
            return publicKey != null && !publicKey.isEmpty();
        }

        /** Returns true if this identity has API credentials for making calls. */
        public boolean hasApiCredentials() {
            return apiKey != null && !apiKey.isEmpty()
                    && apiSecret != null && !apiSecret.isEmpty();
        }

        @Override
        public String toString() {
            return name + " (identity " + index + ")";
        }
    }
}
