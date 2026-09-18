import type { RouteRecordRaw } from 'vue-router'
import { SafetyOutlined } from '@antdv-next/icons'

export const firewallRoutes: RouteRecordRaw[] = [
  {
    path: 'firewall',
    name: 'Firewall',
    component: () => import('@/views/firewall/FirewallView.vue'),
    meta: {
      name: () => $gettext('Firewall'),
      icon: SafetyOutlined,
    },
  },
]
