import i18n from 'i18next'
import { initReactI18next } from 'react-i18next'

const resources = {
  en: { translation: { welcome: 'Welcome' } },
  hi: { translation: { welcome: 'स्वागत है' } },
  ta: { translation: { welcome: 'வரவேற்பு' } }
}

i18n.use(initReactI18next).init({
  resources,
  lng: 'en',
  interpolation: { escapeValue: false }
})

export default i18n
