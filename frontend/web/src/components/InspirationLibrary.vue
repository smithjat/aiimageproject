<script setup lang="ts">
import { ref, computed } from 'vue'

interface Inspiration {
  id: number
  type: 'text' | 'image'
  prompt?: string
  imageUrl?: string
  thumbnail?: string
  title: string
  tags?: string[]
  category?: string
}

const props = defineProps<{
  type: 'text' | 'image' | 'edit' | 'cross-border' | 'food'
}>()

const emit = defineEmits<{
  select: [inspiration: Inspiration]
}>()

const searchKeyword = ref('')
const activeCategory = ref('all')

const categories = [
  { key: 'all', label: '全部' },
  { key: 'product', label: '产品' },
  { key: 'scene', label: '场景' },
  { key: 'portrait', label: '人物' },
  { key: 'food', label: '美食' },
]

const textInspirations = ref<Inspiration[]>([
  { id: 1, type: 'text', prompt: '高端电子产品展示，简约白色背景，专业产品摄影，细节清晰', title: '电子产品主图', tags: ['产品', '电商'], category: 'product' },
  { id: 2, type: 'text', prompt: '时尚服装模特展示，自然光线，专业摄影棚拍摄，高清细节', title: '服装模特图', tags: ['服装', '模特'], category: 'portrait' },
  { id: 3, type: 'text', prompt: '精致美食摆盘，餐厅氛围灯光，专业美食摄影，诱人食欲', title: '美食摄影', tags: ['美食', '餐饮'], category: 'food' },
  { id: 4, type: 'text', prompt: '家居生活场景，温馨舒适氛围，自然光线，生活化展示', title: '家居场景', tags: ['家居', '场景'], category: 'scene' },
  { id: 5, type: 'text', prompt: '化妆品产品展示，高级质感，专业打光，品牌调性', title: '美妆产品', tags: ['美妆', '产品'], category: 'product' },
  { id: 6, type: 'text', prompt: '珠宝首饰特写，奢华质感，专业微距摄影，细节展现', title: '珠宝首饰', tags: ['珠宝', '产品'], category: 'product' },
  { id: 7, type: 'text', prompt: '跨境电商主图，白底产品图，符合平台规范，专业拍摄', title: '电商主图', tags: ['电商', '产品'], category: 'product' },
  { id: 8, type: 'text', prompt: '户外运动场景，活力四射，专业运动摄影，品牌宣传', title: '运动场景', tags: ['运动', '场景'], category: 'scene' },
])

const imageInspirations = ref<Inspiration[]>([
  { id: 1, type: 'image', imageUrl: 'https://picsum.photos/200/200?random=1', thumbnail: 'https://picsum.photos/200/200?random=1', title: '产品摄影', tags: ['产品', '电商'], category: 'product' },
  { id: 2, type: 'image', imageUrl: 'https://picsum.photos/200/200?random=2', thumbnail: 'https://picsum.photos/200/200?random=2', title: '人物肖像', tags: ['人物', '肖像'], category: 'portrait' },
  { id: 3, type: 'image', imageUrl: 'https://picsum.photos/200/200?random=3', thumbnail: 'https://picsum.photos/200/200?random=3', title: '美食图片', tags: ['美食', '餐饮'], category: 'food' },
  { id: 4, type: 'image', imageUrl: 'https://picsum.photos/200/200?random=4', thumbnail: 'https://picsum.photos/200/200?random=4', title: '场景展示', tags: ['场景', '家居'], category: 'scene' },
  { id: 5, type: 'image', imageUrl: 'https://picsum.photos/200/200?random=5', thumbnail: 'https://picsum.photos/200/200?random=5', title: '电商主图', tags: ['电商', '产品'], category: 'product' },
  { id: 6, type: 'image', imageUrl: 'https://picsum.photos/200/200?random=6', thumbnail: 'https://picsum.photos/200/200?random=6', title: '品牌宣传', tags: ['品牌', '宣传'], category: 'scene' },
])

const allInspirations = computed(() => {
  if (props.type === 'text' || props.type === 'cross-border' || props.type === 'food') {
    return textInspirations.value
  }
  return imageInspirations.value
})

const filteredInspirations = computed(() => {
  let result = allInspirations.value
  
  if (activeCategory.value !== 'all') {
    result = result.filter(item => item.category === activeCategory.value)
  }
  
  if (searchKeyword.value.trim()) {
    const keyword = searchKeyword.value.toLowerCase().trim()
    result = result.filter(item => {
      const titleMatch = item.title.toLowerCase().includes(keyword)
      const promptMatch = item.prompt?.toLowerCase().includes(keyword)
      const tagsMatch = item.tags?.some(tag => tag.toLowerCase().includes(keyword))
      return titleMatch || promptMatch || tagsMatch
    })
  }
  
  return result
})

const handleSelect = (item: Inspiration) => {
  emit('select', item)
}

const clearSearch = () => {
  searchKeyword.value = ''
}
</script>

<template>
  <div class="inspiration-library">
    <div class="library-header">
      <span class="header-title">💡 灵感库</span>
      <span class="header-count">{{ filteredInspirations.length }} 个模板</span>
    </div>
    
    <div class="search-section">
      <div class="search-box">
        <span class="search-icon">🔍</span>
        <input
          v-model="searchKeyword"
          type="text"
          class="search-input"
          placeholder="搜索模板..."
        />
        <button v-if="searchKeyword" class="clear-btn" @click="clearSearch">✕</button>
      </div>
    </div>
    
    <div class="category-tabs">
      <div
        v-for="cat in categories"
        :key="cat.key"
        :class="['tab-item', { active: activeCategory === cat.key }]"
        @click="activeCategory = cat.key"
      >
        {{ cat.label }}
      </div>
    </div>
    
    <div class="library-content">
      <div v-if="filteredInspirations.length === 0" class="empty-state">
        <span class="empty-icon">🔍</span>
        <p class="empty-text">未找到相关模板</p>
        <button class="btn-reset" @click="searchKeyword = ''; activeCategory = 'all'">重置筛选</button>
      </div>
      
      <div
        v-for="item in filteredInspirations"
        :key="item.id"
        class="inspiration-card"
        @click="handleSelect(item)"
      >
        <div v-if="item.type === 'image'" class="card-image">
          <img :src="item.thumbnail" :alt="item.title" />
        </div>
        <div v-else class="card-text">
          <p>{{ item.prompt }}</p>
        </div>
        <div class="card-info">
          <div class="card-title">{{ item.title }}</div>
          <div v-if="item.tags" class="card-tags">
            <span v-for="tag in item.tags.slice(0, 2)" :key="tag" class="tag">{{ tag }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.inspiration-library {
  width: 320px;
  height: 100vh;
  background: #fff;
  border-left: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  position: fixed;
  right: 0;
  top: 0;
}

.library-header {
  padding: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #e2e8f0;
}

.header-title {
  font-size: 16px;
  font-weight: 600;
  color: #1e293b;
}

.header-count {
  font-size: 13px;
  color: #64748b;
}

.search-section {
  padding: 12px 16px;
  border-bottom: 1px solid #e2e8f0;
}

.search-box {
  display: flex;
  align-items: center;
  background: #f1f5f9;
  border-radius: 6px;
  padding: 10px 12px;
  gap: 10px;
  transition: all 0.2s;
}

.search-box:focus-within {
  background: #fff;
  box-shadow: 0 0 0 2px rgba(37, 99, 235, 0.2);
}

.search-icon {
  font-size: 14px;
  opacity: 0.5;
}

.search-input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 14px;
  color: #1e293b;
  outline: none;
}

.search-input::placeholder {
  color: #94a3b8;
}

.clear-btn {
  width: 18px;
  height: 18px;
  border: none;
  background: #cbd5e1;
  border-radius: 50%;
  font-size: 10px;
  color: #64748b;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.clear-btn:hover {
  background: #94a3b8;
  color: #fff;
}

.category-tabs {
  display: flex;
  padding: 12px 16px;
  gap: 8px;
  border-bottom: 1px solid #e2e8f0;
  overflow-x: auto;
}

.tab-item {
  padding: 6px 14px;
  background: #f1f5f9;
  border-radius: 4px;
  font-size: 13px;
  color: #64748b;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s;
}

.tab-item:hover {
  background: #e2e8f0;
}

.tab-item.active {
  background: #2563eb;
  color: #fff;
}

.library-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  text-align: center;
}

.empty-icon {
  font-size: 40px;
  margin-bottom: 12px;
  opacity: 0.4;
}

.empty-text {
  font-size: 14px;
  color: #64748b;
  margin-bottom: 16px;
}

.btn-reset {
  padding: 8px 20px;
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
}

.btn-reset:hover {
  background: #1d4ed8;
}

.inspiration-card {
  background: #fff;
  border-radius: 6px;
  overflow: hidden;
  margin-bottom: 12px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid #e2e8f0;
}

.inspiration-card:hover {
  border-color: #2563eb;
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.1);
}

.card-image {
  width: 100%;
  height: 140px;
  overflow: hidden;
}

.card-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.card-text {
  padding: 14px;
  background: #f8fafc;
  min-height: 80px;
}

.card-text p {
  font-size: 13px;
  line-height: 1.6;
  margin: 0;
  color: #475569;
}

.card-info {
  padding: 10px 14px;
  border-top: 1px solid #f1f5f9;
}

.card-title {
  font-size: 13px;
  font-weight: 500;
  color: #1e293b;
  margin-bottom: 6px;
}

.card-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.tag {
  font-size: 11px;
  padding: 2px 6px;
  background: #f1f5f9;
  border-radius: 3px;
  color: #64748b;
}
</style>
