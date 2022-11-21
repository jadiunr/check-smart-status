package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/anatol/smart.go"
	"github.com/iancoleman/strcase"
	"github.com/jaypipes/ghw"
	"github.com/sensu/sensu-go/types"
	"github.com/sensu/sensu-plugin-sdk/sensu"
)

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	
}

type InfluxDBLine struct {
	Metrics map[string][]*InfluxDBMetric
}

type InfluxDBMetric struct {
	Tags map[string]string
	Value float64
}

func (s InfluxDBLine) generateInfluxDBLine() {
	var metricNamesSorted []string
	for key, _ := range s.Metrics {
		metricNamesSorted = append(metricNamesSorted, key)
	}
	sort.Strings(metricNamesSorted)

	for _, metricName := range metricNamesSorted {
		for _, metric := range s.Metrics[metricName] {
			var tags []string
			for tagName, tagValue := range metric.Tags {
				tags = append(tags, fmt.Sprintf("%s=%s", strcase.ToSnake(tagName), tagValue))
			}
			sort.Strings(tags)
			fmt.Printf("smart_%s,%s value=%.2f\n", strcase.ToSnake(metricName), strings.Join(tags, ","), metric.Value)
		}
	}
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
	influxDBLine := InfluxDBLine{
		Metrics: make(map[string][]*InfluxDBMetric),
	}
	block, err := ghw.Block()
	if err != nil {
		return sensu.CheckStateCritical, err
	}

	for _, disk := range block.Disks {
		dev, err := smart.Open("/dev/" + disk.Name)
		if err != nil {
			continue
		}
		defer dev.Close()
		switch sm := dev.(type) {
		case *smart.SataDevice:
			data, err := sm.ReadSMARTData()
			if err != nil {
				return sensu.CheckStateCritical, err
			}
			for _, attr := range data.Attrs {
				influxDBMetric := InfluxDBMetric{
					Tags: make(map[string]string),
				}
				if attr.Type == smart.AtaDeviceAttributeTypeTempMinMax {
					temp, _, _, _, err := attr.ParseAsTemperature()
					if err != nil {
						return sensu.CheckStateCritical, err
					}
					influxDBMetric.Tags["name"] = disk.Name
					influxDBMetric.Tags["drive_type"] = disk.DriveType.String()
					influxDBMetric.Tags["vendor"] = disk.Vendor
					influxDBMetric.Tags["model"] = disk.Model
					influxDBMetric.Tags["serial_number"] = disk.SerialNumber
					influxDBMetric.Tags["storage_controller"] = disk.StorageController.String()
					influxDBMetric.Value = float64(temp)
					influxDBLine.Metrics[attr.Name] = append(influxDBLine.Metrics[attr.Name], &influxDBMetric)
				} else {
					influxDBMetric.Tags["name"] = disk.Name
					influxDBMetric.Tags["drive_type"] = disk.DriveType.String()
					influxDBMetric.Tags["vendor"] = disk.Vendor
					influxDBMetric.Tags["model"] = disk.Model
					influxDBMetric.Tags["serial_number"] = disk.SerialNumber
					influxDBMetric.Tags["storage_controller"] = disk.StorageController.String()
					influxDBMetric.Value = float64(attr.ValueRaw)
					influxDBLine.Metrics[attr.Name] = append(influxDBLine.Metrics[attr.Name], &influxDBMetric)
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

	influxDBLine.generateInfluxDBLine()
	return sensu.CheckStateOK, nil
}
