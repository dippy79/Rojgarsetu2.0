import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import { useAuth } from '../../hooks/useAuth'

const AdminDashboard = dynamic(
  () => import('../admin/AdminDashboard'),
  { ssr: false }
)

export default function AdminDashboardPage() {
  const router = useRouter()
  const { user, initialized } = useAuth()

  useEffect(() => {
    if (initialized) {
      if (!user) {
        router.push('/login')
      } else if (user.role !== 'admin') {
        router.push('/unauthorized')
      }
    }
  }, [user, initialized, router])

  if (!initialized || !user || user.role !== 'admin') {
    return null; // Or a loading spinner
  }

  return <AdminDashboard />
}
