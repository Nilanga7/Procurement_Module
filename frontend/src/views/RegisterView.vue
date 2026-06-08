<template>
  <div class="min-h-screen bg-[#f8f9ff] flex items-center justify-center p-6">
    <div class="w-full max-w-md">
      <div class="bg-white border border-[#c8c5d3] rounded-2xl p-10 ambient-shadow">

        <!-- Brand -->
        <div class="text-center mb-8">
          <div class="w-14 h-14 bg-[#1a146b] rounded-2xl flex items-center justify-center mx-auto mb-4">
            <span class="material-symbols-outlined text-white text-[28px]"
                  style="font-variation-settings:'FILL' 1">person_add</span>
          </div>
          <h1 class="text-[24px] font-bold text-[#0b1c30] mb-1">Create an account</h1>
          <p class="text-[14px] text-[#777682]">ProcurePro ERP — Procurement Module</p>
        </div>

        <!-- Success state -->
        <div v-if="success" class="text-center py-6">
          <div class="w-16 h-16 bg-[#6cf8bb] rounded-full flex items-center justify-center mx-auto mb-4">
            <span class="material-symbols-outlined text-[#00714d] text-[32px]"
                  style="font-variation-settings:'FILL' 1">check_circle</span>
          </div>
          <p class="text-[18px] font-bold text-[#0b1c30] mb-2">Account Created!</p>
          <p class="text-[14px] text-[#777682] mb-6">You can now sign in with your credentials.</p>
          <button id="go-to-login" @click="router.push('/login')"
            class="w-full bg-[#1a146b] text-white py-3 rounded-xl font-bold text-[14px] hover:opacity-90 transition-all">
            Go to Sign In
          </button>
        </div>

        <!-- Form -->
        <form v-else @submit.prevent="handleRegister" class="space-y-4">
          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Full Name</label>
            <input id="reg-name" v-model="form.name" type="text" required autofocus
              placeholder="e.g. John Silva"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
          </div>

          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Email</label>
            <input id="reg-email" v-model="form.email" type="email" required
              placeholder="you@example.com"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
          </div>

          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">
              Password <span class="font-normal text-[#777682]">(min. 6 characters)</span>
            </label>
            <input id="reg-password" v-model="form.password" type="password" required minlength="6"
              placeholder="••••••••"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
          </div>

          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Role</label>
            <select id="reg-role" v-model="form.role" required
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10 bg-white">
              <option value="" disabled>Select your role…</option>
              <option value="procurement_officer">Procurement Officer</option>
              <option value="procurement_manager">Procurement Manager</option>
              <option value="supplier">Supplier</option>
            </select>
          </div>

          <!-- Company name — only for suppliers -->
          <div v-if="form.role === 'supplier'">
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Company / Business Name *</label>
            <input id="reg-company" v-model="form.company_name" type="text"
              :required="form.role === 'supplier'"
              placeholder="e.g. ABC Suppliers Ltd"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
            <p class="mt-2 flex items-center gap-1.5 text-[12px] text-[#006c49]">
              <span class="material-symbols-outlined text-[14px]" style="font-variation-settings:'FILL' 1">check_circle</span>
              Your supplier profile will be created automatically on registration.
            </p>
          </div>

          <div v-if="error"
            class="flex items-center gap-2 px-4 py-3 bg-[#ffdad6] border border-[#ba1a1a]/30
                   rounded-lg text-[13px] text-[#93000a]">
            <span class="material-symbols-outlined text-[16px]">error</span>
            {{ error }}
          </div>

          <button id="register-submit" type="submit" :disabled="loading"
            class="w-full bg-[#1a146b] text-white py-3 rounded-xl font-bold text-[14px]
                   hover:opacity-90 disabled:opacity-50 transition-all flex items-center justify-center gap-2">
            <svg v-if="loading" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
            </svg>
            {{ loading ? 'Creating account…' : 'Create Account' }}
          </button>
        </form>

        <p class="mt-6 text-center text-[13px] text-[#777682]">
          Already have an account?
          <router-link to="/login" class="text-[#1a146b] font-bold hover:underline ml-1">Sign in</router-link>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

const router  = useRouter()
const loading = ref(false)
const error   = ref('')
const success = ref(false)

const form = reactive({ name: '', email: '', password: '', role: '', company_name: '' })

async function handleRegister() {
  loading.value = true
  error.value   = ''
  try {
    await api.post('/api/auth/register', {
      name: form.name, email: form.email, password: form.password,
      role: form.role, company_name: form.company_name,
    })
    success.value = true
  } catch (e) {
    error.value = e.response?.data?.error || 'Registration failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
