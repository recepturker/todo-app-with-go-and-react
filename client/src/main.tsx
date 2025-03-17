import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { ChakraProvider } from "@chakra-ui/react"
import React from 'react'
import { Provider } from './components/ui/provider.tsx'

createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Provider>
      <App />
    </Provider>
  </React.StrictMode>,
)
