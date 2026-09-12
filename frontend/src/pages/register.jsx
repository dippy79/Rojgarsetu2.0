import { useEffect } from 'react'
import { useRouter } from 'next/router'
import dynamic from 'next/dynamic'

const LoginPage = dynamic(
  () => import('./auth/LoginPage'),
  { ssr: false }
)

export default function Register() {
  const router = useRouter()
  useEffect(() => {
    router.replace('/login')
  }, [router])
  return <LoginPage />
}
