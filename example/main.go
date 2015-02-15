package main

import (
	"github.com/fym201/loggo"
	"runtime"
	"strconv"
	"time"

	"fmt"
)

func log(lg *loggo.Logger, i int) {
	lg.Debug("Debug>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	lg.Debugf("Debug>>>>>>>>>>>>>>>>>>>>>>%s"+strconv.Itoa(i), "format test")
	lg.Info("Info>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	lg.Warn("Warn>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	lg.Error("Error>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	lg.Fatal("Fatal>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
}

func consoleLog() {
	//单纯的控制台log，不写文件
	lg := loggo.NewConsoleLogger()

	for i := 10000; i > 0; i-- {
		go log(lg, i)
		time.Sleep(1000 * time.Millisecond)
	}
}

func rollingFileLog() {
	//指定日志文件备份方式为文件大小的方式
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

func rollingDailyLog() {
	//指定日志文件备份方式为日期的方式
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

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("start")
	go consoleLog()
	go rollingFileLog()
	go rollingDailyLog()

	wait := make(chan int)
	<-wait
}
