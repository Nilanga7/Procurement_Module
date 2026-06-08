<template>
  <div class="min-h-screen bg-[#f8f9ff] flex items-center justify-center p-6">
    <div class="w-full max-w-md">
      <!-- Card -->
      <div class="bg-white border border-[#c8c5d3] rounded-2xl p-10 ambient-shadow">
        <!-- Brand -->
        <div class="text-center mb-8">
          <div class="w-14 h-14 bg-[#1a146b] rounded-2xl flex items-center justify-center mx-auto mb-4">
            <span class="material-symbols-outlined text-white text-[28px]"
                  style="font-variation-settings:'FILL' 1">hub</span>
          </div>
          <h1 class="text-[24px] font-bold text-[#0b1c30] mb-1">Welcome back</h1>
          <p class="text-[14px] text-[#777682]">Sign in to ProcurePro ERP</p>
        </div>

        <!-- Form -->
        <form @submit.prevent="handleLogin" class="space-y-4">
          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Email</label>
            <input id="login-email" v-model="email" type="email" required autofocus
              placeholder="you@company.com"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10
                     transition-colors" />
          </div>

          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Password</label>
            <input id="login-password" v-model="password" type="password" required
              placeholder="••••••••"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10
                     transition-colors" />
          </div>

          <!-- Error -->
          <div v-if="error"
            class="flex items-center gap-2 px-4 py-3 bg-[#ffdad6] border border-[#ba1a1a]/30
                   rounded-lg text-[13px] text-[#93000a]">
            <span class="material-symbols-outlined text-[16px]">error</span>
            {{ error }}
          </div>

          <button id="login-submit" type="submit" :disabled="loading"
            class="w-full bg-[#1a146b] text-white py-3 rounded-xl font-bold text-[14px]
                   hover:opacity-90 disabled:opacity-50 transition-all flex items-center justify-center gap-2">
            <svg v-if="loading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
            </svg>
            {{ loading ? 'Signing in…' : 'Sign In' }}
          </button>
        </form>

        <p class="mt-6 text-center text-[13px] text-[#777682]">
          Don't have an account?
          <router-link to="/register" class="text-[#1a146b] font-bold hover:underline ml-1">Create one</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router   = useRouter()
const auth     = useAuthStore()
const email    = ref('')
const password = ref('')
const loading  = ref(false)
const error    = ref('')

async function handleLogin() {
  loading.value = true
  error.value   = ''
  try {
    await auth.login(email.value, password.value)
    if (auth.isSupplier) {
      router.push('/supplier')
    } else {
      router.push('/')
    }
  } catch (e) {
    error.value = e.response?.data?.error || 'Invalid email or password.'
  } finally {
    loading.value = false
  }
}
</script>
