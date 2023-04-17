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

执行 `python ./dataset/evaluator/evaluator.py -a=./dataset/evaluator/answers.json -p=./dataset/evaluator/predictions.txt`

BLEU: 0.77, EM: 0.0（奇怪的结果，改进在语雀中描述了）


