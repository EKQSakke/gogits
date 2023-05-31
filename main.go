package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	list := tview.NewList().ShowSecondaryText(false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetText("Searching Git folders..."), 1, 1, false).
		AddItem(list, 0, 1, true)

	list.SetSelectedFunc(func(index int, _ string, _ string, _ rune) {
		path, _ := list.GetItemText(index)
		changeDirectory(path)
		app.Stop()
	})

	if err := app.SetRoot(flex, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}

	listGitFoldersFromFile("C:/settings/gitfolders.txt", list)
}

func listGitFoldersFromFile(filePath string, list *tview.List) {
	list.Clear()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		path := strings.TrimSpace(scanner.Text())
		listGitFolders(path, list)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file: %v\n", err)
		return
	}
}

func listGitFolders(rootPath string, list *tview.List) {
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			gitFolder := filepath.Dir(path)
			list.AddItem(gitFolder, "", 0, nil)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func changeDirectory(path string) {
	cmd := exec.Command("cmd", "/c", "start", "cmd", "/k", "cd", "/d", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Error changing directory: %v\n", err)
	}
}
