import dynamic from 'next/dynamic'
import ProtectedRoute from '../../components/ProtectedRoute'

const CompanyDashboard = dynamic(
  () => import('../company/CompanyDashboard'),
  { ssr: false }
)

export default function CompanyDashboardPage() {
  return (
    <ProtectedRoute allowedRoles={['company', 'admin', 'employer']}>
      <CompanyDashboard />
    </ProtectedRoute>
  )
}
