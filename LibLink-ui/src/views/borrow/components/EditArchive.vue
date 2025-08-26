<template>
  <el-drawer
    v-model="visible"
    title="编辑档案"
    direction="rtl"
    size="40%"
    :before-close="handleClose"
  >
    <!-- 编辑表单 -->
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
    >
      <el-form-item label="档案编号" prop="file_no">
        <el-input v-model="form.file_no" placeholder="请输入档案编号" />
      </el-form-item>

      <el-form-item label="档案类型" prop="arc_type">
        <el-input v-model="form.arc_type" placeholder="请输入档案类型" />
      </el-form-item>

      <el-form-item label="合同编号" prop="contract_no">
        <el-input v-model="form.contract_no" placeholder="请输入合同编号" />
      </el-form-item>

      <el-form-item label="姓名" prop="name">
        <el-input v-model="form.name" placeholder="请输入姓名" />
      </el-form-item>

      <el-form-item label="身份证号" prop="id_card">
        <el-input v-model="form.id_card" placeholder="请输入身份证号" />
      </el-form-item>

      <el-form-item label="网点编号" prop="inst_no">
        <el-input v-model="form.inst_no" placeholder="请输入网点编号" />
      </el-form-item>

      <el-form-item label="客户经理" prop="manager">
        <el-input v-model="form.manager" placeholder="请输入客户经理" />
      </el-form-item>

      <el-form-item label="合同金额" prop="amount">
        <el-input v-model="form.amount" placeholder="请输入合同金额" />
      </el-form-item>

      <el-form-item label="入库日期" prop="storage_date">
        <el-date-picker
          v-model="form.storage_date"
          type="date"
          placeholder="选择入库日期"
          style="width: 100%"
        />
      </el-form-item>
    </el-form>

    <!-- 底部操作 -->
    <template #footer>
      <div style="text-align: right; padding: 10px 0">
        <el-button @click="handleClose">取 消</el-button>
        <el-button type="primary" @click="handleSubmit">确 定</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue"
import { ElMessage } from "element-plus"
import { updateArchive } from "@/api/archives"

const visible = ref(false)
const formRef = ref()

// 表单数据
const form = reactive<any>({
  file_no: "",
  arc_type: "",
  contract_no: "",
  name: "",
  id_card: "",
  inst_no: "",
  manager: "",
  amount: "",
  storage_date: "",
})

// 校验规则
const rules = {
  file_no: [{ required: true, message: "档案编号不能为空", trigger: "blur" }],
  contract_no: [{ required: true, message: "合同编号不能为空", trigger: "blur" }],
}

// 关闭抽屉
const handleClose = () => {
  visible.value = false
}

// 提交表单
const emit = defineEmits(["success"])
const handleSubmit = () => {
  formRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    try {
      console.log(form.contract_no)
      await updateArchive(form.contract_no, form)  // 调用 API 更新
      ElMessage.success("档案更新成功")
      handleClose()
      emit("success")  // 通知父组件刷新数据
    } catch (error) {
      ElMessage.error("档案更新失败")
      console.error(error)
    }
  })
}

// 打开并填充数据
const open = (data: any) => {
  Object.assign(form, data)  // 直接覆盖
  visible.value = true
}

// 向父组件暴露方法
defineExpose({
  open,
})
</script>
