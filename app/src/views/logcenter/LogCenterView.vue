<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
const { message } = App.useApp()
import logcenterApi, { type LogSource } from '@/api/logcenter'

const sources = ref<LogSource[]>([])
const selectedSource = ref('nginx_access')
const keyword = ref('')
const limit = ref(200)
const loading = ref(false)
const autoRefresh = ref(false)
const logLines = ref<string[]>([])
let timer: any = null

async function fetchSources() {
  try {
    const res = await logcenterApi.getSources()
    sources.value = res.data || res || []
  } catch (e: any) {
    console.error(e)
  }
}

async function fetchLogs() {
  loading.value = true
  try {
    const res = await logcenterApi.queryLogs(selectedSource.value, keyword.value, limit.value)
    const data = res.data || res
    logLines.value = data.lines || []
  } catch (e: any) {
    message.error('拉取日志失败: ' + (e.message || ''))
  } finally {
    loading.value = false
  }
}

function handleAutoRefreshChange(val: boolean) {
  if (val) {
    timer = setInterval(fetchLogs, 3000)
  } else if (timer) {
    clearInterval(timer)
    timer = null
  }
}

onMounted(async () => {
  await fetchSources()
  await fetchLogs()
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="p-4 space-y-4">
    <ACard title="统一日志中心 (Log Center)" :bordered="false" class="shadow-sm">
      <AFlex justify="space-between" align="center" wrap="wrap" gap="middle" class="mb-4">
        <AFlex align="center" gap="middle" wrap="wrap">
          <div>
            <span class="text-xs text-gray-500 mr-2">选择日志源:</span>
            <ASelect
              v-model:value="selectedSource"
              style="width: 280px"
              @change="() => fetchLogs()"
            >
              <ASelectOption v-for="s in sources" :key="s.id" :value="s.id">
                {{ s.name }}
              </ASelectOption>
            </ASelect>
          </div>

          <div>
            <span class="text-xs text-gray-500 mr-2">关键词过滤:</span>
            <AInput
              v-model:value="keyword"
              placeholder="搜索包含文本..."
              style="width: 200px"
              allow-clear
              @pressEnter="fetchLogs"
            />
          </div>

          <div>
            <span class="text-xs text-gray-500 mr-2">行数:</span>
            <ASelect v-model:value="limit" style="width: 90px" @change="fetchLogs">
              <ASelectOption :value="100">100</ASelectOption>
              <ASelectOption :value="200">200</ASelectOption>
              <ASelectOption :value="500">500</ASelectOption>
              <ASelectOption :value="1000">1000</ASelectOption>
            </ASelect>
          </div>
        </AFlex>

        <AFlex align="center" gap="middle">
          <AFlex align="center" gap="small">
            <span class="text-xs text-gray-500">自动滚动追踪:</span>
            <ASwitch
              v-model:checked="autoRefresh"
              size="small"
              @change="(val: any) => handleAutoRefreshChange(val)"
            />
          </AFlex>

          <AButton type="primary" :loading="loading" @click="fetchLogs">
            检索 / 刷新
          </AButton>
        </AFlex>
      </AFlex>

      <!-- Terminal Log Screen -->
      <div class="bg-gray-950 text-gray-200 p-4 rounded font-mono text-xs h-[650px] overflow-y-auto border border-gray-800 space-y-1 select-text">
        <div v-if="logLines.length === 0" class="text-gray-500 text-center py-20">
          暂无匹配的日志记录
        </div>
        <div
          v-for="(line, idx) in logLines"
          :key="idx"
          class="flex items-start hover:bg-gray-900 px-1 py-0.5 rounded leading-relaxed"
        >
          <span class="text-gray-600 select-none w-12 text-right pr-3 shrink-0">{{ idx + 1 }}</span>
          <span class="text-gray-300 break-all whitespace-pre-wrap flex-1">{{ line }}</span>
        </div>
      </div>
    </ACard>
  </div>
</template>
