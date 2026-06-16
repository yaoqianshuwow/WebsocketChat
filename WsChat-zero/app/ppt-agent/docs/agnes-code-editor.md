# Agnes Code Editor

这个项目新增了一套基于 `agnes-2.0-flash` 的代码编辑插件骨架，包含：

- 一个 Codex skill：指导何时用 Agnes 进行代码编辑
- 一个本地 MCP 服务：先生成候选修改，校验通过后才真正写入文件

插件位置：
- `C:\Users\28407\plugins\agnes-code-editor`

需要的环境变量：

```powershell
$env:AGNES_API_KEY="你的 Agnes API Key"
```

MCP 提供两个工具：

- `preview_code_edit`
  生成候选修改并做本地校验，只返回 diff，不写文件
- `apply_code_edit`
  生成候选修改并做本地校验，只有校验通过才写文件；失败会自动带着错误重试

支持的内建校验：

- `.go`：`gofmt -w`
- `.py`：`python -m py_compile`
- `.json`：`python -m json.tool`

也可以额外传入 `validation_command`：

- 命令里可使用 `{candidate}` 占位符
- 或读取环境变量 `CODE_EDIT_CANDIDATE`

示例思路：

1. 先调用 `preview_code_edit`
2. 查看 diff 与校验结果
3. 确认后再调用 `apply_code_edit`

