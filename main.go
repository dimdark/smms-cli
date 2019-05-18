package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	smms_api "github.com/dimdark/smms-cli/smms-api"
	"github.com/urfave/cli"
)

func Clear(c *cli.Context) error {
	rst, err := smms_api.ClearUploadHistory()
	if err != nil {
		return err
	}
	fmt.Println(rst.Code)
	return nil
}

func Delete(c *cli.Context) error {
	delPictureUrl := c.Args().First()
	rst := smms_api.DeleteUnloadedPicture(delPictureUrl)
	fmt.Println(rst)
	return nil
}

func List(c *cli.Context) error {
	rst, err := smms_api.ListUploadHistory()
	if err != nil {
		return err
	}
	for idx, record := range rst.Data {
		fmt.Printf("----------------------------\n")
		fmt.Printf("Id = %d\n", idx+1)
		fmt.Printf("Filename: %s\n", record.FileName)
		fmt.Printf("FileInfo: %d x %d\n", record.Width, record.Height)
		fmt.Printf("StoreName: %s\n", record.StoreName)
		// fmt.Printf("Size: %d\n", record.Size)
		fmt.Printf("Path: %s\n", record.Path)
		// fmt.Printf("Hash: %s\n", record.Hash)
		fmt.Printf("TimeStamp: %d\n", record.TimeStamp)
		fmt.Printf("Url: %s\n", record.Url)
		fmt.Printf("Delete url link: %s\n", record.DeleteUrl)
		fmt.Printf("----------------------------\n")
	}
	return nil
}

func Upload(c *cli.Context) error {
	filename := c.Args().First()

	if isDir(filename) {
		fileInfos, err := ioutil.ReadDir(filename)
		if err != nil {
			return err
		}
		for _, fileInfo := range fileInfos {
			if isPicture, _ := smms_api.CheckFileSuffix(fileInfo.Name()); isPicture {
				var absFilename string
				if strings.HasSuffix(filename, "/") {
					absFilename = filename + fileInfo.Name()
				} else {
					absFilename = filename + "/" + fileInfo.Name()
				}
				err = uploadPicture(absFilename)
				if err != nil {
					return err
				}
			}
		}
	} else {
		err := uploadPicture(filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadPicture(absFilename string) error {
	rst, err := smms_api.UploadPicture(absFilename)
	if err != nil {
		return err
	}
	if rst.Msg != "" {
		fmt.Printf("msg: %s\n", rst.Msg)
	} else {
		pictureMarkdownUrl := fmt.Sprintf("![%s](%s)", rst.Data.FileName, rst.Data.Url)
		fmt.Printf("md: %s\n", pictureMarkdownUrl)
		fmt.Printf("del: %s\n", rst.Data.DeleteUrl)
	}
	return nil
}

func isDir(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

// config smms client commands
func configCliCommands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:    "clear",
			Aliases: []string{"c"},
			Usage:   "clear your sm.ms picture upload history",
			Action:  Clear,
		},
		{
			Name:    "delete",
			Aliases: []string{"d"},
			Usage:   "delete picture that you uploaded to sm.ms",
			Action:  Delete,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list upload history about pictures that uploaded to sm.ms",
			Action:  List,
		},
		{
			Name:    "upload",
			Aliases: []string{"u"},
			Usage:   "upload the picture or the pictures in the folder to sm.ms",
			Action:  Upload,
		},
	}
}

func main() {
	app := cli.NewApp()

	app.Name = "smms"
	app.Compiled = time.Now()
	app.Version = "v0.1.0"
	app.Authors = []cli.Author{
		{
			Name:  "dimdark",
			Email: "13760693284@163.com",
		},
	}
	app.Copyright = "(c) 2019 dimdark<13760693284@163.com>"
	app.Usage = "simple client for sm.ms"

	// 配置命令
	configCliCommands(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
