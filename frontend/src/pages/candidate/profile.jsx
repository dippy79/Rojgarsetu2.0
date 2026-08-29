import dynamic from 'next/dynamic'
import { useRouter } from 'next/router'
import { useEffect } from 'react'
import { useAuth } from '../../hooks/useAuth'

const CandidateProfile = dynamic(
  () => import('./CandidateProfile'),
  { ssr: false }
)

export default function CandidateProfilePage() {
  const router = useRouter()
  const { user, initialized } = useAuth()

  useEffect(() => {
    if (initialized && !user) {
      router.push('/login')
    }
  }, [user, initialized, router])

  if (!initialized || !user) return null

  return <CandidateProfile />
}
