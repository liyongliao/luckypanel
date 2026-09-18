import { http } from '@uozi-admin/request'

export interface DockerStatus {
  available: boolean
  server_version: string
  containers: number
  containers_running: number
  containers_paused: number
  containers_stopped: number
  images: number
  message: string
}

export interface DockerContainer {
  id: string
  names: string[]
  image: string
  state: string
  status: string
  ports: string[]
  created: number
}

export interface DockerImage {
  id: string
  repo_tags: string[]
  size: number
  created: number
}

export interface ComposeStack {
  name: string
  path: string
  content: string
  updated_at: string
  status: string
}

export const dockerApi = {
  getStatus: () => http.get<DockerStatus>('/docker/status'),
  getContainers: () => http.get<{ data: DockerContainer[] }>('/docker/containers'),
  containerAction: (id: string, action: string) => http.post(`/docker/containers/${id}/action`, { action }),
  getLogs: (id: string, tail?: string) => http.get<{ logs: string }>(`/docker/containers/${id}/logs?tail=${tail || '200'}`),
  getImages: () => http.get<{ data: DockerImage[] }>('/docker/images'),
  getCompose: () => http.get<ComposeStack[]>('/docker/compose'),
  saveCompose: (name: string, content: string) => http.post('/docker/compose', { name, content }),
  composeAction: (name: string, action: string) => http.post<{ output: string }>(`/docker/compose/${name}/action`, { action }),
}

export default dockerApi
