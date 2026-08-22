import Link from 'next/link'

export default function UnauthorizedPage() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-slate-50 text-slate-900 p-6 text-center">
      <div className="w-16 h-16 bg-red-100 text-red-600 rounded-full flex items-center justify-center mb-4 text-2xl font-bold">
        !
      </div>
      <h1 className="text-3xl font-extrabold mb-2">403 - Access Denied</h1>
      <p className="text-slate-600 max-w-md mb-6 text-sm">
        You do not have the required permissions to access this page. Please log in with an authorized account.
      </p>
      <Link
        href="/login"
        className="px-6 py-3 bg-slate-900 text-white font-semibold rounded-xl hover:bg-slate-800 transition-all text-sm shadow-sm"
      >
        Return to Login
      </Link>
    </div>
  )
}
