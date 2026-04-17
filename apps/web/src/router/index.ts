import { createRouter as newRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '../stores/auth'

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
        component: () => import('../views/users/UserListView.vue')
      },
      {
        path: 'departments',
        name: 'departments',
        component: () => import('../views/departments/DepartmentListView.vue')
      },
      {
        path: 'roles',
        name: 'roles',
        component: () => import('../views/roles/RoleListView.vue')
      },
      {
        path: 'menus',
        name: 'menus',
        component: () => import('../views/menus/MenuListView.vue')
      },
      {
        path: 'dictionary',
        name: 'dictionary',
        component: () => import('../views/dictionary/DictionaryListView.vue')
      },
      {
        path: 'system-params',
        name: 'system-params',
        component: () => import('../views/params/SystemParamListView.vue')
      },
      {
        path: 'logs',
        name: 'logs',
        component: () => import('../views/logs/LogCenterView.vue')
      },
      {
        path: 'monitor',
        name: 'monitor',
        component: () => import('../views/monitor/MonitorView.vue')
      },
      {
        path: 'positions',
        name: 'positions',
        component: () => import('../views/positions/PositionListView.vue')
      },
      {
        path: 'notices',
        name: 'notices',
        component: () => import('../views/notices/NoticeListView.vue')
      },
      {
        path: 'online-users',
        name: 'online-users',
        component: () => import('../views/online-users/OnlineUserView.vue')
      },
      {
        path: 'scheduled-jobs',
        name: 'scheduled-jobs',
        component: () => import('../views/jobs/ScheduledJobView.vue')
      },
      {
        path: 'generated/department',
        name: 'generated-department',
        component: () => import('../generated/department/ListView.generated.vue')
      }
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
    if (to.path === '/login' && auth.isLoggedIn) {
      return '/'
    }
    return true
  })

  return router
}
