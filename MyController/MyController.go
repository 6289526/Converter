package MyController

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"Converter/img"
	"Converter/sdl"
	"Converter/ttf"
)

type ControllerType int

const (
	Fighting_Commander = iota
	Horipad_S
)

type StickType int

const (
	L_X = iota
	L_Y
	R_X
	R_Y
)

type Controller struct {
	window            *sdl.Window
	renderer          *sdl.Renderer
	controller        *sdl.GameController
	joy               *sdl.Joystick
	button            map[int]sdl.GameControllerButton
	con_type          ControllerType
	name              string
	guid              string
	analog_input_num  int
	digital_input_num int
	hats_input_num    int
	connecting        bool
}

// Seach Controller Mapping Data in Data Base Text File
func SeachMappingData(file_name string, controller_name string) string {
	fp, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		str := scanner.Text()
		arr := strings.Split(str, ",")

		if controller_name == arr[1] {
			return str[0 : len(str)-len(arr[len(arr)-2])-2]
		}
	}

	return ""
}

// Set Connected Controller Data
func (c *Controller) SetControllerData() {
	c.name = c.controller.Name()

	c.joy = c.controller.Joystick()
	var g sdl.JoystickGUID = c.joy.GUID()
	var guid_s string = sdl.JoystickGetGUIDString(g)
	c.guid = guid_s

	c.analog_input_num = c.joy.NumAxes()
	c.digital_input_num = c.joy.NumButtons()
	c.hats_input_num = c.joy.NumHats()

	if c.name == "HORI Fighting Commander" {
		c.con_type = Fighting_Commander
	}
}

// Connect NonConnected Controller
func (c *Controller) AddController(index int) error {
	// Already Connected
	if c.controller != nil {
		c.connecting = true
		return nil
	}

	// Open Controller
	c.controller = sdl.GameControllerOpen(index)

	if c.controller != nil {
		// Success
		c.SetControllerData()
		c.connecting = true
		return nil
	} else {
		// Failure
		c.connecting = false
		sdl.Log("sdl.GameControllerOpen Failure")
		return sdl.GetError()
	}
}

func (c *Controller) RemoveController() {
	// Close Controller
	if c.controller != nil {
		c.controller.Close()
	}

	c.controller = nil
	c.connecting = false
	fmt.Println("hoge")
}

func (c *Controller) GetTrigger(t int16, flag sdl.GameControllerAxis) bool {
	var value int16 = c.controller.Axis(flag)
	if t < value {
		return true
	} else {
		return false
	}
}

func (c *Controller) ShowWindow() error {
	var err error

	c.window, err = sdl.CreateWindow("Converter Program", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 0, 0, sdl.WINDOW_SHOWN)
	if err != nil {
		sdl.Log("sdl.CreateWindow Failure")
		return err
	}

	c.renderer, err = sdl.CreateRenderer(c.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		sdl.Log("sdl.CreateRenderer Failure")
		return err
	}

	tex1, err := img.LoadTexture(c.renderer, "image/im1.png")
	defer tex1.Destroy()
	if err != nil {
		sdl.Log("img.LoadTexture Failure")
		return err
	}

	rect1 := sdl.Rect{0, 0, 200, 200}
	err = c.renderer.Copy(tex1, nil, &rect1)
	if err != nil {
		sdl.Log("renderer.Copy Failure")
		return err
	}

	font, err := ttf.OpenFont("MyricaM/MyricaM.TTC", 24)
	if err != nil {
		sdl.Log("ttf.OpenFont Failure")
		return err
	}

	color := sdl.Color{10, 200, 10, 255}

	var text [13]string

	text[0] = "            Welcome to Converter Program !!!!            "
	text[1] = " "
	text[2] = "     Controller      : "
	text[3] = " Analog   Input Num  : "
	text[4] = " Digital  Input Num  : "
	text[5] = "  Hats    Input Num  : "
	text[6] = " "
	text[7] = "==========================Usage=========================="
	text[8] = " Finish Program      :    Shuts Down Automatically       "
	text[9] = "                          Choose One of Three            "
	text[10] = "1. Press Plus Button and Minus Button and Up Cross Button"
	text[11] = "2. Disconnected Controller                               "
	text[12] = "3. Press Any Key on Connected Keyboard                   "

	var text2 [4]string
	text2[0] = " " + c.name
	text2[1] = " " + strconv.Itoa(c.analog_input_num)
	text2[2] = " " + strconv.Itoa(c.digital_input_num)
	text2[3] = " " + strconv.Itoa(c.hats_input_num)

	w, h := c.window.GetSize()

	for i := 0; i < len(text); i++ {
		sur1, err := font.RenderUTF8Blended(text[i], color)
		if err != nil {
			sdl.Log("font.RenderUTF8Blended Failure")
			return err
		}

		tex2, err := c.renderer.CreateTextureFromSurface(sur1)
		if err != nil {
			sdl.Log("renderer.CreateTextureFromSurface Failure")
			return err
		}

		if i == 2 || i == 3 || i == 4 || i == 5 {
			for len(text2[i-2]) < 20 {
				text2[i-2] += " "
			}

			rect2 := sdl.Rect{10, int32(int32((h-20)/13)*int32(i)) + 10, (w - 20) * 2 / 5, int32(h/13) - 10}
			err = c.renderer.Copy(tex2, nil, &rect2)
			if err != nil {
				sdl.Log("renderer.Copy Failure")
				return err
			}

			sur2, err := font.RenderUTF8Blended(text2[i-2], color)
			if err != nil {
				sdl.Log("font.RenderUTF8Blended Failure")
				return err
			}

			tex2, err := c.renderer.CreateTextureFromSurface(sur2)
			if err != nil {
				sdl.Log("renderer.CreateTextureFromSurface Failure")
				return err
			}

			rect3 := sdl.Rect{(10 + (w - 20)) * 2 / 5, int32(int32((h-20)/13)*int32(i)) + 10, (w - 20) * 3 / 5, int32(h/13) - 10}
			err = c.renderer.Copy(tex2, nil, &rect3)
			if err != nil {
				sdl.Log("renderer.Copy Failure")
				return err
			}
		} else {
			rect2 := sdl.Rect{10, int32(int32((h-20)/13)*int32(i)) + 10, w - 20, int32(h/13) - 10}
			err = c.renderer.Copy(tex2, nil, &rect2)
			if err != nil {
				sdl.Log("renderer.Copy Failure")
				return err
			}
		}
	}

	c.renderer.Present()

	sdl.ShowCursor(sdl.DISABLE)

	return nil
}

func (c *Controller) Init() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		sdl.Log("sdl.Init Failure")
		return err
	}
	if err := img.Init(img.INIT_PNG); err != nil {
		sdl.Log("img.Init Failure")
		return err
	}
	if err := ttf.Init(); err != nil {
		sdl.Log("ttf.Init Failure")
		return err
	}

	c.controller = nil
	c.joy = nil

	c.button = make(map[int]sdl.GameControllerButton)

	c.button[0] = sdl.CONTROLLER_BUTTON_A
	c.button[1] = sdl.CONTROLLER_BUTTON_B
	c.button[2] = sdl.CONTROLLER_BUTTON_X
	c.button[3] = sdl.CONTROLLER_BUTTON_Y
	c.button[4] = sdl.CONTROLLER_BUTTON_LEFTSHOULDER
	c.button[5] = sdl.CONTROLLER_BUTTON_RIGHTSHOULDER
	c.button[6] = sdl.CONTROLLER_BUTTON_PADDLE1
	c.button[7] = sdl.CONTROLLER_BUTTON_PADDLE2
	c.button[8] = sdl.CONTROLLER_BUTTON_DPAD_UP
	c.button[9] = sdl.CONTROLLER_BUTTON_DPAD_DOWN
	c.button[10] = sdl.CONTROLLER_BUTTON_DPAD_LEFT
	c.button[11] = sdl.CONTROLLER_BUTTON_DPAD_RIGHT
	c.button[12] = sdl.CONTROLLER_BUTTON_START
	c.button[13] = sdl.CONTROLLER_BUTTON_BACK
	c.button[14] = sdl.CONTROLLER_BUTTON_GUIDE
	c.button[15] = sdl.CONTROLLER_BUTTON_MISC1
	c.button[16] = sdl.CONTROLLER_BUTTON_LEFTSTICK
	c.button[17] = sdl.CONTROLLER_BUTTON_RIGHTSTICK

	c.con_type = Horipad_S
	c.analog_input_num = 0
	c.digital_input_num = 0
	c.hats_input_num = 0
	c.connecting = false

	return nil
}

func (c *Controller) Final() {
	if c.controller != nil {
		c.controller.Close()
	}

	c.controller = nil

	c.renderer.Destroy()
	c.window.Destroy()
	ttf.Quit()
	img.Quit()
	sdl.Quit()
}

func (c *Controller) Connect() error {
	// Connected Controller Num
	var joystick_num int = sdl.NumJoysticks()

	// Connect Controller Num Error
	if joystick_num == 0 {
		return errors.New("Controller Not Connected")
	}

	// Connect Controller Error
	for i := 0; i < joystick_num; i++ {
		// Ignore Unsupported Controller
		if !sdl.IsGameController(i) {
			continue
		}

		if err := c.AddController(i); err == nil {
			// Success
			break
		} else {
			// Failure
			return err
		}
	}

	str := SeachMappingData("gamecontrollerdb.txt", c.name)

	if len(str) != 0 {
		if sdl.GameControllerAddMapping(str) == -1 {
			sdl.Log("sdl.GameControllerAddMapping Failure")
			return sdl.GetError()
		}
	} else {
		f, err := os.OpenFile("gamecontrollerdb.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		str := sdl.GameControllerMappingForGUID(c.joy.GUID()) + ",platform:Linux,"
		_, err = fmt.Fprintln(f, str)
		if err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (c *Controller) Update() error {
	for {
		event := sdl.PollEvent()
		if event == nil {
			break
		}

		switch event.GetType() {
		case sdl.CONTROLLERDEVICEADDED:
			if err := c.Connect(); err != nil {
				return err
			}
		case sdl.CONTROLLERDEVICEREMOVED:
			c.RemoveController()
		case sdl.QUIT:
			return errors.New("Escape_Click")
		case sdl.KEYDOWN:
			return errors.New("Escape_Click")
		}
	}

	return nil
}

func (c *Controller) IsConnecting() bool {
	return c.connecting
}

func (c *Controller) GetButton() uint32 {
	var ret uint32 = 0

	for k, v := range c.button {
		if c.controller.Button(v) == 1 {
			ret |= (1 << k)
		}
	}

	if c.con_type == Fighting_Commander {
		if c.GetTrigger(16384, sdl.CONTROLLER_AXIS_TRIGGERLEFT) {
			ret |= (1 << 4)
		}
		if c.GetTrigger(16384, sdl.CONTROLLER_AXIS_TRIGGERRIGHT) {
			ret |= (1 << 5)
		}
	}

	return ret

}

func (c *Controller) GetStick(flag StickType) int16 {
	var ret int16 = 0
	switch flag {
	case L_X:
		ret = c.controller.Axis(sdl.CONTROLLER_AXIS_LEFTX)
		break
	case L_Y:
		ret = c.controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
		break
	case R_X:
		ret = c.controller.Axis(sdl.CONTROLLER_AXIS_RIGHTX)
		break
	case R_Y:
		ret = c.controller.Axis(sdl.CONTROLLER_AXIS_RIGHTY)
		break
	}
	return ret
}

func (c *Controller) Print_Device_Data() {
	fmt.Printf("\n-----------------------------------------------------------\n")
	fmt.Printf("  Controler  Name  : %s \n", c.name)
	fmt.Printf("       GUID        : %s \n", c.guid)
	fmt.Printf(" Analog  Input Num : %d \n", c.analog_input_num)
	fmt.Printf(" Digital Input Num : %d \n", c.digital_input_num)
	fmt.Printf(" Hats    Input Num : %d \n", c.hats_input_num)
	fmt.Printf("-----------------------------------------------------------\n")
}
