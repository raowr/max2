<template>
  <div class="index">
    <div id="right">
      <div class="right1"></div>
      <div class="right2"></div>
      <div class="right3" @click="right()"></div>
    </div>
    <div class="yourenchang" id="charubiaoqian1">
      <img class="title_bg" src="@/assets/img/ui/title_bg.png" alt="标题背景" />
      <img class="friendroom" src="@/assets/img/ui/txt_friendroom.png" />
      <div class="img_return2" @onclick="return2()"><img src="@/assets/img/ui/img_return2.png" alt="标题背景"  /></div>
      <br /><br /><a href="room.html"><img src="@/assets/img/ui/bg_abmatch.png" class="bg_abmatch1" /></a>
      <router-link to="/room"><img src="@/assets/img/ui/btn_create_room.png" class="btn_create_room" /></router-link>
      <img class="tips" src="@/assets/img/ui/tips.png" />
      <p class="free">限时免费</p>
      <a href="#"><img src="@/assets/img/ui/bg_abmatch.png" class="bg_abmatch" @onclick="inroom()" /></a>
      <a href="#"><img src="@/assets/img/ui/w_joinroom.png" class="w_joinroom" @onclick="inroom()" /></a>
      <form action="api/room/" method="get">
        <input type="text" placeholder="房间号" name="roomid" class="roominput" id="inroom" />
        <input type="submit" value="确定" class="roomsubmit" id="btn_inroom" />
      </form>
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
      <div class="setting-btn btn-set"></div>
      <div class="setting-btn btn-guide"></div>
      <div class="setting-btn btn-trophy"></div>
    </div>

    <!--  人物信息 -->
    <div class="character-info">
      <div class="character-bg"></div>
      <div class="avatar"></div>
      <div class="title"></div>
      <div class="nickname">
        <p style="font-size: 25px; color: white;">帅哥</p>
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
        9999
      </p>

      <p class="diamond-amount">
        9999
      </p>
    </div>
  </div>
</template>

<script setup>

   const right =() =>{
      // 实现right方法逻辑，例如显示/隐藏房间面板
      const panel = document.getElementById("charubiaoqian1");
      panel.style.display = panel.style.display === "block" ? "none" : "block";

      const right = document.getElementById("right");
      right.style.display = "none"
    }
      const   return2 =() =>{
      // 实现right方法逻辑，例如显示/隐藏房间面板
      const panel = document.getElementById("charubiaoqian1");
      panel.style.display =  "none";

      const right = document.getElementById("right");
      right.style.display = "block"
      
    }
   const  inroom =() =>{
      // 实现加入房间逻辑
      const input = document.getElementById("inroom");
      const submitBtn = document.getElementById("btn_inroom");
      input.style.display = "block";
      submitBtn.style.display = "block";
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
  width: 23%;
  top: 15%;
  position: absolute;
  right: 13%;
  float: right;
}

.right2 {
  background: url("@/assets/img/ui/btn_dajiangsai.png");
  background-size: 100% 100%;
  height: 20%;
  width: 25%;
  top: 35%;
  position: absolute;
  right: 11%;
  float: right;
}

.right3 {
  background: url("@/assets/img/ui/btn_yourenchang.png");
  background-size: 100% 100%;
  height: 20%;
  width: 25%;
  top: 55%;
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
  background: url('@/assets/img/lihui/full15342.png');
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
</style>
