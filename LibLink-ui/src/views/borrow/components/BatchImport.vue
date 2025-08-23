<template>
  <el-dialog v-model="visible" title="批量新增档案" width="500px" @close="handleClose">
    <div style="margin-bottom: 15px">
      请先下载模板文件并填写数据：
      <el-link type="primary" :underline="false" :href="templateUrl" download>下载模板</el-link>
    </div>

    <!-- 文件上传区域 -->
    <el-upload
      ref="uploadRef"
      :auto-upload="false"
      :limit="1"
      :before-upload="handleBeforeUpload"
      :on-change="handleFileChange"
      :file-list="fileList"
    >
      <el-button type="primary">选择文件</el-button>
    </el-upload>

    <template #footer>
      <span class="dialog-footer">
        <el-button @click="close">取消</el-button>
        <el-button type="primary" @click="submitUpload">上传</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request' 

const props = defineProps<{
  modelValue: boolean
}>()

const uploadUrl = ""
const templateUrl = `${import.meta.env.VITE_BASE_URL_API}/static/templates/archives_template.xlsx`

console.log("templateUrl:", templateUrl)

const emit = defineEmits(['update:modelValue', 'success'])

const visible = ref(props.modelValue)
watch(() => props.modelValue, val => visible.value = val)

const fileList = ref<any[]>([])  // 存储选中的文件

const close = () => {
  emit('update:modelValue', false)
}

const handleClose = () => {
  emit('update:modelValue', false)
}

const handleFileChange = (file: any, files: any[]) => {
  fileList.value = files
}

const handleBeforeUpload = () => {
  // 阻止 el-upload 默认的自动上传行为
  return false
}

const submitUpload = async () => {
  if (fileList.value.length === 0) {
    ElMessage.warning('请先选择文件')
    return
  }

  const formData = new FormData()
  formData.append('file', fileList.value[0].raw) 

  try {
    const res = await request({
      url: uploadUrl,
      method: 'post',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    ElMessage.success('批量导入成功')
    emit('success', res)
    close()
  } catch (err) {
    console.error('上传失败:', err)
    ElMessage.error('批量导入失败，请检查文件格式')
  }
}
</script>
