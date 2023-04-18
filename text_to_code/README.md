# 使用ChatGPT生成代码

**注意：** 当前工作目录为此README.md文件所在目录！请先cd到`text_to_code`目录下！

## 项目介绍
- 使用`./dataset/test_shuffled_with_path_and_id_concode.json`中的nl作为chatGPT的输入去查询代码
- 原始查询结果在`./dataset/test_shuffled_with_path_and_id_concode_response.json`中，在该json文件中
  - 其中`query`为查询输入
  - `flag`为是否成功返回java代码的标记，`code`为chatGPT响应中的代码（flag为true时code才有内容），`message`为chatGPT响应内容
  - 注意，这并不是一个符合语法的json文件，该文件中的每一行为一个json对象
- 为了用于评估，将`./dataset/test_shuffled_with_path_and_id_concode.json`中的nl与code部分提取到`./dataset/evaluator/answers.json`文件中，将`./dataset/test_shuffled_with_path_and_id_concode_response.json`中的code部分提取到`./dataset/evaluator/predictions.txt`中。
  - 其中answers.json为标准答案，predictions.txt为ChatGPT生成结果

## 使用

**注意：** 项目已经运行成功，直接使用已生成的数据集即可，若无自定义需求，请不要执行下述代码。

1. 在项目根路径下，运行`go run ./text_to_code/main.go`，或在`text_to_code`目录下运行`go run main.go`
2. 程序在运行大约需要一小时，若中途报错，请删除掉`text_to_code/dataset/test_shuffled_with_path_and_id_concode_`开头的所有文件，注意不包括`code_to_code/dataset/test_shuffled_with_path_and_id_concode.json`文件。

## 评估

以该`README.md`所在目录为工作目录。

1. 计算ACC和BLEU：执行 `python ./dataset/evaluator/evaluator.py -a=./dataset/evaluator/answers.json -p=./dataset/evaluator/predictions.txt`  
   > BLEU: 0.77, EM: 0.0
2. 计算CodeBLEU：执行`cd ./dataset/evaluator/CodeBLEU`,然后 `python ./calc_code_bleu.py --refs ../answers.json --hyp ../predictions.txt --lang java`  
   > WARNING: There is no reference data-flows extracted from the whole corpus, and the data-flow match score degenerates to 0. Please consider ignoring this score.
   > ngram match: 0.031972361577153885, weighted ngram match: 0.035803209671719745, syntax_match: 0.10433854907539118, dataflow_match: 0
   > CodeBLEU score:  0.043028530081066205

## 问题及解决

1. 如果计算CodeBLEU库时出错，你需要`pip install tree-sitter`
2. 如果tree-sitter库无法安装，
   1. 如果你是Windows操作系统，请在当前工作目录打开Git Bash，确保可以执行脚本
   2. `cd ./dataset/evaluator/CodeBLEU/my_parser`
   3. `sh build.sh`, 如果执行脚本时报错：distutils.errors.DistutilsPlatformError: Microsoft Visual C++ 14.0 or greater is required. Get it with "Microsoft C++ Build Tools": https://visualstudio.microsoft.com/visual-cpp-build-tools/ ，去官网下载构建工具即可
   4. `cd ../../../../`
3. 计算CodeBLEU时若报错OSError: \[WinError 193\] %1 不是有效的 Win32 应用程序，请参考<https://github.com/microsoft/CodeXGLUE/issues/116>


