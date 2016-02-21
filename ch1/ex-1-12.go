package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
)

var palette = []color.Color{color.White, color.Black}

const (
	whiteIndex = 0 // first color in palette
	blackIndex = 1 // next color in palette
)

func main() {
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println("Request")
	params := make(map[string]int)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		val, err := strconv.Atoi(v[0])
		if err == nil {
			params[k] = val
		}
	}
	lissajous(w, params)
}

func get(m map[string]int, key string, defaultValue int) int {
	value, ok := m[key]
	if ok {
		return value
	} else {
		return defaultValue
	}
}

func lissajous(out io.Writer, params map[string]int) {
	const (
		res = 0.001 // angular resolution
	)

	delay := get(params, "delay", 8)            // delay between frames in 10ms units
	size := get(params, "size", 100)            // image canvas covers [-size..+size]
	nframes := get(params, "nframes", 64)       // number of animation frames
	cycles := float64(get(params, "cycles", 5)) // number of complete x oscillator revolutions

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	sizef := float64(size)
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*sizef+0.5), size+int(y*sizef+0.5),
				blackIndex)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
}
