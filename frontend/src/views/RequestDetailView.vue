<template>
  <div class="detail-page">

    <div v-if="loading" class="loading">Loading...</div>

    <div v-else-if="!request" class="loading">Request not found.</div>

    <template v-else>

      <!-- Page header -->
      <div class="page-header">
        <button class="back-btn" @click="router.push('/requests')">← Back</button>
        <div>
          <h1>{{ request.title }}</h1>
          <span class="req-num">{{ request.request_number }}</span>
        </div>
        <StatusBadge :status="request.status" />
      </div>

      <!-- ── Action buttons — shown based on role + status ── -->
      <div class="action-bar">

        <!-- Officer sees this when status is draft -->
        <button
          v-if="auth.isOfficer && request.status === 'draft'"
          class="btn blue"
          @click="sendToFinance"
          :disabled="acting"
        >
          📤 Send to Finance
        </button>

        <!-- Manager sees these when budget is approved -->
        <template v-if="auth.isManager && request.status === 'budget_approved'">
          <button class="btn green" @click="openApproveModal" :disabled="acting">
            ✓ Approve
          </button>
          <button class="btn red" @click="openRejectModal" :disabled="acting">
            ✕ Reject
          </button>
        </template>

        <!-- Officer sees this when request is open for bidding or a supplier has been selected -->
        <button
          v-if="auth.isOfficer && (request.status === 'open_for_bidding' || request.status === 'supplier_selected')"
          class="btn purple"
          @click="router.push(`/bids/${request.id}`)"
        >
          🏷 View Bids →
        </button>

                <!-- Edit — only on draft requests -->
        <button
          v-if="request.status === 'draft'"
          class="btn grey"
          @click="openEditModal"
        >✏ Edit</button>

        <!-- Delete — only on draft requests -->
        <button
          v-if="request.status === 'draft'"
          class="btn danger"
          @click="deleteRequest"
          :disabled="acting"
        >🗑 Delete</button>

      </div>

      <div class="content-grid">

        <!-- ── Left column: request info + items ── -->
        <div class="left">

          <!-- Request info card -->
          <div class="card">
            <h2>Request Information</h2>
            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">Department</span>
                <span>{{ request.department || '—' }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Required Date</span>
                <span>{{ formatDate(request.required_date) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Estimated Total</span>
                <span class="amount">LKR {{ formatCurrency(request.estimated_total) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">Created On</span>
                <span>{{ formatDate(request.created_at) }}</span>
              </div>
            </div>
            <div v-if="request.description" class="description">
              <span class="info-label">Description</span>
              <p>{{ request.description }}</p>
            </div>
          </div>

          <!-- Items table -->
          <div class="card">
            <h2>Items ({{ request.items?.length || 0 }})</h2>
            <table>
              <thead>
                <tr>
                  <th>Item Name</th>
                  <th>Qty</th>
                  <th>Unit Price</th>
                  <th>Subtotal</th>
                  <th>Category</th>
                  <th>Type</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in request.items" :key="item.id">
                  <td>{{ item.item_name }}</td>
                  <td>{{ item.quantity }}</td>
                  <td>{{ formatCurrency(item.estimated_price) }}</td>
                  <td>{{ formatCurrency(item.quantity * item.estimated_price) }}</td>
                  <td>{{ item.category || '—' }}</td>
                  <td>
                    <span :class="['type-badge', item.item_type]">
                      {{ item.item_type }}
                    </span>
                  </td>
                </tr>
                <tr v-if="!request.items?.length">
                  <td colspan="6" class="empty">No items</td>
                </tr>
              </tbody>
            </table>
          </div>

        </div>

        <!-- ── Right column: timeline + approval info ── -->
        <div class="right">

          <!-- Status timeline -->
          <div class="card">
            <h2>Status Timeline</h2>
            <div class="timeline">
              <div
                v-for="step in timeline"
                :key="step.status"
                :class="['timeline-step', getStepState(step.status)]"
              >
                <div class="step-dot"></div>
                <div class="step-info">
                  <span class="step-label">{{ step.label }}</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Approval info — shown only if approved or rejected -->
          <div class="card"
            v-if="request.approved_at || request.approval_remarks"
          >
            <h2>Approval Details</h2>
            <div class="info-item">
              <span class="info-label">Decision Date</span>
              <span>{{ formatDate(request.approved_at) }}</span>
            </div>
            <div class="info-item" v-if="request.approval_remarks">
              <span class="info-label">Remarks</span>
              <span>{{ request.approval_remarks }}</span>
            </div>
          </div>

        </div>
      </div>
    </template>

    <!-- ── Approve modal ── -->
    <div v-if="showApproveModal" class="modal-overlay">
      <div class="modal">
        <h3>Approve Request</h3>
        <p>Add optional remarks before approving.</p>
        <textarea v-model="remarksInput" rows="3" placeholder="Remarks (optional)"></textarea>
        <div class="modal-actions">
          <button class="cancel-btn" @click="showApproveModal = false">Cancel</button>
          <button class="btn green" @click="confirmApprove" :disabled="acting">
            Confirm Approve
          </button>
        </div>
      </div>
    </div>

    <!-- ── Reject modal ── -->
    <div v-if="showRejectModal" class="modal-overlay">
      <div class="modal">
        <h3>Reject Request</h3>
        <p>You must provide a reason for rejection.</p>
        <textarea v-model="remarksInput" rows="3" placeholder="Reason for rejection *" required></textarea>
        <div class="modal-actions">
          <button class="cancel-btn" @click="showRejectModal = false">Cancel</button>
          <button class="btn red" @click="confirmReject" :disabled="acting || !remarksInput">
            Confirm Reject
          </button>
        </div>
      </div>
    </div>

  </div>

    <!-- Edit modal -->
    <div v-if="showEditModal" class="modal-overlay">
      <div class="modal">
        <h3>Edit Request</h3>
        <div class="form-group">
          <label>Title</label>
          <input v-model="editForm.title" type="text" />
        </div>
        <div class="form-group">
          <label>Description</label>
          <textarea v-model="editForm.description" rows="3"></textarea>
        </div>
        <div class="form-group">
          <label>Department</label>
          <input v-model="editForm.department" type="text" />
        </div>
        <div class="form-group">
          <label>Required Date</label>
          <input v-model="editForm.required_date" type="date" />
        </div>
        <div class="form-group">
          <label>Estimated Total (LKR)</label>
          <input v-model.number="editForm.estimated_total" type="number" />
        </div>
        <div class="modal-actions">
          <button class="cancel-btn" @click="showEditModal = false">Cancel</button>
          <button class="btn blue" @click="confirmEdit" :disabled="acting">Save Changes</button>
        </div>
      </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../api'
import { useAuthStore } from '../stores/auth'
import StatusBadge from '../components/StatusBadge.vue'

const route  = useRoute()
const router = useRouter()
const auth   = useAuthStore()

const request         = ref(null)
const loading         = ref(false)
const acting          = ref(false)   // true while an API call is in progress
const showApproveModal = ref(false)
const showRejectModal  = ref(false)
const remarksInput     = ref('')

// Edit & Delete
const showEditModal = ref(false)
const editForm = ref({
  title: '', description: '', department: '',
  required_date: '', estimated_total: 0, status: 'draft'
})

function openEditModal() {
  // Pre-fill form with current values
  editForm.value = {
    title:           request.value.title,
    description:     request.value.description || '',
    department:      request.value.department || '',
    required_date:   request.value.required_date || '',
    estimated_total: request.value.estimated_total,
    status:          request.value.status,
  }
  showEditModal.value = true
}

async function confirmEdit() {
  acting.value = true
  await api.put(`/api/procurement/requests/${request.value.id}`, editForm.value)
  showEditModal.value = false
  await fetchRequest()   // reload to show updated values
  acting.value = false
}

async function deleteRequest() {
  if (!confirm('Delete this request? This cannot be undone.')) return
  acting.value = true
  await api.delete(`/api/procurement/requests/${request.value.id}`)
  router.push('/requests')   // go back to list after deletion
}
// The ordered list of statuses for the timeline
const timeline = [
  { status: 'draft',                      label: 'Draft Created' },
  { status: 'sent_to_finance',            label: 'Sent to Finance' },
  { status: 'budget_approved',            label: 'Budget Approved' },
  { status: 'manager_approved',           label: 'Manager Approved' },
  { status: 'open_for_bidding',           label: 'Open for Bidding' },
  { status: 'supplier_selected',          label: 'Supplier Selected' },
  { status: 'purchase_order_created',     label: 'PO Created' },
  { status: 'goods_received',             label: 'Goods Received' },
  { status: 'completed',                  label: 'Completed' },
]

// Order map so we can compare where we are in the flow
const statusOrder = {
  draft: 0, sent_to_finance: 1, budget_approved: 2,
  budget_rejected: 2, manager_approved: 3, manager_rejected: 3,
  open_for_bidding: 4, supplier_selected: 5, purchase_order_created: 6,
  goods_received: 7, sent_to_asset_module: 7,
  sent_to_finance_for_payment: 8, completed: 8,
}

function getStepState(stepStatus) {
  const current = statusOrder[request.value?.status] ?? 0
  const step    = statusOrder[stepStatus] ?? 0
  if (step < current)  return 'done'
  if (step === current) return 'active'
  return 'pending'
}

// Fetch the request (includes items — backend returns them together)
async function fetchRequest() {
  loading.value = true
  try {
    const { data } = await api.get(`/api/procurement/requests/${route.params.id}`)
    request.value = data
  } catch {
    request.value = null
  } finally {
    loading.value = false
  }
}

// ── Actions ──

async function sendToFinance() {
  acting.value = true
  await api.post(`/api/procurement/requests/${request.value.id}/send-to-finance`)
  await fetchRequest()   // reload to show updated status
  acting.value = false
}

function openApproveModal() { remarksInput.value = ''; showApproveModal.value = true }
function openRejectModal()  { remarksInput.value = ''; showRejectModal.value  = true }

async function confirmApprove() {
  acting.value = true
  await api.post(`/api/procurement/requests/${request.value.id}/approve`,
    { remarks: remarksInput.value })
  showApproveModal.value = false
  await fetchRequest()
  acting.value = false
}

async function confirmReject() {
  if (!remarksInput.value) return
  acting.value = true
  await api.post(`/api/procurement/requests/${request.value.id}/reject`,
    { remarks: remarksInput.value })
  showRejectModal.value = false
  await fetchRequest()
  acting.value = false
}

function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-GB')
}
function formatCurrency(n) {
  return Number(n || 0).toLocaleString('en-LK')
}

onMounted(fetchRequest)
</script>

<style scoped>

.btn.grey   { background: #f3f4f6; color: #374151; }
.btn.danger { background: #fee2e2; color: #991b1b; }
.form-group { display: flex; flex-direction: column; gap: 5px; margin-bottom: 14px; }
.form-group label { font-size: 13px; font-weight: 500; color: #374151; }
.form-group input,
.form-group textarea { padding: 9px 12px; border: 1px solid #d1d5db;
  border-radius: 6px; font-size: 14px; outline: none; }
.form-group input:focus,
.form-group textarea:focus { border-color: #4f46e5; }
.detail-page { padding: 24px; max-width: 1100px; margin: 0 auto; }
.loading { text-align: center; padding: 40px; color: #6b7280; }

.page-header { display: flex; align-items: center; gap: 16px; margin-bottom: 16px; flex-wrap: wrap; }
h1 { font-size: 20px; font-weight: 600; }
h2 { font-size: 15px; font-weight: 600; margin-bottom: 14px; }
.req-num { font-size: 13px; color: #9ca3af; }
.back-btn { background: none; border: 1px solid #d1d5db; border-radius: 6px;
  padding: 6px 12px; cursor: pointer; font-size: 13px; flex-shrink: 0; }

.action-bar { display: flex; gap: 10px; margin-bottom: 20px; flex-wrap: wrap; }

.btn { padding: 9px 20px; border: none; border-radius: 6px; font-size: 14px;
  font-weight: 500; cursor: pointer; }
.btn:disabled { opacity: 0.6; cursor: not-allowed; }
.btn.blue   { background: #dbeafe; color: #1d4ed8; }
.btn.green  { background: #dcfce7; color: #166534; }
.btn.red    { background: #fee2e2; color: #991b1b; }
.btn.purple { background: #ede9fe; color: #5b21b6; }

.content-grid { display: grid; grid-template-columns: 1fr 300px; gap: 20px; }
@media(max-width:800px){ .content-grid { grid-template-columns: 1fr; } }
.left, .right { display: flex; flex-direction: column; gap: 20px; }

.card { background: #fff; border: 1px solid #e5e7eb; border-radius: 10px; padding: 20px; }

.info-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; margin-bottom: 14px; }
.info-item { display: flex; flex-direction: column; gap: 3px; }
.info-label { font-size: 12px; color: #9ca3af; text-transform: uppercase; }
.amount { font-weight: 700; font-size: 16px; color: #111827; }
.description p { font-size: 14px; color: #374151; margin-top: 4px; line-height: 1.6; }

table { width: 100%; border-collapse: collapse; font-size: 13px; }
th { text-align: left; padding: 8px 10px; color: #6b7280;
  background: #f9fafb; border-bottom: 1px solid #e5e7eb; font-size: 11px; text-transform: uppercase; }
td { padding: 10px; border-bottom: 1px solid #f3f4f6; }
.empty { text-align: center; color: #9ca3af; padding: 20px; }

.type-badge { padding: 2px 8px; border-radius: 10px; font-size: 11px; font-weight: 500; }
.type-badge.asset      { background: #dbeafe; color: #1d4ed8; }
.type-badge.consumable { background: #f0fdf4; color: #166534; }

/* Timeline */
.timeline { display: flex; flex-direction: column; gap: 0; }
.timeline-step { display: flex; align-items: center; gap: 10px; padding: 8px 0;
  position: relative; }
.timeline-step:not(:last-child)::after {
  content: ''; position: absolute; left: 7px; top: 26px;
  width: 2px; height: 100%; background: #e5e7eb; z-index: 0;
}
.step-dot { width: 16px; height: 16px; border-radius: 50%; flex-shrink: 0;
  border: 2px solid #d1d5db; background: #fff; z-index: 1; }
.timeline-step.done   .step-dot { background: #16a34a; border-color: #16a34a; }
.timeline-step.active .step-dot { background: #4f46e5; border-color: #4f46e5; }
.step-label { font-size: 13px; color: #6b7280; }
.timeline-step.done   .step-label { color: #16a34a; font-weight: 500; }
.timeline-step.active .step-label { color: #4f46e5; font-weight: 600; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4);
  display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: #fff; border-radius: 12px; padding: 28px;
  width: 420px; max-width: 90vw; }
.modal h3 { font-size: 16px; font-weight: 600; margin-bottom: 8px; }
.modal p  { font-size: 13px; color: #6b7280; margin-bottom: 12px; }
.modal textarea { width: 100%; padding: 10px; border: 1px solid #d1d5db;
  border-radius: 6px; font-size: 14px; resize: vertical; }
.modal-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 16px; }
.cancel-btn { padding: 8px 16px; border: 1px solid #d1d5db; border-radius: 6px;
  background: #fff; cursor: pointer; font-size: 13px; }
</style>
