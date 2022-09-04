package main

import (
	"fmt"
	"reflect"

	"github.com/anatol/smart.go"
	"github.com/jaypipes/ghw"
	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-smart-status",
			Short:    "S.M.A.R.T. status check for Sensu",
			Keyspace: "sensu.io/plugins/check-smart-status/config",
		},
	}

	options = []*sensu.PluginConfigOption{}
)

func main() {
	check := sensu.NewGoCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	block, err := ghw.Block()
	if err != nil {
		return sensu.CheckStateCritical, err
	}

	for _, disk := range block.Disks {
		dev, err := smart.Open("/dev/" + disk.Name)
		if err != nil {
			return sensu.CheckStateCritical, err
		}
		defer dev.Close()

		fmt.Printf("%s: ", disk.Name)
		fmt.Printf("%s, ", disk.DriveType.String())
		fmt.Printf("%s\n\n", disk.StorageController.String())
		switch sm := dev.(type) {
		case *smart.SataDevice:
			data, err := sm.ReadSMARTData()
			if err != nil {
				return sensu.CheckStateCritical, err
			}
			for _, attr := range data.Attrs {
				fmt.Printf("%s %d: ", attr.Name, attr.Type)
				if attr.Type == smart.AtaDeviceAttributeTypeTempMinMax {
					temp, _, _, _, err := attr.ParseAsTemperature()
					if err != nil {
						return sensu.CheckStateCritical, err
					}
					fmt.Printf("%d\n", temp)
				} else {
					fmt.Printf("%d\n", attr.ValueRaw)
				}
			}
		case *smart.NVMeDevice:
			data, err := sm.ReadSMART()
			if err != nil {
				return sensu.CheckStateCritical, err
			}
			v := reflect.ValueOf(data)
			vData := reflect.Indirect(v)
			for i := 0; i < vData.NumField(); i++ {
				if vData.Type().Field(i).Name == "_" {
					continue
				}
				if vData.Field(i).CanUint() {
					fmt.Printf("%s: ", vData.Type().Field(i).Name)
					fmt.Printf("%d\n", vData.Field(i).Uint())
				} else if vData.Field(i).Type().String() == "smart.Uint128" {
					fmt.Printf("%s: ", vData.Type().Field(i).Name)
					fmt.Printf("%d\n", vData.Field(i).FieldByName("Val").Index(0).Uint())
				} else {
					continue
				}
			}
		case *smart.ScsiDevice:
			continue
		}
	}

	return sensu.CheckStateOK, nil
}
