<template>
  <el-drawer
    v-model="visible"
    title="新增文献"
    direction="rtl"
    size="40%"
    :before-close="handleClose"
  >
    <el-form :model="form" label-width="80px">
      <el-form-item label="柜号">
        <el-input v-model="form.cabinetNo" />
      </el-form-item>
      <el-form-item label="盒号">
        <el-input v-model="form.boxNo" />
      </el-form-item>
      <el-form-item label="编号">
        <el-input v-model="form.innerNo" />
      </el-form-item>
      <el-form-item label="档案类型">
        <el-input v-model="form.arcType" />
      </el-form-item>
      <el-form-item label="合同编号">
        <el-input v-model="form.contractNo" />
      </el-form-item>
      <el-form-item label="姓名">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="身份证号">
        <el-input v-model="form.idCard" />
      </el-form-item>
      <el-form-item label="网点编号">
        <el-input v-model="form.branchNo" />
      </el-form-item>
      <el-form-item label="客户经理">
        <el-input v-model="form.manager" />
      </el-form-item>
      <el-form-item label="合同金额">
        <el-input-number v-model="form.amount" :min="0" style="width: 100%" />
      </el-form-item>
      <el-form-item label="权限类型">
        <el-input v-model="form.groupPermission" />
      </el-form-item>
      <el-form-item label="入库日期">
        <el-date-picker
          v-model="form.storageDate"
          type="date"
          placeholder="选择日期"
          style="width: 100%"
        />
      </el-form-item>
    </el-form>

    <!-- 底部操作按钮 -->
    <template #footer>
      <div style="text-align: right;">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit">提交</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, reactive, defineExpose } from "vue"

const visible = ref(false)

// 表单数据
const form = reactive({
  fileNo: "",
  cabinetNo: "",
  boxNo: "",
  innerNo: "",
  arcType: "",
  contractNo: "",
  name: "",
  idCard: "",
  branchNo: "",
  manager: "",
  amount: 0,
  groupPermission: "",
  storageDate: "",
  borrowStatus: 0, // 默认未借阅
})

// 向父组件暴露方法
const open = () => {
  visible.value = true
}

const handleClose = () => {
  visible.value = false
}

// 提交表单
const emit = defineEmits(["submit"])
const handleSubmit = () => {
  emit("submit", { ...form })
  handleClose()
}

// 让父组件可以调用 open()
defineExpose({
  open,
})
</script>
