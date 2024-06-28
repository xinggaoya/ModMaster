<template>
  <div class="home">
    <n-scrollbar style="height: 600px">
      <n-space>
        <n-card v-for="(item,index) in gameStore.gameList"
                content-style="padding: 10px"
                header-style="padding: 0 10px"
                :key="index" hoverable style="width: 200px;">
          <template v-slot:cover>
            <n-image :src="item.img" width="200" height="200"/>
          </template>
          <template v-slot:header>
            <n-ellipsis style="font-size: 16px;min-height:50px" line-clamp="2">{{ item.name }}</n-ellipsis>
          </template>
          <n-space justify="end">
            <n-button @click="runTheGame(item.path)" type="primary">
              <n-icon size="20">
                <GameControllerOutline/>
              </n-icon>
            </n-button>
            <n-popconfirm @positive-click="deleteGame(item.path)">
              <template v-slot:trigger>
                <n-button type="error">
                  <n-icon size="20">
                    <CloseCircleOutline/>
                  </n-icon>
                </n-button>
              </template>
              请确认删除？
            </n-popconfirm>
          </n-space>
        </n-card>
      </n-space>
    </n-scrollbar>
  </div>
</template>

<script lang="ts" setup>
import {DeleteGame, RunGame} from "~/go/internal/App";
import {useMessage} from "naive-ui";
import {useGameStore} from "@/stores/GameStore";
import {CloseCircleOutline, GameControllerOutline} from "@vicons/ionicons5";


const message = useMessage()
const gameStore = useGameStore()

// 运行游戏
function runTheGame(path: string) {
  RunGame(path).finally(() => {
    message.warning('执行完成')
  }).catch(() => {
    message.error('执行失败')
  })
}

// 删除游戏
function deleteGame(path: string) {
  DeleteGame(path).then(() => {
    gameStore.deleteGame(path)
    message.success('删除成功')
  }).catch(() => {
    message.error('删除失败')
  })
}


</script>

<style scoped>

.home {
  padding: 10px;
}
</style>
