<template>
  <AppLayout>
    <router-view />
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import AppLayout from '@/layouts/index.vue'
import { useSettingsStore } from '@/stores/settings'

const settingsStore = useSettingsStore()

onMounted(() => {
  // Apply saved settings (theme, etc.) on app start
  settingsStore.initSettings()
})
</script>

<style>
/* Global styles are imported in main.css */
</style>
