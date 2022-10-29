import '../styles/globals.css'
import type { AppProps } from 'next/app'

import { PlantUMLProvider } from '../components/plantuml'

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <PlantUMLProvider value={process.env.NEXT_PUBLIC_PLANTUML_SERVER as string}>
      <Component {...pageProps} />
    </PlantUMLProvider>
  )
}

export default MyApp
