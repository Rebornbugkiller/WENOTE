// 8-bit 音效 & BGM 播放器
export const AudioEngine = {
  ctx: null,
  isPlaying: false,
  bgmAudio: null,

  init() {
    if (!this.ctx) {
      this.ctx = new (window.AudioContext || window.webkitAudioContext)()
    }
    if (this.ctx.state === 'suspended') {
      this.ctx.resume()
    }
  },

  playSFX(type) {
    this.init()
    const ctx = this.ctx
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain)
    gain.connect(ctx.destination)
    const now = ctx.currentTime

    if (type === 'type') {
      osc.type = 'square'
      osc.frequency.setValueAtTime(800 + Math.random() * 200, now)
      osc.frequency.exponentialRampToValueAtTime(100, now + 0.05)
      gain.gain.setValueAtTime(0.03, now)
      gain.gain.exponentialRampToValueAtTime(0.001, now + 0.05)
      osc.start(now)
      osc.stop(now + 0.05)
    } else if (type === 'hover') {
      osc.type = 'triangle'
      osc.frequency.setValueAtTime(400, now)
      gain.gain.setValueAtTime(0.02, now)
      gain.gain.linearRampToValueAtTime(0, now + 0.03)
      osc.start(now)
      osc.stop(now + 0.03)
    } else if (type === 'start') {
      osc.type = 'sawtooth'
      osc.frequency.setValueAtTime(220, now)
      osc.frequency.linearRampToValueAtTime(880, now + 0.4)
      gain.gain.setValueAtTime(0.05, now)
      gain.gain.exponentialRampToValueAtTime(0.001, now + 0.4)
      osc.start(now)
      osc.stop(now + 0.4)
    } else if (type === 'win') {
      const notes = [523.25, 659.25, 783.99, 1046.50]
      notes.forEach((freq, i) => {
        const o = ctx.createOscillator()
        const g = ctx.createGain()
        o.type = 'square'
        o.frequency.value = freq
        o.connect(g)
        g.connect(ctx.destination)
        const startTime = now + i * 0.1
        g.gain.setValueAtTime(0.05, startTime)
        g.gain.exponentialRampToValueAtTime(0.001, startTime + 0.4)
        o.start(startTime)
        o.stop(startTime + 0.4)
      })
    } else if (type === 'error') {
      osc.type = 'sawtooth'
      osc.frequency.setValueAtTime(200, now)
      osc.frequency.exponentialRampToValueAtTime(50, now + 0.3)
      gain.gain.setValueAtTime(0.1, now)
      gain.gain.exponentialRampToValueAtTime(0.001, now + 0.3)
      osc.start(now)
      osc.stop(now + 0.3)
    } else if (type === 'success') {
      osc.type = 'sine'
      osc.frequency.setValueAtTime(600, now)
      osc.frequency.linearRampToValueAtTime(900, now + 0.1)
      gain.gain.setValueAtTime(0.04, now)
      gain.gain.exponentialRampToValueAtTime(0.001, now + 0.15)
      osc.start(now)
      osc.stop(now + 0.15)
    } else if (type === 'firework') {
      // 胜利fanfare音效 - 更响亮更明显
      const notes = [523, 659, 784, 1047, 784, 1047]
      notes.forEach((freq, i) => {
        const o = ctx.createOscillator()
        const g = ctx.createGain()
        o.connect(g)
        g.connect(ctx.destination)
        o.type = 'square'
        const t = now + i * 0.12
        o.frequency.setValueAtTime(freq, t)
        g.gain.setValueAtTime(0.15, t)
        g.gain.exponentialRampToValueAtTime(0.01, t + 0.25)
        o.start(t)
        o.stop(t + 0.25)
      })
    } else if (type === 'achievement') {
      // 成就解锁音效 - 胜利号角
      const notes = [523, 659, 784, 1047, 1319, 1568]
      notes.forEach((freq, i) => {
        const o = ctx.createOscillator()
        const g = ctx.createGain()
        o.connect(g)
        g.connect(ctx.destination)
        o.type = 'square'
        const t = now + i * 0.1
        o.frequency.setValueAtTime(freq, t)
        g.gain.setValueAtTime(0.1, t)
        g.gain.exponentialRampToValueAtTime(0.01, t + 0.3)
        o.start(t)
        o.stop(t + 0.3)
      })
    } else if (type === 'streak') {
      // 连续天数音效 - 火焰音效
      osc.type = 'sawtooth'
      osc.frequency.setValueAtTime(100, now)
      osc.frequency.exponentialRampToValueAtTime(400, now + 0.1)
      osc.frequency.exponentialRampToValueAtTime(100, now + 0.2)
      gain.gain.setValueAtTime(0.05, now)
      gain.gain.exponentialRampToValueAtTime(0.001, now + 0.2)
      osc.start(now)
      osc.stop(now + 0.2)
    } else if (type === 'goalComplete') {
      // 目标完成音效 - 庆祝铃声
      const notes = [784, 988, 1175, 1568]
      notes.forEach((freq, i) => {
        const o = ctx.createOscillator()
        const g = ctx.createGain()
        o.connect(g)
        g.connect(ctx.destination)
        o.type = 'sine'
        const t = now + i * 0.15
        o.frequency.setValueAtTime(freq, t)
        g.gain.setValueAtTime(0.08, t)
        g.gain.exponentialRampToValueAtTime(0.001, t + 0.4)
        o.start(t)
        o.stop(t + 0.4)
      })
    } else if (type === 'switch') {
      // 切换音效
      osc.type = 'square'
      osc.frequency.setValueAtTime(440, now)
      osc.frequency.setValueAtTime(550, now + 0.05)
      gain.gain.setValueAtTime(0.03, now)
      gain.gain.exponentialRampToValueAtTime(0.001, now + 0.1)
      osc.start(now)
      osc.stop(now + 0.1)
    }
  },

  toggleBGM() {
    this.init()
    if (this.isPlaying) {
      this.stopBGM()
    } else {
      this.startBGM()
    }
    return this.isPlaying
  },

  startBGM() {
    if (this.isPlaying) return
    if (!this.bgmAudio) {
      this.bgmAudio = new Audio('/bgm.mp3')
      this.bgmAudio.loop = true
      this.bgmAudio.volume = 0.3
    }
    this.isPlaying = true
    const playPromise = this.bgmAudio.play()
    if (playPromise !== undefined) {
      playPromise.catch((error) => {
        console.log('BGM autoplay was prevented:', error)
        this.isPlaying = false
      })
    }
  },

  stopBGM() {
    if (this.bgmAudio) {
      this.bgmAudio.pause()
    }
    this.isPlaying = false
  },

  getUserMusicPreference() {
    const pref = localStorage.getItem('wenote_music_enabled')
    return pref === null ? true : pref === 'true'
  },

  setUserMusicPreference(enabled) {
    localStorage.setItem('wenote_music_enabled', enabled.toString())
  }
}
