# Task 03: Load .env with godotenv

> Retro-written on 2026-10-21 from the GO-001 closure evidence. The redo checklist is a draft for you to confirm or edit.

- **Item:** GO-001
- **Status:** PASSED
- **Covers:** AC12, AC14, AC15, AC16

## Goal
Values in `go-backend/.env` are loaded into the process environment before configuration is read. A missing file is fine, a malformed one stops startup, and the shell wins over the file.

## Concepts
What `godotenv.Load` does (it sets environment variables and does not override existing ones), `errors.Is` with `fs.ErrNotExist`, the working directory, and adding a third-party module.

## Check
`.env` with `PORT=4600` and the shell unset listens on 4600; shell `PORT=4700` wins and listens on 4700; with no `.env` and exported variables the service starts; a malformed `.env` exits 1.

## Result
2026-10-17: listening on 4600, then 4700 with the shell value, then 4600 with no `.env` and exported variables. A malformed line exited 1 with the godotenv error. `go.mod` lists `github.com/joho/godotenv v1.5.1` as a direct requirement.

## Redo checklist
1. Run `go get github.com/joho/godotenv` — this adds the module to `go.mod` and `go.sum`.
2. Call `godotenv.Load()` first in `main`, before `config.Load()` — configuration must see the values.
3. Ignore the error only when `errors.Is(err, fs.ErrNotExist)`; fatal on any other error — a missing file is normal, a malformed one is a mistake.
4. Run the service from `go-backend/` — `Load()` reads `.env` from the working directory.
5. Add `.env` to `.gitignore` and commit a `.env.example` with placeholders — secrets stay out of git.
