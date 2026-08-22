import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'

const EmployerDashboard = dynamic(
  () => import('../employer/EmployerDashboard'),
  { ssr: false }
)

export default function CompanyDashboardPage() {
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('token') || localStorage.getItem('rojgar_token')
    const role = localStorage.getItem('userRole') || localStorage.getItem('role')

    if (!token) {
      router.push('/login')
    } else if (role !== 'company' && role !== 'employer') {
      router.push('/unauthorized')
    }
  }, [router])

  return <EmployerDashboard />
}
