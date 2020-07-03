## expr-CLI
安装说明:
```
git clone https://github.com/sukrfu/expr-CLI.git
```
在项目目录下,
```
go run ./
```

---

使用说明：

命令: [ops] + [field_name] + [value]

ops: get, set, delete, use, print, list, quit 

field_name: 当前对象的字段名(支持内嵌, 各级用‘.'分割)

value: eval字符串

注意： 

1. 首先需要使用use命令选定当前操作对象, 
2. delete命令只支持slice和map对象