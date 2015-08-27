## 手机号码库
根据手机号的前7位,查出对应的省份,城市,邮编,区号,手机卡的类型
想法来源: [lovedboy/phone](https://github.com/lovedboy/phone)

## 安装

```bash
go get -u github.com/deloz/phone
```

## 例子

### 引入
```go
 phoneInfo, err := phone.Find("1888888")
	if err != nil {
 		fmt.Println(err)
 } else {
 		fmt.Println(phoneInfo)
 }
```

### 测试
```go
go test -v
```

### 输出
```bash
&{Phone:1888888 Province:北京 City:北京 ZipCode:100000 AreaCode:0 PhoneType:移动 PhoneRecordCount:307990}
```

## 支持号段
13\* , 15\* , 18\* , 14[5,7] , 17[0,6,7,8]

## phone.dat文件格式

```

        | 4 bytes |                     <- phone.dat 版本号
        ------------
        | 4 bytes |                     <-  第一个索引的偏移
        -----------------------
        |  offset - 8            |      <-  记录区
        -----------------------
        |  index                 |      <-  索引区
        -----------------------

```

* `头部` 头部为`8`个字节，版本号为`4`个字节，第一个索引的偏移为`4`个字节。
* `记录区` 中每条记录的格式为`<省份>|<城市>|<邮编>|<长途区号>\0`。 每条记录以`\0`结束。
* `索引区` 中每条记录的格式为`<手机号前七位(长4字节)><记录区的偏移(长4字节)><卡类型(长1字节)>`，每个索引的长度为`9`个字节。

### 解析步骤:

 * 解析头部`8`个字节，得到索引区的第一条索引的偏移。
 * 在索引区用二分查找得出`手机号在记录区的记录偏移`。
 * 在记录区从上一步得到的`记录偏移处取数据`，直到遇到`\0`。


### 定义的卡类型为:

* 1 移动
* 2 联通
* 3 电信
* 4 电信虚拟运营商
* 5 联通虚拟运营商
* 6 移动虚拟运营商

## Contributing

1. Fork it ( https://github.com/deloz/phone/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request

## LICENSE

The MIT License (MIT) Copyright (c) 2015 Deloz
