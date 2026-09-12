import dynamic from 'next/dynamic'
import ProtectedRoute from '../../components/ProtectedRoute'

const AdminDashboard = dynamic(
  () => import('../admin/AdminDashboard'),
  { ssr: false }
)

export default function AdminDashboardPage() {
  return (
    <ProtectedRoute allowedRoles={['admin']}>
      <AdminDashboard />
    </ProtectedRoute>
  )
}
