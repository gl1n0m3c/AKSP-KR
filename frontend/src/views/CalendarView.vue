<template>
  <div class="layout">
    <section class="hero">
      <div>
        <div class="eyebrow">Планирование встреч</div>
        <h1>Календарь недели</h1>
        <p class="sub">Создавай события, выбирай участников и управляй командой в одном месте.</p>
        <div class="actions">
          <button class="primary" @click="openCreate">Новое событие</button>
          <button class="ghost" @click="today">Сегодня</button>
        </div>
      </div>
      <div class="week-nav">
        <div class="nav-block">
          <button class="circle" @click="prevWeek">‹</button>
          <div class="range">{{ weekRangeLabel }}</div>
          <button class="circle" @click="nextWeek">›</button>
        </div>
        <div class="user-pill" v-if="auth.user">
          {{ auth.user.first_name }} {{ auth.user.last_name }}
        </div>
        <button v-else class="primary ghost" @click="showAuth = true">Войти / Регистрация</button>
      </div>
    </section>

    <div class="calendar-shell">
      <aside class="sidebar">
        <div class="panel">
          <div class="panel-title">Месяц</div>
          <div class="month">{{ currentMonth }}</div>
        </div>
        <div class="panel">
          <div class="panel-title">Команды</div>
          <div class="pill-list">
            <span class="pill" v-for="t in teams" :key="t.id">{{ t.name }}</span>
          </div>
        </div>
        <div class="panel">
          <div class="panel-title">Должности</div>
          <div class="pill-list">
            <span class="pill" v-for="p in positions" :key="p.id">{{ p.name }}</span>
          </div>
        </div>
      </aside>

      <section class="calendar">
        <div class="hours">
          <div v-for="h in 14" :key="h" class="hour">{{ h + 7 }}:00</div>
        </div>
        <div class="days">
          <div v-for="day in days" :key="day.iso" class="day-col">
            <div class="day-header">{{ day.label }}</div>
            <div class="slots">
              <div
                v-for="ev in eventsByDay(day.iso)"
                :key="ev.id"
                class="event-card"
                @click="openEdit(ev)"
              >
                <div class="event-title">{{ ev.title }}</div>
                <div class="event-time">{{ formatTime(ev.scheduled_at) }}</div>
                <div class="event-meta">
                  <span>{{ (ev.participants || []).length }} участников</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>

    <EventModal
      v-if="showEventModal"
      :event="editingEvent"
      :users="users"
      :teams="teams"
      @close="closeEvent"
      @save="saveEvent"
    />
    <AuthModal
      v-if="showAuth"
      :positions="positions"
      :teams="teams"
      :units="units"
      :users="users"
      @authed="onAuthed"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import dayjs from 'dayjs'
import api from '../api'
import { useAuth } from '../composables/useAuth'
import type { EventItem, User, Team, Position, Unit } from '../types'
import EventModal from '../components/EventModal.vue'
import AuthModal from '../components/AuthModal.vue'

const auth = useAuth()
const showAuth = ref(false)
const showEventModal = ref(false)
const editingEvent = ref<EventItem | null>(null)
const events = ref<EventItem[]>([])
const users = ref<User[]>([])
const teams = ref<Team[]>([])
const positions = ref<Position[]>([])
const units = ref<Unit[]>([])

const currentWeekStart = ref(dayjs().startOf('week').add(1, 'day')) // Monday

const days = computed(() => {
  return Array.from({ length: 7 }).map((_, i) => {
    const d = currentWeekStart.value.add(i, 'day')
    return {
      iso: d.format('YYYY-MM-DD'),
      label: d.format('dd, D MMM'),
    }
  })
})

const weekRangeLabel = computed(() => {
  const start = currentWeekStart.value.format('D MMM')
  const end = currentWeekStart.value.add(6, 'day').format('D MMM')
  return `${start} — ${end}`
})

const currentMonth = computed(() => currentWeekStart.value.format('MMMM YYYY'))

const eventsByDay = (iso: string) => {
  return events.value.filter((e) => dayjs(e.scheduled_at).format('YYYY-MM-DD') === iso)
}

const formatTime = (ts: string) => dayjs(ts).format('HH:mm')

const loadEvents = async () => {
  const res = await api.get<EventItem[]>('/api/events')
  events.value = res.data
}

const loadUsers = async () => {
  try {
    const res = await api.get<User[]>('/api/users')
    users.value = res.data
  } catch (e: any) {
    if (e?.response?.status === 401 || e?.response?.status === 403) {
      showAuth.value = true
      return
    }
    throw e
  }
}

const loadMeta = async () => {
  const [t, p, u] = await Promise.all([
    api.get<Team[]>('/api/teams'),
    api.get<Position[]>('/api/positions'),
    api.get<Unit[]>('/api/units'),
  ])
  teams.value = t.data
  positions.value = p.data
  units.value = u.data
}

const openCreate = () => {
  if (!auth.user.value) {
    showAuth.value = true
    return
  }
  editingEvent.value = null
  showEventModal.value = true
}

const openEdit = (ev: EventItem) => {
  editingEvent.value = ev
  showEventModal.value = true
}

const closeEvent = () => {
  showEventModal.value = false
  editingEvent.value = null
}

const saveEvent = async (payload: any) => {
  if (editingEvent.value) {
    await api.put(`/api/events/${editingEvent.value.id}`, payload)
  } else {
    await api.post('/api/events', payload)
  }
  await loadEvents()
  closeEvent()
}

const onAuthed = async () => {
  showAuth.value = false
  await auth.loadMe()
  if (auth.user.value) {
    await loadUsers()
    await loadMeta()
    await loadEvents()
  }
}

const prevWeek = () => {
  currentWeekStart.value = currentWeekStart.value.subtract(1, 'week')
}
const nextWeek = () => {
  currentWeekStart.value = currentWeekStart.value.add(1, 'week')
}

const today = () => {
  currentWeekStart.value = dayjs().startOf('week').add(1, 'day')
}

onMounted(async () => {
  await loadMeta()
  await loadUsers()
  await auth.loadMe()
  if (!auth.user.value) {
    showAuth.value = true
    return
  }
  await loadEvents()
})
</script>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.hero {
  background: linear-gradient(135deg, rgba(79, 70, 229, 0.2), rgba(20, 184, 166, 0.2));
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 16px;
  padding: 20px;
  display: flex;
  justify-content: space-between;
  gap: 12px;
}
.eyebrow {
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-size: 12px;
  color: #a5b4fc;
}
.hero h1 {
  margin: 4px 0;
}
.sub {
  color: #94a3b8;
  margin: 6px 0 12px;
}
.actions {
  display: flex;
  gap: 8px;
}
.primary {
  background: linear-gradient(135deg, #6366f1, #14b8a6);
  border: none;
  color: #fff;
  padding: 10px 14px;
  border-radius: 12px;
  cursor: pointer;
}
.ghost {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.14);
  color: #e5e7eb;
  padding: 10px 14px;
  border-radius: 12px;
  cursor: pointer;
}
.week-nav {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-end;
}
.nav-block {
  display: flex;
  gap: 8px;
  align-items: center;
}
.circle {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 1px solid rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.04);
  color: #fff;
  cursor: pointer;
}
.range {
  font-weight: 600;
}
.user-pill {
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.08);
}
.calendar-shell {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 12px;
}
.sidebar {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.panel {
  background: #0f172a;
  border: 1px solid #1f2937;
  border-radius: 12px;
  padding: 12px;
}
.panel-title {
  font-size: 13px;
  color: #94a3b8;
  margin-bottom: 6px;
}
.month {
  font-weight: 700;
  font-size: 16px;
}
.pill-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}
.pill {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  padding: 6px 10px;
  font-size: 12px;
  border: 1px solid rgba(255, 255, 255, 0.08);
}
.calendar {
  display: grid;
  grid-template-columns: 70px 1fr;
  gap: 8px;
  background: #0f172a;
  border: 1px solid #1f2937;
  border-radius: 16px;
  padding: 10px;
}
.hours {
  display: grid;
  grid-template-rows: repeat(14, 1fr);
  gap: 6px;
  color: #9ca3af;
  font-size: 12px;
}
.days {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 6px;
}
.day-col {
  background: rgba(255, 255, 255, 0.02);
  border: 1px dashed rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  min-height: 620px;
  display: flex;
  flex-direction: column;
}
.day-header {
  padding: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-weight: 600;
}
.slots {
  flex: 1;
  padding: 6px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}
.event-card {
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.4), rgba(20, 184, 166, 0.35));
  border: 1px solid rgba(255, 255, 255, 0.12);
  color: #fff;
  padding: 10px;
  border-radius: 10px;
  cursor: pointer;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.25);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}
.event-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.35);
}
.event-title {
  font-weight: 700;
}
.event-time {
  font-size: 12px;
  opacity: 0.9;
}
.event-meta {
  font-size: 12px;
  color: #cbd5e1;
  margin-top: 4px;
}
</style>

