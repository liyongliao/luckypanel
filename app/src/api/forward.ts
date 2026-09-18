import { http } from '@uozi-admin/request'

export interface PortForwardRule {
  id: number
  name: string
  protocol: string
  listen_ip: string
  listen_port: number
  target_ip: string
  target_port: number
  enabled: boolean
  enable_upnp: boolean
  auto_open_firewall: boolean
  rx_bytes: number
  tx_bytes: number
  description: string
  active_conns?: number
  is_running?: boolean
  created_at: string
}

export const forwardApi = {
  getRules: () => http.get<PortForwardRule[]>('/forward/rules'),
  createRule: (data: Partial<PortForwardRule>) => http.post<PortForwardRule>('/forward/rules', data),
  updateRule: (id: number, data: Partial<PortForwardRule>) => http.put<PortForwardRule>(`/forward/rules/${id}`, data),
  deleteRule: (id: number) => http.delete(`/forward/rules/${id}`),
  toggleRule: (id: number) => http.post<PortForwardRule>(`/forward/rules/${id}/toggle`),
}

export default forwardApi
