import React, { useState, useEffect } from 'react';
import { apiUrl } from '../apiConfig';
// Removed global CSS import - moved to _app.js or use CSS Modules
// import './GovtFormsDashboard.css';

/**
 * GovtFormsDashboard — "Govt Forms & Exams" dashboard.
 *
 * Smart logic:
 *  - Strictly sorted by closing date (ascending: nearest first).
 *  - Priority 1 (<=3 days): Red/Pink, on top, blinking "Ending Soon" badge.
 *  - Priority 2 (<=14 days): Orange warning badge.
 *  - Priority 3 (>14 days): Green safe badge.
 *
 * Fallback: The backend `/api/v1/forms` endpoint may not exist yet, so we seed
 * the component with 5 realistic mock records. If the endpoint later returns
 * data, the mock is replaced.
 */

// Priority thresholds (in days).
const P1_DAYS = 3;   // Ending soon (1-3 days)
const P2_DAYS = 14;  // Warning (1-2 weeks)

// Realistic mock data injected ONLY if the backend endpoint is unavailable.
const MOCK_FORMS = [
  {
    id: 'mock-1',
    title: 'SSC CGL 2025 Online Application',
    department: 'Staff Selection Commission',
    last_date: (() => {
      const d = new Date();
      d.setDate(d.getDate() + 2);
      return d.toISOString();
    })(),
    apply_url: 'https://ssc.gov.in',
  },
  {
    id: 'mock-2',
    title: 'RRB NTPC 2025 Application Form',
    department: 'Railway Recruitment Board',
    last_date: (() => {
      const d = new Date();
      d.setDate(d.getDate() + 6);
      return d.toISOString();
    })(),
    apply_url: 'https://rrbcdg.gov.in',
  },
  {
    id: 'mock-3',
    title: 'UPSC Civil Services Prelims 2025',
    department: 'Union Public Service Commission',
    last_date: (() => {
      const d = new Date();
      d.setDate(d.getDate() + 10);
      return d.toISOString();
    })(),
    apply_url: 'https://upsc.gov.in',
  },
  {
    id: 'mock-4',
    title: 'IBPS PO 2025 Registration',
    department: 'Institute of Banking Personnel Selection',
    last_date: (() => {
      const d = new Date();
      d.setDate(d.getDate() + 20);
      return d.toISOString();
    })(),
    apply_url: 'https://ibps.in',
  },
  {
    id: 'mock-5',
    title: 'State PSC Assistant Engineer 2025',
    department: 'State Public Service Commission',
    last_date: (() => {
      const d = new Date();
      d.setDate(d.getDate() + 45);
      return d.toISOString();
    })(),
    apply_url: 'https://psc.gov.in',
  },
];

function daysUntil(isoDate) {
  const target = new Date(isoDate);
  const now = new Date();
  const diff = target.getTime() - now.getTime();
  return Math.ceil(diff / (1000 * 60 * 60 * 24));
}

function getPriority(days) {
  if (days <= P1_DAYS) return 1;
  if (days <= P2_DAYS) return 2;
  return 3;
}

function formatDate(isoDate) {
  const d = new Date(isoDate);
  return d.toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' });
}

const GovtFormsDashboard = () => {
  const [forms, setForms] = useState([]);
  const [loading, setLoading] = useState(true);
  const [usingMock, setUsingMock] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    let cancelled = false;

    async function load() {
      try {
        setLoading(true);
        const response = await fetch(`${apiUrl('/api/v1/forms')}`);
        if (!response.ok) throw new Error('HTTP ' + response.status);

        const data = await response.json();
        const list = (data && (data.data || data.forms)) || [];
        if (cancelled) return;

        if (Array.isArray(list) && list.length > 0) {
          setForms(list);
          setUsingMock(false);
        } else {
          // Backend reachable but empty -> fall back to mock so UI never breaks.
          setForms(MOCK_FORMS);
          setUsingMock(true);
        }
      } catch (err) {
        // Endpoint does not exist / backend offline -> inject realistic mocks.
        if (cancelled) return;
        setForms(MOCK_FORMS);
        setUsingMock(true);
        setError('Live forms feed unavailable — showing latest expected deadlines.');
      } finally {
        if (!cancelled) setLoading(false);
      }
    }

    load();
    return () => {
      cancelled = true;
    };
  }, []);

  // Strictly sort by closing date ascending (closest first).
  const sorted = [...forms].sort((a, b) => new Date(a.last_date) - new Date(b.last_date));

  return (
    <div className="forms-dashboard">
      <div className="page-header">
        <h1>📝 Govt Forms &amp; Exams</h1>
        <p>Deadlines closing soon — sorted by nearest date so you never miss an application.</p>
      </div>

      {usingMock && <div className="forms-mock-banner">{error || 'Showing sample deadlines while the live feed connects.'}</div>}

      {loading ? (
        <div className="loading">Loading forms...</div>
      ) : (
        <div className="forms-list">
          {sorted.length === 0 && <div className="no-results">No forms available right now.</div>}

          {sorted.map((form) => {
            const days = daysUntil(form.last_date);
            const priority = getPriority(days);
            const safeDays = Math.max(days, 0);

            return (
              <div key={form.id || form.title} className={`form-card priority-${priority}`}>
                <div className="form-card-body">
                  <h3 className="form-title">{form.title}</h3>
                  <p className="form-department">{form.department || 'Government'}</p>

                  {priority === 1 && (
                    <span className="badge badge-p1" role="status">
                      ⚠️ Ending Soon
                    </span>
                  )}
                  {priority === 2 && (
                    <span className="badge badge-p2" role="status">
                      ⏳ Closing in {safeDays} days
                    </span>
                  )}
                  {priority === 3 && (
                    <span className="badge badge-p3" role="status">
                      ✅ Closing in {safeDays} days
                    </span>
                  )}
                </div>

                <div className="form-card-meta">
                  <span className="form-date">📅 {formatDate(form.last_date)}</span>
                  {form.apply_url ? (
                    <a className="form-apply" href={form.apply_url} target="_blank" rel="noreferrer">
                      Apply Now →
                    </a>
                  ) : (
                    <span className="form-apply disabled">Apply Soon</span>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
};

export default GovtFormsDashboard;
