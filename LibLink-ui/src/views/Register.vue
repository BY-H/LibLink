<template>
    <div class="login-wrapper">
        <div class="login-container">
            <div class="form-header">
                <h2>用户注册</h2>
                <p>您好，请注册您的账号</p>
            </div>
            <form class="floating-form" @submit.prevent="handleRegister">
                <div class="input-group">
                    <input id="username" type="text" autocomplete="off" required v-model="form.email" />
                    <label for="username">邮箱</label>
                    <span class="highlight"></span>
                </div>
                <div class="input-group">
                    <input id="password" type="password" autocomplete="off" required v-model="form.password" />
                    <label for="password">密码</label>
                    <span class="highlight"></span>
                </div>
                <button type="submit" class="submit-btn" :disabled="!isFormValid">
                    <span>注册</span>
                    <i class="arrow-icon"></i>
                </button>
                <div class="form-footer">
                    <span>已经有账号？</span>
                    <el-link @click="jumpToLogin">立即登录</el-link>
                </div>
            </form>
        </div>
    </div>
</template>
<script lang="ts" name="Register" setup>
import { ref, reactive, computed } from 'vue'
import { register } from '@/api/user'
import { setToken } from '@/utils/auth'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
const router = useRouter()

const form = reactive({
    email: '',
    password: ''
})

const isFormValid = computed(() => {
    return form.email && form.password
})

const handleRegister = async () => {
    const user = {
        email: form.email,
        password: form.password
    }
    const response = await register(user)
    console.log(response)
    if (response.status === 200) {
        ElMessage.info('注册成功，请登录')
        router.push({ path: '/login' })
    } else {
        ElMessage.error('注册失败，请重试')
    }
}

const jumpToLogin = () => {
    router.push({ path: '/login' })
}
</script>
<style scoped>
.login-wrapper {
    min-height: 93.5vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
}

.login-container {
    width: 100%;
    max-width: 480px;
    background: white;
    border-radius: 20px;
    padding: 40px;
    padding-right: 80px;
    box-shadow: 0 15px 35px rgba(0, 0, 0, 0.1);
}

.form-header {
    text-align: center;
    margin-bottom: 40px;
}

.form-header h2 {
    color: #2c3e50;
    font-size: 32px;
    margin-bottom: 10px;
    font-weight: 700;
}

.form-header p {
    color: #95a5a6;
    font-size: 16px;
}

.floating-form .input-group {
    position: relative;
    margin-bottom: 30px;
}

.input-group input {
    width: 100%;
    padding: 15px;
    border: 2px solid #e0e0e0;
    border-radius: 12px;
    font-size: 16px;
    transition: all 0.3s ease;
    background: transparent;
}

.input-group label {
    position: absolute;
    left: 15px;
    top: 50%;
    transform: translateY(-50%);
    background: white;
    padding: 0 5px;
    color: #95a5a6;
    font-size: 16px;
    transition: all 0.3s ease;
    pointer-events: none;
}

.input-group input:focus,
.input-group input:valid {
    border-color: #3498db;
}

.input-group input:focus + label,
.input-group input:valid + label {
    top: 0;
    font-size: 14px;
    color: #3498db;
}

.submit-btn {
    width: 100%;
    padding: 15px;
    margin-left: 15px;
    background: linear-gradient(to right, #3498db, #2980b9);
    color: white;
    border: none;
    border-radius: 12px;
    font-size: 18px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 10px;
}

.submit-btn:hover {
    transform: translateY(-2px);
    box-shadow: 0 5px 15px rgba(52, 152, 219, 0.3);
}

.arrow-icon {
    border: solid white;
    border-width: 0 2px 2px 0;
    display: inline-block;
    padding: 3px;
    transform: rotate(-45deg);
}

.form-footer {
    text-align: center;
    margin-top: 20px;
    color: #95a5a6;
}

.form-footer a {
    color: #3498db;
    text-decoration: none;
    margin-left: 5px;
    font-weight: 600;
}

.form-footer a:hover {
    text-decoration: underline;
}

.error-message {
    color: #f56c6c;
    font-size: 14px;
    text-align: center;
    margin-bottom: 20px;
}

@media (max-width: 480px) {
    .login-container {
        padding: 20px;
    }

    .form-header h2 {
        font-size: 24px;
    }

    .input-group input {
        padding: 12px;
    }
}

@media (max-width: 768px) {
    .login-container {
        max-width: 400px;
        padding: 30px;
    }

    .form-header h2 {
        font-size: 28px;
    }

    .form-header p {
        font-size: 14px;
    }
}

@media (max-width: 480px) {
    .login-container {
        padding: 20px;
        margin: 10px;
        max-width: 100%;
    }

    .form-header h2 {
        font-size: 24px;
    }

    .form-header p {
        font-size: 14px;
    }

    .input-group input {
        padding: 12px;
        font-size: 14px;
    }

    .input-group label {
        font-size: 14px;
    }

    .submit-btn {
        padding: 12px;
        font-size: 16px;
    }
}

@media (max-width: 320px) {
    .login-container {
        padding: 15px;
    }

    .form-header h2 {
        font-size: 20px;
    }

    .input-group {
        margin-bottom: 20px;
    }
}
</style>
