<template>
  <div class="p-10 max-w-4xl mx-auto">
    <!-- Page header -->
    <div class="flex items-center gap-3 mb-8">
      <button @click="router.push('/requests')"
        class="flex items-center gap-1 px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
               text-[#474651] hover:bg-[#eff4ff] transition-colors">
        <span class="material-symbols-outlined text-[16px]">arrow_back</span>
        Back
      </button>
      <div>
        <h1 class="text-[24px] font-bold text-[#0b1c30]">New Procurement Request</h1>
        <p class="text-[13px] text-[#777682]">Fill in the details below to submit a new request</p>
      </div>
    </div>

    <form @submit.prevent="submitForm" class="space-y-6">

      <!-- Section 1: Basic Details -->
      <div class="bg-white border border-[#c8c5d3] rounded-xl p-6 ambient-shadow">
        <h2 class="text-[16px] font-bold text-[#0b1c30] mb-5 flex items-center gap-2">
          <span class="material-symbols-outlined text-[20px] text-[#1a146b]">description</span>
          Request Details
        </h2>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Title *</label>
            <input v-model="form.title" type="text" required placeholder="e.g. Laptop Purchase Q2"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
          </div>
          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Department</label>
            <input v-model="form.department" type="text" placeholder="e.g. IT, Finance, HR"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
          </div>
        </div>

        <div class="mb-4">
          <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Description</label>
          <textarea v-model="form.description" rows="3"
            placeholder="Describe the purpose of this request..."
            class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                   focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10 resize-none"></textarea>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Required Date</label>
            <input v-model="form.required_date" type="date"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
          </div>
          <div>
            <label class="block text-[12px] font-semibold text-[#474651] mb-1.5">Estimated Total (LKR)</label>
            <input v-model.number="form.estimated_total" type="number" min="0" step="0.01"
              class="w-full px-4 py-2.5 border border-[#c8c5d3] rounded-lg text-[14px]
                     focus:outline-none focus:border-[#1a146b] focus:ring-2 focus:ring-[#1a146b]/10" />
          </div>
        </div>
      </div>

      <!-- Section 2: Items -->
      <div class="bg-white border border-[#c8c5d3] rounded-xl p-6 ambient-shadow">
        <div class="flex justify-between items-center mb-5">
          <h2 class="text-[16px] font-bold text-[#0b1c30] flex items-center gap-2">
            <span class="material-symbols-outlined text-[20px] text-[#1a146b]">inventory_2</span>
            Items
          </h2>
          <button type="button" @click="addItem"
            class="flex items-center gap-1.5 px-4 py-2 bg-[#eff4ff] border border-[#c8c5d3]
                   rounded-lg text-[12px] font-bold text-[#1a146b] hover:bg-[#dce9ff] transition-colors">
            <span class="material-symbols-outlined text-[16px]">add</span>
            Add Item
          </button>
        </div>

        <div v-if="form.items.length === 0"
          class="py-10 text-center border-2 border-dashed border-[#c8c5d3] rounded-xl">
          <span class="material-symbols-outlined text-[36px] text-[#c8c5d3] block mb-2">add_shopping_cart</span>
          <p class="text-[14px] text-[#777682]">No items yet. Click "Add Item" to start.</p>
        </div>

        <div v-for="(item, index) in form.items" :key="index"
          class="mb-4 p-4 bg-[#f8f9ff] border border-[#c8c5d3] rounded-xl">
          <div class="flex items-center justify-between mb-3">
            <span class="text-[11px] font-bold text-[#1a146b] bg-[#e2dfff] px-2 py-0.5 rounded-full">
              Item {{ index + 1 }}
            </span>
            <button type="button" @click="removeItem(index)"
              class="text-[#ba1a1a] hover:bg-[#ffdad6] p-1 rounded-lg transition-colors">
              <span class="material-symbols-outlined text-[18px]">delete</span>
            </button>
          </div>
          <div class="grid grid-cols-2 md:grid-cols-5 gap-3">
            <div class="col-span-2">
              <label class="block text-[11px] font-semibold text-[#474651] mb-1">Item Name *</label>
              <input v-model="item.item_name" type="text" required placeholder="e.g. Laptop"
                class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                       focus:outline-none focus:border-[#1a146b]" />
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-[#474651] mb-1">Qty *</label>
              <input v-model.number="item.quantity" type="number" min="1" required @input="recalcTotal"
                class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                       focus:outline-none focus:border-[#1a146b]" />
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-[#474651] mb-1">Unit Price *</label>
              <input v-model.number="item.estimated_price" type="number" min="0" step="0.01" required @input="recalcTotal"
                class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                       focus:outline-none focus:border-[#1a146b]" />
            </div>
            <div>
              <label class="block text-[11px] font-semibold text-[#474651] mb-1">Type *</label>
              <select v-model="item.item_type" required
                class="w-full px-3 py-2 border border-[#c8c5d3] rounded-lg text-[13px]
                       focus:outline-none focus:border-[#1a146b] bg-white">
                <option value="">Select</option>
                <option value="asset">Asset</option>
                <option value="consumable">Consumable</option>
              </select>
            </div>
          </div>
        </div>

        <div v-if="form.items.length > 0"
          class="mt-4 flex justify-end">
          <div class="bg-[#e5eeff] border border-[#c8c5d3] rounded-xl px-6 py-3">
            <span class="text-[12px] text-[#474651] font-semibold">Calculated Total: </span>
            <span class="text-[18px] font-bold text-[#1a146b]">LKR {{ formatCurrency(form.estimated_total) }}</span>
          </div>
        </div>
      </div>

      <!-- Error -->
      <div v-if="error"
        class="flex items-center gap-2 px-4 py-3 bg-[#ffdad6] border border-[#ba1a1a]/30
               rounded-xl text-[13px] text-[#93000a]">
        <span class="material-symbols-outlined text-[16px]">error</span>
        {{ error }}
      </div>

      <!-- Actions -->
      <div class="flex justify-end gap-3">
        <button type="button" @click="router.push('/requests')"
          class="px-6 py-2.5 border border-[#c8c5d3] rounded-xl text-[13px] font-semibold
                 text-[#474651] hover:bg-[#eff4ff] transition-colors">
          Cancel
        </button>
        <button type="submit" :disabled="submitting"
          class="px-8 py-2.5 bg-[#1a146b] text-white rounded-xl text-[13px] font-bold
                 hover:opacity-90 disabled:opacity-50 transition-all flex items-center gap-2">
          <svg v-if="submitting" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"/>
          </svg>
          {{ submitting ? 'Submitting…' : 'Create Request' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../api'

const router    = useRouter()
const submitting = ref(false)
const error      = ref('')

const form = ref({
  title: '', description: '', department: '',
  required_date: '', estimated_total: 0, items: [],
})

function addItem() {
  form.value.items.push({ item_name: '', quantity: 1, estimated_price: 0, category: '', item_type: '' })
}

function removeItem(index) {
  form.value.items.splice(index, 1)
  recalcTotal()
}

function recalcTotal() {
  form.value.estimated_total = form.value.items.reduce(
    (sum, item) => sum + (item.quantity * item.estimated_price || 0), 0
  )
}

function formatCurrency(n) {
  return Number(n || 0).toLocaleString('en-LK')
}

async function submitForm() {
  if (form.value.items.length === 0) { error.value = 'Please add at least one item.'; return }
  submitting.value = true
  error.value = ''
  try {
    const { data } = await api.post('/api/procurement/requests', form.value)
    if (!data.id) { error.value = 'Server did not return a request ID. Please try again.'; return }
    router.push(`/requests/${data.id}`)
  } catch (e) {
    error.value = e.response?.data?.error || 'Failed to create request. Please try again.'
  } finally {
    submitting.value = false
  }
}
</script>
