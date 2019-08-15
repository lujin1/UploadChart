package main

import (
	"encoding/base64"
	"fmt"
	"github.com/ddliu/go-httpclient"
	"github.com/urfave/cli"
	"log"
	"os"
)

const (
	USERAGENT        = "lujin httpclient"
	TIMEOUT          = 30
	SERVER           = "https://harbor.wise-paas.io/api/chartrepo/tekton/charts"
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

func main() {
	app := cli.NewApp()
	app.Name = "UploadChart"
	app.Usage = "upload chart .tgz to harbor"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Load helm chart .tgz file",
			Value: "test2.tgz",
		},
		cli.StringFlag{
			Name:  "username, u",
			Usage: "username for harbor",
			Value: "admin",
		},
		cli.StringFlag{
			Name:  "password, p",
			Usage: "passworld for harbor",
			Value: "@dvantecH_2019",
		},
		cli.StringFlag{
			Name:  "harborurl, url",
			Usage: "harborurl for harbor",
			Value: "https://harbor.wise-paas.io/api/chartrepo/tekton/charts",
		},
	}
	app.Action = func(c *cli.Context) error {
		file := c.String("file")
		harborurl := c.String("harborurl")
		username := c.String("username")
		password := c.String("password")
		code, text := postfile(harborurl, file, username, password)
		//fmt.Println(code,text)
		fmt.Printf("code: %d \nresponse: %s\n", code,text)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
