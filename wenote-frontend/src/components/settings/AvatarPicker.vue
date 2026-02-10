<script setup>
import { useI18n } from 'vue-i18n'
import { AVATAR_STYLES, AVATAR_COLORS } from '../../stores/user'

const props = defineProps({
  style: { type: String, default: 'cat' },
  color: { type: String, default: '#fbbf24' }
})

const emit = defineEmits(['update:style', 'update:color'])
const { t } = useI18n()

const avatarList = Object.entries(AVATAR_STYLES)
</script>

<template>
  <div class="space-y-4">
    <!-- Avatar Style -->
    <div>
      <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
        {{ t('settings.selectAvatar') }}
      </label>
      <div class="grid grid-cols-5 sm:grid-cols-6 md:grid-cols-8 gap-2">
        <button v-for="[key, emoji] in avatarList" :key="key"
                @click="emit('update:style', key)"
                class="w-10 h-10 md:w-12 md:h-12 rounded-xl border-4 text-xl md:text-2xl flex items-center justify-center transition-all hover:scale-110"
                :class="style === key
                  ? 'border-black bg-slate-100 shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]'
                  : 'border-slate-200 hover:border-slate-400'">
          {{ emoji }}
        </button>
      </div>
    </div>

    <!-- Avatar Color -->
    <div>
      <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
        {{ t('settings.avatarColor') }}
      </label>
      <div class="flex flex-wrap gap-2 md:gap-3">
        <button v-for="c in AVATAR_COLORS" :key="c"
                @click="emit('update:color', c)"
                class="w-9 h-9 md:w-10 md:h-10 rounded-xl border-4 transition-all hover:scale-110"
                :class="color === c
                  ? 'border-black shadow-[2px_2px_0px_0px_rgba(0,0,0,1)]'
                  : 'border-slate-200 hover:border-slate-400'"
                :style="{ backgroundColor: c }">
        </button>
      </div>
    </div>
  </div>
</template>
