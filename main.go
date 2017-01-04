package main

import (
	"bytes"
	"flag"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-fingerprint/fingerprint"
	"github.com/go-fingerprint/gochroma"
)

func main() {
	httpbind := flag.String("bind", ":6464", "the HTTP bind address")
	flag.Parse()

	http.HandleFunc("/chromaprint", chromaPrintGen)

	log.Printf("Error running HTTP server %s", http.ListenAndServe(*httpbind, nil))
	//

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

	if req.URL.Query().Get("png") != "" {
		fprint, err := fpcalc.RawFingerprint(
			fingerprint.RawInfo{
				Src:        inputSamples,
				Channels:   1,
				Rate:       44100,
				MaxSeconds: 120,
			})

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
		fprint, err := fpcalc.Fingerprint(
			fingerprint.RawInfo{
				Src:        inputSamples,
				Channels:   1,
				Rate:       44100,
				MaxSeconds: 120,
			})

		if err != nil {
			http.Error(rw, "Unable to fingerprint "+err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Write([]byte(fprint))
	}

}
