<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ArrowLeft, Gamepad2 } from 'lucide-vue-next'
import { ElMessage } from 'element-plus'
import { useUserStore } from '../stores/user'
import ProfileSection from '../components/settings/ProfileSection.vue'
import SecuritySection from '../components/settings/SecuritySection.vue'
import DangerZone from '../components/settings/DangerZone.vue'
import FloatingIcons from '../components/login/FloatingIcons.vue'

const router = useRouter()
const userStore = useUserStore()
const { t } = useI18n()

const isGameMode = ref(true)

const goBack = () => router.push('/')

const handleProfileUpdate = () => {
  ElMessage.success(t('settings.profileUpdated'))
}

const handlePasswordChanged = () => {
  ElMessage.success(t('settings.passwordChanged'))
}

const handleAccountDeleted = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<template>
  <div class="min-h-screen bg-green-50 font-sans overflow-hidden selection:bg-green-200 relative">
    <!-- Game Mode Overlay Elements -->
    <template v-if="isGameMode">
      <!-- CRT Scanlines -->
      <div class="fixed inset-0 z-50 pointer-events-none opacity-5"
           style="background: linear-gradient(rgba(18,16,16,0) 50%, rgba(0,0,0,0.25) 50%), linear-gradient(90deg, rgba(255,0,0,0.06), rgba(0,255,0,0.02), rgba(0,0,255,0.06)); background-size: 100% 2px, 3px 100%;"></div>

      <!-- Floating Background Icons -->
      <div class="fixed inset-0 pointer-events-none z-0">
        <FloatingIcons />
      </div>
    </template>

    <!-- Grid Background -->
    <div class="fixed inset-0 pointer-events-none overflow-hidden">
      <div class="grid-bg absolute inset-0 opacity-30" />
    </div>

    <!-- Header -->
    <header class="sticky top-0 z-20 p-6 bg-green-50/80 backdrop-blur-sm border-b-4 border-black">
      <div class="max-w-4xl mx-auto flex items-center justify-between">
        <div class="flex items-center gap-4">
          <button @click="goBack"
                  class="p-2 bg-white border-4 border-black rounded-xl hover:bg-slate-100 shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:-translate-y-1 transition-all active:translate-y-1 active:shadow-none">
            <ArrowLeft class="w-6 h-6" />
          </button>
          <h1 class="text-3xl font-black text-slate-800">
            {{ t('settings.title') }}
          </h1>
        </div>

        <!-- Game Mode Toggle -->
        <button @click="isGameMode = !isGameMode"
                class="p-2 rounded-xl border-2 border-black transition-all hover:scale-105 active:scale-95"
                :class="isGameMode ? 'bg-green-500 text-white shadow-[2px_2px_0px_0px_#000]' : 'bg-white text-slate-400'">
          <Gamepad2 class="w-6 h-6" />
        </button>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-4xl mx-auto p-6 space-y-8 relative z-10 pb-20">
      <!-- Profile Section -->
      <ProfileSection @update="handleProfileUpdate" />

      <!-- Security Section -->
      <SecuritySection @password-changed="handlePasswordChanged" />

      <!-- Danger Zone -->
      <DangerZone @account-deleted="handleAccountDeleted" />
    </main>

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
</style>
