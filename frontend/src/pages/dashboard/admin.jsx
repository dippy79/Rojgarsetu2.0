import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'

const AdminDashboard = dynamic(
  () => import('../admin/AdminDashboard'),
  { ssr: false }
)

export default function AdminDashboardPage() {
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('token') || localStorage.getItem('rojgar_token')
    const role = localStorage.getItem('userRole') || localStorage.getItem('role')

    if (!token) {
      router.push('/login')
    } else if (role !== 'admin') {
      router.push('/unauthorized')
    }
  }, [router])

  return <AdminDashboard />
}
