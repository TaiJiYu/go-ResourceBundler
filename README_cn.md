# 资源打包器
- en [English](README.md)
- zh_CN [简体中文](README_cn.md)
## 1. 用法
### 打包资源
1. 执行sh .\build.sh
2. 生成tool目录，在tool/resource下存放自己的资源文件，可以包含文件夹
3. 运行tool目录下的packer.exe自动打包resource下的文件
    - 也可以使用命令行执行packer.exe，提供如下可选参数：
      -   -n string
           产物名，示例：packer.exe -n my_package
      -  -o string
           产物目录，示例：packer.exe -o ./outcome
      -  -r string
           资源所在目录，示例：packer.exe -r ./resource
      -  -s string
           密码，必须小于等于16个字节，示例：packer.exe -s codesec
4. 打包后在outcome目录下生成fqm文件
5. 调用.fqm可用如下方法：
    ```go
    option := component.UnpackOption{
        FqmFilePath: "tool/outcome/my_resource.fqm", // 可提供绝对目录或相对目录
        SecretKey:   []byte("hjaslkdh"),  // 密码
	}
	err := api.UnpackerInit(o)
	if err != nil {
        panic(err)  
	}
	api.Show()  // 可查看fqm文件基础信息
    fmt.Println(api.Key("filepath/xx.txt"))

	api.Close() // 结束后调用
    ```

## 2. 文件结构
|序号|示例|含义|字节数|
|:--:|:--:|:--:|:--:|
|1|66 71 6D|文件头|3|
|2|00 00 00 00|对加密后所有字节的CRC-32校验|4|
|3|01 00|读取文件所需最小sdk版本1.0|2|
|4|01|是否加密,01加密,00未加密|1|
|5|07 E8|文件创建的年份|2|
|6|0B 08|文件夹创建的月份[1] 日期[1]|2|
|7|00 00 00 00|索引信息总字节数|4|
|8|00 00 00 00 00 00|索引名总字节数|6|
|9|00 00 00 00 00 00|数据总字节数|6|
|10*i|00 00 00 00|key_i开始的字节位置|4|
|10*i+1|00 00|key_i的字节数|2|
|10*i+2|00 00 00 00 00 00|key_i对应数据的开始字节位置|6|
|10*i+3|00 00 00 00 00 00|key_i对应数据的字节数|6|
||** ** ** ** **|key的排列||
||** ** ** ** **|数据排列||

