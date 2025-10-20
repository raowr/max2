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

    <!--player1 牌 -->
    <!-- 移除容器的内联静态样式，保留类名 -->
    <div class="player1_card">
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

      <img src="@/assets/img/touxiang/bighead15718.png" width="100px"
        style="position: absolute;bottom:38.2%;left:3%;border-radius: 25px;">
      <div style="width:140px;height:40px;text-align:center;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥1</p>
      </div>
    </div>
    <!--player1 信息结束 -->

    <!--player2 信息 -->
    <div class="player2-container" style="">
    <!-- 修改为圆形倒计时 -->
    <div v-if="state.countdownPlayer == 2"
      style="position: absolute; top: -5%; left: 34%; transform: translateX(-50%); z-index: 999; text-align: center">
      <svg :width="100" :height="100">
        <circle cx="50" cy="50" r="45" stroke="#eee" stroke-width="8" fill="transparent" />
        <circle cx="50" cy="50" r="45" :stroke="countdownPlayer2 > (state.outCardTimeout / 3) ? '#4CAF50' : '#ff5722'" stroke-width="8"
          fill="transparent" :style="{
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
      <img src="@/assets/img/touxiang/bighead15419.png" width="90px"
        style="position: absolute;bottom:42.2%;left:3.1%;border-radius: 25px;">
      <div style="width:140px;height:40px;text-align:center;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥2</p>
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:-70%;top:-2%'>
      <p style="z-index:1;font-size:16px; color:white;position: absolute;top:66%;left:-60%;">剩{{ state.player2CardsNum }}张
      </p>
    </div>

    <!--player2 信息结束 -->

    <!--player4 信息 -->
    <div class="player4-container" style="">
 <!-- 修改为圆形倒计时 -->
    <div v-if="state.countdownPlayer == 4"
      style="position: absolute; top: -4%; left: 40%; transform: translateX(-50%); z-index: 999; text-align: center">
      <svg :width="100" :height="100">
        <circle cx="50" cy="50" r="45" stroke="#eee" stroke-width="8" fill="transparent" />
        <circle cx="50" cy="50" r="45" :stroke="countdownPlayer4 > (state.outCardTimeout / 3) ? '#4CAF50' : '#ff5722'" stroke-width="8"
          fill="transparent" :style="{
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
      <img src="@/assets/img/touxiang/bighead15339.png" width="85px"
        style="position: absolute;bottom:42.2%;left:3.1%;border-radius: 25px;">
      <div style="width:100px;height:40px;text-align:center;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥4</p>
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:110%;top:0%'>
      <p style="z-index:1;font-size:16px; color:white;position: absolute;top:67%;left:130%;width: 60px;">
        剩{{ state.player4CardsNum }}张</p>
    </div>

    <!--player4 信息结束 -->

    <!--player3 信息 -->
    <div class="player3-container" style="">
    <!-- 修改为圆形倒计时 -->
    <div v-if="state.countdownPlayer == 3"
      style="position: absolute; top: 8%; left: 50%; transform: translateX(-50%); z-index: 999; text-align: center">
      <svg :width="100" :height="100">
        <circle cx="50" cy="50" r="45" stroke="#eee" stroke-width="8" fill="transparent" />
        <circle cx="50" cy="50" r="45" :stroke="countdownPlayer3 > (state.outCardTimeout / 3) ? '#4CAF50' : '#ff5722'" stroke-width="8"
          fill="transparent" :style="{
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
        <img src="@/assets/img/touxiang/bighead15729.png" width="90px"
          style="bottom:31.2%;left:3.1%;border-radius: 25px;">
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:120px;'>
      <p style="position: absolute;font-size:16px; color:white;top:90px;left:135px;width: 60px;">
        剩{{ state.player3CardsNum }}张</p>
      <div style="position: absolute;width:105px;height:40px;text-align:center;top:90px;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥3</p>
      </div>

    </div>

    <!--player3 信息结束 -->

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
        :stroke="countdownPlayer1 > (state.outCardTimeout / 3) ? '#4CAF50' : '#ff5722'" :fill="countdownPlayer1 > (state.outCardTimeout / 3) ? '#81C784' : '#FF9800'"
        :style="{
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

      <!-- 不叫按钮添加类名：bujiao-btn -->
      <img v-if="state.mustPid != 0" src='@/assets/img/btn_bujiao2.png' class="bujiao-btn" :style="{
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

        <!-- 修改按钮区域，交换按钮顺序 -->
        <div class="buttons-container">
          <button class="room-btn" @click="goToRoom()">
            返回房间
          </button>
          <button class="restart-btn" @click="state.player1Point === 0 ? goToHome() : restartGame()">
            {{ state.player1Point === 0 ? '返回首页' : '再来一局' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRouter } from 'vue-router' // 添加这行导入
import { audioManager } from '@/utils/audio'
import { websocket } from '@/utils/websocket'
const state = reactive({
  deck: [],
  countdownPlayer: 0,
  cards: [],
  player2CardsNum: 13,
  player3CardsNum: 13,
  player4CardsNum: 13,
  outCards: [],
  lastmsg: "",
  mustPid: 0,
  player1Point: 0, // 初始瓜子数
  outCardTimeout: 30, // 出牌最大时间(单位秒) /s
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

// 导入背景图片（新增代码）
import winBg from '@/assets/img/ui/win_bg.png'
import loseBg from '@/assets/img/ui/lose_bg.png'


onMounted(() => {
  initDeck()
  const bgmUrl = new URL('@/assets/music/game_bg1.mp3', import.meta.url).href;
  audioManager.preload('bgm', bgmUrl); // 使用解析后的 URL
  audioManager.playBGM('bgm')
  websocket.send({ "type": "play", "data": "", "name": "" })
  websocket.on('message', handleMessage)
})

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
        Value: 3 + i,
        Suit: suit,
        Name: suitName + valueName,
        Id: cardId,
        Rank: valueName,
        SuitName: suitName
      };
      state.deck.push(card);
    }
  }
  console.log(state.deck);
}

const handleMessage = (data) => {
  // 处理接收到的消息
  console.log('Received message:', data)
  if (data.type == "showCard") {
    data = JSON.parse(data.data)
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
    startCountdown(data.current + 1)
    //如果必出是玩家，记录下必出玩家的pid
    state.mustPid = data.current
    //更新玩家总瓜子数
    state.player1Point = data.playerPoint
    // 重置倒计时
    countdownPlayer1.value = state.outCardTimeout
    //清空选中的牌
    selectedCards.value = []
  }
  if (data.type == "outCard") {
    data = JSON.parse(data.data)
    // 刷新玩家牌数
    let cardsMsg = ""
    switch (data.pid) {
      case 0:
        if (data.code == 0) {
          // state.must = data.must
          // selectedCards.value = []
          // state.outCards=(data.cards)
          // startCountdown(2)
          for (let i = 0; i < data.cardIds.length; i++) {
            for (let j = 0; j < state.deck.length; j++) {
              if (state.deck[j].Id == data.cardIds[i]) {
                cardsMsg += state.deck[j].Name + " "
              }
            }
          }
          state.lastmsg = "玩家1出了：" + cardsMsg

          // 过滤掉state.cards中存在于data.cards的牌
          state.cards = state.cards.filter(item =>
            !data.cardIds.some(cardId => cardId === item)  // 取反some的结果，排除包含的牌
          );
          selectedCards.value = []
        }
        break
      case 1:
        state.player2CardsNum = data.cards_num
        // state.outCards=(data.cards)
        // startCountdown(3)
        for (let i = 0; i < data.cardIds.length; i++) {
          for (let j = 0; j < state.deck.length; j++) {
            if (state.deck[j].Id == data.cardIds[i]) {
              cardsMsg += state.deck[j].Name + " "
            }
          }
        }
        state.lastmsg = "玩家2出了：" + cardsMsg
        break
      case 2:
        state.player3CardsNum = data.cards_num
        // state.outCards=(data.cards)
        // startCountdown(4)
        for (let i = 0; i < data.cardIds.length; i++) {
          for (let j = 0; j < state.deck.length; j++) {
            if (state.deck[j].Id == data.cardIds[i]) {
              cardsMsg += state.deck[j].Name + " "
            }
          }
        }
        state.lastmsg = "玩家3出了：" + cardsMsg
        break
      case 3:
        state.player4CardsNum = data.cards_num
        // state.outCards=(data.cards)
        // startCountdown(1)
        for (let i = 0; i < data.cardIds.length; i++) {
          for (let j = 0; j < state.deck.length; j++) {
            if (state.deck[j].Id == data.cardIds[i]) {
              cardsMsg += state.deck[j].Name + " "
            }
          }
        }
        state.lastmsg = "玩家4出了：" + cardsMsg
        //如果必出是玩家
        break
    }

    //如果成功再改内容
    if (data.code == 0) {
      state.outCards = (data.cardIds)
      // 更新倒计时时长（如果服务端提供）
      if (data.outCardTimeout !== undefined) {
        state.outCardTimeout = data.outCardTimeout
      }
      startCountdown(data.current + 1)
      //如果必出是玩家，记录下必出玩家的pid
      state.mustPid = data.mustPid
    }

  }
  if (data.type == "pass") {
    data = JSON.parse(data.data)
    state.mustPid = data.mustPid
    startCountdown(data.current + 1)

  }
  if (data.type == "over") {
    clearInterval(timer1)
    clearInterval(timer2)
    clearInterval(timer3)
    clearInterval(timer4)
    state.countdownPlayer = 0
    data = JSON.parse(data.data)
    state.lastmsg = "游戏结束，玩家: " + data.winName + "胜利,赢得:" + data.win + "颗瓜子"

    state.player1Point = data.playerPoint
    // 显示游戏结束弹窗
    showGameOverModal.value = true
    // 判断当前玩家是否胜利（假设当前玩家是"帅哥1"）
    isWinner.value = data.winner === 0
    gameOverMessage.value = isWinner.value
      ? `哇！恭喜你赢得了${data.playerWin}颗瓜子!`
      : `咦！你很菜这局你输了${data.playerWin}颗瓜子`

    // 重置弃牌
    state.outCards = []
    // 移除自动重新开始的逻辑，由用户点击按钮触发
  }
}

// 添加重新开始游戏的方法
const restartGame = () => {
  showGameOverModal.value = false
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

const checkOut = () => {
  if (selectedCards.value.length === 0) {
    return false
  } else {
    return true
  }
}




const startCountdown = (pid) => {
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
      countdownPlayer1.value = state.outCardTimeout // 使用从服务端获取的倒计时时长
      timer1 = setInterval(() => {
        if (countdownPlayer1.value > 0) {
          countdownPlayer1.value--
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
      countdownPlayer2.value = state.outCardTimeout // 使用从服务端获取的倒计时时长
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
      countdownPlayer3.value = state.outCardTimeout // 使用从服务端获取的倒计时时长
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
      countdownPlayer4.value = state.outCardTimeout // 使用从服务端获取的倒计时时长
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

const chupai = () => {
  if (selectedCards.value.length === 0) {
    return false
  } else {
    let cards = []
    selectedCards.value.forEach(item => {
      state.cards.forEach(card => {
        if (card == item) {
          cards.push(card)
        }
      })
    })
    let data = {
      type: "playCard",
      data: JSON.stringify({
        pid: state.countdownPlayer - 1,
        cardIds: cards,
      }),
    }
    // console.log(data)
    websocket.send(data)
    // selectedCards.value = []
  }
}

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
const goToRoom = () => {
  showGameOverModal.value = false
  router.push('/room')  // 跳转到房间路由
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
  left: 3%;
}

.player3-container {
  position: absolute;
  top: 5%;
  left: 40%;
  width: 90px;
}

/* 弃牌堆容器样式 */
.discard-pile-container {
  position: absolute;
  top: 30%;
  left: 47%;
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
  left: 57%;
  transition: transform 0.1s ease;
  /* 过渡效果移至CSS */
}

/* 禁用状态出牌按钮静态样式 */
.chupai-btn-disabled {
  position: absolute;
  top: 0;
  left: 54%;
}

/* 不叫按钮静态样式 */
.bujiao-btn {
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
  margin-top: 20px;
  flex-wrap: nowrap; /* 强制不换行 */
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
/* 横屏模式专属样式 */
@media (max-width: 998px) {
  .player1-container {
    transform: scale(0.8);
    /* 缩小至80% */
    bottom: -4%;
    /* 向下移动（从原来的5%调整为8%） */
    transform-origin: bottom left;
    /* 确保缩放从左下角开始，避免水平偏移 */
  }

  .player2-container {
    transform: scale(0.8);
    /* 整体缩小至80% */
    transform-origin: top right;
    /* 以右上角为缩放原点，避免位置偏移 */
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

  /* 不叫按钮静态样式 */
  .bujiao-btn {
    left: 22%;
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
  
  .restart-btn, .room-btn {
    padding: 10px 20px;
    font-size: 16px;
    min-width: 120px;
    height: 44px;
    line-height: 24px;
  }

}
</style>
