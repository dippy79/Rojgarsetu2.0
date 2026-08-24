import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import ProtectedRoute from '../../components/ProtectedRoute'

const JobApplicants = dynamic(
  () => import('./JobApplicants'),
  { ssr: false }
)

export default function ApplicantsPage() {
  return (
    <ProtectedRoute allowedRoles={['company', 'admin', 'employer']}>
      <JobApplicants />
    </ProtectedRoute>
  )
}
