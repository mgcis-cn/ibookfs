import { useToastStore } from '@/shared/stores/toast'
import { useSettingsStore } from '@/features/settings/store'
import type { Settings } from '@/features/settings/store'
import './SettingsPage.css'

export default function SettingsPage() {
  const toast = useToastStore(s => s.add)
  const { settings, updateSettings } = useSettingsStore()

  const update = (patch: Partial<Settings>, msg: string) => {
    updateSettings(patch)
    toast('success', msg)
  }

  return (
    <div className="page-settings">
      <h1 className="page-title">设置</h1>

      <div className="settings-sections">
        {/* Appearance */}
        <section className="settings-section">
          <h2 className="section-title">外观</h2>
          <div className="setting-item">
            <div className="setting-info">
              <label>主题模式</label>
              <p className="setting-description">选择浅色或深色主题</p>
            </div>
            <select className="setting-select" value={settings.theme} onChange={e => update({ theme: e.target.value as Settings['theme'] }, '主题设置已更新')}>
              <option value="light">浅色</option>
              <option value="dark">深色</option>
              <option value="system">跟随系统</option>
            </select>
          </div>
        </section>

        {/* Upload */}
        <section className="settings-section">
          <h2 className="section-title">上传</h2>
          <div className="setting-item">
            <div className="setting-info">
              <label>图片质量</label>
              <p className="setting-description">上传时的压缩质量</p>
            </div>
            <select className="setting-select" value={settings.uploadQuality} onChange={e => update({ uploadQuality: e.target.value as Settings['uploadQuality'] }, '上传质量设置已更新')}>
              <option value="original">原图</option>
              <option value="high">高质量</option>
              <option value="standard">标准</option>
            </select>
          </div>
          <div className="setting-item">
            <div className="setting-info">
              <label>默认起始页码</label>
              <p className="setting-description">新建书籍时的默认起始页码</p>
            </div>
            <input type="number" min={1} className="setting-input" value={settings.defaultStartPage} onChange={e => update({ defaultStartPage: Number(e.target.value) }, '默认起始页码已更新')} />
          </div>
          <div className="setting-item">
            <div className="setting-info">
              <label>自动优化图片</label>
              <p className="setting-description">上传后自动进行图片优化处理</p>
            </div>
            <label className="toggle-switch">
              <input type="checkbox" checked={settings.autoOptimize} onChange={e => update({ autoOptimize: e.target.checked }, '自动优化设置已更新')} />
              <span className="toggle-slider" />
            </label>
          </div>
        </section>

        {/* About */}
        <section className="settings-section">
          <h2 className="section-title">关于</h2>
          <div className="setting-item">
            <div className="setting-info">
              <label>版本</label>
              <p className="setting-description">iBookFS v1.0.0</p>
            </div>
          </div>
          <div className="setting-item">
            <div className="setting-info">
              <label>开源协议</label>
              <p className="setting-description">MIT License</p>
            </div>
          </div>
        </section>
      </div>
    </div>
  )
}
