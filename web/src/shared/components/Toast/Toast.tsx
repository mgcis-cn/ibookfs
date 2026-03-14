import { useToastStore } from '@/shared/stores/toast'
import { X, CheckCircle, AlertTriangle, AlertCircle, Info } from 'lucide-react'

const icons = {
  success: CheckCircle,
  warning: AlertTriangle,
  error: AlertCircle,
  info: Info,
}

const colors = {
  success: 'bg-green-50 border-green-200 text-green-800',
  warning: 'bg-amber-50 border-amber-200 text-amber-800',
  error: 'bg-red-50 border-red-200 text-red-800',
  info: 'bg-blue-50 border-blue-200 text-blue-800',
}

export function ToastContainer() {
  const toasts = useToastStore(s => s.toasts)
  const remove = useToastStore(s => s.remove)

  if (toasts.length === 0) return null

  return (
    <div className="fixed top-4 right-4 z-[800] flex flex-col gap-2">
      {toasts.map(toast => {
        const Icon = icons[toast.type]
        return (
          <div
            key={toast.id}
            className={`flex items-center gap-3 px-4 py-3 rounded-lg border shadow-lg animate-in slide-in-from-right ${colors[toast.type]}`}
          >
            <Icon className="w-5 h-5 shrink-0" />
            <span className="text-sm font-medium">{toast.message}</span>
            <button
              onClick={() => remove(toast.id)}
              className="ml-2 p-0.5 rounded hover:bg-black/10 transition-colors"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        )
      })}
    </div>
  )
}
