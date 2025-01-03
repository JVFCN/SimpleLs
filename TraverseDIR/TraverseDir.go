package TraverseDIR

import (
	"fmt"
	"log"
	"ls/GetTerminalSize"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	DirNum            = 0
	FileNum           = 0
	AllSize           = 0.0
	maxLength         = 0
	outputNum         = 0
	dir, err          = os.Getwd()
	files, errReadDir = os.ReadDir(dir)
)

func List(showTime bool, mod bool, re *regexp.Regexp, all bool) {
	if len(files) == 0 {
		return
	}
	if err != nil {
		fmt.Println("无法获取当前目录：", err)
		return
	}
	if errReadDir != nil {
		fmt.Printf("Error reading directory %s: %s\n", dir, err)
		return
	}

	if all {
		showAll(showTime, mod, re)
	} else {
		showNoAll(showTime, mod, re)
	}
}
func Default(showTime bool, re *regexp.Regexp, all bool) {
	if len(files) == 0 {
		return
	}
	width, err := GetTerminalSize.GetTerminalWidth()
	if err != nil {
		fmt.Println("无法获取终端宽度:", err)
	}
	if err != nil {
		fmt.Println("无法获取当前目录：", err)
		return
	}
	if errReadDir != nil {
		fmt.Printf("Error reading directory %s: %s\n", dir, err)
		return
	}

	var allFileLen int

	// 获取名字最长的文件名字长度
	for _, entry := range files {
		if !isHidden(entry) {
			fileLen := len(entry.Name())
			allFileLen += fileLen + 1
			if length := fileLen; length > maxLength {
				maxLength = length
			}
		}
	}
	// 计算一行最多能放下多少文件
	LineMaxFileNum := width / maxLength
	//fmt.Println("一行最多", LineMaxFileNum, "个文件")

	start := time.Now()

	//fmt.Println("allFileLen:", allFileLen, "\nwidth:", width)
	if width >= allFileLen {
		//fmt.Println("能一行打印完")
		for _, file := range files {
			if !isHidden(file) {
				if file.IsDir() {
					fmt.Print("\033[32m", file.Name(), "/ ")
				} else {
					fmt.Print("\033[36m", file.Name(), "* ")
				}
			}
			fmt.Print("\033[39m")
		}
		if showTime {
			fmt.Println("\n", time.Now().Sub(start))
		}
	} else {
		for _, file := range files {
			if isHidden(file) {
				continue
			}
			fileName := file.Name()
			filePath := string(dir) + "\\" + fileName
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				log.Fatal(err)
				return
			}

			if re.MatchString(fileName) {
				outputNum++
				if file.IsDir() {
					DirNum++
					if outputNum < LineMaxFileNum {
						fmt.Print("\033[32m", fileName, "/")
					} else {
						fmt.Print("\n\033[32m", fileName, "/")
						outputNum = 0
					}
				} else {
					fileSizeF64 := float64(fileInfo.Size())
					FileNum++
					AllSize += fileSizeF64
					if outputNum < LineMaxFileNum {
						fmt.Print("\033[36m", fileName, "*")
					} else {
						fmt.Print("\n\033[36m", fileName, "*")
						outputNum = 0
					}
					fmt.Print("\033[39m")
				}
				fmt.Print(strings.Repeat(" ", maxLength-len(fileName)+1))
			}
		}
		fmt.Print("\033[39m")
		fmt.Print("\nDIR:" + strconv.Itoa(DirNum) + "	File:" + strconv.Itoa(FileNum))
		fmt.Println("\nSize:", fmtSize(AllSize))

		if showTime {
			fmt.Println(time.Now().Sub(start))
		}
	}
}

func fmtSize(size float64) string {
	if size > 1_073_741_824.0 {
		return strconv.FormatFloat(size/1024.0/1024.0/1024.0, 'f', 2, 64) + "G"
	} else if size > 1_048_576.0 {
		return strconv.FormatFloat(size/1024.0/1024.0, 'f', 2, 64) + "MB"
	} else if size > 1024.0 {
		return strconv.FormatFloat(size/1024.0, 'f', 2, 64) + "KB"
	} else {
		return strconv.FormatFloat(size, 'f', 2, 64) + "KB"
	}
}

func showAll(showTime bool, mod bool, re *regexp.Regexp) {
	start := time.Now()

	for _, file := range files {
		fileName := file.Name()
		filePath := string(dir) + "\\" + fileName
		fileInfo, err := os.Stat(filePath)
		fileModTime := fileInfo.ModTime()

		if err != nil {
			log.Fatal(err)
			return
		}

		if re.MatchString(fileName) {
			if file.IsDir() {
				DirNum++
				if mod {
					fmt.Printf(fileModTime.Format("2006-01-02 15:04:05")+"\033[32m%10s	%s", "DIR", fileName+"/\033[39m\n")
				} else {
					fmt.Printf("\033[32m%10s	%s", "DIR", fileName+"/\033[39m\n")
				}
			} else {
				fileSizeF64 := float64(fileInfo.Size())
				FileNum++
				AllSize += fileSizeF64
				if mod {
					fmt.Printf(fileModTime.Format("2006-01-02 15:04:05")+"%10s\033[36m	%s", fmtSize(fileSizeF64), fileName+"*\033[39m\n")
				} else {
					fmt.Printf("%10s\033[36m	%s", fmtSize(fileSizeF64), fileName+"*\033[39m\n")
				}
			}
		}
	}
	fmt.Print("\nDIR:" + strconv.Itoa(DirNum) + "	File:" + strconv.Itoa(FileNum))
	fmt.Println("\nSize:", fmtSize(AllSize))

	if showTime {
		fmt.Println(time.Now().Sub(start))
	}
}

func showNoAll(showTime bool, mod bool, re *regexp.Regexp) {
	start := time.Now()

	for _, file := range files {
		if isHidden(file) {
			continue
		}

		fileName := file.Name()
		filePath := string(dir) + "\\" + fileName
		fileInfo, err := os.Stat(filePath)
		fileModTime := fileInfo.ModTime()

		if err != nil {
			log.Fatal(err)
			return
		}

		if re.MatchString(fileName) {
			if file.IsDir() {
				DirNum++
				if mod {
					fmt.Printf(fileModTime.Format("2006-01-02 15:04:05")+"\033[32m%10s	%s", "DIR", fileName+"/\033[39m\n")
				} else {
					fmt.Printf("\033[32m%10s	%s", "DIR", fileName+"/\033[39m\n")
				}
			} else {
				fileSizeF64 := float64(fileInfo.Size())
				FileNum++
				AllSize += fileSizeF64
				if mod {
					fmt.Printf(fileModTime.Format("2006-01-02 15:04:05")+"%10s\033[36m	%s", fmtSize(fileSizeF64), fileName+"*\033[39m\n")
				} else {
					fmt.Printf("%10s\033[36m	%s", fmtSize(fileSizeF64), fileName+"*\033[39m\n")
				}
			}
		}
	}
	fmt.Print("\nDIR:" + strconv.Itoa(DirNum) + "	File:" + strconv.Itoa(FileNum))
	fmt.Println("\nSize:", fmtSize(AllSize))

	if showTime {
		fmt.Println(time.Now().Sub(start))
	}
}
