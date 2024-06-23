<script setup lang="ts">
import {RouterView} from 'vue-router'
import themeOverrides from '@/assets/theme-overrides.json'
</script>

<template>
  <n-config-provider :theme-overrides="themeOverrides">
    <n-modal-provider>
      <n-message-provider>
        <RouterView/>
      </n-message-provider>
    </n-modal-provider>
  </n-config-provider>
</template>

<style scoped>

</style>
