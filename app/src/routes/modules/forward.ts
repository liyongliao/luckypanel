import type { RouteRecordRaw } from 'vue-router'
import { SwapOutlined } from '@antdv-next/icons'

export const forwardRoutes: RouteRecordRaw[] = [
  {
    path: 'forward',
    name: 'PortForward',
    component: () => import('@/views/forward/PortForwardView.vue'),
    meta: {
      name: () => $gettext('Port Forward'),
      icon: SwapOutlined,
    },
  },
]
