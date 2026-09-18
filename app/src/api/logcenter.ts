import { http } from '@uozi-admin/request'

export interface LogSource {
  id: string
  name: string
  type: string
  description: string
}

export interface LogQueryResult {
  source: string
  keyword: string
  count: number
  lines: string[]
}

export const logcenterApi = {
  getSources: () => http.get<LogSource[]>('/logcenter/sources'),
  queryLogs: (source: string, keyword = '', limit = 200) =>
    http.get<LogQueryResult>(`/logcenter/query?source=${source}&keyword=${encodeURIComponent(keyword)}&limit=${limit}`),
}

export default logcenterApi
