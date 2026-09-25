package com.taurushq.sdk.protect.client.testutil;

import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.auth.ApiKeyTPV1Auth;
import okhttp3.HttpUrl;
import okhttp3.Interceptor;
import okhttp3.MediaType;
import okhttp3.Protocol;
import okhttp3.Request;
import okhttp3.Response;
import okhttp3.ResponseBody;
import okio.Buffer;

import java.io.IOException;
import java.util.ArrayDeque;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.Deque;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * The transport stub for unit tests: an OkHttp interceptor that records every request and
 * answers from a queue of canned JSON bodies, without touching the network.
 * <p>
 * It is installed on a real {@link ApiClient} through {@code setHttpClient}, so the
 * GENERATED client runs end to end: its parameter names, its query encoding, its Gson
 * (de)serializers. Stubbing above the generated layer is how SDKs came to read reply fields
 * the generated models do not have.
 * <pre>{@code
 * StubTransport stub = StubTransport.replying("{\"result\":[],\"totalItems\":\"0\"}");
 * WalletService wallets = new WalletService(stub.client(), new ApiExceptionMapper());
 * wallets.getWallets(0, 0);
 * assertEquals(Arrays.asList(new String[]{"limit", "20"}), ...stub.only().query());
 * }</pre>
 * A call that reads several endpoints routes them by path instead, answering each request
 * from what it asks for:
 * <pre>{@code
 * StubTransport stub = StubTransport.replying()
 *         .route("/api/rest/v1/addresses", r -> page(r.values("addressIds")));
 * }</pre>
 */
public final class StubTransport implements Interceptor {

    /**
     * Answers one routed request.
     */
    @FunctionalInterface
    public interface Responder {

        /**
         * Returns the JSON body for a request.
         *
         * @param request the recorded request
         * @return the reply body
         */
        String reply(Recorded request);
    }

    /**
     * A placeholder host: the stub answers every request itself.
     */
    public static final String BASE_PATH = "https://stub.invalid";

    private static final MediaType JSON = MediaType.get("application/json; charset=utf-8");

    private final Deque<String> replies = new ArrayDeque<>();
    private final List<Recorded> requests = new ArrayList<>();
    private final Map<String, Responder> routes = new HashMap<>();
    private final Map<String, Integer> routeStatus = new HashMap<>();
    private final ApiClient client;
    private String lastReply = "{}";
    private int status = 200;

    private StubTransport(final String... bodies) {
        replies.addAll(Arrays.asList(bodies));
        client = new ApiClient();
        client.setBasePath(BASE_PATH);
        ((ApiKeyTPV1Auth) client.getAuthentication("ApiKeyTPV1")).setBearerToken("stub-token");
        client.setHttpClient(client.getHttpClient().newBuilder().addInterceptor(this).build());
    }

    /**
     * Creates a stub answering with the given bodies in order; the last one repeats.
     *
     * @param bodies the JSON reply bodies, none for {@code {}}
     * @return the stub
     */
    public static StubTransport replying(final String... bodies) {
        return new StubTransport(bodies);
    }

    /**
     * Makes every reply carry this HTTP status.
     *
     * @param code the status code
     * @return this stub
     */
    public StubTransport withStatus(final int code) {
        this.status = code;
        return this;
    }

    /**
     * Answers every request to a path from the request itself; unrouted paths keep the
     * queue of replies.
     *
     * @param path      the decoded request path
     * @param responder builds the reply body
     * @return this stub
     */
    public StubTransport route(final String path, final Responder responder) {
        return route(path, 200, responder);
    }

    /**
     * Answers every request to a path with an HTTP status and a body built from the request.
     *
     * @param path      the decoded request path
     * @param code      the status code
     * @param responder builds the reply body
     * @return this stub
     */
    public StubTransport route(final String path, final int code, final Responder responder) {
        routes.put(path, responder);
        routeStatus.put(path, code);
        return this;
    }

    /**
     * Returns the requests recorded to one path, in order.
     *
     * @param path the decoded request path
     * @return the requests
     */
    public List<Recorded> requestsTo(final String path) {
        List<Recorded> matching = new ArrayList<>();
        for (Recorded request : requests) {
            if (request.path().equals(path)) {
                matching.add(request);
            }
        }
        return matching;
    }

    /**
     * Returns the API client wired to this stub.
     *
     * @return the client
     */
    public ApiClient client() {
        return client;
    }

    /**
     * Returns every request recorded so far, in order.
     *
     * @return the requests
     */
    public List<Recorded> requests() {
        return Collections.unmodifiableList(requests);
    }

    /**
     * Returns the only request recorded, failing when there is not exactly one.
     *
     * @return the request
     */
    public Recorded only() {
        if (requests.size() != 1) {
            throw new AssertionError("expected exactly one request, got " + requests.size());
        }
        return requests.get(0);
    }

    @Override
    public Response intercept(final Chain chain) throws IOException {
        Request request = chain.request();
        Recorded recorded = new Recorded(request);
        requests.add(recorded);
        String body;
        int code;
        Responder responder = routes.get(recorded.path());
        if (responder != null) {
            body = responder.reply(recorded);
            code = routeStatus.get(recorded.path());
        } else {
            if (!replies.isEmpty()) {
                lastReply = replies.poll();
            }
            body = lastReply;
            code = status;
        }
        return new Response.Builder()
                .request(request)
                .protocol(Protocol.HTTP_1_1)
                .code(code)
                .message(code == 200 ? "OK" : "stub")
                .body(ResponseBody.create(body, JSON))
                .build();
    }

    /**
     * One recorded request.
     */
    public static final class Recorded {

        private final String method;
        private final HttpUrl url;
        private final String body;

        Recorded(final Request request) throws IOException {
            this.method = request.method();
            this.url = request.url();
            if (request.body() == null) {
                this.body = null;
            } else {
                Buffer buffer = new Buffer();
                request.body().writeTo(buffer);
                this.body = buffer.readUtf8();
            }
        }

        /**
         * Returns the HTTP method.
         *
         * @return the method
         */
        public String method() {
            return method;
        }

        /**
         * Returns the decoded path.
         *
         * @return the path
         */
        public String path() {
            return url.encodedPath();
        }

        /**
         * Returns the query exactly as sent, still percent-encoded.
         *
         * @return the raw query, or null when there is none
         */
        public String rawQuery() {
            return url.encodedQuery();
        }

        /**
         * Returns the decoded query pairs, sorted by name then value, so two requests compare
         * as multisets.
         *
         * @return the pairs as {name, value}
         */
        public List<List<String>> query() {
            List<List<String>> pairs = new ArrayList<>();
            for (int i = 0; i < url.querySize(); i++) {
                pairs.add(Arrays.asList(url.queryParameterName(i), url.queryParameterValue(i)));
            }
            pairs.sort((a, b) -> {
                int byName = a.get(0).compareTo(b.get(0));
                return byName != 0 ? byName : String.valueOf(a.get(1)).compareTo(String.valueOf(b.get(1)));
            });
            return pairs;
        }

        /**
         * Returns the value of the only occurrence of a query parameter.
         *
         * @param name the parameter name
         * @return the decoded value, or null when absent
         */
        public String param(final String name) {
            List<String> values = url.queryParameterValues(name);
            if (values.size() > 1) {
                throw new AssertionError(name + " was sent " + values.size() + " times");
            }
            return values.isEmpty() ? null : values.get(0);
        }

        /**
         * Returns every value of a repeated query parameter, in order.
         *
         * @param name the parameter name
         * @return the decoded values, empty when absent
         */
        public List<String> values(final String name) {
            return url.queryParameterValues(name);
        }

        /**
         * Returns the request body.
         *
         * @return the body text, or null for a request without one
         */
        public String body() {
            return body;
        }
    }
}
