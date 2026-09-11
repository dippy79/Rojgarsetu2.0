---
name: rojgarsetu-production-audit
description: "Use when auditing or repairing RojgarSetu production workflows end to end, especially authentication, navigation, jobs, courses, videos, government forms, API integrations, crawlers, broken links, empty states, or pixel-accurate UI. Requires local evidence, real data, executable verification, and explicit approval before advancing phases."
argument-hint: "Audit or repair a RojgarSetu phase, starting with authentication"
---

# RojgarSetu Production Audit

Perform a staged, evidence-driven full-stack audit of RojgarSetu 2.0 against enterprise US/UK expectations. Work with real local services and persisted data. Do not invent records, stub successful responses, or claim a fix without running a relevant check.

## Operating Rules

- Start with the phase named by the user. If no phase is named, start with Phase 1.
- Before editing, inspect the owning frontend route/component, backend handler/service, persistence path, and nearest relevant test or executable check.
- State a concise local hypothesis about the failure and one check that could disconfirm it.
- Prefer the smallest root-cause change that preserves existing APIs and conventions.
- Allow mocks only inside isolated unit tests for testing an individual function offline. Strictly prohibit mocks everywhere else, including E2E tests, integration tests, API Gateway responses, database layers, and frontend components. Never use fake data or mocked integrations to make a workflow appear healthy; test against the configured local backend, database, crawler, or external service behavior.
- Protect secrets. Use existing environment configuration and never print tokens, passwords, API keys, or connection strings.
- Do not perform destructive database, migration, deployment, or git operations without explicit confirmation.
- After every substantive edit, immediately run the narrowest relevant validation before broadening the investigation.
- Keep unrelated user changes intact.
- Report exact findings with workspace-relative file links, observed responses or counts, commands run, and remaining risks.

## Phase Gates

Treat each phase as a separately approved work unit:

1. Inspect and reproduce the phase's failures.
2. Present concrete findings and the proposed minimal repair.
3. Implement the repair.
4. Run focused executable verification, then relevant integration checks.
5. Summarize results and stop.
6. Continue to the next phase only after the user explicitly replies with approval, such as `Approved`.

A phase is complete only when its UI path, API contract, persistence or upstream data path, error state, and relevant links have been checked. If an upstream dependency is unavailable, distinguish unavailable verification from a verified fix and stop at that boundary.

## Phase 1: Authentication

Audit `/login` and `/register` end to end.

### Frontend

- Locate the actual login and registration routes, forms, validation schema, API client, loading state, error rendering, and success navigation.
- Verify required fields, email format, password strength, confirmation matching, accessible labels, keyboard use, disabled/loading states, and clear actionable errors.
- Confirm validation behavior matches the backend contract rather than silently hiding server errors.

### API and Backend

- Identify the configured auth base URL and confirm the frontend calls the intended `/api/v1/auth/...` endpoints.
- Start only the required local services if they are not running, using the repository's documented commands.
- Test registration, login, invalid credentials, duplicate identity, malformed input, and protected access using real local requests. Use `curl` or the repository's existing API test tooling.
- Trace failures into the Go handlers, service layer, database queries, migrations, configuration, and middleware. Fix the controlling layer, not just the displayed error.
- Confirm status codes and JSON error/success shapes are consistent with the frontend client.

### Session State

- Verify the token/session mechanism, expiry behavior, refresh or logout behavior, storage security, and authenticated route protection.
- Prefer httpOnly, Secure, appropriately SameSite cookies where the existing architecture supports them. If local storage is required by the architecture, document the risk and verify the narrowest safe implementation.
- Confirm a successful login updates the navbar to authenticated actions such as Dashboard and Logout, and that logout clears the session and returns the user to the correct public state.

### Phase 1 Completion Checks

- Valid registration creates one real persisted account and returns the documented response.
- Valid login succeeds and establishes a session that authorizes a protected endpoint or page.
- Invalid, duplicate, and malformed requests fail with intentional status codes and visible messages.
- Refreshing the application preserves the intended session behavior.
- Logout invalidates the session and updates global navigation.
- Focused backend tests, frontend checks, and a manual or automated browser path pass.

Unit-test mocks may be used only when the test remains offline and targets one function in isolation. Auth integration, E2E, gateway, database, and frontend behavior checks must use the real configured components and persisted or externally returned data.

Stop after reporting Phase 1. Do not start Phase 2 until the user explicitly approves.

## Phase 2: Navbar and Global Routing

- Locate the actual navbar and router implementation, including responsive and authenticated variants.
- Verify every visible link resolves to an existing route and preserves intended query parameters.
- Test click and keyboard behavior for menus and dropdowns; verify open, close, outside-click, Escape, focus, and mobile behavior.
- Confirm active states derive from the current pathname and relevant search parameters.
- Verify examples such as UPSC and IT job links map to the correct route and query values.
- Test authenticated and unauthenticated navigation, unknown routes, back/forward navigation, and direct deep links.
- Complete the phase only after checking all visible navbar links and recording any intentionally external destinations.

## Phase 3: Government Jobs Pipeline

- Query the configured database for the real `jobs_government` count and inspect representative rows without exposing sensitive data.
- Trace crawler sources in `services/crawler-go`, scheduling, parsing, deduplication, database insertion, and failure logging for UPSC, SSC, RRB, and configured sources.
- Verify the backend `/api/v1/gov-jobs` contract, filters, pagination, ordering, and empty/error behavior.
- Verify the frontend maps the actual response shape, renders loading/error/empty/data states, and applies filters.
- Confirm Apply Now links use the persisted official URL, validate URL handling, and open in a new tab with appropriate security attributes.
- Do not seed synthetic jobs to hide an empty pipeline. If no upstream data is available, report the blocked boundary separately.

## Phase 4: Private Jobs Pipeline

- Check the real `jobs_private` data count and representative records.
- Trace RemoteOK, Adzuna, and other configured providers through credentials, requests, rate limits, parsing, normalization, deduplication, and persistence.
- Verify `/api/v1/priv-jobs` response shape, filters, pagination, and errors.
- Verify Location, Experience, and Salary controls trigger correctly encoded API requests and update the UI without stale results.
- Test provider failures, no-result filters, malformed records, and external Apply links.

## Phase 5: Courses Pipeline

- Locate the provider list and course card components plus the API/database contract.
- Verify Coursera, Udemy, Swayam, W3Schools, and configured providers render with aligned logos, labels, counts, loading states, and errors.
- Trace each course URL from the `courses` table or upstream API through serialization to the card.
- Fix missing or malformed anchors so clicking a course opens its real `course.url` safely and visibly; do not substitute placeholder URLs.
- Test keyboard activation, external navigation, invalid URLs, and an empty provider result.

## Phase 6: Videos Pipeline

- Search the actual Videos route and components for Netflix assets, labels, SVGs, or icon imports and replace only the incorrect branding with a standard YouTube treatment consistent with the design system.
- Trace video discovery through YouTube RSS or the configured YouTube Data API, including search terms such as government job tutorials and interview preparation, credentials, quotas, parsing, persistence, and refresh behavior.
- Verify the backend response maps `thumbnail_url`, title, channel, duration, and watch URL correctly.
- Verify the tutorial list renders real records with loading, error, and empty states.
- Test Watch behavior using either a safe YouTube link or an accessible modal with a correctly scoped iframe, close behavior, and responsive sizing.

## Phase 7: Government Forms Pipeline

- Trace Admit Cards, Results, and Notifications from configured crawlers or APIs through persistence and backend responses.
- Build or repair the frontend tabs/table/cards using the actual response fields and clear loading, error, and empty states.
- Verify each PDF or official notification URL is present, clickable, safe to open, and points to the real source.
- Test malformed URLs, unavailable documents, duplicate records, and mobile layout.

## Verification and Reporting

For every phase, report:

- Reproduction steps and the exact observed failure.
- Root cause and affected files or services.
- Real data/API/database evidence, with secrets redacted.
- Minimal code changes made.
- Focused commands or tests run and their outcomes.
- Remaining blockers, assumptions, and residual risks.
- The exact approval needed before proceeding.

Use a compact status format:

```text
Phase: <number and name>
Status: <blocked | repaired | verified>
Findings:
- ...
Changes:
- ...
Validation:
- <command>: <result>
Open risks:
- ...
Next gate: Reply `Approved` to begin Phase <next number>, or specify a correction.
```
