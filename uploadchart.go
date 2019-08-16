package main

import (
	"encoding/base64"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"github.com/mholt/archiver"
	"github.com/urfave/cli"
	"log"
	"os"
)

const (
	USERAGENT        = "lujin httpclient"
	TIMEOUT          = 30
)

func postfile (url string, filename string, username string, password string) (code int, text string) {
	httpclient.Defaults(httpclient.Map {
		"ops_useragent": USERAGENT,
		"ops_timeoue": TIMEOUT,
		"Accept-Encoding": "gzip,deflate,sdch",
	})
	d := username + ":" + password
	auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(d))
	//fmt.Println(auth)
	res, err := httpclient.
		WithHeader("Authorization",auth).
		WithHeader("accept", "application/json").
		//WithHeader("Content-Type","multipart/form-data").
		//Get("https://harbor.wise-paas.io/api/chartrepo/tekton/charts")
		Post(url, map[string]string {
		"@chart": filename,
		})

	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	//fmt.Printf("code: %d\n", res.StatusCode)
	//fmt.Println(res.ToString())
	text, _ = res.ToString()
	//fmt.Println(text)
	return res.StatusCode, text
}


func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}



func tgzfile(file string, filetgzname string) error {
	if IsDir(file)  {
		// archive format is determined by file extension
		err := archiver.Archive([]string{file}, filetgzname)
		if err != nil {
    		return err
		}
	}else{
		fmt.Printf("%s is not a dir\n", file)
	}
	return nil

}


func main() {
	app := cli.NewApp()
	app.Name = "UploadChart"
	app.Usage = "upload chart .tgz to harbor"
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Load helm chart .tgz file",
			Value: "test.tgz",
		},
		cli.StringFlag{
			Name:  "username, u",
			Usage: "username for harbor",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "passworld for harbor",
		},
		cli.StringFlag{
			Name:  "harborurl, url",
			Usage: "harborurl for harbor",
			Value: "https://harbor/api/chartrepo/tekton/charts",
		},
	}
	app.Action = func(c *cli.Context) error {
		file := c.String("file")
		harborurl := c.String("harborurl")
		username := c.String("username")
		password := c.String("password")
		//fmt.Println(file,harborurl,username,password)

		var filetgzname string

		if IsDir(file) {
			filetgzname = file + ".tgz"
		}else{
			filetgzname = file
		}


		if file != "" {
			err := tgzfile(file,filetgzname)
			if err != nil {
				fmt.Println(err)
			}
			code, text := postfile(harborurl, filetgzname, username, password)
			fmt.Printf("code: %d \nresponse: %s\n", code,text)
		}else {
			//time.Sleep(time.Second * 10)
			fmt.Printf("--help, -h            show help \n")
			//time.Sleep(time.Minute * 10)
		}
		//fmt.Println(code,text)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
