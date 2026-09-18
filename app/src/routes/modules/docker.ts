import type { RouteRecordRaw } from 'vue-router'
import { AppstoreOutlined } from '@antdv-next/icons'

export const dockerRoutes: RouteRecordRaw[] = [
  {
    path: 'docker',
    name: 'Docker',
    component: () => import('@/views/docker/DockerView.vue'),
    meta: {
      name: () => $gettext('Docker'),
      icon: AppstoreOutlined,
    },
  },
]
