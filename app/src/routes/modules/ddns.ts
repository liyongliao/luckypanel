import type { RouteRecordRaw } from 'vue-router'
import { GlobalOutlined } from '@antdv-next/icons'

export const ddnsRoutes: RouteRecordRaw[] = [
  {
    path: 'ddns',
    name: 'DDNS',
    component: () => import('@/views/ddns/DDNSView.vue'),
    meta: {
      name: () => $gettext('DDNS'),
      icon: GlobalOutlined,
    },
  },
]
