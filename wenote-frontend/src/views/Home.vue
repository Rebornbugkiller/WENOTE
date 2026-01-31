<script setup>
import { ref, onMounted, watch, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { Search, Plus, Menu, X, FolderOpen, Trash2, LogOut, Gamepad2, Globe, BarChart3 } from 'lucide-vue-next'
import confetti from 'canvas-confetti'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import { useGamificationStore } from '../stores/gamification'
import { useNotes } from '../composables/useNotes'
import Sidebar from '../components/Sidebar.vue'
import NoteCard from '../components/notes/NoteCard.vue'
import FloatingIcons from '../components/login/FloatingIcons.vue'
import StatsDashboard from '../components/stats/StatsDashboard.vue'
import AchievementUnlockModal from '../components/gamification/AchievementUnlockModal.vue'
import { AudioEngine } from '../components/login/AudioEngine'

const router = useRouter()
const userStore = useUserStore()
const gamificationStore = useGamificationStore()
const { locale, t } = useI18n()

// Language toggle
const toggleLocale = () => {
  locale.value = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  localStorage.setItem('locale', locale.value)
}

// Game Mode State
const isGameMode = ref(true)
const isPlayingMusic = ref(true)

// Stats Panel State
const showStats = ref(false)

// Achievement Notification State
const showAchievementModal = ref(false)
const currentAchievement = ref(null)

// Watch for new achievements
watch(() => gamificationStore.pendingNotifications, (newAchievements) => {
  if (newAchievements.length > 0 && !showAchievementModal.value) {
    currentAchievement.value = newAchievements[0]
    showAchievementModal.value = true
    if (isGameMode.value) {
      AudioEngine.playSFX('achievement')
      confetti({
        particleCount: 150,
        spread: 100,
        origin: { y: 0.5 },
        colors: ['#fbbf24', '#f59e0b', '#d97706', '#22c55e', '#ec4899']
      })
    }
  }
}, { immediate: true })

const handleAchievementClose = async (achievementId) => {
  if (achievementId) {
    await gamificationStore.dismissNotification(achievementId)
  }
  showAchievementModal.value = false
  currentAchievement.value = null

  // Check if there are more achievements to show
  if (gamificationStore.pendingNotifications.length > 0) {
    setTimeout(() => {
      currentAchievement.value = gamificationStore.pendingNotifications[0]
      showAchievementModal.value = true
      if (isGameMode.value) {
        AudioEngine.playSFX('achievement')
        confetti({ particleCount: 150, spread: 100, origin: { y: 0.5 } })
      }
    }, 500)
  }
}

// Audio Helpers
const playSound = (type) => {
  if (isGameMode.value) {
    AudioEngine.playSFX(type)
  }
}

const toggleMusic = () => {
  if (isGameMode.value) {
    const playing = AudioEngine.toggleBGM()
    isPlayingMusic.value = playing
  }
}

// Use notes composable
const {
  notes,
  notebooks,
  tags,
  isLoading,
  currentView,
  searchQuery,
  fetchInitialData,
  fetchNotes,
  handleUpdateNote,
  handleDeleteNote,
  handleRestoreNote,
  handleToggleStatus,
  handleCreateNotebook,
  handleUpdateNotebook,
  handleDeleteNotebook,
  handleCreateTag,
  handleUpdateTag,
  handleDeleteTag,
  handleBatchDelete,
  handleBatchRestore,
  handleEmptyTrash,
  filterTagId,
  setFilterTag,
  setView,
  setSearch,
  // Pagination
  page,
  pageSize,
  total,
  setPage
} = useNotes()

// 批量选择
const selectedIds = ref([])
const isAllSelected = computed(() => notes.value.length > 0 && selectedIds.value.length === notes.value.length)

const toggleSelectAll = () => {
  selectedIds.value = isAllSelected.value ? [] : notes.value.map(n => n.id)
}

const toggleSelect = (id) => {
  const idx = selectedIds.value.indexOf(id)
  idx === -1 ? selectedIds.value.push(id) : selectedIds.value.splice(idx, 1)
}

const doBatchDelete = async () => {
  if (await handleBatchDelete(selectedIds.value)) {
    selectedIds.value = []
  }
}

const doBatchRestore = async () => {
  if (await handleBatchRestore(selectedIds.value)) {
    selectedIds.value = []
  }
}

// 切换视图时清空选择
watch(currentView, () => { selectedIds.value = [] })

// UI State
const sidebarOpen = ref(true)
const searchInput = ref('')

// View title
const viewTitle = computed(() => {
  if (currentView.value === 'trash') return `🗑️ ${t('home.trash')}`
  if (currentView.value === 'active') return `🎒 ${t('home.allNotes')}`
  if (currentView.value === 'starred') return `⭐ ${t('home.favorites')}`
  if (typeof currentView.value === 'number') {
    const nb = notebooks.value.find(n => n.id === currentView.value)
    if (nb) {
      const displayName = nb.is_default ? t('sidebar.uncategorized') : nb.name
      return `📒 ${displayName}`
    }
    return `📒 ${t('sidebar.notebooks')}`
  }
  return '📝 Notes'
})

// Create new note - 直接跳转到新建页面，不预先创建笔记
const createNote = () => {
  playSound('start')
  // Trigger confetti
  if (isGameMode.value) {
    confetti({
      particleCount: 100,
      spread: 70,
      origin: { y: 0.6 },
      colors: ['#22c55e', '#eab308', '#ec4899']
    })
  }
  // Navigate to new editor page
  router.push('/editor/new')
}

// Search handler with debounce
let searchTimeout = null
const handleSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    setSearch(searchInput.value)
  }, 300)
}

// Logout
const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}

// Handle note click
const handleNoteClick = (note) => {
  playSound('click')
  router.push(`/editor/${note.id}`)
}

// Initialize
onMounted(async () => {
  await userStore.fetchUser()
  await fetchInitialData()
  await fetchNotes()

  // Fetch gamification status
  await gamificationStore.fetchStatus()

  // Init audio context and auto-play BGM on first user interaction
  window.addEventListener('click', () => {
    if (isGameMode.value) {
      AudioEngine.init()
      // Auto-play BGM if isPlayingMusic is true
      if (isPlayingMusic.value && !AudioEngine.isPlaying) {
        AudioEngine.startBGM()
      }
    }
  }, { once: true })
})

// Cleanup
onUnmounted(() => {
  // 保持音乐状态，不停止BGM
})
</script>

<template>
  <div class="min-h-screen bg-green-50 font-sans overflow-hidden flex selection:bg-green-200 relative">
    <!-- Game Mode Overlay Elements -->
    <template v-if="isGameMode">
      <!-- CRT Scanlines -->
      <div class="fixed inset-0 z-50 pointer-events-none opacity-5" style="background: linear-gradient(rgba(18,16,16,0) 50%, rgba(0,0,0,0.25) 50%), linear-gradient(90deg, rgba(255,0,0,0.06), rgba(0,255,0,0.02), rgba(0,0,255,0.06)); background-size: 100% 2px, 3px 100%;"></div>
      
      <!-- Floating Background Icons -->
      <div class="fixed inset-0 pointer-events-none z-0">
        <FloatingIcons />
      </div>
    </template>

    <!-- 动态网格背景 -->
    <div class="fixed inset-0 pointer-events-none overflow-hidden">
      <div class="grid-bg absolute inset-0 opacity-30" />
      <div v-for="i in 12" :key="i" class="particle" :style="{ '--i': i }" />
    </div>

    <!-- Sidebar -->
    <Sidebar
      :notebooks="notebooks"
      :tags="tags"
      :current-view="currentView"
      :filter-tag-id="filterTagId"
      :user="userStore.user"
      :sidebar-open="sidebarOpen"
      :game-mode="isGameMode"
      @change-view="(v) => { setView(v); showStats = false }"
      @create-notebook="handleCreateNotebook"
      @update-notebook="handleUpdateNotebook"
      @delete-notebook="handleDeleteNotebook"
      @create-tag="handleCreateTag"
      @update-tag="handleUpdateTag"
      @delete-tag="handleDeleteTag"
      @filter-by-tag="(id) => { setFilterTag(id); showStats = false }"
    />

    <!-- Main Content -->
    <main class="flex-1 h-screen overflow-y-auto relative z-10">
      <!-- Header -->
      <header class="sticky top-0 z-20 p-6 flex items-center justify-between bg-green-50/80 backdrop-blur-sm">
        <div class="flex items-center gap-4 w-full max-w-2xl">
          <!-- Mobile menu toggle -->
          <button
            @click="sidebarOpen = !sidebarOpen; playSound('click')"
            class="md:hidden p-2 bg-white border-2 border-black rounded-lg active:translate-y-1 transition-transform"
          >
            <X v-if="sidebarOpen" />
            <Menu v-else />
          </button>

          <!-- Search -->
          <div class="relative w-full group">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Search class="h-5 w-5 text-slate-400 group-focus-within:text-green-500 transition-colors" />
            </div>
            <input
              v-model="searchInput"
              type="text"
              class="block w-full pl-10 pr-3 py-3 border-4 border-transparent focus:border-black rounded-xl leading-5 bg-white shadow-sm placeholder-slate-400 focus:outline-none font-bold text-lg transition-all"
              :class="{ 'focus:shadow-[4px_4px_0px_0px_rgba(34,197,94,1)]': isGameMode }"
              :placeholder="currentView === 'trash' ? t('home.searchTrash') : t('home.searchNotes')"
              @input="handleSearch"
              @keyup.enter="handleSearch"
              @focus="playSound('hover')"
            />
            <div class="absolute right-3 top-3 px-2 py-0.5 bg-slate-100 rounded text-xs font-bold text-slate-400 border border-slate-300">
              CTRL+K
            </div>
          </div>
        </div>

        <div class="flex items-center gap-4">
          <!-- Language Toggle -->
          <button
            @click="toggleLocale"
            class="px-3 py-2 flex items-center gap-2 bg-white border-2 border-black rounded-xl font-bold text-sm hover:bg-slate-100 shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] hover:shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] hover:-translate-y-0.5 transition-all"
          >
            <Globe class="w-4 h-4" />
            {{ locale === 'zh-CN' ? 'EN' : '中文' }}
          </button>

           <!-- Game Mode Toggle -->
           <button
            @click="isGameMode = !isGameMode; playSound('switch')"
            class="p-2 rounded-xl border-2 border-black transition-all hover:scale-105 active:scale-95"
            :class="isGameMode ? 'bg-green-500 text-white shadow-[2px_2px_0px_0px_#000]' : 'bg-white text-slate-400'"
            title="Toggle Game Mode"
          >
            <Gamepad2 class="w-6 h-6" />
          </button>

          <!-- Music Widget (Only in Game Mode) -->
          <div v-if="isGameMode" class="relative">
             <button
              @click="toggleMusic"
              class="p-2 rounded-xl border-2 border-black bg-black text-white hover:bg-slate-800 transition-all active:translate-y-1"
              :class="{ 'animate-pulse': isPlayingMusic }"
            >
              <span v-if="isPlayingMusic">🎵 ON</span>
              <span v-else>🔇 OFF</span>
            </button>
          </div>

          <!-- Stats -->
          <button
            @click="showStats = !showStats; playSound('click')"
            class="px-4 py-2 flex items-center gap-2 bg-white border-2 border-black rounded-xl font-bold text-slate-600 hover:bg-blue-500 hover:text-white hover:border-blue-500 shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-[3px] hover:translate-y-[3px] transition-all"
            @mouseenter="playSound('hover')"
          >
            <BarChart3 class="w-4 h-4" />
            {{ t('stats.title') }}
          </button>

          <!-- Logout -->
          <button
            @click="handleLogout"
            class="px-4 py-2 flex items-center gap-2 bg-white border-2 border-black rounded-xl font-bold text-slate-600 hover:bg-red-500 hover:text-white hover:border-red-500 shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-[3px] hover:translate-y-[3px] transition-all"
            @mouseenter="playSound('hover')"
          >
            <LogOut class="w-4 h-4" />
            {{ t('common.logout') }}
          </button>
        </div>
      </header>

      <!-- Content -->
      <div class="p-6 pt-0 max-w-7xl mx-auto">
        <!-- Stats Dashboard -->
        <div v-if="showStats" class="mb-6">
          <StatsDashboard />
        </div>
        <!-- Title & New Note Button -->
        <div v-if="!showStats" class="flex justify-between items-end mb-8">
          <div>
            <h2 class="text-4xl font-black text-slate-800 mb-1">
              {{ viewTitle }}
            </h2>
            <p class="text-slate-500 font-bold ml-1">
              {{ isLoading ? t('common.loading') : t('home.noteCount', { count: notes.length }) }}
            </p>
          </div>

          <button
            v-if="currentView !== 'trash'"
            @click="createNote"
            class="flex items-center gap-2 px-6 py-4 bg-black text-white rounded-2xl font-black text-lg shadow-[6px_6px_0px_0px_rgba(34,197,94,1)] hover:shadow-[8px_8px_0px_0px_rgba(34,197,94,1)] hover:-translate-y-1 transition-all border-2 border-transparent active:shadow-none active:translate-y-1"
          >
            <Plus class="w-6 h-6" />
            <span>{{ t('home.newNote') }}</span>
          </button>
        </div>

        <!-- 回收站批量操作栏 -->
        <div v-if="!showStats && currentView === 'trash' && notes.length > 0" class="flex items-center gap-4 mb-4 p-3 bg-white rounded-xl border-2 border-slate-200">
          <label class="flex items-center gap-2 cursor-pointer">
            <input type="checkbox" :checked="isAllSelected" @change="toggleSelectAll" class="w-5 h-5 accent-green-500" />
            <span class="font-bold text-slate-600">{{ t('home.selectAll') }}</span>
          </label>
          <span class="text-slate-400">{{ t('home.selected', { count: selectedIds.length }) }}</span>
          <template v-if="selectedIds.length > 0">
            <button
              @click="doBatchRestore"
              class="ml-auto flex items-center gap-2 px-4 py-2 bg-green-500 text-white rounded-lg font-bold hover:bg-green-600 transition-colors"
            >
              {{ t('home.batchRestore') }}
            </button>
            <button
              @click="doBatchDelete"
              class="flex items-center gap-2 px-4 py-2 bg-red-500 text-white rounded-lg font-bold hover:bg-red-600 transition-colors"
            >
              <Trash2 class="w-4 h-4" />
              {{ t('home.batchDelete') }}
            </button>
          </template>
          <button
            v-if="selectedIds.length === 0"
            @click="handleEmptyTrash"
            class="ml-auto flex items-center gap-2 px-4 py-2 bg-red-500 text-white rounded-lg font-bold hover:bg-red-600 transition-colors"
          >
            <Trash2 class="w-4 h-4" />
            {{ t('home.emptyTrash') }}
          </button>
        </div>

        <!-- Notes Grid -->
        <div v-if="!showStats" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6 pb-20">
          <NoteCard
            v-for="(note, index) in notes"
            :key="note.id"
            :note="note"
            :index="index"
            :is-trash="currentView === 'trash'"
            :selected="selectedIds.includes(note.id)"
            @click="handleNoteClick(note)"
            @delete="handleDeleteNote"
            @restore="handleRestoreNote"
            @toggle-status="handleToggleStatus"
            @toggle-select="toggleSelect"
            @mouseenter="playSound('hover')"
          />

          <!-- Empty State -->
          <div
            v-if="!isLoading && notes.length === 0"
            class="col-span-full text-center py-20 text-slate-400 font-bold text-xl opacity-50 flex flex-col items-center"
          >
            <FolderOpen class="w-16 h-16 mb-4 text-slate-300" />
            {{ currentView === 'trash' ? t('home.emptyTrashHint') : t('home.emptyNotes') }}
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="total > pageSize" class="flex justify-center mt-6 pb-6">
          <el-pagination
            :current-page="page"
            :page-size="pageSize"
            :total="total"
            layout="prev, pager, next"
            @current-change="setPage"
            background
          />
        </div>
      </div>
    </main>

    <!-- Achievement Unlock Modal -->
    <AchievementUnlockModal
      :show="showAchievementModal"
      :achievement="currentAchievement"
      @close="handleAchievementClose"
    />

    <!-- Dot Pattern Background -->
    <div
      class="fixed inset-0 z-[-1] opacity-20 pointer-events-none"
      style="background-image: radial-gradient(#22c55e 2px, transparent 2px); background-size: 30px 30px;"
    />
  </div>
</template>

<style scoped>
.grid-bg {
  background-image:
    linear-gradient(rgba(34,197,94,0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(34,197,94,0.1) 1px, transparent 1px);
  background-size: 40px 40px;
  animation: grid-scroll 20s linear infinite;
}

@keyframes grid-scroll {
  0% { transform: translate(0, 0); }
  100% { transform: translate(40px, 40px); }
}

.particle {
  position: absolute;
  width: 6px;
  height: 6px;
  background: #22c55e;
  border-radius: 50%;
  opacity: 0.4;
  animation: float 15s infinite ease-in-out;
  left: calc(var(--i) * 8%);
  top: calc(var(--i) * 7%);
  animation-delay: calc(var(--i) * -1.2s);
}

@keyframes float {
  0%, 100% { transform: translate(0, 0) scale(1); opacity: 0.4; }
  50% { transform: translate(30px, -40px) scale(1.5); opacity: 0.7; }
}

</style>
