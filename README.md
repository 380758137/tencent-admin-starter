# Tencent Admin Starter

传统 CRUD 后台管理系统骨架，包含：

- 前端：Vue3 + TDesign
- 后端：Go + Gin + GORM + MySQL + JWT + RBAC
- 基础平台底座：菜单管理、角色管理、字典管理、参数中心、日志中心、系统监控
- 传统平台补齐：岗位管理、通知公告、在线用户、定时任务（含执行日志）
- 第三批基础能力：角色权限码（按钮级权限）、用户数据权限（按部门及子部门/本人）、用户 CSV 导入导出
- 平台闭环增强：菜单改为后端数据驱动（`/api/auth/menus`），并在各业务页统一应用按钮级权限控制
- 用户角色升级：支持一用户多角色（角色可多选分配，权限取并集，数据范围按 all > dept > self 合并）
- 数据权限可配置化：角色支持 dataScope（all/dept/self）并实时作用于用户数据访问范围
- 模块扩展：`module.spec` 驱动生成前后端代码（默认 dry-run），并自动接入生成模块路由、菜单与权限点
- 技能集成：内置 Claude/OpenCode 技能与其他客户端提示模板

## 目录

- `apps/web`：后台前端
- `apps/api`：后台 API
- `tools/modulegen`：模块生成器
- `specs/modules`：模块 spec
- `.claude/skills`、`.opencode/skills`：可安装技能

## 快速开始

```bash
cd /Users/liusheng/demos/tencent-admin-starter
npm install
cd apps/web && npm install
cd ../api && go mod tidy
```

## 启动

```bash
# 启动前端
npm run dev:web

# 启动后端
npm run dev:api
```

后端默认账号（首次自动种子）：

- 用户名：`admin`
- 密码：`Admin@123456`

## 模块生成（以部门模块为例）

```bash
# 1) 预览变更（dry-run）
npm run module:plan -- --spec specs/modules/department.module.spec.json

# 2) 确认后应用
npm run module:apply -- --spec specs/modules/department.module.spec.json
```

## 环境变量

后端参考：`apps/api/.env.example`

注意：`JWT_SECRET` 为必填，未配置时后端会拒绝启动。
