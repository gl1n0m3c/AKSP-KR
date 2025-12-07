<template>
  <div class="shell">
    <header class="topbar">
      <div class="brand">
        <span class="logo">🤙</span>
        <div class="brand-text">
          <div class="name">feemous</div>
          <div class="tagline">events & meetings</div>
        </div>
      </div>
      <nav class="nav">
        <RouterLink to="/">Календарь</RouterLink>
        <RouterLink to="/admin">Админка</RouterLink>
      </nav>
    </header>
    <main class="content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
</script>

<style scoped>
.shell {
  min-height: 100vh;
  background: radial-gradient(circle at 20% 20%, rgba(79, 70, 229, 0.12), transparent 35%),
    radial-gradient(circle at 80% 0%, rgba(16, 185, 129, 0.14), transparent 30%),
    #0b0d11;
  color: #e5e7eb;
}
.topbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 18px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(10px);
  position: sticky;
  top: 0;
  z-index: 20;
}
.brand {
  display: flex;
  gap: 10px;
  align-items: center;
}
.logo {
  font-size: 24px;
}
.brand-text .name {
  font-weight: 700;
  letter-spacing: 0.04em;
}
.brand-text .tagline {
  font-size: 12px;
  color: #94a3b8;
}
.nav {
  display: flex;
  gap: 12px;
}
.nav a {
  padding: 8px 12px;
  border-radius: 10px;
  text-decoration: none;
  color: #e5e7eb;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: all 0.2s ease;
}
.nav a.router-link-active {
  background: linear-gradient(135deg, #6366f1, #14b8a6);
  border-color: transparent;
}
.nav a:hover {
  transform: translateY(-1px);
}
.content {
  padding: 14px 18px 32px;
}
</style>

