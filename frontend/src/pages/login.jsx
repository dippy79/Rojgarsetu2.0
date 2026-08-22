import dynamic from 'next/dynamic'

const LoginPage = dynamic(
  () => import('./auth/LoginPage'),
  { ssr: false }
)

export default function Login() {
  return <LoginPage />
}
