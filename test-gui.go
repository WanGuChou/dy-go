package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	log.Println("Testing Fyne GUI...")
	fmt.Println("If you see this, terminal output works!")
	
	log.Println("Creating Fyne app...")
	myApp := app.New()
	
	log.Println("Creating window...")
	myWindow := myApp.NewWindow("GUI Test")
	
	log.Println("Setting content...")
	myWindow.SetContent(widget.NewLabel("Hello! GUI works! 🎉"))
	
	log.Println("Showing window...")
	myWindow.ShowAndRun()
	
	log.Println("Window closed")
}
