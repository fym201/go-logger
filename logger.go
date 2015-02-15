package loggo

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	_VER string = "1.0.0"
)

type LEVEL int32

const DATEFORMAT = "2006-01-02"

type BYTE int64

const (
	_       = iota
	KB BYTE = 1 << (iota * 10)
	MB
	GB
	TB
)

const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	OFF
)

var _Prefix = map[LEVEL]string{DEBUG: "debug", INFO: "info", WARN: "warn", ERROR: "error", FATAL: "fatal"}

type _FILE struct {
	dir          string
	filename     string
	_suffix      int
	isCover      bool
	_date        *time.Time
	mu           *sync.RWMutex
	logfile      *os.File
	lg           *log.Logger
	maxFileSize  int64 //日志文件大小
	maxFileCount int32 //日志文件最大数量
	dailyRolling bool  //是否是按日期写文件
	rollingFile  bool  //是否是按文件大小的方式
}

type Logger struct {
	LogLevel        LEVEL //default:DEBUG,日志等级 ALL，DEBUG，INFO，WARN，ERROR，FATAL，OFF 级别由低到高
	ConsoleAppender bool  //default:true,是否要在控制台上输出
	logObj          *_FILE
}

//控制台log，不写文件
func NewConsoleLogger() (lg *Logger) {
	lg = &Logger{LogLevel: DEBUG, ConsoleAppender: true}
	return lg
}

//指定日志文件备份方式为文件大小的方式
//第一个参数为日志文件存放目录
//第二个参数为日志文件命名
//第三个参数为备份文件最大数量
//第四个参数为备份文件大小,单位为byte
func NewRollingFileLogger(fileDir, fileName string, maxNumber int32, maxSize BYTE) (lg *Logger, err error) {
	if !isExist(fileDir) {
		if err = os.MkdirAll(fileDir, 0x644); err != nil {
			return
		}
	}

	logObj := &_FILE{dir: fileDir,
		filename:     fileName,
		isCover:      false,
		mu:           new(sync.RWMutex),
		maxFileSize:  int64(maxSize),
		maxFileCount: maxNumber,
		dailyRolling: false,
		rollingFile:  true}

	for i := 1; i <= int(maxNumber); i++ {
		if isExist(fileDir + "/" + fileName + "." + strconv.Itoa(i)) {
			logObj._suffix = i
		} else {
			break
		}
	}
	if !logObj.isMustRename() {
		logObj.logfile, err = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
		if err != nil {
			return
		}
		logObj.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logObj.rename()
	}

	lg = &Logger{LogLevel: DEBUG,
		ConsoleAppender: true,
		logObj:          logObj}

	go logObj.fileMonitor()
	return
}

//指定日志文件备份方式为日期的方式
//第一个参数为日志文件存放目录
//第二个参数为日志文件命名
func NewRollingDailyLogger(fileDir, fileName string) (lg *Logger, err error) {

	if !isExist(fileDir) {
		if err = os.MkdirAll(fileDir, 0x644); err != nil {
			return
		}
	}

	t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	logObj := &_FILE{dir: fileDir,
		filename:     fileName,
		_date:        &t,
		isCover:      false,
		mu:           new(sync.RWMutex),
		dailyRolling: true,
		rollingFile:  false}

	if !logObj.isMustRename() {
		logObj.logfile, err = os.OpenFile(fileDir+"/"+fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0)
		if err != nil {
			return
		}
		logObj.lg = log.New(logObj.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		logObj.rename()
	}

	lg = &Logger{LogLevel: DEBUG,
		ConsoleAppender: true,
		logObj:          logObj}

	return
}

func (lg *Logger) console(s ...interface{}) {
	if lg.ConsoleAppender {
		_, file, line, _ := runtime.Caller(2)
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		log.Println(file+":"+strconv.Itoa(line), s)
	}
}

func catchError() {
	if err := recover(); err != nil {
		log.Println("err", err)
		//panic(err)
	}
}

func (lg *Logger) log(level LEVEL, v ...interface{}) {
	if lg.LogLevel <= level {
		if lg.logObj != nil {
			if lg.logObj.dailyRolling {
				lg.logObj.fileCheck()
			}
			defer catchError()
			lg.logObj.mu.RLock()
			defer lg.logObj.mu.RUnlock()

			lg.logObj.lg.Output(2, fmt.Sprintln(_Prefix[level], v))
		}

		lg.console(_Prefix[level], v)
	}
}

func (lg *Logger) logf(level LEVEL, format string, v ...interface{}) {
	if lg.LogLevel <= level {
		str := fmt.Sprintf(format, v...)
		if lg.logObj != nil {
			if lg.logObj.dailyRolling {
				lg.logObj.fileCheck()
			}
			defer catchError()
			lg.logObj.mu.RLock()
			defer lg.logObj.mu.RUnlock()

			lg.logObj.lg.Output(2, fmt.Sprintln(_Prefix[level], str))
		}

		lg.console(_Prefix[level], str)
	}
}

func (lg *Logger) Debug(v ...interface{}) {
	lg.log(DEBUG, v...)
}

func (lg *Logger) Debugf(format string, v ...interface{}) {
	lg.logf(DEBUG, format, v...)
}

func (lg *Logger) Info(v ...interface{}) {
	lg.log(INFO, v...)
}

func (lg *Logger) Infof(format string, v ...interface{}) {
	lg.logf(INFO, format, v...)
}

func (lg *Logger) Warn(v ...interface{}) {
	lg.log(WARN, v...)
}

func (lg *Logger) Warnf(format string, v ...interface{}) {
	lg.logf(WARN, format, v...)
}

func (lg *Logger) Error(v ...interface{}) {
	lg.log(ERROR, v...)
}

func (lg *Logger) Errorf(format string, v ...interface{}) {
	lg.logf(ERROR, format, v...)
}

func (lg *Logger) Fatal(v ...interface{}) {
	lg.log(FATAL, v...)
}

func (lg *Logger) Fatalf(format string, v ...interface{}) {
	lg.logf(FATAL, format, v...)
}

func (f *_FILE) isMustRename() bool {
	if f.dailyRolling {
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		if t.After(*f._date) {
			return true
		}
	} else {
		if f.maxFileCount > 1 {
			if fileSize(f.dir+"/"+f.filename) >= f.maxFileSize {
				return true
			}
		}
	}
	return false
}

func (f *_FILE) rename() {
	if f.dailyRolling {
		fn := f.dir + "/" + f.filename + "." + f._date.Format(DATEFORMAT)
		if !isExist(fn) && f.isMustRename() {
			if f.logfile != nil {
				f.logfile.Close()
			}
			err := os.Rename(f.dir+"/"+f.filename, fn)
			if err != nil {
				f.lg.Println("rename err", err.Error())
			}
			t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
			f._date = &t
			f.logfile, _ = os.Create(f.dir + "/" + f.filename)
			f.lg = log.New(f.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
		}
	} else {
		f.coverNextOne()
	}
}

func (f *_FILE) nextSuffix() int {
	return int(f._suffix%int(f.maxFileCount) + 1)
}

func (f *_FILE) coverNextOne() {
	f._suffix = f.nextSuffix()
	if f.logfile != nil {
		f.logfile.Close()
	}
	if isExist(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix))) {
		os.Remove(f.dir + "/" + f.filename + "." + strconv.Itoa(int(f._suffix)))
	}
	os.Rename(f.dir+"/"+f.filename, f.dir+"/"+f.filename+"."+strconv.Itoa(int(f._suffix)))
	f.logfile, _ = os.Create(f.dir + "/" + f.filename)
	f.lg = log.New(f.logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
}

func (f *_FILE) fileCheck() {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	if f.isMustRename() {
		f.mu.Lock()
		defer f.mu.Unlock()
		f.rename()
	}
}

func (f *_FILE) fileMonitor() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-timer.C:
			f.fileCheck()
		}
	}
}

func fileSize(file string) int64 {
	fmt.Println("fileSize", file)
	f, e := os.Stat(file)
	if e != nil {
		fmt.Println(e.Error())
		return 0
	}
	return f.Size()
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
