# TDesign 交互与组件统一标准

> 依据来源  
> 1. 官方模板：`Tencent/tdesign-vue-next-starter`  
> 2. 官方组件站：`tdesign.tencent.com/vue-next`  
> 3. 模板样式变量（参考）：`src/style/variables.less` 中动画节奏与动效曲线定义

## A. 组件统一使用规则

1. 表单输入必须使用 TDesign 组件：`t-input / t-select / t-form / t-form-item`。  
2. 数据展示必须使用：`t-table / t-tag / t-pagination`。  
3. 弹窗统一使用：`t-dialog`，不再使用自定义弹层。  
4. 操作按钮统一使用：`t-button / t-link`。  
5. 标题与文本优先使用：`t-typography-title` 等官方排版组件。

## B. 动效与节奏规则

1. 页面切换使用轻量过渡（`0.24s`，ease 曲线），避免夸张位移。  
2. 卡片 hover 仅允许轻微阴影和 1-2px 位移。  
3. 弹窗动画由 `t-dialog` 内置动效统一控制，不自定义冲突动画。  
4. 禁止引入重动画背景、粒子效果、复杂全局动画。

## C. 布局规则（后台）

1. 标准结构：深色侧栏 + 浅色头部 + 内容区。  
2. 内容区卡片统一 `6px` 圆角、`1px` 边框、标准间距（16/24）。  
3. 顶部品牌区、菜单区、头部栏必须对齐（同宽/同内边距体系）。  
4. 侧栏必须支持收起/展开。

## D. 验收清单

1. 视觉风格是否贴近 `tdesign-vue-next-starter`。  
2. 所有交互组件是否均为 TDesign 官方组件。  
3. 页面动画是否克制且一致（含弹窗）。  
4. 侧栏与头部对齐、收起展开是否正常。

