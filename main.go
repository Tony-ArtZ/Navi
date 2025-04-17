package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/navi/constants"
	"github.com/navi/filemanager"
	"github.com/navi/utils"
)

func showSplashScreen() {
	splash := `
    ╔══════════════════════════════════════╗
    ║                                      ║
    ║    ███╗   ██╗ █████╗ ██╗   ██╗██╗    ║
    ║    ████╗  ██║██╔══██╗██║   ██║██║    ║
    ║    ██╔██╗ ██║███████║██║   ██║██║    ║
    ║    ██║╚██╗██║██╔══██║╚██╗ ██╔╝██║    ║
    ║    ██║ ╚████║██║  ██║ ╚████╔╝ ██║    ║
    ║    ╚═╝  ╚═══╝╚═╝  ╚═╝  ╚═══╝  ╚═╝    ║
    ║                                      ║
    ║      Terminal File Manager v1.0      ║
    ║         github.com/Tony-ArtZ         ║
    ║                                      ║
    ╚══════════════════════════════════════╝
    `
	fmt.Print("\033[H\033[J") // Clear screen
	fmt.Printf("%s%s%s%s%s\n\n",
		constants.BLUE_FG,
		constants.BOLD,
		splash,
		constants.RESET_COLOR,
		constants.RESET_COLOR)
	time.Sleep(1 * time.Second)
}

func main() {
	fm, err := filemanager.New()
	if err != nil {
		fmt.Println("Error initializing file manager:", err)
		return
	}

	utils.EnableRawMode()
	defer utils.DisableRawMode()

	showSplashScreen()

	for {
		fm.Render()

		key, err := utils.ReadKey()
		if err != nil {
			break
		}

		if fm.InputMode {
			switch key {
			case '\n':
				fm.HandleInput(fm.InputBuffer)
			case 27:
				fm.InputMode = false
				fm.InputBuffer = ""
				fm.InputPrompt = ""
				fm.InputHandler = nil
			case 127:
				if len(fm.InputBuffer) > 0 {
					fm.InputBuffer = fm.InputBuffer[:len(fm.InputBuffer)-1]
				}
			default:
				if key >= 32 && key <= 126 {
					fm.InputBuffer += string(key)
				}
			}
			continue
		}

		switch key {
		case 'q':
			fmt.Print("\033[H\033[J")
			return
		case 'c':
			fm.PathBuffer = filepath.Join(fm.CurrentPath, fm.Files[fm.Cursor])
			fm.IsMoving = false
			fm.StatusMsg = "File copied to buffer"
			fm.StatusTime = time.Now()
		case 'x':
			fm.PathBuffer = filepath.Join(fm.CurrentPath, fm.Files[fm.Cursor])
			fm.IsMoving = true
			fm.StatusMsg = "File cut to buffer"
			fm.StatusTime = time.Now()
		case 'v':
			fm.PasteFile()
		case 'w':
			if err := os.Chdir(fm.CurrentPath); err != nil {
				fm.StatusMsg = "Error changing working directory: " + err.Error()
			} else {
				fm.StatusMsg = "Working directory changed to: " + fm.CurrentPath
			}
			fm.StatusTime = time.Now()
		case 'n':
			fm.CreateNewFile()
		case 'N':
			fm.CreateNewFolder()
		case 'r':
			fm.RenameFile()
		case 'd':
			fm.DeleteFileOrFolder()
		case 'p':
			fm.TogglePreview()
		case 'o':
			fm.Open()
		case '\n':
			selected := fm.Files[fm.Cursor]
			newPath := filepath.Join(fm.CurrentPath, selected)
			info, err := os.Stat(newPath)
			if err == nil && info.IsDir() {
				fm.CurrentPath = newPath
				fm.Files = utils.ListFiles(fm.CurrentPath)
				fm.Cursor = 0
			}
		case 27:
			if next, err := utils.ReadKey(); err == nil && next == '[' {
				if arrow, err := utils.ReadKey(); err == nil {
					switch arrow {
					case 'A':
						if fm.Cursor > 0 {
							fm.Cursor--
						}
					case 'B':
						if fm.Cursor < len(fm.Files)-1 {
							fm.Cursor++
						}
					}
				}
			}
		default:
			if fm.DeleteConfirmPending {
				fm.DeleteConfirmPending = false
			}
		}
	}
}
