// https://mtosak-tech.hatenablog.jp/entry/2020/08/22/114622

package main

import (
	"Converter/MyController"
	"Converter/nscon"
	"os"
	"os/exec"
	"os/signal"
	"time"
)

func setInputButton(input *uint8) {
	*input++
	time.AfterFunc(5*time.Millisecond, func() {
		*input--
	})
}

func main() {
	var con MyController.Controller
	if err := con.Init(); err != nil {
		panic(err)
	}
	defer con.Final()

	if err := con.Connect(); err != nil {
		panic(err)
	}

	if err := con.ShowWindow(); err != nil {
		panic(err)
	}

	con.Print_Device_Data()

	target := "/dev/hidg0"
	ncon := nscon.NewController(target)
	ncon.LogLevel = 0
	defer ncon.Close()
	ncon.Connect()

	// Set tty break for read keyboard input directly
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "-cbreak").Run()
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	ch1 := make(chan string)

	go func() {
		for {
			err := con.Update()
			if err != nil {
				ch1 <- "end"
			}

			if !con.IsConnecting() {
				ch1 <- "end"
			}

			ncon.Input.Stick.Left.X = float64(con.GetStick(MyController.L_X)) / 32767
			ncon.Input.Stick.Left.Y = float64(con.GetStick(MyController.L_Y)) / -32767
			ncon.Input.Stick.Right.X = float64(con.GetStick(MyController.R_X)) / 32767
			ncon.Input.Stick.Right.Y = float64(con.GetStick(MyController.R_Y)) / -32767

			var b uint32

			b = con.GetButton()

			if b != 0 {
				if (b & (1 << 0)) != 0 {
					setInputButton(&ncon.Input.Button.A)
				}
				if (b & (1 << 1)) != 0 {
					setInputButton(&ncon.Input.Button.B)
				}
				if (b & (1 << 2)) != 0 {
					if 0 < ncon.Input.Stick.Left.X {
						// Right
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.X = 1.0
						ncon.Wait()
						ncon.Input.Stick.Left.X = 0.0

						// Down
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.Y = -1.0
						ncon.Wait()
						ncon.Input.Stick.Left.X = 0.0

						// Right Down + A
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.X = 1.0
						ncon.Input.Stick.Left.Y = -1.0
						for !ncon.GetInputOK() {
							setInputButton(&ncon.Input.Button.A)
							time.Sleep(600 * time.Microsecond)
						}
					} else if ncon.Input.Stick.Left.X < 0 {
						// Left
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.X = -1.0
						ncon.Wait()
						ncon.Input.Stick.Left.X = 0.0

						// Down
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.Y = -1.0
						ncon.Wait()
						ncon.Input.Stick.Left.X = 0.0

						// Left Down + A
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.X = -1.0
						ncon.Input.Stick.Left.Y = -1.0
						for !ncon.GetInputOK() {
							setInputButton(&ncon.Input.Button.A)
							time.Sleep(600 * time.Microsecond)
						}
					}
				}
				if (b & (1 << 3)) != 0 {
					setInputButton(&ncon.Input.Button.Y)
				}
				if (b & (1 << 4)) != 0 {
					setInputButton(&ncon.Input.Button.L)
				}
				if (b & (1 << 5)) != 0 {
					setInputButton(&ncon.Input.Button.R)
				}
				if (b & (1 << 6)) != 0 {
					if 0 < ncon.Input.Stick.Left.X {
						// Neutral
						ncon.SetInputOK(false)
						ncon.Reset_Input()
						ncon.Wait()

						// Right
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.X = 1.0
						ncon.Wait()
						ncon.Input.Stick.Left.X = 0.0

						// Down
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.Y = -1.0
						ncon.Wait()
						ncon.Input.Stick.Left.X = 0.0

						// Right Down
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.X = 1.0
						ncon.Input.Stick.Left.Y = -1.0
						ncon.Wait()
					} else if ncon.Input.Stick.Left.X < 0 {
						// Neutral
						ncon.SetInputOK(false)
						ncon.Reset_Input()
						ncon.Wait()

						// Left
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.X = -1.0
						ncon.Wait()
						ncon.Input.Stick.Left.X = 0.0

						// Down
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.Y = -1.0
						ncon.Wait()
						ncon.Input.Stick.Left.X = 0.0

						// Left Down
						ncon.SetInputOK(false)
						ncon.Input.Stick.Left.X = -1.0
						ncon.Input.Stick.Left.Y = -1.0
						ncon.Wait()
					}

					ncon.Reset_Input()
				}
				if (b & (1 << 7)) != 0 {
					setInputButton(&ncon.Input.Button.ZR)
				}
				if (b & (1 << 8)) != 0 {
					setInputButton(&ncon.Input.Dpad.Up)
				}
				if (b & (1 << 9)) != 0 {
					setInputButton(&ncon.Input.Dpad.Down)
				}
				if (b & (1 << 10)) != 0 {
					setInputButton(&ncon.Input.Dpad.Left)
				}
				if (b & (1 << 11)) != 0 {
					setInputButton(&ncon.Input.Dpad.Right)
				}
				if (b & (1 << 12)) != 0 {
					setInputButton(&ncon.Input.Button.Plus)
				}
				if (b & (1 << 13)) != 0 {
					setInputButton(&ncon.Input.Button.Minus)
				}
				if (b & (1 << 14)) != 0 {
					setInputButton(&ncon.Input.Button.Home)
				}
				if (b & (1 << 15)) != 0 {
					setInputButton(&ncon.Input.Button.Capture)
				}
				if (b & (1 << 16)) != 0 {
					setInputButton(&ncon.Input.Stick.Left.Press)
				}
				if (b & (1 << 17)) != 0 {
					setInputButton(&ncon.Input.Stick.Right.Press)
				}

				// End Command
				if (b&(1<<12)) != 0 && (b&(1<<13)) != 0 && (b&(1<<8)) != 0 {
					ch1 <- "end"
				}
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	ch2 := make(chan os.Signal, 1)
	signal.Notify(ch2, os.Interrupt)

	select {
	case <-ch1:
		return
	case <-ch2:
		return
	}
}
