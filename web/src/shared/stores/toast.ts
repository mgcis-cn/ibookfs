import { create } from 'zustand'

interface ToastItem {
  id: number
  type: 'success' | 'warning' | 'error' | 'info'
  message: string
}

interface ToastState {
  toasts: ToastItem[]
  add: (type: ToastItem['type'], message: string, duration?: number) => void
  remove: (id: number) => void
}

let nextId = 0

export const useToastStore = create<ToastState>((set, get) => ({
  toasts: [],
  add: (type, message, duration = 3000) => {
    const id = nextId++
    set({ toasts: [...get().toasts, { id, type, message }] })
    if (duration > 0) {
      setTimeout(() => get().remove(id), duration)
    }
  },
  remove: (id) => set({ toasts: get().toasts.filter(t => t.id !== id) }),
}))
