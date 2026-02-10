<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { User, Save } from 'lucide-vue-next'
import { useUserStore, AVATAR_STYLES, AVATAR_COLORS } from '../../stores/user'
import AvatarPicker from './AvatarPicker.vue'

const emit = defineEmits(['update'])
const { t } = useI18n()
const userStore = useUserStore()

const loading = ref(false)
const form = ref({
  nickname: '',
  email: '',
  bio: '',
  avatar_style: 'cat',
  avatar_color: '#fbbf24'
})

// 初始化表单
watch(() => userStore.user, (user) => {
  if (user) {
    form.value = {
      nickname: user.nickname || '',
      email: user.email || '',
      bio: user.bio || '',
      avatar_style: user.avatar_style || 'cat',
      avatar_color: user.avatar_color || '#fbbf24'
    }
  }
}, { immediate: true })

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('zh-CN')
}

const handleSubmit = async () => {
  loading.value = true
  try {
    await userStore.updateProfile(form.value)
    emit('update')
  } catch (e) {
    console.error('Update profile failed:', e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="bg-white border-4 border-black rounded-2xl shadow-[4px_4px_0px_0px_rgba(34,197,94,1)] md:shadow-[8px_8px_0px_0px_rgba(34,197,94,1)] overflow-hidden">
    <!-- Section Header -->
    <div class="bg-green-500 border-b-4 border-black px-4 md:px-6 py-3 md:py-4">
      <h2 class="text-lg md:text-xl font-black text-white flex items-center gap-2">
        <User class="w-5 h-5 md:w-6 md:h-6" />
        {{ t('settings.profile') }}
      </h2>
    </div>

    <div class="p-4 md:p-6 space-y-4 md:space-y-6">
      <!-- Avatar & User Info -->
      <div class="flex items-center gap-4 md:gap-6">
        <div class="w-16 h-16 md:w-20 md:h-20 rounded-2xl border-4 border-black flex items-center justify-center text-3xl md:text-4xl shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] md:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] flex-shrink-0"
             :style="{ backgroundColor: form.avatar_color }">
          {{ AVATAR_STYLES[form.avatar_style] || '🐱' }}
        </div>
        <div class="min-w-0">
          <p class="font-black text-lg md:text-xl text-slate-800 truncate">{{ userStore.user?.username }}</p>
          <p class="text-xs md:text-sm text-slate-500 font-bold">
            {{ t('settings.memberSince') }} {{ formatDate(userStore.user?.created_at) }}
          </p>
          <div class="flex flex-wrap gap-2 md:gap-4 mt-2 text-xs font-bold text-slate-400">
            <span>{{ userStore.user?.total_notes || 0 }} {{ t('settings.notes') }}</span>
            <span>{{ (userStore.user?.total_chars || 0).toLocaleString() }} {{ t('settings.chars') }}</span>
            <span>🔥 {{ userStore.user?.current_streak || 0 }} {{ t('settings.streak') }}</span>
          </div>
        </div>
      </div>

      <!-- Form Fields -->
      <div class="grid gap-4">
        <!-- Nickname -->
        <div>
          <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
            {{ t('settings.nickname') }}
          </label>
          <input v-model="form.nickname" type="text" maxlength="100"
                 autocomplete="off"
                 class="block w-full px-4 py-3 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(34,197,94,1)] transition-all font-bold"
                 :placeholder="t('settings.nicknamePlaceholder')" />
        </div>

        <!-- Email -->
        <div>
          <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
            {{ t('settings.email') }}
          </label>
          <input v-model="form.email" type="email" maxlength="255"
                 autocomplete="off"
                 class="block w-full px-4 py-3 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(34,197,94,1)] transition-all font-bold"
                 :placeholder="t('settings.emailPlaceholder')" />
        </div>

        <!-- Bio -->
        <div>
          <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
            {{ t('settings.bio') }}
          </label>
          <textarea v-model="form.bio" maxlength="500" rows="3"
                    class="block w-full px-4 py-3 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(34,197,94,1)] transition-all font-bold resize-none"
                    :placeholder="t('settings.bioPlaceholder')"></textarea>
        </div>

        <!-- Avatar Picker -->
        <AvatarPicker v-model:style="form.avatar_style" v-model:color="form.avatar_color" />
      </div>

      <!-- Submit Button -->
      <div class="flex justify-center md:justify-end">
        <button @click="handleSubmit"
                :disabled="loading"
                class="flex items-center gap-2 px-5 md:px-6 py-2.5 md:py-3 bg-green-500 text-white border-4 border-black rounded-xl font-black text-base md:text-lg shadow-[3px_3px_0px_0px_rgba(0,0,0,1)] md:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] md:hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:-translate-y-1 transition-all active:shadow-none active:translate-y-1 disabled:opacity-50 disabled:cursor-not-allowed">
          <Save class="w-4 h-4 md:w-5 md:h-5" />
          {{ loading ? t('common.saving') : t('settings.saveChanges') }}
        </button>
      </div>
    </div>
  </section>
</template>
