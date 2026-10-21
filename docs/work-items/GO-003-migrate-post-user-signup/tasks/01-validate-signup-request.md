# Task 01: Decode and validate the signup request

- **Item:** GO-003
- **Status:** OPEN
- **Covers:** AC5 (validation failures return 400 and `{"msg":"Invalid inputs"}`)

## Goal
A bad signup body gets `400` with `{"msg":"Invalid inputs"}`. A good body passes validation and stops there, with no database, hashing or token yet.

Each step below has an expectation (from `backend/routes/user.js:15-20`, the zod schema) and a check. Finish one step, pass its check, then move on.

## Steps

- [ ] **1. Create the request struct.**
  - Expect: the struct lives in the handler. Four `string` fields, with JSON keys `firstname`, `lastname`, `username`, `password` (as in `user.js:16-19`).
  - Check: `go build ./internal/user` exits 0.

- [ ] **2. Decode the JSON body.**
  - Expect: a body with the four keys decodes. A body with a non-string value (for example `"firstname": 1`) fails to decode.
  - Check: a test decodes one valid body and one with a number, and sees success and failure.

- [ ] **3. Check `firstname` and `lastname` are non-empty.**
  - Expect: a validation function takes plain `string` values, not the struct. An empty string fails (`min(1)`, `user.js:16-17`). Whitespace-only passes, because the JS code does not trim.
  - Check: test rows for `""` (fails), `"   "` (passes), and `"A"` (passes).

- [ ] **4. Check the password has at least 6 characters.**
  - Expect: `"abcde"` fails and `"abcdef"` passes (`min(6)`, `user.js:19`). Count characters, not bytes: Go's `len` counts bytes, so a multi-byte password can look longer than it is.
  - Check: test rows for `"abcde"`, `"abcdef"`, and a 6-character non-ASCII password (passes).

- [ ] **5. Check `username` is an email.**
  - Expect: a validation function takes the `username` string. It follows the practical email rule decided by the owner (option 2), which requires:
    - a local part of 1 to 64 characters, made of letters, digits and `. _ % + -`, with no leading, trailing or doubled dot;
    - exactly one `@`;
    - a domain of 1 to 253 characters, made of labels of letters, digits and hyphens separated by dots, with no label starting or ending with a hyphen;
    - at least one dot in the domain, and a last label (top-level domain) of at least 2 letters;
    - no spaces, and no more than 254 characters in total.
  - Check: test rows for one valid address and for each rule above that an address can break.

- [ ] **6. Register the route and return 400.**
  - Expect: `POST /api/v1/user/signup` is registered. Any failure above returns status `400`, `Content-Type: application/json`, and body `{"msg":"Invalid inputs"}`.
  - Check: `curl -i -X POST localhost:4500/api/v1/user/signup -H "Content-Type: application/json" -d '{"firstname":"A","lastname":"B","username":"not-an-email","password":"secret1"}'` returns `400` and `{"msg":"Invalid inputs"}`.

## Open decisions (ask before the step that needs them)
1. **Email rule (step 5): decided by owner (option 2, a practical email rule).** The rule is written out in step 5.
2. **Malformed JSON (step 2): decided by owner.** A body that is not valid JSON returns `400` with `{"msg":"invalid request body"}`. The JS differs: Express returns its default 400 response before the route runs (`backend/routes/user.js:22-24`). Record this as an intentional difference when the item is closed.

## Result
_Not yet checked._

## Redo checklist
_Filled after the check passes._
