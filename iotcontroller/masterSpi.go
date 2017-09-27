package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecheney/gpio"
	"github.com/davecheney/gpio/rpi"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"time"
	"strconv"
)

var (
	dataIn, dataOut, dataClk, dataReset gpio.Pin
	delay                               time.Duration = 10
	err                               error
	limits [][]float32
)

type IOTData struct {
	BoardId          int       `json:"BoardId"`
	Description string    `json:"description"`
	Temperature []float32 `json:"temperature"`
	DigitalIOA  uint8     `json:"digitalIOA"`
	DigitalIOB  uint8     `json:"digitalIOB"`
	Time        time.Time `json:"timestamp"`
}

func initSpi() {
	// set the spi pins and mode
	dataIn, err = gpio.OpenPin(rpi.GPIO23, gpio.ModeInput)     // miso
	dataOut, err = gpio.OpenPin(rpi.GPIO27, gpio.ModeOutput)    // mosi
	dataClk, err = gpio.OpenPin(rpi.GPIO22, gpio.ModeOutput)   // scl

	if err != nil {
		logger.Error(fmt.Sprintf("Error opening pin! %s\n", err))
		return
	}

	dataOut.Clear()
	dataClk.Clear()

	logger.Info("GPIO data pins initialised ")
}

func writeSpi(data uint8) {
	dataOut.Clear()
	var mask uint8 = 0x80

	for mask != 0 {
		if data&mask == mask {
			dataOut.Set()
		} else {
			dataOut.Clear()
		}
		time.Sleep(delay * time.Microsecond)
		dataClk.Set()
		time.Sleep(delay * time.Microsecond)
		dataClk.Clear()
		time.Sleep(delay * time.Microsecond)
		mask = mask >> 1
		dataOut.Clear()
		time.Sleep(delay * time.Microsecond)
	}
}

func readSpi() uint8 {
	var ret, mask uint8 = 0x00, 0x80
	for mask != 0 {
		dataClk.Clear()
		time.Sleep(delay * time.Microsecond)
		dataClk.Set()
		time.Sleep(delay * time.Microsecond)
		if dataIn.Get() {
			ret |= mask
		}
		dataClk.Clear()
		time.Sleep(delay * time.Microsecond)
		mask = mask >> 1
	}
	return ret
}

func checkError(err error) {
	if err != nil {
		logger.Error(fmt.Sprintf("%s ", err.Error()))
	}
}

func processIOT(tr *http.Transport, config Config) {
	t := []float32{0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0}
	var ioChannelA, ioChannelB uint8 = 0x00,0x00

	logger.Level = config.Level

		// get all hi and low bytes from adc
		writeSpi(0x01)
		time.Sleep(1 * time.Millisecond)
		for i := 0; i < 8; i++ {
			readLo := readSpi()
			time.Sleep(2 * time.Millisecond)
			readHi := readSpi()
			time.Sleep(2 * time.Millisecond)
			logger.Debug(fmt.Sprintf("masterSPI adc lo byte from IOTboard channel : %d value : %d", i, readLo))
			logger.Debug(fmt.Sprintf("masterSPI adc hi byte from IOTboard channel : %d value : %d", i, readHi))

			// calculate temperature and compare against hi and lo
			result := uint16((256 * uint16(readHi)) + uint16(readLo))

			lower,_ := strconv.ParseFloat(config.Limits[i].Lower, 32)
			upper,_ := strconv.ParseFloat(config.Limits[i].Upper, 32)
			logger.Debug(fmt.Sprintf("masterSPI temperature lower limit channel : %d value : %f", i, lower))
			logger.Debug(fmt.Sprintf("masterSPI temperature upper limit channel : %d value : %f", i, upper))
			if result != 0 {
				logger.Debug(fmt.Sprintf("masterSPI adc byte from IOTboard channel : %d value : %d", i, result))
				t[i] = (float32((1100.0/1023.0)*float32(result)) - 500.0) / 10.0
				logger.Debug(fmt.Sprintf("masterSPI temperature calc channel : %d value : %f", i, t[i]))
				// add hysterisis
				if t[i] > float32(upper) {
					ioChannelA |= (1 << uint8(i))
				}
				if t[i] < float32(lower) {
					ioChannelB |= (1 << uint8(i))
				}
			}
		}

		time.Sleep(2 * time.Millisecond)
		writeSpi(ioChannelA);
		time.Sleep(2 * time.Millisecond)
		writeSpi(ioChannelB);
		time.Sleep(2 * time.Millisecond)

		
		logger.Debug(fmt.Sprintf("masterSPI sent ioChannel A data : %d B data : %d", ioChannelA, ioChannelB))
		iotdata := &IOTData{BoardId: 1, Description: "Board IOT", Temperature: t, DigitalIOA: ioChannelA, DigitalIOB: ioChannelB, Time: time.Now()}
		b, _ := json.Marshal(iotdata)
		logger.Trace(string(b))

		req, err := http.NewRequest("POST", config.Url + "?apikey=" + config.Apikey, bytes.NewBuffer(b))
		req.Header.Set("X-Custom-Header", config.Apikey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		logger.Trace("response Status:" + resp.Status)
		logger.Trace("response Headers: ")
		for k, _ := range resp.Header {
			logger.Trace("  key:" + string(k))
		}
		body, _ := ioutil.ReadAll(resp.Body)
		logger.Info("response Body:" + string(body))
		time.Sleep(1000 * time.Millisecond)
}

func cleanup(c *cron.Cron) {
	logger.Info("\nClearing and unexporting pins.\n")
	dataOut.Clear()
	dataClk.Clear()

	logger.Warn("Cleanup resources")
	logger.Info("Terminating")
	c.Stop()
}
