<template>
  <div class="modal-backdrop">
    <div class="modal">
      <h2 v-if="mode === 'register'">Регистрация</h2>
      <h2 v-else>Вход</h2>

      <form @submit.prevent="submit">
        <div v-if="mode === 'register'" class="grid">
          <label>Имя <input v-model="form.first_name" required /></label>
          <label>Фамилия <input v-model="form.last_name" required /></label>
          <label>Отчество <input v-model="form.patronymic" required /></label>
          <label>Email <input v-model="form.email" type="email" required /></label>
          <label>Телефон <input v-model="form.phone" /></label>
          <label>Telegram <input v-model="form.telegram" /></label>
          <label>
            Должность
            <select v-model="form.position_id" required>
              <option value="">Выберите</option>
              <option v-for="p in positions" :key="p.id" :value="p.id">{{ p.name }}</option>
            </select>
          </label>
          <label>
            Команда
            <select v-model="form.team_id" required>
              <option value="">Выберите</option>
              <option v-for="t in teams" :key="t.id" :value="t.id">{{ t.name }}</option>
            </select>
          </label>
          <label>
            Юнит
            <select v-model="form.unit_id" required>
              <option value="">Выберите</option>
              <option v-for="u in units" :key="u.id" :value="u.id">{{ u.name }}</option>
            </select>
          </label>
          <label>
            Руководитель (опционально)
            <select v-model="form.head_id">
              <option value="">Не выбран</option>
              <option v-for="u in users" :key="u.id" :value="u.id">
                {{ u.last_name }} {{ u.first_name }}
              </option>
            </select>
          </label>
        </div>

        <label>Логин <input v-model="form.login" required /></label>
        <label>Пароль <input v-model="form.password" type="password" required /></label>

        <div class="actions">
          <button type="submit">{{ mode === 'register' ? 'Зарегистрироваться' : 'Войти' }}</button>
          <button type="button" class="ghost" @click="toggleMode">
            {{ mode === 'register' ? 'У меня уже есть аккаунт' : 'Зарегистрироваться' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useAuth } from '../composables/useAuth'
import type { Position, Team, Unit, User } from '../types'

const props = defineProps<{
  positions: Position[]
  teams: Team[]
  units: Unit[]
  users: User[]
}>()

const emit = defineEmits(['authed'])
const { login, register } = useAuth()

const mode = ref<'login' | 'register'>('register')

const form = reactive({
  login: '',
  password: '',
  first_name: '',
  last_name: '',
  patronymic: '',
  email: '',
  phone: '',
  telegram: '',
  position_id: '',
  team_id: '',
  unit_id: '',
  head_id: '',
  status: 'active',
})

const submit = async () => {
  if (mode.value === 'register') {
    await register({
      login: form.login,
      password: form.password,
      first_name: form.first_name,
      last_name: form.last_name,
      patronymic: form.patronymic,
      email: form.email,
      phone: form.phone,
      telegram: form.telegram,
      position_id: form.position_id || null,
      team_id: form.team_id || null,
      unit_id: form.unit_id || null,
      head_id: form.head_id || null,
      status: form.status,
    })
  } else {
    await login(form.login, form.password)
  }
  emit('authed')
}

const toggleMode = () => {
  mode.value = mode.value === 'login' ? 'register' : 'login'
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
  background: #0f172a;
  color: #fff;
  padding: 24px;
  border-radius: 12px;
  width: 520px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.08);
}
form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 14px;
}
input,
select {
  background: #111827;
  border: 1px solid #1f2937;
  border-radius: 8px;
  padding: 8px 10px;
  color: #fff;
}
.actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}
button {
  background: #4f46e5;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 8px 12px;
  cursor: pointer;
}
.ghost {
  background: transparent;
  border: 1px solid #444;
}
.grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 8px;
}
</style>

