import React from 'react';
import { act, render, screen } from '@testing-library/react';
import Home from '../pages';

jest.mock('next/link', () => ({
  __esModule: true,
  default: function MockLink({ href, children, ...props }) {
    return require('react').createElement('a', { href, ...props }, children);
  },
}));

describe('Home page', () => {
  let resolveStatsRequest;

  beforeEach(() => {
    global.fetch = jest.fn(
      () =>
        new Promise((resolve) => {
          resolveStatsRequest = resolve;
        })
    );
  });

  it('renders the career discovery message and core opportunity links', () => {
    render(<Home />);

    expect(
      screen.getByRole('heading', { name: /your gateway to digital empowerment/i })
    ).toBeInTheDocument();
    expect(screen.getByText(/AI-Powered Career Portal for India/i)).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Browse Gov Jobs' })).toHaveAttribute(
      'href',
      '/gov-jobs'
    );
    expect(screen.getByRole('link', { name: 'Private Sector' })).toHaveAttribute(
      'href',
      '/private-jobs'
    );
  });

  it('renders routes to learning resources and account creation', () => {
    render(<Home />);

    expect(screen.getByRole('link', { name: /courses/i })).toHaveAttribute('href', '/courses');
    expect(screen.getByRole('link', { name: /video prep/i })).toHaveAttribute('href', '/videos');
    expect(screen.getByRole('link', { name: /join thousands of candidates today/i })).toHaveAttribute(
      'href',
      '/register'
    );
  });

  afterEach(async () => {
    if (resolveStatsRequest) {
      await act(async () => {
        resolveStatsRequest({
          ok: true,
          json: async () => ({ data: {} }),
        });
      });
      resolveStatsRequest = undefined;
    }
    jest.restoreAllMocks();
  });
});