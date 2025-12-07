<template>
  <div class="layout">
    <header class="topbar">
      <RouterLink to="/" class="link">← Календарь</RouterLink>
      <div class="title">Админ-панель</div>
    </header>

    <div class="grid">
      <section class="card">
        <div class="card-header">
          <h3>Команды</h3>
          <p>Создай или удали команду</p>
        </div>
        <form @submit.prevent="createTeam" class="form-row">
          <input v-model="team.name" placeholder="Название" required />
          <input v-model="team.description" placeholder="Описание" />
          <button type="submit" class="primary">Создать</button>
        </form>
        <ul class="list">
          <li v-for="t in teams" :key="t.id" class="list-item">
            <div>
              <div class="item-title">{{ t.name }}</div>
              <div class="item-sub">{{ t.description }}</div>
            </div>
            <button class="ghost" @click="deleteTeam(t.id)">Удалить</button>
          </li>
        </ul>
      </section>

      <section class="card">
        <div class="card-header">
          <h3>Должности</h3>
          <p>Справочник позиций</p>
        </div>
        <form @submit.prevent="createPosition" class="form-row">
          <input v-model="position.name" placeholder="Название" required />
          <input v-model="position.description" placeholder="Описание" />
          <button type="submit" class="primary">Создать</button>
        </form>
        <ul class="list">
          <li v-for="p in positions" :key="p.id" class="list-item">
            <div>
              <div class="item-title">{{ p.name }}</div>
              <div class="item-sub">{{ p.description }}</div>
            </div>
            <button class="ghost" @click="deletePosition(p.id)">Удалить</button>
          </li>
        </ul>
      </section>

      <section class="card">
        <div class="card-header">
          <h3>Юниты</h3>
          <p>Организационные единицы</p>
        </div>
        <form @submit.prevent="createUnit" class="form-row">
          <input v-model="unit.name" placeholder="Название" required />
          <input v-model="unit.description" placeholder="Описание" />
          <button type="submit" class="primary">Создать</button>
        </form>
        <ul class="list">
          <li v-for="u in units" :key="u.id" class="list-item">
            <div>
              <div class="item-title">{{ u.name }}</div>
              <div class="item-sub">{{ u.description }}</div>
            </div>
            <button class="ghost" @click="deleteUnit(u.id)">Удалить</button>
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { RouterLink } from 'vue-router'
import api from '../api'
import type { Team, Position, Unit } from '../types'
import { v4 as uuid } from 'uuid'

const teams = ref<Team[]>([])
const positions = ref<Position[]>([])
const units = ref<Unit[]>([])

const team = reactive({ name: '', description: '' })
const position = reactive({ name: '', description: '' })
const unit = reactive({ name: '', description: '' })

const loadAll = async () => {
  const [t, p, u] = await Promise.all([
    api.get<Team[]>('/api/teams'),
    api.get<Position[]>('/api/positions'),
    api.get<Unit[]>('/api/units'),
  ])
  teams.value = t.data
  positions.value = p.data
  units.value = u.data
}

const createTeam = async () => {
  await api.post('/api/teams', { id: uuid(), name: team.name, description: team.description })
  team.name = ''
  team.description = ''
  await loadAll()
}
const deleteTeam = async (id: string) => {
  await api.delete(`/api/teams/${id}`)
  await loadAll()
}

const createPosition = async () => {
  await api.post('/api/positions', { id: uuid(), name: position.name, description: position.description })
  position.name = ''
  position.description = ''
  await loadAll()
}
const deletePosition = async (id: string) => {
  await api.delete(`/api/positions/${id}`)
  await loadAll()
}

const createUnit = async () => {
  await api.post('/api/units', { id: uuid(), name: unit.name, description: unit.description })
  unit.name = ''
  unit.description = ''
  await loadAll()
}
const deleteUnit = async (id: string) => {
  await api.delete(`/api/units/${id}`)
  await loadAll()
}

onMounted(loadAll)
</script>

<style scoped>
.layout {
  min-height: 100vh;
  color: #f1f5f9;
}
.topbar {
  display: flex;
  gap: 12px;
  align-items: center;
  padding: 12px 16px;
  background: #111827;
  border-bottom: 1px solid #1f2937;
}
.title {
  font-weight: 700;
  font-size: 18px;
}
.link {
  color: #c7d2fe;
  text-decoration: none;
}
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 16px;
  padding: 16px;
}
.card {
  background: #0f172a;
  border: 1px solid #1f2937;
  border-radius: 14px;
  padding: 14px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.28);
}
.card-header p {
  margin: 0;
  color: #94a3b8;
}
.form-row {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 10px;
}
input {
  background: #1c1c1c;
  border: 1px solid #333;
  border-radius: 8px;
  padding: 6px 10px;
  color: #fff;
  flex: 1;
  min-width: 140px;
}
button {
  background: #4f46e5;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 6px 10px;
  cursor: pointer;
}
.primary {
  background: linear-gradient(135deg, #6366f1, #14b8a6);
}
.ghost {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.16);
}
.list {
  list-style: none;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.item-title {
  font-weight: 700;
}
.item-sub {
  color: #94a3b8;
  font-size: 13px;
}
</style>

