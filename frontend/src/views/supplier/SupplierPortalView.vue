<template>
  <!-- Standalone page — no sidebar, own top bar -->
  <div class="min-h-screen bg-[#f8f9ff]">

    <!-- Top bar -->
    <nav class="fixed top-0 left-0 w-full z-50 flex justify-between items-center px-10 h-16
                bg-[#f8f9ff] border-b border-[#c8c5d3]">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 bg-[#1a146b] rounded-lg flex items-center justify-center">
          <span class="material-symbols-outlined text-white text-[16px]"
                style="font-variation-settings:'FILL' 1">storefront</span>
        </div>
        <span class="text-[16px] font-bold text-[#1a146b]">Supplier Portal</span>
      </div>
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2 px-3 py-1.5 rounded-full bg-[#eff4ff]">
          <div class="w-6 h-6 rounded-full bg-[#1a146b] flex items-center justify-center">
            <span class="text-white text-[10px] font-bold">{{ initials }}</span>
          </div>
          <span class="text-[12px] font-semibold text-[#1a146b]">{{ auth.user?.name }}</span>
        </div>
        <button @click="handleLogout"
          class="px-3 py-1.5 rounded-lg border border-[#c8c5d3] text-[12px] font-semibold
                 text-[#474651] hover:bg-[#eff4ff] transition-colors">
          Logout
        </button>
      </div>
    </nav>

    <!-- Page content -->
    <main class="pt-16 p-10 max-w-7xl mx-auto">

      <!-- Welcome banner -->
      <div class="mb-8 bg-gradient-to-r from-[#1a146b] to-[#3e3c8f] rounded-2xl p-6 text-white">
        <p class="text-[12px] font-semibold tracking-wide opacity-70 mb-1">Welcome back</p>
        <h1 class="text-[24px] font-bold">{{ auth.user?.name }}</h1>
        <p class="text-[14px] opacity-70 mt-1">Browse open procurement requests and manage your bids below.</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">

        <!-- ── LEFT: Open requests ── -->
        <div>
          <div class="flex items-center justify-between mb-4">
            <div>
              <h2 class="text-[18px] font-bold text-[#0b1c30]">Open for Bidding</h2>
              <p class="text-[13px] text-[#777682]">These requests are accepting bids right now.</p>
            </div>
            <span class="px-2.5 py-1 bg-[#3e3c8f] text-[#e2dfff] text-[11px] font-bold rounded-full">
              {{ openRequests.length }} active
            </span>
          </div>

          <div v-if="loading" class="bg-white border border-[#c8c5d3] rounded-xl p-12 text-center ambient-shadow">
            <span class="material-symbols-outlined text-[36px] text-[#c8c5d3] block mb-2">hourglass_empty</span>
            <p class="text-[14px] text-[#777682]">Loading requests…</p>
          </div>

          <div v-else-if="openRequests.length === 0"
            class="bg-white border border-[#c8c5d3] rounded-xl p-12 text-center ambient-shadow">
            <span class="material-symbols-outlined text-[40px] text-[#c8c5d3] block mb-2">gavel</span>
            <p class="text-[14px] text-[#474651] font-semibold">No open requests</p>
            <p class="text-[13px] text-[#777682] mt-1">No procurement requests are currently open for bidding.</p>
          </div>

          <div v-for="req in openRequests" :key="req.id"
            class="bg-white border border-[#c8c5d3] rounded-xl p-5 mb-3 ambient-shadow hover-shadow transition-all">
            <div class="flex justify-between items-start mb-3">
              <div>
                <span class="text-[11px] font-bold text-[#3e3c8f] bg-[#e2dfff]/50 px-2 py-0.5 rounded-full">
                  {{ req.request_number }}
                </span>
                <h3 class="text-[15px] font-bold text-[#0b1c30] mt-1.5">{{ req.title }}</h3>
              </div>
              <button
                :disabled="alreadyBid(req.id)"
                @click="openBidForm(req.id)"
                :class="alreadyBid(req.id)
                  ? 'bg-[#6cf8bb]/30 text-[#006c49] border-[#006c49]/30 cursor-not-allowed'
                  : 'bg-[#1a146b] text-white hover:opacity-90'"
                class="flex-shrink-0 flex items-center gap-1.5 px-4 py-2 rounded-xl text-[12px]
                       font-bold border transition-all">
                <span class="material-symbols-outlined text-[16px]">
                  {{ alreadyBid(req.id) ? 'check_circle' : 'add_circle' }}
                </span>
                {{ alreadyBid(req.id) ? 'Bid Submitted' : 'Submit Bid' }}
              </button>
            </div>
            <div class="flex flex-wrap gap-3">
              <span v-if="req.department"
                class="flex items-center gap-1 text-[12px] text-[#474651] bg-[#f8f9ff] px-2.5 py-1 rounded-full border border-[#c8c5d3]">
                <span class="material-symbols-outlined text-[14px]">business</span>
                {{ req.department }}
              </span>
              <span v-if="req.required_date"
                class="flex items-center gap-1 text-[12px] text-[#474651] bg-[#f8f9ff] px-2.5 py-1 rounded-full border border-[#c8c5d3]">
                <span class="material-symbols-outlined text-[14px]">calendar_today</span>
                Required by {{ formatDate(req.required_date) }}
              </span>
              <span class="flex items-center gap-1 text-[12px] text-[#474651] bg-[#f8f9ff] px-2.5 py-1 rounded-full border border-[#c8c5d3]">
                <span class="material-symbols-outlined text-[14px]">payments</span>
                Est. LKR {{ formatCurrency(req.estimated_total) }}
              </span>
            </div>
          </div>
        </div>

        <!-- ── RIGHT: My submitted bids ── -->
        <div>
          <div class="flex items-center justify-between mb-4">
            <div>
              <h2 class="text-[18px] font-bold text-[#0b1c30]">My Submitted Bids</h2>
              <p class="text-[13px] text-[#777682]">Track the status of your submissions.</p>
            </div>
            <span class="px-2.5 py-1 bg-[#e5eeff] text-[#1a146b] text-[11px] font-bold rounded-full">
              {{ myBids.length }} bid{{ myBids.length !== 1 ? 's' : '' }}
            </span>
          </div>

          <div v-if="myBids.length === 0"
            class="bg-white border border-[#c8c5d3] rounded-xl p-12 text-center ambient-shadow">
            <span class="material-symbols-outlined text-[40px] text-[#c8c5d3] block mb-2">receipt_long</span>
            <p class="text-[14px] text-[#474651] font-semibold">No bids yet</p>
            <p class="text-[13px] text-[#777682] mt-1">Submit a bid on an open request to get started.</p>
          </div>

          <div v-for="bid in myBids" :key="bid.id"
            :class="bid.status === 'accepted' ? 'border-[#006c49]/40 bg-[#6cf8bb]/5' :
                    bid.status === 'rejected' ? 'opacity-60' : ''"
            class="bg-white border border-[#c8c5d3] rounded-xl p-5 mb-3 ambient-shadow transition-all">

            <div class="flex justify-between items-center mb-3">
              <span class="text-[11px] font-bold text-[#3e3c8f] bg-[#e2dfff]/50 px-2 py-0.5 rounded-full">
                {{ bid.request_number || 'Request' }}
              </span>
              <span :class="bid.status === 'accepted' ? 'bg-[#6cf8bb] text-[#00714d]' :
                             bid.status === 'rejected' ? 'bg-[#d3e4fe] text-[#474651]' :
                             'bg-[#ffddb8] text-[#653e00]'"
                class="px-2.5 py-1 rounded-full text-[11px] font-bold">
                {{ statusLabel(bid.status) }}
              </span>
            </div>

            <div class="grid grid-cols-3 gap-3 mb-3">
              <div class="bg-[#f8f9ff] rounded-lg p-3 text-center">
                <p class="text-[11px] text-[#777682] font-semibold mb-0.5">Your Price</p>
                <p class="text-[13px] font-bold text-[#1a146b]">LKR {{ formatCurrency(bid.offered_price) }}</p>
              </div>
              <div class="bg-[#f8f9ff] rounded-lg p-3 text-center">
                <p class="text-[11px] text-[#777682] font-semibold mb-0.5">Delivery</p>
                <p class="text-[13px] font-bold text-[#0b1c30]">{{ bid.delivery_days }}d</p>
              </div>
              <div class="bg-[#f8f9ff] rounded-lg p-3 text-center">
                <p class="text-[11px] text-[#777682] font-semibold mb-0.5">Submitted</p>
                <p class="text-[13px] font-bold text-[#0b1c30]">{{ formatDate(bid.created_at) }}</p>
              </div>
            </div>

            <div v-if="bid.status === 'accepted'"
              class="flex items-start gap-2 px-4 py-3 bg-[#6cf8bb]/20 border border-[#006c49]/20 rounded-xl text-[13px] text-[#006c49]">
              <span class="material-symbols-outlined text-[16px] flex-shrink-0 mt-0.5"
                    style="font-variation-settings:'FILL' 1">celebration</span>
              Congratulations! Your bid has been accepted. A purchase order will be created soon.
            </div>
            <div v-if="bid.status === 'rejected'"
              class="px-4 py-3 bg-[#d3e4fe] border border-[#c8c5d3] rounded-xl text-[13px] text-[#474651]">
              Thank you for your bid. Your offer was not selected this time.
            </div>
          </div>
        </div>

      </div>
    </main>

    <!-- ── Bid submission modal ── -->
    <Teleport to="body">
      <div v-if="showModal"
        class="fixed inset-0 z-[100] flex items-center justify-center bg-black/40 backdrop-blur-sm"
        @click.self="showModal = false">
        <div class="bg-white w-full max-w-md mx-4 rounded-2xl border border-[#c8c5d3] overflow-hidden"
             style="box-shadow:0 25px 60px rgba(0,0,0,0.15)">
          <!-- Modal header -->
          <div class="flex justify-between items-center px-6 py-4 border-b border-[#c8c5d3] bg-[#eff4ff]">
            <div>
              <h3 class="text-[16px] font-bold text-[#0b1c30]">Submit Your Bid</h3>
              <p class="text-[12px] text-[#777682]">Request: <strong>{{ selectedRequestNumber }}</strong></p>
            </div>
            <button @click="showModal = false"
              class="p-1.5 rounded-lg hover:bg-[#dce9ff] transition-colors text-[#474651]">
              <span class="material-symbols-outlined text-[20px]">close</span>
            </button>
          </div>

          <!-- Modal body -->
          <div class="px-6 py-5 space-y-4">
            <div>
              <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Offered Price (LKR) *</label>
              <input v-model.number="bidForm.offered_price" type="number" min="0" step="0.01"
                placeholder="e.g. 240000"
                class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                       focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
            </div>
            <div>
              <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Delivery Time (days) *</label>
              <input v-model.number="bidForm.delivery_days" type="number" min="1"
                placeholder="e.g. 14"
                class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                       focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
            </div>
            <div>
              <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Notes</label>
              <textarea v-model="bidForm.notes" rows="3"
                placeholder="Any additional information about your offer..."
                class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                       focus:outline-none focus:border-[#1a146b] resize-none"></textarea>
            </div>

            <div v-if="submitError"
              class="flex items-center gap-2 px-4 py-3 bg-[#ffdad6] border border-[#ba1a1a]/30
                     rounded-lg text-[13px] text-[#93000a]">
              <span class="material-symbols-outlined text-[16px]">error</span>
              {{ submitError }}
            </div>
          </div>

          <!-- Modal footer -->
          <div class="flex gap-3 px-6 py-4 bg-[#f8f9ff] border-t border-[#c8c5d3]">
            <button @click="showModal = false"
              class="flex-1 py-2.5 border border-[#c8c5d3] rounded-xl text-[13px] font-semibold
                     text-[#474651] hover:bg-[#eff4ff] transition-colors">Cancel</button>
            <button :disabled="submitting || !bidForm.offered_price || !bidForm.delivery_days"
              @click="submitBid"
              class="flex-1 py-2.5 bg-[#1a146b] text-white rounded-xl text-[13px] font-bold
                     hover:opacity-90 disabled:opacity-50 transition-all flex items-center justify-center gap-2">
              <svg v-if="submitting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
              </svg>
              {{ submitting ? 'Submitting…' : 'Submit Bid' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../../api'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const auth   = useAuthStore()

const openRequests = ref([])
const myBids       = ref([])
const loading      = ref(false)
const submitting   = ref(false)
const submitError  = ref('')
const showModal    = ref(false)

const bidForm = ref({ request_id: '', offered_price: 0, delivery_days: 0, notes: '' })

const initials = computed(() => {
  const name = auth.user?.name || ''
  return name.split(' ').map(w => w[0]).join('').toUpperCase().slice(0, 2)
})

async function loadData() {
  loading.value = true
  try {
    const [reqRes, bidRes] = await Promise.all([
      api.get('/api/procurement/requests', { params: { status: 'open_for_bidding' } }),
      api.get('/api/supplier/bids'),
    ])
    openRequests.value = Array.isArray(reqRes.data) ? reqRes.data : []
    myBids.value       = Array.isArray(bidRes.data)  ? bidRes.data  : []
  } catch {
    openRequests.value = []
    myBids.value       = []
  } finally {
    loading.value = false
  }
}

function openBidForm(requestId) {
  bidForm.value = { request_id: requestId, offered_price: 0, delivery_days: 0, notes: '' }
  submitError.value = ''
  showModal.value = true
}

function alreadyBid(requestId) {
  return myBids.value.some(b => b.request_id === requestId)
}

const selectedRequestNumber = computed(() => {
  const req = openRequests.value.find(r => r.id === bidForm.value.request_id)
  return req?.request_number || ''
})

async function submitBid() {
  submitting.value  = true
  submitError.value = ''
  try {
    await api.post('/api/supplier/bids', bidForm.value)
    showModal.value = false
    await loadData()
  } catch (e) {
    submitError.value = e.response?.data?.error || 'Failed to submit bid. Please try again.'
  } finally {
    submitting.value = false
  }
}

function statusLabel(status) {
  return { pending: 'Under Review', accepted: 'Accepted ✓', rejected: 'Not Selected' }[status] || status
}

function formatDate(d) { if (!d) return '—'; return new Date(d).toLocaleDateString('en-GB') }
function formatCurrency(n) { return Number(n || 0).toLocaleString('en-LK') }

function handleLogout() { auth.logout(); router.push('/login') }

onMounted(loadData)
</script>