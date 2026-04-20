import { createRouter as newRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { generatedModuleRoutes } from '../generated/module-routes.generated'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/LoginView.vue')
  },
  {
    path: '/',
    component: () => import('../layouts/AdminLayout.vue'),
    children: [
      {
        path: '',
        name: 'dashboard',
        component: () => import('../views/DashboardView.vue')
      },
      {
        path: 'users',
        name: 'users',
        component: () => import('../views/users/UserListView.vue'),
        meta: { perm: 'user:list' }
      },
      {
        path: 'departments',
        name: 'departments',
        component: () => import('../views/departments/DepartmentListView.vue'),
        meta: { perm: 'department:list' }
      },
      {
        path: 'roles',
        name: 'roles',
        component: () => import('../views/roles/RoleListView.vue'),
        meta: { perm: 'role:list' }
      },
      {
        path: 'menus',
        name: 'menus',
        component: () => import('../views/menus/MenuListView.vue'),
        meta: { perm: 'menu:list' }
      },
      {
        path: 'dictionary',
        name: 'dictionary',
        component: () => import('../views/dictionary/DictionaryListView.vue'),
        meta: { perm: 'dictionary:list' }
      },
      {
        path: 'system-params',
        name: 'system-params',
        component: () => import('../views/params/SystemParamListView.vue'),
        meta: { perm: 'param:list' }
      },
      {
        path: 'logs',
        name: 'logs',
        component: () => import('../views/logs/LogCenterView.vue'),
        meta: { perm: 'log:operation:list' }
      },
      {
        path: 'monitor',
        name: 'monitor',
        component: () => import('../views/monitor/MonitorView.vue'),
        meta: { perm: 'monitor:view' }
      },
      {
        path: 'positions',
        name: 'positions',
        component: () => import('../views/positions/PositionListView.vue'),
        meta: { perm: 'position:list' }
      },
      {
        path: 'notices',
        name: 'notices',
        component: () => import('../views/notices/NoticeListView.vue'),
        meta: { perm: 'notice:list' }
      },
      {
        path: 'online-users',
        name: 'online-users',
        component: () => import('../views/online-users/OnlineUserView.vue'),
        meta: { perm: 'online-user:list' }
      },
      {
        path: 'scheduled-jobs',
        name: 'scheduled-jobs',
        component: () => import('../views/jobs/ScheduledJobView.vue'),
        meta: { perm: 'job:list' }
      },
      ...generatedModuleRoutes
    ]
  }
]

export function createRouter() {
  const router = newRouter({
    history: createWebHistory(),
    routes
  })

  router.beforeEach(async (to) => {
    const auth = useAuthStore()
    if (!auth.user && auth.isLoggedIn) {
      try {
        await auth.loadMe()
      } catch {
        auth.logout()
      }
    }
    if (to.path !== '/login' && !auth.isLoggedIn) {
      return '/login'
    }
    const routePerm = typeof to.meta?.perm === 'string' ? to.meta.perm : ''
    if (routePerm && !auth.hasPerm(routePerm)) {
      return '/'
    }
    if (to.path === '/login' && auth.isLoggedIn) {
      return '/'
    }
    return true
  })

  return router
}
