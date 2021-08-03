package demo

import (
	"os/signal"
	"syscall"

	nanogui "../.."

	"context"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	vnc "github.com/amitbet/vnc2video"
	"github.com/amitbet/vnc2video/logger"
	"github.com/donomii/glim"
)

func VncWin(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {

	window := nanogui.NewWindow(screen, "Vnc Window")

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = "VncWin"

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1
	window.SetPosition(545, 15)
	nanogui.NewResize(window, window)
	window.SetLayout(nanogui.NewGroupLayout())

	img := nanogui.NewImageView(window)
	img.SetPolicy(nanogui.ImageSizePolicyExpand)
	img.SetFixedSize(350, 350)

	go doVnc("192.168.178.39:5900", app, img)
	return window

}
func doVnc(server string, app *nanogui.Application, img *nanogui.ImageView) {
	runtime.GOMAXPROCS(4)

	// Establish TCP connection to VNC server.
	nc, err := net.DialTimeout("tcp", server, 5*time.Second)
	if err != nil {
		logger.Fatalf("Error connecting to VNC host. %v", err)
	}

	logger.Tracef("starting up the client, connecting to: %s", server)
	// Negotiate connection with the server.
	cchServer := make(chan vnc.ServerMessage)
	cchClient := make(chan vnc.ClientMessage)
	errorCh := make(chan error)

	ccfg := &vnc.ClientConfig{
		SecurityHandlers: []vnc.SecurityHandler{
			//&vnc.ClientAuthATEN{Username: []byte(os.Args[2]), Password: []byte(os.Args[3])}
			&vnc.ClientAuthVNC{Password: []byte("aaaaaa")},
			&vnc.ClientAuthNone{},
		},
		DrawCursor:      true,
		PixelFormat:     vnc.PixelFormat32bit,
		ClientMessageCh: cchClient,
		ServerMessageCh: cchServer,
		Messages:        vnc.DefaultServerMessages,
		Encodings: []vnc.Encoding{
			&vnc.RawEncoding{},
			&vnc.TightEncoding{},
			&vnc.HextileEncoding{},
			&vnc.ZRLEEncoding{},
			&vnc.CopyRectEncoding{},
			&vnc.CursorPseudoEncoding{},
			&vnc.CursorPosPseudoEncoding{},
			&vnc.ZLibEncoding{},
			&vnc.RREEncoding{},
		},
		ErrorCh: errorCh,
	}

	cc, err := vnc.Connect(context.Background(), nc, ccfg)
	screenImage := cc.Canvas
	if err != nil {
		logger.Fatalf("Error negotiating connection to VNC host. %v", err)
	}

	counter := 0

	//screenImage := vnc.NewVncCanvas(int(cc.Width()), int(cc.Height()))
	//rect := image.Rect(0, 0, int(cc.Width()), int(cc.Height()))
	//screenImage := image.NewRGBA64(rect)

	log.Printf("Using screen image %+v\n", screenImage.Bounds())
	for _, enc := range ccfg.Encodings {
		myRenderer, ok := enc.(vnc.Renderer)

		if ok {
			log.Printf("Supported encoding: %v", enc.Type().String())
			myRenderer.SetTargetImage(screenImage)
		}
	}

	logger.Tracef("connected to: %s", server)
	defer cc.Close()

	cc.SetEncodings([]vnc.EncodingType{
		vnc.EncCursorPseudo,
		vnc.EncPointerPosPseudo,
		vnc.EncCopyRect,
		vnc.EncTight,
		vnc.EncZRLE,
		//	vnc.EncHextile,
		//	vnc.EncZlib,
		//	vnc.EncRRE,
	})
	//rect := image.Rect(0, 0, int(cc.Width()), int(cc.Height()))
	//screenImage := image.NewRGBA64(rect)
	// Process messages coming in on the ServerMessage channel.

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	frameBufferReq := 0

	reqMsg := vnc.FramebufferUpdateRequest{Inc: 0, X: 0, Y: 0, Width: cc.Width(), Height: cc.Height()}
	//cc.ResetAllEncodings()
	reqMsg.Write(cc)

	for {
		select {
		case err := <-errorCh:
			panic(err)
		case msg := <-cchClient:
			logger.Tracef("Received client message type:%v msg:%v\n", msg.Type(), msg)
		case msg := <-cchServer:
			fmt.Printf("Received server message type:%v msg:%v\n", msg.Type(), msg)

			if msg.Type() == vnc.FramebufferUpdateMsgType {

				fmt.Printf("Received FramebufferUpdateMsgType message type:%v msg:%+v\n", msg.Type(), msg)
				//secsPassed := time.Now().Sub(timeStart).Seconds()
				frameBufferReq++
				//reqPerSec := float64(frameBufferReq) / secsPassed
				counter++
				/*result := effect.EdgeDetection(screenImage, 1.0)
				//result = transform.Resize(result, 350, 350, transform.Linear)
				result = effect.Grayscale(result)
				result = effect.Dilate(result, 1)
				result = adjust.Contrast(result, 2.0)
				*/

				p, w, h := glim.GFormatToImage(screenImage.Image, nil, 0, 0)
				p = glim.ForceAlpha(p, 255)
				o := glim.ImageToGFormat(w, h, p)
				/*
									out, err := os.Create("./output" + strconv.Itoa(counter) + ".jpg")
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

							jpeg.Encode(out, o, &jpeg.Options{100})

							f, err := os.Create("output" + strconv.Itoa(counter) + ".png")
							if err != nil {
								panic(err)
							}
							if err := png.Encode(f, o); err != nil {
								panic(err)
							}
							if err := f.Close(); err != nil {
								panic(err)
							}
				*/
				app.MainThreadThunker <- func() {

					ctx := app.Screen.NVGContext()
					//gr := ctx.CreateImageFromGoImage(0, nanogui.StripChart(dt.Series[0][1]))
					gr := ctx.CreateImageFromGoImage(0, o)
					img.SetImage(gr)
					log.Println("Frame")

				}

				reqMsg := vnc.FramebufferUpdateRequest{Inc: 0, X: 0, Y: 0, Width: cc.Width(), Height: cc.Height()}
				//cc.ResetAllEncodings()
				reqMsg.Write(cc)
			}
		}
	}
}
