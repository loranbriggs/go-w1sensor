// Pacakge w1thermsensor contains functons for reading the temperature from a
// one wire sensor suchs as the DS18B20 1-wire sensor.
package w1thermsensor

import (
  "bufio"
  "io/ioutil"
  "log"
  "os"
  "strings"
  "strconv"
)

const (
  BASE_DIRECTORY = "/sys/bus/w1/devices"
  SENSOR_DS18S20 = "10"
  SENSOR_DS1822 = "22"
  SENSOR_DS18B20 = "28"
  SLAVE_FILE = "w1_slave"
)

func Temperature() float64 {
  return ReadSensor(Devices()[0])
}

func TemperatureF() float64 {
  cel := ReadSensor(Devices()[0])
  fah := cel*1.8 + 32
  return fah
}

// reads a sensor and returns degrees in celcius
func ReadSensor(s string) float64 {
  file, err := os.Open(BASE_DIRECTORY + "/" + s + "/" + SLAVE_FILE)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines) 
  scanner.Scan() // jump to 2nd line
  scanner.Scan() // jump to 2nd line
  tempString := strings.Split(scanner.Text(), "=")[1]
  tempRaw, err := strconv.ParseFloat(tempString, 64)
  if err != nil {
    log.Fatal(err)
  }
  cel := tempRaw/1000
  //return strconv.FormatFloat(fah, 'f', 6, 64)
  return cel
}

// Returns list of Device names found
func Devices() []string {
  filesInfo, err := ioutil.ReadDir(BASE_DIRECTORY)
  if err != nil {
    log.Fatal(err)
  }
  names :=  make([]string, 0)
  for _, file := range filesInfo {
    if isSensor(file.Name()) {
      names = append(names, file.Name())
    }
  }
  return names
}

// Returns true if the string's prefix matches a sensor name
func isSensor(s string) bool {
  for _, v := range []string{SENSOR_DS18S20, SENSOR_DS1822, SENSOR_DS18B20} {
    if strings.HasPrefix(s,v) {
      return true
    }
  }
  return false
}

