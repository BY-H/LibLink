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
        <el-input v-model="form.cabinet_no" />
      </el-form-item>
      <el-form-item label="盒号">
        <el-input v-model="form.box_no" />
      </el-form-item>
      <el-form-item label="编号">
        <el-input v-model="form.inner_no" />
      </el-form-item>
      <el-form-item label="档案类型">
        <el-input v-model="form.arc_type" />
      </el-form-item>
      <el-form-item label="合同编号">
        <el-input v-model="form.contract_no" />
      </el-form-item>
      <el-form-item label="姓名">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="身份证号">
        <el-input v-model="form.id_card" />
      </el-form-item>
      <el-form-item label="网点编号">
        <el-input v-model="form.inst_no" />
      </el-form-item>
      <el-form-item label="客户经理">
        <el-input v-model="form.manager" />
      </el-form-item>
      <el-form-item label="合同金额">
        <el-input-number v-model="form.amount" :min="0" style="width: 100%" />
      </el-form-item>
      <!-- <el-form-item label="权限类型">
        <el-input v-model="form.group_permission" />
      </el-form-item> -->
      <el-form-item label="入库日期">
        <el-date-picker
          v-model="form.storage_date"
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

// 初始表单
const initialForm = {
  file_no: "",
  cabinet_no: "",
  box_no: "",
  inner_no: "",
  arc_type: "",
  contract_no: "",
  name: "",
  id_card: "",
  inst_no: "",
  manager: "",
  amount: "0",
  group_permission: "",
  storage_date: "",
  borrow_state: "0",
}

// 表单数据
const form = reactive({ ...initialForm })

// 重置方法
const resetForm = () => {
  Object.assign(form, initialForm)
}

// 向父组件暴露方法
const open = () => {
  visible.value = true
}

const handleClose = () => {
  visible.value = false
  resetForm()
}

// 提交表单
const emit = defineEmits(["submit"])
const handleSubmit = () => {
  // 拼接 file_no
  form.file_no = `${form.cabinet_no}${form.box_no}-${form.inner_no}`
  form.amount = String(form.amount) // 转为字符串
  emit("submit", { ...form })
  handleClose()
}

// 让父组件可以调用 open()
defineExpose({
  open,
})
</script>
