<template>
  <div class="page-settings">
    <h1 class="page-title">设置</h1>

    <div class="settings-sections">
      <!-- Appearance Settings -->
      <section class="settings-section">
        <h2 class="section-title">外观</h2>
        <div class="setting-item">
          <div class="setting-info">
            <label>主题模式</label>
            <p class="setting-description">选择浅色或深色主题</p>
          </div>
          <select id="theme-setting" name="theme-setting" v-model="localSettings.theme" 
            class="setting-select" @change="handleThemeChange">
            <option value="light">浅色</option>
            <option value="dark">深色</option>
            <option value="system">跟随系统</option>
          </select>
        </div>
      </section>

      <!-- Upload Settings -->
      <section class="settings-section">
        <h2 class="section-title">上传</h2>
        <div class="setting-item">
          <div class="setting-info">
            <label>图片质量</label>
            <p class="setting-description">上传时的压缩质量</p>
          </div>
          <select
            id="quality-setting" name="quality-setting"
            v-model="localSettings.uploadQuality"
            class="setting-select"
            @change="handleUploadQualityChange"
          >
            <option value="original">原图</option>
            <option value="high">高质量</option>
            <option value="standard">标准</option>
          </select>
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label>默认起始页码</label>
            <p class="setting-description">新建书籍时的默认起始页码</p>
          </div>
          <input
            v-model.number="localSettings.defaultStartPage"
            type="number"
            min="1"
            class="setting-input"
            @change="handleStartPageChange"
          />
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label>自动优化图片</label>
            <p class="setting-description">上传后自动进行图片优化处理</p>
          </div>
          <label class="toggle-switch">
            <input
              v-model="localSettings.autoOptimize"
              type="checkbox"
              @change="handleAutoOptimizeChange"
            />
            <span class="toggle-slider"></span>
          </label>
        </div>
      </section>

      <!-- About -->
      <section class="settings-section">
        <h2 class="section-title">关于</h2>
        <div class="setting-item">
          <div class="setting-info">
            <label>版本</label>
            <p class="setting-description">iBookFS v1.0.0</p>
          </div>
        </div>
        <div class="setting-item">
          <div class="setting-info">
            <label>开源协议</label>
            <p class="setting-description">MIT License</p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useToastStore } from '@/stores/toast'
import type { Settings } from '@/types'

const settingsStore = useSettingsStore()
const toastStore = useToastStore()

const localSettings = ref<Settings>({
  theme: 'light',
  uploadQuality: 'high',
  autoOptimize: true,
  defaultStartPage: 1,
})

onMounted(async () => {
  await settingsStore.fetchSettings()
  localSettings.value = { ...settingsStore.settings }
})

const handleThemeChange = async () => {
  const success = await settingsStore.updateSettings({
    theme: localSettings.value.theme,
  })
  if (success) {
    toastStore.success('主题设置已更新')
  }
}

const handleUploadQualityChange = async () => {
  const success = await settingsStore.updateSettings({
    uploadQuality: localSettings.value.uploadQuality,
  })
  if (success) {
    toastStore.success('上传质量设置已更新')
  }
}

const handleAutoOptimizeChange = async () => {
  const success = await settingsStore.updateSettings({
    autoOptimize: localSettings.value.autoOptimize,
  })
  if (success) {
    toastStore.success('自动优化设置已更新')
  }
}

const handleStartPageChange = async () => {
  const success = await settingsStore.updateSettings({
    defaultStartPage: localSettings.value.defaultStartPage,
  })
  if (success) {
    toastStore.success('默认起始页码已更新')
  }
}
</script>

<style scoped>
.page-settings {
  max-width: 800px;
  margin: 0 auto;
  padding: var(--space-8);
}

.page-title {
  font-size: var(--text-4xl);
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--space-8);
  color: var(--text-primary);
}

.settings-sections {
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

.settings-section {
  background: var(--bg-elevated);
  border-radius: var(--radius-xl);
  padding: var(--space-6);
  box-shadow: var(--shadow-sm);
}

.section-title {
  font-size: var(--text-xl);
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-4);
  padding-bottom: var(--space-3);
  border-bottom: 1px solid var(--border-subtle);
  color: var(--text-primary);
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4) 0;
}

.setting-item:not(:last-child) {
  border-bottom: 1px solid var(--border-subtle);
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.setting-info label {
  font-weight: var(--font-weight-medium);
  color: var(--text-primary);
}

.setting-description {
  font-size: var(--text-sm);
  color: var(--text-secondary);
}



.setting-select,
.setting-input {
  padding: var(--space-2) var(--space-4);
  font-size: var(--text-base);
  color: var(--text-primary);
  background: var(--bg-primary);
  border: 1px solid var(--border-medium);
  border-radius: var(--radius-md);
  min-width: 150px;
  transition: all var(--duration-fast) var(--ease-out);
}

.setting-select:focus,
.setting-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: var(--shadow-focus);
}

.setting-input {
  width: 100px;
  text-align: center;
}

/* Toggle Switch */
.toggle-switch {
  position: relative;
  display: inline-block;
  width: 52px;
  height: 28px;
}

.toggle-switch input {
  opacity: 0;
  width: 0;
  height: 0;
}

.toggle-slider {
  position: absolute;
  cursor: pointer;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: var(--bg-secondary);
  border: 1px solid var(--border-medium);
  transition: all var(--duration-fast) var(--ease-out);
  border-radius: var(--radius-full);
}

.toggle-slider:before {
  position: absolute;
  content: '';
  height: 20px;
  width: 20px;
  left: 3px;
  bottom: 3px;
  background-color: white;
  transition: all var(--duration-fast) var(--ease-out);
  border-radius: var(--radius-full);
  box-shadow: var(--shadow-sm);
}

.toggle-switch input:checked + .toggle-slider {
  background-color: var(--color-accent);
  border-color: var(--color-accent);
}

.toggle-switch input:checked + .toggle-slider:before {
  transform: translateX(24px);
}

.toggle-switch:hover .toggle-slider {
  box-shadow: var(--shadow-md);
}

@media (max-width: 768px) {
  .page-settings {
    padding: var(--space-4);
  }

  .page-title {
    font-size: var(--text-3xl);
  }

  .setting-item {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-3);
  }

  .setting-select,
  .setting-input {
    width: 100%;
  }
}
</style>
