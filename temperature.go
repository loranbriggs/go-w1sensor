// Pacakge w1thermsensor contains functons for reading the temperature from a
// one wire sensor suchs as the DS18B20 1-wire sensor.
package w1thermsensor

import (
  "bufio"
  "io/ioutil"
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

func Temperatrue() string {
  return ReadSensor(Devices()[0])
}

// reads a sensor
func ReadSensor(s string) string {
  file, _ := os.Open(BASE_DIRECTORY + "/" + s + "/" + SLAVE_FILE)
  defer file.Close()
  scanner := bufio.NewScanner(file)
  scanner.Split(bufio.ScanLines) 
  scanner.Scan() // jump to 2nd line
  scanner.Scan() // jump to 2nd line
  tempS := strings.Split(scanner.Text(), "=")[1]
  temp, _ := strconv.ParseFloat(tempS, 64)
  cel := temp/1000
  fah := cel*1.8 + 32
  return strconv.FormatFloat(fah, 'f', 6, 64)
}

// Returns list of Device names
func Devices() []string {
  filesInfo, _ := ioutil.ReadDir(BASE_DIRECTORY)
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

