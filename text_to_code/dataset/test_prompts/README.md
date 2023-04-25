# test_prompts的评估结果

## 执行评估

在`text_to_code/dataset/evaluator`目录下

`python evaluator.py --ref ../ref/references.txt --pre ../test_prompts/predictions.txt`
   
> BLEU: 1.72 ; Acc: 0.0
 
在`text_to_code/dataset/evaluator/CodeBLEU`目录下

`python calc_code_bleu.py --refs "../../ref/references.txt" --hyp "../../test_prompts/predictions.txt" --lang java`

> ngram match: 0.016694208683044658, weighted ngram match: 0.025459404610297758, syntax_match: 0.39537507050197407, dataflow_match: 0.49206349206349204  
> CodeBLEU score:  0.23239804396470215

# 观察

ngram match和weighted ngram match分数很低，syntax_match和dataflow_match尤其是dataflow_match分数很高

在计算CodeBLEU时候有参数`--params 0.25,0.25,0.25,0.25(default)`，如果调整参数比重，执行评估：`python calc_code_bleu.py --refs "../../ref/references.txt" --hyp "../../test_prompts/predictions.txt" --lang java --params 0.1,0.1,0.3,0.5`

> ngram match: 0.016694208683044658, weighted ngram match: 0.025459404610297758, syntax_match: 0.39537507050197407, dataflow_match: 0.4973544973544973<br/>
> CodeBLEU score:  0.3715051311571751
