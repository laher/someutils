package someutils

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func init() {
	Register(Util{
		"wget",
		Wget})
}

type WgetOptions struct {
	IsContinue *bool
	Filename *string
}

const (
	FILEMODE os.FileMode = 0660
)

func Wget(call []string) error {

	options := WgetOptions{}
	flagSet := flag.NewFlagSet("wget", flag.ContinueOnError)
	options.IsContinue = flagSet.Bool("c", false, "continue")
	options.Filename = flagSet.String("o", "", "output filename")
	helpFlag := flagSet.Bool("help", false, "Show this help")

	err := flagSet.Parse(splitSingleHyphenOpts(call[1:]))
	if err != nil {
		return err
	}
	if *helpFlag {
		println("wget [options] URL")
		flagSet.PrintDefaults()
		return nil
	}
	args := flagSet.Args()
	if len(args) < 1 {
		flagSet.PrintDefaults()
		return errors.New("Not enough args")
	}
	if len(args) > 0 {
		links := args
		return wget(links, options)
	} else {
		if IsPipingStdin() {
			//check STDIN
			return wget([]string{}, options)
		} else {
			//NOT piping.
			return errors.New("Not enough args")
		}
	}
}

func wget(links []string, options WgetOptions) error {
	for _, link := range links {
		err := wgetOne(link, options)
		if err != nil {
			return err
		}
	}
	return nil
}

func tidyFilename(filename string) string {
	//invalid filenames ...
	if filename == "" || filename == "/" ||filename == "\\" || filename == "." {
		filename = "index"
	}
	return filename
}

func wgetOne(link string, options WgetOptions) error {
	if !strings.Contains(link, ":") {
		link = "http://" + link
	}
	startTime := time.Now()
	request, err := http.NewRequest("GET", link, nil)
	//resp, err := http.Get(link)
	if err != nil {
		return err
	}

	filename := ""
	if *options.Filename != "" {
		filename = *options.Filename
	}
	client := &http.Client{}
	//continue from where we left off ...
	if *options.IsContinue {
		if filename == "" {
			filename = filepath.Base(request.URL.Path)
			filename = tidyFilename(filename)
			if !strings.Contains(filename, ".") {
				filename = filename + ".html"
			}
		}
		fi, err := os.Stat(filename)
		if err != nil {
			return err
		}
		from := fi.Size()
		headRequest, err := http.NewRequest("HEAD", link, nil)
		if err != nil {
			return err
		}
		headResp, err := client.Do(headRequest)
		if err != nil {
			return err
		}
		cl := headResp.Header.Get("Content-Length")
		if cl != "" {
		rangeHeader := fmt.Sprintf("bytes %d-%s", from, cl)
		fmt.Printf("Adding range header: %s\n", rangeHeader)
		request.Header.Add("Range", rangeHeader)
		} else {
			fmt.Println("Could not find file length using HEAD request")
		}
	}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("Http response: %s\n", resp.Status)
	
	lenS := resp.Header.Get("Content-Length")
	len := int64(-1)
	if lenS != "" {
		len, err = strconv.ParseInt(lenS, 10, 32)
		if err != nil {
			return err
		}
	}
	typ := resp.Header.Get("Content-Type")
	fmt.Printf("Length: %v [%s]\n", len, typ)
	
	defer resp.Body.Close()
	if filename == "" {	
		filename, err = getFilename(resp)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Saving to: '%v'\n\n", filename)
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	
	buf := make([]byte, 4068)
	tot := int64(0)
	i := 0
	
	for {
        // read a chunk
        n, err := resp.Body.Read(buf)
        if err != nil && err != io.EOF { 
			return err
		}
        if n == 0 { break }
		tot += int64(n)

        // write a chunk
        if _, err := out.Write(buf[:n]); err != nil {
            return err
        }
		i+=1
		if len > -1 {
			if len < 1 {
				fmt.Printf("\r     [ <=>                                  ] %d\t-.--KB/s eta ?s             ", tot)
			} else {
				//show percentage
				perc := (100 * tot) / len
				prog := progress(perc)
				nowTime := time.Now()
				totTime := nowTime.Sub(startTime)
				spd := float64(tot / 1000) / totTime.Seconds()
				remKb := float64(len - tot) / float64(1000)
				eta :=  remKb / spd
				fmt.Printf("\r%3d%% [%s] %d\t%0.2fKB/s eta %0.1fs             ", perc, prog, tot, spd, eta)
			}
		} else {
			//show dots
			if math.Mod(float64(i), 20) == 0 {
				fmt.Print(".")
			}
		}
    }
	perc := (100 * tot) / len
	prog := progress(perc)
	nowTime := time.Now()
	totTime := nowTime.Sub(startTime)
	spd := float64(tot / 1000) / totTime.Seconds()
	if len < 1 {
		fmt.Printf("\r     [ <=>                                  ] %d\t-.--KB/s in %0.1fs             ", tot, totTime.Seconds())
		fmt.Printf("\n (%0.2fKB/s) - '%v' saved [%v]\n", spd, filename, tot)
	} else {
		fmt.Printf("\r%3d%% [%s] %d\t%0.2fKB/s in %0.1fs             ", perc, prog, tot, spd, totTime.Seconds())
		fmt.Printf("\n '%v' saved [%v/%v]\n", filename, tot, len)
	}
	if err != nil {
		return err
	}
	err = out.Close()
	return err
}

func progress(perc int64) string {
	equalses := perc * 38 / 100 
	if equalses < 0 {
		equalses = 0
	}
	spaces := 38 - equalses
	if spaces < 0 {
		spaces = 0
	}
	prog := strings.Repeat("=", int(equalses)) + ">" + strings.Repeat(" ", int(spaces))
	return prog 
}

func getFilename(resp *http.Response) (string, error) {
	filename := filepath.Base(resp.Request.URL.Path)
	filename = tidyFilename(filename)

	if !strings.Contains(filename, ".") {
		ct := resp.Header.Get("Content-Type")
		//println(ct)
		ext := "htm"
		mediatype, _, err := mime.ParseMediaType(ct)
		if err != nil {
			fmt.Printf("mime error: %v\n", err)
		} else {
			fmt.Printf("mime type: %v (from Content-Type %v)\n", mediatype, ct)
			slash := strings.Index(mediatype, "/")
			if slash != -1 {
				_, sub := mediatype[:slash], mediatype[slash+1:]
				if sub != "" {
					ext = sub
				}
			}
		}
		filename = filename + "." + ext
	}
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return filename, nil
		} else {
			return "", err
		}
	} else {
		num := 1
		//just stop after 100
		for num < 100 {
			filenameNew := filename + "." + strconv.Itoa(num)
			_, err := os.Stat(filenameNew)
			if err != nil {
				if os.IsNotExist(err) {
					return filenameNew, nil
				} else {
					return "", err
				}
			}
			num += 1
		}
		return filename, errors.New("Stopping after trying 100 filename variants")
	}
}
