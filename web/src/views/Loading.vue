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
export default {
  name: 'LoadingPage',
  data() {
    return {
      allResources: [], // 存储所有资源的 URL 路径
      loadedCount: 0,
      totalResources: 0,
      progressPercentage: 0
    };
  },
  mounted() {
    // 1. 自动收集项目资源（根据实际需求调整匹配规则）
    this.collectAllResources();
    // 2. 开始预加载所有资源
    this.preloadAllResources();
  },
  methods: {
    // 自动收集资源（支持图片、音频等多种类型）
    collectAllResources() {
      // 使用 Vite glob 导入：匹配 assets 目录下所有常见资源文件
      // 格式说明：../assets/**/* 匹配 assets 下所有子目录，{png,jpg,jpeg,mp3,wav} 指定文件类型
      const resourceModules = import.meta.glob('@/assets/**/*.{png,jpg,jpeg,mp3,wav}', { 
        eager: true, // 立即导入（非懒加载）
        import: 'default' // 只获取资源 URL（默认导出）
      });

      // 提取资源 URL 到 allResources 数组
      this.allResources = Object.values(resourceModules);
      this.totalResources = this.allResources.length;
    },

    // 预加载所有资源（自动区分图片/音频类型）
    preloadAllResources() {
      this.allResources.forEach(resourceUrl => {
        // 根据文件后缀判断资源类型
        if (/\.(png|jpg|jpeg)$/.test(resourceUrl)) {
          // 图片资源：用 Image 对象预加载
          const img = new Image();
          img.src = resourceUrl;
          img.onload = () => this.updateProgress();
          img.onerror = () => {
            console.error(`图片加载失败: ${resourceUrl}`);
            this.updateProgress();
          };
        } else if (/\.(mp3|wav)$/.test(resourceUrl)) {
          // 音频资源：用 Audio 对象预加载
          const audio = new Audio();
          audio.src = resourceUrl;
          audio.oncanplaythrough = () => this.updateProgress(); // 确保音频可流畅播放
          audio.onerror = () => {
            console.error(`音频加载失败: ${resourceUrl}`);
            this.updateProgress();
          };
          audio.load(); // 显式触发加载
        }
      });
    },

    // 更新加载进度
    updateProgress() {
      this.loadedCount++;
      this.progressPercentage = Math.round((this.loadedCount / this.totalResources) * 100);
      
      if (this.loadedCount === this.totalResources) {
        setTimeout(() => {
          this.$router.push('/index'); // 全部加载完成后跳转
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
  background-size: cover; /* 覆盖全屏 */
  background-position: center; /* 居中显示 */
  background-repeat: no-repeat; /* 不重复 */
  display: flex;
  /* 2. 将加载内容定位到底部（默认居中，改为底部对齐） */
  align-items: flex-end; /* 垂直方向底部对齐 */
  padding-bottom: 80px; /* 距离底部 80px（可调整数值控制"上一点"的距离） */
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
  background: #000; /* 原背景：#f0f0f0 */
  border-radius: 2px;
  /* overflow: hidden; */
}

.progress-bar {
  height: 100%;
  background-color: white; /* 进度条主体为白色 */
  position: relative; /* 相对定位，用于顶点图标的绝对定位 */
  transition: width 0.3s ease;
}
/* 进度条顶点图标（末端指示） */
.progress-bar::after {
  content: '';
  position: absolute;
  right: -16px; /* 向右偏移，使图标部分超出进度条末端（根据图标宽度调整） */
  top: 50%;
  transform: translateY(-50%); /* 垂直居中 */
  width: 32px; /* 图标宽度（根据实际图片尺寸调整） */
  height: 32px; /* 图标高度（根据实际图片尺寸调整） */
  background-image: url('@/assets/img/ui/bar23182.png'); /* 顶点图标路径 */
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
    /* 强制显示完整图片（可选，解决图片被意外裁剪） */
  clip: auto;
}


.loading-text {
  margin-top: 16px;
  text-align: center;
  color: #ffffff; /* 保持原白色 */
  font-size: 16px;
  font-weight: bold; /* 添加粗体样式 */
}
</style>
