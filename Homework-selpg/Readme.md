# Centos下使用Golang开发Selpg

---

## 一、实验前准备  

参考资料：  

- [开发Linux命令行实用教程](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)    
- [Go语言实现selpg](https://blog.csdn.net/kunailin/article/details/78262456)  
- [Golang之使用Flag和Pflag](https://o-my-chenjian.com/2017/09/20/Using-Flag-And-Pflag-With-Golang/)  

这些资料能让我们对命令行主要涉及内容能有比较好的理解，虽然我看完了之后，我也是一头雾水……  

接下来我们需要安装Pflag，这个很简单，只需要在gowork的pkg文件夹下使用命令：  
```
go get github.com/spf13/pflag
```
  
然后我们可以开始正式的编码了。  

---

## 二、编码设计   

我的代码主要分为四个大函数：  

- main函数
- ReceiveArgs函数
- CheckArgs函数
- HandleArgs函数  

### 1、main函数  

这个函数是最简单的，因为不需要什么实现，只需要接受参数，然后调用后面三个具体实现的函数。  

具体的代码如下：  
```Go  
func main() {
    args := new(selpg_args)
    ReceiveArgs(args)
    CheckArgs(args)
    HandleArgs(args)
}
```  

### 2、ReceiveArgs函数  

这个函数主要功能是将得到的参数进行分割，用到的函数是pflag这个库中的各种函数，其中，这里对每个没有赋值的参数，定义了缺省的值。  

具体的实现代码如下：  
```Go
func ReceiveArgs(args *selpg_args) {
	pflag.Usage = usage;
	//add error information
	pflag.IntVar(&(args.startPage), "s", -1, "start page")
	pflag.IntVar(&(args.endPage), "e", -1, "end page")
	pflag.IntVar(&(args.pageLen), "l", 72, "page len")
	pflag.StringVar(&(args.printDestination), "d", "", "print destionation")
	pflag.BoolVar(&(args.pageType), "f", false, "type of print")
	pflag.Parse()  
	//parse for input file names
	othersArg := pflag.Args()
	if len(othersArg) > 0 {
		args.inFile = othersArg[0]
	} else {
		args.inFile = ""
	}
} 
```  

### 3、CheckArgs函数  

这部分函数主要是完成对传入的参数的检查。  
需要检查的部分是：  

- 是否传入了开始页面和结束页面这两个参数；
- 开始页面和结束页面是否有效；
- 检查页面的长度是否有效。  

具体的实现代码如下：  
```Go
func CheckArgs(args *selpg_args) {
    if args.startPage == -1 || args.endPage == -1 {
        os.Stderr.Write([]byte("You should input --s --e at least\n"))
        pflag.Usage()
        os.Exit(0)
    }
    if args.startPage < 1 || args.startPage > (math.MaxInt32-1) {
        os.Stderr.Write([]byte("Invalid start page\n"))
        pflag.Usage()
        os.Exit(1)
    }
    if args.endPage < 1 || args.endPage > (math.MaxInt32-1) || args.endPage < args.startPage {
        os.Stderr.Write([]byte("Invalid end page\n"))
        pflag.Usage()
        os.Exit(2)
    }
    if args.pageLen < 1 || args.pageLen > (math.MaxInt32-1) {
        os.Stderr.Write([]byte("Invalid page length\n"))
        pflag.Usage()
        os.Exit(3)
    }
}
```  

### 4、HandleArgs函数  

这个函数的实现主要分成了以下几步：  

- 初始化Reader，当没有输入的文件路径的时候，我们选择标准输入。
- 初始化Writer，当没有输出的文件路径的时候，我们选择标准输出；当有路径的时候，使用管道输出。
- 处理其他的参数：如-d -l -f -l等。  

---

## 三、实验结果  

### 1、处理命令：selpg --s=1 --e=1 in.txt    

![image1](https://github.com/EmilyBlues/ServiceComputing/blob/master/Homework-selpg/image/1.png)

### 2、处理命令：selpg --s=1 --e=1 < in.txt    

![image2](https://github.com/EmilyBlues/ServiceComputing/blob/master/Homework-selpg/image/2.png)

### 3、处理命令：selpg --s=1 --e=1 in.txt > out.txt    

![image3](https://github.com/EmilyBlues/ServiceComputing/blob/master/Homework-selpg/image/3.png)

### 4、处理命令：selpg --s=1 --e=1 in.txt 2> error.txt    

![image4](https://github.com/EmilyBlues/ServiceComputing/blob/master/Homework-selpg/image/4.png)

### 5、处理命令：selpg --s=1 --e=1 --l=6 in.txt    

![image5](https://github.com/EmilyBlues/ServiceComputing/blob/master/Homework-selpg/image/5.png)

### 6、使用cat指令    

![image6](https://github.com/EmilyBlues/ServiceComputing/blob/master/Homework-selpg/image/6.png)

### 7、使用ps指令  

![image7](https://github.com/EmilyBlues/ServiceComputing/blob/master/Homework-selpg/image/7.png)
