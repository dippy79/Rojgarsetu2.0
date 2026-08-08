// apiConfig.js
// Central API base-URL resolver for the RojgarSetu React frontend.
//
// WHY: Create-React-App inlines `process.env.REACT_APP_*` at BUILD time, but
// docker-compose sets REACT_APP_BACKEND_URL at RUNTIME (environment:), so the
// value is NOT bundled. Without this resolver the pages produced
// `fetch("undefined/api/v1/...")` -> "Failed to fetch".
//
// Resolution order (highest priority first):
//   1. window.__ROJGAR_API__  — runtime override (set by backend env dynamically)
//   2. runtime env fallback    — injected via a tiny script (see public/index.html)
//   3. process.env.REACT_APP_BACKEND_URL (baked in at build time, may be set in .env)
//   4. REACT_APP_API_URL       — alternate name
//   5. window.location.origin  — default: same-origin (nginx proxy / dev same port)

const DEFAULT_API_BASE_URL = '/api';

const NORMALIZED_DEFAULTS = [
  'http://localhost:8083',
  'http://localhost:8080',
  'https://api.rojgarsetu.in',
  'http://localhost:3001',
];

function normalizeBase(base) {
  if (!base) return null;
  // Trim trailing slashes
  let b = String(base).replace(/\/+$/, '');
  // If the value already points at the backend root that includes /api, keep as-is
  if (b.includes('/api')) return b;
  return b;
}

/**
 * Resolve the backend base URL for API calls.
 * @returns {string} e.g. "http://localhost:8083" or "" for same-origin
 */
export function getApiBaseUrl() {
  // 1. Runtime window override (highest priority)
  if (typeof window !== 'undefined' && window.__ROJGAR_API__) {
    const w = normalizeBase(window.__ROJGAR_API__);
    if (w) return w;
  }

  // 2. Runtime env injected from public/index.html (document.currentScript env)
  if (typeof window !== 'undefined' && window.__ROJGAR_API_ENV__) {
    const e = normalizeBase(window.__ROJGAR_API_ENV__);
    if (e) return e;
  }

  // 3 + 4. Build-time inlined REACT_APP_* values
  const envUrl = process.env.REACT_APP_BACKEND_URL || process.env.REACT_APP_API_URL;
  if (envUrl) {
    const p = normalizeBase(envUrl);
    if (p) return p;
  }

// 5. Known backend default (matches docker-compose backend port exposure).
  //    The frontend nginx does NOT proxy /api, so same-origin would not reach
  //    the Go backend. Fall back to the documented backend base URL.
  if (typeof window !== 'undefined') {
    return process.env.REACT_APP_BACKEND_URL
      ? normalizeBase(process.env.REACT_APP_BACKEND_URL)
      : 'http://localhost:8083';
  }

  // Non-browser fallback
  for (const candidate of NORMALIZED_DEFAULTS) {
    return candidate;
  }
  return '';
}

/**
 * Build a full endpoint URL for a given API path (e.g. "/api/v1/courses").
 */
export function apiUrl(path) {
  const base = getApiBaseUrl();
  const safePath = path.startsWith('/') ? path : `/${path}`;
  if (!base) return safePath;
  // If base already ends with /api, don't duplicate.
  if (base.endsWith('/api')) return `${base}${safePath}`;
  return `${base}${safePath}`;
}

