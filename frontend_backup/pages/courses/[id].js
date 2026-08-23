import { useRouter } from 'next/router'
import useSWR from 'swr'
import { fetcher } from '../../lib/api'

export default function CoursePage(){
  const router = useRouter()
  const { id } = router.query
  const { data: course } = useSWR(id ? `/courses/${id}?lang=en` : null, fetcher)

  if (!course) return <div>Loading...</div>

  return (
    <div style={{padding:20}}>
      <h1>{course.name}</h1>
      <p>{course.fees_structure}</p>
      <a href={course.apply_link} target="_blank" rel="noreferrer">Apply</a>
    </div>
  )
}
