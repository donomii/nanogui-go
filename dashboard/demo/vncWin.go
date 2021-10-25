package demo

import (
	"errors"
	"image"
	"image/jpeg"
	"os/signal"
	"syscall"

	nanogui "../.."

	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/amitbet/vnc2video"
	vnc "github.com/amitbet/vnc2video"
	"github.com/amitbet/vnc2video/logger"
	"github.com/donomii/glim"
	"github.com/shibukawa/glfw"
)

func VncAuth(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {
	return AuthWin(app, screen, "Login to Vnc", "VncAuth", [][]string{
		[]string{"Server :", "vnc-server", "localhost"},
		[]string{"Password :", "vnc-password", "aaaaaa"},
		[]string{"Port :", "vnc-port", "5900"},
	})
}

func VncWin(app *nanogui.Application, screen *nanogui.Screen) *nanogui.Window {

	window := nanogui.NewWindow(screen, "Vnc Window")
	window.SetFixedSize(100, 100)
	window.SetSize(200, 200)

	if WindowList == nil {
		WindowList = []*nanogui.Window{}
	}

	WindowList = append(WindowList, window)

	actor := NewActor(window)
	actor.WinType = "VncWin"

	window.WidgetId = fmt.Sprintf("%v", nextWindowId)
	nextWindowId += 1
	window.SetPosition(0, 0)

	window.SetLayout(nanogui.NewStrictLayout())
	nanogui.NewResize(window, window)
	img := nanogui.NewImageView(window)
	img.SetPolicy(nanogui.ImageSizePolicyExpand)
	//img.SetFixedSize(800, 600)
	//img.SetSize(800, 600)
	nanogui.NewResize(window, window)

	go func() {
		for {
			doVnc(app.GetGlobal("vnc-server")+":"+app.GetGlobal("vnc-port"), app, img, window)
		}
	}()
	return window

}

func connectVnc(ctx context.Context, c net.Conn, cfg *vnc.ClientConfig) (conn *vnc.ClientConn, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			conn = nil
			err = errors.New("vnc connection failed")
		}
	}()
	conn, err = vnc.Connect(context.Background(), c, cfg)
	return conn, err
}

func queueWork(p chan func(), f func()) {
	go func() {
		p <- f
	}()
}

var mjpeg = false

func doVnc(server string, app *nanogui.Application, img *nanogui.ImageView, window *nanogui.Window) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in doVnc", r)
		}
	}()

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
	var cc *vnc.ClientConn
	//var err error
	for cc, err = connectVnc(context.Background(), nc, ccfg); err != nil; cc, err = connectVnc(context.Background(), nc, ccfg) {
	}
	screenImage := cc.Canvas
	if err != nil {
		logger.Fatalf("Error negotiating connection to VNC host. %v", err)
	}

	//MJPEG server
	if mjpeg {
		go func() {
			for {
				servAddr := "192.168.178.39:1001"
				tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
				if err != nil {
					println("ResolveTCPAddr failed:", err.Error())

				} else {

					conn, err := net.DialTCP("tcp", nil, tcpAddr)
					if err != nil {
						println("Dial failed:", err.Error())

					} else {

						for err == nil {
							time.Sleep(50 * time.Millisecond)
							var o image.Image
							o, err = jpeg.Decode(conn)

							if err == nil {
								println("read image")
								app.MainThreadThunker <- func() {

									ctx := app.Screen.NVGContext()
									gr := ctx.CreateImageFromGoImage(0, o)
									img.SetImage(gr)
									log.Println("Updated image")

								}
							} else {
								log.Printf("Received corrupted jpeg: %v", err)
							}
						}
						conn.Close()
					}
				}
			}
		}()
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
		//vnc.EncCursorPseudo,
		//vnc.EncPointerPosPseudo,
		//vnc.EncCopyRect,
		vnc.EncTight,
		vnc.EncZRLE,
		vnc.EncHextile,
		vnc.EncZlib,
		vnc.EncRRE,
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

	serverOutPipe := make(chan func(), 20)
	go func() {
		for {
			f := <-serverOutPipe
			f()
		}
	}()

	if !mjpeg {
		go func() {
			for {

				reqMsg1 := vnc.FramebufferUpdateRequest{Inc: 1, X: 0, Y: 0, Width: cc.Width(), Height: cc.Height()}
				queueWork(serverOutPipe, func() {
					//cc.ResetAllEncodings()
					reqMsg1.Write(cc)
				})
				time.Sleep(100 * time.Millisecond)

			}
		}()
	}

	reqMsg := vnc.FramebufferUpdateRequest{Inc: 0, X: 0, Y: 0, Width: cc.Width(), Height: cc.Height()}
	//cc.ResetAllEncodings()
	queueWork(serverOutPipe, func() { reqMsg.Write(cc) })

	img.SetCallback(func(x, y int, button glfw.MouseButton, down bool, modifier glfw.ModifierKey) {
		if down {
			imw, _ := img.Size()

			aspect := float64(cc.Width()) / float64(cc.Height())
			asHeight := 1.0 / aspect * float64(imw)
			rX := int(x-img.WidgetPosX) * int(cc.Width()) / int(imw)
			rY := int(y-img.WidgetPosY) * int(cc.Height()) / int(asHeight)
			b := uint8(button)
			switch button {
			case 2:
				b = 2
			case 1:
				b = 4
			case 0:
				b = 1

			}
			//fmt.Printf("Image %+v\n", img)

			queueWork(serverOutPipe, func() {
				reqMsg := vnc.PointerEvent{Mask: b, X: uint16(rX), Y: uint16(rY)}
				reqMsg.Write(cc)
				fmt.Println("click at ", x, ",", y, "resized ", rX, ",", rY, "button:", b, "remote screen", cc.Width(), ",", cc.Height())
			})

			reqMsg1 := vnc.FramebufferUpdateRequest{Inc: 1, X: 0, Y: 0, Width: cc.Width(), Height: cc.Height()}
			queueWork(serverOutPipe, func() {
				//cc.ResetAllEncodings()
				reqMsg1.Write(cc)

			})
		}
	})

	img.SetMotionCallback(func(x, y, relX, relY int, button int, modifier glfw.ModifierKey) {
		imw, _ := img.Size()

		aspect := float64(cc.Width()) / float64(cc.Height())
		asHeight := 1.0 / aspect * float64(imw)
		//fmt.Printf("Image size: %v,%v\n", imw, asHeight)
		rX := int(x-img.WidgetPosX) * int(cc.Width()) / int(imw)
		rY := int(y-img.WidgetPosY) * int(cc.Height()) / int(asHeight)
		//fmt.Println("Remote scree ", cc.Width(), ",", cc.Height())
		//fmt.Println("Callback motion at ", rX, ",", rY)
		reqMsg := vnc.PointerEvent{Mask: uint8(button), X: uint16(rX), Y: uint16(rY)}
		queueWork(serverOutPipe, func() { reqMsg.Write(cc) })

	})

	img.SetKeyboardEventCallback(func(key glfw.Key, scanCode int, action glfw.Action, modifier glfw.ModifierKey) {

		out := vnc.Space
		fmt.Printf("Got key %+v\n", glfw.Key(scanCode))
		switch scanCode {
		case 36:
			out = vnc.Return
		case 51:
			out = vnc.BackSpace
		case 53:
			out = vnc.Escape
		case 60:
			out = vnc.ShiftRight
		case 56:
			out = vnc.ShiftLeft
		case 57:
			out = vnc.ShiftLock
		case 59:
			out = vnc.ControlLeft
		case 55:
			out = vnc.ControlLeft
		case 58:
			out = vnc.AltLeft
		case 48:
			out = vnc.Tab
		case 123:
			out = vnc.Left
		case 124:
			out = vnc.Right
		case 126:
			out = vnc.Up
		case 125:
			out = vnc.Down
		}

		if out == vnc.Space {
			return
		}
		if action == glfw.Press {
			queueWork(serverOutPipe, func() {
				reqMsg := vnc.KeyEvent{Down: 1, Key: out}
				reqMsg.Write(cc)
				fmt.Printf("Sent key %+v\n", reqMsg)
			})
		} else {

			queueWork(serverOutPipe, func() {
				reqMsg := vnc.KeyEvent{Down: 0, Key: out}
				reqMsg.Write(cc)
				fmt.Printf("Sent key %+v\n", reqMsg)
			})
		}
	})

	img.SetKeyboardCharacterEventCallback(func(c rune) {

		queueWork(serverOutPipe, func() {
			reqMsg := vnc.KeyEvent{Down: 1, Key: vnc2video.Key(uint32(c))}
			reqMsg.Write(cc)
			fmt.Printf("Sent key %+v\n", reqMsg)
		})

		queueWork(serverOutPipe, func() {
			reqMsg := vnc.KeyEvent{Down: 0, Key: vnc2video.Key(uint32(c))}
			reqMsg.Write(cc)
			fmt.Printf("Sent key %+v\n", reqMsg)
		})
	})

	for {
		select {
		case err := <-errorCh:
			panic(err)
		case msg := <-cchClient:
			logger.Tracef("Received client message type:%v msg:%v\n", msg.Type(), msg)
		case msg := <-cchServer:
			logger.Tracef("Received server message type:%v msg:%v\n", msg.Type(), msg)
			//fmt.Printf("Received server message type:%v msg:%v\n", msg.Type(), msg)
			if msg.Type() == vnc.FramebufferUpdateMsgType {

				//fmt.Printf("Received FramebufferUpdateMsgType message type:%v msg:%+v\n", msg.Type(), msg)
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
				//img.SetFixedSize(ww, wh)
				p, w, h := glim.GFormatToImage(screenImage.Image, nil, 0, 0)
				p = glim.ForceAlpha(p, 255)
				p = glim.FlipUp(w, h, p)
				o := glim.ImageToGFormat(w, h, p)
				/*	out, err := os.Create("./output" + strconv.Itoa(counter) + ".jpg")
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					//i := screenImage.Image
					//fmt.Printf("Gimage %+v\n", i.Bounds())

					//jpeg.Encode(out, o, &jpeg.Options{100})

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

				//log.Println("Queueing vnc frame")
				app.MainThreadThunker <- func() {

					ctx := app.Screen.NVGContext()
					gr := ctx.CreateImageFromGoImage(0, o)
					img.SetImage(gr)
					//log.Println("Updated image")

				}

			}
		}
		//log.Println("Finished select, looping")
	}
}
