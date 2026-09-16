# MASTER SYSTEM INVENTORY - ROJGARSETU 2.0

## 1. ROUTE INVENTORY (FRONTEND)
| ID | PATH | SOURCE FILE | DESCRIPTION | AUTH REQ |
|---|---|---|---|---|
| RT-001 | `/` | `src/pages/index.js` | Landing Page | No |
| RT-002 | `/login` | `src/pages/login.jsx` | Auth Entry | No |
| RT-003 | `/register` | `src/pages/register.jsx` | User Registration | No |
| RT-004 | `/dashboard/admin` | `src/pages/dashboard/admin.jsx` | Admin Hub | Yes (Admin) |
| RT-005 | `/dashboard/candidate` | `src/pages/dashboard/candidate.jsx` | Candidate Hub | Yes (Candidate) |
| RT-006 | `/dashboard/company` | `src/pages/dashboard/company.jsx` | Recruiter Hub | Yes (Company) |
| RT-007 | `/gov-jobs` | `src/pages/gov-jobs/index.jsx` | Government Jobs | No |
| RT-008 | `/private-jobs` | `src/pages/private-jobs/index.jsx` | Private Jobs | No |
| RT-009 | `/courses` | `src/pages/courses/index.jsx` | Courses | No |
| RT-010 | `/videos` | `src/pages/videos/index.jsx` | Video Tutorials | No |
| RT-011 | `/govt-forms` | `src/pages/govt-forms/index.jsx` | Official Forms | No |
| RT-012 | `/candidate/profile` | `src/pages/candidate/profile.jsx` | Profile Settings | Yes |
| RT-013 | `/candidate/applications` | `src/pages/candidate/applications.jsx` | Application History | Yes |
| RT-014 | `/jobs/[id]` | `src/pages/jobs/[id].js` | Job Details | No |

## 2. API INVENTORY (GATEWAY PROXIED)
| ID | ENDPOINT | SERVICE | AUTH | DESCRIPTION |
|---|---|---|---|---|
| API-001 | `/api/auth/login` | `auth-java` | No | Authenticate user |
| API-002 | `/api/v1/gov-jobs` | `backend_go` | No | Fetch govt jobs |
| API-003 | `/api/v1/priv-jobs` | `backend_go` | No | Fetch private jobs |
| API-004 | `/api/v1/courses` | `backend_go` | No | Fetch academic courses |
| API-005 | `/api/v1/videos` | `backend_go` | No | Fetch curated videos |
| API-006 | `/api/v1/stats` | `backend_go` | No | Public stats |
| API-007 | `/api/v1/candidates/me` | `backend_go` | Yes | Personal profile data |

## 3. DATABASE INVENTORY
| ID | TABLE | TYPE | DESCRIPTION |
|---|---|---|---|
| DB-001 | `users` | Core | Central identity table |
| DB-002 | `jobs_government` | Dynamic | Scraped govt data |
| DB-003 | `jobs_private` | Dynamic | Aggregated private data |
| DB-004 | `courses` | Content | Educational data |
| DB-005 | `youtube_videos` | Content | Curated media |
| DB-006 | `email_queue` | System | Background worker state |

## 4. INFRASTRUCTURE INVENTORY
| ID | COMPONENT | PORT | HEALTHCHECK |
|---|---|---|---|
| INF-001 | `api-gateway` | 3001 | `/health` |
| INF-002 | `backend` | 8083 | `/health` |
| INF-003 | `postgres` | 5432 | `pg_isready` |
| INF-004 | `redis` | 6379 | `redis-cli ping` |
| INF-005 | `crawler` | 8082 | `/health` |
| INF-006 | `ai-engine` | 8000 | `/health` |
