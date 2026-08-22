import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'

const EmployerDashboard = dynamic(
  () => import('./EmployerDashboard'),
  { ssr: false }
)

export default function EmployerDashboardPage() {
  const router = useRouter()

  useEffect(() => {
    const token = localStorage.getItem('token') || localStorage.getItem('rojgar_token')
    if (!token) router.push('/login')
  }, [router])

  return <EmployerDashboard />
}
