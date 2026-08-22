import '../styles/globals.css'
import { appWithTranslation } from 'react-i18next'
import '../lib/i18n'

function MyApp({ Component, pageProps }) {
  return <Component {...pageProps} />
}

export default appWithTranslation(MyApp)
