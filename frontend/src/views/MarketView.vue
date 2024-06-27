<template>
  <n-scrollbar ref="marketScrollbar">
    <n-spin :show="loading" style="height: 600px" description="加载中...">
      <n-space>
        <n-card v-for="(item,index) in gameList"
                content-style="padding: 10px"
                header-style="padding: 0 10px"
                :key="index" hoverable style="width: 200px">
          <template v-slot:cover>
            <n-image :src="item.img" width="200" height="200"/>
          </template>
          <template v-slot:header>
            <n-ellipsis style="font-size: 16px;min-height:50px" line-clamp="2">{{ item.name }}</n-ellipsis>
          </template>
          <n-space justify="end">
            <n-button type="primary" @click="downloadGame(item.url, item.img)">
              <n-icon size="20">
                <DownloadOutline/>
              </n-icon>
            </n-button>
          </n-space>
        </n-card>
      </n-space>

      <n-space justify="center" style="padding: 20px 0">
        <n-button type="primary" @click="nextPage">加载更多</n-button>
      </n-space>
    </n-spin>
  </n-scrollbar>
</template>

<script setup lang="ts">
import {onMounted, ref} from "vue";
import {App} from "~/ModMaster/internal/service/index";
import {DownloadOutline} from "@vicons/ionicons5";
import {type GameInfo, useGameStore} from "@/stores/GameStore";
import {useMessage} from "naive-ui";

const gameList = ref<GameInfo[]>([])
const marketScrollbar = ref()
const loading = ref(false)
const message = useMessage()
const gameStore = useGameStore()
const page = ref(1)

// 获取分页列表
function getList() {
  loading.value = true
  App.GetGameListPage(page.value).then((data: any) => {
    gameList.value = data
    loading.value = false
    marketScrollbar.value.scrollTo({
      top: 0
    })
  })
}

// 下一页
function nextPage() {
  page.value++
  getList()
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


onMounted(() => {
  getList()
})

</script>

<style scoped>

</style>