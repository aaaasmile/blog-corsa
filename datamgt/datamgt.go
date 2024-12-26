package datamgt

import (
	"corsa-blog/idl"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type DataMgt struct {
	datafileName string
	items        []idl.CmtItem
}

func (dmgt *DataMgt) ReadData() error {
	fname := dmgt.datafileName
	log.Println("load json data ", fname)
	if fname == "" {
		return fmt.Errorf("data file is empty")
	}
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&dmgt.items)
	if err != nil {
		return err
	}
	log.Println("Loaded scheduler from file ", fname, len(dmgt.items))
	return nil
}
