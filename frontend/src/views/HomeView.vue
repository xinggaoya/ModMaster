<template>
  <n-card>
    <n-tabs type="line" animated default-value="0">
      <n-tab-pane name="0" tab="我的修改器">
        <GameView/>
      </n-tab-pane>
      <n-tab-pane name="2" tab="修改器市场">
        <MarketView/>
      </n-tab-pane>
      <n-tab-pane name="1" tab="搜索修改器">
        <SearchView/>
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<script lang="ts" setup>
import GameView from "@/views/GameView.vue";
import MarketView from "@/views/MarketView.vue";
import SearchView from "@/views/SearchView.vue";
</script>

<style scoped>

</style>