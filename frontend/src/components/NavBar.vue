<template>
  <nav class="navbar">
    <div class="nav-brand" @click="router.push('/')">
      📦 Procurement
    </div>

    <div class="nav-links">
      <router-link to="/">Dashboard</router-link>
      <router-link to="/requests">Requests</router-link>
      <!-- Bid Evaluations: only relevant to officers who evaluate competing bids -->
      <router-link
        v-if="auth.isOfficer"
        to="/requests?status=open_for_bidding"
        class="bid-eval-link"
      >🏷 Bid Evaluations</router-link>
      <router-link to="/purchase-orders">Purchase Orders</router-link>
      <router-link to="/goods-receipt">Goods Receipt</router-link>
      <router-link to="/transactions">Transactions</router-link>
    </div>

    <div class="nav-user">
      <div class="user-info">
        <span class="user-name">{{ auth.user?.name }}</span>
        <span class="user-role">{{ roleLabel }}</span>
      </div>
      <button class="logout-btn" @click="handleLogout">Logout</button>
    </div>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth   = useAuthStore()

const roleLabel = computed(() => {
  const labels = {
    admin:                'Admin',
    procurement_officer:  'Procurement Officer',
    procurement_manager:  'Procurement Manager',
    supplier:             'Supplier',
  }
  return labels[auth.user?.role] || auth.user?.role
})

function handleLogout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.navbar {
  display: flex;
  align-items: center;
  gap: 24px;
  padding: 0 24px;
  height: 56px;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 50;
}

.nav-brand {
  font-size: 16px;
  font-weight: 700;
  color: #4f46e5;
  cursor: pointer;
  white-space: nowrap;
  flex-shrink: 0;
}

.nav-links {
  display: flex;
  gap: 4px;
  flex: 1;
}
.nav-links a {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  color: #6b7280;
  text-decoration: none;
  transition: all .15s;
}
.nav-links a:hover        { background: #f3f4f6; color: #111827; }
.nav-links a.router-link-active { background: #ede9fe; color: #4f46e5; }

.nav-user {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
  margin-left: auto;
}
.user-info {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}
.user-name { font-size: 13px; font-weight: 500; color: #111827; }
.user-role { font-size: 11px; color: #9ca3af; }

.logout-btn {
  padding: 6px 14px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  background: #fff;
  font-size: 13px;
  cursor: pointer;
}
.logout-btn:hover { background: #f9fafb; }
</style>
