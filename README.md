# IOT Brewery

The iot brewery 

## Outline

iot brewery is made of a webservice and an iotcontroller and an electron desktop dashboard

### webservice

A simple rest api webservice written on go. It stores iot data from the iotcontroller client (see iotcontroller). The data is then stored in a mongodb database. 
The data can then be queried from any device, i.e mobile or more specifically the electron dashboard that views all data in  a graphical easy to read format 
It also contains a Dockerfile based on a golang 1.9

Run this command to test 

```bash

docker run -it   -p 9001:9001  --network="host" <image-id>

```

### iotcontroller
Interfaces to the webservice using a simple rest post call. The iotcontroller uses an config file that gets set and via the electron dashboard. The config has a cron
setting to read the custom designed hardware (8 to 16 channels of analog data) and then do some specific convertion (temperature in this case). Depending
on the upper and lower limit settings the iotcontroller will the send data to 2X8 bit digital io expander that in turn can be used to control relays etc.

### electron-iot-desktop
This is a simple electron desktop app that reads a config.json file to visualize and read data from the webservice for graphical viewing and monitoring
