package genesis

import (
	"encoding/json"
	"fmt"
	"oasisTracker/smodels"
	"os"
)

const DefaultGenesisFileName = "genesis.json"

func ReadGenesisFile(fileName string) (gen smodels.GenesisDocument, err error) {

	//Use root folder as default
	file, err := os.Open(fmt.Sprint("./", fileName))
	if err != nil {
		return gen, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&gen)
	if err != nil {
		return gen, err
	}

	return gen, nil
}
