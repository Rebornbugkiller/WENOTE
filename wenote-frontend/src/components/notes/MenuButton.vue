<script setup>
defineProps({
  icon: Object,
  label: String,
  count: [Number, String],
  active: Boolean,
  className: String
})
defineEmits(['click'])
</script>

<template>
  <button
    @click="$emit('click')"
    class="w-full flex items-center justify-between px-4 py-3 rounded-xl font-bold transition-all border-2"
    :class="[
      active
        ? 'bg-yellow-400 border-black text-black shadow-[3px_3px_0px_0px_rgba(0,0,0,1)]'
        : 'bg-transparent border-transparent text-slate-600 hover:bg-slate-100 hover:text-black',
      className
    ]"
  >
    <div class="flex items-center gap-3">
      <component :is="icon" class="w-5 h-5" />
      <span>{{ label }}</span>
    </div>
    <span
      v-if="count !== undefined"
      class="text-xs px-2 py-0.5 rounded-full"
      :class="active ? 'bg-black/10' : 'bg-slate-200'"
    >
      {{ count }}
    </span>
  </button>
</template>
