package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/cheggaaa/pb"
	"time"
)

/*
	usage = usage text
	version = current number
	help use Sprintf
	about (me :))
	*cliUrl from cmd
	*cliVersion from cmd
	*cliHelp from cmd
	*cliAbout from cmd
*/
var (
	usage      = "Usage: ./gofret -url=http://some/do.zip"
	version    = "Version: 0.1"
	about 	   = "Coded by Ali GOREN, aligoren.com"
	help       = fmt.Sprintf("\n\n  %s\n\n\n  %s\n\n\n  %s", usage, version, about)
	cliUrl     *string
	cliVersion *bool
	cliHelp    *bool
	cliAbout   *bool
)

func init() {
	/*
		if *cliUrl != "" {
			fmt.Println(*cliUrl)
		}

		./gofret -url=http://somesite.com/somefile.zip
		./gofret -url=https://github.com/aligoren/syspy/archive/master.zip
	*/
	cliUrl = flag.String("url", "", usage)

	/*
		else if *cliVersion{
			fmt.Println(flag.Lookup("version").Usage)
		}

		./gofret -version
	*/
	cliVersion = flag.Bool("version", false, version)

	/*
		if *cliHelp {
			fmt.Println(flag.Lookup("help").Usage)
		}

		./gofret -help
	*/
	cliHelp = flag.Bool("help", false, help)

	/*
		if *cliAbout {
			fmt.Println(flag.Lookup("about").Usage)
		}
	*/
	cliAbout = flag.Bool("about", false, about)
}

func main() {

	/*
		Parse all flags
	*/
	flag.Parse()

	if *cliUrl != "" {
		fmt.Println("\nDownloading file...\n")

		/* parse url from *cliUrl */
		fileUrl, err := url.Parse(*cliUrl)

		if err != nil {
			panic(err)
		}

		/* get path from *cliUrl */
		filePath := fileUrl.Path

		/*
			seperate file.
			http://+site.com/+(file.zip)
		*/
		segments := strings.Split(filePath, "/")

		/*
			file.zip filename lenth -1
		*/
		fileName := segments[len(segments)-1]

		/*
			Create new file.
			Filename from fileName variable
		*/
		file, err := os.Create(fileName)

		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer file.Close()

		/*
			check status and CheckRedirect
		*/
		checkStatus := http.Client{
			CheckRedirect: func(r *http.Request, via []*http.Request) error {
				r.URL.Opaque = r.URL.Path
				return nil
			},
		}

		/*
			Get Response: 200 OK?
		*/
		response, err := checkStatus.Get(*cliUrl)

		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		defer response.Body.Close()
		fmt.Printf("Request Status: %s\n\n", response.Status) // Example: 200 OK

		/*
			filesize example: 12572 bytes
		*/
		filesize := response.ContentLength
		/*
			go func == goroutines
			more on my topic: http://stackoverflow.com/a/30534837/3821823
		*/
		go func() {
	        n, err := io.Copy(file, response.Body)
	        if n != filesize {
	            fmt.Println("Truncated")
	        }
	        if err != nil {
	            fmt.Printf("Error: %v", err)
	        }
    	}()

    	/*
    		countSize == fileSize := 20000 / 1000
    		=> 20
    	*/
		countSize := int(filesize / 1000)
		bar := pb.StartNew(countSize) // start new progressbar
		var fi os.FileInfo // get file information from os
		for fi == nil || fi.Size() < filesize { // for like while
			fi, _ = file.Stat() // File status
			bar.Set(int(fi.Size() / 1000)) // File size / 1000
			time.Sleep(time.Millisecond) // wait millisecond
		}
		finishMessage := fmt.Sprintf("\n%s with %v bytes downloaded",
		 fileName, filesize)
		bar.FinishPrint(finishMessage) // finished messages

		if err != nil {
			panic(err)
		}

	} else if *cliVersion {
		/*
			lookup version flag's usage text
		*/
		fmt.Println(flag.Lookup("version").Usage)
	} else if *cliHelp {
		/*
			lookup help flag's usage text
		*/
		fmt.Println(flag.Lookup("help").Usage)
	} else if *cliAbout{
		/*
			lookup about flag's usage text
		*/
		fmt.Println("\n\n"+flag.Lookup("about").Usage)
	} else {
		/*
			using help's usage text for handling other status
		*/
		fmt.Println(flag.Lookup("help").Usage)
	}
}