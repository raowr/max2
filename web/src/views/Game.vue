<template>
  <div class="game">

    <div class="msg" style="position: absolute; top: 1%;left: 4%; color:gold;">
      {{ state.lastmsg }}
    </div>

    <!-- 移除内联样式，仅保留必要的动态属性 -->
    <div class="discard-pile-container">
      <img v-for="cardId in state.outCards" :key="cardId" :src="getCardImage(cardId)" class="discard-card">
    </div>
    <!--弃牌堆结束 -->

    <!-- 出牌动画过渡容器 -->
    <div class="card-transition-container">
      <transition-group name="card-move" tag="div">
        <img v-for="(cardId, index) in state.movingCards" :key="`moving-${cardId}`" :src="getCardImage(cardId)"
          class="moving-card" :style="getMovingCardStyle(index)">
      </transition-group>
    </div>

    <!--player1 牌 -->
    <!-- 移除容器的内联静态样式，保留类名 -->
    <div class="player1_card">
       <span  v-show="player1pass"  style="
        position: absolute;
        left: 50%;
        color: #FF6600;
        font-size: 28px;
        top: -50px;
        ">过</span>
      <!-- 为卡牌添加类名.player1-card-item，抽取静态样式 -->
      <img v-for="(cardId, index) in state.cards" :key="cardId" :src="getCardImage(cardId)" class="player1-card-item"
        :style="{
          // 保留动态样式：选中状态的位移、层级，以及基于cardId的 translateX
          transform: selectedCards.includes(cardId) ? 'translateY(-15px)' : 'translateX(calc(-1% * (13 - ' + index + ')))',
          zIndex: selectedCards.includes(cardId) ? 10 : 1
        }" @click="toggleCard(cardId)">
    </div>
    <!--player1 牌结束 -->

    <!--player1 信息 -->
    <div class="player1-container">
      <img src="@/assets/img/ui/chatlog.png" width="110px">

      <!-- 总瓜子数信息 (头像上方) -->
      <div style="position: absolute; top: -15%; left: 50%; transform: translateX(-50%); display: flex; align-items: center; 
              white-space: nowrap; min-width: 100px;"> <!-- 增大最小宽度至100px -->
        <img src="@/assets/img/ui/bf_heart.png" width='20px' style='margin-right: 5px;'>
        <p style="z-index:1;font-size:14px; color:yellow; font-weight: bold;">{{ state.player1Point }}瓜子</p>
      </div>

      <img src="@/assets/img/touxiang/smallhead.png" width="100px"
        style="position: absolute;bottom:10.2%;left:4%;border-radius: 25px;">
      <div style="width:110px;height:40px;text-align:center;position: absolute;">
        <p style="z-index:1;font-size:16px; color:white;">{{ state.player1Name }}</p>
      </div>
    </div>
    <!--player1 信息结束 -->

    <!--player2 信息 -->
    <div class="player2-container" style="">
      <!-- 修改为圆形倒计时 -->
      <div v-if="state.countdownPlayer == 2"
        style="position: absolute; top: -5%; left: 50%; transform: translateX(-50%); z-index: 999; text-align: center">
        <svg :width="100" :height="100">
          <circle cx="50" cy="50" r="45" stroke="#eee" stroke-width="8" fill="transparent" />
          <circle cx="50" cy="50" r="45" :stroke="countdownPlayer2 > (state.outCardTimeout / 3) ? '#4CAF50' : '#ff5722'"
            stroke-width="8" fill="transparent" :style="{
              strokeDasharray: 283,
              strokeDashoffset: 283 * (1 - countdownPlayer2 / state.outCardTimeout),
              transition: 'stroke-dashoffset 1s linear'
            }" />
        </svg>
        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); 
               font-size: 24px; color: white; text-shadow: 0 0 5px rgba(0,0,0,0.5)">
          {{ countdownPlayer2 }}
        </div>
      </div>
      <img src="@/assets/img/ui/chatlog.png" width="90px">
      <img :src="player2touxiang" width="90px" style="position: absolute;bottom:10.2%;left:1.1%;border-radius: 25px;">
      <div style="width:89px;height:40px;text-align:center;position: absolute;">
        <p style="z-index:1;font-size:16px; color:white;">{{ state.player2Name }}</p>
      </div>
      <span v-show="player2pass" style="
            position: absolute;
            right: 135%;
            color: #FF6600;
            font-size: 28px;
            top: -40%;
        ">过</span>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:-95%;top:-2%'>
      <p style="z-index:1;font-size:16px; color:white;position: absolute;top:96%;left:-75%;">剩<span style="color:#FF6600;">{{ state.player2CardsNum }} </span>张</p>
    </div>

    <!--player2 信息结束 -->

        <!--player3 信息 -->
    <div class="player3-container" style="">
      <!-- 修改为圆形倒计时 -->
      <div v-if="state.countdownPlayer == 3"
        style="position: absolute; top: -9%; left: 50%; transform: translateX(-50%); z-index: 999; text-align: center">
        <svg :width="100" :height="100">
          <circle cx="50" cy="50" r="45" stroke="#eee" stroke-width="8" fill="transparent" />
          <circle cx="50" cy="50" r="45" :stroke="countdownPlayer3 > (state.outCardTimeout / 3) ? '#4CAF50' : '#ff5722'"
            stroke-width="8" fill="transparent" :style="{
              strokeDasharray: 283,
              strokeDashoffset: 283 * (1 - countdownPlayer3 / state.outCardTimeout),
              transition: 'stroke-dashoffset 1s linear'
            }" />
        </svg>
        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); 
               font-size: 24px; color: white; text-shadow: 0 0 5px rgba(0,0,0,0.5)">
          {{ countdownPlayer3 }}
        </div>
      </div>
      <div width="90px" :style="{
        backgroundImage: `url(${chatlogBgUrl})`,  // 绑定解析后的路径
        backgroundSize: '90px',
        position: 'absolute',
        width: '90px',
        height: '90px'
      }">
        <img :src="player3touxiang" width="90px" style="bottom:31.2%;left:3.1%;border-radius: 25px;">
      </div>
        <span  v-show="player3pass"  style="
        position: absolute;
        left: 200%;
        color: #FF6600;
        font-size: 28px;
        top: 24px;
        ">过</span>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:100px;'>
      <p style="position: absolute;font-size:16px; color:white;top:90px;left:118px;width: 60px;">
        剩<span style="color:#FF6600;">{{ state.player3CardsNum }} </span>张</p>
      <div style="position: absolute;width:105px;height:40px;text-align:center;top:90px;">
        <p style="z-index:1;font-size:16px; color:white;">{{ state.player3Name }}</p>
      </div>

    </div>

    <!--player3 信息结束 -->

    <!--player4 信息 -->
    <div class="player4-container" style="">
      <!-- 修改为圆形倒计时 -->
      <div v-if="state.countdownPlayer == 4"
        style="position: absolute; top: -4%; left: 46%; transform: translateX(-50%); z-index: 999; text-align: center">
        <svg :width="100" :height="100">
          <circle cx="50" cy="50" r="45" stroke="#eee" stroke-width="8" fill="transparent" />
          <circle cx="50" cy="50" r="45" :stroke="countdownPlayer4 > (state.outCardTimeout / 3) ? '#4CAF50' : '#ff5722'"
            stroke-width="8" fill="transparent" :style="{
              strokeDasharray: 283,
              strokeDashoffset: 283 * (1 - countdownPlayer4 / state.outCardTimeout),
              transition: 'stroke-dashoffset 1s linear'
            }" />
        </svg>
        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); 
               font-size: 24px; color: white; text-shadow: 0 0 5px rgba(0,0,0,0.5)">
          {{ countdownPlayer4 }}
        </div>
      </div>
      <img src="@/assets/img/ui/chatlog.png" width="90px">
      <img :src="player4touxiang" width="85px" style="position: absolute;bottom:36.2%;left:3.1%;border-radius: 25px;">
      <div style="width:100px;height:40px;text-align:center;">
        <p style="z-index:1;font-size:16px; color:white;">{{ state.player4Name }}</p>
      </div>
      <span v-show="player4pass"  style="
            position: absolute;
            left: 120%;
            color: #FF6600;
            font-size: 28px;
            top: -28%;
        ">过</span>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:94%;top:0%'>
      <p style="z-index:1;font-size:16px; color:white;position: absolute;top:67%;left:112%;width: 60px;">
        剩<span style="color:#FF6600;">{{ state.player4CardsNum }} </span>张</p>
    </div>

    <!--player4 信息结束 -->



    <!--玩家身份信息 -->
    <div style="width:140px;height:40px;text-align:center;position: absolute;bottom:27%;left:3%;">
      <p style="z-index:1;font-size:20px; color:white;"></p>
    </div>
    <div style="width:140px;height:40px;text-align:center;position: absolute;top:14%;right:5%;">
      <p style="z-index:1;font-size:20px; color:white;"></p>
    </div>
    <div style="width:140px;height:40px;text-align:center;position: absolute;top:14%;left:3%;">
      <p style="z-index:1;font-size:20px; color:white;"></p>
    </div>
    <!--玩家身份信息结束 -->

    <!--功能区 -->


    <!-- 修改为圆形倒计时 -->
    <!-- 外层容器添加类名：countdown-circle-container -->
    <div v-if="state.countdownPlayer == 1" class="countdown-circle-container">
      <!-- SVG尺寸保持属性形式（非样式） -->
      <svg width="80" height="80">
        <!-- 背景圆添加类名：countdown-bg-circle -->
        <circle class="countdown-bg-circle" cx="40" cy="40" r="36" />
        <!-- 进度圆添加类名：countdown-progress-circle -->
        <circle class="countdown-progress-circle" cx="40" cy="40" r="36"
          :stroke="countdownPlayer1 > (state.outCardTimeout / 3) ? '#4CAF50' : '#ff5722'"
          :fill="countdownPlayer1 > (state.outCardTimeout / 3) ? '#81C784' : '#FF9800'" :style="{
            strokeDashoffset: 226.19 * (1 - countdownPlayer1 / state.outCardTimeout)  // 保留动态进度偏移
          }" />
      </svg>
      <!-- 文本区域添加类名：countdown-text -->
      <div class="countdown-text">
        {{ countdownPlayer1 }}
      </div>
    </div>



    <!-- 后手出牌时的功能 !-->
    <!-- 容器添加类名：backhand-actions-container -->
    <div v-if="state.countdownPlayer == 1" class="backhand-actions-container">

      <!-- 出牌按钮添加类名：chupai-btn -->
      <img v-if="checkOut()" src='@/assets/img/btn_chupai.png' class="chupai-btn" :style="{
        transform: isChupaiBtnPressed ? 'scale(0.9)' : 'scale(1)'  // 保留动态样式
      }" @click='chupai()' @mousedown="isChupaiBtnPressed = true" @mouseup="isChupaiBtnPressed = false"
        @mouseleave="isChupaiBtnPressed = false" @touchstart="isChupaiBtnPressed = true"
        @touchend="isChupaiBtnPressed = false">

      <!-- 禁用出牌按钮添加类名：chupai-btn-disabled -->
      <img v-else src='@/assets/img/btn_chupai_hui.png' class="chupai-btn-disabled">

      <!-- 不出按钮添加类名：buchu-btn -->
      <img v-if="state.mustPid != currentPlayerId" src='@/assets/img/btn_buchu.png' class="buchu-btn" :style="{
        transform: isPassBtnPressed ? 'scale(0.9)' : 'scale(1)'  // 保留动态样式
      }" @click='pass()' @mousedown="isPassBtnPressed = true" @mouseup="isPassBtnPressed = false"
        @mouseleave="isPassBtnPressed = false" @touchstart="isPassBtnPressed = true"
        @touchend="isPassBtnPressed = false">
    </div>
    <!-- 后手出牌时的功能 !-->

    <!-- 后手出牌时的功能 !-->


    <!-- 游戏结束弹窗 -->
    <div v-if="showGameOverModal" class="game-over-modal">
      <div class="modal-content">
        <!-- 胜利/失败标题 (移至背景图外部顶部) -->
        <h2 class="result-title">{{ isWinner ? '恭喜胜利!' : '游戏失败' }}</h2>

        <!-- 缩小的背景图容器 -->
        <div class="result-bg" :style="{
          backgroundImage: isWinner ? `url(${winBg})` : `url(${loseBg})`
        }">
          <!-- 移除原内部标题 -->
        </div>

        <!-- 积分信息 (背景图下方) -->
        <p class="score-info" :style="{
          // 胜利状态样式 (保持不变)
          fontSize: isWinner ? '24px' : '22px',
          fontWeight: isWinner ? 'bold' : 'bold',
          color: isWinner ? '#FF3333' : '#FF4500',
          textShadow: isWinner ? '0 0 10px rgba(255, 215, 0, 0.8)' : '0 0 8px rgba(255, 69, 0, 0.7)',
          padding: '10px',
          borderRadius: '8px',
          background: isWinner ? 'linear-gradient(45deg, #FFD700, #FFA500)' : 'linear-gradient(45deg, #FF6347, #FFA07A)'
        }">{{ gameOverMessage }}</p>

        <!-- 修改按钮区域，交换按钮顺序 只有房主才能操作-->
        <div class="buttons-container">
          <button class="room-btn" @click="toRoom()">
            返回房间
          </button>
          <button class="restart-btn" v-if="currentPlayerId === 0" @click="state.player1Point === 0 ? goToHome() : restartGame()">
            {{ state.player1Point === 0 ? '返回首页' : '再来一局' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, onBeforeUnmount, computed } from 'vue'
import { useRouter } from 'vue-router' // 添加这行导入
import { audioManager } from '@/utils/audio'
import { websocket } from '@/utils/websocket'
import { cardUtil, CARD_TYPE } from '@/utils/card';
import { storage } from '@/utils/storage'
import { getTouxiang } from '@/utils/touxiang'//随机返回一个头像
// 导入背景图片（新增代码）
import winBg from '@/assets/img/ui/win_bg.png'
import loseBg from '@/assets/img/ui/lose_bg.png'

import { reconnectWebSocket } from '@/utils/websocketReconnect';

const state = reactive({
  deck: [],
  countdownPlayer: 0,
  cards: [],
  player2CardsNum: 13,
  player3CardsNum: 13,
  player4CardsNum: 13,
  outCards: [],
  lastOutCards: [],//上一手出牌
  lastmsg: "",
  mustPid: 0,
  player1Point: 0, // 初始瓜子数
  outCardTimeout: 30, // 出牌最大时间(单位秒) /s
  containerAtEndPosition: false, // 控制容器是否在结束位置

  player1Name: "",
  player2Name: "",
  player3Name: "",
  player4Name: "",
})
const isPassBtnPressed = ref(false)  // 添加：跟踪按钮按下状态
const isChupaiBtnPressed = ref(false)  // 添加：跟踪"出牌"按钮按下状态

// 添加游戏结束弹窗相关状态
const showGameOverModal = ref(false)
const isWinner = ref(false)
const gameOverMessage = ref('')
const router = useRouter() // 添加这行获取router实例

const selectedCards = ref([])  // 改为数组存储选中状态

// 新增倒计时逻辑
const countdownPlayer1 = ref(state.outCardTimeout) // 初始30秒
const countdownPlayer2 = ref(state.outCardTimeout) // 初始30秒
const countdownPlayer3 = ref(state.outCardTimeout) // 初始30秒
const countdownPlayer4 = ref(state.outCardTimeout) // 初始30秒
let timer1 = null
let timer2 = null
let timer3 = null
let timer4 = null

// 在原有状态变量中添加
const isPlayingCard = ref(false) // 标记是否正在出牌过程中
const pendingCards = ref([]) // 存储等待验证的牌
const playCardTimer = ref(null) // 出牌超时定时器

// 响应式变量：当前屏幕宽度
const screenWidth = ref(window.innerWidth);
// 缩放比例（默认100%，小于998px时为60%）
const scaleFactor = ref(screenWidth.value < 998 ? 0.6 : 1);

// 添加这行来定义bgmInitialized变量
const bgmInitialized = ref(false) // 跟踪BGM是否已经初始化

const player2touxiang = ref("")
const player3touxiang = ref("")
const player4touxiang = ref("")

const player1pass = ref(false)
const player2pass = ref(false)
const player3pass = ref(false)
const player4pass = ref(false)

const currentPlayerUid = ref(storage.local.get('user_id'))
const currentPlayerId = ref(0)
const proactivePass = ref(false) // 主动过牌标记

const roomType = ref(0)


// 监听屏幕尺寸变化
const handleResize = () => {
  screenWidth.value = window.innerWidth;
  scaleFactor.value = screenWidth.value < 998 ? 0.6 : 1; // 更新缩放比例
};

const fetchGameDataFromWebSocket = () => {
  const sendGetInfo = (retryCount = 0) => {
    console.log('game sending getInfo (retry:', retryCount, ')');
    
    // 检查连接状态
    if (!websocket.ws || websocket.ws.readyState !== WebSocket.OPEN) {
      console.log('game WebSocket not open, waiting...');
      // 如果还没到最大重试次数，继续重试
      if (retryCount < 5) {
        setTimeout(() => sendGetInfo(retryCount + 1), 100);
      }
      return;
    }
    
    websocket.send({ "type": "getInfo", "data": "", "name": "" });
    
    // 设置超时检查，如果500ms内没收到响应就重试
    const timeout = setTimeout(() => {
      if (retryCount < 5) {
        console.log('game getInfo 超时，重试 (retry:', retryCount + 1, ')');
        sendGetInfo(retryCount + 1);
      }
    }, 500);
    
    // 在收到 getInfo 响应时清除超时
    const handleMessageOnce = (data) => {
      try {
        const parsedData = typeof data === 'string' ? JSON.parse(data) : data;
        if (parsedData.type === "getInfo") {
          clearTimeout(timeout);
          websocket.off('message', handleMessageOnce);
          console.log('game getInfo 响应已收到，取消重试');
        }
      } catch (e) {}
    };
    websocket.on('message', handleMessageOnce);
  };

  // 监听 open 事件（处理未来的连接打开）
  websocket.on('open', sendGetInfo);

  // 立即检查当前连接状态
  if (websocket.ws) {
    console.log('game WebSocket 状态:', websocket.ws.readyState, '(0=CONNECTING, 1=OPEN, 2=CLOSING, 3=CLOSED)');
    if (websocket.ws.readyState === WebSocket.OPEN) {
      // 已连接，立即发送
      sendGetInfo();
    } else if (websocket.ws.readyState === WebSocket.CONNECTING) {
      // 正在连接中，等待 open 事件触发 sendGetInfo
      console.log('game WebSocket 正在连接中，等待 open 事件');
    } else {
      // 未连接或已关闭，等待 reconnectWebSocket 重新连接
      console.log('game WebSocket 未连接，等待重连');
    }
  }
}

onMounted(() => {

  initDeck()
   // 页面刷新时自动重连
  reconnectWebSocket(websocket);
  
  const bgmUrl = new URL('@/assets/music/game_bg1.mp3', import.meta.url).href;
  audioManager.preload('bgm', bgmUrl); // 使用解析后的 URL
  // 添加用户交互事件监听器来初始化BGM
  const initBGMOnInteraction = () => {
    if (!bgmInitialized.value) {
      audioManager.playBGM('bgm')
      bgmInitialized.value = true
      // 移除监听器，避免重复触发
      document.removeEventListener('click', initBGMOnInteraction)
      document.removeEventListener('touchstart', initBGMOnInteraction)
    }
  }

  // 添加多种交互方式的支持
  document.addEventListener('click', initBGMOnInteraction)
  document.addEventListener('touchstart', initBGMOnInteraction)

  // 1. 尝试从 localStorage 读取缓存的游戏数据
  const cachedData = storage.local.get('roomInfo');
  if (cachedData) {
    try {
      const gameData = JSON.parse(cachedData);
      console.log('使用缓存数据初始化游戏');
      initGameWithData(gameData);
      // 清除缓存，避免下次进入时重复使用（可选）
      storage.local.remove('roomInfo');
      // 不需要再发送 getInfo
    } catch (e) {
      console.error('解析缓存数据失败', e);
      fetchGameDataFromWebSocket();
    }
  } else {
    fetchGameDataFromWebSocket();
  }

  websocket.on('message', handleMessage)
  window.addEventListener('resize', handleResize)

  // 从storage中获取player2avatar
  player2touxiang.value = storage.local.get('player2avatar')
  player3touxiang.value = storage.local.get('player3avatar')
  player4touxiang.value = storage.local.get('player4avatar')

})

// 组件卸载时移除回调（关键：防止内存泄漏）
onUnmounted(() => {
  websocket.off('open', () => { console.log('off open'); });
  websocket.off('message', () => { console.log('off message'); });
  websocket.off('error', () => { console.log('off error'); });
  websocket.off('error', () => { console.log('off error'); });

  // 移除BGM初始化相关的事件监听器
  document.removeEventListener('click', () => { });
  document.removeEventListener('touchstart', () => { });
  window.removeEventListener('resize', handleResize);
});

onUnmounted(() => window.removeEventListener('resize', handleResize));

// 初始化一副牌
const initDeck = () => {
  const suits = ["♦", "♣", "♥", "♠"];
  const values = ["3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"];
  let cardId = 0;
  for (let i = 0; i < values.length; i++) {
    const valueName = values[i];
    for (let suit = 0; suit < suits.length; suit++) {
      const suitName = suits[suit];
      cardId++;
      const card = {
        rank: 3 + i,
        suit: suit,
        name: suitName + valueName,
        id: cardId,
        rankName: valueName,
        suitName: suitName
      };
      state.deck.push(card);
    }
  }
  // console.log(state.deck);
}

const initGameWithData =(data) =>{
    roomType.value = data.type
    //确定玩家pid
    for (let i = 0; i < data.players.length; i++) {
      if (data.players[i].UserId == currentPlayerUid.value){
        currentPlayerId.value = data.players[i].ID
      }
    }
    for (let i = 0; i < data.players.length; i++) {
      //设置玩家名字
      switch (getPlayerPosition(data.players[i].ID) + 1) {
        case 1:
          state.player1Name = data.players[i].Name
          break;
        case 2:
          state.player2Name = data.players[i].Name
          break;
        case 3:
          state.player3Name = data.players[i].Name
          break;
        case 4:
          state.player4Name = data.players[i].Name
          break;
      }
    }
    //对应位置
    //先确定差值1-0
    // 假如：玩家id：1，玩家就是确定在底部0位置，2号玩家：4-|(2-1)|,读作：4减去2-1的绝对值
    //当减去差值大于等于0，直接用，当减去差值小于0需要用4减去差值的绝对值

    //是否游戏中
    if (data.isPlaying) {
      state.outCards = (data.outCards)
      state.lastOutCards = (data.outCards)
      //如果必出是玩家，记录下必出玩家的pid
      state.mustPid = data.mustPid
      //玩家牌排序
      state.cards = data.cards.sort((a, b) => b - a);
      //设置牌数
      for (let i = 0; i < data.cardsNum.length; i++) {
        switch (getPlayerPosition(data.cardsNum[i].id) + 1) {
          case 2:
            state.player2CardsNum = data.cardsNum[i].cardNum
            break;
          case 3:
            state.player3CardsNum = data.cardsNum[i].cardNum
            break;
          case 4:
            state.player4CardsNum = data.cardsNum[i].cardNum
            break;
        }
      }
      // 更新倒计时时长（如果服务端提供）
      if (data.outCardTimeout !== undefined) {
        state.outCardTimeout = data.outCardTimeout
      }
      //先出牌的开始倒计时
      startCountdown(getPlayerPosition(data.current) + 1, data.remainOutCardTimeout)
      //如果必出是玩家，记录下必出玩家的pid
      //更新玩家总瓜子数
      state.player1Point = data.playerPoint
      let cardsMsg = ""
      for (let i = 0; i < data.outCards.length; i++) {
        for (let j = 0; j < state.deck.length; j++) {
          if (state.deck[j].id == data.outCards[i]) {
            cardsMsg += state.deck[j].name + " "
          }
        }
      }
      switch (getPlayerPosition(data.lastPid)) {
        case 0:
          state.lastmsg = "座位1 "+state.player1Name+"出了：" + cardsMsg
          break;
        case 1:
          state.lastmsg = "座位2 "+state.player2Name+"出了：" + cardsMsg
          break;
        case 2:
          state.lastmsg = "座位3 "+state.player3Name+"出了：" + cardsMsg
          break;
        case 3:
          state.lastmsg = "座位4 "+state.player4Name+"出了：" + cardsMsg
          break;
      }
    } else {
      switch (data.status) {
          //   case 0://未开始
          //     websocket.send({ "type": "play", "data": "", "name": "" })
          //     break;
        case 2://结算中
          console.log("结算中");
          router.push('/index')  // 跳转到首页路由
          break;
      }

    }
}

const handleMessage = (data) => {
  // 处理接收到的消息
  console.log('Received message:', data)
  const parsedData = JSON.parse(data)
  if (parsedData.type == "getInfo") {
    data = JSON.parse(parsedData.data)
    console.log("getInfo", data);
    initGameWithData(data)
  }
  if (parsedData.type == "showCard") {

    resetAvatar()

    data = JSON.parse(parsedData.data)
    console.log(data);
    //玩家牌排序
    state.cards = data.cards.sort((a, b) => b - a);
    //重置牌数
    state.player2CardsNum = 13
    state.player3CardsNum = 13
    state.player4CardsNum = 13
    // 更新倒计时时长（如果服务端提供）
    if (data.outCardTimeout !== undefined) {
      state.outCardTimeout = data.outCardTimeout
    }
    //先出牌的开始倒计时
    startCountdown(getPlayerPosition(data.current) + 1, state.outCardTimeout)
    //如果必出是玩家，记录下必出玩家的pid
    state.mustPid = data.current
    //更新玩家总瓜子数
    state.player1Point = data.playerPoint
    // 重置倒计时
    countdownPlayer1.value = state.outCardTimeout
    //清空选中的牌
    selectedCards.value = []

    player1pass.value = false
    player2pass.value = false
    player3pass.value = false
    player4pass.value = false
  }
  if (parsedData.type == "outCard") {
    data = JSON.parse(parsedData.data)
    // 刷新玩家牌数
    let cardsMsg = ""
    switch (getPlayerPosition(data.pid)) {
      case 0:
        if (isPlayingCard.value) {
          isPlayingCard.value = false
          // 无论成功失败都清除定时器
          console.log('清除出牌定时器')
          clearTimeout(playCardTimer.value)
          //等待动画完成后清空移动中的牌
          state.movingCards = []
        }
        if (data.code == 0) {
          // state.must = data.must
          // selectedCards.value = []
          // state.outCards=(data.cards)
          // startCountdown(2)
          for (let i = 0; i < data.cardIds.length; i++) {
            for (let j = 0; j < state.deck.length; j++) {
              if (state.deck[j].id == data.cardIds[i]) {
                cardsMsg += state.deck[j].name + " "
              }
            }
          }
          state.lastmsg = "座位1 "+state.player1Name+"出了：" + cardsMsg

          // 过滤掉state.cards中存在于data.cards的牌
          state.cards = state.cards.filter(item =>
            !data.cardIds.some(cardId => cardId === item)  // 取反some的结果，排除包含的牌
          );
          selectedCards.value = []
          // 出牌成功，清空 pending
          pendingCards.value = []
        //隐藏过
        player1pass.value = false
        } else {
          // 出牌失败，恢复牌
          console.log('出牌失败，恢复牌')
          state.cards = [...state.cards, ...pendingCards.value]
          state.cards.sort((a, b) => b - a) // 保持排序一致
          //恢复上一手出牌
          state.outCards = state.lastOutCards
        }
        break
      case 1:
        state.player2CardsNum = data.cards_num
        // state.outCards=(data.cards)
        // startCountdown(3)
        for (let i = 0; i < data.cardIds.length; i++) {
          for (let j = 0; j < state.deck.length; j++) {
            if (state.deck[j].id == data.cardIds[i]) {
              cardsMsg += state.deck[j].name + " "
            }
          }
        }
        state.lastmsg = "座位2 "+state.player2Name+"出了：" + cardsMsg
                //隐藏过
        player2pass.value = false
        break
      case 2:
        state.player3CardsNum = data.cards_num
        // state.outCards=(data.cards)
        // startCountdown(4)
        for (let i = 0; i < data.cardIds.length; i++) {
          for (let j = 0; j < state.deck.length; j++) {
            if (state.deck[j].id == data.cardIds[i]) {
              cardsMsg += state.deck[j].name + " "
            }
          }
        }
        state.lastmsg = "座位3 "+state.player3Name+"出了：" + cardsMsg
        //隐藏过
        player3pass.value = false
        break
      case 3:
        state.player4CardsNum = data.cards_num
        // state.outCards=(data.cards)
        // startCountdown(1)
        for (let i = 0; i < data.cardIds.length; i++) {
          for (let j = 0; j < state.deck.length; j++) {
            if (state.deck[j].id == data.cardIds[i]) {
              cardsMsg += state.deck[j].name + " "
            }
          }
        }
        state.lastmsg = "座位4 "+state.player4Name+"出了：" + cardsMsg
        //隐藏过
        player4pass.value = false
        break
    }

    //如果成功再改内容
    if (data.code == 0) {
      state.outCards = (data.cardIds)
      state.lastOutCards = (data.cardIds)
      // 更新倒计时时长（如果服务端提供）
      if (data.outCardTimeout !== undefined) {
        state.outCardTimeout = data.outCardTimeout
      }
      startCountdown(getPlayerPosition(data.current) + 1, state.outCardTimeout)
      //如果必出是玩家，记录下必出玩家的pid
      state.mustPid = data.mustPid

      //播放出牌声音
      const cardType = cardUtil.getCardType(state.outCards.map(id => state.deck.find(c => c.id === id)).filter(Boolean))
      const playedCards = state.outCards.map(id => state.deck.find(c => c.id === id)).filter(Boolean)
      console.log(state.outCards)
      console.log(cardType)
      console.log(playedCards)
      let musicPath = ""
      switch (cardType) {
        case CARD_TYPE.SINGLE:
          musicPath = "single/" + playedCards[0].id
          playSound(musicPath)
          break
        case CARD_TYPE.PAIR:
          musicPath = "pair/" + playedCards[0].rank
          playSound(musicPath)
          break
        case CARD_TYPE.STRAIGHT:
          musicPath = "straight/straight"
          playSound(musicPath)
          break
        case CARD_TYPE.SUIT:
          musicPath = "suit/suit"
          playSound(musicPath)
          break
        case CARD_TYPE.FULL_HOUSE:
          let countFull = cardUtil.countRanks(playedCards)
          let rankFull = cardUtil.getRankByCount(countFull, 3)
          musicPath = "full_house/" + rankFull
          playSound(musicPath)
          break
        case CARD_TYPE.FOUR_OF_A_KIND:
          let countFour = cardUtil.countRanks(playedCards)
          let rankFour= cardUtil.getRankByCount(countFour, 4)
          musicPath = "four/" + rankFour
          playSound(musicPath)
          break
        case CARD_TYPE.STRAIGHT_FLUSH:
          musicPath = "straight_flush/straight_flush"
          playSound(musicPath)
          break
        default:
          playSound('default')
          break
      }
          //先隐藏过
     if (getPlayerPosition(data.current) == 0) {
        player1pass.value = false
      } else if (getPlayerPosition(data.current) == 1) {
        player2pass.value = false
      } else if (getPlayerPosition(data.current) == 2) {
        player3pass.value = false
      } else if (getPlayerPosition(data.current) == 3) {
        player4pass.value = false
      }
    }

  }
  if (parsedData.type == "pass") {
    data = JSON.parse(parsedData.data)
    state.mustPid = data.mustPid

    if (!proactivePass.value) {
      const musicPath = "guo"
      playSound(musicPath)
      startCountdown(getPlayerPosition(data.current) + 1, state.outCardTimeout)
    }
    
    if (data.pid == currentPlayerId.value ){
      proactivePass.value = false
    }


    //判断谁过
     if (getPlayerPosition(data.pid) == 0) {
      player1pass.value = true
    } else if (getPlayerPosition(data.pid) == 1) {
      player2pass.value = true
    } else if (getPlayerPosition(data.pid) == 2) {
      player3pass.value = true
    } else if (getPlayerPosition(data.pid) == 3) {
      player4pass.value = true
    }

              //先隐藏过
    if (getPlayerPosition(data.current) == 0) {
      player1pass.value = false
    } else if (getPlayerPosition(data.current) == 1) {
      player2pass.value = false
    } else if (getPlayerPosition(data.current) == 2) {
      player3pass.value = false
    } else if (getPlayerPosition(data.current) == 3) {
      player4pass.value = false
    }

  }
  if (parsedData.type == "over") {
    clearInterval(timer1)
    clearInterval(timer2)
    clearInterval(timer3)
    clearInterval(timer4)
    state.countdownPlayer = 0
    data = JSON.parse(parsedData.data)
    data.forEach(winer=>{
        if (winer.win > 0 ){
          state.lastmsg = "游戏结束，玩家: " + winer.winName + "胜利,赢得:" + winer.win + "颗瓜子"
        }
      if (winer.winner == currentPlayerId.value ){
        var win = Math.abs(winer.win)
        isWinner.value = winer.win > 0
        gameOverMessage.value = isWinner.value
                ? `哇！恭喜你赢得了${win}颗瓜子!`
                : `咦！你很菜,这局输了${win}颗瓜子`

        //更新总瓜子数
        state.player1Point = winer.point
      }
    })

    // 显示游戏结束弹窗
    showGameOverModal.value = true
    // 判断当前玩家是否胜利（假设当前玩家是"帅哥1"）
    // isWinner.value = data.winner === 0
    // gameOverMessage.value = isWinner.value
    //   ? `哇！恭喜你赢得了${data.playerWin}颗瓜子!`
    //   : `咦！你很菜,这局输了${data.playerWin}颗瓜子`

    // 重置弃牌
    state.outCards = []
    // 移除自动重新开始的逻辑，由用户点击按钮触发
    if (isPlayingCard.value) {
      isPlayingCard.value = false
      pendingCards.value = []
      // 游戏结束也清除定时器
      clearTimeout(playCardTimer.value)
    }
    
    player1pass.value = false
    player2pass.value = false
    player3pass.value = false
    player4pass.value = false
    
    //玩家输完，跳到首页
    if (state.player1Point <= 0) {
      router.push({ path: '/index' })
    }

    //好友房，五秒后跳到房间页面
    // if (roomType.value == 2) {
    //   setTimeout(()=>{
    //     router.push({path:'/room'})
    //   },5000)
    // }

    //修改玩家积分
    var userInfo = storage.local.get('user')
    userInfo.point = state.player1Point
    if (roomType.value == 1) {
      storage.local.set('user', userInfo)
    }
    
  }
  if (parsedData.type === "play") {
    showGameOverModal.value = false
  }
}

// 在文件顶部添加音频预加载代码
// 预加载 music 目录下所有 .mp3 音频文件（Vite 特有的资源导入方式）
// eager: true 表示立即加载，import: 'default' 获取音频 URL
const audioFiles = import.meta.glob('@/assets/music/**/*.mp3', { eager: true, import: 'default' });
console.log(audioFiles)

// 动态获取音频文件路径（利用预加载的音频映射）
const getAudioUrl = (musicPath) => {
  try {
    // 构造预加载音频的路径 key（与 glob 匹配的路径格式）
    const audioKey = `/src/assets/music/${musicPath}.mp3`;
    // 返回预加载的音频 URL（开发/生产环境路径自动适配）
    return audioFiles[audioKey] || ''; // 若找不到对应音频，返回空
  } catch (error) {
    console.error('获取音频URL失败:', error);
    return '';
  }
};

const playSound = (musicPath) => {
  try {
    // 直接使用 AudioManager 的 playOnce 方法播放音频
    audioManager.playOnce(musicPath);
    
    // 备选方案：如果 AudioManager 中没有预加载该音频，则尝试直接加载播放
    if (!audioManager.audioElements.has(musicPath)) {
      // 获取音频 URL
      const audioUrl = getAudioUrl(musicPath);
      if (audioUrl) {
        // 创建新的 Audio 对象并播放
        const audio = new Audio(audioUrl);
        audio.play().catch(err => {
          console.warn('播放音频失败:', err);
        });
      }
    }
  } catch (error) {
    console.error('播放声音时发生错误:', error);
  }
};

// 添加重新开始游戏的方法
const restartGame = () => {
  // showGameOverModal.value = false
  websocket.send({ "type": "play", "data": "", "name": "" })
}



const toggleCard = (n) => {
  const index = selectedCards.value.indexOf(n)
  if (index > -1) {
    // 如果已选中则移除
    selectedCards.value.splice(index, 1)
  } else {
    // 添加新选中项（保留旧项）
    selectedCards.value.push(n)
  }
}

// 修改checkOut方法，正确处理currentCards为选中牌的数组
const checkOut = () => {
  if (selectedCards.value.length === 0) {
    return false;
  }

  // 构建currentCards数组（选中的牌对象数组）
  const currentCards = selectedCards.value.map(cardId => {
    return state.deck.find(card => card.id === cardId);
  }).filter(Boolean); // 过滤无效牌

  // 构建lastHand参数
  const lastHand = {
    cards: state.outCards.map(cardId => {
      return state.deck.find(card => card.id === cardId);
    }).filter(Boolean),
    isSelf: state.mustPid === currentPlayerId.value, // 是否是自己出的最后一手牌
    type: state.outCards.length > 0
      ? cardUtil.getCardType(state.outCards.map(id => state.deck.find(c => c.id === id)).filter(Boolean))
      : null
  };
  // 调用牌型校验工具
  let res = cardUtil.canPlayCards(currentCards, lastHand);
  console.log("验证出牌：", currentCards, lastHand, res)
  return res
};




const startCountdown = (pid, remainOutCardTimeout) => {
  clearInterval(timer1)
  clearInterval(timer2)
  clearInterval(timer3)
  clearInterval(timer4)
  state.countdownPlayer = pid
  switch (pid) {
    case 1:
      // 关闭其他倒计时
      clearInterval(timer2)
      clearInterval(timer3)
      clearInterval(timer4)
      //显示并开始，当前倒计时
      countdownPlayer1.value = remainOutCardTimeout // 使用从服务端获取的倒计时时长
      let isPlayKuaidian = false
      timer1 = setInterval(() => {
        if (countdownPlayer1.value > 0) {
          countdownPlayer1.value--
          if (countdownPlayer1.value < remainOutCardTimeout/2) {
            if (!isPlayKuaidian) {
              isPlayKuaidian = true
              playSound('kuaidian')
            }
          }
        } else {
          clearInterval(timer1)
          // 触发超时逻辑
          console.log('出牌超时')
          state.countdownPlayer = 0
          // startCountdown(2)
        }
      }, 1000)

      break
    case 2:
      //关闭其他倒计时
      clearInterval(timer1)
      clearInterval(timer3)
      clearInterval(timer4)
      countdownPlayer2.value = remainOutCardTimeout // 使用从服务端获取的倒计时时长
      timer2 = setInterval(() => {
        if (countdownPlayer2.value > 0) {
          countdownPlayer2.value--
        } else {
          clearInterval(timer2)
          // 触发超时逻辑
          console.log('出牌超时')
          state.countdownPlayer = 0
          // startCountdown(3)
        }
      }, 1000)
      break
    case 3:
      // 关闭其他倒计时
      clearInterval(timer1)
      clearInterval(timer2)
      clearInterval(timer4)
      countdownPlayer3.value = remainOutCardTimeout // 使用从服务端获取的倒计时时长
      timer3 = setInterval(() => {
        if (countdownPlayer3.value > 0) {
          countdownPlayer3.value--
        } else {
          clearInterval(timer3)
          // 触发超时逻辑
          console.log('出牌超时')
          state.countdownPlayer = 0
          // startCountdown(4)
        }
      }, 1000)
      break
    case 4:
      // 关闭其他倒计时
      clearInterval(timer1)
      clearInterval(timer2)
      clearInterval(timer3)
      countdownPlayer4.value = remainOutCardTimeout // 使用从服务端获取的倒计时时长
      timer4 = setInterval(() => {
        if (countdownPlayer4.value > 0) {
          countdownPlayer4.value--
        } else {
          clearInterval(timer4)
          // 触发超时逻辑
          console.log('出牌超时')
          state.countdownPlayer = 0
          // startCountdown(1)
        }
      }, 1000)
      break
  }
}

const resetCountdown = () => {
  countdown.value = 30
  if (!timer) startCountdown()
}


onBeforeUnmount(() => {
  clearInterval(timer)
})


// 修改 chupai 函数，添加牌移除和定时器逻辑
const chupai = () => {
  if (isPlayingCard.value) {
    console.log("正在处理出牌，请等待响应")
    return false
  }

  if (selectedCards.value.length === 0) {
    return false
  } else {
    // 1. 保存选中的牌用于后续恢复
    const cardsToPlay = [...selectedCards.value]
    pendingCards.value = cardsToPlay
    isPlayingCard.value = true

    // 2. 先从手牌中移除牌
    const originalCards = [...state.cards]
    state.cards = state.cards.filter(card => !cardsToPlay.includes(card))

    // 3. 准备并发送数据到服务端
    let data = {
      type: "playCard",
      data: JSON.stringify({
        pid: state.countdownPlayer - 1,
        cardIds: cardsToPlay,
      }),
    }
    websocket.send(data)

    // 4. 记录正在移动的牌
    state.movingCards = [...cardsToPlay]

    // 5. 清空选中状态
    selectedCards.value = []

    // 6. 等待动画完成后清空移动中的牌
    setTimeout(() => {
      state.movingCards = []
      state.outCards = [...cardsToPlay]
    }, 500)

    // 7. 设置30秒超时定时器，超时未响应则恢复牌
    playCardTimer.value = setTimeout(() => {
      console.log('出牌超时，未收到服务响应，恢复牌')
      if (isPlayingCard.value) {
        // 恢复牌
        state.cards = [...originalCards]
        state.cards.sort((a, b) => b - a) // 保持排序一致
        // state.lastmsg = "出牌超时，服务未响应"
      }
      isPlayingCard.value = false
      pendingCards.value = []
      //等待动画完成后清空移动中的牌
      state.movingCards = []
    }, 30000) // 30秒超时
  }
}

// 修改getMovingCardStyle函数，使移动中的牌排成一行并重叠4px
const getMovingCardStyle = (index) => {
  // 获取当前移动中的卡牌总数
  const totalCards = state.movingCards.length;

  // 设置卡牌宽度为120px，重叠4px
  const cardWidth = 120 * scaleFactor.value;
  const cardHeight = 140 * scaleFactor.value;
  const overlap = 10;
  const effectiveWidth = cardWidth - overlap; // 有效宽度为116px，确保重叠4px

  // 计算整行卡牌的总宽度（包括重叠）
  const totalWidth = cardWidth + (totalCards - 1) * effectiveWidth;

  // 计算起始偏移量，使整行卡牌居中显示
  const startOffset = -totalWidth / 2;

  // 计算当前卡牌的水平偏移量
  const horizontalOffset = startOffset + index * effectiveWidth;

  // 保持卡牌水平，不添加旋转角度
  const rotation = 0;

  return {
    width: `${cardWidth}px`,
    height: `${cardHeight}px`,
    transform: `translateX(${horizontalOffset}px) rotate(${rotation}deg)`
  };
};

//不出牌
const pass = () => {
  let data = {
    type: "playCard",
    data: JSON.stringify({
      pid: state.countdownPlayer - 1,
      pass: 1,
    }),
  }
  websocket.send(data)

  //隐藏不出按钮
  state.countdownPlayer = 0
  //显示"过"
  player1pass.value = true
  //播过声音
  const musicPath = "guo"
  playSound(musicPath)

  //玩家出牌后下一位开始倒计时
  startCountdown(2, state.outCardTimeout)
  //隐藏2号的"过"
  player2pass.value = false

  //主动过牌
  proactivePass.value = true
}

// 添加返回首页的方法
const goToHome = () => {
  showGameOverModal.value = false
  router.push('/index')  // 跳转到首页路由
}

// 预加载 cards 目录下所有 .png 图片（Vite 特有的资源导入方式）
// eager: true 表示立即加载，import: 'default' 获取图片 URL
const cardImages = import.meta.glob('@/assets/img/cards/*.png', { eager: true, import: 'default' });

// 动态获取卡牌图片路径（利用预加载的图片映射）
const getCardImage = (cardId) => {
  const cardIdStr = String(cardId).trim();
  // 构造预加载图片的路径 key（与 glob 匹配的路径格式）
  const imageKey = `/src/assets/img/cards/${cardIdStr}.png`;
  // 返回预加载的图片 URL（开发/生产环境路径自动适配）
  return cardImages[imageKey] || ''; // 若找不到对应图片，返回空（可添加默认图）
};


// 计算属性：解析聊天框背景图路径（处理 @ 别名）
const chatlogBgUrl = computed(() => {
  // 使用 new URL() 构造路径，Vite 会自动解析 @ 别名
  return new URL('@/assets/img/ui/chatlog.png', import.meta.url).href;
});

// 添加返回房间的方法
const toRoom = () => {
  showGameOverModal.value = false
  //如果是好友房直接返回
  if (roomType.value !== 2){
    websocket.send({"type":"initRoom","data":"","name":""});
  }
  router.push({path:'/room'})
}

//重新设置头像
const resetAvatar = () => {
  //头像处理
  if (Math.random() > 0.5) {
    let player2avatar = getTouxiang()
    player2touxiang.value = player2avatar
    storage.local.set('player2avatar', player2avatar)
  }
  if (Math.random() > 0.5) {
    let player3avatar = getTouxiang()
    player3touxiang.value = player3avatar
    storage.local.set('player3avatar', player3avatar)
  }

  if (Math.random() > 0.5) {
    let player4avatar = getTouxiang()
    player4touxiang.value = player4avatar
    storage.local.set('player4avatar', player4avatar)
  }
}

const getPlayerPosition = (playerId) => {
  // 计算差值
  // const diff =playerId- currentPlayerId.value
  // // 当差值大于等于0，直接用，当差值小于0需要用4减去差值的绝对值
  // const position = diff >= 0 ? diff : 4 - Math.abs(diff)
  // return position

  //座位号 = (目标玩家索引 - 当前玩家索引 + 总人数) % 总人数
  const position = (playerId - currentPlayerId.value + 4) % 4
  return position
}
</script>

<style scoped>
:global(body) {
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.game {
  background: url('@/assets/img/bg/Table_Dif12324.png') no-repeat;
  background-size: cover;
  background-position: center;
  position: absolute;
  width: 100%;
  height: 100%;
  min-height: 100vh;
  overflow: hidden;
}

.player1_card img:hover {
  transform: translateY(-20px) translateX(calc(-1% * (13 - var(--n))));
  z-index: 10;
}

:global(body) {
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.game {
  background: url('@/assets/img/bg/Table_Dif12324.png') no-repeat;
  background-size: cover;
  background-position: center;
  position: absolute;
  width: 100%;
  height: 100%;
  min-height: 100vh;
  overflow: hidden;
}

.player1_card img:hover {
  transform: translateY(-20px) translateX(calc(-1% * (13 - var(--n))));
  z-index: 10;
}


/* 游戏结束弹窗样式 */
.game-over-modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  /* 移除原背景图设置 */
  padding: 20px;
  /* 减少内边距 */
  border-radius: 15px;
  text-align: center;
  width: 320px;
  /* 调整弹窗宽度 */
  display: flex;
  flex-direction: column;
  align-items: center;
}

/* 胜利/失败标题 (背景图外部顶部) */
.result-title {
  color: white;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.8);
  margin: 0 0 0 0;
  /* 底部 margin 分隔标题与背景图 */
  font-size: 26px;
  text-align: center;
}

/* 缩小的背景图容器 (移除内部标题相关样式) */
.result-bg {
  width: 220px;
  height: 160px;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  margin-bottom: 20px;
  /* 背景图与积分信息间距 */
  border-radius: 10px;
  /* 移除原 position: relative (不再需要内部定位) */
}

/* 积分信息 (背景图下方) */
.score-info {
  color: #333;
  font-size: 18px;
  margin-bottom: 25px;
  text-shadow: none;
  white-space: nowrap;
  /* 添加：强制不换行 */
}



.player1-container {
  position: absolute;
  bottom: 5%;
  left: 3%;
}

.player2-container {
  position: absolute;
  top: 30%;
  right: 5%;
}

.player4-container {
  position: absolute;
  top: 30%;
  left: 5%;
}

.player3-container {
  position: absolute;
  top: 5%;
  left: 40%;
  width: 90px;
  height: 90px;

}

/* 弃牌堆容器样式 */
.discard-pile-container {
  position: absolute;
  top: 30%;
  left: 50%;
  transform: translateX(-50%);
  display: inline-flex;
  align-items: center;
  flex-wrap: nowrap;
  padding: 10px 0;
}

/* 弃牌堆卡牌样式 */
.discard-card {
  width: 120px;
  height: 140px;
  margin-right: -15px;
  transition: transform 0.3s ease;
  flex-shrink: 0;
}

/* player1牌容器样式 */
.player1_card {
  position: absolute;
  top: 78%;
  left: 2%;
  right: 2%;
  display: flex;
  justify-content: center;
}

/* player1卡牌静态样式 */
.player1-card-item {
  width: 7%;
  /* 原内联width */
  height: 120px;
  /* 原内联height */
  margin-right: -1%;
  /* 原内联marginRight */
  transition: transform 0.3s ease;
  /* 原内联transition */
}


/* 功能区容器静态样式 */
.backhand-actions-container {
  position: absolute;
  top: 60%;
  left: 0;
  width: 100%;
}

/* 出牌按钮静态样式 */
.chupai-btn {
  position: absolute;
  top: 0;
  left: 54%;
  transition: transform 0.1s ease;
  /* 过渡效果移至CSS */
}

/* 禁用状态出牌按钮静态样式 */
.chupai-btn-disabled {
  position: absolute;
  top: 0;
  left: 54%;
}

/* 不出按钮静态样式 */
.buchu-btn {
  position: absolute;
  top: 0;
  left: 32%;
  transition: transform 0.1s ease;
  /* 过渡效果移至CSS */
}

/* 倒计时容器样式（抽取定位相关静态样式） */
.countdown-circle-container {
  position: absolute;
  top: 59%;
  left: 51%;
  transform: translateX(-50%);
  z-index: 999;
  text-align: center;
}

/* 背景圆样式（抽取静态视觉样式） */
.countdown-bg-circle {
  stroke: #eee;
  stroke-width: 8;
  fill: #555;
  /* 实心背景 */
}

/* 进度圆样式（抽取静态基础样式） */
.countdown-progress-circle {
  stroke-width: 8;
  stroke-dasharray: 226.19;
  /* 固定圆周长度（2πr = 2*3.14*36 ≈ 226.19） */
  transition: stroke-dashoffset 1s linear;
  /* 进度过渡动画 */
}

/* 倒计时文本样式（抽取静态文本样式） */
.countdown-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 26px;
  color: white;
  text-shadow: 0 0 5px rgba(0, 0, 0, 0.5);
}

/* 按钮容器样式 - 确保在所有屏幕尺寸下都不换行 */
.buttons-container {
  display: flex;
  gap: 15px;
  justify-content: center;
  margin-top: 0px;
  flex-wrap: nowrap;
  /* 强制不换行 */
  align-items: center;
}

/* 再来一局按钮样式 */
.restart-btn {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 12px 25px;
  font-size: 18px;
  font-weight: bold;
  border-radius: 8px;
  cursor: pointer;
  min-width: 140px;
  width: auto;
  box-shadow: 0 4px 12px rgba(76, 175, 80, 0.4);
  transition: all 0.3s ease;
  white-space: nowrap;
  text-align: center;
  height: 48px;
  line-height: 24px;
}

.restart-btn:hover {
  background-color: #45a049;
  transform: scale(1.05);
  box-shadow: 0 6px 16px rgba(76, 175, 80, 0.6);
}

/* 返回房间按钮样式 - 与再来一局按钮完全相同的尺寸和圆角 */
.room-btn {
  padding: 12px 25px;
  font-size: 18px;
  font-weight: bold;
  color: white;
  background: linear-gradient(45deg, #4169E1, #1E90FF);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  min-width: 140px;
  width: auto;
  box-shadow: 0 4px 12px rgba(30, 144, 255, 0.4);
  transition: all 0.3s ease;
  white-space: nowrap;
  text-align: center;
  height: 48px;
  line-height: 24px;
}

.room-btn:hover {
  transform: scale(1.05);
  box-shadow: 0 6px 16px rgba(30, 144, 255, 0.6);
}


/* 卡牌过渡动画容器 */
.card-transition-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 99;
}


/* 移动中的卡牌样式 - 确保初始位置正确 */
.moving-card {
  position: absolute;
  width: 80px;
  height: 112px;
  z-index: 100;
  top: 80%;
  /* 新增：默认位置设为玩家手牌区域 */
  left: 50%;
  transform: translateX(-50%);
}

/* 卡牌移动动画 */
.card-move-enter-active {
  transition: all 500ms cubic-bezier(0.25, 0.8, 0.25, 1);
}

.card-move-leave-active {
  display: none;
}

.card-move-enter-from {
  /* 开始位置：玩家手牌区域，添加opacity: 0防止闪烁 */
  top: 80%;
  left: 50%;
  transform: translateX(-50%) scale(1);
  opacity: 0;
  /* 新增：动画开始时透明 */
}

.card-move-enter-to {
  /* 结束位置：屏幕中间（弃牌堆位置） */
  top: 32%;
  left: 50%;
  transform: translateX(-50%);
  opacity: 1;
  /* 新增：动画结束时显示 */
}

/* 添加这个新类，处理动画结束后移除元素 */
.card-move-enter-to.card-move-complete {
  display: none;
  /* 动画完成后隐藏元素 */
}

/* 横屏模式专属样式 */
@media (max-width: 998px) {
  .player1-container {
    transform: scale(0.8);
    /* 缩小至80% */
    bottom: 6%;
    /* 向下移动（从原来的5%调整为8%） */
    transform-origin: bottom left;
  }

  .player2-container {
    transform: scale(0.8);
    /* 整体缩小至80% */
    transform-origin: top right;
  }

  .player4-container {
    transform: scale(0.8);
    /* 整体缩小至80% */
  }

  .player3-container {
    transform: scale(0.8);
    /* 整体缩小至80% */
  }

  .discard-pile-container {

    transform: translateX(-50%) scale(0.6);
    transform-origin: top center;
  }

  .player1_card {
    position: absolute;
    top: 72%;
    /* 原78% → 减少2%实现向上移动 */
    left: -9%;
    /* 原2% → 增加1%实现向右移动 */
    right: -19%;
    /* 保持右边界不变 */
    display: flex;
    justify-content: center;
    transform: scale(0.8);
  }

  .backhand-actions-container {
    transform: scale(0.6);
    /* 整体缩小至60% */
    transform-origin: top center;
    /* 以顶部中心为缩放原点，避免位置偏移 */
  }

  /* 不出按钮静态样式 */
  .buchu-btn {
    left: 22%;
  }

    /* 出牌按钮静态样式 */
  .chupai-btn {
    left: 57%;
  }

  /* 禁用状态出牌按钮静态样式 */
  .chupai-btn-disabled {
    left: 57%;
  }

  .countdown-circle-container {
    transform: translateX(-50%) scale(0.6);
    /* 保持水平居中并缩小至60% */
    transform-origin: top center;
    /* 以顶部中心为缩放原点，避免位置偏移 */
  }

  .buttons-container {
    gap: 12px;
  }

  .restart-btn,
  .room-btn {
    padding: 10px 20px;
    font-size: 16px;
    min-width: 120px;
    height: 44px;
    line-height: 24px;
  }

  /* 移动中的卡牌在小屏幕下缩小到0.4倍 */
  .moving-card {
    transform: scale(0.1);
  }

  /* 调整卡牌移动动画的位置，适配缩小后的尺寸 */
  .card-move-enter-to {
    top: 32%;
    left: 50%;
    transform: translateX(-50%) scale(0.1);
    opacity: 1;
    /* 新增：动画结束时显示 */
  }

}
</style>
