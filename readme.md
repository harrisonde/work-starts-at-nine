# WSAN — Work Starts At Nine

WSAN is a fictitious REST API in the spirit of FOAAS that exists for one reason: to tell people, in many tones and flavors, that work starts at nine. It's built on the Cidekar [adele](https://github.com/cidekar/adele) framework in Go 1.24 and serves a small catalogue of JSON-returning endpoints under `/api`.

## Run

```sh
make build && make run
```

The API is served under `/api`.

## Endpoints

| Method & Path | Description |
| --- | --- |
| `GET /api/` | version & info |
| `GET /api/operations` | list all operations |
| `GET /api/nine/{name}/{from}` | basic reminder |
| `GET /api/late/{name}/{from}` | you're late |
| `GET /api/reminder/{name}/{from}` | friendly reminder |
| `GET /api/strict/{name}/{from}` | strict tone |
| `GET /api/boss/{name}/{from}` | from the boss |
| `GET /api/hr/{name}/{from}` | HR formal |
| `GET /api/polite/{name}/{from}` | extremely polite |
| `GET /api/rude/{name}/{from}` | rude tone |
| `GET /api/monday/{name}/{from}` | Monday edition |
| `GET /api/meeting/{name}/{from}` | 9 AM meeting joke |
| `GET /api/earlybird/{name}/{from}` | early bird twist |
| `GET /api/coffee/{name}/{from}` | finish coffee on way in |
| `GET /api/standup/{name}/{from}` | standup edition |
| `GET /api/deadline/{name}/{from}` | deadline was 9 |
| `GET /api/random/{name}/{from}` | random tone per request |

Each operation is now backed by a procedural generator (`handlers/generator`),
so every request returns a different message assembled from per-tone fragment
pools. For reproducible output (handy for tests, screenshots, and bug reports)
append `?seed=<non-zero int64>`:

```sh
curl http://localhost:8080/api/nine/Dave/Alice?seed=42
```

The same seed against the same path will always return the same message.

**Seed semantics:**

- `?seed=<non-zero int64>` — deterministic, reproducible output.
- `seed=0` is **reserved** (the internal "time-seeded" sentinel) and is rejected.
- Any invalid seed (non-integer, out of range, or zero) returns **HTTP 400**
  with the WSAN envelope `{"message":"Invalid seed. Use a non-zero int64.","subtitle":"— WSAN"}`.
- Omit `?seed=` entirely to get fresh random output per request.

## Example

```sh
curl http://localhost:8080/api/nine/Dave/Alice
```


```json
{
  "message": "Hey Dave, work starts at 9.",
  "subtitle": "— Alice"
}
```

Example output; actual messages vary per request. Use `?seed=<non-zero int64>`
for reproducible output.

## Development

- `make fmt` — format all Go sources
- `make lint` — run golangci-lint
- `make test` — run tests with race detector and coverage
- `make ci` — run fmt-check, vet, lint, and tests
