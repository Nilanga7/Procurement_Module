<template>
  <!-- ── Full-screen layout for authenticated internal staff ── -->
  <div v-if="showNav" class="min-h-screen bg-[#f8f9ff]">

    <!-- ── Top Nav ── -->
    <nav class="fixed top-0 left-0 w-full z-50 flex justify-between items-center px-10 h-16
                bg-[#f8f9ff] border-b border-[#c8c5d3]">
      <div class="flex items-center gap-6">
        <span class="text-[20px] font-bold text-[#1a146b] leading-7">ProcurePro ERP</span>
      </div>
      <div class="flex items-center gap-3">
        <!-- User pill -->
        <div class="flex items-center gap-2 px-3 py-1.5 rounded-full bg-[#eff4ff] cursor-pointer
                    hover:bg-[#dce9ff] transition-colors">
          <div class="w-7 h-7 rounded-full bg-[#1a146b] flex items-center justify-center">
            <span class="text-white text-[11px] font-bold">{{ initials }}</span>
          </div>
          <span class="text-[12px] font-semibold text-[#1a146b]">{{ auth.user?.name }}</span>
          <span class="text-[11px] text-[#777682]">{{ roleLabel }}</span>
        </div>
        <button @click="handleLogout"
          class="px-3 py-1.5 rounded-lg border border-[#c8c5d3] text-[12px] font-semibold
                 text-[#474651] hover:bg-[#eff4ff] transition-colors">
          Logout
        </button>
      </div>
    </nav>

    <!-- ── Sidebar ── -->
    <aside class="fixed left-0 top-16 bottom-0 w-64 flex flex-col p-4
                  bg-[#eff4ff] border-r border-[#c8c5d3] z-40">

      <!-- Brand block -->
      <div class="flex items-center gap-3 px-2 py-4 mb-4">
        <div class="w-10 h-10 bg-[#1a146b] rounded-xl flex items-center justify-center flex-shrink-0">
          <span class="material-symbols-outlined text-white text-[20px]"
                style="font-variation-settings:'FILL' 1">hub</span>
        </div>
        <div>
          <p class="text-[14px] font-bold text-[#1a146b] leading-tight">Procurement Hub</p>
          <p class="text-[11px] text-[#777682]">Enterprise Tier</p>
        </div>
      </div>

      <!-- Create button — only for officers -->
      <button v-if="auth.isOfficer"
        @click="router.push('/requests/new')"
        class="mb-6 bg-[#1a146b] text-white font-bold py-3 px-4 rounded-xl
               flex items-center justify-center gap-2 hover:opacity-90 transition-all
               ambient-shadow text-[12px]">
        <span class="material-symbols-outlined text-[18px]">add_circle</span>
        Create New Request
      </button>

      <!-- Nav links -->
      <nav class="flex-1 space-y-1">
        <router-link to="/" class="nav-item" active-class="sidebar-active" exact>
          <span class="material-symbols-outlined text-[20px]">dashboard</span>
          <span>Dashboard</span>
        </router-link>

        <router-link to="/requests" class="nav-item" active-class="sidebar-active">
          <span class="material-symbols-outlined text-[20px]">request_quote</span>
          <span>Requests</span>
        </router-link>

        <!-- Bid Evaluations — officer only -->
        <router-link v-if="auth.isOfficer"
          to="/requests?status=open_for_bidding"
          class="nav-item" active-class="sidebar-active">
          <span class="material-symbols-outlined text-[20px]">gavel</span>
          <span>Bid Evaluations</span>
        </router-link>

        <router-link to="/purchase-orders" class="nav-item" active-class="sidebar-active">
          <span class="material-symbols-outlined text-[20px]">shopping_cart</span>
          <span>Purchase Orders</span>
        </router-link>

        <router-link to="/goods-receipt" class="nav-item" active-class="sidebar-active">
          <span class="material-symbols-outlined text-[20px]">inventory_2</span>
          <span>Goods Receipt</span>
        </router-link>

        <router-link to="/transactions" class="nav-item" active-class="sidebar-active">
          <span class="material-symbols-outlined text-[20px]">receipt_long</span>
          <span>Transactions</span>
        </router-link>
      </nav>

      <!-- Bottom logout -->
      <div class="border-t border-[#c8c5d3] pt-4 mt-auto">
        <button @click="handleLogout"
          class="nav-item w-full text-[#ba1a1a] hover:bg-[#ffdad6]">
          <span class="material-symbols-outlined text-[20px]">logout</span>
          <span>Logout</span>
        </button>
      </div>
    </aside>

    <!-- ── Page content ── -->
    <main class="ml-64 mt-16 min-h-[calc(100vh-64px)]">
      <router-view />
    </main>
  </div>

  <!-- ── Public / supplier pages — no sidebar ── -->
  <router-view v-else />
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const route  = useRoute()
const router = useRouter()
const auth   = useAuthStore()

const showNav = computed(() => !route.meta.public && !route.meta.hideNav)

const initials = computed(() => {
  const name = auth.user?.name || ''
  return name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
})

const roleLabel = computed(() => {
  const labels = {
    admin:               'Admin',
    procurement_officer: 'Officer',
    procurement_manager: 'Manager',
    supplier:            'Supplier',
  }
  return labels[auth.user?.role] || auth.user?.role
})

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<style>
/* Nav item base style */
.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 16px;
  border-radius: 0.5rem;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.05em;
  color: #474651;
  text-decoration: none;
  transition: background 0.15s, color 0.15s;
  cursor: pointer;
  border: none;
  background: none;
  width: 100%;
  text-align: left;
}
.nav-item:hover {
  background-color: #dce9ff;
  color: #1a146b;
}
</style>