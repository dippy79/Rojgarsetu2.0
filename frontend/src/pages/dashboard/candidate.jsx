import dynamic from 'next/dynamic'
import ProtectedRoute from '../../components/ProtectedRoute'

const CandidateDashboard = dynamic(
  () => import('../candidate/CandidateDashboard'),
  { ssr: false }
)

export default function CandidateDashboardPage() {
  return (
    <ProtectedRoute allowedRoles={['candidate']}>
      <CandidateDashboard />
    </ProtectedRoute>
  )
}
