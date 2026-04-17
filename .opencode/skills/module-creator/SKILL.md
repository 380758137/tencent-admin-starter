---
name: module-creator
description: 使用 spec 驱动生成前后端 CRUD 模块，默认 dry-run
---

# Module Creator (OpenCode)

当开发者输入模块需求时：

1. 先写入 `specs/modules/<module>.module.spec.json`
2. 执行 dry-run
3. 展示变更计划
4. 经确认后执行 apply

命令：

```bash
npm run module:plan -- --spec specs/modules/<module>.module.spec.json
npm run module:apply -- --spec specs/modules/<module>.module.spec.json
```

