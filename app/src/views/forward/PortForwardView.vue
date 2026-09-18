<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
const { message } = App.useApp()
import forwardApi, { type PortForwardRule } from '@/api/forward'

const loading = ref(false)
const rules = ref<PortForwardRule[]>([])
const showModal = ref(false)
const isEditing = ref(false)
const editId = ref(0)
let timer: any = null

const form = ref({
  name: '',
  protocol: 'tcp',
  listen_ip: '0.0.0.0',
  listen_port: 8080,
  target_ip: '127.0.0.1',
  target_port: 80,
  enabled: true,
  enable_upnp: false,
  auto_open_firewall: true,
  description: '',
})

function formatBytes(bytes: number): string {
  if (!bytes || bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

async function fetchRules() {
  try {
    const res = await forwardApi.getRules()
    rules.value = res.data || res || []
  } catch (e: any) {
    console.error(e)
  }
}

function openAddModal() {
  isEditing.value = false
  editId.value = 0
  form.value = {
    name: '',
    protocol: 'tcp',
    listen_ip: '0.0.0.0',
    listen_port: 8080,
    target_ip: '127.0.0.1',
    target_port: 80,
    enabled: true,
    enable_upnp: false,
    auto_open_firewall: true,
    description: '',
  }
  showModal.value = true
}

function openEditModal(record: PortForwardRule) {
  isEditing.value = true
  editId.value = record.id
  form.value = {
    name: record.name,
    protocol: record.protocol,
    listen_ip: record.listen_ip,
    listen_port: record.listen_port,
    target_ip: record.target_ip,
    target_port: record.target_port,
    enabled: record.enabled,
    enable_upnp: record.enable_upnp,
    auto_open_firewall: record.auto_open_firewall,
    description: record.description,
  }
  showModal.value = true
}

async function handleSubmit() {
  if (!form.value.name || !form.value.listen_port || !form.value.target_ip || !form.value.target_port) {
    message.warning('请补全端口转发参数')
    return
  }
  try {
    if (isEditing.value) {
      await forwardApi.updateRule(editId.value, form.value)
      message.success('规则更新成功')
    } else {
      await forwardApi.createRule(form.value)
      message.success('规则创建成功')
    }
    showModal.value = false
    fetchRules()
  } catch (e: any) {
    message.error('保存失败: ' + (e.message || ''))
  }
}

async function handleToggle(record: PortForwardRule) {
  try {
    await forwardApi.toggleRule(record.id)
    message.success('状态已切换')
    fetchRules()
  } catch (e: any) {
    message.error('切换失败: ' + (e.message || ''))
  }
}

async function handleDelete(id: number) {
  try {
    await forwardApi.deleteRule(id)
    message.success('规则已删除')
    fetchRules()
  } catch (e: any) {
    message.error('删除失败: ' + (e.message || ''))
  }
}

const columns = [
  { title: '规则名称', dataIndex: 'name', key: 'name', width: 150 },
  { title: '协议', dataIndex: 'protocol', key: 'protocol', width: 100 },
  { title: '监听地址:端口', key: 'listen', width: 180 },
  { title: '目标主机:端口', key: 'target', width: 200 },
  { title: '活跃连接', key: 'active_conns', width: 100 },
  { title: '累积流量 (下行/上行)', key: 'traffic', width: 180 },
  { title: 'UPnP', key: 'upnp', width: 80 },
  { title: '状态', key: 'enabled', width: 100 },
  { title: '操作', key: 'action', width: 140 },
]

onMounted(() => {
  fetchRules()
  timer = setInterval(fetchRules, 3000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="p-4 space-y-4">
    <ACard title="TCP / UDP 端口转发与 NAT 中继" :bordered="false" class="shadow-sm">
      <template #extra>
        <AFlex gap="small">
          <AButton @click="fetchRules">刷新</AButton>
          <AButton type="primary" @click="openAddModal">+ 添加转发规则</AButton>
        </AFlex>
      </template>

      <AAlert
        type="info"
        show-icon
        class="mb-4"
        message="端口转发与防火墙自动联动"
        description="新建端口转发时若勾选【自动在防火墙放行该端口】，平台将自动在系统防火墙打通对应监听端口，避免二次配置。"
      />

      <ATable
        :columns="columns"
        :data-source="rules"
        row-key="id"
        :pagination="{ pageSize: 10 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <span class="font-medium">{{ record.name }}</span>
            <div class="text-xs text-gray-400">{{ record.description || '无备注' }}</div>
          </template>

          <template v-else-if="column.key === 'protocol'">
            <ATag :color="record.protocol === 'tcp' ? 'blue' : record.protocol === 'udp' ? 'purple' : 'cyan'">
              {{ record.protocol.toUpperCase() }}
            </ATag>
          </template>

          <template v-else-if="column.key === 'listen'">
            <span class="font-mono text-blue-600 font-semibold">{{ record.listen_ip }}:{{ record.listen_port }}</span>
          </template>

          <template v-else-if="column.key === 'target'">
            <span class="font-mono text-emerald-600 font-semibold">{{ record.target_ip }}:{{ record.target_port }}</span>
          </template>

          <template v-else-if="column.key === 'active_conns'">
            <ATag :color="record.active_conns > 0 ? 'green' : 'default'">
              {{ record.active_conns || 0 }}
            </ATag>
          </template>

          <template v-else-if="column.key === 'traffic'">
            <div class="text-xs font-mono">
              <div>↓ {{ formatBytes(record.rx_bytes) }}</div>
              <div>↑ {{ formatBytes(record.tx_bytes) }}</div>
            </div>
          </template>

          <template v-else-if="column.key === 'upnp'">
            <ATag :color="record.enable_upnp ? 'green' : 'default'">
              {{ record.enable_upnp ? 'ON' : 'OFF' }}
            </ATag>
          </template>

          <template v-else-if="column.key === 'enabled'">
            <ASwitch
              :checked="record.enabled"
              size="small"
              @change="() => handleToggle(record)"
            />
          </template>

          <template v-else-if="column.key === 'action'">
            <AFlex gap="small">
              <AButton type="link" size="small" @click="openEditModal(record)">编辑</AButton>
              <APopconfirm
                title="确定删除此转发规则吗？"
                @confirm="handleDelete(record.id)"
              >
                <AButton type="link" danger size="small">删除</AButton>
              </APopconfirm>
            </AFlex>
          </template>
        </template>
      </ATable>
    </ACard>

    <!-- Add/Edit Modal -->
    <AModal
      v-model:open="showModal"
      :title="isEditing ? '编辑端口转发规则' : '添加端口转发规则'"
      ok-text="保存"
      cancel-text="取消"
      @ok="handleSubmit"
    >
      <AForm layout="vertical" class="mt-4">
        <AFormItem label="规则名称" required>
          <AInput v-model:value="form.name" placeholder="例如: Web项目端口转发、Minecraft服务器中继" />
        </AFormItem>

        <AFlex gap="middle">
          <AFormItem label="传输协议" class="flex-1" required>
            <ASelect v-model:value="form.protocol">
              <ASelectOption value="tcp">TCP</ASelectOption>
              <ASelectOption value="udp">UDP</ASelectOption>
              <ASelectOption value="both">TCP/UDP (双协议)</ASelectOption>
            </ASelect>
          </AFormItem>

          <AFormItem label="监听端口" class="flex-1" required>
            <AInputNumber v-model:value="form.listen_port" class="w-full" :min="1" :max="65535" />
          </AFormItem>
        </AFlex>

        <AFlex gap="middle">
          <AFormItem label="目标主机 / 局域网 IP" class="flex-1" required>
            <AInput v-model:value="form.target_ip" placeholder="例如: 127.0.0.1 或 192.168.1.100" />
          </AFormItem>

          <AFormItem label="目标端口" class="flex-1" required>
            <AInputNumber v-model:value="form.target_port" class="w-full" :min="1" :max="65535" />
          </AFormItem>
        </AFlex>

        <AFlex gap="middle" class="py-2">
          <ACheckbox v-model:checked="form.auto_open_firewall">
            自动在防火墙放行该端口 (推荐)
          </ACheckbox>
          <ACheckbox v-model:checked="form.enable_upnp">
            尝试 UPnP 路由器映射
          </ACheckbox>
        </AFlex>

        <AFormItem label="备注说明">
          <AInput v-model:value="form.description" placeholder="用途说明" />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>
