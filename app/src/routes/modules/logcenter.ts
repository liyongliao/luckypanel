import type { RouteRecordRaw } from 'vue-router'
import { FileTextOutlined } from '@antdv-next/icons'

export const logcenterRoutes: RouteRecordRaw[] = [
  {
    path: 'logcenter',
    name: 'LogCenter',
    component: () => import('@/views/logcenter/LogCenterView.vue'),
    meta: {
      name: () => $gettext('Log Center'),
      icon: FileTextOutlined,
    },
  },
]
