---
name: module-creator
description: 根据自然语言需求生成 module.spec，并触发前后端 CRUD 模块代码生成
---

# Module Creator Skill

## Goal
将“需求描述”转换为 `specs/modules/<module>.module.spec.json`，再调用项目内生成器完成落地。

## Workflow
1. 提取模块名、展示名、字段、权限点。
2. 生成 `module.spec.json` 并校验字段命名与类型。
3. 先执行 dry-run：
   `npm run module:plan -- --spec specs/modules/<module>.module.spec.json`
4. 用户确认后执行 apply：
   `npm run module:apply -- --spec specs/modules/<module>.module.spec.json`

## Field types
- `string`
- `int`
- `bool`
- `datetime`

## Constraints
- `moduleName` 仅允许小写字母、数字、下划线，且以字母开头
- 不要直接跳过 dry-run
- 出现冲突时必须要求确认覆盖

