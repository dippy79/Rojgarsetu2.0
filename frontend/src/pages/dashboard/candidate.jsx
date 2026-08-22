import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'

const CandidateDashboard = dynamic(
  () => import('../candidate/CandidateDashboard'),
  { ssr: false }
)

export default function CandidateDashboardPage() {
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('token') || localStorage.getItem('rojgar_token')
    const role = localStorage.getItem('userRole') || localStorage.getItem('role')

    if (!token) {
      router.push('/login')
    } else if (role !== 'candidate') {
      router.push('/unauthorized')
    }
  }, [router])

  return <CandidateDashboard />
}
