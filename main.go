package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/blueicesir/clash/config"
	C "github.com/blueicesir/clash/constant"
	"github.com/blueicesir/clash/hub"

	log "github.com/sirupsen/logrus"

	// add by BlueICE
	// "strings"
	// "golang.org/x/sys/windows/svc"
	// "github.com/kardianos/service"

	// add by BlueICE Windows systray
	"io/ioutil"
	"github.com/getlantern/systray"
	"os/exec"

	"encoding/json"
	"net/http"
	"bytes"

)

const svcName="Clash for BlueICE"

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
	mEdit:=systray.AddMenuItem("Edit config.yml","Edit Clash Config")
	mReloadConfig:=systray.AddMenuItem("Reload config.yml","Reload Clash Config")
	systray.AddSeparator()
	mQuit:=systray.AddMenuItem("Exit","Exit then Clash")
	go ClashMain()

	// 退出按钮只能点击一次
	go func(){
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	// 事件循环，如果是编辑配置文件菜单，可以多次触发。
	go func(){
		for {
			select {
				case <-mEdit.ClickedCh:
					home:=getUserDir()
					config_path:=filepath.Join(home,".config\\clash\\config.yml")
					c:=exec.Command("D:\\Tools\\EditPlus\\editplus.exe",config_path)
					c.Start()

				case <-mReloadConfig.ClickedCh:
					go ReloadConfig()
			}

		}
	}()
}

// 获取用户目录，未来可以适配Unix以及Linux系统
func getUserDir() string {
	var home string
	if "windows"==runtime.GOOS {
		home=os.Getenv("HOMEDRIVE")+os.Getenv("HOMEPATH")
		if home == ""{
			home=os.Getenv("USERPROFILE")
		}
	}
	return home
}

// 重新加载配置的Tray菜单
func ReloadConfig(){
	home:=getUserDir()
	config_path:=filepath.Join(home,".config\\clash\\config.yml")
	url:="http://127.0.0.1:9090/configs"
	params:=map[string]string{"path":config_path,"force":"true"}
	jsonBody,_:=json.Marshal(params)
	request,_:=http.NewRequest("PUT",url,bytes.NewBuffer(jsonBody))
	request.Header.Add("Content-Type","application/json;charset=utf-8")
	client:=&http.Client{}
	resp,_:=client.Do(request)
	defer resp.Body.Close()
	fmt.Printf("ReloadConfig Response Code:%d",resp.StatusCode)
	body,_:=ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
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
