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

    <!--player2 信息 -->
    <div style="position: absolute;top:30%;right:5%;">
      <img src="@/assets/img/ui/chatlog.png" width="90px">
      <img src="@/assets/img/touxiang/bighead15419.png" width="90px"
        style="position: absolute;bottom:42.2%;left:3.1%;border-radius: 25px;">
      <div style="width:140px;height:40px;text-align:center;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥2</p>
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:-70%;top:5%'>
      <p style="z-index:1;font-size:16px; color:white;position: absolute;top:74%;left:-60%;">剩10张</p>
    </div>

    <!--player2 信息结束 -->

    <!--player3 信息 -->
    <div style="position: absolute;top:30%;left:3%;">
      <img src="@/assets/img/ui/chatlog.png" width="90px">
      <img src="@/assets/img/touxiang/bighead15339.png" width="85px"
        style="position: absolute;bottom:42.2%;left:3.1%;border-radius: 25px;">
      <div style="width:100px;height:40px;text-align:center;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥3</p>
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:110%;top:0%'>
      <p style="z-index:1;font-size:16px; color:white;position: absolute;top:67%;left:130%;width: 60px;">剩10张</p>
    </div>

    <!--player3 信息结束 -->

    <!--player4 信息 -->
    <div style="position: absolute;top:5%;left:40%;">
      <div
        style="background-image: url('src/assets/img/ui/chatlog.png'); background-size:90px;position: absolute;width: 90px;height: 90px;">
        <img src="@/assets/img/touxiang/bighead15729.png" width="90px"
          style="bottom:31.2%;left:3.1%;border-radius: 25px;">
      </div>
      <img src='@/assets/img/54.png' width='80px' style='position: absolute;z-index:1;left:120px;'>
      <p style="position: absolute;font-size:16px; color:white;top:90px;left:135px;width: 60px;">剩10张</p>
      <div style="position: absolute;width:105px;height:40px;text-align:center;top:90px;">
        <p style="z-index:1;font-size:16px; color:white;">帅哥4</p>
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
import { ref,  onMounted } from 'vue'
import { audioManager } from '@/utils/audio'
import { websocket } from '@/utils/websocket'

onMounted(() => {
  audioManager.preload('bgm', 'src/assets/music/game_bg1.mp3')
   audioManager.playBGM('bgm')
   websocket.send({"type":"play","data":"","name":""})
   websocket.on('message', handleMessage)
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
