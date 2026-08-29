package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/LuigiD5555/origami/internal/reference"
)

func main() {
	experiment := flag.String("experiment", "experiments/EXP-001-relational-state/experiment.json", "experiment specification")
	out := flag.String("out", "", "optional evidence output path")
	flag.Parse()
	e, err := reference.Load(*experiment); if err != nil { fatal(err) }
	ev := reference.RunExperiment(e)
	b, err := json.MarshalIndent(ev,"","  "); if err != nil { fatal(err) }
	if *out != "" { if err:=os.WriteFile(*out,append(b,'\n'),0644); err!=nil { fatal(err) } } else { fmt.Println(string(b)) }
}
func fatal(err error){ fmt.Fprintln(os.Stderr,"origami-reference:",err); os.Exit(1) }
