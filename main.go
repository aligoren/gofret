package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

/*
	usage = usage text
	version = current number
	help use Sprintf
	*cliUrl from cmd
	*cliVersion from cmd
	*cliHelp * from cmd
*/
var (
	usage      = "Usage: ./gofret -url=http://some/do.zip"
	version    = "Version: 0.1"
	help       = fmt.Sprintf("\n\n  %s\n\n\n  %s", usage, version)
	cliUrl     *string
	cliVersion *bool
	cliHelp    *bool
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
}

func main() {

	/*
		Parse all flags
	*/
	flag.Parse()

	if *cliUrl != "" {
		fmt.Println("Downloading file")

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
		fmt.Println(response.Status) // Example: 200 OK

		/*
			fileSize example: 12572 bytes
		*/
		fileSize, err := io.Copy(file, response.Body)

		if err != nil {
			panic(err)
		}

		fmt.Printf("%s with %v bytes downloaded", fileName, fileSize)

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
	} else {
		/*
			using help's usage text for handling other status
		*/
		fmt.Println(flag.Lookup("help").Usage)
	}
}
