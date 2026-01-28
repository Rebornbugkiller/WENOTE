<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { LayoutGrid, Star, Book, Trash2, Plus, X, Settings, Edit2 } from 'lucide-vue-next'
import MenuButton from './notes/MenuButton.vue'
import CreateModal from './notes/CreateModal.vue'
import EditModal from './notes/EditModal.vue'
import { ElMessageBox, ElDropdown, ElDropdownMenu, ElDropdownItem } from 'element-plus'
import { AudioEngine } from './login/AudioEngine'
import { tagColors, getTagColor } from '../utils/color'
import { AVATAR_STYLES } from '../stores/user'

const router = useRouter()
const { t } = useI18n()

const props = defineProps({
  notebooks: { type: Array, default: () => [] },
  tags: { type: Array, default: () => [] },
  currentView: { type: [String, Number], default: 'active' },
  filterTagId: { type: Number, default: null },
  user: { type: Object, default: () => ({ username: 'Guest' }) },
  sidebarOpen: { type: Boolean, default: true },
  gameMode: { type: Boolean, default: false }
})

const playSound = (type) => {
  if (props.gameMode) {
    AudioEngine.playSFX(type)
  }
}

const emit = defineEmits(['change-view', 'create-notebook', 'update-notebook', 'delete-notebook', 'create-tag', 'update-tag', 'delete-tag', 'filter-by-tag', 'toggle-sidebar'])

// Modal state
const showNotebookModal = ref(false)
const showTagModal = ref(false)
const showEditNotebookModal = ref(false)
const showEditTagModal = ref(false)
const editingNotebook = ref(null)
const editingTag = ref(null)

const randomColor = () => tagColors[Math.floor(Math.random() * tagColors.length)]

// Handle create
const handleCreate = (name, type) => {
  if (type === 'notebook') {
    emit('create-notebook', name)
    showNotebookModal.value = false
  } else {
    emit('create-tag', name, randomColor())
    showTagModal.value = false
  }
}

const goToSettings = () => {
  playSound('click')
  router.push('/settings')
}

// Handle notebook actions
const handleNotebookAction = (command, notebook) => {
  if (command === 'rename') {
    editingNotebook.value = notebook
    showEditNotebookModal.value = true
  } else if (command === 'delete') {
    ElMessageBox.confirm(
      t('sidebar.deleteNotebookConfirm', { name: notebook.name }),
      t('sidebar.deleteNotebook'),
      { type: 'warning', confirmButtonText: t('common.delete'), cancelButtonText: t('common.cancel') }
    ).then(() => emit('delete-notebook', notebook.id))
      .catch(() => {})
  }
}

// Handle tag actions
const handleTagAction = (command, tag) => {
  if (command === 'edit') {
    editingTag.value = tag
    showEditTagModal.value = true
  } else if (command === 'delete') {
    ElMessageBox.confirm(
      t('sidebar.deleteTagConfirm'),
      t('sidebar.deleteTag'),
      { type: 'warning', confirmButtonText: t('common.delete'), cancelButtonText: t('common.cancel') }
    ).then(() => emit('delete-tag', tag.id))
      .catch(() => {})
  }
}

// Handle save notebook
const handleSaveNotebook = (data) => {
  emit('update-notebook', editingNotebook.value.id, data.name)
  showEditNotebookModal.value = false
  editingNotebook.value = null
}

// Handle save tag
const handleSaveTag = (data) => {
  emit('update-tag', editingTag.value.id, data.name, data.color)
  showEditTagModal.value = false
  editingTag.value = null
}
</script>

<template>
  <!-- Sidebar -->
  <aside
    class="fixed md:relative z-20 w-64 h-screen bg-white border-r-4 border-slate-900 shadow-2xl flex flex-col shrink-0 transition-transform duration-300"
    :class="sidebarOpen ? 'translate-x-0' : '-translate-x-[240px]'"
  >
    <!-- Logo -->
    <div class="p-6 border-b-4 border-slate-100 bg-slate-50">
      <div
        class="flex items-center gap-2 cursor-pointer hover:scale-105 transition-transform"
        @click="emit('change-view', 'active'); playSound('click')"
        @mouseenter="playSound('hover')"
      >
        <svg class="w-10 h-10" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
          <rect x="15" y="10" width="70" height="80" rx="10" fill="#00b894" stroke="#2d3436" stroke-width="4"/>
          <path d="M15 25 L85 25" stroke="#2d3436" stroke-width="4"/>
          <rect x="25" y="25" width="50" height="40" rx="4" fill="#ffffff" stroke="#2d3436" stroke-width="4"/>
          <line x1="32" y1="35" x2="68" y2="35" stroke="#dfe6e9" stroke-width="3" stroke-linecap="round"/>
          <line x1="32" y1="45" x2="68" y2="45" stroke="#dfe6e9" stroke-width="3" stroke-linecap="round"/>
          <line x1="32" y1="55" x2="55" y2="55" stroke="#dfe6e9" stroke-width="3" stroke-linecap="round"/>
          <path d="M35 42 L48 55 L75 28" fill="none" stroke="#00b894" stroke-width="6" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M30 75 h10 v-5 h5 v5 h5 v5 h-5 v5 h-5 v-5 h-10 z" fill="#2d3436"/>
          <circle cx="70" cy="78" r="4" fill="#d63031" stroke="#2d3436" stroke-width="2"/>
          <g transform="translate(75, 10) rotate(15)">
            <rect x="0" y="0" width="12" height="40" rx="2" fill="#ffeaa7" stroke="#2d3436" stroke-width="3"/>
            <path d="M0 40 L6 50 L12 40 Z" fill="#2d3436"/>
            <rect x="0" y="-5" width="12" height="8" rx="2" fill="#fab1a0" stroke="#2d3436" stroke-width="3"/>
          </g>
        </svg>
        <h1 class="text-2xl font-black tracking-tighter text-slate-800">
          WE<span class="text-green-500">NOTE</span>
        </h1>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 p-4 space-y-6 overflow-y-auto">
      <!-- Main Menu -->
      <div class="space-y-2">
        <MenuButton
          :icon="LayoutGrid"
          :label="t('home.allNotes')"
          :active="currentView === 'active'"
          @click="emit('change-view', 'active'); playSound('click')"
          @mouseenter="playSound('hover')"
        />
        <MenuButton
          :icon="Star"
          :label="t('home.favorites')"
          :active="currentView === 'starred'"
          @click="emit('change-view', 'starred'); playSound('click')"
          @mouseenter="playSound('hover')"
        />
      </div>

      <!-- Notebooks -->
      <div>
        <h3 class="px-4 text-xs font-black text-slate-400 uppercase tracking-widest mb-2 flex justify-between items-center group">
          {{ t('sidebar.notebooks') }}
          <Plus
            class="w-4 h-4 cursor-pointer hover:text-green-500 transition-colors hover:scale-125"
            @click="showNotebookModal = true"
          />
        </h3>
        <div class="space-y-1">
          <div
            v-for="nb in notebooks"
            :key="nb.id"
            class="flex items-center pr-2 group"
          >
            <MenuButton
              :icon="Book"
              :label="nb.name"
              :count="nb.note_count"
              :active="currentView === nb.id"
              class="flex-1"
              @click="emit('change-view', nb.id); playSound('click')"
              @mouseenter="playSound('hover')"
            />
            <!-- Settings dropdown -->
            <el-dropdown
              trigger="click"
              @command="(cmd) => handleNotebookAction(cmd, nb)"
            >
              <button
                class="ml-1 p-0.5 rounded text-slate-400 hover:text-slate-600 transition-all opacity-0 group-hover:opacity-100"
                @click.stop
              >
                <Settings class="w-3 h-3" />
              </button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="rename">
                    <Edit2 class="w-4 h-4 inline mr-2" />
                    {{ t('sidebar.renameNotebook') }}
                  </el-dropdown-item>
                  <el-dropdown-item v-if="!nb.is_default" command="delete" class="text-red-500">
                    <Trash2 class="w-4 h-4 inline mr-2" />
                    {{ t('sidebar.deleteNotebook') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <div
            v-if="notebooks.length === 0"
            class="px-4 text-xs text-slate-400 italic"
          >
            {{ t('sidebar.noNotebooks') }}
          </div>
        </div>
      </div>

      <!-- Tags -->
      <div>
        <h3 class="px-4 text-xs font-black text-slate-400 uppercase tracking-widest mb-2 flex justify-between items-center group">
          {{ t('sidebar.tags') }}
          <Plus
            class="w-4 h-4 cursor-pointer hover:text-green-500 transition-colors hover:scale-125"
            @click="showTagModal = true"
          />
        </h3>
        <div class="px-4 flex flex-wrap gap-2">
          <span
            v-for="tag in tags"
            :key="tag.id"
            class="text-xs font-bold px-2 py-1 rounded-md border cursor-pointer transition-all flex items-center gap-1 group relative"
            :class="filterTagId === tag.id ? 'bg-black text-white border-black' : 'bg-slate-50 text-slate-700 border-slate-200 hover:border-black'"
          >
            <span
              class="w-2 h-2 rounded-full border border-black/20"
              :style="{ backgroundColor: getTagColor(tag) }"
              @click="emit('filter-by-tag', tag.id); playSound('click')"
              @mouseenter="playSound('hover')"
            />
            <span
              @click="emit('filter-by-tag', tag.id); playSound('click')"
              @mouseenter="playSound('hover')"
            >
              {{ tag.name }}
            </span>
            <!-- Settings dropdown -->
            <el-dropdown
              trigger="click"
              @command="(cmd) => handleTagAction(cmd, tag)"
            >
              <button
                class="ml-1 opacity-0 group-hover:opacity-100 transition-opacity"
                @click.stop
              >
                <Settings class="w-3 h-3 text-slate-400 hover:text-slate-600" />
              </button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">
                    <Edit2 class="w-4 h-4 inline mr-2" />
                    {{ t('sidebar.editTag') }}
                  </el-dropdown-item>
                  <el-dropdown-item command="delete" class="text-red-500">
                    <Trash2 class="w-4 h-4 inline mr-2" />
                    {{ t('sidebar.deleteTag') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </span>
          <div
            v-if="tags.length === 0"
            class="text-xs text-slate-400 italic"
          >
            {{ t('sidebar.noTags') }}
          </div>
        </div>
      </div>

      <!-- Trash -->
      <div class="pt-4 border-t-2 border-slate-100">
        <MenuButton
          :icon="Trash2"
          :label="t('home.trash')"
          :active="currentView === 'trash'"
          class-name="text-red-500 hover:bg-red-50"
          @click="emit('change-view', 'trash'); playSound('click')"
          @mouseenter="playSound('hover')"
        />
      </div>
    </nav>

    <!-- User Info -->
    <div class="p-4 border-t-2 border-slate-100">
      <div
        class="flex items-center gap-3 p-2 rounded-xl bg-slate-100 border-2 border-slate-200 hover:border-black transition-colors cursor-pointer group"
        @click="goToSettings"
        @mouseenter="playSound('hover')"
      >
        <div
          class="w-10 h-10 rounded-full border-2 border-black overflow-hidden flex items-center justify-center font-bold text-lg"
          :style="{ backgroundColor: user?.avatar_color || '#fbbf24' }"
        >
          {{ AVATAR_STYLES[user?.avatar_style] || user?.username?.charAt(0) || '🐱' }}
        </div>
        <div class="overflow-hidden flex-1">
          <p class="font-bold text-sm truncate">{{ user?.nickname || user?.username || 'Guest' }}</p>
          <div class="text-xs text-slate-500 flex items-center gap-1">
            <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
            {{ t('sidebar.online') }}
          </div>
        </div>
        <Settings class="w-5 h-5 text-slate-400 group-hover:text-slate-600 group-hover:rotate-90 transition-all" />
      </div>
    </div>
  </aside>

  <!-- Modals -->
  <CreateModal
    :show="showNotebookModal"
    type="notebook"
    @close="showNotebookModal = false"
    @create="(name) => handleCreate(name, 'notebook')"
  />
  <CreateModal
    :show="showTagModal"
    type="tag"
    @close="showTagModal = false"
    @create="(name) => handleCreate(name, 'tag')"
  />
  <EditModal
    :show="showEditNotebookModal"
    type="notebook"
    :initial-name="editingNotebook?.name"
    @close="showEditNotebookModal = false; editingNotebook = null"
    @save="handleSaveNotebook"
  />
  <EditModal
    :show="showEditTagModal"
    type="tag"
    :initial-name="editingTag?.name"
    :initial-color="editingTag?.color"
    @close="showEditTagModal = false; editingTag = null"
    @save="handleSaveTag"
  />
</template>
