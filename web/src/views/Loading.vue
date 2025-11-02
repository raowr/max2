<template>
  <div class="loading-page" :style="backgroundStyle">
    <div class="loading-content">
      <div class="progress-container">
        <div class="progress-bar" :style="{ width: progressPercentage + '%' }"></div>
      </div>
      <p class="loading-text">加载中... {{ progressPercentage }}%</p>
      <!-- login Button -->
      <div class="login-btn-container" v-if=isReady>
        <img src="@/assets/img/ui/login.png" class="login-btn" alt="进入主页" @click="toIndex()">
      </div>
    </div>
  </div>
</template>

<script>

// 1. 显式引入所有需要预加载的图片（确保 Vite 打包时包含这些资源）
import cg4Img from '@/assets/img/cg/cg4.png';
//index.vue 预加载资源
import titleBgImg from '@/assets/img/ui/title_bg.png';
import friendRoomImg from '@/assets/img/ui/txt_friendroom.png';
import return2Img from '@/assets/img/ui/img_return2.png';
import abMatchBgImg from '@/assets/img/ui/bg_abmatch.png';
import createRoomBtnImg from '@/assets/img/ui/btn_create_room.png';
import tipsImg from '@/assets/img/ui/tips.png';
import joinRoomImg from '@/assets/img/ui/w_joinroom.png';
import bgChunyingImg from '@/assets/img/bg/beijing_chunying.png';
import yibanBtnImg from '@/assets/img/ui/btn_yibanchang.png';
import daJiangSaiBtnImg from '@/assets/img/ui/btn_dajiangsai.png';
import yourenChangBtnImg from '@/assets/img/ui/btn_yourenchang.png';
import btn1Img from '@/assets/img/ui/btn1.png';
import btn4Img from '@/assets/img/ui/btn4.png';
import btn0Img from '@/assets/img/ui/btn0.png';
import btn2Img from '@/assets/img/ui/btn2.png';
import btn5Img from '@/assets/img/ui/btn5.png';
import btn3Img from '@/assets/img/ui/btn3.png';
import btn6Img from '@/assets/img/ui/btn6.png';
import setBtnImg from '@/assets/img/ui/btn_set.png';
import xinshouYindaoBtnImg from '@/assets/img/ui/btn_xinshouyindao.png';
import trophyBtnImg from '@/assets/img/ui/btn_trophy.png';
import bg2Img from '@/assets/img/ui/bg_2.png';
import bigheadImg from '@/assets/img/touxiang/bighead.png';
import zuichuqueshengImg from '@/assets/img/ui/zuichuquesheng.png';
import creatorImg from '@/assets/img/creator.jpg'
import actBgImg from '@/assets/img/ui/act_bg.png'

//room.vue 预加载资源
import bar23182Img from '@/assets/img/ui/bar23182.png';
import starImg from '@/assets/img/ui/star.png';
import starDarkImg from '@/assets/img/ui/star_dark.png';
import bighead15339 from '@/assets/img/touxiang/bighead15339.png';
import full15920 from '@/assets/img/lihui/full15920.png';
import full16020 from '@/assets/img/lihui/full16020.png';
import bighead15419 from '@/assets/img/touxiang/bighead15419.png';
import full15418 from '@/assets/img/lihui/full15418.png';
import indoorBgImg from '@/assets/img/bg/indoor.png';
import return1BgImg from '@/assets/img/ui/img_return1_bg.png';
import return1Img from '@/assets/img/ui/img_return1.png';
// import friendRoomImg from '@/assets/img/ui/txt_friendroom.png';
import roomOwnerImg from '@/assets/img/ui/roomowner.png';
import frameGoldImg from '@/assets/img/ui/frame_gold.png';
import characterBgImg from '@/assets/img/ui/character_bg.png';
import nameBgImg from '@/assets/img/ui/namebg.png';
import readyImg from '@/assets/img/ui/ready.png';
import btnStartNoImg from '@/assets/img/ui/btn_start_no.png';
import btnStartImg from '@/assets/img/ui/btn_start.png';
import full15703 from '@/assets/img/lihui/full15703.png';
import btnSmallImg from '@/assets/img/ui/btn_small.png';
//game.vue 预加载资源
import chatlogImg from '@/assets/img/ui/chatlog.png'; // 新增：聊天日志
import bfHeartImg from '@/assets/img/ui/bf_heart.png'; // 新增：BF 心形图标
import bighead15718 from '@/assets/img/touxiang/bighead15718.png'; // 新增：大头像15718
import img54 from '@/assets/img/54.png'; // 新增：54 图片
import bighead15729 from '@/assets/img/touxiang/bighead15729.png'; // 新增：大头像15729
import btnChupaiImg from '@/assets/img/btn_chupai.png'; // 新增：出牌按钮
import btnChupaiHuiImg from '@/assets/img/btn_chupai_hui.png'; // 新增：出牌按钮（隐藏）
import btnBujiao2Img from '@/assets/img/btn_buchu.png'; // 新增：放弃按钮
// 导入背景图片（新增代码）
import winBg from '@/assets/img/ui/win_bg.png'
import loseBg from '@/assets/img/ui/lose_bg.png'
import tableBgImg from '@/assets/img/bg/Table_Dif12324.png' // 新增：桌子背景
import gameBgMusic from '@/assets/music/game_bg.mp3';
import gameBgMusic1 from '@/assets/music/game_bg1.mp3' // 新增：游戏背景音乐1
import card1Img from '@/assets/img/cards/1.png'
import card2Img from '@/assets/img/cards/2.png'
import card3Img from '@/assets/img/cards/3.png'
import card4Img from '@/assets/img/cards/4.png'
import card5Img from '@/assets/img/cards/5.png'
import card6Img from '@/assets/img/cards/6.png'
import card7Img from '@/assets/img/cards/7.png'
import card8Img from '@/assets/img/cards/8.png'
import card9Img from '@/assets/img/cards/9.png'
import card10Img from '@/assets/img/cards/10.png'
import card11Img from '@/assets/img/cards/11.png'
import card12Img from '@/assets/img/cards/12.png'
import card13Img from '@/assets/img/cards/13.png'
import card14Img from '@/assets/img/cards/14.png'
import card15Img from '@/assets/img/cards/15.png'
import card16Img from '@/assets/img/cards/16.png'
import card17Img from '@/assets/img/cards/17.png'
import card18Img from '@/assets/img/cards/18.png'
import card19Img from '@/assets/img/cards/19.png'
import card20Img from '@/assets/img/cards/20.png'
import card21Img from '@/assets/img/cards/21.png'
import card22Img from '@/assets/img/cards/22.png'
import card23Img from '@/assets/img/cards/23.png'
import card24Img from '@/assets/img/cards/24.png'
import card25Img from '@/assets/img/cards/25.png'
import card26Img from '@/assets/img/cards/26.png'
import card27Img from '@/assets/img/cards/27.png'
import card28Img from '@/assets/img/cards/28.png'
import card29Img from '@/assets/img/cards/29.png'
import card30Img from '@/assets/img/cards/30.png'
import card31Img from '@/assets/img/cards/31.png'
import card32Img from '@/assets/img/cards/32.png'
import card33Img from '@/assets/img/cards/33.png'
import card34Img from '@/assets/img/cards/34.png'
import card35Img from '@/assets/img/cards/35.png'
import card36Img from '@/assets/img/cards/36.png'
import card37Img from '@/assets/img/cards/37.png'
import card38Img from '@/assets/img/cards/38.png'
import card39Img from '@/assets/img/cards/39.png'
import card40Img from '@/assets/img/cards/40.png'
import card41Img from '@/assets/img/cards/41.png'
import card42Img from '@/assets/img/cards/42.png'
import card43Img from '@/assets/img/cards/43.png'
import card44Img from '@/assets/img/cards/44.png'
import card45Img from '@/assets/img/cards/45.png'
import card46Img from '@/assets/img/cards/46.png'
import card47Img from '@/assets/img/cards/47.png'
import card48Img from '@/assets/img/cards/48.png'
import card49Img from '@/assets/img/cards/49.png'
import card50Img from '@/assets/img/cards/50.png'
import card51Img from '@/assets/img/cards/51.png'
import card52Img from '@/assets/img/cards/52.png'
//出牌音频
import single1 from '@/assets/music/single/1.mp3';
import single2 from '@/assets/music/single/2.mp3';
import single3 from '@/assets/music/single/3.mp3';
import single4 from '@/assets/music/single/4.mp3';
import single5 from '@/assets/music/single/5.mp3';
import single6 from '@/assets/music/single/6.mp3';
import single7 from '@/assets/music/single/7.mp3';
import single8 from '@/assets/music/single/8.mp3';
import single9 from '@/assets/music/single/9.mp3';
import single10 from '@/assets/music/single/10.mp3';
import single11 from '@/assets/music/single/11.mp3';
import single12 from '@/assets/music/single/12.mp3';
import single13 from '@/assets/music/single/13.mp3';
import single14 from '@/assets/music/single/14.mp3';
import single15 from '@/assets/music/single/15.mp3';
import single16 from '@/assets/music/single/16.mp3';
import single17 from '@/assets/music/single/17.mp3';
import single18 from '@/assets/music/single/18.mp3';
import single19 from '@/assets/music/single/19.mp3';
import single20 from '@/assets/music/single/20.mp3';
import single21 from '@/assets/music/single/21.mp3';
import single22 from '@/assets/music/single/22.mp3';
import single23 from '@/assets/music/single/23.mp3';
import single24 from '@/assets/music/single/24.mp3';
import single25 from '@/assets/music/single/25.mp3';
import single26 from '@/assets/music/single/26.mp3';
import single27 from '@/assets/music/single/27.mp3';
import single28 from '@/assets/music/single/28.mp3';
import single29 from '@/assets/music/single/29.mp3';
import single30 from '@/assets/music/single/30.mp3';
import single31 from '@/assets/music/single/31.mp3';
import single32 from '@/assets/music/single/32.mp3';
import single33 from '@/assets/music/single/33.mp3';
import single34 from '@/assets/music/single/34.mp3';
import single35 from '@/assets/music/single/35.mp3';
import single36 from '@/assets/music/single/36.mp3';
import single37 from '@/assets/music/single/37.mp3';
import single38 from '@/assets/music/single/38.mp3';
import single39 from '@/assets/music/single/39.mp3';
import single40 from '@/assets/music/single/40.mp3';
import single41 from '@/assets/music/single/41.mp3';
import single42 from '@/assets/music/single/42.mp3';
import single43 from '@/assets/music/single/43.mp3';
import single44 from '@/assets/music/single/44.mp3';
import single45 from '@/assets/music/single/45.mp3';
import single46 from '@/assets/music/single/46.mp3';
import single47 from '@/assets/music/single/47.mp3';
import single48 from '@/assets/music/single/48.mp3';
import single49 from '@/assets/music/single/49.mp3';
import single50 from '@/assets/music/single/50.mp3';
import single51 from '@/assets/music/single/51.mp3';
import single52 from '@/assets/music/single/52.mp3';
import pair3 from '@/assets/music/pair/3.mp3';
import pair4 from '@/assets/music/pair/4.mp3';
import pair5 from '@/assets/music/pair/5.mp3';
import pair6 from '@/assets/music/pair/6.mp3';
import pair7 from '@/assets/music/pair/7.mp3';
import pair8 from '@/assets/music/pair/8.mp3';
import pair9 from '@/assets/music/pair/9.mp3';
import pair10 from '@/assets/music/pair/10.mp3';
import pair11 from '@/assets/music/pair/11.mp3';
import pair12 from '@/assets/music/pair/12.mp3';
import pair13 from '@/assets/music/pair/13.mp3';
import pair14 from '@/assets/music/pair/14.mp3';
import pair15 from '@/assets/music/pair/15.mp3';
import fullHouse3 from "@/assets/music/full_house/3.mp3";
import fullHouse4 from "@/assets/music/full_house/4.mp3";
import fullHouse5 from "@/assets/music/full_house/5.mp3";
import fullHouse6 from "@/assets/music/full_house/6.mp3";
import fullHouse7 from "@/assets/music/full_house/7.mp3";
import fullHouse8 from "@/assets/music/full_house/8.mp3";
import fullHouse9 from "@/assets/music/full_house/9.mp3";
import fullHouse10 from "@/assets/music/full_house/10.mp3";
import fullHouse11 from "@/assets/music/full_house/11.mp3";
import fullHouse12 from "@/assets/music/full_house/12.mp3";
import fullHouse13 from "@/assets/music/full_house/13.mp3";
import fullHouse14 from "@/assets/music/full_house/14.mp3";
import fullHouse15 from "@/assets/music/full_house/15.mp3";
import four3 from "@/assets/music/four/3.mp3";
import four4 from "@/assets/music/four/4.mp3";
import four5 from "@/assets/music/four/5.mp3";
import four6 from "@/assets/music/four/6.mp3";
import four7 from "@/assets/music/four/7.mp3";
import four8 from "@/assets/music/four/8.mp3";
import four9 from "@/assets/music/four/9.mp3";
import four10 from "@/assets/music/four/10.mp3";
import four11 from "@/assets/music/four/11.mp3";
import four12 from "@/assets/music/four/12.mp3";
import four13 from "@/assets/music/four/13.mp3";
import four14 from "@/assets/music/four/14.mp3";
import four15 from "@/assets/music/four/15.mp3";
import straight from "@/assets/music/straight/straight.mp3";
import suit from "@/assets/music/suit/suit.mp3";
import straightFlush from "@/assets/music/straight_flush/straight_flush.mp3";
import guo from "@/assets/music/guo.mp3";
import kuaidian from "@/assets/music/kuaidian.mp3";







export default {
  name: 'LoadingPage',
  data() {
    return {
      // 所有需要预加载的图片路径（根据实际项目补充）
      imagePaths: [
        // cg4Img,       // 直接使用 import 后的变量
        titleBgImg,   // 新增：标题背景
        friendRoomImg, // 新增：好友房间
        return2Img,   // 新增：返回2图标
        abMatchBgImg,
        createRoomBtnImg,
        tipsImg,
        joinRoomImg,
        bgChunyingImg,
        yibanBtnImg,
        daJiangSaiBtnImg,
        yourenChangBtnImg,
        btn1Img,
        btn4Img,
        btn0Img,
        btn2Img,
        btn5Img,
        btn3Img,
        btn6Img,
        setBtnImg,
        xinshouYindaoBtnImg,
        trophyBtnImg,
        bg2Img,
        bigheadImg,
        zuichuqueshengImg, // 新增：最上游程
        creatorImg,   // 新增：创建者头像
        actBgImg,     // 新增：规则背景
        bar23182Img,  // 新增：进度条图片
        starImg,      // 新增：星星图标
        starDarkImg,  // 新增：暗星星图标
        bighead15339, // 新增：大头像15339
        full16020,    // 新增：满16020
        bighead15419, // 新增：大头像15419
        full15418,    // 新增：满15418
        indoorBgImg,  // 新增：室内背景
        return1BgImg, // 新增：返回1背景
        return1Img,   // 新增：返回1图标
        // friendRoomImg,// 新增：好友房间
        roomOwnerImg, // 新增：房间所有者
        frameGoldImg, // 新增：金色框架
        characterBgImg,// 新增：角色背景
        nameBgImg,    // 新增：姓名背景
        readyImg,     // 新增：准备图标
        btnStartNoImg,// 新增：开始按钮（未准备）
        btnStartImg,  // 新增：开始按钮（已准备）
        // gameBgMusic,  // 新增：游戏背景音乐
        full15703,    // 新增：满15703
        btnSmallImg,  // 新增：小按钮
        chatlogImg,   // 新增：聊天日志
        bfHeartImg,   // 新增：BF 心形图标
        bighead15718, // 新增：大头像15718
        img54,        // 新增：54 图片
        bighead15729, // 新增：大头像15729
        btnChupaiImg, // 新增：出牌按钮
        btnChupaiHuiImg, // 新增：出牌按钮（隐藏）
        btnBujiao2Img, // 新增：放弃按钮2
        winBg,        // 新增：胜利背景
        loseBg,       // 新增：失败背景
        // gameBgMusic1, // 新增：游戏背景音乐1
        tableBgImg,   // 新增：桌子背景
        card1Img,       // 新增：卡片1
        card2Img,       // 新增：卡片2
        card3Img,       // 新增：卡片3
        card4Img,       // 新增：卡片4
        card5Img,       // 新增：卡片5
        card6Img,       // 新增：卡片6
        card7Img,       // 新增：卡片7
        card8Img,       // 新增：卡片8
        card9Img,       // 新增：卡片9
        card10Img,      // 新增：卡片10
        card11Img,      // 新增：卡片11
        card12Img,      // 新增：卡片12
        card13Img,      // 新增：卡片13
        card14Img,      // 新增：卡片14
        card15Img,      // 新增：卡片15
        card16Img,      // 新增：卡片16
        card17Img,      // 新增：卡片17
        card18Img,      // 新增：卡片18
        card19Img,      // 新增：卡片19
        card20Img,      // 新增：卡片20
        card21Img,      // 新增：卡片21
        card22Img,      // 新增：卡片22
        card23Img,      // 新增：卡片23
        card24Img,      // 新增：卡片24
        card25Img,      // 新增：卡片25
        card26Img,      // 新增：卡片26
        card27Img,      // 新增：卡片27
        card28Img,      // 新增：卡片28
        card29Img,      // 新增：卡片29
        card30Img,      // 新增：卡片30
        card31Img,      // 新增：卡片31
        card32Img,      // 新增：卡片32
        card33Img,      // 新增：卡片33
        card34Img,      // 新增：卡片34
        card35Img,      // 新增：卡片35
        card36Img,      // 新增：卡片36
        card37Img,      // 新增：卡片37
        card38Img,      // 新增：卡片38
        card39Img,      // 新增：卡片39
        card40Img,      // 新增：卡片40
        card41Img,      // 新增：卡片41
        card42Img,      // 新增：卡片42
        card43Img,      // 新增：卡片43
        card44Img,      // 新增：卡片44
        card45Img,      // 新增：卡片45
        card46Img,      // 新增：卡片46
        card47Img,      // 新增：卡片47
        card48Img,      // 新增：卡片48
        card49Img,      // 新增：卡片49
        card50Img,      // 新增：卡片50
        card51Img,      // 新增：卡片51
        card52Img,      // 新增：卡片52
      ],
      audioPaths: [
          // gameBgMusic,       // 引入的背景音乐
          // gameBgMusic1,      // 引入的点击音效
          // 出牌音频
          single1, single2, single3, single4, single5, single6, single7, single8, single9, single10,
          single11, single12, single13, single14, single15, single16, single17, single18, single19, single20,
          single21, single22, single23, single24, single25, single26, single27, single28, single29, single30,
          single31, single32, single33, single34, single35, single36, single37, single38, single39, single40,
          single41, single42, single43, single44, single45, single46, single47, single48, single49, single50,
          single51, single52,
          // 对子音频
          pair3, pair4, pair5, pair6, pair7, pair8, pair9, pair10,
          pair11, pair12, pair13, pair14, pair15,
          // 葫芦音频
          fullHouse3, fullHouse4, fullHouse5, fullHouse6, fullHouse7, fullHouse8, fullHouse9, fullHouse10,
          fullHouse11, fullHouse12, fullHouse13, fullHouse14, fullHouse15,
          // 四条音频
          four3, four4, four5, four6, four7, four8, four9, four10,
          four11, four12, four13, four14, four15,
          // 其他牌型音频
          straight, suit, straightFlush,
          // 功能音频
          guo, kuaidian
        ],
      loadedCount: 0,
      totalFiles: 0,       // 新增：总资源数（图片+音频）
      progressPercentage: 0,
      progressBarStyle: {
        // 确保图片完整缓存到浏览器
        backgroundSize: 'contain', // 保持图片比例完整显示
        backgroundRepeat: 'no-repeat'
      },
      isReady: false,
      backgroundLoaded:false,
      isCapacitor: window.Capacitor && window.Capacitor.isPluginAvailable('CapacitorHttp'),
    };
  },
  computed: {
    // 新增：动态计算背景样式
    backgroundStyle() {
      return {
        backgroundImage: this.backgroundLoaded ? `url(${cg4Img})` : 'none',
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat'
      };
    }
  },
  mounted() {
// 优先预加载背景图
    this.preloadBackgroundImage().then(() => {
      // 背景加载完成后，开始加载其他资源
      this.totalFiles = this.imagePaths.length + this.audioPaths.length;
      this.preloadAllResources();
    });
  },
  methods: {
    // 优先预加载背景图片
    preloadBackgroundImage() {
      return new Promise((resolve) => {
        const img = new Image();
        img.onload = () => {
          console.log('背景图片优先加载完成');
          this.backgroundLoaded = true; // 更新加载状态
          resolve();
        };
        img.onerror = () => {
          console.error('背景图片加载失败');
          this.backgroundLoaded = false;
          resolve(); // 即使失败也继续
        };
        img.src = cg4Img; // 直接使用导入的变量，不需要 require
      });
    },
    // 新增：统一预加载图片和音频
    preloadAllResources() {
      // 预加载图片（复用原有逻辑）
      this.imagePaths.forEach(path => {
        const img = new Image();
        img.src = path;
        img.onload = () => this.updateProgress();
        img.onerror = () => {
          console.error(`图片加载失败: ${path}`);
          this.updateProgress();
        };
      });

      // 新增：预加载音频
      this.audioPaths.forEach(path => {
        const audio = new Audio();
        audio.src = path;
        // 音频加载完成事件（可使用 canplaythrough 确保可播放）
        audio.oncanplaythrough = () => this.updateProgress();
        // 音频加载失败处理
        audio.onerror = () => {
          console.error(`音频加载失败: ${path}`);
          this.updateProgress();
        };
        // 触发加载（部分浏览器需显式调用 load()）
        audio.load();
      });
    },
    // 更新进度（调整为基于总资源数计算）
    updateProgress() {
      this.loadedCount++;
      this.progressPercentage = Math.round((this.loadedCount / this.totalFiles) * 100);

      if (this.loadedCount === this.totalFiles) {
        this.isReady = true
      }
    },
    toIndex() {
      this.$router.push('/index');
    }
  }
};
</script>

<style scoped>
.loading-page {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  /* 1. 添加背景图（路径根据实际位置调整，此处假设图片在 src/assets/ 下） */
  /* background-image: url('@/assets/img/cg/cg4.png'); */
  background-size: cover;
  /* 覆盖全屏 */
  background-position: center;
  /* 居中显示 */
  background-repeat: no-repeat;
  /* 不重复 */
  display: flex;
  /* 2. 将加载内容定位到底部（默认居中，改为底部对齐） */
  align-items: flex-end;
  /* 垂直方向底部对齐 */
  padding-bottom: 80px;
  /* 距离底部 80px（可调整数值控制"上一点"的距离） */
  transition: background-image 0.3s ease;
}

.loading-content {
  width: 80%;
  max-width: 500px;
  /* 水平居中（保持原居中逻辑） */
  margin: 0 auto;
}

.progress-container {
  width: 100%;
  height: 4px;
  /* 进度条底部背景改为黑色 */
  background: #000;
  /* 原背景：#f0f0f0 */
  border-radius: 2px;
  /* overflow: hidden; */
}

.progress-bar {
  height: 100%;
  background-color: white;
  /* 进度条主体为白色 */
  position: relative;
  /* 相对定位，用于顶点图标的绝对定位 */
  transition: width 0.3s ease;
}

/* 进度条顶点图标（末端指示） */
.progress-bar::after {
  content: '';
  position: absolute;
  right: -16px;
  /* 向右偏移，使图标部分超出进度条末端（根据图标宽度调整） */
  top: 50%;
  transform: translateY(-50%);
  /* 垂直居中 */
  width: 32px;
  /* 图标宽度（根据实际图片尺寸调整） */
  height: 32px;
  /* 图标高度（根据实际图片尺寸调整） */
  background-image: url('@/assets/img/ui/bar23182.png');
  /* 顶点图标路径 */
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
  /* 强制显示完整图片（可选，解决图片被意外裁剪） */
  clip: auto;
}


.loading-text {
  margin-top: 16px;
  text-align: center;
  color: #ffffff;
  /* 保持原白色 */
  font-size: 16px;
  font-weight: bold;
  /* 添加粗体样式 */
}

.login-btn-container {
  position: absolute;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
}

.login-btn {
  cursor: pointer;
  transition: transform 0.2s;
}

.login-btn:hover {
  transform: scale(1.1);
}
</style>
