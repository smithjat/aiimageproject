<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const menuItems = [
  { key: 'text-to-image', label: '文本生图', icon: '📝', path: '/text-to-image', desc: 'AI智能生成' },
  { key: 'image-edit', label: '图片编辑', icon: '🖼️', path: '/image-edit', desc: '专业修图' },
  { key: 'cross-border', label: '跨境电商', icon: '🌍', path: '/cross-border', desc: '电商专用' },
  { key: 'food-business', label: '餐饮商家', icon: '🍽️', path: '/food-business', desc: '美食营销' },
]

const activeKey = ref(route.path)

const handleMenuClick = (item: typeof menuItems[0]) => {
  activeKey.value = item.path
  router.push(item.path)
}
</script>

<template>
  <div class="sidebar">
    <div class="logo">
      <div class="logo-icon">
        <span class="logo-letter">AI</span>
      </div>
      <div class="logo-text">
        <span class="logo-title">AI 创作平台</span>
        <span class="logo-subtitle">企业级图像解决方案</span>
      </div>
    </div>
    <nav class="nav-menu">
      <div
        v-for="item in menuItems"
        :key="item.key"
        :class="['nav-item', { active: activeKey === item.path }]"
        @click="handleMenuClick(item)"
      >
        <span class="nav-icon">{{ item.icon }}</span>
        <div class="nav-content">
          <span class="nav-label">{{ item.label }}</span>
          <span class="nav-desc">{{ item.desc }}</span>
        </div>
        <span v-if="activeKey === item.path" class="nav-arrow">›</span>
      </div>
    </nav>
    <div class="sidebar-footer">
      <div class="footer-item">
        <span class="footer-icon">📞</span>
        <span class="footer-text">客服支持</span>
      </div>
      <div class="footer-item">
        <span class="footer-icon">📖</span>
        <span class="footer-text">使用教程</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.sidebar {
  width: 240px;
  height: 100vh;
  background: #1e2a3b;
  color: #fff;
  display: flex;
  flex-direction: column;
  position: fixed;
  left: 0;
  top: 0;
}

.logo {
  padding: 24px 20px;
  display: flex;
  align-items: center;
  gap: 14px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.logo-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-letter {
  font-size: 18px;
  font-weight: 700;
  color: #fff;
}

.logo-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.logo-title {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
}

.logo-subtitle {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
}

.nav-menu {
  padding: 16px 12px;
  flex: 1;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 4px;
  position: relative;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.nav-item.active {
  background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  transform: translateY(-50%);
  width: 3px;
  height: 24px;
  background: #60a5fa;
  border-radius: 0 2px 2px 0;
}

.nav-icon {
  font-size: 20px;
  width: 24px;
  text-align: center;
}

.nav-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-label {
  font-size: 14px;
  font-weight: 500;
}

.nav-desc {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.5);
}

.nav-arrow {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.5);
}

.sidebar-footer {
  padding: 16px 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.footer-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
}

.footer-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.footer-icon {
  font-size: 16px;
}

.footer-text {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
}
</style>
