<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { AlertTriangle, Trash2 } from 'lucide-vue-next'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteAccount } from '../../api/user'

const emit = defineEmits(['account-deleted'])
const { t } = useI18n()

const loading = ref(false)
const showModal = ref(false)
const form = ref({
  password: '',
  confirm: ''
})

const openModal = () => {
  form.value = { password: '', confirm: '' }
  showModal.value = true
}

const handleDelete = async () => {
  // 验证
  if (!form.value.password) {
    ElMessage.warning(t('settings.enterPassword'))
    return
  }
  if (form.value.confirm !== 'DELETE') {
    ElMessage.warning(t('settings.typeDelete'))
    return
  }

  loading.value = true
  try {
    await deleteAccount({
      password: form.value.password,
      confirm: form.value.confirm
    })
    showModal.value = false
    emit('account-deleted')
  } catch (e) {
    console.error('Delete account failed:', e)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="bg-white border-4 border-red-500 rounded-2xl shadow-[8px_8px_0px_0px_rgba(239,68,68,1)] overflow-hidden">
    <!-- Section Header -->
    <div class="bg-red-500 border-b-4 border-red-600 px-6 py-4">
      <h2 class="text-xl font-black text-white flex items-center gap-2">
        <AlertTriangle class="w-6 h-6" />
        {{ t('settings.dangerZone') }}
      </h2>
    </div>

    <div class="p-6">
      <div class="flex items-center justify-between">
        <div>
          <p class="font-bold text-slate-800">{{ t('settings.deleteAccount') }}</p>
          <p class="text-sm text-slate-500">{{ t('settings.deleteAccountWarning') }}</p>
        </div>
        <button @click="openModal"
                class="flex items-center gap-2 px-4 py-2 bg-red-500 text-white border-4 border-black rounded-xl font-bold shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:-translate-y-1 transition-all active:shadow-none active:translate-y-1">
          <Trash2 class="w-4 h-4" />
          {{ t('settings.deleteAccount') }}
        </button>
      </div>
    </div>
  </section>

  <!-- Delete Confirmation Modal -->
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="showModal" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 backdrop-blur-sm" @click.self="showModal = false">
        <div class="bg-white border-4 border-black rounded-2xl p-6 w-[90vw] max-w-md shadow-[8px_8px_0px_0px_rgba(239,68,68,1)]">
          <div class="text-center mb-6">
            <div class="text-6xl mb-4">⚠️</div>
            <h3 class="font-black text-xl text-slate-800 mb-2">{{ t('settings.confirmDelete') }}</h3>
            <p class="text-slate-500 text-sm">{{ t('settings.deleteConfirmText') }}</p>
          </div>

          <div class="space-y-4">
            <!-- Password -->
            <div>
              <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
                {{ t('settings.password') }}
              </label>
              <input v-model="form.password"
                     type="password"
                     class="block w-full px-4 py-3 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(239,68,68,1)] transition-all font-bold"
                     :placeholder="t('settings.enterYourPassword')" />
            </div>

            <!-- Confirm DELETE -->
            <div>
              <label class="block text-xs font-black uppercase tracking-wider text-slate-500 mb-2">
                {{ t('settings.typeDeleteToConfirm') }}
              </label>
              <input v-model="form.confirm"
                     type="text"
                     class="block w-full px-4 py-3 border-4 border-black rounded-xl bg-slate-50 focus:bg-white focus:outline-none focus:shadow-[4px_4px_0px_0px_rgba(239,68,68,1)] transition-all font-bold"
                     placeholder="DELETE" />
            </div>
          </div>

          <div class="flex gap-3 mt-6">
            <button @click="showModal = false"
                    class="flex-1 py-3 border-4 border-black rounded-xl font-bold hover:bg-slate-100 transition-colors">
              {{ t('common.cancel') }}
            </button>
            <button @click="handleDelete"
                    :disabled="loading || form.confirm !== 'DELETE'"
                    class="flex-1 py-3 bg-red-500 text-white border-4 border-black rounded-xl font-bold shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-none hover:translate-x-1 hover:translate-y-1 transition-all disabled:opacity-50 disabled:cursor-not-allowed">
              {{ loading ? t('common.deleting') : t('settings.confirmDeleteBtn') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: all 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
