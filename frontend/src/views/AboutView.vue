<template>
  <div class="about">
    <div>
      <n-input v-model:value="gameName"/>
    </div>
    <n-space>
      <n-button @click="onClick">Default</n-button>
    </n-space>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue'
import {GetGameList} from "../../wailsjs/go/main/App";

const gameName = ref('')
const gameList = ref<any>([])

async function onClick() {
  gameList.value = await GetGameList(gameName.value)
}
</script>

<style>
@media (min-width: 1024px) {
  .about {
    min-height: 100vh;
    display: flex;
    align-items: center;
  }
}
</style>
