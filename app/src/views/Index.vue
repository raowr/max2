<template>
  <div class="index">
    <div v-show="userIdShow" class="msg" style="position: absolute; top: 0%;left: 11%; color:gold;">
      {{ storage.local.get('user_id') }}
    </div>
    <div id="right">
      <div class="right1" @click="toRoom()"></div>
      <div class="right2"  @click="right()" ></div>
      <div class="right3"></div>
    </div>
    <div class="yourenchang" id="charubiaoqian1">
      <img class="title_bg" src="@/assets/img/ui/title_bg.png" alt="标题背景" />
      <img class="friendroom" src="@/assets/img/ui/txt_friendroom.png" />
      <div class="img_return2" @click="return2()"><img src="@/assets/img/ui/img_return2.png" alt="标题背景" /></div>
      <br /><br /><a @click="inroom()"><img src="@/assets/img/ui/bg_abmatch.png" class="bg_abmatch1" /></a>
      <img src="@/assets/img/ui/btn_create_room.png" class="btn_create_room"  @click="createRoom()"/>
      <img class="tips" src="@/assets/img/ui/tips.png" />
      <p class="free">限时免费</p>
      <img src="@/assets/img/ui/bg_abmatch.png" class="bg_abmatch" @click="inroom()" />
      <img src="@/assets/img/ui/w_joinroom.png" class="w_joinroom" @click="inroom()" />
      <div>
        <input type="text" placeholder="房间号" name="roomid" class="roominput" id="inroom" />
        <input type="submit" value="确定" class="roomsubmit" id="btn_inroom" @click="joinroom()"/>
      </div>
      <img src="@/assets/img/ui/tips.png" class="tips2" />
      <p class="free2">限时免费</p>
      <!-- </right> -->
    </div>

    <!--  底部按钮 -->
    <div class="bottom-buttons" style="display: flex">
      <a href=""></a>
      <div class="btn btn-1"></div>
      <div class="btn btn-2"></div>
      <div class="btn btn-3"></div>
      <div class="btn btn-4"></div>
      <div class="btn btn-5"></div>
      <div class="btn btn-6"></div>
      <div class="btn btn-7"></div>
    </div>
    <!--  右边设置栏 -->
    <div class="right-settings">
      <div class="setting-btn btn-set" @click="toggleSettingImage()"></div>
      <div class="setting-btn btn-guide"  @click="toggleGuideModal()"></div>
      <div class="setting-btn btn-trophy"></div>
      <!-- 弹出图片容器 -->
      <div v-if="showSettingImage" class="setting-image-container">
        <img src="@/assets/img/creator.jpg" alt="设置图片" class="setting-image">
        <!-- 可选：点击图片外部关闭 -->
        <div class="setting-image-overlay" @click="toggleSettingImage()"></div>
      </div>
            <!-- 规则说明弹窗 -->
      <div v-if="showGuideModal" class="guide-modal-container">
        <img src="@/assets/img/ui/act_bg.png" alt="规则说明背景" class="guide-modal-bg">
        <div class="guide-modal-content">
          <h2>游戏规则说明</h2>
          <div class="guide-content">
            <p>1. 游戏基本玩法介绍: 核心规则是 4 人对战，，先出完所有牌者获胜。</p>
            <p>2. 卡牌规则说明: 使用一副 52 张扑克牌,每个玩家开始游戏时，会随机分配13张牌，<br/>基本牌型：单张、对子、顺子、5张同花、三带二、四带一、5张同花顺.<br/>大小规则：单张牌，点数从大到小为 2、A、K、Q、J、10…3，对子类似，5张牌:同花顺 > 四带一 > 三带二 > 同花 > 顺子 。其中牌数相同才能打
              <br/>出牌顺序：拥有最小牌方块3先出牌，之后按逆时针顺序，下家需出比上家大的牌型，或选择不出
            </p>
            <p>3. 胜利条件:任意一名玩家先出完所有牌则该获胜</p>
            <p>4. 获胜奖励：游戏结束时，结算所有玩家剩余牌总数就是获胜瓜子数</p>
            <!-- 可以根据实际游戏规则添加更多内容 -->
          </div>
          <button class="guide-close-btn" @click="toggleGuideModal()">关闭</button>
        </div>
        <div class="guide-modal-overlay" @click="toggleGuideModal()"></div>
      </div>
    </div>

    <!--  人物信息 -->
    <div class="character-info" @click="info()">
      <div class="character-bg"></div>
      <div class="avatar"></div>
      <div class="title"></div>
      <div class="nickname">
        <p style="font-size: 25px; color: white;">{{ userName }}</p>
      </div>
      <div class="star star-dark"></div>
      <div class="star star-1"></div>
      <div class="star star-2"></div>
    </div>
    <!--  立绘 -->
    <div class="character-art" id="lihui"></div>
    <!--  金币/钻石 -->
    <div class="currency-container">
      <div class="currency-bg gold-bg"></div>
      <div class="currency-icon gold-icon"></div>
      <div class="currency-bg diamond-bg"></div>
      <div class="currency-icon diamond-icon"></div>
      <p class="gold-amount">
        {{ point }}
      </p>

      <p class="diamond-amount">
        9999
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { audioManager } from '@/utils/audio'
import { storage } from '@/utils/storage'
import { useRoute, useRouter } from 'vue-router'
import { websocket } from '@/utils/websocket'

import { reconnectWebSocket } from '@/utils/websocketReconnect';

const userName = ref('')
const point = ref(0)

const router = useRouter()
const route = useRoute()
const showSettingImage = ref(false)  // 弹出图片容器显示与隐藏
const showGuideModal = ref(false)   // 规则说明弹窗显示与隐藏
// 页面可见性变化处理函数
const handleVisibilityChange = () => {
  if (document.hidden) {
    // 页面不可见时暂停音乐
    audioManager.pauseBGM()
    console.log('页面不可见，暂停音乐')
  } else {
    // 页面重新可见时恢复音乐
    if (!audioManager.isBGMPlaying()) {
      audioManager.playBGM('bgm')
      console.log('页面可见，恢复音乐')
    }
  }
}
const toggleSettingImage = () => {
  showSettingImage.value = !showSettingImage.value
}
const toggleGuideModal = () => {
  showGuideModal.value = !showGuideModal.value
}

const userIdShow = ref(false)  // 玩家di显示与隐藏
onMounted(() => {
  init()
  // 添加页面可见性变化监听
  document.addEventListener('visibilitychange', handleVisibilityChange)

   // 页面刷新时自动重连
  reconnectWebSocket(websocket);

  var userInfo = storage.local.get('user')
  userName.value = userInfo.username
  point.value = userInfo.point

  console.log(userName.value, point.value)
})
// 组件卸载时移除监听
onUnmounted(() => {
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})
const init = async () => {
  try {
    // 动态解析音频路径（替换字符串路径为 new URL() 构造的 URL）
    const bgmUrl = new URL('@/assets/music/yuanshanshaonian.mp3', import.meta.url).href;
    audioManager.preload('bgm', bgmUrl); // 使用解析后的 URL
    audioManager.playBGM('bgm')

  } catch (error) {
    console.error('初始化背景音乐失败:', error);
  }
  //当前重连
  websocket.on('open', () => {
    console.log('index on open');
    // 连接成功后再执行 toggleReady
    websocket.send({ "type": "getInfo", "data": "", "name": "" })
  });
  // 检查当前连接状态，如果已经连接，则直接执行 toggleReady
  if (websocket.ws && websocket.ws.readyState === WebSocket.OPEN) {
    console.log('index play on open');
    websocket.send({ "type": "getInfo", "data": "", "name": "" })
  }
  websocket.on('message', handleMessage)
}


const right = () => {
  // 实现right方法逻辑，例如显示/隐藏房间面板
  const panel = document.getElementById("charubiaoqian1");
  panel.style.display = panel.style.display === "block" ? "none" : "block";

  const right = document.getElementById("right");
  right.style.display = "none"
}
const return2 = () => {
  // 实现right方法逻辑，例如显示/隐藏房间面板
  const panel = document.getElementById("charubiaoqian1");
  panel.style.display = "none";

  const right = document.getElementById("right");
  right.style.display = "block"

}
const inroom = () => {
  // 实现加入房间逻辑
  const input = document.getElementById("inroom");
  const submitBtn = document.getElementById("btn_inroom");
  input.style.display = "block";
  submitBtn.style.display = "block";
}

const createRoom = () => {
  // 实现创建房间逻辑
  websocket.send({ "type": "createRoom", "data": "", "name": "" });
  router.push({ path: '/room' })
}
const joinroom = () => {
  // 实现加入房间逻辑
  const input = document.getElementById("inroom");
  const roomid = input.value
  const data = JSON.stringify({roomID:roomid})
  websocket.send({ "type": "joinRoom", "data": data, "name": "" });
  router.push({ path: '/room' })
}
const info = () => {
  userIdShow.value = !userIdShow.value
}

const handleMessage = (data) => {
  // 处理接收到的消息
  console.log('index Received message:', data)
  const parsedData = JSON.parse(data)
  if (parsedData.type == "getInfo") {
    data = JSON.parse(parsedData.data)
    //如果游戏中跳到游戏界面
    if (data.isPlaying && route.name !== "Game") {
      router.push('/game')  // 跳转到游戏页面
    }
  }
}
const toRoom = () => {
  websocket.send({ "type": "initRoom", "data": "", "name": "" });
  router.push({ path: '/room' })
}
</script>


<style scoped>
.index {
  background: url("@/assets/img/bg/beijing_chunying.png") no-repeat;
  background-size: 100% 100%;
  height: 100vh;
  /* 使用vh单位确保占满视口高度 */
  background-attachment: fixed;
  overflow: hidden;
  position: relative;
  /* 添加相对定位作为子元素定位参考 */
}

.right1 {
  background: url("@/assets/img/ui/btn_yibanchang.png");
  background-size: 100% 100%;
  height: 20%;
  width: 25%;
  top: 20%;
  position: absolute;
  right: 11%;
  float: right;
}

.right2 {
  background: url("@/assets/img/ui/btn_yourenchang.png");
  background-size: 100% 100%;
  height: 20%;
  width: 25%;
  top: 42%;
  position: absolute;
  right: 11%;
  float: right;
}

.right3 {
  background: url("@/assets/img/ui/btn_dajiangsai.png");
  background-size: 100% 100%;
  height: 20%;
  width: 25%;
  top: 61%;
  position: absolute;
  right: 11%;
  float: right;
}

.yourenchang {
  height: 200px;

  top: 15%;
  left: 60%;
  position: absolute;

  display: none;
}

.title_bg {
  height: 60px;
  width: 500px;
  margin-top: 60px;
}

.friendroom {
  z-index: 1;
  position: absolute;
  left: 40%;
  top: 35%;
  height: 40px;
}

.img_return2 {
  z-index: 1;
  position: absolute;
  left: 8%;
  top: 35%;
  height: 40px;
}

.btn_create_room {
  z-index: 1;
  position: absolute;
  left: 20%;
  top: 90%;
}

.tips {
  z-index: 1;
  position: absolute;
  left: 85%;
  top: 89%;
}

.free {
  z-index: 1;
  position: absolute;
  left: 39%;
  top: 132%;
  font-size: 25px;
  color: white;
}

.bg_abmatch1 {
  z-index: 1;
  position: absolute;
  left: 0%;
  top: 72%;
}

.bg_abmatch {
  z-index: 1;
  position: absolute;
  left: 0%;
  top: 152%;
}

.w_joinroom {
  z-index: 1;
  position: absolute;
  left: 12%;
  top: 158%;
}

.roominput {
  display: none;
  z-index: 1;
  position: absolute;
  left: 30%;
  top: 166%;
  height: 30px;
  width: 200px;
}

.roomsubmit {
  display: none;
  z-index: 1;
  position: absolute;
  left: 40%;
  top: 190%;
  height: 30px;
  width: 100px;
  background: red;
  color: white;
}

.tips2 {
  z-index: 1;
  position: absolute;
  left: 85%;
  top: 162%;
}

.free2 {
  z-index: 1;
  position: absolute;
  left: 39%;
  top: 214%;
  font-size: 25px;
  color: white;
}

.bottom-buttons .btn {
  height: 9%;
  width: 6%;
  position: absolute;
  bottom: 1%;
}

.bottom-buttons .btn-1 {
  background: url("@/assets/img/ui/btn1.png");
  background-size: 100% 100%;
  left: 55%;
}

.bottom-buttons .btn-2 {
  background: url("@/assets/img/ui/btn4.png");
  background-size: 100% 100%;
  left: 61%;
}

.bottom-buttons .btn-3 {
  background: url("@/assets/img/ui/btn0.png");
  background-size: 100% 100%;
  left: 67%;
}

.bottom-buttons .btn-4 {
  background: url("@/assets/img/ui/btn2.png");
  background-size: 100% 100%;
  left: 73%;
}

.bottom-buttons .btn-5 {
  background: url("@/assets/img/ui/btn5.png");
  background-size: 100% 100%;
  left: 79%;
}

.bottom-buttons .btn-6 {
  background: url("@/assets/img/ui/btn3.png");
  background-size: 100% 100%;
  left: 85%;
}

.bottom-buttons .btn-7 {
  background: url("@/assets/img/ui/btn6.png");
  background-size: 100% 100%;
  left: 91%;
}

.right-settings .setting-btn {
  position: absolute;
  right: 2%;
}

.right-settings .btn-set {
  background: url("@/assets/img/ui/btn_set.png");
  background-size: 100% 100%;
  height: 10%;
  width: 5.5%;
  top: 2%;
}

.right-settings .btn-guide {
  background: url("@/assets/img/ui/btn_xinshouyindao.png");
  background-size: 100% 100%;
  height: 9%;
  width: 5%;
  top: 12%;
}

.right-settings .btn-trophy {
  background: url("@/assets/img/ui/btn_trophy.png");
  background-size: 100% 100%;
  height: 9%;
  width: 5%;
  top: 22%;
}

.character-info .character-bg {
  background: url('@/assets/img/ui/bg_2.png');
  background-size: 100% 100%;
  height: 17%;
  width: 25%;
  position: absolute;
  left: 2%;
  top: 2%;
}

.character-info .avatar {
  background: url('@/assets/img/touxiang/bighead.png');
  background-size: 100% 100%;
  height: 13%;
  width: 7%;
  position: absolute;
  left: 3%;
  top: 4%;
  z-index: 1;
  border-radius: 50px;
}

.character-info .title {
  background: url('@/assets/img/ui/zuichuquesheng.png');
  background-size: 100% 100%;
  height: 5%;
  width: 15%;
  position: absolute;
  left: 10%;
  top: 5%;
  z-index: 1;
}

.character-info .nickname {
  position: absolute;
  width: 15%;
  height: 5%;
  z-index: 1;
  overflow: none;
  left: 15%;
  top: 9.5%;
}

.star {
  position: absolute;
}

.star-dark {
  background: url('@/assets/img/ui/star_dark.png');
  background-size: 100% 100%;
  height: 5%;
  width: 3%;
  left: 1%;
  top: 10%;
  z-index: 1;
}

.star-1 {
  background: url('@/assets/img/ui/star.png');
  background-size: 100% 100%;
  height: 3.5%;
  width: 2%;
  position: absolute;
  left: 3%;
  top: 14%;
  z-index: 1;
}

.star-2 {
  background: url('@/assets/img/ui/star.png');
  background-size: 100% 100%;
  height: 2.5%;
  width: 1.5%;
  position: absolute;
  left: 4.5%;
  top: 17%;
  z-index: 1;
}

#lihui {
  background: url('@/assets/img/lihui/full15920.png');
  background-size: 100% 100%;
  height: 100%;
  width: 40%;
  position: absolute;
  left: 4.5%;
  top: 16.5%;
  z-index: 1;
}

.currency-container .currency-bg {
  position: absolute;
  z-index: 1;
}

.currency-container .currency-icon {
  position: absolute;
  z-index: 1;
}

.currency-container .gold-bg {
  background: url('@/assets/img/ui/bg_0.png');
  background-size: 100% 100%;
  height: 6.5%;
  width: 11%;
  top: 2%;
  left: 50%;
}

.currency-container .gold-icon {
  background: url('@/assets/img/ui/gold0.png');
  background-size: 100% 100%;
  height: 6%;
  width: 4%;
  top: 2%;
  left: 48%;
}

.currency-container .diamond-bg {
  background: url('@/assets/img/ui/bg_0.png');
  background-size: 100% 100%;
  height: 6.5%;
  width: 11%;
  top: 2%;
  left: 63%;
}

.currency-container .diamond-icon {
  background: url('@/assets/img/ui/gold1.png');
  background-size: 100% 100%;
  height: 6%;
  width: 4%;
  top: 2%;
  left: 61%;
}

.currency-container .gold-amount {
  z-index: 1;
  position: absolute;
  left: 52.5%;
  top: 1.8%;
  font-size: 25px;
  color: white;
}

.currency-container .diamond-amount {
  z-index: 1;
  position: absolute;
  left: 65.5%;
  top: 1.8%;
  font-size: 25px;
  color: white;
}

/* 设置图片容器样式 */
.setting-image-container {
  position: absolute;
  right: 7%;
  /* 定位在设置按钮左侧 */
  top: 2%;
  /* 与设置按钮顶部对齐 */
  z-index: 1000;
  /* 确保显示在最上层 */
  margin-right: 10px;
  /* 与设置按钮保持一定距离 */
}

/* 设置图片样式 */
.setting-image {
  max-width: 200px;
  /* 设置图片最大宽度 */
  max-height: 200px;
  /* 设置图片最大高度 */
  border-radius: 8px;
  /* 可选：添加圆角 */
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  /* 可选：添加阴影效果 */
}

/* 可选：点击外部关闭的遮罩层 */
.setting-image-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: transparent;
  z-index: 999;
}

/* 规则说明弹窗样式 */
.guide-modal-container {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 2000;
  display: flex;
  justify-content: center;
  align-items: center;
}

.guide-modal-bg {
  position: absolute;
  width: 80%;
  height: 80%;
  /* object-fit: cover; */
  opacity: 0.9;
}

.guide-modal-content {
  position: relative;
  /* background: rgba(0, 0, 0, 0.4); */
  border-radius: 12px;
  padding: 30px;
  max-width: 60vw; /* 修改为屏幕宽度的60% */
  max-height: 60vh;
  overflow-y: auto;
  z-index: 10;
  color: white;
}

/* 自定义滚动条样式 */
.guide-modal-content::-webkit-scrollbar {
  width: 12px;
  /* 滚动条宽度 */
}
.guide-modal-content::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.3);
  /* 滚动条背景 */
  border-radius: 6px;
  /* 滚动条背景圆角 */
}

.guide-modal-content::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, #999, #666);
  /* 滚动条滑块灰色渐变 */
  border-radius: 6px;
  /* 滚动条滑块圆角 */
  border: 2px solid rgba(0, 0, 0, 0.2);
  /* 滚动条滑块边框 */
}

.guide-modal-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, #aaa, #777);
  /* 滚动条滑块悬停状态灰色渐变 */
}

/* Firefox 滚动条样式 */
.guide-modal-content {
  scrollbar-width: thin;
  /* 滚动条宽度 */
  scrollbar-color: #888 rgba(0, 0, 0, 0.3);
  /* 滚动条颜色（滑块颜色 背景颜色） */
}

.guide-modal-content h2 {
  text-align: center;
  color: gold;
  margin-bottom: 20px;
  font-size: 28px;
}

.guide-content {
  font-size: 18px;
  line-height: 1.8;
}

.guide-content p {
  margin-bottom: 15px;
}

.guide-close-btn {
  display: block;
  margin: 20px auto 0;
  padding: 10px 30px;
  background: gold;
  color: #333;
  border: none;
  border-radius: 20px;
  font-size: 18px;
  font-weight: bold;
  cursor: pointer;
  transition: all 0.3s;
}

.guide-close-btn:hover {
  background: #ffcc00;
  transform: scale(1.05);
}

.guide-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: transparent;
  z-index: 9;
}


/* 横屏模式专属样式 */
@media (max-width: 998px) {

  .yourenchang {
    height: 100%;
    width: 32%;
    top: 15%;
    left: 60%;
    position: absolute;

    display: none;
  }

  .title_bg {
    height: 10%;
    width: 100%;
    margin-top: 22px;
  }

  .friendroom {
    z-index: 1;
    position: absolute;
    left: 49%;
    top: 6%;
    height: 30px;
  }

  .img_return2 {
    z-index: 1;
    position: absolute;
    left: 8%;
    top: 6%;
    height: 40px;
  }

  .img_return2 img {
    width: 80%;
  }


  .btn_create_room {
    z-index: 1;
    position: absolute;
    left: 14%;
    top: 21%;
    width: 70%;
  }

  .tips {
    z-index: 2;
    position: absolute;
    left: 85%;
    top: 16%;
  }

  .tips2 {
  z-index: 1;
  position: absolute;
  left: 85%;
  top: 41%;
}

  .free {
    z-index: 1;
    position: absolute;
    left: 35%;
    top: 33%;
    font-size: 20px;
    color: white;
  }

  .free2 {
    z-index: 1;
    position: absolute;
    left: 35%;
    top: 58%;
    font-size: 20px;
    color: white;
  }

  .bg_abmatch1 {
    z-index: 1;
    position: absolute;
    left: 0%;
    top: 17%;
    width: 100%;
  }

  .bg_abmatch {
    z-index: 1;
    position: absolute;
    left: 0%;
    top: 42%;
    width: 100%;
  }

  .w_joinroom {
    z-index: 1;
    position: absolute;
    left: 6%;
    top: 42%;
    width: 90%;
  }

  .roominput {
    display: none;
    z-index: 1;
    position: absolute;
    left: 16%;
    top: 41%;
    height: 30px;
    width: 200px;
  }

  .roomsubmit {
    display: none;
    z-index: 1;
    position: absolute;
    left: 32%;
    top: 51%;
    height: 30px;
    width: 100px;
    background: red;
    color: white;
  }

  .setting-image {
    max-width: 200px;
    /* 横屏模式下缩小图片 */
    max-height: 200px;
  }

  .guide-modal-content {
    max-width: 60vw; /* 修改为屏幕宽度的60% */
    max-height: 60vh;
    margin: 20px;
    padding: 20px;
    margin-top: 50px; /* 添加顶部外边距使其向下移动 */
  }
  
  .guide-modal-content h2 {
    font-size: 24px;
  }
  
  .guide-content {
    font-size: 16px;
  }
}
</style>
