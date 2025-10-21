<template>
  <div class="loading-page">
    <div class="loading-content">
      <div class="progress-container">
        <div class="progress-bar" :style="{ width: progressPercentage + '%' }"></div>
      </div>
      <p class="loading-text">加载中... {{ progressPercentage }}%</p>
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



export default {
  name: 'LoadingPage',
  data() {
    return {
      // 所有需要预加载的图片路径（根据实际项目补充）
      imagePaths: [
        cg4Img,       // 直接使用 import 后的变量
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
      // 新增：MP3 文件路径列表
      audioPaths: [
        // gameBgMusic,       // 引入的背景音乐
        // gameBgMusic1,     // 引入的点击音效
        // ... 补充其他 MP3 路径
      ],
      loadedCount: 0,
      totalFiles: 0,       // 新增：总资源数（图片+音频）
      progressPercentage: 0,
      progressBarStyle: {
        // 确保图片完整缓存到浏览器
        backgroundSize: 'contain', // 保持图片比例完整显示
        backgroundRepeat: 'no-repeat'
      }
    };
  },
  mounted() {
    // 2. 计算总资源数（图片数 + 音频数）
    this.totalFiles = this.imagePaths.length + this.audioPaths.length;

    // 3. 预加载所有资源（图片 + 音频）
    this.preloadAllResources();
  },
  methods: {
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
        setTimeout(() => {
          this.$router.push('/index');
        }, 300);
      }
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
  background-image: url('@/assets/img/cg/cg4.png');
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
</style>
