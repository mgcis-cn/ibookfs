import { defineStore } from 'pinia'
import type { Toast } from '@/types'

interface ToastState {
  toasts: Toast[]
  nextId: number
}

export const useToastStore = defineStore('toast', {
  state: (): ToastState => ({
    toasts: [],
    nextId: 1,
  }),

  actions: {
    addToast(message: string, type: Toast['type'] = 'info', duration: number = 3000) {
      const id = this.nextId++
      const toast: Toast = { id, message, type, duration }
      this.toasts.push(toast)

      if (duration > 0) {
        setTimeout(() => {
          this.removeToast(id)
        }, duration)
      }

      return id
    },

    removeToast(id: number) {
      const index = this.toasts.findIndex((t) => t.id === id)
      if (index !== -1) {
        this.toasts.splice(index, 1)
      }
    },

    clearAll() {
      this.toasts = []
    },

    success(message: string, duration?: number) {
      return this.addToast(message, 'success', duration)
    },

    warning(message: string, duration?: number) {
      return this.addToast(message, 'warning', duration)
    },

    error(message: string, duration?: number) {
      return this.addToast(message, 'error', duration)
    },

    info(message: string, duration?: number) {
      return this.addToast(message, 'info', duration)
    },
  },
})
