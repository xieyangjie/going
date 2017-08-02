# Golang 编码规范

## Gofmt
为保证代码格式一致，编码过程中或者提交代码之前，应该使用 gofmt 工具对代码进行格式化。

## 注释
Golang 支持 "//" 和 "/* */" 两种注释方式, 为了保持代码的统一和减少错误，统一采用双反斜线作为注释方式。

* 注释必须是完整的句子，以需要描述的内容作为开头，句点作为结尾；
* 注释应该出现在被注释内容的上方，切忌出现在被描述代码尾部及下方；
* 如果是针对 struct 的属性写注释内容，则可以忽略第二条规定，可以在代码尾部写注释，但是必须将属性的注释内容对齐;
* 注释与被描述内容之间不得出现空行，注释与其它不相关内容之间应该换行以作区分；
* “//” 与“注释内容”之间应该有一空格。

如下所示：

	// String returns the source text used to compile the regular expression.
	func (re *Regexp) String() string {
		return re.expr
	}
	
	// Compile parses a regular expression and returns, if successful, a Regexp
	// object that can be used to match against text.
	func Compile(str string) (regexp *Regexp, err error) {
	}

不强制要求对所有内容都书写注释，但是对于一些比较关键或者难以理解的点，应该配有注释内容。

## import
对 import 的包进行分组管理，一般分为四组：第一组为 Golang 标准库；第二组为第三方库；第三组为自有库；第四组为项目其它 package。
	
	import (
		//Golang 标准库
		"fmt"
		
		//第三访库
		"gopkg.in/mgo.v2/bson"
		
		//自己写的工具库
		"github.com/smartwalle/going/time"
		
		//项目相关
		"coffee/database"
		"coffee/app/models"
	)

关于 import, 有以下需要注意的几点：

* 拒绝使用单行引入，即使只引入一个包，也请务必采用 import () 的形式；
* 尽量避免使用匿名引入，所有引入的包，必须显示设置别名或者使用默认别名；
* 如果使用 Gofmt 工具对代码进行格式化，可以不按上面的描述对引入包进行分组；

## 命名
命名应该采用有意义的字符，尽量做到见名识意。

#### 包名
包名统一采用小写单词，不得使用下划线或者混合大小写。

#### 接口名
第一种形式，接口名统一以大写字母 "I" 开头，后续采用驼峰结构。
		
	type IHuman interface{
	}
	
第二种形式，可以不以大写字母 “I” 开头，结合是否可导出确定首字母是否需要大小写。

#### 结构体名
结构体名统一采用驼峰命名结构，不得出现下划线，结合是否可导出确定首字母是否需要大小写。

	type Human struct {
	}

#### 常量
常量统一采用大写字母，单词之间使用下划线进行分隔。如果是不可导出的常量，可在其名字前加上 "k\_" 作为前缀。可导出常量与不可导出常量应该分开声明，不得出现在同一常量声明块内，并且应该在其名字前加上 "K\_" 作为前缀。
	
	// 不可导出常量
	const (
		k_GENDER_MALE   = 1
		k_GEDNER_FEMALE = 2
	)
	
	// 可导出常量
	const (
		K_GENDER_MALE   = 1
		K_GENDER_FEMALE = 2
	)

#### 变量

##### 全局变量
采用驼峰结构命名，结合是否可导出确定首字母大小写，不得出现下划线。可导出变量与不可导出变量应该分开声明，不得出现在同一变量声明块内。

##### 形参
采用驼峰结构命名，首字母必须小写，不得出现下划线。

##### 局部变量
采用驼峰结构命名，首字母必须小写，不得出现下线线。

#### 函数（方法）
* 采用驼峰结构命名，结合是否可导出确定首字母大小写，不得出现下划线；
* 返回值必须命名，采用驼峰结构命名，首字母必须小写，不得出现下划线；
* 方法的接收者统一命名 this，接收者类型统一采用指针，特殊情况除外；
* 如果接收者是 map, slice 或者 chan，则不要用指针传递；

如下所示：

	// addition 这是一个包可见的函数，用于计算两个 int 类型的和
	func addition(param1, param2 int) (sum int) {
		sum = param1 + pram2
		return sum
	}
	
## 空格
利用空格（换行）来规范代码的格式。

### 对齐
代码中如果出现多行赋值语句，应该以等号为基准对齐:

	var name     = "Guest"
	var age      = 10
	var birthday = "2015.11.30"
	
如果是结构体属性，则应该以属性类型作为基准对齐:

	type Human struct {
		Name     string
		Age      int
		Birthday string
	}


### 代码块
各代码块之间只能有一行空格；

	package main
	
	import (
	)
	
	type IHuman interface {
	}
	
	func init() {
	}

### 顶格
最外层代码块应该顶格编写，不得留有任何空白，子级代码块相对于父级应该缩进四个空格或者一个 Tab；
	
	package main
	
	import (
		"fmt"
	)
	
	type IHuman interface {
		SayHi()
	}

### 圆括弧
* import、const、var 和方法接收者的左括弧("(")应该与相关关键字import、const、var 和 func 留有一个空格；
* 函数（方法）名与其后的左括弧不能留空格；
* 返回值与方法名之间必须保留一个空格；

例如：

	// import 与 "(" 之间留有一个空格.
	import (
		"fmt"
	)
	
	const (
	)
	
	var (
	)
	
	// 函数名与其后的"("挨在一起，不留空格.
	func addition(param1, param2 int) (sum int) {
	}

### 花括弧
用于代码块的花括弧必须与相关内容保留一个空格，花括弧内的内容必须顶行编写，块内的首行和尾行不能为空行；

	// interface 和 "{" 之间留有一个空格
	type IHuman interface {
	}
	
	func addition(param1, param2 int) (sum int) {
		// 顶行编写，不能留有空白
		sum = param1 + param2
		// 尾行不能为空行
		return sum
	}
	
### 逗号
逗号作为（变量声明、常量声明、形参、实参等）分隔符，其应该紧贴上一项，下一项与逗号之间应该保留一个空格。

	var a, b, c = 1, 2, 3
	
	var sum = addition(1, 2)
	
### 等号
等号两边的表达式应该与“＝”之间保留一个空格。

### 分号
任何单条语句结尾不得出现分号，分号只应该出现在 if、for 等条件语句中。条件语句中，“;”后面应该留有一个空格。

	for i := 0; i < 10; i++ {
	}
	
## 分隔
原则上应该是每一个独立的模块存在于独立的文件中，但是如果特殊需要将不同的模块整合在同一文件中，或者需要分隔一下代码块。可以使用 80 个 "/" 作为分隔符，但需要注意，"/" 应该与上一代码块保留一行空行，下一代码块与"/"之间不留空行。

	type IHuman interface {
	}

	////////////////////////////////////////////////////////////////////////////////
	type Human struct {
	}
	

