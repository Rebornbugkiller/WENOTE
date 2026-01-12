<template>
  <div class="gameboy-container">
    <div class="gameboy-body">
      <!-- Screen Area -->
      <div class="screen-bezel">
        <div class="power-led"></div>
        <div class="screen-content">
          <div class="snake-container">
            <div class="snake-grid"></div>
            <div
              v-for="(part, i) in snake"
              :key="i"
              class="snake-part"
              :class="{ 'snake-head': i === 0 }"
              :style="getPartStyle(part)"
            ></div>
            <div class="snake-food" :style="getFoodStyle()"></div>
            <div class="snake-score">SCORE: {{ score }}</div>
          </div>
        </div>
        <div class="brand-text">NINTENDO</div>
      </div>
      
      <!-- Controls -->
      <div class="controls-area">
        <div class="d-pad">
          <div class="d-pad-h"></div>
          <div class="d-pad-v"></div>
        </div>
        <div class="action-buttons">
          <div class="btn-group">
            <div class="btn btn-b"></div>
            <div class="btn-label">B</div>
          </div>
          <div class="btn-group">
            <div class="btn btn-a"></div>
            <div class="btn-label">A</div>
          </div>
        </div>
      </div>
      
      <div class="start-select">
        <div class="pill-btn"></div>
        <div class="pill-btn"></div>
      </div>
      
      <div class="speaker-grill">
        <div class="grill-slot"></div>
        <div class="grill-slot"></div>
        <div class="grill-slot"></div>
        <div class="grill-slot"></div>
        <div class="grill-slot"></div>
        <div class="grill-slot"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const ROWS = 15
const COLS = 20

const snake = ref([{ x: 5, y: 5 }, { x: 4, y: 5 }, { x: 3, y: 5 }])
const food = ref({ x: 10, y: 8 })
const score = ref(0)

let timer = null

const getPartStyle = (part) => ({
  left: `${(part.x / COLS) * 100}%`,
  top: `${(part.y / ROWS) * 100}%`,
  width: `${100 / COLS}%`,
  height: `${100 / ROWS}%`
})

const getFoodStyle = () => ({
  left: `${(food.value.x / COLS) * 100}%`,
  top: `${(food.value.y / ROWS) * 100}%`,
  width: `${100 / COLS}%`,
  height: `${100 / ROWS}%`
})

const move = () => {
  const head = snake.value[0]
  const maxLen = Math.floor(ROWS * COLS * 0.8) // 占满80%时重开

  // 蛇太长了，重开
  if (snake.value.length >= maxLen) {
    snake.value = [{ x: 5, y: 5 }, { x: 4, y: 5 }, { x: 3, y: 5 }]
    score.value = 0
    return
  }

  // 检查某个位置是否安全（不是蛇身、不出界）
  const isSafe = (x, y) => {
    if (x < 0 || x >= COLS || y < 0 || y >= ROWS) return false
    const body = snake.value.slice(0, -1)
    return !body.some(p => p.x === x && p.y === y)
  }

  // 计算从某点能到达的空格数（用于判断是否会把自己困住）
  const countReachable = (startX, startY) => {
    const visited = new Set()
    const queue = [{ x: startX, y: startY }]
    visited.add(`${startX},${startY}`)
    while (queue.length > 0) {
      const { x, y } = queue.shift()
      for (const [dx, dy] of [[1,0],[-1,0],[0,1],[0,-1]]) {
        const nx = x + dx, ny = y + dy
        const key = `${nx},${ny}`
        if (!visited.has(key) && isSafe(nx, ny)) {
          visited.add(key)
          queue.push({ x: nx, y: ny })
        }
      }
    }
    return visited.size
  }

  const directions = [
    { dx: 1, dy: 0 },
    { dx: -1, dy: 0 },
    { dx: 0, dy: 1 },
    { dx: 0, dy: -1 }
  ]

  const safeMoves = directions.filter(d => isSafe(head.x + d.dx, head.y + d.dy))

  if (safeMoves.length === 0) {
    snake.value = [{ x: 5, y: 5 }, { x: 4, y: 5 }, { x: 3, y: 5 }]
    score.value = 0
    return
  }

  // 评估每个方向：优先选择空间大的，其次选择离食物近的
  let best = safeMoves[0]
  let bestScore = -Infinity

  for (const m of safeMoves) {
    const nx = head.x + m.dx
    const ny = head.y + m.dy
    const reachable = countReachable(nx, ny)
    const distToFood = Math.abs(nx - food.value.x) + Math.abs(ny - food.value.y)
    // 空间权重高，距离权重低
    const score = reachable * 10 - distToFood
    if (score > bestScore) {
      bestScore = score
      best = m
    }
  }

  const newHead = { x: head.x + best.dx, y: head.y + best.dy }

  // 检查是否撞到自己（双重保险）
  if (snake.value.some(p => p.x === newHead.x && p.y === newHead.y)) {
    snake.value = [{ x: 5, y: 5 }, { x: 4, y: 5 }, { x: 3, y: 5 }]
    score.value = 0
    return
  }

  const newSnake = [newHead, ...snake.value]

  if (newHead.x === food.value.x && newHead.y === food.value.y) {
    score.value++
    // 生成新食物，确保不在蛇身上
    let newFood
    do {
      newFood = {
        x: Math.floor(Math.random() * COLS),
        y: Math.floor(Math.random() * ROWS)
      }
    } while (newSnake.some(p => p.x === newFood.x && p.y === newFood.y))
    food.value = newFood
  } else {
    newSnake.pop()
  }

  snake.value = newSnake
}

onMounted(() => {
  timer = setInterval(move, 50)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<style scoped>
.gameboy-container {
  width: 280px;
  height: 480px;
  background: #c0c0c0;
  border-radius: 10px 10px 40px 10px;
  box-shadow: 
    -5px 0 0 #a0a0a0,
    10px 10px 20px rgba(0,0,0,0.5),
    inset 2px 2px 5px rgba(255,255,255,0.5);
  padding: 20px;
  position: relative;
  border: 2px solid #909090;
}

.gameboy-body {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.screen-bezel {
  background: #777;
  border-radius: 10px 10px 30px 10px;
  padding: 20px 30px;
  box-shadow: inset 2px 2px 5px rgba(0,0,0,0.5);
  position: relative;
  margin-bottom: 30px;
}

.power-led {
  width: 8px;
  height: 8px;
  background: #f00;
  border-radius: 50%;
  position: absolute;
  top: 40%;
  left: 10px;
  box-shadow: 0 0 5px #f00;
  animation: pulse-led 2s infinite;
}

@keyframes pulse-led {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.screen-content {
  background: #9ca04c;
  width: 100%;
  height: 160px;
  box-shadow: inset 2px 2px 5px rgba(0,0,0,0.2);
  position: relative;
  overflow: hidden;
}

.snake-container {
  width: 100%;
  height: 100%;
  position: relative;
}

.snake-grid {
  position: absolute;
  inset: 0;
  opacity: 0.1;
  background-image: linear-gradient(#000 1px, transparent 1px),
    linear-gradient(90deg, #000 1px, transparent 1px);
  background-size: 10px 10px;
}

.snake-part {
  position: absolute;
  background: #0f380f;
  border-radius: 1px;
}

.snake-head {
  background: #0f380f;
  z-index: 10;
}

.snake-food {
  position: absolute;
  background: #0f380f;
  border-radius: 2px;
  animation: blink 0.5s infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

.snake-score {
  position: absolute;
  top: 4px;
  right: 8px;
  font-size: 10px;
  font-weight: 900;
  color: #0f380f;
  font-family: monospace;
}

.brand-text {
  color: #555;
  font-family: sans-serif;
  font-weight: 900;
  font-style: italic;
  font-size: 10px;
  margin-top: 5px;
  letter-spacing: 1px;
}

/* Controls */
.controls-area {
  display: flex;
  justify-content: space-between;
  padding: 0 10px;
  margin-bottom: 30px;
}

.d-pad {
  width: 90px;
  height: 90px;
  position: relative;
}

.d-pad-h {
  position: absolute;
  top: 30px;
  left: 0;
  width: 90px;
  height: 30px;
  background: #333;
  border-radius: 4px;
  box-shadow: 0 2px 0 #111;
}

.d-pad-v {
  position: absolute;
  top: 0;
  left: 30px;
  width: 30px;
  height: 90px;
  background: #333;
  border-radius: 4px;
  box-shadow: 0 2px 0 #111;
}

.d-pad-v::after {
  content: '';
  position: absolute;
  top: 30px;
  left: 5px;
  width: 20px;
  height: 30px;
  background: radial-gradient(circle at center, #222 2px, transparent 2px);
  background-size: 6px 6px;
  opacity: 0.5;
}

.action-buttons {
  display: flex;
  gap: 15px;
  align-items: flex-end;
  transform: rotate(-15deg);
}

.btn-group {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 5px;
}

.btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #a9163f;
  box-shadow: 2px 2px 0 #680d26;
  position: relative;
}

.btn::after {
  content: '';
  position: absolute;
  top: 5px;
  left: 5px;
  width: 10px;
  height: 10px;
  background: rgba(255,255,255,0.2);
  border-radius: 50%;
}

.btn:active {
  box-shadow: none;
  transform: translate(2px, 2px);
}

.btn-label {
  color: #000088;
  font-weight: 900;
  font-family: sans-serif;
  font-size: 12px;
}

.start-select {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-bottom: 30px;
}

.pill-btn {
  width: 40px;
  height: 12px;
  background: #999;
  border-radius: 10px;
  transform: rotate(-25deg);
  border: 1px solid #777;
  box-shadow: 1px 1px 0 #555;
}

.speaker-grill {
  position: absolute;
  bottom: 20px;
  right: 20px;
  display: flex;
  gap: 6px;
  transform: rotate(-25deg);
}

.grill-slot {
  width: 6px;
  height: 60px;
  background: rgba(0,0,0,0.1);
  border-radius: 3px;
  box-shadow: inset 1px 1px 2px rgba(0,0,0,0.2);
}
</style>
