package radio

import (
	"encoding/json"
	"os"

	"gitlab.com/AgentNemo/goradios"
)

const (
	stationsJsonFilename = "stations.json"
)

func GetStations() []goradios.Station {
	stations, err := getStationsFromJson(stationsJsonFilename)
	if err != nil {
		stations = goradios.FetchStations(goradios.StationsByCodecExact, "mp3")
		writeStationsToJson(stationsJsonFilename, stations)
	}

	return stations
}

func getStationsFromJson(filename string) ([]goradios.Station, error) {
	var stations []goradios.Station
	contents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(contents, &stations)
	if err != nil {
		return nil, err
	}

	return stations, nil
}

func writeStationsToJson(filename string, stations []goradios.Station) error {
	contents, err := json.Marshal(stations)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, contents, 0666)
	if err != nil {
		return err
	}

	return nil
}
