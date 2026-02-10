<template>
  <div class="fixed bottom-0 w-full bg-black border-t-4 border-green-500 py-2 overflow-hidden z-50 hidden md:block">
    <div ref="trackRef" class="marquee-track whitespace-nowrap text-green-400 font-mono text-sm font-bold">
      <span ref="contentRef" class="marquee-content">{{ text }}</span>
      <span class="marquee-content" aria-hidden="true">{{ text }}</span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()
const trackRef = ref(null)
const contentRef = ref(null)

const text = computed(() => {
  const items = [
    `★ ${t('footer.welcome')}`,
    `★ ${t('footer.highScore')}`,
    `★ ${t('footer.collectNotes')}`,
    `★ ${t('footer.poweredBy')}`,
    `★ ${t('footer.insertCoin')}`,
    `★ ${t('footer.levelUp')}`,
    `★ ${t('footer.readyPlayer')}`,
    `★ ${t('footer.combo')}`,
    `★ ${t('footer.achievement')}`,
    `★ ${t('footer.pressStart')}`,
    `★ ${t('footer.playerReady')}`,
    `★ ${t('footer.loadingCreativity')}`,
    `★ ${t('footer.savePoint')}`,
    `★ ${t('footer.newQuest')}`,
    `★ ${t('footer.bonusStage')}`,
    `★ ${t('footer.ggwp')}`,
  ]
  return items.join('    ') + '    '
})

function updateAnimationDuration() {
  nextTick(() => {
    if (contentRef.value && trackRef.value) {
      const w = contentRef.value.offsetWidth
      const duration = w / 60
      trackRef.value.style.animationDuration = `${duration}s`
      if (!trackRef.value.classList.contains('marquee-animate')) {
        trackRef.value.classList.add('marquee-animate')
      }
    }
  })
}

watch(locale, updateAnimationDuration)

onMounted(() => updateAnimationDuration())
</script>

<style scoped>
.marquee-track {
  display: inline-flex;
  will-change: transform;
}
.marquee-content {
  flex-shrink: 0;
}
.marquee-animate {
  animation: marquee-scroll var(--duration, 30s) linear infinite;
}
@keyframes marquee-scroll {
  0% { transform: translateX(0); }
  100% { transform: translateX(-50%); }
}
</style>
