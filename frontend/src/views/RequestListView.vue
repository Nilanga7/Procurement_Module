<template>
  <div class="p-10">
    <!-- Header -->
    <header class="mb-8 flex justify-between items-end">
      <div>
        <nav class="flex items-center gap-1 text-[12px] font-semibold tracking-wide text-[#777682] mb-2">
          <span>Requests</span>
          <span class="material-symbols-outlined text-[14px]">chevron_right</span>
          <span class="text-[#1a146b]">All Requests</span>
        </nav>
        <h1 class="text-[32px] font-bold leading-10 tracking-tight text-[#0b1c30]">Procurement Requests</h1>
      </div>
      <button v-if="auth.isOfficer" @click="router.push('/requests/new')"
        class="bg-[#1a146b] text-white px-5 py-2.5 rounded-xl text-[12px] font-bold
               flex items-center gap-2 hover:opacity-90 transition-all ambient-shadow">
        <span class="material-symbols-outlined text-[18px]">add_circle</span>
        New Request
      </button>
    </header>

    <!-- Filter bar -->
    <div class="bg-white border border-[#c8c5d3] rounded-xl p-4 mb-6 flex flex-wrap gap-4 items-center ambient-shadow">
      <div class="flex items-center gap-2">
        <span class="material-symbols-outlined text-[18px] text-[#777682]">filter_list</span>
        <span class="text-[12px] font-semibold text-[#474651]">Filter by status:</span>
      </div>
      <select v-model="statusFilter"
        class="border border-[#c8c5d3] rounded-lg px-3 py-1.5 text-[13px] text-[#0b1c30]
               focus:outline-none focus:border-[#1a146b] bg-white">
        <option value="">All Statuses</option>
        <option value="draft">Draft</option>
        <option value="sent_to_finance">Sent to Finance</option>
        <option value="budget_approved">Budget Approved</option>
        <option value="budget_rejected">Budget Rejected</option>
        <option value="manager_approved">Manager Approved</option>
        <option value="manager_rejected">Manager Rejected</option>
        <option value="open_for_bidding">Open for Bidding</option>
        <option value="supplier_selected">Supplier Selected</option>
        <option value="purchase_order_created">PO Created</option>
        <option value="goods_received">Goods Received</option>
        <option value="completed">Completed</option>
      </select>
      <span class="text-[12px] text-[#777682] ml-auto">{{ requests.length }} result{{ requests.length !== 1 ? 's' : '' }}</span>
    </div>

    <!-- Error banner -->
    <div v-if="fetchError"
      class="mb-4 px-4 py-3 bg-[#ffdad6] border border-[#ba1a1a]/30 rounded-xl text-[14px] text-[#93000a] flex items-center gap-2">
      <span class="material-symbols-outlined text-[18px]">error</span>
      {{ fetchError }}
    </div>

    <!-- Loading -->
    <div v-if="loading" class="bg-white border border-[#c8c5d3] rounded-xl p-16 text-center ambient-shadow">
      <span class="material-symbols-outlined text-[40px] text-[#c8c5d3] block mb-3">hourglass_empty</span>
      <p class="text-[14px] text-[#777682]">Loading requests…</p>
    </div>

    <!-- Empty -->
    <div v-else-if="requests.length === 0 && !fetchError"
      class="bg-white border border-[#c8c5d3] rounded-xl p-16 text-center ambient-shadow">
      <span class="material-symbols-outlined text-[48px] text-[#c8c5d3] block mb-3">inbox</span>
      <p class="text-[14px] text-[#474651] font-semibold">No requests found</p>
      <p class="text-[13px] text-[#777682] mt-1">
        {{ statusFilter ? 'Try a different status filter.' : 'Create the first procurement request to get started.' }}
      </p>
    </div>

    <!-- Table -->
    <div v-else class="bg-white border border-[#c8c5d3] rounded-xl ambient-shadow overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead class="bg-[#d3e4fe]/30">
            <tr>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Request No.</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Title</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Department</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Est. Total</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Status</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Date</th>
              <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase text-right">Action</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-[#c8c5d3]/30">
            <tr v-for="r in requests" :key="r.id"
                class="hover:bg-[#d3e4fe]/10 transition-colors cursor-pointer"
                @click="router.push(`/requests/${r.id}`)">
              <td class="px-6 py-4 text-[14px] font-bold text-[#1a146b]">{{ r.request_number }}</td>
              <td class="px-6 py-4 text-[14px] text-[#0b1c30] max-w-[200px] truncate">{{ r.title }}</td>
              <td class="px-6 py-4 text-[14px] text-[#474651]">{{ r.department || '—' }}</td>
              <td class="px-6 py-4 text-[14px] text-[#474651]">{{ formatCurrency(r.estimated_total) }}</td>
              <td class="px-6 py-4"><StatusBadge :status="r.status" /></td>
              <td class="px-6 py-4 text-[14px] text-[#474651]">{{ formatDate(r.created_at) }}</td>
              <td class="px-6 py-4 text-right">
                <button @click.stop="router.push(`/requests/${r.id}`)"
                  class="text-[#1a146b] font-bold text-[12px] hover:underline">View →</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="px-6 py-3 border-t border-[#c8c5d3]/50 bg-[#eff4ff]/30">
        <p class="text-[12px] text-[#777682]">Showing {{ requests.length }} request{{ requests.length !== 1 ? 's' : '' }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import api from '../api'
import StatusBadge from '../components/StatusBadge.vue'

const router       = useRouter()
const route        = useRoute()
const auth         = useAuthStore()
const requests     = ref([])
const statusFilter = ref(route.query.status || '')
const loading      = ref(false)
const fetchError   = ref('')

async function fetchRequests() {
  loading.value    = true
  fetchError.value = ''
  try {
    const params = statusFilter.value ? { status: statusFilter.value } : {}
    const { data } = await api.get('/api/procurement/requests', { params })
    requests.value = Array.isArray(data) ? data : []
  } catch (e) {
    fetchError.value = e.response?.data?.error || 'Failed to load requests. Please try again.'
    requests.value = []
  } finally {
    loading.value = false
  }
}

function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-GB')
}

function formatCurrency(amount) {
  if (!amount) return '—'
  return 'LKR ' + Number(amount).toLocaleString('en-LK')
}

onMounted(fetchRequests)
watch(statusFilter, fetchRequests)
watch(() => route.query.status, (val) => { statusFilter.value = val || '' })
</script>
