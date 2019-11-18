package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	LINE   = "\n"
	C      = "c"
	CPP    = "c++"
	GOLANG = "go"
	JS     = "js"
)

type XrmLan struct {
	lanName      string
	ignoreFiles  []string
	ignoreDirs   []string
	fileNameExts []string
}

type Xrm struct {
	xrmLans   map[string]*XrmLan
	lans      map[string]bool
	filePaths []string
}

func MakeXrm() *Xrm {
	result := &Xrm{}
	result.xrmLans = make(map[string]*XrmLan)
	result.xrmLans[C] = &XrmLan{
		lanName:      C,
		fileNameExts: []string{"c", "h"},
	}
	result.xrmLans[CPP] = &XrmLan{
		lanName:      CPP,
		fileNameExts: []string{"cpp", "h", "hpp", "cxx", "c++", "cc"},
	}
	result.xrmLans[GOLANG] = &XrmLan{
		lanName:      GOLANG,
		fileNameExts: []string{"go"},
		ignoreDirs:   []string{"vendor"},
	}
	result.xrmLans[JS] = &XrmLan{
		lanName:      JS,
		fileNameExts: []string{"js"},
	}
	result.lans = make(map[string]bool)
	result.filePaths = []string{}
	return result
}

func (self *Xrm) Config(lanName string) bool {
	if _, ok := self.xrmLans[lanName]; !ok {
		log.Println("[error]", "not support", lanName)
		return false
	}
	_, ok := self.lans[lanName]
	if !ok {
		self.lans[lanName] = true
	}
	return true
}

func (self *Xrm) ConfigCodeDir(dir []string) {
	self.filePaths = dir
}

func (self *Xrm) Execute() bool {
	for lanName, _ := range self.lans {
		xfiles := []string{}
		xrmLan, ok := self.xrmLans[lanName]
		if !ok {
			log.Println("[error]", "not support", lanName)
			continue
		}
		for i := 0; i < len(self.filePaths); i++ {
			for j := 0; j < len(xrmLan.fileNameExts); j++ {
				files, err := WalkDir(self.filePaths[i], xrmLan.fileNameExts[i])
				if err != nil {
					log.Println("[error]", err)
					continue
				}
				xfiles = append(xfiles, files...)
			}
		}

		for _, fileName := range xfiles {
			if CheckIgnoreFileName(fileName, xrmLan) {
				continue
			}
			log.Println(fileName + " XrmComment_SW start")
			XrmComment_SW(fileName, "//")
			log.Println(fileName + " XrmComment_SW end")
			log.Println(fileName + " XrmComment_SE start")
			XrmComment_SE(fileName, "/*", "*/")
			log.Println(fileName + " XrmComment_SE end")
		}
	}
	return true
}

func CheckIgnoreFileName(fileName string, xrmLan *XrmLan) bool {
	ignore := false
	for _, ignoreFileName := range xrmLan.ignoreFiles {
		if fileName == ignoreFileName {
			ignore = true
			break
		}
	}
	if ignore {
		return true
	}
	ignore = false
	for _, ignoreFileName := range xrmLan.ignoreDirs {
		if strings.Contains(fileName, ignoreFileName) {
			ignore = true
			break
		}
	}
	return ignore
}

func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Panic(err)
		}
		if fi.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}

func XrmComment_SW(fileName, prefix string) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	fileStr := string(bytes)
	lines := strings.Split(fileStr, LINE)
	totalFileStr := ""

	for _, s := range lines {
		commentCount := strings.Count(s, prefix)
		if commentCount <= 0 {
			// line := XrmComment_AdjustEmptyLine(s)
			// if line != "" {
			totalFileStr += s + LINE
			// }
			continue
		}

		for {
			if strings.Contains(s, "+build") || strings.Contains(s, "export") {
				break
			}
			index := strings.Index(s, prefix)
			if index < 0 {
				break
			}
			line := s[0:index]
			testStr := strings.Replace(line, "\\\"", "", -1)
			if strings.Count(testStr, "\"")%2 != 0 {
				break
			}

			httpLen := len(line) - len("http:")
			if httpLen > 0 {
				testStr := line[httpLen:]
				if testStr == "http:" {
					break
				}
			}

			httpsLen := len(line) - len("https:")
			if httpsLen > 0 {
				testStr := line[httpsLen:]
				if testStr == "https:" {
					break
				}
			}

			s = line
		}

		line := XrmComment_AdjustEmptyLine(s)
		if line != "" {
			totalFileStr += line + LINE
		}
		if strings.Contains(line, "+build") {
			totalFileStr += LINE
		}
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	if _, err := file.Write([]byte(totalFileStr)); err != nil {
		log.Println("[error]", err)
		return
	}

	if err := file.Close(); err != nil {
		log.Println("[error]", err)
		return
	}
}

func XrmComment_SE(fileName, prefix, suffix string) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	totalFileStr := string(bytes)

	for {
		preIndex := strings.Index(totalFileStr, prefix)
		if preIndex < 0 {
			break
		}

		sufIndex := strings.Index(totalFileStr, suffix)
		if sufIndex < 0 {
			break
		}

		preUsefulStr := totalFileStr[0:preIndex]
		testStr := strings.Replace(preUsefulStr, "\\\"", "", -1)
		if strings.Count(testStr, "\"")%2 != 0 {
			break
		}

		sufUsefulStr := totalFileStr[sufIndex+len(suffix):]
		middleStr := totalFileStr[preIndex:sufIndex]
		if strings.Contains(middleStr, "+build") {
			break
		}
		if strings.Contains(middleStr, "#include") || strings.Contains(middleStr, "extern") {
			break
		}

		totalFileStr = preUsefulStr + sufUsefulStr
	}

	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Println("[error]", err)
		return
	}

	if _, err := file.Write([]byte(totalFileStr)); err != nil {
		log.Println("[error]", err)
		return
	}

	if err := file.Close(); err != nil {
		log.Println("[error]", err)
		return
	}
}

func XrmComment_AdjustEmptyLine(line string) string {
	result := ""

	s := strings.Replace(line, " ", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	if s != "" {
		result = line
	}

	return result
}
