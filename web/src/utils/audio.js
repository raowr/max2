class AudioManager {
  constructor() {
    this.audioElements = new Map()
    this.currentBGM = null
  }

  preload(key, url) {
    const audio = new Audio(url)
    audio.loop = true
    this.audioElements.set(key, audio)
  }

  playBGM(key) {
    if (this.currentBGM) {
      this.currentBGM.pause()
    }
    const audio = this.audioElements.get(key)
    if (audio) {
      audio.play().catch(error => console.error('BGM播放失败:', error))
      this.currentBGM = audio
    }
  }

  playEffect(key) {
    const audio = this.audioElements.get(key)
    if (audio) {
      const clone = audio.cloneNode(true)
      clone.play()
    }
  }

      // 添加播放一次的方法
  playOnce(key) {
    const audio = this.audioElements.get(key)
    if (audio) {
      // 克隆音频元素，避免影响原始音频
      const clone = audio.cloneNode(true)
      // 确保设置为非循环播放
      clone.loop = false
      // 播放音频并处理可能的错误
      clone.play().catch(error => console.error('单次播放失败:', error))
    }
  }

    // 添加暂停背景音乐的方法
  pauseBGM() {
    if (this.currentBGM) {
      this.currentBGM.pause()
    }
  }
  
  // 添加获取当前BGM状态的方法
  isBGMPlaying() {
    return this.currentBGM && !this.currentBGM.paused
  }
}

export const audioManager = new AudioManager()
