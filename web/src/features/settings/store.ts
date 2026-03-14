import { create } from 'zustand'

const STORAGE_KEY = 'ibookfs_settings'

export interface Settings {
  theme: 'light' | 'dark' | 'system'
  uploadQuality: 'original' | 'high' | 'standard'
  autoOptimize: boolean
  defaultStartPage: number
}

const defaultSettings: Settings = {
  theme: 'light',
  uploadQuality: 'high',
  autoOptimize: true,
  defaultStartPage: 1,
}

function loadFromStorage(): Settings {
  try {
    const stored = localStorage.getItem(STORAGE_KEY)
    if (stored) {
      return { ...defaultSettings, ...JSON.parse(stored) }
    }
  } catch (e) {
    console.error('Failed to load settings from localStorage:', e)
  }
  return defaultSettings
}

function saveToStorage(settings: Settings): void {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings))
  } catch (e) {
    console.error('Failed to save settings to localStorage:', e)
  }
}

function applyTheme(theme: Settings['theme']): void {
  if (theme === 'system') {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
    document.documentElement.setAttribute('data-theme', prefersDark ? 'dark' : 'light')
  } else {
    document.documentElement.setAttribute('data-theme', theme)
  }
}

interface SettingsState {
  settings: Settings
  initSettings: () => void
  updateSettings: (data: Partial<Settings>) => void
}

export const useSettingsStore = create<SettingsState>((set, get) => ({
  settings: loadFromStorage(),

  initSettings: () => {
    const s = loadFromStorage()
    set({ settings: s })
    applyTheme(s.theme)
  },

  updateSettings: (data) => {
    const updated = { ...get().settings, ...data }
    set({ settings: updated })
    saveToStorage(updated)
    applyTheme(updated.theme)
  },
}))
