import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import { useAuth } from '../../hooks/useAuth'

const MyApplications = dynamic(
  () => import('./MyApplications'),
  { ssr: false }
)

export default function MyApplicationsPage() {
  const router = useRouter()
  const { user, initialized } = useAuth()

  useEffect(() => {
    if (initialized && !user) {
      router.push('/login')
    }
  }, [user, initialized, router])

  if (!initialized || !user) return null

  return <MyApplications />
}
