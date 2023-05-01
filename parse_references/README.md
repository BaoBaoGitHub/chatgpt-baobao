# 解析ground truth

解析ground truth中调用的方法和是否包括异常处理。

## dataset

包括code_generation和code_translation目录。

### code_generation目录

references.txt 为从concode原始数据集中解析出的ground truth。

references_api.json为请求ChatGPT references.txt文件中每行代码中方法调用的原始数据，references_api.txt为references_api.json中的code部分，如果没有调用api那么该行为空行。  
references_api.txt被复制粘贴到text_to_code/dataset/ref目录下，以供代码生成使用。


references_exception.json为请求ChatGPt references.txt文件中每行代码是否包含异常处理的相关数据，references_exception.txt为references_exception.json中的code部分。<br/>
references_exception.txt被复制粘贴到text_to_code/dataset/ref目录下，以供代码生成使用。

test_references.txt为references.txt的一小部分数据，拿来测试用的。

### code_translation目录

结构同code_generation目录。

### JSON格式说明

JSON文件的每一行是一个JSON对象，代表一次ChatGPT请求、响应与结果解析。

- query是发送给ChatGPT的请求内容

- message是ChatGPT的原始响应内容

- flag代表ChatGPT的响应内容中是否有预期结果

- code代表从ChatGPT的响应结果中解析出的预期结果

## parse目录

parse.go为解析的主体方法。

## main.go

TODO说明的部分为需要注意、修改的部分。

refPath为要解析的ground truth的路径，parseMode控制解析模式。