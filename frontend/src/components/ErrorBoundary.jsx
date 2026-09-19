import React from 'react';

class ErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error) {
    // Update state so the next render will show the fallback UI.
    return { hasError: true, error };
  }

  componentDidCatch(error, errorInfo) {
    // You can also log the error to an error reporting service
    console.error("[ErrorBoundary] Caught error:", error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      // You can render any custom fallback UI
      return (
        <div className="min-h-screen flex items-center justify-center bg-slate-50 px-6 py-24">
          <div className="max-w-md w-full bg-white rounded-3xl shadow-xl border border-slate-200 p-10 text-center space-y-6">
            <div className="w-20 h-20 bg-rose-50 text-rose-500 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="w-10 h-10" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </div>
            <h2 className="text-2xl font-black text-slate-900 tracking-tight">Something went wrong.</h2>
            <p className="text-slate-500 font-medium leading-relaxed">
              We encountered an unexpected error while rendering this page. Our engineers have been notified.
            </p>
            <div className="pt-4 flex flex-col gap-3">
              <button
                onClick={() => window.location.reload()}
                className="w-full py-4 bg-slate-900 text-white font-black rounded-2xl hover:bg-slate-800 transition-all shadow-xl shadow-slate-900/10 uppercase tracking-widest text-[11px]"
              >
                Refresh Page
              </button>
              <button
                onClick={() => window.location.href = '/'}
                className="w-full py-4 bg-white text-slate-900 border-2 border-slate-100 font-black rounded-2xl hover:border-slate-200 transition-all uppercase tracking-widest text-[11px]"
              >
                Go to Homepage
              </button>
            </div>
            {process.env.NODE_ENV === 'development' && (
              <details className="mt-8 text-left bg-slate-50 p-4 rounded-xl border border-slate-100 overflow-hidden">
                <summary className="text-[10px] font-black uppercase text-slate-400 cursor-pointer">Technical Details</summary>
                <pre className="mt-4 text-[10px] text-rose-600 font-mono overflow-auto max-h-48 whitespace-pre-wrap">
                  {this.state.error?.toString()}
                </pre>
              </details>
            )}
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
