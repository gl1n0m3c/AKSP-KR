<template>
  <div class="modal-backdrop" @click.self="close">
    <div class="modal">
      <div class="header-row">
        <h2>{{ event ? 'Редактирование события' : 'Создать событие' }}</h2>
        <div class="subtext">Заполни детали встречи и выбери участников</div>
      </div>

      <form @submit.prevent="save" class="form-grid">
        <label class="wide">
          <span>Название</span>
          <input v-model="form.title" placeholder="Синк, демо, 1:1..." required />
        </label>

        <label class="wide">
          <span>Описание</span>
          <textarea v-model="form.description" placeholder="Повестка, цели, ссылки..." />
        </label>

        <div class="row">
          <label class="half">
            <span>Дата и время</span>
            <input type="datetime-local" v-model="form.scheduled_at" required />
          </label>
          <label class="half">
            <span>Команда</span>
            <select v-model="form.team_id">
              <option value="">Без команды</option>
              <option v-for="t in teams" :key="t.id" :value="t.id">{{ t.name }}</option>
            </select>
          </label>
        </div>

        <div class="participants">
          <div class="participants-header">Участники</div>
          <div class="list">
            <label v-for="u in users" :key="u.id" class="user-item">
              <input type="checkbox" :value="u.id" v-model="form.participants" />
              <div class="user-meta">
                <div class="user-name">{{ u.last_name }} {{ u.first_name }}</div>
                <div class="user-extra">{{ u.email || 'нет email' }}</div>
              </div>
            </label>
          </div>
        </div>

        <div class="actions">
          <button type="submit">Сохранить</button>
          <button type="button" class="ghost" @click="close">Закрыть</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, watch } from 'vue'
import type { EventItem, User, Team } from '../types'

const props = defineProps<{
  event?: EventItem | null
  users: User[]
  teams: Team[]
}>()

const emit = defineEmits(['close', 'save'])

const form = reactive({
  title: '',
  description: '',
  scheduled_at: '',
  team_id: '',
  participants: [] as string[],
})

watch(
  () => props.event,
  (ev) => {
    if (ev) {
      form.title = ev.title
      form.description = ev.description
      form.scheduled_at = ev.scheduled_at.slice(0, 16)
      form.team_id = ev.team_id || ''
      form.participants = [...ev.participants]
    } else {
      form.title = ''
      form.description = ''
      form.scheduled_at = ''
      form.team_id = ''
      form.participants = []
    }
  },
  { immediate: true }
)

const close = () => emit('close')
const save = () => {
  emit('save', {
    title: form.title,
    description: form.description,
    scheduled_at: new Date(form.scheduled_at).toISOString(),
    team_id: form.team_id || null,
    participants: form.participants,
  })
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 50;
}
.modal {
  background: radial-gradient(circle at 20% 20%, rgba(99, 102, 241, 0.12), transparent 35%),
    radial-gradient(circle at 80% 0%, rgba(20, 184, 166, 0.14), transparent 32%),
    #0b1220;
  color: #fff;
  padding: 24px;
  border-radius: 16px;
  width: 720px;
  max-height: 90vh;
  overflow-y: auto;
  border: 1px solid rgba(255, 255, 255, 0.06);
  box-shadow: 0 24px 80px rgba(0, 0, 0, 0.55);
}
.header-row {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-bottom: 12px;
}
.header-row h2 {
  margin: 0;
}
.subtext {
  color: #94a3b8;
  font-size: 14px;
}
form.form-grid {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  color: #cbd5e1;
  font-size: 14px;
}
input,
textarea,
select {
  background: #0f172a;
  border: 1px solid #1f2937;
  border-radius: 8px;
  padding: 10px 12px;
  color: #e2e8f0;
  transition: border 0.15s ease, box-shadow 0.15s ease;
}
textarea {
  min-height: 80px;
}
input:focus,
textarea:focus,
select:focus {
  outline: none;
  border-color: #6366f1;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.18);
}
.row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}
.half {
  width: 100%;
}
.participants {
  border: 1px solid #1f2937;
  border-radius: 12px;
  padding: 12px;
  background: rgba(15, 23, 42, 0.82);
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.participants-header {
  margin-bottom: 2px;
  font-weight: 600;
  color: #e2e8f0;
}
.list {
  max-height: 200px;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 8px;
}
.user-item {
  display: flex;
  gap: 8px;
  align-items: center;
  padding: 8px 10px;
  border: 1px solid #1e293b;
  border-radius: 10px;
  background: rgba(30, 41, 59, 0.6);
}
.user-meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.user-name {
  font-weight: 600;
  color: #e5e7eb;
}
.user-extra {
  color: #94a3b8;
  font-size: 12px;
}
.actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-top: 4px;
}
button {
  background: linear-gradient(135deg, #6366f1, #14b8a6);
  color: #fff;
  border: none;
  border-radius: 10px;
  padding: 10px 14px;
  cursor: pointer;
  font-weight: 600;
}
.ghost {
  background: transparent;
  border: 1px solid #444;
  color: #e2e8f0;
}
</style>

