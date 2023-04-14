# 自然语言代码搜索

- 使用`./dataset/test_shuffled_with_path_and_id_concode.json`中的nl作为chatGPT的输入去查询代码
- 查询结果在`./dataset/test_shuffled_with_path_and_id_concode_response.json`中，在该json文件中
  - 其中`query`为查询输入
  - `flag`为是否成功返回java代码的标记，`code`为chatGPT响应中的代码（flag为true时code才有内容），`message`为chatGPT响应内容
  - 注意，这并不是一个符合语法的json文件，该文件中的每一行为一个json对象

# 使用

1. 在项目根路径下，运行`go run ./text_to_code/main.go`，或在`text_to_code`目录下运行`go run main.go`

2. 程序在运行大约需要一小时，若中途报错，请删除掉`text_to_code/dataset/test_shuffled_with_path_and_id_concode_`开头的所有文件，注意不包括`code_to_code/dataset/test_shuffled_with_path_and_id_concode.json`文件。
