<template>
  <span :class="badge.cls" class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-[11px] font-semibold whitespace-nowrap">
    <span class="w-1.5 h-1.5 rounded-full" :class="badge.dot"></span>
    {{ badge.label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ status: String })

const MAP = {
  draft:                    { label: 'Draft',              cls: 'bg-[#d3e4fe] text-[#1a146b]',         dot: 'bg-[#1a146b]' },
  sent_to_finance:          { label: 'Sent to Finance',    cls: 'bg-[#ffddb8] text-[#543300]',         dot: 'bg-[#543300]' },
  budget_approved:          { label: 'Budget Approved',    cls: 'bg-[#6cf8bb]/40 text-[#006c49]',      dot: 'bg-[#006c49]' },
  budget_rejected:          { label: 'Budget Rejected',    cls: 'bg-[#ffdad6] text-[#93000a]',         dot: 'bg-[#93000a]' },
  manager_approved:         { label: 'Mgr Approved',       cls: 'bg-[#c3c0ff]/40 text-[#1a146b]',      dot: 'bg-[#1a146b]' },
  manager_rejected:         { label: 'Mgr Rejected',       cls: 'bg-[#ffdad6] text-[#93000a]',         dot: 'bg-[#93000a]' },
  open_for_bidding:         { label: 'Open for Bidding',   cls: 'bg-[#3e3c8f] text-[#e2dfff]',         dot: 'bg-[#9c9af4]' },
  supplier_selected:        { label: 'Supplier Selected',  cls: 'bg-[#6cf8bb]/40 text-[#005236]',      dot: 'bg-[#005236]' },
  purchase_order_created:   { label: 'PO Created',         cls: 'bg-[#e5eeff] text-[#1a146b]',         dot: 'bg-[#5654a8]' },
  goods_received:           { label: 'Goods Received',     cls: 'bg-[#6cf8bb]/30 text-[#006c49]',      dot: 'bg-[#006c49]' },
  sent_to_asset_module:     { label: 'Sent to Assets',     cls: 'bg-[#ffddb8] text-[#372000]',         dot: 'bg-[#372000]' },
  sent_to_finance_for_payment: { label: 'Sent for Payment', cls: 'bg-[#ffddb8] text-[#543300]',       dot: 'bg-[#543300]' },
  completed:                { label: 'Completed',          cls: 'bg-[#6cf8bb] text-[#00714d]',         dot: 'bg-[#006c49]' },
  pending:                  { label: 'Pending',            cls: 'bg-[#ffddb8] text-[#653e00]',         dot: 'bg-[#653e00]' },
  accepted:                 { label: 'Accepted ✓',         cls: 'bg-[#6cf8bb] text-[#00714d]',         dot: 'bg-[#006c49]' },
  rejected:                 { label: 'Not Selected',       cls: 'bg-[#d3e4fe] text-[#474651]',         dot: 'bg-[#777682]' },
}

const badge = computed(() =>
  MAP[props.status] || { label: props.status, cls: 'bg-[#d3e4fe] text-[#474651]', dot: 'bg-[#777682]' }
)
</script>