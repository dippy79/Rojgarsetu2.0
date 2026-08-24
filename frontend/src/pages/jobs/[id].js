import { useRouter } from 'next/router'
import { useEffect, useState } from 'react'
import { fetcher } from '../../lib/api'

export default function JobPage(){
  const router = useRouter()
  const { id } = router.query
  const [job, setJob] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (id) {
      fetcher(`/api/jobs/${id}`)
        .then(data => setJob(data))
        .catch(err => console.error(err))
        .finally(() => setLoading(false))
    }
  }, [id])

  if (loading) return <div className="p-10 text-center font-bold">Loading Position Data...</div>
  if (!job) return <div className="p-10 text-center font-bold">Position Not Found</div>

  return (
    <div className="max-w-4xl mx-auto p-10 font-sans">
      <h1 className="text-4xl font-black text-slate-900 mb-4">{job.title || job.data?.title}</h1>
      <div className="bg-white border border-slate-200 rounded-[2rem] p-8 shadow-sm">
         <p className="text-slate-600 leading-relaxed mb-8">{job.description || job.data?.description}</p>
         <a
           href={job.apply_link || job.data?.apply_link}
           target="_blank"
           rel="noreferrer"
           className="inline-block px-10 py-4 bg-slate-900 text-white font-black rounded-2xl hover:bg-indigo-600 transition-all uppercase text-xs tracking-widest"
         >
           Apply Now
         </a>
      </div>
    </div>
  )
}
