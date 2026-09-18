import { http } from '@uozi-admin/request'

export interface FirewallStatus {
  type: string
  is_active: boolean
  available: boolean
  message: string
}

export interface FirewallRule {
  id: number
  port: string
  protocol: string
  strategy: string
  source: string
  description: string
  created_at: string
  updated_at: string
}

export const firewallApi = {
  getStatus: () => http.get<FirewallStatus>('/firewall/status'),
  toggleStatus: (enable: boolean) => http.post('/firewall/toggle', { enable }),
  getRules: () => http.get<FirewallRule[]>('/firewall/rules'),
  openPort: (data: { port: string, protocol: string, source?: string, description?: string }) =>
    http.post<FirewallRule>('/firewall/rules', data),
  closePort: (id: number) => http.delete(`/firewall/rules/${id}`),
}

export default firewallApi
