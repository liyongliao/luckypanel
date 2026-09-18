import { http } from '@uozi-admin/request'

export interface DDNSTask {
  id: number
  name: string
  provider: string
  domains: string
  ip_type: string
  ip_method: string
  ip_interface: string
  ip_url: string
  access_key_id?: string
  access_key_secret?: string
  webhook_url?: string
  interval_seconds: number
  enabled: boolean
  last_ipv4?: string
  last_ipv6?: string
  last_run_at?: string
  last_status?: string
  last_error?: string
}

export interface InterfaceInfo {
  name: string
  ipv4s: string[]
  ipv6s: string[]
}

export const ddnsApi = {
  getTasks: () => http.get<DDNSTask[]>('/ddns/tasks'),
  createTask: (data: Partial<DDNSTask>) => http.post<DDNSTask>('/ddns/tasks', data),
  updateTask: (id: number, data: Partial<DDNSTask>) => http.put<DDNSTask>(`/ddns/tasks/${id}`, data),
  deleteTask: (id: number) => http.delete(`/ddns/tasks/${id}`),
  runTaskNow: (id: number) => http.post<DDNSTask>(`/ddns/tasks/${id}/run`),
  getInterfaces: () => http.get<InterfaceInfo[]>('/ddns/interfaces'),
}

export default ddnsApi
