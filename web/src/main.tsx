import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import App from './app/App'
import './shared/styles/index.css'
import { useSettingsStore } from './features/settings/store'

// Apply saved settings (theme, etc.) on app start
useSettingsStore.getState().initSettings()

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>,
)
