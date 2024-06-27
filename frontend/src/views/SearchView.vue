<template>
  <div>
    <div style="margin-bottom: 20px">
      <n-input-group>
        <n-input v-model:value="gameName" placeholder="输入游戏名称" @keydown.enter="onSearch"/>
        <n-button type="primary" @click="onSearch">搜 索</n-button>
      </n-input-group>
    </div>
    <div>
      <h3>搜索列表</h3>
    </div>
    <n-spin :show="searchGameLoading" description="加载中...">
      <n-scrollbar style="height: 400px">
        <n-space vertical>
          <n-card v-for="(item,index) in gameList" :key="index" header-style="padding: 0 20px">
            <template #header>
              <n-space>
                <n-image :src="item.img" width="60" height="60" fit="fit"/>
                <n-ellipsis style="max-width: 360px">{{ item.name }}</n-ellipsis>
              </n-space>
            </template>
            <template #header-extra>
              <n-button type="tertiary" @click="downloadGame(item.url,item.img)">下载</n-button>
            </template>
          </n-card>
        </n-space>
      </n-scrollbar>
    </n-spin>
  </div>
</template>
<script setup lang="ts">
import {ref} from "vue";
import {App} from "~/ModMaster/internal/service/index";
import {useMessage} from "naive-ui";
import type {GameInfo} from "@/stores/GameStore"
import {useGameStore} from "@/stores/GameStore";

const gameName = ref('')
const searchGameLoading = ref(false)
const gameList = ref<any>([])
const message = useMessage()
const gameStore = useGameStore()

// 搜索游戏
function onSearch() {
  searchGameLoading.value = true
  App.GetGameList(gameName.value).then((res: any) => {
    gameList.value = res
    searchGameLoading.value = false
  }).catch(() => {
    searchGameLoading.value = false
  })
}

// 下载游戏
function downloadGame(url: string, img: string) {
  App.GetGameInfo(url, img).then((info: GameInfo) => {
    gameStore.addGame(info)
    message.success('下载成功')
  }).catch(() => {
    message.error('下载失败')
  })
}

</script>

<style scoped>

</style>