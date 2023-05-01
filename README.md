# chatgpt-for-se-tasks

本项目包括ChatGPT调用、使用ChatGPT解析源代码中调用的方法与是否有异常处理、使用ChatGPT执行代码生成、使用ChatGPT执行代码翻译三个部分。

## 项目路径介绍

`ChatGPT`目录中存储了使用ChatGPT的相关代码。

`text_to_code`目录中存储了使用ChatGPT执行用自然语言生成java代码的代码与数据。

`code_to_code`目录中存储了使用ChatGPT执行代码翻译（c# -> java）的代码与数据。

`parse_references`目录中存储了使用ChatGPT解析源代码中调用的方法与是否有异常处理的代码与数据。

`utils`目录中为使用的工具包。

## 依赖关系

utils目录（utils包）为底层工具内容，提供包括文件读取、JSON处理等基本内容，被下列四个包使用。

ChatGPT目录（chat包）提供ChatGPT访问工具和代码生成功能的Prompts生成功能，被下列三个包使用。

parse_references目录（parse包）解析了ground truth代码中调用的方法与是否有异常处理功能，依赖上面两个包，其结果被下面两个包使用。

text_to_code目录（code_generation包）使用parse包生成的结果以及nl等信息执行代码生成任务，依赖上述三个包。

code_to_code目录（code_translation包）使用parse包生成的结果以及c#代码等信息执行代码翻译任务，依赖上述除code_generation外的三个包。

## 使用

请跳转到对应目录查看README文件，建议先看utils，再看chatGPT，再看parse_references，再看text_to_code和code_to_code。

## go.mod报错

（建议安装中文插件）  
在goland中双击shift，搜索`go 模块`设置，打开`启用go模块集成设置`即可。

## GitHub下载或Goland网络问题

在goland中双击shift，搜索`http`设置，找到`外观与行为`->`系统设置`->`HTTP代理`，选择`手动代理配置`，选择`HTTP`,编辑主机名为`127.0.0.1`，编辑端口号为`7890`。  
打开`Clash for Windows`代理。
