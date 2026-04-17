<script setup lang="ts">
import { reactive, ref } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const auth = useAuthStore()
const loading = ref(false)
const form = reactive({
  username: 'admin',
  password: 'Admin@123456'
})

async function onSubmit() {
  loading.value = true
  try {
    await auth.signIn(form)
    router.push('/')
  } catch {
    MessagePlugin.error('登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-shell">
    <section class="login-panel">
      <t-card class="shell-surface" title="登录后台系统">
        <t-form @submit="onSubmit" layout="vertical">
          <t-form-item label="用户名" name="username">
            <t-input v-model="form.username" placeholder="请输入用户名" />
          </t-form-item>
          <t-form-item label="密码" name="password">
            <t-input v-model="form.password" type="password" placeholder="请输入密码" />
          </t-form-item>
          <t-form-item style="margin-top: 6px">
            <t-button theme="primary" type="submit" :loading="loading" block>登录</t-button>
          </t-form-item>
          <div style="color: #86909c; font-size: 12px; margin-top: 6px">
            默认账号：admin / Admin@123456
          </div>
        </t-form>
      </t-card>
    </section>
  </div>
</template>
