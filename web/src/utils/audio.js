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
}

export const audioManager = new AudioManager()
