package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/openconfig/ygot/ygot"

	oopt "github.com/osrg/oopt/pkg/gnmi"
	model "github.com/osrg/oopt/pkg/model_gnmi"
)

var (
	current *model.Component
	git_dir = "/etc/oopt"
)

const (
	CONFIG_FILE = "config_gnmi.json"
	opticalModuleNum int = 8
)

func initConfig() error {
	name := fmt.Sprintf("Opt3")
	d := &model.Component{
		Name: &name,
		OpticalChannel: &model.Component_OpticalChannel{
		Frequency: ygot.Uint64(0),
		},
	}
	/*d := model.Component{}
	for i := 1; i <= opticalModuleNum; i++ {
		key := fmt.Sprintf("Opt%d", i)
		d[i] = &model.Component{
			Name: &key,
			OpticalChannel: &model.Component_OpticalChannel{
				Frequency: ygot.Uint64(0),
			},
		}
	}*/
	json, err := ygot.EmitJSON(d, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
	})
	if err != nil {
		return err
	}
	fmt.Printf(json)
	file, err := os.Create(fmt.Sprintf("%s/%s", git_dir, CONFIG_FILE))
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(([]byte)(json))
	file.Write([]byte("\n"))
	return nil
}

func callback(newConfig ygot.ValidatedGoStruct) error {
	buf, err := ygot.EmitJSON(newConfig, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
	})
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	file, err := os.Create(fmt.Sprintf("%s/%s", git_dir, CONFIG_FILE))
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer file.Close()
	file.Write(([]byte)(buf))
	file.Write([]byte("\n"))

	return nil
}

func main() {
	port := flag.Int64("port", 10164, "Listen port")
	flag.Parse()

	err := initConfig()
	if err != nil {
		panic(fmt.Sprintf("init: %v", err))
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", git_dir, CONFIG_FILE))
	if err != nil {
		panic(fmt.Sprintf("open: %v", err))
	}
	current = &model.Component{}
	model.Unmarshal(data, current)

	servermodel := oopt.NewModel(
		oopt.ModelData,
		reflect.TypeOf((*model.Component)(nil)),
		model.SchemaTree["PacketTransponder"],
		model.Unmarshal,
		model.ΛEnum,
	)

	if err = current.Validate(); err != nil {
		panic(fmt.Sprintf("validation failed: %v", err))
	}

	json, err := ygot.EmitJSON(current, &ygot.EmitJSONConfig{
		Format: ygot.RFC7951,
		Indent: "  ",
		RFC7951Config: &ygot.RFC7951JSONConfig{
			AppendModuleName: true,
		},
	})
	if err != nil {
		panic(fmt.Sprintf("EmitJSON failed: %v", err))
	}

	srv, err := oopt.NewServer(servermodel, []byte(json), *port, callback, nil)
	if err != nil {
		panic(fmt.Sprintf("NewServer() failed: %v", err))
	}
	srv.Serve()
}
