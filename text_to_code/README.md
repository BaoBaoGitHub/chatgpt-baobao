# 自然语言代码搜索

- 使用`./dataset/test_shuffled_with_path_and_id_concode.json`中的nl作为chatGPT的输入去查询代码
- 查询结果在`./dataset/response.json`中，在该json文件中
  - 其中`query`为查询输入
  - `flag`为是否成功返回java代码的标记，`code`为chatGPT响应中的代码（flag为true时code才有内容），`message`为chatGPT响应内容
