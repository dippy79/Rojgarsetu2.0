import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import { useAuth } from '../../hooks/useAuth'

const CandidateDashboard = dynamic(
  () => import('../candidate/CandidateDashboard'),
  { ssr: false }
)

export default function CandidateDashboardPage() {
  const router = useRouter()
  const { user, initialized } = useAuth()

  useEffect(() => {
    if (initialized) {
      if (!user) {
        router.push('/login')
      } else if (user.role !== 'candidate') {
        router.push('/unauthorized')
      }
    }
  }, [user, initialized, router])

  if (!initialized || !user || user.role !== 'candidate') {
    return null;
  }

  return <CandidateDashboard />
}
