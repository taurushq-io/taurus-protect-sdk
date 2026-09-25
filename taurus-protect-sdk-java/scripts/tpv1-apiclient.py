#!/usr/bin/env python3
"""Post-generation edits to the generated ApiClient.java.

Replaces the old `patch < openapi-tpv1.patch` step: that patch pinned context lines the generator no
longer emits ("Authorisation" became "Authorization"), and `patch` is not installed everywhere this
script runs. Every edit is anchored and fails loudly if the generator's output moved.
"""
import sys

AUTH_IMPORTS = (
    "import com.taurushq.sdk.protect.openapi.auth.Authentication;\n"
    "import com.taurushq.sdk.protect.openapi.auth.HttpBasicAuth;\n"
    "import com.taurushq.sdk.protect.openapi.auth.HttpBearerAuth;\n"
    "import com.taurushq.sdk.protect.openapi.auth.ApiKeyAuth;\n"
)
API_KEY_AUTH = 'authentications.put("ApiKeyTPV1", new ApiKeyAuth("header", "Authorization"));'
TPV1_HELPERS = '''    /**
     * Helper method to set API key value for the first API key authentication.
     *
     * @param apiKey API key
     */
    public void setApiKeyTPV1(String apiKey) {
        for (Authentication auth : authentications.values()) {
            if (auth instanceof ApiKeyTPV1Auth) {
                ((ApiKeyTPV1Auth) auth).setApiKey(apiKey);
                return;
            }
        }
        throw new RuntimeException("No API key authentication configured!");
    }

    /**
     * Helper method to set API secret value for the first API key authentication.
     *
     * @param apiSecret API key
     */
    public void setApiSecretTPV1(String apiSecret) throws ApiKeyTPV1Exception {
        for (Authentication auth : authentications.values()) {
            if (auth instanceof ApiKeyTPV1Auth) {
                ((ApiKeyTPV1Auth) auth).setApiSecret(apiSecret);
                return;
            }
        }
        throw new RuntimeException("No API secret authentication configured!");
    }

'''
BASE_PATH_DOC = "    /**\n     * Get base path\n"
COLLECTION_BRANCH = "        } else if (param instanceof Collection) {\n"
# The server decodes a bytes query parameter as base64 (the JSON form of the same field); without
# this branch String.valueOf(byte[]) sends "[B@1b6d3586" and every bytes cursor is unusable.
BYTES_BRANCH = (
    "        } else if (param instanceof byte[]) {\n"
    "            return java.util.Base64.getEncoder().encodeToString((byte[]) param);\n"
)


def replace_once(src, old, new, what):
    count = src.count(old)
    if count != 1:
        sys.exit(f"ERROR: expected exactly one {what} in ApiClient.java, found {count}")
    return src.replace(old, new)


def main(path):
    src = open(path, encoding="utf-8").read()
    src = replace_once(src, AUTH_IMPORTS, "import com.taurushq.sdk.protect.openapi.auth.*;\n", "auth import block")
    if src.count(API_KEY_AUTH) != 2:
        sys.exit(f"ERROR: expected two ApiKeyTPV1 registrations, found {src.count(API_KEY_AUTH)}")
    src = src.replace(API_KEY_AUTH, 'authentications.put("ApiKeyTPV1", new ApiKeyTPV1Auth());')
    src = replace_once(src, BASE_PATH_DOC, TPV1_HELPERS + BASE_PATH_DOC, "'Get base path' javadoc")
    src = replace_once(src, COLLECTION_BRANCH, BYTES_BRANCH + COLLECTION_BRANCH, "parameterToString Collection branch")
    open(path, "w", encoding="utf-8").write(src)
    print("post-generation: ApiClient.java TPV1 auth + byte[] query parameters")


if __name__ == "__main__":
    main(sys.argv[1])
