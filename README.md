# chatgpt-baobao

## 项目路径介绍

`chatGPT`目录中存储了使用chatGPT的相关代码

`test_to_code`目录中存储了使用chatGPT执行自然语言代码搜索的代码与数据。

`code_to_code`目录中存储了使用chatGPT执行代码翻译（c# -> java）的代码与数据。

`utils`目录中为使用的工具包。

## 使用

**注意：** 当前的工作路径应为项目根路径，关于`chatGPT, code_to_code, test_to_code`的具体信息，请到相应目录下查看`README.md`文件。

若想使用`chatGPT`，请执行`go run ./chatGPT/main.go`

若想使用`code_to_code`，请执行`go run ./code_to_code/main.go`

若想使用`test_to_code`，请执行`go run ./text_to_code/main.go`

# go.mod报错

在goland中双击shift，搜索`go 模块`设置，打开`启用go 模块集成设置`即可。