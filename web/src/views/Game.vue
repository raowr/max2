<template>
  <div class="game">

<!--弃牌堆 -->
    <div >
    <img id="" src='@/assets/img/1.png' width='7%' height='120px' style='position: absolute; top: 30%;left:27%'>
    
</div>
<!--弃牌堆结束 -->

<!--player1 牌 -->
<div class="player1_card" style="position: absolute; top:78%; left: 2%; right: 2%; display: flex; justify-content: center;">
  <img 
    v-for="n in 13" 
    :key="n" 
    :src="'src/assets/img/'+n+'.png'"
    :style="{
      width: '7%',
      height: '120px',
      marginRight: '-1%',
      transform: selectedCards.includes(n) ? 'translateY(-20px)' : 'translateX(calc(-1% * (13 - ' + n + ')))',
      zIndex: selectedCards.includes(n) ? 10 : 1,
      transition: 'transform 0.3s ease'
    }"
    @click="toggleCard(n)"
  >
</div>

    <!--player1 信息 -->

    <div style="position: absolute;bottom:5%;left:3%;">
      <img src="@/assets/img/ui/chatlog.png" width="110px">
      <img src="@/assets/img/touxiang/bighead15718.png" width="100px"
        style="position: absolute;bottom:38.2%;left:3%;border-radius: 25px;">
      <div style="width:140px;height:40px;text-align:center;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥1</p>
      </div>

    </div>

    <!--player1 信息结束 -->

    <!--player4 信息 -->
    <div style="position: absolute;top:30%;right:5%;">
          <!-- 修改为圆形倒计时 -->
    <div v-if="state.countdownPlayer == 4" style="position: absolute; top: -5%; left: 34%; transform: translateX(-50%); z-index: 999; text-align: center">
      <svg :width="100" :height="100">
        <circle 
          cx="50" 
          cy="50" 
          r="45" 
          stroke="#eee" 
          stroke-width="8" 
          fill="transparent"
        />
        <circle
          cx="50"
          cy="50"
          r="45"
          :stroke="countdownPlayer4 > 10 ? '#4CAF50' : '#ff5722'"
          stroke-width="8"
          fill="transparent"
          :style="{
            strokeDasharray: 283,
            strokeDashoffset: 283 * (1 - countdownPlayer4/30),
            transition: 'stroke-dashoffset 1s linear'
          }"
        />
      </svg>
      <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); 
                 font-size: 24px; color: white; text-shadow: 0 0 5px rgba(0,0,0,0.5)">
        {{ countdownPlayer4 }}
      </div>
    </div>
      <img src="@/assets/img/ui/chatlog.png" width="90px">
      <img src="@/assets/img/touxiang/bighead15419.png" width="90px"
        style="position: absolute;bottom:42.2%;left:3.1%;border-radius: 25px;">
      <div style="width:140px;height:40px;text-align:center;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥4</p>
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:-70%;top:-2%'>
      <p style="z-index:1;font-size:16px; color:white;position: absolute;top:74%;left:-60%;">剩10张</p>
    </div>

    <!--player4 信息结束 -->

    <!--player2 信息 -->
    <div style="position: absolute;top:30%;left:3%;">
          <!-- 修改为圆形倒计时 -->
    <div v-if="state.countdownPlayer == 2" style="position: absolute; top: -4%; left: 40%; transform: translateX(-50%); z-index: 999; text-align: center">
      <svg :width="100" :height="100">
        <circle 
          cx="50" 
          cy="50" 
          r="45" 
          stroke="#eee" 
          stroke-width="8" 
          fill="transparent"
        />
        <circle
          cx="50"
          cy="50"
          r="45"
          :stroke="countdownPlayer2 > 10 ? '#4CAF50' : '#ff5722'"
          stroke-width="8"
          fill="transparent"
          :style="{
            strokeDasharray: 283,
            strokeDashoffset: 283 * (1 - countdownPlayer2/30),
            transition: 'stroke-dashoffset 1s linear'
          }"
        />
      </svg>
      <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); 
                 font-size: 24px; color: white; text-shadow: 0 0 5px rgba(0,0,0,0.5)">
        {{ countdownPlayer2 }}
      </div>
    </div>
      <img src="@/assets/img/ui/chatlog.png" width="90px">
      <img src="@/assets/img/touxiang/bighead15339.png" width="85px"
        style="position: absolute;bottom:42.2%;left:3.1%;border-radius: 25px;">
      <div style="width:100px;height:40px;text-align:center;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥2</p>
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:110%;top:0%'>
      <p style="z-index:1;font-size:16px; color:white;position: absolute;top:67%;left:130%;width: 60px;">剩10张</p>
    </div>

    <!--player2 信息结束 -->

    <!--player3 信息 -->
    <div style="position: absolute;top:5%;left:40%;width: 90px;" >
          <!-- 修改为圆形倒计时 -->
    <div v-if="state.countdownPlayer == 3" style="position: absolute; top: 8%; left: 50%; transform: translateX(-50%); z-index: 999; text-align: center">
      <svg :width="100" :height="100">
        <circle 
          cx="50" 
          cy="50" 
          r="45" 
          stroke="#eee" 
          stroke-width="8" 
          fill="transparent"
        />
        <circle
          cx="50"
          cy="50"
          r="45"
          :stroke="countdownPlayer3 > 10 ? '#4CAF50' : '#ff5722'"
          stroke-width="8"
          fill="transparent"
          :style="{
            strokeDasharray: 283,
            strokeDashoffset: 283 * (1 - countdownPlayer3/30),
            transition: 'stroke-dashoffset 1s linear'
          }"
        />
      </svg>
      <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); 
                 font-size: 24px; color: white; text-shadow: 0 0 5px rgba(0,0,0,0.5)">
        {{ countdownPlayer3 }}
      </div>
    </div>
      <div width="90px"
        style="background-image: url('src/assets/img/ui/chatlog.png'); background-size:90px;position: absolute;width: 90px;height: 90px;">
        <img src="@/assets/img/touxiang/bighead15729.png" width="90px"
          style="bottom:31.2%;left:3.1%;border-radius: 25px;">
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:120px;'>
      <p style="position: absolute;font-size:16px; color:white;top:90px;left:135px;width: 60px;">剩10张</p>
      <div style="position: absolute;width:105px;height:40px;text-align:center;top:90px;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥3</p>
      </div>

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
    <div v-if="state.countdownPlayer == 1" style="position: absolute; top: 58%; left: 68%; transform: translateX(-50%); z-index: 999; text-align: center">
      <svg :width="100" :height="100">
        <circle 
          cx="50" 
          cy="50" 
          r="45" 
          stroke="#eee" 
          stroke-width="8" 
          fill="transparent"
        />
        <circle
          cx="50"
          cy="50"
          r="45"
          :stroke="countdownPlayer1 > 10 ? '#4CAF50' : '#ff5722'"
          stroke-width="8"
          fill="transparent"
          :style="{
            strokeDasharray: 283,
            strokeDashoffset: 283 * (1 - countdownPlayer1/30),
            transition: 'stroke-dashoffset 1s linear'
          }"
        />
      </svg>
      <div style="position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); 
                 font-size: 24px; color: white; text-shadow: 0 0 5px rgba(0,0,0,0.5)">
        {{ countdownPlayer1 }}
      </div>
    </div>

    <!-- 先手出牌时的功能 !-->
    <img src='@/assets/img/btn_chupai.png' style='position: absolute;top:60%;left:35%' onclick='chupai()'>
    <!-- 先手出牌时的功能 !-->


    <!-- 后手出牌时的功能 !-->
    <img src='@/assets/img/btn_chupai.png' style='position: absolute;top:60%;left:35%' onclick='chupai()'>
    <img src='@/assets/img/btn_bujiao2.png' style='position: absolute;top:60%;left: 48%' onclick='buyao()'>
    <!-- 后手出牌时的功能 !-->


    <!-- 要不起不能出牌时的功能 !-->
    <img src='@/assets/img/btn_chupai_hui.png' style='position: absolute;top:60%;left:35%'>
    <img src='@/assets/img/btn_bujiao2.png' style='position: absolute;top:60%;left: 48%' onclick='buyao()'>
    <!-- 要不起不能出牌时的功能 !-->

  </div>
</template>

<script setup>
import { ref,  reactive,onMounted,onBeforeUnmount  } from 'vue'
import { audioManager } from '@/utils/audio'
import { websocket } from '@/utils/websocket'
const state = reactive({
  countdownPlayer:1
})

onMounted(() => {
  audioManager.preload('bgm', 'src/assets/music/game_bg1.mp3')
   audioManager.playBGM('bgm')
   websocket.send({"type":"play","data":"","name":""})
   websocket.on('message', handleMessage)
   startCountdown(1)
})

const handleMessage = (data) => {
  // 处理接收到的消息
  console.log('Received message:', data)
  if(data.type == "over"){
    // 重复玩
     websocket.send({"type":"play","data":"","name":""})
  }
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
          startCountdown(2)
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
          startCountdown(3)
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
          startCountdown(4)
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
          startCountdown(1)
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

</style>
