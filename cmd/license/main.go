package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"lisence/pkg/config"
	"lisence/pkg/iocgo"
	"lisence/pkg/router"
	"lisence/pkg/version"

	"lisence/pkg/libs/bufferpool"
	"lisence/pkg/libs/perf"
	"lisence/pkg/libs/rotatefd"

	log "github.com/sirupsen/logrus"
)

const (
	serviceName = `lisence`
	HttpTimeout = 30
)

/**
 * name		命令行参数名称，比如 -b, -help
 * value	默认值，未显式指定的参数，给出隐式的默认值
 * usage	提示信息，如果给出的参数不正确或者需要查看帮助 -help，那么会给出这里指定的字符串
 */
var cfgfile = flag.String("c", "config.yaml", "configuration file, default to config.yaml")

// 返回命令行参数是否未指定了version参数
var ver = flag.Bool("version", false, "Output version and exit")

func main() {
	// args parse,flag使用前，必须首先解析：
	flag.Parse()

	version.Service = serviceName

	// 没有指定version参数，退出
	if *ver {
		fmt.Println(version.Service, ": ", version.Version)
		return
	}

	// 使用viper进行配置文件解析
	cfg, err := config.NewTpaasConfig(*cfgfile)
	if err != nil {
		fmt.Printf(" load conf %s err:%v", *cfgfile, err)
		return
	}

	// 滚动日志
	rf := rotatefd.NewRotateFile(cfg.GetString("common.service_log"), 100*rotatefd.MiB)
	defer rf.Close()
	log.SetFormatter(&log.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {

			s := strings.Split(f.Function, "/")
			funcname := s[len(s)-1]
			buf := bufferpool.Get()
			flen := len(f.File)
			if flen > 20 {
				buf.WriteString("...")
				buf.WriteString(f.File[flen-20+3:])
				//filename = filename[flen-20:]
				//filename = "..." + filename[3:]
			}
			buf.WriteString(":")
			buf.WriteString(strconv.Itoa(f.Line))
			filename := buf.String()
			buf.Reset()
			bufferpool.Put(buf)
			return funcname, filename
		},
	})
	log.SetOutput(rf)
	log.SetReportCaller(true)

	//  log level
	log.SetLevel(log.WarnLevel)
	switch strings.ToLower(cfg.GetString("common.log_level")) {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	}

	/*
		//仿照Linux管理方式，存储所有的进程PID到文件中，后续可以通过读文件指定杀进程
		err = InitPidfile(cfg)
		if err != nil {
			fmt.Printf("initpidfile err:%v", err)
			return

		}
		defer QuitPidFile(cfg)
	*/
	//////////////////////////////////////
	//     service init
	//  init the dependency service
	//  初始化各种Pool
	err = InitDependencyService(cfg)
	if err != nil {
		log.Errorf("init dependency service err:%v", err)
		fmt.Printf("init dependency service err:%v", err)
		return
	}
	log.Infof("init dependency service ok")

	//  when service stopping, close the dependency service
	//  使用各个Pool注册时的Close方法将其资源进行释放
	defer CloseDependencyService()

	///////////////////////////////////////
	////  http service
	// perf service
	// 性能测试
	if cfg.IsSet("http.pprof_addr") {
		perf.Init(cfg.GetString("http.pprof_addr"))
		log.Infof("http pprof service init ok")
	}

	if !cfg.IsSet("http.http_addr") {
		log.Errorf("http address is not in config")
		fmt.Printf("http address is not in config")
		return
	}

	httptimeout := cfg.GetInt("http.http_timeout")
	if httptimeout == 0 {
		httptimeout = HttpTimeout
	}

	// http service init
	router.Api(). // singleTon api
			ConfigDefault().
			SetTimeout(time.Duration(httptimeout) * time.Second).
			SetLog(rf).
			InitRouter().
			Runapi(cfg.GetString("http.http_addr"))
	log.Infof("http service init ok")

	////////
	log.Infof("          ________                                                     ")
	log.Infof("       __/_/      |______   %s.%s is running                ", version.Service, version.Version)
	log.Infof("      / O O O O O O O O O ...........................................  ")
	log.Infof("                                                                       ")
	log.Infof("      %s", time.Now().String())
	log.Infof("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")

	////////
	// signal
	InitSignal()
}

func InitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	//log.Infof("Wait for signal.......")
	for {
		s := <-c
		log.Infof("service[%s] get a signal %s", version.Version, s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
			GracefulQuit()
			return
		// case syscall.SIGUSR2:
		// 	//todo: Define your signal processing functions
		// case syscall.SIGUSR1:
		// 	// todo: Define signal processing functions
		// 	config.ReloadGloableCfg()
		// 	iocgo.ReloadEngine(config.GloableCfg())
		// 	//return
		default:
			return
		}
	}
}

func GracefulQuit() {
	log.Infof("service make a graceful quit !!!!!!!!!!!!!!")
	router.Api().Shutdown() // close http service
	// close your service here

	//防止写日志的未读盘还在缓冲区中
	time.Sleep(1 * time.Second)
}

func InitDependencyService(cfg *config.TpaasConfig) error {
	return iocgo.LaunchEngine(cfg)
}

func CloseDependencyService() error {
	return iocgo.StopEngine()
}

/*
func InitPidfile(cfg *config.TpaasConfig) error {
	//pid file
	pidfile := ""
	if !cfg.IsSet("common.pid_file") {
		return nil
	} else {
		pidfile = cfg.GetString("common.pid_file")
	}
	contents, err := ioutil.ReadFile(pidfile)
	if err == nil {
		pid, err := strconv.Atoi(strings.TrimSpace(string(contents)))
		if err != nil {
			log.Errorf("Error reading proccess id from pidfile '%s': %s",
				pidfile, err)
			return errors.WithMessage(err, "reading proccess id from pidfile")
		}

		process, err := os.FindProcess(pid)

		// on Windows, err != nil if the process cannot be found
		if runtime.GOOS == "windows" {
			if err == nil {
				log.Errorf("Process %d is already running.", pid)
				return fmt.Errorf("already running")
			}
		} else if process != nil {
			// err is always nil on POSIX, so we have to send the process
			// a signal to check whether it exists
			if err = process.Signal(syscall.Signal(0)); err == nil {
				log.Errorf("Process %d is already running.", pid)
				return fmt.Errorf("already running")
			}
		}
	}
	if err = ioutil.WriteFile(pidfile, []byte(strconv.Itoa(os.Getpid())),
		0644); err != nil {

		log.Errorf("Unable to write pidfile '%s': %s", pidfile, err)
		return err
	}
	log.Infof("Wrote pid to pidfile '%s'", pidfile)
	return nil
}
*/
func QuitPidFile(cfg *config.TpaasConfig) {
	pidfile := ""
	if !cfg.IsSet("common.pid_file") {
		return
	} else {
		pidfile = cfg.GetString("common.pid_file")
	}

	if err := os.Remove(pidfile); err != nil {
		log.Errorf("Unable to remove pidfile '%s': %s", pidfile, err)
	}
	return
}
