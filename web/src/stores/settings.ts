import { defineStore } from 'pinia'
import type { Settings } from '@/types'

const STORAGE_KEY = 'ibookfs_settings'

interface SettingsState {
  settings: Settings
  loading: boolean
  error: string | null
}

const defaultSettings: Settings = {
  theme: 'light',
  uploadQuality: 'high',
  autoOptimize: true,
  defaultStartPage: 1,
}

// Load settings from localStorage
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

// Save settings to localStorage
function saveToStorage(settings: Settings): void {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings))
  } catch (e) {
    console.error('Failed to save settings to localStorage:', e)
  }
}

export const useSettingsStore = defineStore('settings', {
  state: (): SettingsState => ({
    settings: loadFromStorage(),
    loading: false,
    error: null,
  }),

  getters: {
    currentTheme: (state) => state.settings.theme,
  },

  actions: {
    // Initialize settings from localStorage
    initSettings() {
      this.settings = loadFromStorage()
      this.applyTheme()
    },

    // Update settings and save to localStorage
    updateSettings(data: Partial<Settings>) {
      this.settings = { ...this.settings, ...data }
      saveToStorage(this.settings)
      this.applyTheme()
      return true
    },

    applyTheme() {
      const theme = this.settings.theme

      if (theme === 'system') {
        const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
        document.documentElement.setAttribute('data-theme', prefersDark ? 'dark' : 'light')
      } else {
        document.documentElement.setAttribute('data-theme', theme)
      }
    },

    setTheme(theme: Settings['theme']) {
      this.updateSettings({ theme })
    },

    setUploadQuality(quality: Settings['uploadQuality']) {
      this.updateSettings({ uploadQuality: quality })
    },

    setAutoOptimize(enabled: boolean) {
      this.updateSettings({ autoOptimize: enabled })
    },

    setDefaultStartPage(page: number) {
      this.updateSettings({ defaultStartPage: page })
    },
  },
})
