// Pacakge w1sensor contains functons for reading the temperature from a
// one wire sensor suchs as the DS18B20 1-wire sensor.
package w1sensor

import (
  "errors"
  "fmt"
  "io/ioutil"
  "os"
  "regexp"
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

var (
  tempRegExp = regexp.MustCompile(`t=\d+`) // looking for `t=` followed by one or more digits
)

type Sensor struct {
  name string
  filePath string
}

func FirstAvailableSensor() (*Sensor, error) {
  sensors, err := AvailableSensors()
  if err != nil {
    return nil, err
  }
  if sensors != nil {
    return &sensors[0], nil
  }
  return nil, errors.New("found no sensor")
}

func AvailableSensors() ([]Sensor, error) {
  subDirectories, err := ioutil.ReadDir(BASE_DIRECTORY)
  if err != nil {
    return nil, err
  }
  sensors := make([]Sensor, 0)
  for _, dir := range subDirectories {
    name := dir.Name()
    if isSensor(name) {
      path := fmt.Sprintf("%s/%s/%s", BASE_DIRECTORY, name, SLAVE_FILE)
      sensor := Sensor{
        name: name,
        filePath: path,
      }
      sensors = append(sensors, sensor)
    }
  }
  return sensors, nil
}

// Returns the reading in celsius
func (s *Sensor) Read() (float64, error) {
  file, err := os.Open(s.filePath)
  if err != nil {
    return 0, err
  }

  bytes, err := ioutil.ReadAll(file)
  if err != nil {
    return 0, err
  }

  rawData := tempRegExp.FindString(string(bytes[:]))
  if rawData == "" {
    return 0, errors.New("Reading not found")
  }
  split := strings.Split(rawData, "=")
  file.Close()
  tempRaw, err := strconv.ParseFloat(split[1], 64)
  if err != nil {
    return 0, err
  }
  cel := tempRaw/1000
  return cel, nil
}

// Returns the reading in fah
func (s *Sensor) ReadF() (float64, error) {
  cel, err := s.Read()
  if err != nil {
    return 0, err
  }
  fah := cel*1.8 + 32
  return fah, nil
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

