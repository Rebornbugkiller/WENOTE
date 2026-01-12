// 预设颜色
export const tagColors = ['#ef4444', '#f97316', '#eab308', '#22c55e', '#06b6d4', '#3b82f6', '#8b5cf6', '#ec4899']

// 根据标签名生成固定颜色
export const getTagColor = (tag) => {
  if (tag.color && tag.color !== '#6B7280') return tag.color
  let hash = 0
  for (let i = 0; i < tag.name.length; i++) hash = tag.name.charCodeAt(i) + ((hash << 5) - hash)
  return tagColors[Math.abs(hash) % tagColors.length]
}
