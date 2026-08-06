# Security Policy — RojgarSetu 2.0 (v0.2.1)

At **RojgarSetu**, we take the security of our platform, user data, and automated crawlers seriously. 
This security policy applies to version `v0.2.1` and outlines 
supported versions, vulnerability reporting procedures, and standard security guidelines for contributors and maintainers.

---

## 1. Supported Versions

We actively maintain and provide security patches only for the current stable minor release line.

| Version | Supported          | Status |
| :--- | :---: | :--- |
| **v0.2.1** | Yes | Current Active Release (Gated Security Patches) |
| v0.2.0 | No | End of Support (Upgrade to v0.2.1 recommended) |
| < v0.2.0 | No | End of Life (EOL) |

---

## 2. Reporting a Vulnerability

**DO NOT** create a public GitHub Issue or Pull Request to report a security vulnerability or leaked credential.

### Reporting Procedure
If you discover a security flaw or credential exposure in RojgarSetu 2.0, please report it privately:

1. **Email:** Send details to `security@rojgarsetu.dev` (or the repository owner's security contact).
2. **Details to Include:**
   * Affected service (`crawler-go`, `backend_go`, `ai-engine-python`, `java-auth-service`, `frontend`, or `gateway`).
   * Description of the vulnerability or proof-of-concept (PoC).
   * Steps to reproduce the issue safely.
   * Potential impact on systems or candidate data.

### Response SLA
* **Initial Acknowledgment:** Within **24–48 hours**.
* **Status Update & Remediation Plan:** Within **5 business days**.
* **Responsible Disclosure Period:** We request a **30-day embargo** prior to public disclosure while a fix is deployed and verified.

---

## 3. Core Security Standards & Practices (v0.2.1)

### A. Secret Management & Credential Sanitization
* **Zero Hardcoded Secrets:** No API keys, database connection strings, or service tokens may be hardcoded in any version-controlled file.
* **Environment Variables:** All secrets must be loaded via `.env` or system environment variables (`GOOGLE_API_KEY`, `POSTGRES_PASSWORD`, `REDIS_PASSWORD`).
* **Automated Scanning:** All pull requests and commits are scanned using `Gitleaks` pre-commit hooks and GitHub Secret Scanning. Leaked keys are immediately revoked in their respective cloud providers (GCP, AWS, etc.).

### B. Database & Query Safety
* **Parameterized Queries:** SQL injection prevention is strictly enforced using `sqlc` in `backend_go` and explicit bound parameters in `ai-engine-python`.
* **Migration Integrity:** Schema updates must be executed using version-tracked SQL migrations (`backend_go/migrations/`).

### C. Web Crawling & Rate Limiting (Crawler Integrity)
* **Respectful Ingestion:** `crawler-go` enforces rate-limiting and domain-level circuit breakers (e.g., `gobreaker` for ATS portals like Greenhouse/Lever).
* **WAF/Bot Mitigation:** Requests must handle HTTP `429` (Rate Limited) and `403` (Forbidden) gracefully without triggering infinite retries or cascade crashes.

### D. Authentication & Authorization
* Token issuance and identity checks are centralized in the Java Auth module and enforced at the Gateway level using signed JWTs.

---

## 4. Security Checklist for Developers & AI Agents

Before submitting code or triggering automated builds in `v0.2.1`:

- [ ] Ran `gitleaks detect --source . -v` locally to ensure no secrets/keys are staged.
- [ ] Confirmed `.env` and local credential overrides are ignored in `.gitignore`.
- [ ] Confirmed all new SQL queries use parameterized bindings.
- [ ] Verified that external API endpoints handle errors gracefully (no raw stack traces exposed to clients).

---

## 5. Security Hall of Fame & Acknowledgments

We appreciate the efforts of security researchers and open-source contributors who help keep RojgarSetu secure and reliable. Valid, responsibly disclosed security reports will be acknowledged in our release notes.
