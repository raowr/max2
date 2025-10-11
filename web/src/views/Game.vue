<template>
  <div class="game">

    <div class="msg" style="position: absolute; top: 1%; left: 6%; transform: translateX(-50%);color:gold;">
      {{state.lastmsg}}
    </div>

    <!--弃牌堆 -->
    <!-- 通过 left: 50% + transform: translateX(-50%) 实现水平居中 -->
    <div
      style="position: absolute; top: 30%; left: 50%; transform: translateX(-50%); display: flex; align-items: center;">
      <img v-for="cardId in state.outCards" :key="cardId" :src="'@/assets/img/cards/'+cardId+'.png'"
        style="width: 100%; height: 140px; margin-right: -4%; transition: transform 0.3s ease;">
    </div>
    <!--弃牌堆结束 -->

    <!--player1 牌 -->
    <div class="player1_card"
      style="position: absolute; top:78%; left: 2%; right: 2%; display: flex; justify-content: center;">
      <img v-for="cardId in state.cards" :key="cardId" :src="getCardImage(cardId)" :style="{
      width: '7%',
      height: '120px',
      marginRight: '-1%',
      transform: selectedCards.includes(cardId) ? 'translateY(-20px)' : 'translateX(calc(-1% * (13 - ' + cardId + ')))',
      zIndex: selectedCards.includes(cardId) ? 10 : 1,
      transition: 'transform 0.3s ease'
    }" @click="toggleCard(cardId)">
    </div>

<!--player1 信息 -->
<div style="position: absolute;bottom:5%;left:3%;">
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
    <div style="position: absolute;top:30%;right:5%;">
      <!-- 修改为圆形倒计时 -->
      <div v-if="state.countdownPlayer == 2"
        style="position: absolute; top: -5%; left: 34%; transform: translateX(-50%); z-index: 999; text-align: center">
        <svg :width="100" :height="100">
          <circle cx="50" cy="50" r="45" stroke="#eee" stroke-width="8" fill="transparent" />
          <circle cx="50" cy="50" r="45" :stroke="countdownPlayer2 > 10 ? '#4CAF50' : '#ff5722'" stroke-width="8"
            fill="transparent" :style="{
            strokeDasharray: 283,
            strokeDashoffset: 283 * (1 - countdownPlayer2/30),
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
      <p style="z-index:1;font-size:16px; color:white;position: absolute;top:74%;left:-60%;">剩{{state.player2CardsNum}}张
      </p>
    </div>

    <!--player2 信息结束 -->

    <!--player4 信息 -->
    <div style="position: absolute;top:30%;left:3%;">
      <!-- 修改为圆形倒计时 -->
      <div v-if="state.countdownPlayer == 4"
        style="position: absolute; top: -4%; left: 40%; transform: translateX(-50%); z-index: 999; text-align: center">
        <svg :width="100" :height="100">
          <circle cx="50" cy="50" r="45" stroke="#eee" stroke-width="8" fill="transparent" />
          <circle cx="50" cy="50" r="45" :stroke="countdownPlayer4 > 10 ? '#4CAF50' : '#ff5722'" stroke-width="8"
            fill="transparent" :style="{
            strokeDasharray: 283,
            strokeDashoffset: 283 * (1 - countdownPlayer4/30),
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
        剩{{state.player4CardsNum}}张</p>
    </div>

    <!--player4 信息结束 -->

    <!--player3 信息 -->
    <div style="position: absolute;top:5%;left:40%;width: 90px;">
      <!-- 修改为圆形倒计时 -->
      <div v-if="state.countdownPlayer == 3"
        style="position: absolute; top: 8%; left: 50%; transform: translateX(-50%); z-index: 999; text-align: center">
        <svg :width="100" :height="100">
          <circle cx="50" cy="50" r="45" stroke="#eee" stroke-width="8" fill="transparent" />
          <circle cx="50" cy="50" r="45" :stroke="countdownPlayer3 > 10 ? '#4CAF50' : '#ff5722'" stroke-width="8"
            fill="transparent" :style="{
            strokeDasharray: 283,
            strokeDashoffset: 283 * (1 - countdownPlayer3/30),
            transition: 'stroke-dashoffset 1s linear'
          }" />
        </svg>
        <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); 
                 font-size: 24px; color: white; text-shadow: 0 0 5px rgba(0,0,0,0.5)">
          {{ countdownPlayer3 }}
        </div>
      </div>
      <div width="90px"
        style="background-image: url('@/assets/img/ui/chatlog.png'); background-size:90px;position: absolute;width: 90px;height: 90px;">
        <img src="@/assets/img/touxiang/bighead15729.png" width="90px"
          style="bottom:31.2%;left:3.1%;border-radius: 25px;">
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:120px;'>
      <p style="position: absolute;font-size:16px; color:white;top:90px;left:135px;width: 60px;">
        剩{{state.player3CardsNum}}张</p>
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
    <div v-if="state.countdownPlayer == 1"
      style="position: absolute; top: 59%; left: 51%; transform: translateX(-50%); z-index: 999; text-align: center">
      <!-- 尺寸从100调整为80（80%） -->
      <svg :width="80" :height="80">
        <!-- 圆形半径从45调整为36（45*0.8），背景圆改为实心深灰色 -->
        <circle cx="40" cy="40" r="36" stroke="#eee" stroke-width="8" fill="#555" />  <!-- 实心背景 -->
        <!-- 进度圆改为实心，颜色随倒计时变化 -->
        <circle cx="40" cy="40" r="36" 
          :stroke="countdownPlayer1 > 10 ? '#4CAF50' : '#ff5722'" 
          :fill="countdownPlayer1 > 10 ? '#81C784' : '#FF9800'"  
          stroke-width="8"
          :style="{
            strokeDasharray: 226.19,
            strokeDashoffset: 226.19 * (1 - countdownPlayer1/30),
            transition: 'stroke-dashoffset 1s linear'
          }" />
      </svg>
      <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); 
                 font-size: 26px; color: white; text-shadow: 0 0 5px rgba(0,0,0,0.5)">  <!-- 字体从19px放大到26px -->
        {{ countdownPlayer1 }}
      </div>
    </div>



    <!-- 后手出牌时的功能 !-->
    <div v-if="state.countdownPlayer == 1">
      <!-- 修改：出牌按钮添加按下缩小效果 -->
      <img 
        v-if="checkOut()" 
        src='@/assets/img/btn_chupai.png' 
        :style="{
          position: 'absolute',
          top: '60%',
          left: '54%',
          transform: isChupaiBtnPressed ? 'scale(0.9)' : 'scale(1)',  // 按下缩小至90%
          transition: 'transform 0.1s ease'  // 平滑过渡动画
        }" 
        @click='chupai()'
        @mousedown="isChupaiBtnPressed = true" 
        @mouseup="isChupaiBtnPressed = false"   
        @mouseleave="isChupaiBtnPressed = false"  
      >
      <img v-else src='@/assets/img/btn_chupai_hui.png' style='position: absolute;top:60%;left:54%'>

      <img 
        v-if="state.mustPid != 0"
        src='@/assets/img/btn_bujiao2.png' 
        :style="{
          position: 'absolute',
          top: '60%',
          left: '32%',
          transform: isPassBtnPressed ? 'scale(0.9)' : 'scale(1)',  // 按下时缩小到90%
          transition: 'transform 0.1s ease'  // 平滑过渡动画
        }" 
        @click='pass()'
        @mousedown="isPassBtnPressed = true"
        @mouseup="isPassBtnPressed = false"
        @mouseleave="isPassBtnPressed = false"
      >
    </div>

    <!-- 后手出牌时的功能 !-->


 <!-- 游戏结束弹窗 -->
<div v-if="showGameOverModal" class="game-over-modal">
  <div class="modal-content">
    <!-- 胜利/失败标题 (移至背景图外部顶部) -->
    <h2 class="result-title">{{ isWinner ? '恭喜胜利!' : '游戏失败' }}</h2>
    
    <!-- 缩小的背景图容器 -->
    <div 
      class="result-bg" 
      :style="{ 
        backgroundImage: isWinner ? `url(${winBg})` : `url(${loseBg})`
      }"
    >
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
    
      <!-- 修改按钮文本和点击事件 -->
      <button 
        class="restart-btn" 
        @click="state.player1Point === 0 ? goToHome() : restartGame()"
      >
        {{ state.player1Point === 0 ? '返回首页' : '再来一局' }}
      </button>
  </div>
</div>

  </div>
</template>

<script setup>
import { ref,  reactive,onMounted,onBeforeUnmount  } from 'vue'
import { audioManager } from '@/utils/audio'
import { websocket } from '@/utils/websocket'
const state = reactive({
  deck:[],
  countdownPlayer:0,
  cards:[],
  player2CardsNum:13,
  player3CardsNum:13,
  player4CardsNum:13,
  outCards:[],
  lastmsg:"",
  mustPid:0,
  player1Point: 0, // 初始瓜子数
})
const isPassBtnPressed = ref(false)  // 添加：跟踪按钮按下状态
const isChupaiBtnPressed = ref(false)  // 添加：跟踪"出牌"按钮按下状态

// 添加游戏结束弹窗相关状态
const showGameOverModal = ref(false)
const isWinner = ref(false)
const gameOverMessage = ref('')

// 导入背景图片（新增代码）
import winBg from '@/assets/img/ui/win_bg.png'
import loseBg from '@/assets/img/ui/lose_bg.png'


onMounted(() => {
  initDeck()
  const bgmUrl = new URL('@/assets/music/game_bg1.mp3', import.meta.url).href;
  audioManager.preload('bgm', bgmUrl); // 使用解析后的 URL
   audioManager.playBGM('bgm')
   websocket.send({"type":"play","data":"","name":""})
   websocket.on('message', handleMessage)
})

// 初始化一副牌
const initDeck= () => {
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
    //先出牌的开始倒计时
    startCountdown(data.current+1)
    //如果必出是玩家，记录下必出玩家的pid
    state.mustPid = data.current
    //更新玩家总瓜子数
    state.player1Point = data.playerPoint
  }
  if(data.type == "outCard"){
    data = JSON.parse(data.data)
        // 刷新玩家牌数
        let cardsMsg = ""
    switch (data.pid) {
      case 0:
        if(data.code == 0){
          // state.must = data.must
          // selectedCards.value = []
          // state.outCards=(data.cards)
          // startCountdown(2)
          for(let i=0;i<data.cardIds.length;i++){
            for(let j=0;j<state.deck.length;j++){
              if(state.deck[j].Id == data.cardIds[i]){
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
          for(let i=0;i<data.cardIds.length;i++){
            for(let j=0;j<state.deck.length;j++){
              if(state.deck[j].Id == data.cardIds[i]){
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
          for(let i=0;i<data.cardIds.length;i++){
            for(let j=0;j<state.deck.length;j++){
              if(state.deck[j].Id == data.cardIds[i]){
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
          for(let i=0;i<data.cardIds.length;i++){
            for(let j=0;j<state.deck.length;j++){
              if(state.deck[j].Id == data.cardIds[i]){
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
      startCountdown(data.current + 1)
      //如果必出是玩家，记录下必出玩家的pid
      state.mustPid = data.mustPid
    }

  }
  if (data.type == "pass") {
    data = JSON.parse(data.data)
    state.mustPid = data.mustPid
    startCountdown(data.current+1)

  }
  if(data.type == "over"){
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
  websocket.send({"type":"play","data":"","name":""})
}

const selectedCards = ref([])  // 改为数组存储选中状态

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
  }else{
    return true
  }
}


// 新增倒计时逻辑
const countdownPlayer1 = ref(30) // 初始30秒
const countdownPlayer2 = ref(30) // 初始30秒
const countdownPlayer3 = ref(30) // 初始30秒
const countdownPlayer4 = ref(30) // 初始30秒
let timer1 = null
let timer2 = null
let timer3 = null
let timer4 = null

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
      countdownPlayer1.value = 30
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
      countdownPlayer2.value = 30
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
      countdownPlayer3.value = 30
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
      countdownPlayer4.value = 30
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
  if(!timer) startCountdown()
}


onBeforeUnmount(() => {
  clearInterval(timer)
})

const chupai = () => {
  if (selectedCards.value.length === 0) {
    return false
  }else{
    let cards =[]
    selectedCards.value.forEach(item => {
      state.cards.forEach(card => {
        if(card == item){
          cards.push(card)
        }
      })
    })
    let data = {
      type: "playCard",
      data: JSON.stringify({
        pid: state.countdownPlayer-1,
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
      pid: state.countdownPlayer-1,
      pass: 1,
    }),
  }
  websocket.send(data)
}

// 添加返回首页的方法
const goToHome = () => {
  showGameOverModal.value = false
  router.push('/')  // 跳转到首页路由
}

// 定义动态获取卡牌图片路径的方法（解析 @ 别名）
const getCardImage = (cardId) => {
  console.log(cardId)
  // 使用 new URL() 构造路径，Vite 会自动解析 @ 别名并处理资源
  const cardIdStr = String(cardId);
  return new URL(`@/assets/img/cards/${cardIdStr}.png`, import.meta.url).href;
};


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
  padding: 20px;  /* 减少内边距 */
  border-radius: 15px;
  text-align: center;
  width: 320px;  /* 调整弹窗宽度 */
  display: flex;
  flex-direction: column;
  align-items: center;
}

/* 胜利/失败标题 (背景图外部顶部) */
.result-title {
  color: white;
  text-shadow: 0 2px 4px rgba(0,0,0,0.8);
  margin: 0 0 0 0; /* 底部 margin 分隔标题与背景图 */
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
  margin-bottom: 20px; /* 背景图与积分信息间距 */
  border-radius: 10px;
  /* 移除原 position: relative (不再需要内部定位) */
}

/* 积分信息 (背景图下方) */
.score-info {
  color: #333;
  font-size: 18px;
  margin-bottom: 25px;
  text-shadow: none;
  white-space: nowrap; /* 添加：强制不换行 */
}


/* 按钮 (积分信息下方) */
.restart-btn {
  background-color: #4CAF50;
  color: white;
  border: none;
  padding: 12px 24px;
  font-size: 18px;
  border-radius: 8px;
  cursor: pointer;
  width: 180px;
  margin-top: 10px;
}

.restart-btn:hover {
  background-color: #45a049;
}
</style>
