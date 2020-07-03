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

系统命令: ? , list, quit

数据操作命令: [object_name].[ops] + [field_name] + [value]

ops: get, set, delete, print

field_name: 当前对象的字段名(支持内嵌, 各级用‘.'分割)

value: eval字符串

注意： 

1. 首先需要使用use命令选定当前操作对象, 
2. delete命令只支持slice和map对象

---
### Examples

#### 0.
```
player.print
```
output: 
```
    Name=tom
    Id=654321
    Coin=123
    Friends len: 0
```

#### 1.
```
player.get Name
```
output: 
```
field Name:     string=nick
```

#### 2.
```
player.get Friends[0].Name
```
output: 
```
field Friends[0].Name:     string=tom
```

#### 3.
```
player.set Friends[0].Name "test"
```

#### 4. 
注意delete操作只能用于slice或map
```
player.delete Friends[0]
```
