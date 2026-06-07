import { defineStore } from 'pinia'
import api from '../api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user') || 'null'),
    token: localStorage.getItem('token') || null,
  }),
  getters: {
    isLoggedIn: s => !!s.token,
    isManager: s => s.user?.role === 'procurement_manager',
    isOfficer: s => s.user?.role === 'procurement_officer',
    isSupplier: s => s.user?.role === 'supplier',
  },
  actions: {
    async login(email, password) {
      const { data } = await api.post('/api/auth/login', { email, password })
      this.token = data.token
      this.user = data.user
      localStorage.setItem('token', data.token)
      localStorage.setItem('user', JSON.stringify(data.user))
    },
    logout() {
      this.token = null; this.user = null
      localStorage.removeItem('token'); localStorage.removeItem('user')
    }
  }
})