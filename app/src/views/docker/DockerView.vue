<script setup lang="ts">
import { ref, onMounted } from 'vue'
const { message } = App.useApp()
import dockerApi, {
  type DockerStatus,
  type DockerContainer,
  type DockerImage,
  type ComposeStack,
} from '@/api/docker'

const activeTab = ref('containers')
const status = ref<DockerStatus>({
  available: false,
  server_version: '',
  containers: 0,
  containers_running: 0,
  containers_paused: 0,
  containers_stopped: 0,
  images: 0,
  message: '',
})

const containers = ref<DockerContainer[]>([])
const images = ref<DockerImage[]>([])
const stacks = ref<ComposeStack[]>([])

const loading = ref(false)

// Logs modal
const showLogsModal = ref(false)
const currentLogs = ref('')
const currentContainerName = ref('')

// Compose modal
const showComposeModal = ref(false)
const composeForm = ref({
  name: '',
  content: `version: '3.8'
services:
  app:
    image: nginx:alpine
    ports:
      - "8081:80"
    restart: always
`,
})

async function fetchStatus() {
  try {
    const res = await dockerApi.getStatus()
    status.value = res.data || res
  } catch (e: any) {
    console.error(e)
  }
}

async function fetchContainers() {
  loading.value = true
  try {
    const res = await dockerApi.getContainers()
    containers.value = (res.data && res.data.data) ? res.data.data : (res.data || [])
  } catch (e: any) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function fetchImages() {
  try {
    const res = await dockerApi.getImages()
    images.value = (res.data && res.data.data) ? res.data.data : (res.data || [])
  } catch (e: any) {
    console.error(e)
  }
}

async function fetchCompose() {
  try {
    const res = await dockerApi.getCompose()
    stacks.value = res.data || res || []
  } catch (e: any) {
    console.error(e)
  }
}

async function handleAction(containerId: string, action: string) {
  try {
    message.loading({ content: `正在执行 ${action}...`, key: 'action' })
    await dockerApi.containerAction(containerId, action)
    message.success({ content: `容器操作 [${action}] 执行成功`, key: 'action' })
    fetchContainers()
    fetchStatus()
  } catch (e: any) {
    message.error({ content: '操作失败: ' + (e.message || ''), key: 'action' })
  }
}

async function openLogs(record: DockerContainer) {
  currentContainerName.value = record.names[0] || record.id
  currentLogs.value = '正在拉取日志...'
  showLogsModal.value = true
  try {
    const res = await dockerApi.getLogs(record.id)
    currentLogs.value = (res.data && res.data.logs) ? res.data.logs : (res.logs || '暂无日志输出')
  } catch (e: any) {
    currentLogs.value = '拉取日志失败: ' + (e.message || '')
  }
}

async function handleSaveCompose() {
  if (!composeForm.value.name || !composeForm.value.content) {
    message.warning('请填写项目名称与 YAML 配置')
    return
  }
  try {
    await dockerApi.saveCompose(composeForm.value.name, composeForm.value.content)
    message.success('Compose 栈保存成功')
    showComposeModal.value = false
    fetchCompose()
  } catch (e: any) {
    message.error('保存失败: ' + (e.message || ''))
  }
}

async function handleComposeAction(name: string, action: string) {
  try {
    message.loading({ content: `正在部署/执行 ${name} (${action})...`, key: 'compose' })
    const res = await dockerApi.composeAction(name, action)
    message.success({ content: '执行成功', key: 'compose' })
    fetchCompose()
    fetchContainers()
  } catch (e: any) {
    message.error({ content: '执行失败: ' + (e.message || ''), key: 'compose' })
  }
}

function formatSize(bytes: number) {
  if (!bytes) return '0 B'
  const mb = bytes / (1024 * 1024)
  if (mb < 1024) return mb.toFixed(1) + ' MB'
  return (mb / 1024).toFixed(2) + ' GB'
}

onMounted(() => {
  fetchStatus()
  fetchContainers()
  fetchImages()
  fetchCompose()
})
</script>

<template>
  <div class="p-4 space-y-4">
    <!-- Header Status Banner -->
    <ACard :bordered="false" class="shadow-sm">
      <AFlex justify="space-between" align="center" wrap="wrap" gap="middle">
        <AFlex align="center" gap="middle">
          <span class="text-base font-semibold">Docker 守护进程</span>
          <ATag :color="status.available ? 'green' : 'red'">
            {{ status.available ? '已连接 (ONLINE)' : '未连接 (OFFLINE)' }}
          </ATag>
          <span v-if="status.server_version" class="text-xs text-gray-500 font-mono">
            版本: v{{ status.server_version }}
          </span>
          <span class="text-xs text-gray-400">
            {{ status.message }}
          </span>
        </AFlex>

        <AFlex gap="large" v-if="status.available">
          <div class="text-center">
            <div class="text-lg font-bold text-blue-600">{{ status.containers }}</div>
            <div class="text-xs text-gray-400">总容器</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold text-emerald-600">{{ status.containers_running }}</div>
            <div class="text-xs text-gray-400">运行中</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold text-gray-500">{{ status.containers_stopped }}</div>
            <div class="text-xs text-gray-400">已停止</div>
          </div>
          <div class="text-center">
            <div class="text-lg font-bold text-purple-600">{{ status.images }}</div>
            <div class="text-xs text-gray-400">本地镜像</div>
          </div>
        </AFlex>
      </AFlex>
    </ACard>

    <!-- Main Tabs -->
    <ACard :bordered="false" class="shadow-sm">
      <ATabs v-model:activeKey="activeTab">
        <!-- Tab 1: Containers -->
        <ATabPane key="containers" tab="容器管理 (Containers)">
          <div class="flex justify-end mb-3">
            <AButton @click="fetchContainers">刷新容器列表</AButton>
          </div>

          <ATable
            :data-source="containers"
            :loading="loading"
            row-key="id"
            :pagination="{ pageSize: 10 }"
          >
            <ATableColumn title="容器名称" key="names">
              <template #default="{ record }">
                <span class="font-medium text-blue-600">
                  {{ record.names ? record.names.join(', ').replace(/^\//, '') : record.id }}
                </span>
                <div class="font-mono text-xs text-gray-400">ID: {{ record.id }}</div>
              </template>
            </ATableColumn>

            <ATableColumn title="镜像" data-index="image" key="image" />

            <ATableColumn title="运行状态" key="state" width="120">
              <template #default="{ record }">
                <ATag :color="record.state === 'running' ? 'green' : 'default'">
                  {{ record.state ? record.state.toUpperCase() : 'UNKNOWN' }}
                </ATag>
              </template>
            </ATableColumn>

            <ATableColumn title="端口映射" key="ports">
              <template #default="{ record }">
                <span class="font-mono text-xs">{{ record.ports ? record.ports.join(', ') : '-' }}</span>
              </template>
            </ATableColumn>

            <ATableColumn title="操作" key="action" width="220">
              <template #default="{ record }">
                <AFlex gap="small">
                  <AButton
                    v-if="record.state !== 'running'"
                    size="small"
                    type="link"
                    @click="handleAction(record.id, 'start')"
                  >
                    启动
                  </AButton>
                  <AButton
                    v-if="record.state === 'running'"
                    size="small"
                    type="link"
                    danger
                    @click="handleAction(record.id, 'stop')"
                  >
                    停止
                  </AButton>
                  <AButton
                    size="small"
                    type="link"
                    @click="handleAction(record.id, 'restart')"
                  >
                    重启
                  </AButton>
                  <AButton size="small" type="link" @click="openLogs(record)">日志</AButton>
                  <APopconfirm title="确定删除容器？" @confirm="handleAction(record.id, 'remove')">
                    <AButton size="small" type="link" danger>删除</AButton>
                  </APopconfirm>
                </AFlex>
              </template>
            </ATableColumn>
          </ATable>
        </ATabPane>

        <!-- Tab 2: Compose Stacks -->
        <ATabPane key="compose" tab="Compose 栈编排">
          <div class="flex justify-between mb-3 items-center">
            <span class="text-sm text-gray-500">通过统一 YAML 文件对多容器进行版本化管理与一键启停</span>
            <AFlex gap="small">
              <AButton @click="fetchCompose">刷新</AButton>
              <AButton type="primary" @click="showComposeModal = true">+ 新建 Compose 栈</AButton>
            </AFlex>
          </div>

          <ATable :data-source="stacks" row-key="name" :pagination="{ pageSize: 10 }">
            <ATableColumn title="项目名称" data-index="name" key="name">
              <template #default="{ record }">
                <span class="font-semibold text-blue-600">{{ record.name }}</span>
              </template>
            </ATableColumn>

            <ATableColumn title="存放路径" data-index="path" key="path">
              <template #default="{ record }">
                <span class="font-mono text-xs text-gray-500">{{ record.path }}</span>
              </template>
            </ATableColumn>

            <ATableColumn title="操作" key="action" width="220">
              <template #default="{ record }">
                <AFlex gap="small">
                  <AButton size="small" type="primary" @click="handleComposeAction(record.name, 'up')">
                    部署 / 启动
                  </AButton>
                  <AButton size="small" @click="handleComposeAction(record.name, 'restart')">
                    重启
                  </AButton>
                  <AButton size="small" danger @click="handleComposeAction(record.name, 'down')">
                    停止销毁
                  </AButton>
                </AFlex>
              </template>
            </ATableColumn>
          </ATable>
        </ATabPane>

        <!-- Tab 3: Images -->
        <ATabPane key="images" tab="本地镜像 (Images)">
          <div class="flex justify-end mb-3">
            <AButton @click="fetchImages">刷新镜像列表</AButton>
          </div>

          <ATable :data-source="images" row-key="id" :pagination="{ pageSize: 10 }">
            <ATableColumn title="镜像 ID" data-index="id" key="id">
              <template #default="{ record }">
                <span class="font-mono text-xs">{{ record.id }}</span>
              </template>
            </ATableColumn>
            <ATableColumn title="标签 (RepoTags)" key="repo_tags">
              <template #default="{ record }">
                <span class="font-semibold text-blue-600">
                  {{ record.repo_tags ? record.repo_tags.join(', ') : '<none>' }}
                </span>
              </template>
            </ATableColumn>
            <ATableColumn title="镜像大小" key="size">
              <template #default="{ record }">
                <span>{{ formatSize(record.size) }}</span>
              </template>
            </ATableColumn>
          </ATable>
        </ATabPane>
      </ATabs>
    </ACard>

    <!-- Logs Modal -->
    <AModal
      v-model:open="showLogsModal"
      :title="'容器实时日志: ' + currentContainerName"
      width="800px"
      :footer="null"
    >
      <div class="bg-gray-900 text-green-400 p-4 rounded font-mono text-xs h-96 overflow-y-auto whitespace-pre-wrap">
        {{ currentLogs }}
      </div>
    </AModal>

    <!-- Compose Modal -->
    <AModal
      v-model:open="showComposeModal"
      title="新建 / 编辑 Compose 项目"
      width="700px"
      ok-text="保存"
      cancel-text="取消"
      @ok="handleSaveCompose"
    >
      <AForm layout="vertical" class="mt-4">
        <AFormItem label="项目英文标识 (Stack Name)" required>
          <AInput v-model:value="composeForm.name" placeholder="例如: my-nginx, wordpress, redis-cluster" />
        </AFormItem>
        <AFormItem label="docker-compose.yml 配置文件内容" required>
          <AInputTextArea
            v-model:value="composeForm.content"
            :rows="12"
            class="font-mono text-xs"
            placeholder="写入标准 Docker Compose YAML 内容"
          />
        </AFormItem>
      </AForm>
    </AModal>
  </div>
</template>
