import {ref} from 'vue'
import {defineStore} from 'pinia'

export interface GameInfo {
    name: string,
    path: string,
    img: string,
}

export const useGameStore = defineStore('GameStore', () => {
    const gameList = ref<GameInfo[]>([])

    // 新增game
    function addGame(game: GameInfo) {
        gameList.value?.push(game)
    }

    // 删除
    function deleteGame(path: string) {
        gameList.value?.forEach((item, index) => {
            if (item.path === path) {
                gameList.value?.splice(index, 1)
            }
        })
    }

    return {gameList, addGame, deleteGame}
}, {
    persist: {
        storage: localStorage,
    }
})
