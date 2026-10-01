import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import { ShieldCheck, Building2, Calendar, Users, ExternalLink, Loader2 } from 'lucide-react'
import api from '../../lib/api'

export default function JobPage(){
  const router = useRouter()
  const { id } = router.query
  const [job, setJob] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)
  const [jobType, setJobType] = useState('government')

  useEffect(() => {
    if (id) {
      const fetchJob = async () => {
        setLoading(true)
        setError(null)
        try {
          // Try government first, then private
          let res
          try {
            res = await api.get(`/api/v1/gov-jobs/${id}`)
            setJobType('government')
          } catch (govErr) {
            try {
              res = await api.get(`/api/v1/priv-jobs/${id}`)
              setJobType('private')
            } catch (privErr) {
              throw new Error('Job not found')
            }
          }
          setJob(res.data?.data || res.data)
        } catch (err) {
          console.error(err)
          setError('Position not found')
        } finally {
          setLoading(false)
        }
      }
      fetchJob()
    }
  }, [id])

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto p-10 font-sans">
        <div className="flex flex-col items-center justify-center py-20 space-y-4">
          <Loader2 className="w-12 h-12 text-stone-800 animate-spin" />
          <p className="text-stone-500 font-bold uppercase tracking-widest text-[10px]">Loading Position Data...</p>
        </div>
      </div>
    )
  }

  if (error || !job) {
    return (
      <div className="max-w-4xl mx-auto p-10 font-sans">
        <div className="flex flex-col items-center justify-center py-20 text-center">
          <div className="w-16 h-16 bg-stone-100 rounded-full flex items-center justify-center mb-4">
            <ShieldCheck className="w-8 h-8 text-stone-400" />
          </div>
          <h3 className="text-2xl font-bold text-stone-900 mb-2">Position Not Found</h3>
          <p className="text-stone-500 text-sm">This job may have been removed or the ID is invalid</p>
        </div>
      </div>
    )
  }

  const isGov = jobType === 'government'
  const applyUrl = isGov
    ? (job?.apply_url || job?.apply_link || job?.url)
    : (job?.url || job?.apply_url || job?.apply_link)

  const sourceName = isGov
    ? (job?.source || job?.department || 'Government Portal')
    : (job?.company || job?.company_name || 'Private Company')

  const extractDomain = (url) => {
    if (!url) return 'Official Website'
    try {
      const domain = new URL(url).hostname
      return domain.replace('www.', '')
    } catch {
      return 'Official Website'
    }
  }

  return (
    <div className="max-w-4xl mx-auto p-10 font-sans">
      {/* Header with Badge */}
      <div className="mb-6">
        {isGov && (
          <div className="inline-flex items-center gap-2 px-4 py-2 bg-emerald-100 text-emerald-800 rounded-full text-xs font-black uppercase tracking-widest mb-4">
            🏛 Official Govt Source
          </div>
        )}
        <h1 className="text-4xl md:text-5xl font-black text-stone-900 mb-4 tracking-tight leading-tight">
          {job.title || job.data?.title}
        </h1>
        <div className="flex items-center gap-4 text-sm text-stone-500 font-medium">
          <span className="flex items-center gap-2">
            {isGov ? <ShieldCheck className="w-4 h-4" /> : <Building2 className="w-4 h-4" />}
            {sourceName}
          </span>
          {job?.location && (
            <span>• {job.location}</span>
          )}
        </div>
      </div>

      <div className="bg-white border border-stone-200 rounded-[2rem] p-8 shadow-sm space-y-8">
        {/* Description */}
        <div>
          <h2 className="text-lg font-black text-stone-900 mb-4 uppercase tracking-widest">Description</h2>
          <p className="text-stone-600 leading-relaxed">{job.description || job.data?.description || 'No description available'}</p>
        </div>

        {/* Government Specific Fields */}
        {isGov && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {job?.eligibility && (
              <div className="bg-stone-50 rounded-xl p-4">
                <div className="flex items-center gap-2 text-[10px] font-black text-stone-400 uppercase tracking-widest mb-2">
                  <ShieldCheck className="w-3 h-3" />
                  Eligibility
                </div>
                <p className="text-sm text-stone-700 font-medium">{job.eligibility}</p>
              </div>
            )}
            {job?.vacancy_count && (
              <div className="bg-stone-50 rounded-xl p-4">
                <div className="flex items-center gap-2 text-[10px] font-black text-stone-400 uppercase tracking-widest mb-2">
                  <Users className="w-3 h-3" />
                  Vacancies
                </div>
                <p className="text-sm text-stone-700 font-medium">{job.vacancy_count}</p>
              </div>
            )}
            {job?.exam_date && (
              <div className="bg-stone-50 rounded-xl p-4">
                <div className="flex items-center gap-2 text-[10px] font-black text-stone-400 uppercase tracking-widest mb-2">
                  <Calendar className="w-3 h-3" />
                  Exam Date
                </div>
                <p className="text-sm text-stone-700 font-medium">{new Date(job.exam_date).toLocaleDateString()}</p>
              </div>
            )}
          </div>
        )}

        {/* Application Destination Section */}
        <div className="border-t border-stone-100 pt-8">
          <h2 className="text-lg font-black text-stone-900 mb-4 uppercase tracking-widest">Application Destination</h2>
          <div className="bg-stone-50 rounded-xl p-6 flex items-center justify-between">
            <div>
              <p className="text-sm font-bold text-stone-900">{sourceName}</p>
              <p className="text-xs text-stone-500">{extractDomain(applyUrl)}</p>
            </div>
            {applyUrl ? (
              <a
                href={applyUrl}
                target="_blank"
                rel="noreferrer"
                className="inline-flex items-center gap-2 px-6 py-3 bg-stone-900 text-white font-black rounded-xl hover:bg-stone-800 transition-all text-xs uppercase tracking-widest"
              >
                Go to Official Application <ExternalLink className="w-4 h-4" />
              </a>
            ) : (
              <span className="text-stone-400 text-sm">Application URL not available</span>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
