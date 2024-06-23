import {ref} from 'vue'
import {defineStore} from 'pinia'

interface GameInfo {
    name: string,
    url: string,
    img: string,
}

export const useGameStore = defineStore('counter', () => {
    const gameList = ref<GameInfo[]>()

    // 新增game
    function addGame(game: GameInfo) {
        gameList.value?.push(game)
    }

    // 删除
    function deleteGame(name: string) {
        gameList.value?.forEach((item, index) => {
            if (item.name === name) {
                gameList.value?.splice(index, 1)
            }
        })
    }

    return {gameList, addGame, deleteGame}
})
