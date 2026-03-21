# Migration Order Fix TODO

## Steps:
- [ ] Step 1: Read relevant migration files (00002, 00003, 00004 up/down)
- [ ] Step 2: Confirm dependency issue (refresh_tokens REFERENCES users before users created)
- [ ] Step 3: Implement OPTION A: Move users table creation to 00003_refresh_tokens.up.sql BEFORE refresh_tokens
- [ ] Step 4: Update 00004_create_users.up.sql to only indexes/triggers or skip if IF NOT EXISTS
- [ ] Step 5: Update down.sql files if needed
- [ ] Step 6: Verify with search_files for REFERENCES users
- [ ] Step 7: Test with docker compose

Current status: Files read, planning OPTION A
