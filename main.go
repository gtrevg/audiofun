package main

import (
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/oto"
)

func main() {
	var (
		argSampleRate = 11025
		argDuration   = 0.2
		// argVolume     = 10000
		// argFrequency = 440.0
	)
	p, err := oto.NewPlayer(argSampleRate, 1, 1, 256)
	if err != nil {
		log.Fatal("halp", err)
	}
	defer p.Close()
	p.SetUnderrunCallback(func() {
		log.Println("UNDERRUN, YOUR CODE IS SLOW")
	})

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

		p.Write(waveform)
	}
}
