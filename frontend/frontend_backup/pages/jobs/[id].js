import { useRouter } from 'next/router'
import useSWR from 'swr'
import { fetcher } from '../../lib/api'

export default function JobPage(){
  const router = useRouter()
  const { id } = router.query
  const { data: job } = useSWR(id ? `/jobs/${id}?lang=en` : null, fetcher)

  if (!job) return <div>Loading...</div>

  return (
    <div style={{padding:20}}>
      <h1>{job.title}</h1>
      <p>{job.criteria}</p>
      <a href={job.apply_link} target="_blank" rel="noreferrer">Apply</a>
    </div>
  )
}
