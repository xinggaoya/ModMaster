<template>
  <div class="home">
    <n-modal v-model:show="showModel" title="搜索游戏" preset="card" style="width: 60%">
      <template v-slot:default>
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

      </template>
    </n-modal>
    <n-space vertical>
      <n-space justify="space-between">
        <h3>本地游戏</h3>
        <div>
          <n-button type="primary" @click="onShowModel">添加游戏</n-button>
          <n-button type="default" @click="getLocalGame">加载本地</n-button>
        </div>
      </n-space>
      <n-spin :show="loading" description="加载中...">
        <n-card v-for="(item,index) in localGame" :key="index" hoverable>
          <template #header>
            <div class="text-sm opacity-75">{{ item.name }}</div>
          </template>
          <template #header-extra>
            <n-space>
              <n-button type="default" @click="runTheGame(item.path)">启动</n-button>
              <n-popconfirm @positive-click="deleteGame(item.path)">
                <template v-slot:trigger>
                  <n-button type="error">删除</n-button>
                </template>
                请确认删除？
              </n-popconfirm>
            </n-space>
          </template>
        </n-card>
      </n-spin>

    </n-space>
  </div>
</template>

<script lang="ts" setup>
import {ref} from 'vue'
import {DeleteGame, GetGame, GetGameInfo, GetGameList, RunGame} from "~/go/internal/App";
import {useMessage} from "naive-ui";

const message = useMessage()
const gameName = ref('')
const showModel = ref(false)
const searchGameLoading = ref(false)
const loading = ref(false)
const gameList = ref<any>([])
const localGame = ref<any>([])

// 搜索游戏
async function onSearch() {
  searchGameLoading.value = true
  gameList.value = await GetGameList(gameName.value)
  searchGameLoading.value = false
}

// 打开下载弹出框
function onShowModel() {
  gameName.value = ''
  gameList.value = []
  showModel.value = true
}

// 获取本地游戏
async function getLocalGame() {
  loading.value = true
  localGame.value = await GetGame()
  loading.value = false
}

// 运行游戏
async function runTheGame(path: string) {
  RunGame(path).finally(() => {
    message.warning('执行完成')
  })
}

// 删除游戏
async function deleteGame(path: string) {
  DeleteGame(path).finally(() => {
    getLocalGame()
    message.success('删除成功')
  })
}

// 下载游戏
async function downloadGame(url: string, img: string) {
  GetGameInfo(url, img).finally(() => {
    getLocalGame()
    message.success('下载成功')
  })
}

getLocalGame()
</script>

<style scoped>

.home {
  padding: 10px;
}
</style>
