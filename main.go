package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Metar struct {
	City            string
	ObserveTime     string
	WindDirection   string
	WindSpeed       string
	WindOffset      string
	MeteoVisibility string
	Weather         string
	Cloudness       string
	Temperature     string
	DewPoint        string
	AtmoPressure    string
	Visibility      string
	LandingForecast string
	Remark          string
	AirportPressure string
}

var metar Metar
var metar_slice []string
var metar_size int

func decodeCity(code string) string {
	var result string
	airports := map[string]string{
		"UEEE": "Якутск",
	}
	for key, value := range airports {
		if key == code {
			result = value
		}
	}
	return result
}

func clarifyCity(data []string) {
	var trigger bool = false
	city := decodeCity(data[0])
	if len(city) > 0 {
		metar_slice = append(metar_slice, city)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func decodeObserveTime(code string) string {
	var tmp string
	if strings.Contains(code, "Z") {
		tmp = strings.TrimRight(code, "Z")
	}
	return tmp
}

func clarifyObserveTime(data []string) {
	var trigger bool = false
	time := decodeObserveTime(data[1])
	if len(time) > 0 {
		metar_slice = append(metar_slice, time)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func decodeWind(code string) (string, string) {
	var result string
	var direction string
	var speed string
	if strings.Contains(code, "MPS") {
		result = strings.TrimRight(code, "MPS")
		direction = result[:3]
		speed = result[3:]
	}
	return direction, speed
}

func clarifyWind(data []string) {
	var trigger bool = false
	direction, speed := decodeWind(data[2])
	if len(direction) > 0 {
		metar_slice = append(metar_slice, direction)
		metar_slice = append(metar_slice, speed)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
		metar_slice = append(metar_slice, "None")
	}
}

func decodeWindOffset(code string) string {
	var result string
	if strings.Contains(code, "V") {
		tmp := strings.Split(code, "V")
		tmp2 := strings.Join(tmp, "-")
		result = tmp2
	}
	return result
}

// начиная отсюда количество элементов слайса варьируется
func clarifyWindOffset(data []string) {
	var trigger bool = false
	if metar_size == 13 {
		offset := decodeWindOffset(data[3])
		if len(offset) > 0 {
			metar_slice = append(metar_slice, offset)
			trigger = true
		}
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func decodeMeteoVisibility(code string) string {
	var result string
	if _, err := strconv.Atoi(code); err == nil {
		result = code
	}
	return result
}

func clarifyMeteoVisibility(data []string) {
	var trigger bool = false
	var visibility string
	if metar_size == 11 {
		visibility = decodeMeteoVisibility(data[3])
	}
	if metar_size == 12 {
		visibility = decodeMeteoVisibility(data[3])
	}
	if metar_size == 13 {
		visibility = decodeMeteoVisibility(data[4])
	}
	if len(visibility) > 0 {
		metar_slice = append(metar_slice, visibility)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func decodeWeather(code string) string {
	var result string
	weather := map[string]string{
		"DZ":   "морось",
		"RA":   "rain",
		"SN":   "snow",
		"SG":   "снежные зерна",
		"PL":   "ледяная крупа",
		"GS":   "мелкий град или ледяная крупа",
		"GR":   "град",
		"RASN": "дождь со снегом",
		"SNRA": "снег с дождем",
		"SHSN": "ливневый снег",
		"SHRA": "ливневый дождь",
		"SHGR": "град",
		"FZRA": "переохлажденный дождь",
		"FZDZ": "переохлажденная морось",
		"TSRA": "гроза с дождем",
		"TSGR": "гроза с градом",
		"TSGS": "гроза со снежной крупой",
		"TSSN": "гроза со снегом",
		"DS":   "пыльная буря",
		"SS":   "песчаная буря",
		"FG":   "туман",
		"FZFG": "переохлажденный туман",
		"VCFG": "туман в окрестностях аэродрома",
		"MIFG": "поземный туман",
		"PRFG": "аэродром частично покрыт туманом",
		"BCFG": "туман местами",
		"BR":   "дымка",
		"HZ":   "мгла",
		"FU":   "дым",
		"DRSN": "снежный поземок",
		"DRSA": "песчаный поземок",
		"DRDU": "пыльный поземок",
		"DU":   "пыль в воздухе",
		"BLSN": "снежная низовая метель",
		"BLDU": "пыльная низовая метель",
		"SQ":   "шквал",
		"IC":   "ледяные иглы",
		"TS":   "гроза",
		"VCTS": "гроза в окрестности",
		"VA":   "вулканический пепел",
	}
	for key, value := range weather {
		if strings.Contains(code, key) {
			result = value
		}
	}
	return result
}

func clarifyWeather(data []string) {
	var trigger bool = false
	var weather string
	if metar_size == 12 {
		weather = decodeWeather(data[4])
	}
	if metar_size == 13 {
		weather = decodeWeather(data[5])
	}
	if len(weather) > 0 {
		metar_slice = append(metar_slice, weather)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func decodeCloudness(code string) string {
	var result string
	cloudness := map[string]string{
		"SKC":   "ясно",
		"NSC":   "несущественная облачность",
		"FEW":   "местами облачно",
		"SCT":   "небольшая облачность",
		"BKN":   "облачно",
		"OVC":   "сплошная облачность",
		"CAVOK": "ясно",
	}
	for key, value := range cloudness {
		if strings.Contains(code, key) {
			result = value
		}
	}
	return result
}

func clarifyCloudness(data []string) {
	var trigger bool = false
	var cloudness string
	if metar_size <= 10 {
		cloudness = decodeCloudness(data[3])
	}
	if metar_size == 11 {
		cloudness = decodeCloudness(data[4])
	}
	if metar_size == 12 {
		cloudness = decodeCloudness(data[5])
	}
	if metar_size == 13 {
		cloudness = decodeCloudness(data[6])
	}
	if len(cloudness) > 0 {
		metar_slice = append(metar_slice, cloudness)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func decodeTemperature(code string) (string, string) {
	var tmp []string
	tmp = strings.Split(code, "/")
	temp := tmp[0]
	point := tmp[1]
	return temp, point
}

func clarifyTemperature(data []string) {
	var trigger bool = false
	var temperature string
	var dew_point string
	if metar_size <= 10 {
		temperature, dew_point = decodeTemperature(data[4])
	}
	if metar_size == 11 {
		temperature, dew_point = decodeTemperature(data[5])
	}
	if metar_size == 12 {
		temperature, dew_point = decodeTemperature(data[6])
	}
	if metar_size == 13 {
		temperature, dew_point = decodeTemperature(data[7])
	}
	if len(temperature) > 0 {
		metar_slice = append(metar_slice, temperature)
		metar_slice = append(metar_slice, dew_point)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
		metar_slice = append(metar_slice, "None")
	}
}

func decodeAtmoPressure(code string) string {
	var result string
	result = strings.TrimLeft(code, "Q")
	return result
}

func clarifyAtmoPressure(data []string) {
	var trigger bool = false
	var pressure string
	if metar_size <= 10 {
		pressure = decodeAtmoPressure(data[5])
	}
	if metar_size == 11 {
		pressure = decodeAtmoPressure(data[6])
	}
	if metar_size == 12 {
		pressure = decodeAtmoPressure(data[7])
	}
	if metar_size == 13 {
		pressure = decodeAtmoPressure(data[8])
	}
	if len(pressure) > 0 {
		metar_slice = append(metar_slice, pressure)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func decodeVisibility(code string) string {
	var tmp []string
	if strings.Contains(code, "//") {
		tmp = strings.Split(code, "//")
		return tmp[1]
	}
	return "None"
}

func clarifyVisibility(data []string) {
	var trigger bool = false
	var visibility string
	if metar_size == 10 {
		visibility = decodeVisibility(data[6])
	}
	if metar_size == 11 {
		visibility = decodeVisibility(data[7])
	}
	if metar_size == 12 {
		visibility = decodeVisibility(data[8])
	}
	if metar_size == 13 {
		visibility = decodeVisibility(data[9])
	}
	if len(visibility) > 0 {
		metar_slice = append(metar_slice, visibility)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, visibility)
	}
}

func decodeLandingForecast(code string) string {
	var result string
	cloudness := map[string]string{
		"NOSIG": "в ближайшие 2 часа без изменений",
		"BECMG": "ожидаются значительные изменения метеоусловий",
		"TEMPO": "ожидаются временные изменения метеоусловий",
	}
	for key, value := range cloudness {
		if key == code {
			result = value
		}
	}
	return result
}

func clarifyLandingForecast(data []string) {
	var trigger bool = false
	var forecast string
	if metar_size == 10 {
		forecast = decodeLandingForecast(data[7])
	}
	if metar_size == 11 {
		forecast = decodeLandingForecast(data[8])
	}
	if metar_size == 12 {
		forecast = decodeLandingForecast(data[9])
	}
	if metar_size == 13 {
		forecast = decodeLandingForecast(data[10])
	}
	if len(forecast) > 0 {
		metar_slice = append(metar_slice, forecast)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func decodeRemark(code string) string {
	var result string
	if strings.Contains(code, "RMK") {
		result = code
	}
	return result
}

func clarifyRemark(data []string) {
	var trigger bool = false
	var remark string
	if metar_size == 9 {
		remark = decodeRemark(data[7])
	}
	if metar_size == 10 {
		remark = decodeRemark(data[8])
	}
	if metar_size == 11 {
		remark = decodeRemark(data[9])
	}
	if metar_size == 12 {
		remark = decodeRemark(data[10])
	}
	if metar_size == 13 {
		remark = decodeRemark(data[11])
	}
	if len(remark) > 0 {
		metar_slice = append(metar_slice, remark)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func decodeAirportPressure(code string) string {
	tmp := strings.TrimLeft(code, "QFE")
	return tmp
}

func clarifyAiportPressure(data []string) {
	var trigger bool = false
	var pressure string
	if metar_size == 9 {
		pressure = decodeAirportPressure(data[8])
	}
	if metar_size == 10 {
		pressure = decodeAirportPressure(data[9])
	}
	if metar_size == 11 {
		pressure = decodeAirportPressure(data[10])
	}
	if metar_size == 12 {
		pressure = decodeAirportPressure(data[11])
	}
	if metar_size == 13 {
		pressure = decodeAirportPressure(data[12])
	}
	if len(pressure) > 0 {
		metar_slice = append(metar_slice, pressure)
		trigger = true
	}
	if trigger == false {
		metar_slice = append(metar_slice, "None")
	}
}

func parseMetar(url string) []string {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	data := string(responseData)
	array := strings.Split(data, "\n")
	joined_array := strings.Join(array[:], " ")
	joined_array = strings.TrimRight(joined_array, " ")
	parsed_metar := strings.Split(joined_array, " ")[2:]

	return parsed_metar
}

func makeObject() Metar {
	object := Metar{
		City:            metar_slice[0],
		ObserveTime:     metar_slice[1],
		WindDirection:   metar_slice[2],
		WindSpeed:       metar_slice[3],
		WindOffset:      metar_slice[4],
		MeteoVisibility: metar_slice[5],
		Weather:         metar_slice[6],
		Cloudness:       metar_slice[7],
		Temperature:     metar_slice[8],
		DewPoint:        metar_slice[9],
		AtmoPressure:    metar_slice[10],
		Visibility:      metar_slice[11],
		LandingForecast: metar_slice[12],
		Remark:          metar_slice[13],
		AirportPressure: metar_slice[14],
	}
	return object
}

func exportJson(object Metar, dir string) {
	file, _ := json.MarshalIndent(object, "", " ")
	_ = ioutil.WriteFile(dir+"/metar.json", file, 0644)
}

func main() {
	/*
		ObserveTime: ddhhmm
		WindDirection: градус
		WindSpeed: метр в секунду
		WindOffset: диапазон смещения направления ветра (полагаю в градусах)
		MeteoVisibility: метеорологическая видимость (на высоте от 1.5км) в метрах
		Temperature: градус в сельциях
		DewPoint: точка росы в сельциях
		AtmoPressure: атмосферное давление в гектопаскалях
		Visibility: видимость на взлетно-посадочной полосе в метрах
		Remark: ремарка (RMK) после которой идут дополнительный данные
		AirportPressure: давление на взлетно-посадочной полосе (единица ртутного столба)
	*/
	verbPtr := flag.Bool("v", false, "verbose mode")
	flag.Parse()
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	metar_array := parseMetar("https://tgftp.nws.noaa.gov/data/observations/metar/stations/UEEE.TXT")
	metar_size = len(metar_array)
	clarifyCity(metar_array)
	clarifyObserveTime(metar_array)
	clarifyWind(metar_array)
	clarifyWindOffset(metar_array)
	clarifyMeteoVisibility(metar_array)
	clarifyWeather(metar_array)
	clarifyCloudness(metar_array)
	clarifyTemperature(metar_array)
	clarifyAtmoPressure(metar_array)
	clarifyVisibility(metar_array)
	clarifyLandingForecast(metar_array)
	clarifyRemark(metar_array)
	clarifyAiportPressure(metar_array)
	if *verbPtr {
		fmt.Println("Original data:", metar_array)
		fmt.Println("Processed data:", metar_slice)
	}
	object := makeObject()
	exportJson(object, dir)
}
