import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  { path: '/login',    component: () => import('../views/LoginView.vue'),    meta: { public: true } },
  { path: '/register', component: () => import('../views/RegisterView.vue'), meta: { public: true } },
  // ── Supplier portal — completely separate layout, no internal navbar ──
  { path: '/supplier', component: () => import('../views/supplier/SupplierPortalView.vue'), meta: { role: 'supplier', hideNav: true } },
  // ── Internal staff pages ──
  { path: '/',              component: () => import('../views/DashboardView.vue') },
  { path: '/requests',     component: () => import('../views/RequestListView.vue') },
  { path: '/requests/new', component: () => import('../views/CreateRequestView.vue') },
  { path: '/requests/:id', component: () => import('../views/RequestDetailView.vue') },
  { path: '/bids/:requestId',    component: () => import('../views/BidEvaluationView.vue') },
  { path: '/purchase-orders',   component: () => import('../views/PurchaseOrderView.vue') },
  { path: '/goods-receipt',     component: () => import('../views/GoodsReceiptView.vue') },
  { path: '/transactions',      component: () => import('../views/TransactionFilterView.vue') },
]

const router = createRouter({ history: createWebHistory(), routes })

router.beforeEach((to, _, next) => {
  const auth = useAuthStore()

  // 1. Unauthenticated users can only access public routes
  if (!to.meta.public && !auth.isLoggedIn) return next('/login')

  // 2. Suppliers must stay in their own portal — block internal pages
  if (auth.isLoggedIn && auth.user?.role === 'supplier' && to.path !== '/supplier') {
    return next('/supplier')
  }

  // 3. Non-suppliers cannot access the supplier portal
  if (to.meta.role && auth.user?.role !== to.meta.role) return next('/')

  next()
})

export default router