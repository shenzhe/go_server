package comm

import (
    "log"
    "os"
    "config"
    "fmt"
    "time"
    "strconv"
)

const (
    DEBUG = "debug"
    LOG = "log"
    ERROR = "error"
)

var _loggers = make(map[string]*log.Logger , 0)

func getDate(sep string) string {

    year, month, day := time.Now().Date()

    return strconv.Itoa(year) + sep + strconv.Itoa(int(month)) + sep + strconv.Itoa(day)
}

func checkDir(dir string, mode os.FileMode) {
    err := os.Mkdir(dir, mode)
    if nil != err {
        if os.IsNotExist(err) {
            CheckError(err, "mkdir " + dir + " error") 
        }
    }
}

func GetLogger(filename string) *log.Logger {
    logger, ok := _loggers[filename] 
    if !ok {
        dir := config.LogPath + getDate("_")
        checkDir(dir, 0777)
        logfile, err := os.OpenFile(dir + "/"  + filename + ".log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0)
    if err != nil {
        fmt.Printf("%s, %s\r\n", filename, err.Error())
        os.Exit(-1)
    }
    logger = log.New(logfile, "", log.Ldate|log.Ltime|log.Llongfile)
        _loggers[filename] = logger
    }

    return logger
}


func Debug(format string, v ...interface{}) {
    GetLogger(DEBUG).Printf(format, v...)
}

func Log(format string, v ...interface{}) {
    GetLogger(LOG).Printf(format, v...)
}
