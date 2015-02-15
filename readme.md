loggo [![Build Status](https://drone.io/github.com/fym201/loggo/status.png)](https://drone.io/github.com/fym201/loggo/latest) 
=======
loggo 是golang 的日志库 ，是从go-logger上修改而来。 
用法类似java日志工具包log4j 

打印日志有5个方法 Debug(Debugf)，Info(Infof)，Warn(Warnf), Error(Errorf),Fatal(Fatalf)  日志级别由低到高 

设置日志级别的方法为：`logger.LogLevel = loggo.WARN` 
则：logger.Debug(....),logger.Info(...) 日志不会打出
而`logger.Warn(...)`,`logger.Error(...)`,`logger.Fatal(...)`日志会打出。 


设置日志级别的参数有7个，分别为：`ALL,DEBUG,INFO, WARN,ERROR,FATAL,OFF` 
其中`ALL`表示所有调用打印日志的方法都会打出，而`OFF`则表示都不会打出。


日志文件切割有两种类型：1为按日期切分。2为按日志大小切分。
按日期切分时：每天一个备份日志文件，后缀为 `.yyyy-MM-dd`
过0点是生成前一天备份文件
	logger, err := loggo.NewRollingDailyLogger("d:/log_test/rolling_daily", "test.log")


按大小切分是需要3个参数，1为文件大小，2为单位，3为文件数量
文件增长到指定限值时，生成备份文件，结尾为依次递增的自然数。
文件数量增长到指定限制时，新生成的日志文件将覆盖前面生成的同名的备份日志文件。
	logger, err := loggo.NewRollingFileLogger("d:/log_test/rolling_file", "test.log", 10, 5*loggo.KB)

###示例：

//打印日志 
```
func log(lg *loggo.Logger, i int) {
	lg.Debug("Debug>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	lg.Debugf("Debug>>>>>>>>>>>>>>>>>>>>>>%s"+strconv.Itoa(i), "format test")
	lg.Info("Info>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	lg.Warn("Warn>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	lg.Error("Error>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	lg.Fatal("Fatal>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
}
```

//单纯的控制台log，不写文件
```
func consoleLog() {
	lg := loggo.NewConsoleLogger()

	for i := 10000; i > 0; i-- {
		go log(lg, i)
		time.Sleep(1000 * time.Millisecond)
	}
}
```

//指定日志文件备份方式为文件大小的方式 
```
func rollingFileLog() {
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	//第三个参数为备份文件最大数量
	//第四个参数为备份文件大小,单位为byte
	lg, err := loggo.NewRollingFileLogger("d:/log_test/rolling_file", "test.log", 10, 5*loggo.KB)
	if err != nil {
		panic(err)
	}
	lg.LogLevel = loggo.DEBUG //日志级别，默认为logger.DEBUG
	for i := 10000; i > 0; i-- {
		go log(lg, i)
		time.Sleep(1000 * time.Millisecond)
	}
}
```


//指定日志文件备份方式为日期的方式 
```
func rollingDailyLog() {
	//第一个参数为日志文件存放目录
	//第二个参数为日志文件命名
	lg, err := loggo.NewRollingDailyLogger("d:/log_test/rolling_daily", "test.log")
	if err != nil {
		panic(err)
	}
	for i := 10000; i > 0; i-- {
		go log(lg, i)
		time.Sleep(100 * time.Millisecond)
	}
}
```
