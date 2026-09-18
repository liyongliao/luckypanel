<script setup lang="ts">
import { ref, onMounted } from 'vue'
const { message } = App.useApp()
import firewallApi, { type FirewallRule, type FirewallStatus } from '@/api/firewall'

const loading = ref(false)
const status = ref<FirewallStatus>({
  type: 'none',
  is_active: false,
  available: false,
  message: '',
})
const rules = ref<FirewallRule[]>([])
const showModal = ref(false)

const form = ref({
  port: '',
  protocol: 'tcp',
  source: 'any',
  description: '',
})

async function fetchStatus() {
  try {
    const res = await firewallApi.getStatus()
    status.value = res.data || res
  } catch (e: any) {
    console.error(e)
  }
}

async function fetchRules() {
  loading.value = true
  try {
    const res = await firewallApi.getRules()
    rules.value = res.data || res || []
  } catch (e: any) {
    message.error('获取规则失败: ' + (e.message || '网络错误'))
  } finally {
    loading.value = false
  }
}

async function handleToggleStatus(checked: boolean) {
  try {
    await firewallApi.toggleStatus(checked)
    message.success(checked ? '防火墙已启用' : '防火墙已停用')
    await fetchStatus()
  } catch (e: any) {
    message.error('切换防火墙状态失败: ' + (e.message || ''))
  }
}

async function handleOpenPort() {
  if (!form.value.port) {
    message.warning('请输入开放端口')
    return
  }
  try {
    await firewallApi.openPort(form.value)
    message.success('端口放行成功')
    showModal.value = false
    form.value = { port: '', protocol: 'tcp', source: 'any', description: '' }
    fetchRules()
  } catch (e: any) {
    message.error('操作失败: ' + (e.message || ''))
  }
}

async function handleClosePort(id: number) {
  try {
    await firewallApi.closePort(id)
    message.success('端口规则已移除并关闭')
    fetchRules()
  } catch (e: any) {
    message.error('关闭失败: ' + (e.message || ''))
  }
}

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 80 },
  { title: '放行端口', dataIndex: 'port', key: 'port', width: 140 },
  { title: '协议', dataIndex: 'protocol', key: 'protocol', width: 100 },
  { title: '来源限制 (IP/CIDR)', dataIndex: 'source', key: 'source', width: 160 },
  { title: '策略', dataIndex: 'strategy', key: 'strategy', width: 100 },
  { title: '备注说明', dataIndex: 'description', key: 'description' },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 180 },
  { title: '操作', key: 'action', width: 120 },
]

onMounted(() => {
  fetchStatus()
  fetchRules()
})
</script>

<template>
  <div class="p-4 space-y-4">
    <!-- Top Status Card -->
    <ACard title="系统防火墙状态" :bordered="false" class="shadow-sm">
      <template #extra>
        <AFlex gap="small" align="center">
          <AButton @click="() => { fetchStatus(); fetchRules(); }">刷新</AButton>
          <AButton type="primary" @click="showModal = true">+ 开放端口</AButton>
        </AFlex>
      </template>

      <AFlex align="center" gap="middle" wrap="wrap">
        <div>
          <span class="text-gray-500 mr-2">底层组件:</span>
          <ATag color="blue">{{ status.type ? status.type.toUpperCase() : 'NONE' }}</ATag>
        </div>
        <div>
          <span class="text-gray-500 mr-2">运行状态:</span>
          <ASwitch
            :checked="status.is_active"
            checked-children="已启用"
            un-checked-children="已停用"
            @change="(val: any) => handleToggleStatus(val)"
          />
        </div>
        <div class="text-gray-400 text-sm">
          {{ status.message }}
        </div>
      </AFlex>
    </ACard>

    <!-- Security Info -->
    <AAlert
      type="info"
      show-icon
      message="安全防失联保护已开启"
      description="系统 SSH 端口（22）及当前管理面板通信端口受自动防误删保护，禁止直接移除规则，防止操作失误导致远程失联。"
    />

    <!-- Rules Table -->
    <ACard title="放行端口与访问控制规则" :bordered="false" class="shadow-sm">
      <ATable
        :columns="columns"
        :data-source="rules"
        :loading="loading"
        row-key="id"
        :pagination="{ pageSize: 15 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'port'">
            <span class="font-semibold text-blue-600">{{ record.port }}</span>
          </template>

          <template v-else-if="column.key === 'protocol'">
            <ATag :color="record.protocol === 'tcp' ? 'cyan' : record.protocol === 'udp' ? 'purple' : 'geekblue'">
              {{ record.protocol.toUpperCase() }}
            </ATag>
          </template>

          <template v-else-if="column.key === 'strategy'">
            <ATag color="green">ALLOW</ATag>
          </template>

          <template v-else-if="column.key === 'source'">
            <span :class="record.source === 'any' ? 'text-gray-400' : 'text-emerald-600 font-mono'">
              {{ record.source || 'any (0.0.0.0/0)' }}
            </span>
          </template>

          <template v-else-if="column.key === 'created_at'">
            <span class="text-gray-400 text-xs">{{ record.created_at ? record.created_at.slice(0, 19).replace('T', ' ') : '-' }}</span>
          </template>

          <template v-else-if="column.key === 'action'">
            <APopconfirm
              title="确定关闭此端口并移除防火墙规则吗？"
              ok-text="确认关闭"
              cancel-text="取消"
              @confirm="handleClosePort(record.id)"
            >
              <AButton danger size="small" type="link">关闭端口</AButton>
            </APopconfirm>
          </template>
        </template>
      </ATable>
    </ACard>

    <!-- Modal for Open Port -->
    <AModal
      v-model:open="showModal"
      title="开放防火墙端口"
      ok-text="确认放行"
      cancel-text="取消"
      @ok="handleOpenPort"
    >
      <AForm layout="vertical" class="mt-4">
        <AFormItem label="放行端口 (支持单个如 8080 或范围如 9000-9050)" required>
          <AInput v-model:value="form.port" placeholder="例如: 8848 或 8000-8080" />
        </AFormItem>

        <AFormItem label="传输协议" required>
          <ASelect v-model:value="form.protocol">
            <ASelectOption value="tcp">TCP</ASelectOption>
            <ASelectOption value="udp">UDP</ASelectOption>
            <ASelectOption value="both">TCP/UDP (双协议)</ASelectOption>
          </ASelect>
        </AFormItem>

        <AFormItem label="允许来源 IP (留空或 any 允许所有公网访问)">
          <AInput v-model:value="form.source" placeholder="any 或 192.168.1.0/24 / 指定IP" />
        </AFormItem>

        <AFormItem label="规则备注">
          <AInput v-model:value="form.description" placeholder="例如: Web项目端口、测试服务等" />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>
