<template>
  <div class="po-page">

    <!-- ══════════════════════════════════════
         DETAIL VIEW — shown when a PO is selected
         ══════════════════════════════════════ -->
    <div v-if="selectedPO">

      <div class="page-header">
        <button class="back-btn" @click="selectedPO = null">← Back to List</button>
        <div>
          <h1>{{ selectedPO.po_number }}</h1>
          <p class="subtitle">{{ selectedPO.request_title }}</p>
        </div>
        <span :class="['po-status', selectedPO.status]">
          {{ statusLabel(selectedPO.status) }}
        </span>
      </div>

      <div class="detail-grid">

        <!-- Left: items table -->
        <div class="left">
          <div class="card">
            <h2>Items</h2>
            <table>
              <thead>
                <tr>
                  <th>Item Name</th>
                  <th>Qty</th>
                  <th>Unit Price (LKR)</th>
                  <th>Subtotal (LKR)</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in selectedPO.items" :key="item.id">
                  <td>{{ item.item_name }}</td>
                  <td>{{ item.quantity }}</td>
                  <td>{{ formatCurrency(item.unit_price) }}</td>
                  <td>{{ formatCurrency(item.subtotal) }}</td>
                </tr>
                <tr v-if="!selectedPO.items?.length">
                  <td colspan="4" class="empty">No items</td>
                </tr>
              </tbody>
              <tfoot>
                <tr>
                  <td colspan="3" class="total-label">Total</td>
                  <td class="total-value">
                    LKR {{ formatCurrency(selectedPO.total_amount) }}
                  </td>
                </tr>
              </tfoot>
            </table>
          </div>
        </div>

        <!-- Right: supplier info + PO info + actions -->
        <div class="right">

          <div class="card">
            <h2>Supplier</h2>
            <div class="info-item">
              <span class="info-label">Company</span>
              <span>{{ selectedPO.supplier_name }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Email</span>
              <span>{{ selectedPO.supplier_email || '—' }}</span>
            </div>
          </div>

          <div class="card">
            <h2>Purchase Order Info</h2>
            <div class="info-item">
              <span class="info-label">PO Number</span>
              <span>{{ selectedPO.po_number }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Request</span>
              <span>{{ selectedPO.request_number }} — {{ selectedPO.request_title }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Delivery Date</span>
              <span>{{ formatDate(selectedPO.delivery_date) }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Total Amount</span>
              <span class="amount">LKR {{ formatCurrency(selectedPO.total_amount) }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Created On</span>
              <span>{{ formatDate(selectedPO.created_at) }}</span>
            </div>
          </div>

          <!-- Link to Goods Receipt — only shown when PO is in a receivable state -->
          <div
            class="card action-card"
            v-if="selectedPO.status === 'pending' ||
                  selectedPO.status === 'purchase_order_created'"
          >
            <h2>Next Step</h2>
            <p>Once goods are delivered, record the receipt.</p>
            <button
              class="receipt-btn"
              @click="router.push('/goods-receipt')"
            >
              → Go to Goods Receipt
            </button>
          </div>

          <div class="card action-card" v-if="selectedPO.status === 'goods_received'">
            <p class="done-msg">✓ Goods have been received for this PO.</p>
          </div>

        </div>
      </div>
    </div>

    <!-- ══════════════════════════════════════
         LIST VIEW — shown by default
         ══════════════════════════════════════ -->
    <div v-else>

      <div class="page-header">
        <h1>Purchase Orders</h1>
      </div>

      <!-- ── Requests ready for PO creation ── -->
      <div v-if="pendingRequests.length > 0" class="pending-section">
        <h2>⚠ Requests Awaiting PO Creation</h2>
        <p class="section-sub">
          These requests have a supplier selected — create a purchase order for each.
        </p>

        <div class="pending-cards">
          <div
            v-for="req in pendingRequests"
            :key="req.id"
            class="pending-card"
          >
            <div>
              <span class="req-number">{{ req.request_number }}</span>
              <p class="req-title">{{ req.title }}</p>
              <p class="req-dept">{{ req.department }}</p>
            </div>
            <button
              class="create-po-btn"
              @click="openCreateModal(req)"
              :disabled="acting"
            >
              + Create PO
            </button>
          </div>
        </div>
      </div>

      <!-- ── PO list table ── -->
      <div class="card">
        <div class="card-header">
          <h2>All Purchase Orders</h2>
          <span class="count">{{ pos.length }} total</span>
        </div>

        <div v-if="loading" class="empty">Loading...</div>

        <table v-else>
          <thead>
            <tr>
              <th>PO Number</th>
              <th>Request</th>
              <th>Supplier</th>
              <th>Total (LKR)</th>
              <th>Delivery Date</th>
              <th>Status</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="po in pos"
              :key="po.id"
              class="clickable-row"
            >
              <td class="po-num">{{ po.po_number }}</td>
              <td>
                <div>{{ po.request_number }}</div>
                <div class="sub-text">{{ po.request_title }}</div>
              </td>
              <td>{{ po.supplier_name }}</td>
              <td>{{ formatCurrency(po.total_amount) }}</td>
              <td>{{ formatDate(po.delivery_date) }}</td>
              <td>
                <span :class="['po-status', po.status]">
                  {{ statusLabel(po.status) }}
                </span>
              </td>
              <td>
                <button class="view-btn" @click="viewPO(po.id)">View →</button>
              </td>
            </tr>
            <tr v-if="pos.length === 0">
              <td colspan="7" class="empty">No purchase orders yet</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- ══════════════════════════════════════
         Create PO modal
         ══════════════════════════════════════ -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <div class="modal-header">
          <h3>Create Purchase Order</h3>
          <button class="close-btn" @click="showCreateModal = false">✕</button>
        </div>

        <p class="modal-sub">
          Request: <strong>{{ createForm.requestNumber }}</strong>
        </p>

        <div class="form-group">
          <label>Expected Delivery Date</label>
          <input v-model="createForm.delivery_date" type="date" />
        </div>

        <p v-if="createError" class="error">{{ createError }}</p>

        <div class="modal-actions">
          <button class="cancel-btn" @click="showCreateModal = false">Cancel</button>
          <button
            class="confirm-btn"
            :disabled="acting"
            @click="confirmCreatePO"
          >
            {{ acting ? 'Creating...' : 'Create PO' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

const router = useRouter()

const pos             = ref([])        // all purchase orders
const pendingRequests = ref([])        // requests with status=supplier_selected
const selectedPO      = ref(null)      // the PO being viewed in detail
const loading         = ref(false)
const acting          = ref(false)

const showCreateModal = ref(false)
const createError     = ref('')
const createForm      = ref({
  request_id:    '',
  requestNumber: '',  // just for display in the modal
  delivery_date: '',
})

// Load PO list + requests that still need a PO created
async function loadData() {
  loading.value = true
  const [poRes, reqRes] = await Promise.all([
    api.get('/api/procurement/purchase-orders'),
    api.get('/api/procurement/requests', { params: { status: 'supplier_selected' } }),
  ])
  pos.value             = poRes.data
  pendingRequests.value = reqRes.data
  loading.value = false
}

// Fetch full PO detail (includes items) then switch to detail view
async function viewPO(poId) {
  const { data } = await api.get(`/api/procurement/purchase-orders/${poId}`)
  selectedPO.value = data
}

// Open the create modal pre-filled with the request
function openCreateModal(req) {
  createForm.value  = {
    request_id:    req.id,
    requestNumber: req.request_number,
    delivery_date: '',
  }
  createError.value = ''
  showCreateModal.value = true
}

async function confirmCreatePO() {
  acting.value = true
  createError.value = ''
  try {
    await api.post('/api/procurement/purchase-orders', {
      request_id:    createForm.value.request_id,
      delivery_date: createForm.value.delivery_date,
    })
    showCreateModal.value = false
    await loadData()   // refresh list — new PO appears, request disappears from pending
  } catch (e) {
    createError.value = e.response?.data?.error || 'Failed to create PO.'
  } finally {
    acting.value = false
  }
}

function statusLabel(status) {
  const labels = {
    pending:                 'Pending',
    purchase_order_created:  'Active',
    goods_received:          'Goods Received',
    completed:               'Completed',
  }
  return labels[status] || status
}

function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-GB')
}

function formatCurrency(n) {
  return Number(n || 0).toLocaleString('en-LK')
}

onMounted(loadData)
</script>

<style scoped>
.po-page { padding: 24px; max-width: 1100px; margin: 0 auto; }

.page-header {
  display: flex; align-items: center;
  gap: 16px; margin-bottom: 24px; flex-wrap: wrap;
}
h1 { font-size: 20px; font-weight: 600; }
h2 { font-size: 15px; font-weight: 600; margin-bottom: 12px; }
.subtitle { font-size: 13px; color: #9ca3af; margin-top: 3px; }

.back-btn {
  background: none; border: 1px solid #d1d5db; border-radius: 6px;
  padding: 6px 12px; cursor: pointer; font-size: 13px; flex-shrink: 0;
}
.back-btn:hover { background: #f3f4f6; }

/* PO status pill */
.po-status {
  display: inline-block; padding: 3px 10px;
  border-radius: 12px; font-size: 12px; font-weight: 500;
}
.po-status.pending                { background: #fef9c3; color: #854d0e; }
.po-status.purchase_order_created { background: #dbeafe; color: #1d4ed8; }
.po-status.goods_received         { background: #dcfce7; color: #166534; }
.po-status.completed              { background: #166534; color: #fff; }

/* Pending requests banner */
.pending-section {
  background: #fffbeb; border: 1px solid #fde68a;
  border-radius: 10px; padding: 20px; margin-bottom: 24px;
}
.pending-section h2 { font-size: 14px; margin-bottom: 4px; color: #92400e; }
.section-sub { font-size: 13px; color: #92400e; margin-bottom: 14px; }

.pending-cards { display: flex; flex-direction: column; gap: 10px; }
.pending-card {
  display: flex; justify-content: space-between;
  align-items: center; background: #fff;
  border: 1px solid #fde68a; border-radius: 8px; padding: 14px 16px;
}
.req-number { font-size: 11px; color: #9ca3af; }
.req-title  { font-size: 14px; font-weight: 500; margin: 2px 0; }
.req-dept   { font-size: 12px; color: #9ca3af; }

.create-po-btn {
  padding: 8px 18px; background: #4f46e5; color: #fff;
  border: none; border-radius: 6px; font-size: 13px;
  font-weight: 500; cursor: pointer; white-space: nowrap;
}
.create-po-btn:hover:not(:disabled) { background: #4338ca; }
.create-po-btn:disabled { opacity: 0.5; cursor: not-allowed; }

/* PO list table */
.card {
  background: #fff; border: 1px solid #e5e7eb;
  border-radius: 10px; overflow: hidden;
}
.card-header {
  display: flex; justify-content: space-between;
  align-items: center; padding: 16px 20px;
  border-bottom: 1px solid #f3f4f6;
}
.card-header h2 { margin-bottom: 0; }
.count { font-size: 13px; color: #9ca3af; }

table { width: 100%; border-collapse: collapse; }
th {
  text-align: left; padding: 10px 16px; font-size: 11px;
  color: #6b7280; text-transform: uppercase;
  background: #f9fafb; border-bottom: 1px solid #e5e7eb;
}
td { padding: 13px 16px; font-size: 14px; border-bottom: 1px solid #f3f4f6; }
tr:last-child td { border-bottom: none; }

.clickable-row:hover td { background: #f9fafb; }
.po-num   { font-weight: 600; color: #4f46e5; }
.sub-text { font-size: 12px; color: #9ca3af; margin-top: 2px; }
.empty    { text-align: center; color: #9ca3af; padding: 32px; }

.view-btn {
  background: none; border: none; color: #4f46e5;
  cursor: pointer; font-size: 13px; font-weight: 500;
}
.view-btn:hover { text-decoration: underline; }

/* Detail view */
.detail-grid {
  display: grid; grid-template-columns: 1fr 300px; gap: 20px;
}
@media(max-width: 800px) { .detail-grid { grid-template-columns: 1fr; } }
.left, .right { display: flex; flex-direction: column; gap: 16px; }

.card { padding: 20px; }

tfoot td { padding: 12px 16px; border-top: 2px solid #e5e7eb; }
.total-label { font-weight: 600; text-align: right; }
.total-value { font-weight: 700; font-size: 15px; }

.info-item { display: flex; flex-direction: column; gap: 3px; margin-bottom: 12px; }
.info-item:last-child { margin-bottom: 0; }
.info-label { font-size: 11px; color: #9ca3af; text-transform: uppercase; }
.amount { font-weight: 700; font-size: 16px; }

.action-card { background: #f9fafb; }
.action-card p { font-size: 13px; color: #6b7280; margin-bottom: 12px; }
.receipt-btn {
  padding: 9px 18px; background: #4f46e5; color: #fff;
  border: none; border-radius: 6px; font-size: 14px;
  font-weight: 500; cursor: pointer; width: 100%;
}
.receipt-btn:hover { background: #4338ca; }
.done-msg { color: #166534; font-size: 14px; font-weight: 500; margin: 0 !important; }

/* Modal */
.modal-overlay {
  position: fixed; inset: 0; background: rgba(0,0,0,0.4);
  display: flex; align-items: center; justify-content: center; z-index: 100;
}
.modal {
  background: #fff; border-radius: 12px; padding: 28px;
  width: 420px; max-width: 92vw;
}
.modal-header {
  display: flex; justify-content: space-between;
  align-items: center; margin-bottom: 6px;
}
.modal-header h3 { font-size: 16px; font-weight: 600; }
.close-btn {
  background: none; border: none; font-size: 16px;
  cursor: pointer; color: #9ca3af;
}
.modal-sub { font-size: 13px; color: #6b7280; margin-bottom: 16px; }

.form-group { display: flex; flex-direction: column; gap: 5px; margin-bottom: 16px; }
.form-group label { font-size: 13px; font-weight: 500; }
.form-group input {
  padding: 9px 12px; border: 1px solid #d1d5db;
  border-radius: 6px; font-size: 14px; outline: none;
}
.form-group input:focus { border-color: #4f46e5; }

.modal-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 8px; }
.cancel-btn {
  padding: 9px 18px; border: 1px solid #d1d5db;
  border-radius: 6px; background: #fff; cursor: pointer; font-size: 14px;
}
.confirm-btn {
  padding: 9px 22px; background: #4f46e5; color: #fff;
  border: none; border-radius: 6px; font-size: 14px;
  font-weight: 500; cursor: pointer;
}
.confirm-btn:hover:not(:disabled) { background: #4338ca; }
.confirm-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.error { color: #dc2626; font-size: 13px; margin-bottom: 8px; }
</style>