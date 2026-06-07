<template>
  <div class="p-10">
    <!-- Page header -->
    <header class="mb-8 flex justify-between items-end">
      <div>
        <nav class="flex items-center gap-1 text-[12px] font-semibold tracking-wide text-[#777682] mb-2">
          <span>{{ roleLabel }}</span>
          <span class="material-symbols-outlined text-[14px]">chevron_right</span>
          <span class="text-[#1a146b]">Dashboard</span>
        </nav>
        <h1 class="text-[32px] font-bold leading-10 tracking-tight text-[#0b1c30]">
          Procurement Dashboard
        </h1>
      </div>
      <button @click="router.push('/requests/new')" v-if="auth.isOfficer"
        class="bg-[#1a146b] text-white px-5 py-2.5 rounded-xl text-[12px] font-bold
               flex items-center gap-2 hover:opacity-90 transition-all ambient-shadow">
        <span class="material-symbols-outlined text-[18px]">add_circle</span>
        New Request
      </button>
    </header>

    <!-- Stat cards -->
    <section class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <!-- Total -->
      <div class="bg-white border border-[#c8c5d3] rounded-xl p-4 ambient-shadow hover-shadow">
        <div class="flex justify-between items-start mb-3">
          <div class="p-2 bg-[#312e81] text-[#9c9af4] rounded-lg">
            <span class="material-symbols-outlined text-[20px]">folder_open</span>
          </div>
          <span class="text-[11px] font-medium text-[#006c49]">All time</span>
        </div>
        <p class="text-[12px] font-semibold tracking-wide text-[#474651] uppercase mb-1">Total Requests</p>
        <p class="text-[48px] font-bold leading-none tracking-tight text-[#1a146b]">{{ stats.total }}</p>
      </div>

      <!-- Pending -->
      <div class="bg-white border border-[#c8c5d3] rounded-xl p-4 ambient-shadow hover-shadow">
        <div class="flex justify-between items-start mb-3">
          <div class="p-2 bg-[#ffddb8] text-[#653e00] rounded-lg">
            <span class="material-symbols-outlined text-[20px]">pending_actions</span>
          </div>
          <span class="text-[11px] font-medium text-[#777682]">Awaiting review</span>
        </div>
        <p class="text-[12px] font-semibold tracking-wide text-[#474651] uppercase mb-1">Pending Approval</p>
        <p class="text-[48px] font-bold leading-none tracking-tight text-[#e49200]">{{ stats.pending }}</p>
      </div>

      <!-- Open for Bidding -->
      <div class="bg-white border border-[#c8c5d3] rounded-xl p-4 ambient-shadow hover-shadow">
        <div class="flex justify-between items-start mb-3">
          <div class="p-2 bg-[#3e3c8f] text-[#e2dfff] rounded-lg">
            <span class="material-symbols-outlined text-[20px]"
                  style="font-variation-settings:'FILL' 1">gavel</span>
          </div>
          <span class="text-[11px] font-medium text-[#3e3c8f]">Active now</span>
        </div>
        <p class="text-[12px] font-semibold tracking-wide text-[#474651] uppercase mb-1">Open for Bidding</p>
        <p class="text-[48px] font-bold leading-none tracking-tight text-[#3e3c8f]">{{ stats.bidding }}</p>
      </div>

      <!-- Completed -->
      <div class="bg-white border border-[#c8c5d3] rounded-xl p-4 ambient-shadow hover-shadow">
        <div class="flex justify-between items-start mb-3">
          <div class="p-2 bg-[#6cf8bb] text-[#00714d] rounded-lg">
            <span class="material-symbols-outlined text-[20px]">task_alt</span>
          </div>
          <span class="text-[11px] font-medium text-[#006c49]">Done</span>
        </div>
        <p class="text-[12px] font-semibold tracking-wide text-[#474651] uppercase mb-1">Completed</p>
        <p class="text-[48px] font-bold leading-none tracking-tight text-[#006c49]">{{ stats.completed }}</p>
      </div>
    </section>

    <!-- Recent requests table -->
    <section class="bg-white border border-[#c8c5d3] rounded-xl ambient-shadow overflow-hidden">
      <div class="px-6 py-4 border-b border-[#c8c5d3] flex justify-between items-center">
        <h2 class="text-[20px] font-semibold text-[#0b1c30]">Recent Requests</h2>
        <button @click="router.push('/requests')"
          class="text-[12px] font-bold text-[#1a146b] hover:underline flex items-center gap-1">
          View all
          <span class="material-symbols-outlined text-[16px]">arrow_forward</span>
        </button>
      </div>

      <div v-if="loading" class="p-12 text-center text-[#777682] text-[14px]">
        <span class="material-symbols-outlined text-[32px] block mb-2 text-[#c8c5d3]">hourglass_empty</span>
        Loading requests…
      </div>

      <div v-else-if="requests.length === 0" class="p-12 text-center text-[#777682] text-[14px]">
        <span class="material-symbols-outlined text-[40px] block mb-2 text-[#c8c5d3]">inbox</span>
        No requests yet.
        <router-link v-if="auth.isOfficer" to="/requests/new"
          class="text-[#1a146b] font-bold hover:underline ml-1">Create the first one →</router-link>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead class="bg-[#d3e4fe]/30">
            <tr>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Request No.</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Title</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Department</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Status</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Date</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#c8c5d3]/30">
            <tr v-for="r in requests.slice(0, 8)" :key="r.id"
                class="hover:bg-[#d3e4fe]/10 transition-colors cursor-pointer"
                @click="router.push(`/requests/${r.id}`)">
              <td class="px-6 py-4 text-[14px] font-bold text-[#1a146b]">{{ r.request_number }}</td>
              <td class="px-6 py-4 text-[14px] text-[#0b1c30]">{{ r.title }}</td>
              <td class="px-6 py-4 text-[14px] text-[#474651]">{{ r.department || '—' }}</td>
              <td class="px-6 py-4">
                <StatusBadge :status="r.status" />
              </td>
              <td class="px-6 py-4 text-[14px] text-[#474651]">{{ formatDate(r.created_at) }}</td>
              <td class="px-6 py-4 text-right">
                <button @click.stop="router.push(`/requests/${r.id}`)"
                  class="text-[#1a146b] font-bold text-[12px] hover:underline">View</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="requests.length > 0" class="px-6 py-4 border-t border-[#c8c5d3] text-center">
        <button @click="router.push('/requests')"
          class="inline-flex items-center gap-2 px-6 py-2 border border-[#c8c5d3] text-[#1a146b]
                 font-bold text-[12px] rounded-xl hover:bg-[#eff4ff] transition-colors">
          View All Requests
          <span class="material-symbols-outlined text-[16px]">arrow_forward</span>
        </button>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../api'
import StatusBadge from '../components/StatusBadge.vue'

const router  = useRouter()
const auth    = useAuthStore()
const requests = ref([])
const loading  = ref(false)

const roleLabel = computed(() => {
  const labels = {
    admin:               'Admin',
    procurement_officer: 'Officer',
    procurement_manager: 'Manager',
    supplier:            'Supplier',
  }
  return labels[auth.user?.role] || auth.user?.role
})

const stats = computed(() => ({
  total:     requests.value.length,
  pending:   requests.value.filter(r => ['sent_to_finance','budget_approved','manager_approved'].includes(r.status)).length,
  bidding:   requests.value.filter(r => r.status === 'open_for_bidding').length,
  completed: requests.value.filter(r => r.status === 'completed').length,
}))

async function fetchAll() {
  loading.value = true
  try {
    const { data } = await api.get('/api/procurement/requests')
    requests.value = Array.isArray(data) ? data : []
  } catch {
    requests.value = []
  } finally {
    loading.value = false
  }
}

function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-GB')
}

onMounted(fetchAll)
</script>
