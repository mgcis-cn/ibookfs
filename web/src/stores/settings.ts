import { defineStore } from 'pinia'
import { api } from '@/api'
import type { Settings } from '@/types'

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

export const useSettingsStore = defineStore('settings', {
  state: (): SettingsState => ({
    settings: defaultSettings,
    loading: false,
    error: null,
  }),

  getters: {
    currentTheme: (state) => state.settings.theme,
  },

  actions: {
    async fetchSettings() {
      this.loading = true
      this.error = null

      try {
        const response = await api.getSettings()

        if (response.success) {
          this.settings = response.data
          this.applyTheme()
        } else {
          this.error = response.message || 'Failed to fetch settings'
        }
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'An error occurred'
      } finally {
        this.loading = false
      }
    },

    async updateSettings(data: Partial<Settings>) {
      this.loading = true
      this.error = null

      try {
        const response = await api.updateSettings(data)

        if (response.success) {
          this.settings = response.data
          this.applyTheme()
          return true
        } else {
          this.error = response.message || 'Failed to update settings'
          return false
        }
      } catch (err) {
        this.error = err instanceof Error ? err.message : 'An error occurred'
        return false
      } finally {
        this.loading = false
      }
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
      this.settings.theme = theme
      this.applyTheme()
      this.updateSettings({ theme })
    },

    setUploadQuality(quality: Settings['uploadQuality']) {
      this.settings.uploadQuality = quality
      this.updateSettings({ uploadQuality: quality })
    },

    setAutoOptimize(enabled: boolean) {
      this.settings.autoOptimize = enabled
      this.updateSettings({ autoOptimize: enabled })
    },

    setDefaultStartPage(page: number) {
      this.settings.defaultStartPage = page
      this.updateSettings({ defaultStartPage: page })
    },
  },
})
