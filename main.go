package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"flag"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/go-fingerprint/fingerprint"
	"github.com/go-fingerprint/gochroma"
)

var (
	maxtime = flag.Duration("maxchroma", 2*time.Minute,
		"the maximum time the server is allowed to make a chromaprint for")
	httpbind = flag.String("bind", ":6464", "the HTTP bind address")
)

func main() {
	flag.Parse()

	http.HandleFunc("/chromaprint", chromaPrintGen)

	log.Printf("Error running HTTP server %s", http.ListenAndServe(*httpbind, nil))
}

func chromaPrintGen(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(rw, "This endpoint only takes POSTs of content", http.StatusInternalServerError)
		return
	}

	var inputSamples io.Reader

	if req.URL.Query().Get("raw") == "" {
		ffmpeg := exec.Command("ffmpeg", "-v", "quiet",
			"-i", "-", "-f", "s16le", "-ac", "1", "-c:a", "pcm_s16le", "-ar", "44100", "pipe:1")
		stdin, _ := ffmpeg.StdinPipe()
		var err error
		ffmpegOut, err := ffmpeg.StdoutPipe()
		ffmpeg.Stderr = os.Stderr
		if err != nil {
			http.Error(rw, "Unable to spin up ffmpeg to decode. [P] "+err.Error(), http.StatusInternalServerError)
			return
		}
		err = ffmpeg.Start()

		if err != nil {
			http.Error(rw, "Unable to spin up ffmpeg to decode. "+err.Error(), http.StatusInternalServerError)
			return
		}

		go func() {
			io.Copy(stdin, req.Body)
			stdin.Close()
		}()

		alldata, _ := ioutil.ReadAll(ffmpegOut)
		inputSamples = bytes.NewReader(alldata)

		defer ffmpeg.Process.Wait()
		defer ffmpeg.Process.Kill()
	} else {
		inputSamples = req.Body
	}

	fpcalc := gochroma.New(gochroma.AlgorithmDefault)
	defer fpcalc.Close()

	fpoptions := fingerprint.RawInfo{
		Src:        inputSamples,
		Channels:   1,
		Rate:       44100,
		MaxSeconds: uint(maxtime.Seconds()),
	}

	if req.URL.Query().Get("png") != "" {
		fprint, err := fpcalc.RawFingerprint(fpoptions)

		if err != nil {
			http.Error(rw, "Unable to fingerprint "+err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "image/png")

		i := fingerprint.ToImage(fprint)

		if err := png.Encode(rw, i); err != nil {
			http.Error(rw, "Unable to fingerprint "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		fprint, err := fpcalc.RawFingerprint(fpoptions)

		if err != nil {
			http.Error(rw, "Unable to fingerprint "+err.Error(), http.StatusInternalServerError)
			return
		}
		barr := make([]byte, len(fprint)*4)
		by := bytes.NewBuffer(barr)
		for _, v := range fprint {
			// var i int16 = 41
			binary.Write(by, binary.LittleEndian, int32(v))
		}
		rw.Write([]byte(base64.URLEncoding.EncodeToString(by.Bytes())))
	}

}
