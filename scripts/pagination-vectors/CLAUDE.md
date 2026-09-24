# scripts/pagination-vectors — the pagination gate's generators

Two generators, two outputs, loaded by all four SDK unit suites (loader paths: repo-root `CLAUDE.md`
§ "Pagination (cross-SDK)"). Never edit the JSON by hand.

| Generator | Writes | Holds |
|---|---|---|
| `generate.py` | `scripts/resources/pagination-vectors.json` | the reference rules: `parse_count`, `offset_page` (5 rules), `cursor_page`, `token_page`, `resolve_size`/`resolve_offset`; the operation → rule map; the offset / cursor / page-size / offset-input cases |
| `generate_requests.py` | `scripts/resources/list-request-vectors.json` | per operation: method, path, paging shape (read from the shared swagger) and the exact expected query multiset or body per canonical option set |

```bash
python3 scripts/pagination-vectors/generate.py && python3 scripts/pagination-vectors/generate_requests.py
```

Run both, in that order: `generate_requests.py` reads the operation map from `pagination-vectors.json`.
Then update the `counts` constants in all four loaders; each asserts them.

## Adding or changing an operation

1. Put the operationId in exactly one list in `generate.py`: `OFFSET_RULE_OPS` (with its rule),
   `CURSOR_OPS`, `TOKEN_OPS`, or `LIMIT_ONLY_OPS` (size kind, has-total). Add it to `TOTAL_OPS` when a
   cursor reply carries a count. The rule must be verified in tg-validatord **down to the store query or
   paginator** — a controller copying a field proves nothing (the transaction export's `offset` is
   copied and ignored, which is why it is `LIMIT_ONLY_OPS`).
2. In `generate_requests.py`, add what the call cannot omit to `REQUIRED` (options + where they land),
   what every SDK always sends to `ALWAYS`, and path values to `PATH_SENTINELS`. `paging_shape` derives
   the style from the swagger; a `limit_only` rule overrides the spec.
3. Decisions beyond paging (a filter's wire mapping, a rejection) go in `FILTER_VECTORS`.
4. Regenerate both files; update the four loaders' counts and adapter tables. An SDK that does not wrap
   the operation lists it, with a reason, in its loader's not-wrapped table.

## Traps

- **Tokens must be canonical standard base64.** Java's generated type for the token-only cursors
  (`GetWalletTokens`, `GetRulesHistory`) is `byte[]`: the SDK decodes the caller's string and
  `ApiClient.parameterToString` re-encodes it, so `ab+/cd==` (non-zero padding bits) comes back as
  `ab+/cQ==`. Real server tokens are canonical; test tokens must be too, and still carry `+ / =`.
- **Options use canonical snake_case names** (`limit`, `offset`, `page_size`, `cursor`,
  `exclude_disabled`, `from_currency_id`, `to_currency_ids`, `statuses`, `external_request_ids`, …); each loader maps them to its
  own fields, so no SDK's naming wins.
- **Defaults vectors expect ONLY the paging parameters** (plus `REQUIRED`/`ALWAYS`): a method that sends
  an extra parameter unasked fails them. That is deliberate — e.g. an export must not hard-code a format.
