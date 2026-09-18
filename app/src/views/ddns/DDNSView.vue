<script setup lang="ts">
import type { DDNSTask, InterfaceInfo } from '@/api/ddns'
import { onMounted, ref } from 'vue'
import ddnsApi from '@/api/ddns'

const { message } = App.useApp()

const loading = ref(false)
const tasks = ref<DDNSTask[]>([])
const interfaces = ref<InterfaceInfo[]>([])
const showModal = ref(false)
const isEditing = ref(false)
const editId = ref(0)

const form = ref({
  name: '',
  provider: 'cloudflare',
  domains: '',
  ip_type: 'ipv4',
  ip_method: 'url',
  ip_interface: '',
  ip_url: '',
  access_key_id: '',
  access_key_secret: '',
  webhook_url: '',
  interval_seconds: 300,
  enabled: true,
})

async function fetchTasks() {
  loading.value = true
  try {
    const res = await ddnsApi.getTasks()
    tasks.value = res || []
  }
  catch (e: any) {
    console.error(e)
  }
  finally {
    loading.value = false
  }
}

async function fetchInterfaces() {
  try {
    const res = await ddnsApi.getInterfaces()
    interfaces.value = res || []
  }
  catch (e: any) {
    console.error(e)
  }
}

function openAddModal() {
  isEditing.value = false
  editId.value = 0
  form.value = {
    name: '',
    provider: 'cloudflare',
    domains: '',
    ip_type: 'ipv4',
    ip_method: 'url',
    ip_interface: interfaces.value.length > 0 ? interfaces.value[0].name : '',
    ip_url: '',
    access_key_id: '',
    access_key_secret: '',
    webhook_url: '',
    interval_seconds: 300,
    enabled: true,
  }
  showModal.value = true
}

function openEditModal(record: DDNSTask) {
  isEditing.value = true
  editId.value = record.id
  form.value = {
    name: record.name,
    provider: record.provider,
    domains: record.domains,
    ip_type: record.ip_type,
    ip_method: record.ip_method,
    ip_interface: record.ip_interface,
    ip_url: record.ip_url,
    access_key_id: record.access_key_id || '',
    access_key_secret: record.access_key_secret || '',
    webhook_url: record.webhook_url || '',
    interval_seconds: record.interval_seconds,
    enabled: record.enabled,
  }
  showModal.value = true
}

async function handleSubmit() {
  if (!form.value.name || !form.value.domains) {
    message.warning('请填写任务名称及待绑定的域名')
    return
  }
  try {
    if (isEditing.value) {
      await ddnsApi.updateTask(editId.value, form.value)
      message.success('DDNS任务更新成功')
    }
    else {
      await ddnsApi.createTask(form.value)
      message.success('DDNS任务创建成功')
    }
    showModal.value = false
    fetchTasks()
  }
  catch (e: any) {
    message.error(`操作失败: ${e.message || ''}`)
  }
}

async function handleRunNow(id: number) {
  try {
    message.loading({ content: '正在触发同步解析...', key: 'sync' })
    await ddnsApi.runTaskNow(id)
    message.success({ content: '解析记录同步完成', key: 'sync' })
    fetchTasks()
  }
  catch (e: any) {
    message.error({ content: `同步失败: ${e.message || ''}`, key: 'sync' })
  }
}

async function handleDelete(id: number) {
  try {
    await ddnsApi.deleteTask(id)
    message.success('任务已删除')
    fetchTasks()
  }
  catch (e: any) {
    message.error(`删除失败: ${e.message || ''}`)
  }
}

const columns = [
  { title: '任务名称', dataIndex: 'name', key: 'name', width: 140 },
  { title: 'DNS 服务商', dataIndex: 'provider', key: 'provider', width: 120 },
  { title: '同步域名', dataIndex: 'domains', key: 'domains' },
  { title: '类型/获取源', key: 'method', width: 140 },
  { title: '最新解析 IP', key: 'current_ip', width: 160 },
  { title: '上次同步状态', key: 'status', width: 130 },
  { title: '同步间隔', dataIndex: 'interval_seconds', key: 'interval', width: 100 },
  { title: '操作', key: 'action', width: 160 },
]

onMounted(() => {
  fetchTasks()
  fetchInterfaces()
})
</script>

<template>
  <div class="p-4 space-y-4">
    <ACard title="动态域名解析 (DDNS)" :bordered="false" class="shadow-sm">
      <template #extra>
        <AFlex gap="small">
          <AButton @click="fetchTasks">
            刷新
          </AButton>
          <AButton type="primary" @click="openAddModal">
            + 添加 DDNS 任务
          </AButton>
        </AFlex>
      </template>

      <AAlert
        type="info"
        show-icon
        class="mb-4"
        message="智能 IP 变动探测"
        description="支持 IPv4 / IPv6 双栈探测与多网卡过滤，仅在公网 IP 发生实质变动时向云厂商 API 下发变更，安全高效。"
      />

      <ATable
        :columns="columns"
        :data-source="tasks"
        :loading="loading"
        row-key="id"
        :pagination="{ pageSize: 10 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'provider'">
            <ATag color="blue">
              {{ record.provider.toUpperCase() }}
            </ATag>
          </template>

          <template v-else-if="column.key === 'domains'">
            <div class="font-mono text-blue-600 font-semibold">
              {{ record.domains }}
            </div>
          </template>

          <template v-else-if="column.key === 'method'">
            <div>
              <ATag color="cyan">
                {{ record.ip_type.toUpperCase() }}
              </ATag>
              <span class="text-xs text-gray-500">{{ record.ip_method }}</span>
            </div>
          </template>

          <template v-else-if="column.key === 'current_ip'">
            <span class="font-mono font-medium text-emerald-600">
              {{ record.ip_type === 'ipv6' ? (record.last_ipv6 || '未获取') : (record.last_ipv4 || '未获取') }}
            </span>
          </template>

          <template v-else-if="column.key === 'status'">
            <div>
              <ATag :color="record.last_status === 'success' ? 'green' : record.last_status === 'failed' ? 'red' : 'default'">
                {{ record.last_status ? record.last_status.toUpperCase() : 'PENDING' }}
              </ATag>
              <div v-if="record.last_error" class="text-xs text-red-500 truncate max-w-xs" :title="record.last_error">
                {{ record.last_error }}
              </div>
            </div>
          </template>

          <template v-else-if="column.key === 'interval'">
            <span>{{ record.interval_seconds }} 秒</span>
          </template>

          <template v-else-if="column.key === 'action'">
            <AFlex gap="small">
              <AButton type="link" size="small" @click="handleRunNow(record.id)">
                立即同步
              </AButton>
              <AButton type="link" size="small" @click="openEditModal(record)">
                编辑
              </AButton>
              <APopconfirm title="确认删除此任务？" @confirm="handleDelete(record.id)">
                <AButton type="link" danger size="small">
                  删除
                </AButton>
              </APopconfirm>
            </AFlex>
          </template>
        </template>
      </ATable>
    </ACard>

    <!-- Modal -->
    <AModal
      v-model:open="showModal"
      :title="isEditing ? '编辑 DDNS 任务' : '添加 DDNS 任务'"
      ok-text="保存"
      cancel-text="取消"
      @ok="handleSubmit"
    >
      <AForm layout="vertical" class="mt-4">
        <AFormItem label="任务名称" required>
          <AInput v-model:value="form.name" placeholder="例如: 家庭NAS解析、办公网IPv6直连" />
        </AFormItem>

        <AFlex gap="middle">
          <AFormItem label="DNS 服务商" class="flex-1" required>
            <ASelect v-model:value="form.provider">
              <ASelectOption value="cloudflare">
                Cloudflare
              </ASelectOption>
              <ASelectOption value="aliyun">
                阿里云 (AliDNS)
              </ASelectOption>
              <ASelectOption value="tencent">
                腾讯云 (DNSPod)
              </ASelectOption>
              <ASelectOption value="huawei">
                华为云
              </ASelectOption>
              <ASelectOption value="callback">
                自定义 Webhook
              </ASelectOption>
            </ASelect>
          </AFormItem>

          <AFormItem label="IP 类型" class="flex-1" required>
            <ASelect v-model:value="form.ip_type">
              <ASelectOption value="ipv4">
                IPv4 (A 记录)
              </ASelectOption>
              <ASelectOption value="ipv6">
                IPv6 (AAAA 记录)
              </ASelectOption>
            </ASelect>
          </AFormItem>
        </AFlex>

        <AFormItem label="要解析的域名 (多个用逗号隔开)" required>
          <AInput v-model:value="form.domains" placeholder="例如: sub.example.com, home.example.com" />
        </AFormItem>

        <AFlex gap="middle">
          <AFormItem label="IP 探测方式" class="flex-1" required>
            <ASelect v-model:value="form.ip_method">
              <ASelectOption value="url">
                公网探针 API
              </ASelectOption>
              <ASelectOption value="nic">
                本地网卡直读
              </ASelectOption>
            </ASelect>
          </AFormItem>

          <AFormItem v-if="form.ip_method === 'nic'" label="选择网卡设备" class="flex-1">
            <ASelect v-model:value="form.ip_interface">
              <ASelectOption v-for="nic in interfaces" :key="nic.name" :value="nic.name">
                {{ nic.name }} ({{ nic.ipv4s.join(', ') || '无IP' }})
              </ASelectOption>
            </ASelect>
          </AFormItem>
        </AFlex>

        <template v-if="form.provider === 'cloudflare'">
          <AFormItem label="Cloudflare API Token" required>
            <AInputPassword v-model:value="form.access_key_id" placeholder="填入具有 Zone.DNS 编辑权限的 Token" />
          </AFormItem>
        </template>

        <template v-else-if="form.provider === 'aliyun' || form.provider === 'tencent' || form.provider === 'huawei'">
          <AFormItem label="Access Key ID" required>
            <AInput v-model:value="form.access_key_id" placeholder="AccessKey ID" />
          </AFormItem>
          <AFormItem label="Access Key Secret" required>
            <AInputPassword v-model:value="form.access_key_secret" placeholder="AccessKey Secret" />
          </AFormItem>
        </template>

        <template v-else-if="form.provider === 'callback'">
          <AFormItem label="Webhook URL (支持 #{ip}, #{domain}, #{type} 变量)" required>
            <AInput v-model:value="form.webhook_url" placeholder="https://api.myhook.com/update?ip=#{ip}&domain=#{domain}" />
          </AFormItem>
        </template>

        <AFormItem label="同步周期 (秒)">
          <AInputNumber v-model:value="form.interval_seconds" class="w-full" :min="30" :max="86400" />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>
