package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/Dreamacro/clash/config"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/hub"

	log "github.com/sirupsen/logrus"

	// add by BlueICE
	// "strings"
	// "golang.org/x/sys/windows/svc"
	// "github.com/kardianos/service"

	// add by BlueICE Windows systray
	"io/ioutil"
	"github.com/getlantern/systray"
)

const svcName="Clash for BlueICE"

type program struct{}

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) run() {
	// logger.Infof("%s running %v.",svcName,service.Platform())
	go ClashMain()
}

func (p *program) Stop(s service.Service) error {
	// logger.Infof("%s Stopping!",svcName)
	return nil
}



var (
	version bool
	homedir string
)

func init() {
	flag.StringVar(&homedir, "d", "", "set configuration directory")
	flag.BoolVar(&version, "v", false, "show current version of clash")
	flag.Parse()
}

func main(){
	systray.Run(onReady,onExit)
}

func onReady(){
	systray.SetIcon(getIcon("rsc/clash.ico"))
	systray.SetTitle("Clash for BlueICE")
	systray.SetTooltip("Clash for BlueICE")
	go ClashMain()
}

func getIcon(s string) [] byte {
	b,err:=ioutil.ReadFile(s)
	if err!=nil{
		fmt.Print(err)
	}
	return b
}

func onExit(){

}

// func main(){
// 	svcConfig := &service.Config{
// 		Name:        svcName,
// 		DisplayName: svcName,
// 		Description: svcName,
// 	}

// 	prg := &program{}
// 	s, err := service.New(prg, svcConfig)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if len(os.Args) > 1 {
// 		var err error
// 		verb := os.Args[1]
// 		switch verb {
// 		case "install":
// 			err = s.Install()
// 			if err != nil {
// 				log.Errorf("Failed to install: %s\n", err)
// 				return
// 			}
// 			log.Infof("Service \"%s\" installed.\n", svcConfig.DisplayName)
// 		case "uninstall":
// 			err = s.Uninstall()
// 			if err != nil {
// 				log.Errorf("Failed to Uninstall: %s\n", err)
// 				return
// 			}
// 			log.Infof("Service \"%s\" Uninstall.\n", svcConfig.DisplayName)
// 		}
// 		return
// 	} else {
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		err = s.Run()
// 		if err != nil {
// 			log.Error(err)
// 		}
// 		log.Infof("Service \"%s\" started.\n", svcConfig.DisplayName)
// 	}
// }


func ClashMain() {
	if version {
		fmt.Printf("Clash %s %s %s %s\n", C.Version, runtime.GOOS, runtime.GOARCH, C.BuildTime)
		return
	}

	// enable tls 1.3 and remove when go 1.13
	os.Setenv("GODEBUG", os.Getenv("GODEBUG")+",tls13=1")

	if homedir != "" {
		if !filepath.IsAbs(homedir) {
			currentDir, _ := os.Getwd()
			homedir = filepath.Join(currentDir, homedir)
		}
		C.SetHomeDir(homedir)
	}

	fmt.Println(fmt.Sprintf("clash config load by %s ",C.Path.HomeDir()))
	if err := config.Init(C.Path.HomeDir()); err != nil {
		log.Fatalf("Initial configuration directory error: %s", err.Error())
	}

	if err := hub.Parse(); err != nil {
		log.Fatalf("Parse config error: %s", err.Error())
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
