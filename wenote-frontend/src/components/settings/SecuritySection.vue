<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { Lock, Eye, EyeOff } from 'lucide-vue-next'
import { ElMessage } from 'element-plus'
import { changePassword } from '../../api/user'

const emit = defineEmits(['password-changed'])
const { t } = useI18n()

const loading = ref(false)
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

const form = ref({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const handleSubmit = async () => {
  // 验证
  if (!form.value.current_password || !form.value.new_password) {
    ElMessage.warning(t('settings.fillAllFields'))
    return
  }
  if (form.value.new_password.length < 6) {
    ElMessage.warning(t('settings.passwordTooShort'))
    return
  }
  if (form.value.new_password !== form.value.confirm_password) {
    ElMessage.warning(t('settings.passwordMismatch'))
    return
  }

  loading.value = true
  try {
    await changePassword({
      current_password: form.value.current_password,
      new_password: form.value.new_password
    })
    // 清空表单
    form.value = { current_password: '', new_password: '', confirm_password: '' }
    emit('password-changed')
  } catch (e) {
    console.error('Change password failed:', e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="bg-white border-4 border-black rounded-2xl shadow-[8px_8px_0px_0px_rgba(59,130,246,1)] overflow-hidden">
    <!-- Section Header -->
    <div class="bg-blue-500 border-b-4 border-black px-6 py-4">
      <h2 class="text-xl font-black text-white flex items-center gap-2">
        <Lock class="w-6 h-6" />
        {{ t('settings.security') }}
      </h2>
    </div>

    <div class="p-6 space-y-4">
      <!-- Current Password -->
      <div>
        <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
          {{ t('settings.currentPassword') }}
        </label>
        <div class="relative">
          <input v-model="form.current_password"
                 :type="showCurrentPassword ? 'text' : 'password'"
                 class="block w-full px-4 py-3 pr-12 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(59,130,246,1)] transition-all font-bold"
                 :placeholder="t('settings.currentPasswordPlaceholder')" />
          <button @click="showCurrentPassword = !showCurrentPassword"
                  type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600">
            <EyeOff v-if="showCurrentPassword" class="w-5 h-5" />
            <Eye v-else class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- New Password -->
      <div>
        <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
          {{ t('settings.newPassword') }}
        </label>
        <div class="relative">
          <input v-model="form.new_password"
                 :type="showNewPassword ? 'text' : 'password'"
                 class="block w-full px-4 py-3 pr-12 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(59,130,246,1)] transition-all font-bold"
                 :placeholder="t('settings.newPasswordPlaceholder')" />
          <button @click="showNewPassword = !showNewPassword"
                  type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600">
            <EyeOff v-if="showNewPassword" class="w-5 h-5" />
            <Eye v-else class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Confirm Password -->
      <div>
        <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
          {{ t('settings.confirmPassword') }}
        </label>
        <div class="relative">
          <input v-model="form.confirm_password"
                 :type="showConfirmPassword ? 'text' : 'password'"
                 class="block w-full px-4 py-3 pr-12 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(59,130,246,1)] transition-all font-bold"
                 :placeholder="t('settings.confirmPasswordPlaceholder')" />
          <button @click="showConfirmPassword = !showConfirmPassword"
                  type="button"
                  class="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600">
            <EyeOff v-if="showConfirmPassword" class="w-5 h-5" />
            <Eye v-else class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Submit Button -->
      <div class="flex justify-end pt-2">
        <button @click="handleSubmit"
                :disabled="loading"
                class="flex items-center gap-2 px-6 py-3 bg-blue-500 text-white border-4 border-black rounded-xl font-black text-lg shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:-translate-y-1 transition-all active:shadow-none active:translate-y-1 disabled:opacity-50 disabled:cursor-not-allowed">
          <Lock class="w-5 h-5" />
          {{ loading ? t('common.saving') : t('settings.changePassword') }}
        </button>
      </div>
    </div>
  </section>
</template>
