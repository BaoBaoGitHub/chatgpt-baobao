# 代码翻译


1. 使用`./dataset/test.java-cs.txt.cs`中的每一行code，制作chatGPT请求内容（将该文件中的c#代码翻译为java语言代码），具体查询语句如下
> Translate following c# code surrounded \`\`\` to java code.\`\`\`#{code}\`\`\`
2. 查询结果在`./dataset/test.java-cs.txt_response.json`中，在该json文件中
    - 其中`query`为查询输入
    - `flag`为是否成功返回java代码的标记，`code`为chatGPT响应中的代码（flag为true时code才有内容），`message`为chatGPT响应内容

# 使用

1. 在项目根路径下，运行`go run ./code_to_code/main.go`，或在`code_to_code`目录下运行`go run main.go`

2. 程序在运行大约需要一小时，若中途报错，请删除掉`code_to_code/dataset/test.java-cs.txt_`开头的所有文件，注意不包括`code_to_code/dataset/test.java-cs.txt`文件。
