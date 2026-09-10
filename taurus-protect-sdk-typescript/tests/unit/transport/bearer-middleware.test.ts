import type { RequestContext } from "../../../src/internal/openapi/runtime";
import { createBearerMiddleware } from "../../../src/transport/bearer-middleware";

describe("createBearerMiddleware", () => {
  function requestCtx(): RequestContext {
    return {
      fetch: (() => undefined) as unknown,
      url: "https://api.example.com/v1/wallets",
      init: { method: "GET", headers: { "Content-Type": "application/json" } },
    } as unknown as RequestContext;
  }

  it("adds an Authorization: Bearer header, preserving existing headers", async () => {
    const mw = createBearerMiddleware(() => "session-token-xyz");
    const result = await mw.pre!(requestCtx());
    expect(result).toBeDefined();
    const headers = (result as { init: { headers: Record<string, string> } })
      .init.headers;
    expect(headers.Authorization).toBe("Bearer session-token-xyz");
    expect(headers["Content-Type"]).toBe("application/json");
  });

  it("resolves the token per request from the provider", async () => {
    const tokens = ["tok-1", "tok-2"];
    let i = 0;
    const mw = createBearerMiddleware(() => tokens[i++]);
    for (const want of tokens) {
      const result = await mw.pre!(requestCtx());
      const headers = (result as { init: { headers: Record<string, string> } })
        .init.headers;
      expect(headers.Authorization).toBe(`Bearer ${want}`);
    }
  });
});
