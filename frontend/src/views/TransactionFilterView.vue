<template>
  <div class="p-10">
    <!-- Header -->
    <header class="mb-8">
      <nav class="flex items-center gap-1 text-[12px] font-semibold tracking-wide text-[#777682] mb-2">
        <span>Transactions</span>
        <span class="material-symbols-outlined text-[14px]">chevron_right</span>
        <span class="text-[#1a146b]">Transaction Records</span>
      </nav>
      <div class="flex justify-between items-end">
        <h1 class="text-[32px] font-bold leading-10 tracking-tight text-[#0b1c30]">Transaction Records</h1>
        <button :disabled="!results.data?.length" @click="exportCSV"
          class="flex items-center gap-2 px-4 py-2.5 bg-[#6cf8bb]/30 border border-[#006c49]/30
                 text-[#006c49] rounded-xl text-[12px] font-bold hover:bg-[#6cf8bb]/50 transition-colors
                 disabled:opacity-40 disabled:cursor-not-allowed">
          <span class="material-symbols-outlined text-[18px]">download</span>
          Export CSV
        </button>
      </div>
    </header>

    <!-- Filter panel -->
    <div class="bg-white border border-[#c8c5d3] rounded-xl p-6 mb-6 ambient-shadow">
      <h2 class="text-[14px] font-bold text-[#0b1c30] mb-4 flex items-center gap-2">
        <span class="material-symbols-outlined text-[18px] text-[#1a146b]">filter_list</span>
        Filter Records
      </h2>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-[11px] font-semibold text-[#474651] mb-1.5">Request ID</label>
          <input v-model="filters.request_id" placeholder="e.g. PR0001" @input="search"
            class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                   focus:outline-none focus:border-[#1a146b]" />
        </div>
        <div>
          <label class="block text-[11px] font-semibold text-[#474651] mb-1.5">Supplier Name</label>
          <input v-model="filters.supplier" placeholder="e.g. ABC Suppliers" @input="search"
            class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                   focus:outline-none focus:border-[#1a146b]" />
        </div>
        <div>
          <label class="block text-[11px] font-semibold text-[#474651] mb-1.5">Department</label>
          <input v-model="filters.department" placeholder="e.g. IT" @input="search"
            class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                   focus:outline-none focus:border-[#1a146b]" />
        </div>
        <div>
          <label class="block text-[11px] font-semibold text-[#474651] mb-1.5">Status</label>
          <select v-model="filters.status" @change="search"
            class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                   focus:outline-none focus:border-[#1a146b] bg-white">
            <option value="">All</option>
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
        </div>
        <div>
          <label class="block text-[11px] font-semibold text-[#474651] mb-1.5">From Date</label>
          <input v-model="filters.from" type="date" @change="search"
            class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                   focus:outline-none focus:border-[#1a146b]" />
        </div>
        <div>
          <label class="block text-[11px] font-semibold text-[#474651] mb-1.5">To Date</label>
          <input v-model="filters.to" type="date" @change="search"
            class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                   focus:outline-none focus:border-[#1a146b]" />
        </div>
      </div>
      <div class="flex justify-between items-center mt-4 pt-4 border-t border-[#c8c5d3]/50">
        <span class="text-[12px] text-[#777682]">
          {{ results.total_count ?? 0 }} result{{ results.total_count !== 1 ? 's' : '' }}
        </span>
        <button @click="clearFilters"
          class="text-[12px] font-bold text-[#1a146b] hover:underline">Clear Filters</button>
      </div>
    </div>

    <!-- Results -->
    <div class="bg-white border border-[#c8c5d3] rounded-xl ambient-shadow overflow-hidden">
      <div v-if="loading" class="p-12 text-center">
        <span class="material-symbols-outlined text-[36px] text-[#c8c5d3] block mb-2">hourglass_empty</span>
        <p class="text-[14px] text-[#777682]">Searching…</p>
      </div>
      <table v-else class="w-full text-left border-collapse">
        <thead class="bg-[#d3e4fe]/30">
          <tr>
            <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Request No.</th>
            <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Title</th>
            <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Department</th>
            <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Supplier</th>
            <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">PO Number</th>
            <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Status</th>
            <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Date</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-[#c8c5d3]/30">
          <tr v-for="r in results.data" :key="r.id" class="hover:bg-[#d3e4fe]/10 transition-colors">
            <td class="px-6 py-4 text-[14px] font-bold text-[#1a146b]">{{ r.request_number }}</td>
            <td class="px-6 py-4 text-[14px] text-[#0b1c30]">{{ r.title }}</td>
            <td class="px-6 py-4 text-[14px] text-[#474651]">{{ r.department || '—' }}</td>
            <td class="px-6 py-4 text-[14px] text-[#474651]">{{ r.supplier_name || '—' }}</td>
            <td class="px-6 py-4 text-[13px] font-mono text-[#3e3c8f]">{{ r.po_number || '—' }}</td>
            <td class="px-6 py-4"><StatusBadge :status="r.status" /></td>
            <td class="px-6 py-4 text-[14px] text-[#474651]">{{ formatDate(r.created_at) }}</td>
          </tr>
          <tr v-if="!results.data?.length">
            <td colspan="7" class="px-6 py-12 text-center text-[#777682] text-[14px]">
              No results. Use the filters above to search.
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div v-if="results.total_pages > 1"
        class="flex justify-center items-center gap-4 p-4 border-t border-[#c8c5d3]/50">
        <button :disabled="results.page === 1" @click="goToPage(results.page - 1)"
          class="px-4 py-1.5 border border-[#c8c5d3] rounded-lg text-[12px] font-semibold
                 disabled:opacity-40 hover:bg-[#eff4ff] transition-colors">← Prev</button>
        <span class="text-[13px] text-[#777682]">Page {{ results.page }} of {{ results.total_pages }}</span>
        <button :disabled="results.page === results.total_pages" @click="goToPage(results.page + 1)"
          class="px-4 py-1.5 border border-[#c8c5d3] rounded-lg text-[12px] font-semibold
                 disabled:opacity-40 hover:bg-[#eff4ff] transition-colors">Next →</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import api from '../api'
import StatusBadge from '../components/StatusBadge.vue'

const filters = reactive({ request_id: '', supplier: '', status: '', from: '', to: '', department: '' })
const results = ref({ data: [], total_count: 0, total_pages: 0, page: 1 })
const loading = ref(false)
let currentPage = 1

async function search(page = 1) {
  loading.value = true
  currentPage = page
  const params = Object.fromEntries(Object.entries(filters).filter(([_, v]) => v !== ''))
  params.page = page; params.limit = 10
  try {
    const { data } = await api.get('/api/procurement/transactions', { params })
    results.value = data
  } catch { results.value = { data: [], total_count: 0, total_pages: 0, page: 1 } }
  finally { loading.value = false }
}

function goToPage(page) { search(page) }

function clearFilters() {
  Object.keys(filters).forEach(k => filters[k] = '')
  search()
}

function exportCSV() {
  const rows = results.value.data
  if (!rows.length) return
  const headers = ['request_number','title','department','supplier_name','po_number','status','created_at']
  const lines = [headers.join(','), ...rows.map(r => headers.map(h => `"${r[h] ?? ''}"`).join(','))]
  const blob = new Blob([lines.join('\n')], { type: 'text/csv' })
  const a = document.createElement('a')
  a.href = URL.createObjectURL(blob); a.download = 'transactions.csv'; a.click()
}

function formatDate(d) { if (!d) return '—'; return new Date(d).toLocaleDateString('en-GB') }

onMounted(() => search())
</script>