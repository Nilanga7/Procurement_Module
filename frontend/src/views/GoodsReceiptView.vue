<template>
  <div class="page">
    <div class="page-header">
      <h1>Goods Receipts</h1>
      <button class="primary-btn" @click="showForm = true">+ Record Receipt</button>
    </div>

    <!-- ── Success banner — shown after a successful submit ── -->
    <div v-if="successMsg" class="success-banner">
      <div>
        <strong>✓ Goods receipt recorded successfully.</strong>
        <p>{{ successMsg }}</p>
      </div>
      <button @click="successMsg = ''">✕</button>
    </div>

    <!-- ── Existing receipts list ── -->
    <div class="card">
      <div v-if="loading" class="empty">Loading...</div>
      <table v-else>
        <thead>
          <tr>
            <th>PO Number</th>
            <th>Supplier</th>
            <th>Received Date</th>
            <th>Items</th>
            <th>Remarks</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in receipts" :key="r.id">
            <td class="po-num">{{ r.po_number }}</td>
            <td>{{ r.supplier_name }}</td>
            <td>{{ formatDate(r.received_date) }}</td>
            <td>
              <span v-for="item in r.items" :key="item.item_name" class="item-chip">
                {{ item.item_name }} ×{{ item.quantity_received }}
              </span>
            </td>
            <td>{{ r.remarks || '—' }}</td>
          </tr>
          <tr v-if="receipts.length === 0">
            <td colspan="5" class="empty">No receipts recorded yet</td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ── Record Receipt modal ── -->
    <div v-if="showForm" class="modal-overlay" @click.self="closeForm">
      <div class="modal">
        <div class="modal-header">
          <h3>Record Goods Receipt</h3>
          <button class="close-btn" @click="closeForm">✕</button>
        </div>

        <!-- Step 1: Select PO -->
        <div class="form-group">
          <label>Purchase Order *</label>
          <select v-model="form.po_id" @change="onPOSelected" required>
            <option value="">— Select a PO —</option>
            <option v-for="po in availablePOs" :key="po.id" :value="po.id">
              {{ po.po_number }} — {{ po.supplier_name }}
            </option>
          </select>
        </div>

        <!-- Step 2: Date + remarks -->
        <div class="form-row">
          <div class="form-group">
            <label>Received Date *</label>
            <input v-model="form.received_date" type="date" required />
          </div>
          <div class="form-group">
            <label>Remarks</label>
            <input v-model="form.remarks" type="text" placeholder="Optional notes" />
          </div>
        </div>

        <!-- Step 3: Items (auto-populated from PO, qty editable) -->
        <div v-if="form.items.length" class="items-section">
          <label>Items Received</label>
          <div v-for="(item, i) in form.items" :key="i" class="item-row">
            <span class="item-name">{{ item.item_name }}</span>
            <div class="qty-group">
              <label>Qty received</label>
              <input v-model.number="item.quantity_received" type="number" min="0" />
            </div>
          </div>
        </div>

        <p v-if="formError" class="error">{{ formError }}</p>

        <div class="modal-actions">
          <button class="cancel-btn" @click="closeForm">Cancel</button>
          <button
            class="primary-btn"
            :disabled="submitting || !form.po_id || !form.received_date"
            @click="submitReceipt"
          >
            {{ submitting ? 'Saving...' : 'Save Receipt' }}
          </button>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '../api'

const receipts    = ref([])
const availablePOs = ref([])   // POs with status=purchase_order_created
const loading     = ref(false)
const showForm    = ref(false)
const submitting  = ref(false)
const formError   = ref('')
const successMsg  = ref('')

const form = ref({ po_id: '', received_date: '', remarks: '', items: [] })

async function loadData() {
  loading.value = true
  const [recRes, poRes] = await Promise.all([
    api.get('/api/procurement/goods-receipts'),
    api.get('/api/procurement/purchase-orders'),
  ])
  receipts.value     = recRes.data
  // Only show POs that are ready to receive
  availablePOs.value = poRes.data.filter(
    p => p.status === 'purchase_order_created' || p.status === 'pending'
  )
  loading.value = false
}

// When a PO is selected, fetch its items to pre-fill the items list
async function onPOSelected() {
  form.value.items = []
  if (!form.value.po_id) return
  const { data } = await api.get(`/api/procurement/purchase-orders/${form.value.po_id}`)
  // Pre-fill each item with its full ordered quantity as default
  form.value.items = data.items.map(i => ({
    item_name:         i.item_name,
    quantity_received: i.quantity,   // default to full qty, officer can edit down
  }))
}

async function submitReceipt() {
  submitting.value = true
  formError.value  = ''
  try {
    await api.post('/api/procurement/goods-receipts', form.value)

    // Check if any items are assets — tell the officer if so
    const selectedPO = availablePOs.value.find(p => p.id === form.value.po_id)
    successMsg.value =
      'Finance module has been notified for payment. ' +
      'If this PO contained asset-type items, details have been sent to the Asset Management Module.'

    closeForm()
    await loadData()
  } catch (e) {
    formError.value = e.response?.data?.error || 'Failed to save receipt.'
  } finally {
    submitting.value = false
  }
}

function closeForm() {
  showForm.value  = false
  formError.value = ''
  form.value = { po_id: '', received_date: '', remarks: '', items: [] }
}

function formatDate(d) {
  if (!d) return '—'
  return new Date(d).toLocaleDateString('en-GB')
}

onMounted(loadData)
</script>

<style scoped>
.page { padding: 24px; max-width: 1000px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
h1 { font-size: 20px; font-weight: 600; }
h3 { font-size: 16px; font-weight: 600; }

.primary-btn { padding: 9px 20px; background: #4f46e5; color: #fff;
  border: none; border-radius: 6px; font-size: 14px; font-weight: 500; cursor: pointer; }
.primary-btn:hover:not(:disabled) { background: #4338ca; }
.primary-btn:disabled { opacity: 0.5; cursor: not-allowed; }

.success-banner { display: flex; justify-content: space-between; align-items: flex-start;
  background: #f0fdf4; border: 1px solid #bbf7d0; border-radius: 8px;
  padding: 14px 16px; margin-bottom: 20px; font-size: 13px; color: #166534; }
.success-banner p { margin-top: 4px; font-weight: 400; }
.success-banner button { background: none; border: none; cursor: pointer;
  color: #166534; font-size: 16px; }

.card { background: #fff; border: 1px solid #e5e7eb; border-radius: 10px; overflow: hidden; }
table { width: 100%; border-collapse: collapse; }
th { text-align: left; padding: 10px 16px; font-size: 11px; color: #6b7280;
  text-transform: uppercase; background: #f9fafb; border-bottom: 1px solid #e5e7eb; }
td { padding: 12px 16px; font-size: 13px; border-bottom: 1px solid #f3f4f6; }
tr:last-child td { border-bottom: none; }
.po-num { font-weight: 600; color: #4f46e5; }
.empty  { text-align: center; color: #9ca3af; padding: 32px; }
.item-chip { display: inline-block; background: #f3f4f6; border-radius: 4px;
  padding: 2px 7px; font-size: 11px; margin: 2px 2px 2px 0; }

/* Modal */
.modal-overlay { position: fixed; inset: 0; background: rgba(0,0,0,0.4);
  display: flex; align-items: center; justify-content: center; z-index: 100; }
.modal { background: #fff; border-radius: 12px; padding: 28px;
  width: 500px; max-width: 94vw; max-height: 90vh; overflow-y: auto; }
.modal-header { display: flex; justify-content: space-between;
  align-items: center; margin-bottom: 20px; }
.close-btn { background: none; border: none; font-size: 16px;
  cursor: pointer; color: #9ca3af; }

.form-group { display: flex; flex-direction: column; gap: 5px; margin-bottom: 14px; }
.form-group label { font-size: 13px; font-weight: 500; color: #374151; }
.form-group input,
.form-group select { padding: 9px 12px; border: 1px solid #d1d5db;
  border-radius: 6px; font-size: 14px; outline: none; }
.form-group input:focus,
.form-group select:focus { border-color: #4f46e5; }
.form-row { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }

.items-section { margin-bottom: 16px; }
.items-section > label { font-size: 13px; font-weight: 500;
  color: #374151; display: block; margin-bottom: 8px; }
.item-row { display: flex; justify-content: space-between; align-items: center;
  padding: 10px 12px; background: #f9fafb; border-radius: 6px; margin-bottom: 6px; }
.item-name { font-size: 13px; font-weight: 500; }
.qty-group { display: flex; align-items: center; gap: 8px; }
.qty-group label { font-size: 12px; color: #6b7280; }
.qty-group input { width: 70px; padding: 5px 8px; border: 1px solid #d1d5db;
  border-radius: 6px; font-size: 13px; text-align: center; }

.modal-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 20px; }
.cancel-btn { padding: 9px 18px; border: 1px solid #d1d5db;
  border-radius: 6px; background: #fff; cursor: pointer; font-size: 14px; }
.error { color: #dc2626; font-size: 13px; margin-bottom: 8px; }
</style>