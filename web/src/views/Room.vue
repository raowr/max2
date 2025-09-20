<template>
  <div id="room">
    <!-- Background -->
    <div class="background-container">
      <img src="@/assets/img/bg/indoor.png" alt="室内场景背景图">
    </div>

    <!-- Room Container -->
    <div class="room-container">
      <!-- Header with back button -->
      <div class="header">
        <div class="back-btn" @click="goHome">
          <img src="@/assets/img/ui/img_return1_bg.png" alt="返回按钮背景">
          <img src="@/assets/img/ui/img_return1.png" alt="返回" style="position: absolute; left: 0; top: 0;">
          <img src="@/assets/img/ui/txt_friendroom.png" alt="好友房间" class="room-title">
        </div>
      </div>

      <!-- Main Content -->
      <div class="main-content">
        <!-- Players Section -->
        <div class="players-section">
          <div class="player-card" v-for="(player, index) in state.players" :key="index">
            <!-- 房主标识 -->
            <img v-if="player.isOwner" src="@/assets/img/ui/roomowner.png" class="room-owner-badge" alt="房主标志">

            <div class="player-frame">
              <div class="player-bg-container">
                <img src="@/assets/img/ui/frame_gold.png" alt="玩家框" style="width: 100%; height: 100%;">
                <img src="@/assets/img/ui/character_bg.png" alt="玩家背景"
                  style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;">
                <img v-if="player.lihui !== ''" :src="player.lihui" alt="玩家lihui"
                  style="position: absolute; top: 0; left: 0; width: 100%; height: 100%;">
              </div>

              <!-- 玩家内容在玩家框内部 -->
              <div class="player-content">
                <div class="player-inner">
                  <!-- <img :src="player.avatar" class="player-avatar" alt="玩家头像"> -->

                  <!-- 名称容器 - 名称叠在背景上面 -->
                  <div class="player-name-container">
                    <img src="@/assets/img/ui/namebg.png" class="player-nameplate" alt="名称背景">
                    <div class="player-name">{{ player.name }}</div>
                  </div>

                  <div class="player-stars">
                    <img v-for="n in 3" :key="n"
                      :src="n <= player.stars ? 'src/assets/img/ui/star.png' : 'src/assets/img/ui/star_dark.png'"
                      alt="星星">
                  </div>

                  <img src="@/assets/img/ui/ready.png" class="player-ready" :class="{ active: player.ready }"
                    alt="准备状态">
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Controls and Info Section -->
        <div class="controls-section">
          <div class="game-info">
            <div class="info-item">
              <div class="info-label">游戏模式</div>
              <div class="info-value">{{ state.gameType }}</div>
            </div>

            <div class="info-item">
              <div class="info-label">奖励</div>
              <div class="info-value">{{ state.gameReward }}</div>
            </div>

            <div class="info-item">
              <div class="info-label">房间号</div>
              <div class="info-value">{{ state.roomId }}</div>
            </div>
          </div>

          <div class="room-actions">
            <div class="action-btn" @click="addComputer">
              添加电脑
            </div>
            <div class="action-btn" @click="customizeAvatar">
              调整装扮
            </div>
          </div>
        </div>
      </div>

      <!-- Ready Button -->
      <div class="ready-btn-container">
        <img v-if=!isReady() src="@/assets/img/ui/btn_start_no.png" class="ready-btn" alt="准备中" >
        <img v-else src="@/assets/img/ui/btn_start.png" class="ready-btn" alt="已准备" @click="toGame">
      </div>

      <!-- Connection Status Indicator -->
      <div class="connection-status" :class="state.connectionStatusClass">
        {{ state.connectionStatusText }}
      </div>
    </div>
  </div>

</template>
<script setup>
import { ref,reactive, onMounted, computed, onBeforeUnmount } from 'vue'
import { websocket } from '@/utils/websocket'
import { useRouter } from 'vue-router'
import { audioManager } from '@/utils/audio'
const router = useRouter()
const state = reactive({
  players: [
    {
      id: 1,
      name: '玩家一',
      avatar: 'src/assets/img/touxiang/bighead15339.png',
      stars: 3,
      ready: false,
      isOwner: true,
      lihui: 'src/assets/img/lihui/full15342.png',
    },
    {
      id: 2,
      name: '玩家二',
      avatar: 'src/assets/img/touxiang/bighead15419.png',
      stars: 2,
      ready: false,
      isOwner: false,
      lihui: 'src/assets/img/lihui/full15418.png',
    },
    {
      id: 3,
      name: '等待中...',
      avatar: '',
      stars: 0,
      ready: false,
      isOwner: false,
      lihui: '',
    },
    {
      id: 4,
      name: '等待中...',
      avatar: '',
      stars: 0,
      ready: false,
      isOwner: false,
      lihui: '',
    }
  ],
  gameType: '经典模式',
  gameReward: '100金币',
  roomId: '123456',

  // WebSocket related
  ws: null,
  isConnected: false,
  isConnecting: false,
  reconnectAttempts: 0,
  maxReconnectAttempts: 5,

  connectionStatusClass: computed(() => ({
    connected: websocket.socket?.readyState === WebSocket.OPEN,
    disconnected: websocket.socket?.readyState === WebSocket.CLOSED
  })),
  connectionStatusText: computed(() => {
    switch (websocket.socket?.readyState) {
      case WebSocket.OPEN: return '已连接'
      case WebSocket.CONNECTING: return '连接中...'
      default: return '已断开'
    }
  })
})
onMounted(() => {
  init()
})
const init = async () => {
  audioManager.preload('bgm', 'src/assets/music/game_bg.mp3')
   audioManager.playBGM('bgm')

  websocket.on('message', handleMessage)
  websocket.on('connected', () => console.log('Connected'))
  websocket.on('disconnected', () => console.log('Disconnected'))
  toggleReady()

}


const handleMessage = (data) => {
  // 处理接收到的消息
  console.log('Received message:', data)
    if (data.type === "initRoom") {
 try {
      // 解析JSON字符串
      const serverPlayers = JSON.parse(data.data)
      
      // 智能合并数据
      serverPlayers.forEach(serverPlayer => {
        const convertedId = serverPlayer.ID + 1
        // 查找已有玩家
        const existing = state.players.find(p => p.id === convertedId)
        
        if (existing) {
          // 更新已有玩家属性
          existing.name = serverPlayer.Name
          existing.isOwner = serverPlayer.Type === 0
          existing.avatar = getAvatarByType(serverPlayer.Type)
          existing.lihui = getLihuiByType(serverPlayer.Type)
          existing.ready=true
        } else {
          // 添加新玩家（如果需要）
          state.players.push({
            id: convertedId,
            name: serverPlayer.Name,
            avatar: getAvatarByType(serverPlayer.Type),
            stars: 3,
            ready: false,
            isOwner: serverPlayer.Type === 0,
            lihui: getLihuiByType(serverPlayer.Type)
          })
        }
      })
      
      // 保留不在服务端返回中的原有玩家（如等待中的空位）
      // 自动截断超过4个的玩家
      state.players.splice(4)
    } catch (e) {
      console.error('解析玩家数据失败:', e)
    }
  }
}
// 添加工具函数（放在script setup顶部）
const getAvatarByType = (type) => {
  const avatars = [
    'src/assets/img/touxiang/bighead15339.png', // 房主头像
    'src/assets/img/touxiang/bighead15419.png'  // 机器人默认头像
  ]
  return avatars[type] || 'src/assets/img/touxiang/bighead15729.png'
}

const getLihuiByType = (type) => {
  const lihuis = [
    'src/assets/img/lihui/full15342.png', // 房主立绘
    'src/assets/img/lihui/full15418.png'  // 机器人默认立绘
  ]
  return lihuis[type] || 'src/assets/img/lihui/full15703.png'
}
onBeforeUnmount(() => {
  websocket.off('message', handleMessage)

})

// 原有方法保持不变，修改发送消息的方法示例：
const toggleReady = () => {
  websocket.send({"type":"initRoom","data":"","name":""})
  console.log('send ready')
}

const isReady=() => {
  let ready = true
  state.players.forEach(serverPlayer => {
    if (!serverPlayer.ready) {
      ready = false
    }
  })
  return ready
}

const audioUrl = ref('src/assets/music/game_bg1.mp3') // 假设音频文件存放路径
const bgm = ref(null)

const toGame = () => {
  if (isReady()) {
    // websocket.send({"type":"toGame","data":"","name":""})
    //路由跳到游戏页面
    router.push({path:'/game'})
  }
}



</script>
<style scoped>
:root {
  --primary-color: #ffae4c;
  --secondary-color: #5d4037;
  --accent-gold: #ffd700;
  --text-light: #ffffff;
  --text-dark: #333333;
  --bg-overlay: rgba(0, 0, 0, 0.6);
}

/* 新增根容器样式 */
#room {
  position: relative;
  height: 100vh;
  overflow: hidden;
  --primary-color: #ffae4c;
  --secondary-color: #5d4037;
  --accent-gold: #ffd700;
  --text-light: #ffffff;
  --text-dark: #333333;
  --bg-overlay: rgba(0, 0, 0, 0.6);
  color: #ffffff;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  padding: 0;
  margin: 0;
  overflow: hidden;
  height: 100vh;
  font-family: 'Microsoft YaHei', sans-serif;
  color: var(--text-light);
}



/* 调整房间容器层级 */
.room-container {
  position: relative;
  z-index: 1;
  /* 确保内容在背景之上 */
}

.background-container {
  position: fixed;
  /* 改为fixed定位 */
  width: 100vw;
  height: 100vh;
  z-index: 0;
  /* 调整层级为非负值 */
  left: 0;
  top: 0;
}

.background-container img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.room-container {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px;
  position: relative;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  position: relative;
  z-index: 5;
}

.back-btn {
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s;
}

.back-btn:hover {
  transform: scale(1.05);
}

.back-btn img {
  height: 80px;
}

.room-title {
  position: absolute;
  left: 110px;
  top: 15px;
}

.main-content {
  display: flex;
  flex: 1;
  gap: 20px;
}

.players-section {
  flex: 3;
  display: flex;
  flex-wrap: wrap;
  gap: 2%;
  justify-content: flex-start;
}

.player-card {
  width: 23.5%;
  position: relative;
  margin-bottom: 20px;
}

.player-frame {
  width: 100%;
  position: relative;
  height: 80%;
  padding-top: 150%;
  /* 增加高度，从130%增加到150% */
}

.player-bg-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.player-content {
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 50%;
  /* 增加内容区域高度 */
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  padding-bottom: 12%;
  z-index: 5;
}

.player-inner {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.player-avatar {
  width: 45%;
  /* 稍微增大头像 */
  margin-bottom: 6%;
  z-index: 2;
}

/* 名称容器 - 新增 */
.player-name-container {
  position: relative;
  width: 90%;
  margin-bottom: 4%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.player-nameplate {
  width: 100%;
  z-index: 2;
}

.player-name {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 90%;
  text-align: center;
  font-size: 1.3vw;
  /* 稍微增大字体 */
  z-index: 3;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.8);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: bold;
  color: #fff;
}

.player-stars {
  display: flex;
  justify-content: center;
  width: 90%;
  z-index: 3;
}

.player-stars img {
  height: 18px;
  /* 增大星星大小 */
  margin: 0 3px;
}

.player-ready {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  z-index: 4;
  width: 45%;
  /* 增大准备图标 */
  display: none;
}

.player-ready.active {
  display: block;
  animation: pulse 1.5s infinite;
}

/* 房主标识位置 */
.room-owner-badge {
  position: absolute;
  top: -12px;
  left: -12px;
  width: 45%;
  /* 增大房主标识 */
  z-index: 10;
}

.controls-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 15px;
  min-width: 250px;
}

.game-info {
  /* background: var(--bg-overlay); */
  background: rgba(0, 0, 0, 0.6);
  border-radius: 10px;
  padding: 15px;
  backdrop-filter: blur(5px);
}

.info-item {
  margin-bottom: 15px;
  display: flex;
  flex-direction: column;
}

.info-label {
  font-size: 14px;
  color: var(--primary-color);
  --primary-color: #ffae4c;
  margin-bottom: 5px;
}

.info-value {
  font-size: 18px;
  font-weight: bold;
}

.room-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.action-btn {
  background: url('@/assets/img/ui/btn_small.png') center/cover no-repeat;
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
  color: var(--primary-color);
  --primary-color: #ffae4c;
  font-size: 25px;
}

.action-btn:hover {
  transform: scale(1.05);
  color: #ffc107;
}

.ready-btn-container {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
}

.ready-btn {
  cursor: pointer;
  transition: transform 0.2s;
}

.ready-btn:hover {
  transform: scale(1.1);
}

.connection-status {
  position: fixed;
  bottom: 10px;
  right: 10px;
  padding: 5px 10px;
  border-radius: 5px;
  font-size: 12px;
  z-index: 100;
}

.connected {
  background-color: rgba(76, 175, 80, 0.8);
}

.disconnected {
  background-color: rgba(244, 67, 54, 0.8);
}

.connecting {
  background-color: rgba(255, 193, 7, 0.8);
}

/* Animations */
@keyframes pulse {
  0% {
    opacity: 0.6;
    transform: translate(-50%, -50%) scale(0.95);
  }

  50% {
    opacity: 1;
    transform: translate(-50%, -50%) scale(1.05);
  }

  100% {
    opacity: 0.6;
    transform: translate(-50%, -50%) scale(0.95);
  }
}

@keyframes slideIn {
  from {
    transform: translateY(20px);
    opacity: 0;
  }

  to {
    transform: translateY(0);
    opacity: 1;
  }
}

.player-card {
  animation: slideIn 0.5s ease-out;
}

.player-card:nth-child(2) {
  animation-delay: 0.1s;
}

.player-card:nth-child(3) {
  animation-delay: 0.2s;
}

.player-card:nth-child(4) {
  animation-delay: 0.3s;
}

/* Responsive adjustments */
@media (max-width: 1200px) {
  .player-name {
    font-size: 1.5vw;
  }

  .player-frame {
    padding-top: 145%;
  }
}

@media (max-width: 900px) {
  .main-content {
    flex-direction: column;
  }

  .controls-section {
    flex-direction: row;
    flex-wrap: wrap;
  }

  .game-info {
    flex: 2;
  }

  .room-actions {
    flex: 1;
  }

  .player-name {
    font-size: 1.9vw;
  }

  .player-frame {
    padding-top: 140%;
  }
}

@media (max-width: 600px) {
  .players-section {
    gap: 10px;
  }

  .player-card {
    width: 48%;
  }

  .player-name {
    font-size: 2.6vw;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
  }

  .room-title {
    position: static;
    margin-top: 10px;
  }

  /* 移动端调整 */
  .room-owner-badge {
    width: 38%;
    top: -10px;
    left: -10px;
  }

  .player-avatar {
    width: 40%;
  }

  .player-content {
    height: 52%;
    padding-bottom: 10%;
  }

  .player-frame {
    padding-top: 135%;
  }

  .player-stars img {
    height: 16px;
  }
}
</style>
