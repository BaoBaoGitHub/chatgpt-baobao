# 代码翻译（C#->Java）

**注意：** 当前工作目录为此README.md文件所在目录！请先cd到`code_to_code`目录下！

## 项目介绍

1. 使用`./dataset/test.java-cs.txt.cs`中的每一行code，制作chatGPT请求内容（将该文件中的c#代码翻译为java语言代码），具体查询语句如下
2. 查询结果在`./dataset/test.java-cs.txt_response.json`中，在该json文件中
    - 其中`query`为查询输入
    - `flag`为是否成功返回java代码的标记，`code`为chatGPT响应中的代码（flag为true时code才有内容），`message`为chatGPT响应内容
    - 注意，这并不是一个符合语法的json文件，该文件中的每一行为一个json对象
3. 为了用于评估，将`./dataset/test.java-cs.txt.java`修改为txt格式，并拷贝到`./dataset/evaluator/references.txt`文件中，将`./dataset/test.java-cs.txt_response.json`中的code部分提取到`./dataset/evaluator/predictions.txt`中。
   - 其中 references.txt 为C#->Java的标准答案，predictions.txt 为ChatGPT生成结果

# 使用

**注意：** 项目已经运行成功，直接使用已生成的数据集即可，若无自定义需求，请不要执行下述代码。

1. 在项目根路径下，运行`go run ./code_to_code/main.go`，或在`code_to_code`目录下运行`go run main.go`

2. 程序运行大约需要一小时，若中途报错，请删除掉`code_to_code/dataset/test.java-cs.txt_`开头的所有文件，注意不包括`code_to_code/dataset/test.java-cs.txt.cs`和`code_to_code/dataset/test.java-cs.txt.java`文件。

# 评估

当前工作目录为此README.md文件所在目录！请先cd到`code_to_code`目录下！

`python ./dataset/evaluator/evaluator.py -ref ./dataset/evaluator/references.txt -pre ./dataset/evaluator/predictions.txt`

BLEU: 14.25 ; Acc: 0.0

cd 到CodeBLEU目录后，执行`python ./calc_code_bleu.py --refs ../../ref/references.txt --hyp ../../full_prompts/predictions.txt --lang java`

> ngram match: 0.05989724148292398, weighted ngram match: 0.1569376478257052, syntax_match: 0.5502183406113537, dataflow_match: 0.8235294117647058
> CodeBLEU score:  0.39764566042117216

计算CodeBLEU时若出错，请参考text_to_code中的README
