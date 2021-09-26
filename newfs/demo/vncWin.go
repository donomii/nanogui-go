package demo

import (

	nanogui "../.."

	"fmt"
	
	"github.com/shibukawa/nanovgo"
	"github.com/donomii/goof"
	"github.com/donomii/hashare"
	"github.com/kardianos/osext"

	"log"
	
	"time"
	"os/exec"
	"os/user"
	"os"
	"io/ioutil"


)

func VncAuth(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {
	return AuthWin(app, screen, "Login to Vnc", "VncAuth", [][]string{
		[]string{"Server :", "vnc-server", "localhost"},
		[]string{"Password :", "vnc-password", "aaaaaa"},
		[]string{"Port :", "vnc-port", "5900"},
	})
}









var repoDir string
var controlDir string
var debugOutput, traceOutput bool
var useGui bool = true

func loadUserConfig() []string {

	files := goof.Ls(repoDir)
	var out []string
	for i := 0; i < len(files); i++ {
		out = append(out, repoDir+files[i])
	}
	return out
}

func loadRepoDetails(path string) hashare.ClientConfig {

	out, _ := hashare.LoadClientConfig(path)
	return out
}

func mountRepo(config string, debug, trace bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unable to mount %v because %v", config, r)
		}
	}()
	log.Println("Starting ", debug)
	exePath := "vort-fuse.exe"
	if !goof.Exists(exePath) {
		var err error
		exePath, err = exec.LookPath("vort-fuse")
		if err != nil {
			dir, _ := osext.ExecutableFolder()
			exePath = dir + "/vort-fuse"
		}
	}
	cmd := []string{}
	conn := loadRepoDetails(config)
	optList := []string{exePath, "--config", config}
	if debug {
		optList = append(optList, "--debug")
	}
	if trace {
		optList = append(optList, "--trace")
	}
	MountPoint := conn.MountPoint
	if MountPoint == "" {
		MountPoint = nextDrive()
		cmd = []string{exePath, "--config", config, "--mount", MountPoint, "--trace"}
	} else {
		cmd = []string{exePath, "--config", config, "--trace"}
	}
	if debugOutput {
		fmt.Println("Starting mount with command: ", cmd)
		fmt.Println("Mounting to", MountPoint)
	}
	go func() {
		time.Sleep(5 * time.Second)
		//TODO have vort-fuse write a connected status then pop the window asap
		goof.QC([]string{`explorer`, MountPoint + "\\"})
	}()
	goof.QCI(cmd)
	time.Sleep(1 * time.Second)

}



func mountLocal( debug, trace bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unable to mount local share because %v", r)
		}
	}()
	log.Println("Starting ", debug)
	exePath := "vort-pclient.exe"
	if !goof.Exists(exePath) {
		var err error
		exePath, err = exec.LookPath("vort-pclient.exe")
		if err != nil {
			dir, _ := osext.ExecutableFolder()
			exePath = dir + "/vort-pclient.exe"
		}
	}

	cmd := []string{exePath}
	goof.QCI(cmd)
	time.Sleep(1 * time.Second)


}


var nextDriveId = 0

func nextDrive() string {

	drives := []string{"z:", "y:", "x:", "w:", "v:", "u:", "t:"}
	if nextDriveId+1 > len(drives) {
		panic(fmt.Sprintf("Ran out of available drive letters, tried %v", nextDriveId))
	}
	driveStr := drives[nextDriveId]
	nextDriveId = nextDriveId + 1
	return driveStr
}

func mountAllRepositories(debug, trace bool) {
	repositories := loadUserConfig()
	for _, v := range repositories {
		if debugOutput {
			fmt.Println("Opening config ", v)
		}
		go mountRepo(v, debug, trace)

	}

}





func PClientWin(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {

	user, _ := user.Current()
	hDir := user.HomeDir
	repoDir = hDir + "/.vort/repositories/"
	controlDir = hDir + "/.vort/control/"
	os.Remove(controlDir + "shutdown")
	if !goof.Exists(repoDir) {
		os.MkdirAll(repoDir, 0700)
	}
	if !goof.Exists(controlDir) {
		os.MkdirAll(controlDir, 0700)
	}

	window := nanogui.NewWindow(screen, "PClient")
	//window.SetFixedSize(100, 100)
	//window.SetSize(200, 200)

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = "VncWin"

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1
	window.SetPosition(0, 0)




	window.SetLayout(nanogui.NewGroupLayout())


	repositories := loadUserConfig()
	for _, k := range repositories {
		file := k
		//s := gitremind.Repos[file]

		//node := vv.SubNodes[i]
		repoDetail := loadRepoDetails(file)
		name := repoDetail.DisplayName
	

	nanogui.NewLabel(window, "Status:").SetFont("sans-bold")
	status:= nanogui.NewLabel(window, "Disconnected")//.SetFont("sans-bold")
	b2 := nanogui.NewButton(window, "Start "+name)
	b2.SetBackgroundColor(nanovgo.RGBA(0, 0, 255, 25))
	b2.SetIcon(nanogui.IconRocket)
	repoPath:= file
	b2.SetCallback(func() {
		fmt.Println("pushed!")
		status.SetCaption("Connected")
		os.Remove(controlDir + "shutdown")

					go mountRepo(repoPath, true, false)
		
	})

	b3 := nanogui.NewButton(window, "Stop"+name)
	b3.SetBackgroundColor(nanovgo.RGBA(0, 0, 255, 25))
	b3.SetIcon(nanogui.IconRocket)
	b3.SetCallback(func() {
		status.SetCaption("Disconnected")
		ioutil.WriteFile(controlDir+"shutdown", []byte(" "), 0600)
	})
	}
	nanogui.NewResize(window, window)
	img := nanogui.NewImageView(window)
	img.SetPolicy(nanogui.ImageSizePolicyExpand)
	//img.SetFixedSize(800, 600)
	//img.SetSize(800, 600)
	nanogui.NewResize(window, window)


	return window

}
