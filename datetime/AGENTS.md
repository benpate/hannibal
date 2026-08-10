# AGENTS.md — hannibal/datetime

Notes for anyone (human or agent) changing this package. Read [README.md](README.md) first for what it does; this file is about the traps.

## Why this package exists

BUG-003211. ActivityPub work involves **two** timestamp formats, and they are easy to confuse:

| Where | Format | Who does it |
|---|---|---|
| AS2 properties (`published`, `updated`, …) | RFC3339 — `2026-01-02T15:04:05Z` | this package |
| HTTP `Date` header on a signed request | IMF-fixdate — `Fri, 02 Jan 2026 22:04:05 GMT` | `sigs`, unexported |

`published` was being emitted in the HTTP header format, which peers parsing AS2 as ISO 8601 cannot read. `TestFormat_NotAnHTTPDate` is the regression guard for exactly that — do not delete it because it "looks redundant" next to `TestFormat`.

The HTTP formatter used to be an exported `hannibal.TimeFormat` at the repo root, where nothing connected it to its use. It now lives as unexported `dateHeader`/`parseDateHeader` in `sigs`, beside the `FieldDate` constant and the verifier that has to read it back. **Do not add a general-purpose HTTP date formatter here.** The whole point of the split is that someone building an AS2 document never has the wrong function in scope.

## Rules that look like bugs but aren't

**An empty string is a successful return.** Not an error, not a placeholder. It means "this value cannot be expressed as a conformant AS2 date-time, so omit the property." Do not change these to return the Unix epoch, `time.Now()`, or an `error` — callers throughout hannibal and emissary branch on `!= ""`.

**Zero time → `""`.** An absent `published` is well-defined in AS2; a `published` of 1970-01-01 is a false claim about the object.

**Year outside `[1, 9999]` → `""`.** RFC3339 fixes the year at four digits. Go will happily render more: `math.MaxInt64` seconds becomes `"292277026596-12-04T15:30:07Z"`, which no conformant parser accepts. Sentinel timestamps reach this code in practice — that guard is load-bearing.

**UTC normalization is deliberate.** Both `Z` and a numeric offset are legal RFC3339, but pinning to UTC keeps output stable regardless of the host's local zone. Without it the tests are machine-dependent.

**`time.RFC3339`, not `time.RFC3339Nano`.** Fractional seconds are dropped on purpose, keeping output to the plain `date-time` production.

## Traps when editing the tests

**Year 1, January 1, 00:00:00 UTC *is* Go's zero time.** So the low boundary in `TestFormat_OutOfRangeYear` is tested one second later, at `00:00:01`. At midnight the function correctly returns `""` — but for being zero, not for being out of range, which is not what that test is asserting. If you "simplify" that to midnight the test still passes and stops testing the boundary.

**Assert literal strings, not round-trips through this package's own helpers.** `property/time_test.go` used to compare `value.String()` against the same formatter the type calls, so the two could never disagree — and a wrong format would have agreed with itself forever. Assert `"2024-01-22T15:04:05Z"`.

**`FromUnixMilli` vs `FromUnixSeconds` is the mistake these two exist to prevent.** Handing milliseconds to the seconds constructor lands in the year ~58000, which the range guard converts to `""` — so the failure mode is a *silently dropped property*, not visible garbage. `TestFromUnix` pins this.

## Verification

```
go test -cover ./datetime/
```

Coverage is **100% of statements**, and every branch is exercised. Keep it there — the package is four small functions whose entire value is in their edge cases, so an untested branch here is a untested edge case by definition.
