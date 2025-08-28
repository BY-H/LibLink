# ======================
# 1. 构建阶段（Node）
# ======================
FROM node:22-alpine AS builder

# 工作目录
WORKDIR /app

# 复制依赖清单
COPY package.json package-lock.json ./

# 安装依赖（使用 npm ci 保证干净环境）
RUN npm ci

# 复制项目代码
COPY . .

# 解决 esbuild 兼容问题
RUN npm rebuild esbuild

# 构建前端
RUN npm run build-only


# ======================
# 2. 运行阶段（Nginx）
# ======================
FROM nginx:1.27.0-alpine

# 删除默认配置
RUN rm /etc/nginx/conf.d/default.conf

# 添加自定义 Nginx 配置
COPY ./nginx/default.conf /etc/nginx/conf.d/

# 复制前端打包产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 暴露端口
EXPOSE 80

# 启动 nginx
CMD ["nginx", "-g", "daemon off;"]
