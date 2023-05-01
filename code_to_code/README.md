# 代码翻译（C#->Java）

使用chat包下的prompts.go文件中定义的常量控制模式，使用parse包处理所得的dataset/ref/references_api.txt与dataset/ref/references_exception.txt文件内容作为prompts的一部分，请求ChatGPT执行代码翻译任务。

## 目录介绍

### code_translation

该目录下有代码翻译的主要执行的方法。

### dataset

#### ref目录

- test.java-cs.txt.cs为c#文件，test.java-cs.txt.java为java文件
- references_api.txt从parse_references目录下的dataset/code_translation中拷贝得到
- references_exception.txt从parse_references目录下的dataset/code_translation中拷贝得到
- test.cs,test.java,test_references_apit.txt和test_references_exception.txt为上述文件的测试版

### evaluator目录

evaluator.py用于计算BLEU和ACC，CodeBLEU下的calc_code_bleu.py用于计算CodeBLEU。


### prompts目录

不同的prompts目录中有不同的prompts模式生成的结果文件，test为使用test数据生成的结果，round0为完整数据生成的结果。

其中references.txt为ground truth(ref目录下的test.java-cs.txt.java去掉java后缀)，response json为请求与响应的原始数据，predictions.txt为response json中的code部分或message部分。
在json文件中，
- 每行为一个JSON对象，代表一次ChatGPT请求与响应。
- `query`为查询输入、
- `message`为chatGPT原始响应内容
- `flag`为是否成功返回java代码的标记（message中包括```符号就认为返回了代码片段）
- `code`为chatGPT响应中的代码（flag为true时code才有内容）

## 使用

**注意：** 项目已经运行成功，直接使用已生成的数据集即可，若无自定义需求，请不要执行下述代码。

1. 查看并修改main.go文件，其中有两个TODO，一个是控制模式（task、detailed、guided），一个是是否使用测试数据集（前100条）。
2. 程序执行结果输出到dataset下的对应prompt目录中，若重复执行请清空掉prompt下的文件或将其移动到test1、test2、round1或round2目录中。
3. 使用测试数据集运行很快，但使用完整数据集运行会很慢，请耐心等待，我使用了accesstoken池与健壮的错误处理，耐心等待会有实验结果~
4. 若中途停止了运行，请删除掉ref目录与对应prompts目录下的中间文件（文件名中有一堆数字的），重新运行。
5. 执行main.go文件即可运行。

## 评估

1. 计算ACC和BLEU。

   在dataset/evaluator目录下，执行 `python evaluator.py --ref "references.txt路径" --pre "predictions.txt路径"`

2. 计算CodeBLEU。
   在dataset/evaluator/CodeBLEU目录下，执行 `python calc_code_bleu.py --refs "references.txt路径"  --hyp "predictions.txt路径" --lang java`

## 问题及解决

1. 如果计算CodeBLEU库时出错，需要`pip install tree-sitter`
2. 如果tree-sitter库无法安装，
   1. 如果你是Windows操作系统，请在当前工作目录打开Git Bash，确保可以执行脚本
   2. `cd dataset/evaluator/CodeBLEU/my_parser`
   3. `sh build.sh`, 如果执行脚本时报错：distutils.errors.DistutilsPlatformError: Microsoft Visual C++ 14.0 or greater is required. Get it with "Microsoft C++ Build Tools": https://visualstudio.microsoft.com/visual-cpp-build-tools/ ，去官网下载构建工具即可
   4. `cd到CodeBLEU目录下继续评估`
3. 计算CodeBLEU时若报错OSError: \[WinError 193\] %1 不是有效的 Win32 应用程序，请参考<https://github.com/microsoft/CodeXGLUE/issues/116>
