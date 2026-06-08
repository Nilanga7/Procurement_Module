<template>
  <div class="p-10 max-w-5xl mx-auto">
    <!-- Header -->
    <div class="flex items-start gap-4 mb-8">
      <button @click="router.push('/requests')"
        class="flex items-center gap-1 px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
               text-[#474651] hover:bg-[#eff4ff] transition-colors flex-shrink-0">
        <span class="material-symbols-outlined text-[16px]">arrow_back</span> Back
      </button>
      <div>
        <h1 class="text-[24px] font-bold text-[#0b1c30]">Bid Evaluation</h1>
        <p v-if="requestInfo" class="text-[13px] text-[#777682] mt-1">
          {{ requestInfo.request_number }} — {{ requestInfo.title }}
        </p>
      </div>
    </div>

    <div v-if="loading" class="bg-white border border-[#c8c5d3] rounded-xl p-16 text-center ambient-shadow">
      <span class="material-symbols-outlined text-[40px] text-[#c8c5d3] block mb-3">hourglass_empty</span>
      <p class="text-[14px] text-[#777682]">Loading bids…</p>
    </div>

    <div v-else-if="bids.length === 0"
      class="bg-white border border-[#c8c5d3] rounded-xl p-16 text-center ambient-shadow">
      <span class="material-symbols-outlined text-[48px] text-[#c8c5d3] block mb-3">gavel</span>
      <p class="text-[16px] font-semibold text-[#474651]">No bids yet</p>
      <p class="text-[13px] text-[#777682] mt-1">No bids have been submitted for this request yet.</p>
    </div>

    <template v-else>
      <!-- Legend -->
      <div class="flex gap-4 mb-4">
        <span class="flex items-center gap-1.5 text-[12px] text-[#006c49] font-semibold">
          <span class="w-2 h-2 rounded-full bg-[#006c49]"></span> Lowest Price
        </span>
        <span class="flex items-center gap-1.5 text-[12px] text-[#1a146b] font-semibold">
          <span class="w-2 h-2 rounded-full bg-[#1a146b]"></span> Fastest Delivery
        </span>
      </div>

      <!-- Comparison table -->
      <div class="bg-white border border-[#c8c5d3] rounded-xl ambient-shadow overflow-hidden mb-6">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead class="bg-[#d3e4fe]/30">
              <tr>
                <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Supplier</th>
                <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Offered Price</th>
                <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Delivery</th>
                <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Notes</th>
                <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Status</th>
                <th class="px-6 py-4 text-[12px] font-semibold tracking-wide text-[#474651] uppercase">Actions</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-[#c8c5d3]/30">
              <tr v-for="bid in bids" :key="bid.id"
                :class="bid.status === 'accepted' ? 'bg-[#6cf8bb]/10' :
                        bid.status === 'rejected' ? 'opacity-50 bg-[#f8f9ff]' : ''">
                <td class="px-6 py-4">
                  <div class="text-[14px] font-semibold text-[#0b1c30]">{{ bid.supplier_name }}</div>
                  <div class="text-[12px] text-[#777682]">{{ bid.supplier_email }}</div>
                </td>
                <td class="px-6 py-4">
                  <span :class="bid.offered_price === lowestPrice ? 'text-[#006c49] font-bold text-[15px]' : 'text-[14px] text-[#0b1c30]'">
                    LKR {{ formatCurrency(bid.offered_price) }}
                  </span>
                  <span v-if="bid.offered_price === lowestPrice"
                    class="ml-2 px-2 py-0.5 bg-[#6cf8bb] text-[#00714d] text-[10px] font-bold rounded-full">LOWEST</span>
                </td>
                <td class="px-6 py-4">
                  <span :class="bid.delivery_days === fastestDelivery ? 'text-[#1a146b] font-bold' : 'text-[14px] text-[#0b1c30]'">
                    {{ bid.delivery_days }} days
                  </span>
                  <span v-if="bid.delivery_days === fastestDelivery"
                    class="ml-2 px-2 py-0.5 bg-[#c3c0ff] text-[#1a146b] text-[10px] font-bold rounded-full">FASTEST</span>
                </td>
                <td class="px-6 py-4 text-[13px] text-[#777682] max-w-[180px]">{{ bid.notes || '—' }}</td>
                <td class="px-6 py-4"><StatusBadge :status="bid.status" /></td>
                <td class="px-6 py-4">
                  <div v-if="!decisionMade" class="flex gap-2">
                    <button @click="acceptBid(bid)" :disabled="acting"
                      class="px-3 py-1.5 bg-[#6cf8bb]/30 border border-[#006c49]/30 text-[#006c49]
                             rounded-lg text-[12px] font-bold hover:bg-[#6cf8bb]/60 disabled:opacity-40 transition-colors">
                      Accept
                    </button>
                    <button @click="rejectBid(bid.id)" :disabled="acting"
                      class="px-3 py-1.5 bg-[#ffdad6] border border-[#ba1a1a]/30 text-[#93000a]
                             rounded-lg text-[12px] font-bold hover:bg-[#ffdad6]/70 disabled:opacity-40 transition-colors">
                      Reject
                    </button>
                  </div>
                  <span v-else class="text-[#c8c5d3] text-[13px]">—</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Accepted response preview -->
      <div v-if="acceptedBid" class="bg-white border border-[#c8c5d3] rounded-xl p-6 ambient-shadow">
        <h2 class="text-[16px] font-bold text-[#0b1c30] mb-2">Responses Sent</h2>
        <p class="text-[13px] text-[#777682] mb-4">Message sent to suppliers.</p>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div class="bg-[#6cf8bb]/20 border border-[#006c49]/30 rounded-xl p-4">
            <div class="flex justify-between items-center mb-2">
              <span class="text-[11px] font-bold text-[#006c49] tracking-wider">ACCEPTED</span>
              <span class="text-[12px] text-[#777682]">To: {{ acceptedBid.supplier_name }}</span>
            </div>
            <p class="text-[13px] text-[#0b1c30]">Your bid has been accepted. A purchase order will be placed soon.</p>
          </div>
          <div class="bg-[#ffdad6]/30 border border-[#ba1a1a]/20 rounded-xl p-4">
            <div class="flex justify-between items-center mb-2">
              <span class="text-[11px] font-bold text-[#93000a] tracking-wider">REJECTED</span>
              <span class="text-[12px] text-[#777682]">To: {{ bids.length - 1 }} other supplier(s)</span>
            </div>
            <p class="text-[13px] text-[#0b1c30]">Thank you for your bid. Your offer was not selected.</p>
          </div>
        </div>
        <button @click="router.push('/purchase-orders')"
          class="flex items-center gap-2 px-6 py-3 bg-[#1a146b] text-white rounded-xl
                 text-[13px] font-bold hover:opacity-90 transition-all">
          <span class="material-symbols-outlined text-[18px]">shopping_cart</span>
          Go to Purchase Orders to create PO
        </button>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import StatusBadge from '../components/StatusBadge.vue'

const route  = useRoute()
const router = useRouter()
const bids        = ref([])
const requestInfo = ref(null)
const loading     = ref(false)
const acting      = ref(false)

async function fetchBids() {
  loading.value = true
  try {
    const [bidsRes, reqRes] = await Promise.all([
      api.get(`/api/procurement/requests/${route.params.requestId}/bids`),
      api.get(`/api/procurement/requests/${route.params.requestId}`),
    ])
    bids.value        = bidsRes.data
    requestInfo.value = reqRes.data
  } finally { loading.value = false }
}

const lowestPrice    = computed(() => bids.value.length ? Math.min(...bids.value.map(b => b.offered_price)) : null)
const fastestDelivery = computed(() => bids.value.length ? Math.min(...bids.value.map(b => b.delivery_days)) : null)
const decisionMade   = computed(() => bids.value.some(b => b.status === 'accepted'))
const acceptedBid    = computed(() => bids.value.find(b => b.status === 'accepted') || null)

async function acceptBid(bid) {
  if (!confirm(`Accept bid from ${bid.supplier_name}? All other bids will be rejected.`)) return
  acting.value = true
  try { await api.post(`/api/procurement/bids/${bid.id}/accept`); await fetchBids() }
  finally { acting.value = false }
}

async function rejectBid(bidId) {
  acting.value = true
  try { await api.post(`/api/procurement/bids/${bidId}/reject`); await fetchBids() }
  finally { acting.value = false }
}

function formatCurrency(n) { return Number(n || 0).toLocaleString('en-LK') }

onMounted(fetchBids)
</script>