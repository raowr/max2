<template>
  <div class="loading-page" :style="backgroundStyle">
    <div class="loading-content">
      <p class="loading-text">加载中... {{ progressPercentage }}%</p>
      <div class="progress-container">
        <div class="progress-bar" :style="{ width: progressPercentage + '%' }"></div>
      </div>
      <!-- login Button -->
      <div class="login-btn-container" v-if="isReady">
          <div class="login-group">
              <img src="@/assets/img/ui/guest_login.png" class="login-btn" alt="游客进入" @click="toIndex()">
          </div>
          <div class="login-group">
              <img src="@/assets/img/ui/user_login.png" class="login-btn" alt="登录游戏" @click="handleLoginButtonClick()">
              <span class="register-text" @click="toRegister()">注册账号</span>
          </div>
      </div>
    </div>
     <!-- 登录弹窗 -->
    <div class="modal-overlay" v-if="showLoginModal" @click="closeLoginModal">
      <div class="modal-content" @click.stop>
        <h3 class="modal-title">登录游戏</h3>
        <form class="login-form" @submit.prevent="handleLogin">
          <div class="form-group">
            <label for="username">账号名</label>
            <input type="text" id="username" v-model="loginForm.username" placeholder="请输入账号名" required>
          </div>
          <div class="form-group">
            <label for="password">密码</label>
            <input type="password" id="password" v-model="loginForm.password" placeholder="请输入密码" required>
          </div>
          <div class="form-actions">
            <button type="button" class="btn-cancel" @click="closeLoginModal">取消</button>
            <button type="submit" class="btn-confirm">确定</button>
          </div>
        </form>
      </div>
    </div>
        <!-- 注册弹窗 -->
    <div class="modal-overlay" v-if="showRegisterModal" @click="closeRegisterModal">
      <div class="modal-content" @click.stop>
        <h3 class="modal-title">注册账号</h3>
        <form class="login-form" @submit.prevent="handleRegister">
          <div class="form-group">
            <label for="register-username">账号名</label>
            <input type="text" id="register-username" v-model="registerForm.username" placeholder="请输入账号名" required>
          </div>
          <div class="form-group">
            <label for="register-password">密码</label>
            <input type="password" id="register-password" v-model="registerForm.password" placeholder="请输入密码" required>
          </div>
          <div class="form-group">
            <label for="register-confirm-password">确认密码</label>
            <input type="password" id="register-confirm-password" v-model="registerForm.confirmPassword" placeholder="请再次输入密码" required>
          </div>
          <div class="form-actions">
            <button type="button" class="btn-cancel" @click="closeRegisterModal">取消</button>
            <button type="submit" class="btn-confirm">确定</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { login, register } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus';
import { storage } from '@/utils/storage'
import { websocket } from '@/utils/websocket'


import { audioManager } from '@/utils/audio'

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






export default {
  name: 'LoadingPage',
  data() {
    return {
      // 登录弹窗状态
      showLoginModal: false,
      loginForm: {
        username: '',
        password: ''
      },
      // 注册弹窗状态
      showRegisterModal: false,
      registerForm: {
        username: '',
        password: '',
        confirmPassword: ''
      },
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
        // 延迟连接 WebSocket，确保应用先加载
        let wsUrl = import.meta.env.VITE_WS_URL;
      setTimeout(() => {
        if (wsUrl) {
          console.log('user_id:', storage.local.get("user_id"))
    // 初始化 WebSocket 配置（全局一次）
          websocket.init({
            url: wsUrl+"?user_id="+storage.local.get("user_id"), // 后端地址
            reconnectInterval: 5000, // 5秒重连一次
            heartbeatInterval: 10000 // 20秒一次心跳
          });

    // 启动连接（可在登录后再调用，这里直接启动作为示例）
          websocket.connect();
        } else {
          console.warn('VITE_WS_URL is not defined, using default WebSocket URL')
          // 可以在这里设置一个默认的 WebSocket URL
          websocket.connect('ws://127.0.0.1:8000/enter')
        }
      }, 100)
      this.$router.push('/index');
    },
    // 登录按钮点击处理
    handleLoginButtonClick() {
      const userData = storage.local.get('user');
      
      if (!userData) {
        // 如果没有用户信息，弹出登录框
        this.showLoginModal = true;
        return;
      }
      // 如果有用户信息，尝试连接 websocket
      const node = userData.node || userData.Node;
      if (!node) {
        ElMessage.error('用户信息不完整，请重新登录');
        this.showLoginModal = true;
        return;
      }
      // 尝试连接 websocket
      this.tryConnectWebSocket(userData);
    },
       // 尝试连接 WebSocket
    // 尝试连接 WebSocket（使用全局 websocket 实例）
    tryConnectWebSocket(userData) {
        // 从用户数据中提取参数
        const node = userData.node || userData.Node;
        const username = userData.username || userData.Username;
        const token = userData.token || userData.Token;
        
        // 构建 WebSocket URL
        const wsUrl = `ws://${node}/enter?user_name=${encodeURIComponent(username)}&token=${encodeURIComponent(token)}`;
        
        console.log('尝试连接 WebSocket:', wsUrl);

        // 使用全局 websocket 实例
        const ws = this.$websocket;

        // 保存回调引用，用于后续移除
        const openCallback = () => {
            console.log('WebSocket 连接成功');
            // 移除回调避免重复触发
            ws.off('open', openCallback);
            ws.off('error', errorCallback);
            ws.off('close', closeCallback);
            // 连接成功，跳转首页
            this.$router.push('/index');
        };

        const errorCallback = (error) => {
            console.error('WebSocket 连接错误:', error);
            // 移除回调
            ws.off('open', openCallback);
            ws.off('error', errorCallback);
            ws.off('close', closeCallback);
            ElMessage.error('登录已过期，请重新登录');
            this.showLoginModal = true;
        };

        const closeCallback = (event) => {
            console.log('WebSocket 连接关闭:', event.code, event.reason);
            // 移除回调
            ws.off('open', openCallback);
            ws.off('error', errorCallback);
            ws.off('close', closeCallback);
            // 如果不是正常关闭且还在 loading 页面
            if (event.code !== 1000 && this.$route.path === '/loading') {
                ElMessage.error('服务器连接异常，请重新登录');
                this.showLoginModal = true;
            }
        };

        // 注册回调
        ws.on('open', openCallback);
        ws.on('error', errorCallback);
        ws.on('close', closeCallback);

        // 初始化并连接
        ws.init({
            url: wsUrl,
            reconnectInterval: 5000,
            heartbeatInterval: 10000
        }).connect();

        // 超时处理
        setTimeout(() => {
            // 检查是否还在连接中
            if (ws.isConnecting && this.$route.path === '/loading') {
                ws.close();
                ElMessage.error('连接超时，请重新登录');
                this.showLoginModal = true;
            }
        }, 10000);
    },
    // 关闭登录弹窗
    closeLoginModal() {
      this.showLoginModal = false;
      this.loginForm.username = '';
      this.loginForm.password = '';
    },
    // 处理登录
    handleLogin() {
      console.log('登录信息:', this.loginForm);
      // 这里可以添加登录逻辑
      // 登录成功后关闭弹窗并跳转到游戏页面
      // alert('登录成功！');
      login(this.loginForm).then((res)=>{
        console.log(res.data)
        if(res.code === 0){
          ElMessage.success('登录成功！');
          //保存用户信息
          storage.local.set('user', res.data);
                    // 尝试连接 WebSocket
          const userData = res.data;
          if (userData) {
            this.closeLoginModal();
            this.tryConnectWebSocket(userData);
          } else {
            ElMessage.error('服务器信息获取失败');
          }
        }else{
          ElMessage.error(res.message);
        }
      })
      this.closeLoginModal();
      // this.$router.push('/index');
    },
        // 关闭注册弹窗
    closeRegisterModal() {
      this.showRegisterModal = false;
      this.registerForm.username = '';
      this.registerForm.password = '';
      this.registerForm.confirmPassword = '';
    },
    // 跳转到注册（显示注册弹窗）
    toRegister() {
      this.showRegisterModal = true;
    },
    // 处理注册
    handleRegister() {
      // 验证密码是否一致
      if (this.registerForm.password !== this.registerForm.confirmPassword) {
        ElMessage.error('两次输入的密码不一致！');
        return;
      }
      console.log('注册信息:', this.registerForm);
      // 这里可以添加注册逻辑
      // alert('注册成功！');
      let registerData = {
        username: this.registerForm.username,
        password: this.registerForm.password,
        password2: this.registerForm.confirmPassword,
      }
      register(registerData).then((res)=>{
        console.log(res.data)
        if(res.code === 0){
          ElMessage.success('注册成功！');
        }else{
          ElMessage.error(res.message);
        }
      })
      this.closeRegisterModal();
      // 注册成功后可以跳转到登录页面或直接登录
    },
    toLogin() {
    // 添加登录逻辑或跳转到登录页面
    // this.$router.push('/login');
    console.log('登录游戏');
    },
       
// ... existing code ...
    async preloadAllAudios() {
      console.log('开始预加载音频资源...');
      const keys = Object.keys(this.audioResources);
      let loadedCount = 0;
      
      // 分批加载，避免一次性请求过多资源
      const batchSize = 5;
      for (let i = 0; i < keys.length; i += batchSize) {
        const batch = keys.slice(i, i + batchSize);
        const promises = batch.map(musicPath => {
          return new Promise((resolve) => {
            const key = 'preload_' + musicPath;
            try {
              // 修改为：直接使用 audioResources 中导入的音频变量
              // 这些导入的变量会自动解析为编译后的带 hash 的路径
              let cardMusicUrl = this.audioResources[musicPath];
              
              // 如果是 true，表示需要处理（应该避免这种情况）
              if (cardMusicUrl === true) {
                console.warn(`音频路径未正确配置: ${musicPath}，请添加正确的导入语句`);
                resolve();
                return;
              }
              
              // 调试信息，查看实际使用的音频路径
              console.log(`预加载音频: ${musicPath}，URL: ${cardMusicUrl}`);
              
              // 存储预加载的 URL，用于后续检查
              this.preloadedAudioUrls[key] = cardMusicUrl;
              
              // 创建一个临时的 Audio 对象进行预加载
              const tempAudio = new Audio();
              
              // 设置加载完成事件
              tempAudio.oncanplaythrough = () => {
                loadedCount++;
                this.preloadStatus.loaded = loadedCount;
                console.log(`预加载完成: ${musicPath}, 进度: ${loadedCount}/${keys.length}`);
                // 预加载完成后，使用 audioManager 正式加载
                audioManager.preload(key, cardMusicUrl);
                resolve();
              };
              
              // 加载失败处理
              tempAudio.onerror = () => {
                console.error(`预加载失败: ${musicPath}，URL: ${cardMusicUrl}`);
                // 失败也继续，避免阻塞其他资源加载
                loadedCount++;
                this.preloadStatus.loaded = loadedCount;
                resolve();
              };
              
              // 设置音频源并开始加载
              tempAudio.src = cardMusicUrl;
              tempAudio.load();
            } catch (error) {
              console.error(`预加载出错: ${musicPath}`);
              loadedCount++;
              this.preloadStatus.loaded = loadedCount;
              resolve();
            }
          });
        });
        
        // 等待当前批次加载完成
        await Promise.all(promises);
      }
      
      this.preloadStatus.isComplete = true;
      console.log('所有音频资源预加载完成');
    },
// ... existing code ...
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
  padding-bottom: 150px;
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
  bottom: 24px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 10;
  display: flex;
  gap: 20px; /* 按钮之间的间隔 */
}

.login-btn {
  cursor: pointer;
  transition: transform 0.2s;
}

.login-btn:hover {
  transform: scale(1.1);
}
.login-group {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
}

.register-text {
    color: #ffffff;
    font-size: 14px;
    cursor: pointer;
    text-decoration: underline;
    transition: color 0.2s;
}

.register-text:hover {
    color: #ffff00;
}
 
/* 弹窗遮罩层 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}
 
/* 弹窗内容 */
.modal-content {
  background-color: #333;
  border-radius: 12px;
  padding: 30px;
  width: 80%;
  max-width: 400px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
}
 
/* 弹窗标题 */
.modal-title {
  color: #ffffff;
  text-align: center;
  margin-bottom: 20px;
  font-size: 20px;
}
 
/* 登录表单 */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 15px;
}
 
/* 表单组 */
.form-group {
  display: flex;
  flex-direction: column;
  gap: 5px;
}
 
.form-group label {
  color: #ffffff;
  font-size: 14px;
}
 
.form-group input {
  padding: 12px;
  border: 1px solid #555;
  border-radius: 6px;
  background-color: #222;
  color: #ffffff;
  font-size: 14px;
  outline: none;
  transition: border-color 0.2s;
}
 
.form-group input:focus {
  border-color: #4a90d9;
}
 
.form-group input::placeholder {
  color: #666;
}
 
/* 表单按钮组 */
.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 20px;
}
 
/* 按钮样式 */
.btn-cancel,
.btn-confirm {
  flex: 1;
  padding: 12px;
  border: none;
  border-radius: 6px;
  font-size: 16px;
  cursor: pointer;
  transition: background-color 0.2s;
}
 
.btn-cancel {
  background-color: #555;
  color: #ffffff;
}
 
.btn-cancel:hover {
  background-color: #666;
}
 
.btn-confirm {
  background-color: #4a90d9;
  color: #ffffff;
}
 
.btn-confirm:hover {
  background-color: #3a80c9;
}
</style>
