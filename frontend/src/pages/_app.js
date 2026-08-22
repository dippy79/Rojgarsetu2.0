import '../index.css'
import '../components/Navbar.css'
import '../components/JobCard.css'
import '../components/FilterBar.css'
import '../components/VideoCard.css'
import '../components/CourseCard.css'
import '../components/Pagination.css'
import '../components/GovtFormsDashboard.css'
import '../components/FloatingNotification.css'
import './courses/Courses.css'
import './gov-jobs/GovJobs.css'
import './private-jobs/PrivateJobs.css'
import './videos/Videos.css'
import { AuthProvider } from '../context/AuthContext'
import Navbar from '../components/Navbar'
import { useRouter } from 'next/router'

function MyApp({ Component, pageProps }) {
  const router = useRouter()
  const hideNavbar = ['/login'].includes(router.pathname)

  return (
    <AuthProvider>
      {!hideNavbar && <Navbar />}
      <Component {...pageProps} />
    </AuthProvider>
  )
}

export default MyApp
