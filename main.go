package main

import (
	"log"
	"math"
	"math/rand"
	"sync"

	"github.com/hajimehoshi/oto"
)

var argSampleRate = 11025

func main() {
	var (
		out = make(chan []uint8)
		chs = make([]chan []uint8, 3)
	)
	chs[0] = make(chan []uint8)
	chs[1] = make(chan []uint8)
	chs[2] = make(chan []uint8)
	// chs[3] = make(chan []uint8)
	go gen1(chs[0])
	go gen1(chs[1])
	go gen1(chs[2])
	// go gen1(chs[3])
	go agg(chs, out)
	play(out)
}

func gen1(ch chan []uint8) {
	var (
		argDuration = 0.3
		// argVolume     = 10000
		// argFrequency = 440.0
	)

	var (
		samplesPerSecond = uint32(argSampleRate)                           // The number of samples per second
		numberOfSamples  = uint32(float64(samplesPerSecond) * argDuration) // This is the length of the sound file in seconds
	)

	for {
		var (
			phase                     float64
			waveform                  = make([]uint8, numberOfSamples)
			frequency                 = float64(rand.Intn(300) + 100)
			frequencyRadiansPerSample = frequency * 2 * math.Pi / float64(samplesPerSecond)
		)
		for sample := uint32(0); sample < numberOfSamples; sample++ {
			phase += frequencyRadiansPerSample
			sampleValue := (math.Sin(phase) + 1.0) * 127.0
			waveform[sample] = uint8(sampleValue)
			// fmt.Print(strings.Repeat(" ", int(sampleValue/256.0*160)))
			// fmt.Println(".")
			// time.Sleep(5 * time.Millisecond)
		}
		ch <- waveform
	}
}

func agg(chs []chan []uint8, out chan []uint8) {
	for {
		var (
			wg  sync.WaitGroup
			ins = make([][]uint8, len(chs))
		)
		for i, ch := range chs {
			wg.Add(1)
			go func(ch chan []uint8, i int) {
				defer wg.Done()
				ins[i] = <-ch
			}(ch, i)
		}
		wg.Wait()
		out <- doAgg(ins)
	}
}

func doAgg(ins [][]uint8) []uint8 {
	var o = make([]uint8, len(ins[0]))
	for i := range ins[0] {
		var a int
		for j := range ins {
			a += int(ins[j][i])
		}
		o[i] = uint8(a / len(ins))
	}
	return o
}

func play(out chan []uint8) {
	p, err := oto.NewPlayer(argSampleRate, 1, 1, 256)
	if err != nil {
		log.Fatal("halp", err)
	}
	defer p.Close()
	p.SetUnderrunCallback(func() {
		log.Println("UNDERRUN, YOUR CODE IS SLOW")
	})

	for o := range out {
		p.Write(o)
	}
}
