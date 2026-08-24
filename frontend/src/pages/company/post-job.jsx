import dynamic from 'next/dynamic'
import ProtectedRoute from '../../components/ProtectedRoute'

const PostJob = dynamic(
  () => import('./PostJob'),
  { ssr: false }
)

export default function PostJobPage() {
  return (
    <ProtectedRoute allowedRoles={['company', 'admin', 'employer']}>
      <PostJob />
    </ProtectedRoute>
  )
}
