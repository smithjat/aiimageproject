<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

onMounted(() => {
  userStore.fetchProfile()
})
</script>

<template>
  <div class="home">
    <header class="header">
      <h1>AI Image Project</h1>
      <nav>
        <router-link to="/">首页</router-link>
        <router-link to="/images">图库</router-link>
        <template v-if="userStore.isLoggedIn">
          <router-link to="/profile">{{ userStore.user?.nickname || userStore.user?.username }}</router-link>
          <a href="#" @click.prevent="userStore.logout()">退出</a>
        </template>
        <template v-else>
          <router-link to="/login">登录</router-link>
          <router-link to="/register">注册</router-link>
        </template>
      </nav>
    </header>
    <main class="main">
      <h2>欢迎使用 AI Image Project</h2>
      <p>这是一个支持 Web、微信小程序、移动端的前后端分离项目</p>
    </main>
  </div>
</template>

<style scoped>
.home {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 2rem;
  background: #409eff;
  color: white;
}

.header nav {
  display: flex;
  gap: 1rem;
}

.header nav a {
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background 0.3s;
}

.header nav a:hover {
  background: rgba(255, 255, 255, 0.2);
}

.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}
</style>
