package core

import (
	"fmt"
	"sort"
)

type Exploration struct {
	Schema      string       `json:"schema"`
	Keys        []string     `json:"keys"`
	Domain      []Value      `json:"domain"`
	Initials    int          `json:"initial_states"`
	Trajectories []Trajectory `json:"trajectories"`
	Counts      map[string]int `json:"termination_counts"`
}

func EnumerateStates(keys []string, domain []Value) ([]State,error){
	keys=append([]string(nil),keys...);sort.Strings(keys);if len(keys)==0{return nil,fmt.Errorf("enumeration requires keys")};if len(domain)==0{return nil,fmt.Errorf("enumeration requires non-empty domain")};for _,value:=range domain{if !ValidateStatus(value.Status){return nil,fmt.Errorf("invalid domain status %q",value.Status)}}
	out:=[]State{};var rec func(int,State);rec=func(index int,state State){if index==len(keys){out=append(out,CloneState(state));return};key:=keys[index];for _,value:=range domain{state.Values[key]=value;rec(index+1,state)}};rec(0,State{Values:map[string]Value{}});return out,nil
}

func Explore(keys []string,domain []Value,machine Machine,contexts []Context)(Exploration,error){
	initials,err:=EnumerateStates(keys,domain);if err!=nil{return Exploration{},err};report:=Exploration{Schema:SchemaR1+".exploration",Keys:append([]string(nil),keys...),Domain:append([]Value(nil),domain...),Initials:len(initials),Counts:map[string]int{}}
	sort.Strings(report.Keys);for _,initial:=range initials{trajectory,err:=Execute(initial,machine,contexts);if err!=nil{return Exploration{},err};report.Trajectories=append(report.Trajectories,trajectory);report.Counts[trajectory.Terminated]++};return report,nil
}
