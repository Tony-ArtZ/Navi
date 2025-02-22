package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/navi/constants"
	"github.com/navi/utils"
)

type FileManager struct {
	CurrentPath          string
	Files                []string
	Cursor               int
	FileInfo             os.FileInfo
	PathBuffer           string
	IsMoving             bool
	StatusMsg            string
	StatusTime           time.Time
	InputMode            bool
	InputBuffer          string
	InputPrompt          string
	InputHandler         func(string)
	DeleteConfirmPending bool
	DeleteConfirmTime    time.Time
}

func New() (*FileManager, error) {
	currentPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fm := &FileManager{
		CurrentPath: currentPath,
		Files:       utils.ListFiles(currentPath),
	}
	return fm, nil
}

func (fm *FileManager) HandleInput(input string) {
	if input == "" {
		fm.setStatus("Operation cancelled")
		return
	}

	if fm.InputHandler != nil {
		fm.InputHandler(input)
	}

	fm.InputMode = false
	fm.InputBuffer = ""
	fm.InputPrompt = ""
	fm.InputHandler = nil
}

func (fm *FileManager) setStatus(msg string) {
	fm.StatusMsg = msg
	fm.StatusTime = time.Now()
}

func (fm *FileManager) PasteFile() {
	if fm.PathBuffer == "" {
		fm.setStatus("No file in buffer to paste")
		return
	}

	destPath := filepath.Join(fm.CurrentPath, filepath.Base(fm.PathBuffer))

	if _, err := os.Stat(fm.PathBuffer); err != nil {
		fm.setStatus("Source file no longer exists")
		return
	}

	if fm.IsMoving {
		if err := os.Rename(fm.PathBuffer, destPath); err != nil {
			fm.setStatus("Error moving file: " + err.Error())
		} else {
			fm.setStatus("File moved successfully")
			fm.PathBuffer = ""
		}
	} else {
		if err := utils.CopyFile(fm.PathBuffer, destPath); err != nil {
			fm.setStatus("Error copying file: " + err.Error())
		} else {
			fm.setStatus("File copied successfully")
		}
	}

	fm.Files = utils.ListFiles(fm.CurrentPath)
}

func (fm *FileManager) CreateNewFile() {
	fm.InputMode = true
	fm.InputPrompt = "Enter new file name: "
	fm.InputHandler = func(name string) {
		newPath := filepath.Join(fm.CurrentPath, name)
		if err := utils.CreateNewFile(newPath); err != nil {
			fm.setStatus("Error creating file: " + err.Error())
		} else {
			fm.setStatus("File created: " + name)
		}
		fm.Files = utils.ListFiles(fm.CurrentPath)
	}
}

func (fm *FileManager) CreateNewFolder() {
	fm.InputMode = true
	fm.InputPrompt = "Enter new folder name: "
	fm.InputHandler = func(name string) {
		newPath := filepath.Join(fm.CurrentPath, name)
		if err := utils.CreateNewFolder(newPath); err != nil {
			fm.setStatus("Error creating folder: " + err.Error())
		} else {
			fm.setStatus("Folder created: " + name)
		}
		fm.Files = utils.ListFiles(fm.CurrentPath)
	}
}

func (fm *FileManager) RenameFile() {
	if len(fm.Files) == 0 {
		return
	}
	oldName := fm.Files[fm.Cursor]
	fm.InputMode = true
	fm.InputPrompt = "Enter new name: "
	fm.InputHandler = func(newName string) {
		oldPath := filepath.Join(fm.CurrentPath, oldName)
		newPath := filepath.Join(fm.CurrentPath, newName)
		if err := utils.RenameFile(oldPath, newPath); err != nil {
			fm.setStatus("Error renaming: " + err.Error())
		} else {
			fm.setStatus("Renamed " + oldName + " to " + newName)
		}
		fm.Files = utils.ListFiles(fm.CurrentPath)
	}
}

func (fm *FileManager) DeleteFileOrFolder() {
	if len(fm.Files) == 0 {
		return
	}

	if !fm.DeleteConfirmPending {
		fm.DeleteConfirmPending = true
		fm.DeleteConfirmTime = time.Now()
		fm.setStatus("Press delete again to confirm deletion of: " + fm.Files[fm.Cursor])
		return
	}

	if time.Since(fm.DeleteConfirmTime) > 3*time.Second {
		fm.DeleteConfirmPending = false
		fm.setStatus("Delete timeout - press delete again to start over")
		return
	}

	pathToDelete := filepath.Join(fm.CurrentPath, fm.Files[fm.Cursor])
	err := os.RemoveAll(pathToDelete)
	if err != nil {
		fm.setStatus("Error deleting: " + err.Error())
	} else {
		fm.setStatus("Deleted: " + fm.Files[fm.Cursor])
	}

	fm.DeleteConfirmPending = false
	fm.Files = utils.ListFiles(fm.CurrentPath)
	if fm.Cursor >= len(fm.Files) {
		fm.Cursor = len(fm.Files) - 1
		if fm.Cursor < 0 {
			fm.Cursor = 0
		}
	}
}

func (fm *FileManager) Render() {
	fmt.Print("\033[H\033[J")

	fmt.Printf("%s%s%s %s  Path: %s %s%s\n\n",
		constants.HEADER_BG, constants.WHITE_FG, constants.BOLD,
		constants.PATHICON,
		fm.CurrentPath,
		constants.RESET_COLOR,
		constants.RESET_COLOR)

	if time.Since(fm.StatusTime) < 3*time.Second && fm.StatusMsg != "" {
		fmt.Printf("%s%s%s %s%s\n\n",
			constants.DARK_BG, constants.YELLOW_FG, constants.BOLD,
			fm.StatusMsg,
			constants.RESET_COLOR)
	}

	// Render files list
	for i, file := range fm.Files {
		fullPath := filepath.Join(fm.CurrentPath, file)
		fileInfo, err := os.Stat(fullPath)
		if err != nil {
			fmt.Printf("   %s\t%s\n", file, "error")
			continue
		}

		modTime := utils.FormatDate(fileInfo.ModTime())
		size := utils.FormatSize(fileInfo.Size())
		paddedName := fmt.Sprintf("%-30s", file)
		icon := constants.BLUE_FG + constants.FILEICON
		if fileInfo.IsDir() {
			icon = constants.YELLOW_FG + constants.FOLDERICON
		}

		if i == fm.Cursor {
			fmt.Printf("%s%s-> %s %s    %s%s %s    %s%s %s%s\n",
				constants.DARK_BG, constants.BLUE_FG,
				icon, constants.WHITE_FG+paddedName,
				constants.GREEN_FG, constants.SIZEICON+" "+size,
				constants.YELLOW_FG, constants.CLOCKICON+" "+modTime,
				constants.RESET_COLOR, constants.RESET_COLOR, constants.RESET_COLOR)
		} else {
			fmt.Printf("   %s %s    %s%s %s    %s%s %s%s\n",
				icon, constants.WHITE_FG+paddedName,
				constants.GREEN_FG, constants.SIZEICON+" "+size,
				constants.YELLOW_FG, constants.CLOCKICON+" "+modTime,
				constants.RESET_COLOR, constants.RESET_COLOR, constants.RESET_COLOR)
		}
	}

	// Render input mode or help
	if fm.InputMode {
		fm.renderInputPrompt()
	} else {
		fm.renderHelp()
	}
}

func (fm *FileManager) renderInputPrompt() {
	fmt.Printf("\n%s%s%s%s┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄%s\n",
		constants.DARK_BG, constants.BLUE_FG, constants.BOLD, "╍", constants.RESET_COLOR)
	fmt.Printf("%s%s%s %s%s%s\n",
		constants.DARK_BG, constants.WHITE_FG, constants.BOLD,
		fm.InputPrompt, fm.InputBuffer,
		constants.RESET_COLOR)
	fmt.Printf("%s%s%s%s┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄┄%s\n",
		constants.DARK_BG, constants.BLUE_FG, constants.BOLD, "╍", constants.RESET_COLOR)
}

func (fm *FileManager) renderHelp() {
	fmt.Printf("\n%s%s%s %s ↑/↓: Navigate  %s Enter: Open  %s n: New File  %s N: New Folder  %s r: Rename  %s c: Copy  %s x: Cut  %s v: Paste  %s w: Set PWD  %s d: Delete  %s q: Quit%s\n",
		constants.FOOTER_BG, constants.WHITE_FG, constants.BOLD,
		constants.NAVICON,
		constants.FOLDERICON,
		constants.FILEICON,
		constants.FOLDERICON,
		constants.BLUE_FG+"\uf044",
		constants.BLUE_FG+constants.COPYICON,
		constants.YELLOW_FG+constants.CUTICON,
		constants.GREEN_FG+constants.PASTEICON,
		constants.PATHICON,
		constants.RED_FG+"\uf1f8",
		constants.QUITICON,
		constants.RESET_COLOR)

	if len(fm.Files) > 0 {
		fullPath := filepath.Join(fm.CurrentPath, fm.Files[fm.Cursor])
		if fileInfo, err := os.Stat(fullPath); err == nil {
			fmt.Printf("%s%s%sFile: %s    %s%s%s Size: %s    %s%s%s Modified: %s%s\n",
				constants.FOOTER_BG, constants.WHITE_FG, constants.BOLD,
				fm.Files[fm.Cursor],
				constants.GREEN_FG, constants.SIZEICON,
				constants.WHITE_FG,
				utils.FormatSize(fileInfo.Size()),
				constants.YELLOW_FG, constants.CLOCKICON,
				constants.WHITE_FG,
				utils.FormatDate(fileInfo.ModTime()),
				constants.RESET_COLOR)
		}
	}
}
