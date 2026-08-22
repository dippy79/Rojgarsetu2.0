import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'

const MyApplications = dynamic(
  () => import('./MyApplications'),
  { ssr: false }
)

export default function MyApplicationsPage() {
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('token') || localStorage.getItem('rojgar_token')
    if (!token) router.push('/login')
  }, [router])

  return <MyApplications />
}
