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
"net/http"

)

func NFSAuth(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {
	return AuthWin(app, screen, "Login to NFS", "NFSAuth", [][]string{
		[]string{"Server :", "nfs-server", "http://192.168.178.22:8000/"},
		[]string{"Token :", "nfs-token", "12394348765432782912349283435"},
		[]string{"Port :", "nfs-port", "8000"},
	},func()bool{
		
		url:=app.GetGlobal("nfs-server")
		resp, err := http.Get(url+"authenticate")
		if err != nil {
			log.Printf("Get")
			return false
		}
		if resp.StatusCode <300 {
			return true
		}
		return false},
		func()bool{
		
			url:=app.GetGlobal("nfs-server")
			resp, err := http.Get(url+"authenticate")
			if err != nil {
				log.Printf("Get failed %v", err)
				return false
			}
			if resp.StatusCode >299 {
				log.Printf("Get failed with code %v", resp.StatusCode)
				return false
			}
			mountUrl(url)
		return false
		},
	)
}

func TestGet(url string)bool{
	resp, err := http.Get(url+"authenticate")
	if err != nil {
		log.Printf("Get failed:")
		log.Println(err)
		return false
	}
	if resp.StatusCode <300 {
		return true
	}
	log.Println(resp)
	return false}

func AccountWin(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {

	window := nanogui.NewWindow(screen, "Login")

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	
	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1

	window.SetPosition(445, 358)
	layout := nanogui.NewGridLayout(nanogui.Horizontal, 2, nanogui.Middle, 15, 5)
	layout.SetColAlignment(nanogui.Maximum, nanogui.Fill)
	layout.SetColSpacing(10)
	window.SetLayout(layout)

	
		field(window, app, []string{"Account URL :", "earthtide-account", "https://entirety.praeceptamachinae.com/"})
	

	b4 := nanogui.NewButton(window, "Connect")
	b4.SetCallback(func() {
		if TestGet(app.GetGlobal("earthtide-account")){
			b4.SetBackgroundColor(nanovgo.RGBA(0, 255, 0, 255))
		} else {
			b4.SetBackgroundColor(nanovgo.RGBA(255, 0, 0, 255))
		 }
	})

	
	b5 := nanogui.NewButton(window, "Menu")
	b5.SetCallback(func() {
		ControlPanel(app, screen)
		screen.PerformLayout()	

	})
	return window
}



func NFSLocalRepoWin(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {
	
	window := nanogui.NewWindow(screen, "Connect Localnet Repository")

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = "LocalMounterWin"

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1

	window.SetPosition(445, 358)
	layout := nanogui.NewGridLayout(nanogui.Horizontal, 2, nanogui.Middle, 15, 5)
	layout.SetColAlignment(nanogui.Maximum, nanogui.Fill)
	layout.SetColSpacing(10)
	window.SetLayout(layout)


	b5 := nanogui.NewButton(window, "Connect Local")
	b5.SetCallback(func() {
		b5.SetBackgroundColor(nanovgo.RGBA(0, 255, 0, 255))
		if mountLocal(false,false){
			b5.SetBackgroundColor(nanovgo.RGBA(0, 255, 0, 255))
		} else {
			b5.SetBackgroundColor(nanovgo.RGBA(255, 0, 0, 255))
		 }
	})
		
	return window
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

func mountUrl(url string) {
	log.Println("Starting ", url)
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
	MountPoint := nextDrive()
		cmd = []string{exePath, "--url", url, "--mount", MountPoint}

	go func() {
		time.Sleep(5 * time.Second)
		//TODO have vort-fuse write a connected status then pop the window asap
		goof.QC([]string{`explorer`, MountPoint + "\\"})
	}()
	goof.QCI(cmd)

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



func mountLocal( debug, trace bool)bool {
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
return false

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

	b3 := nanogui.NewButton(window, "Stop All")

	statuses:= []*nanogui.Label{}


	repositories := loadUserConfig()
	for _, k := range repositories {
		file := k
		//s := gitremind.Repos[file]

		//node := vv.SubNodes[i]
		repoDetail := loadRepoDetails(file)
		name := repoDetail.DisplayName
	

	st:=nanogui.NewLabel(window, "Status:")
	st.SetFont("sans-bold")
	statuses=append(statuses, st)
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


	}

	b3.SetBackgroundColor(nanovgo.RGBA(0, 0, 255, 25))
	b3.SetIcon(nanogui.IconRocket)
	b3.SetCallback(func() {
		//status.SetCaption("Disconnected")
		ioutil.WriteFile(controlDir+"shutdown", []byte(" "), 0600)
		for _,s := range statuses{
			s.SetCaption("Disconnected")
		}
	})


	nanogui.NewResize(window, window)


	return window

}
