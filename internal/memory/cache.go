package memory

import "sort"

type WorkingSetEntry struct {
	Item ContextItem `json:"item"`
	Tick uint64      `json:"tick"`
}

type WorkingSet struct {
	maxTokens int
	usedTokens int
	tick uint64
	entries map[string]WorkingSetEntry
}

func NewWorkingSet(maxTokens int) *WorkingSet {
	if maxTokens <= 0 { maxTokens = 4000 }
	return &WorkingSet{maxTokens:maxTokens, entries:map[string]WorkingSetEntry{}}
}

func (w *WorkingSet) Put(item ContextItem) {
	if w == nil || item.Address == "" || item.TokenCost <= 0 || item.TokenCost > w.maxTokens { return }
	w.tick++
	if old, ok := w.entries[item.Address]; ok { w.usedTokens -= old.Item.TokenCost }
	w.entries[item.Address] = WorkingSetEntry{Item:item, Tick:w.tick}
	w.usedTokens += item.TokenCost
	w.evict()
}

func (w *WorkingSet) Get(address string) (ContextItem, bool) {
	if w == nil { return ContextItem{}, false }
	entry, ok := w.entries[address]; if !ok { return ContextItem{}, false }
	w.tick++; entry.Tick=w.tick; w.entries[address]=entry
	return entry.Item, true
}

func (w *WorkingSet) UsedTokens() int { if w==nil{return 0}; return w.usedTokens }
func (w *WorkingSet) MaxTokens() int { if w==nil{return 0}; return w.maxTokens }

func (w *WorkingSet) Items() []ContextItem {
	if w == nil { return nil }
	entries:=make([]WorkingSetEntry,0,len(w.entries));for _,entry:=range w.entries{entries=append(entries,entry)}
	sort.Slice(entries,func(i,j int)bool{return entries[i].Tick>entries[j].Tick})
	out:=make([]ContextItem,0,len(entries));for _,entry:=range entries{out=append(out,entry.Item)};return out
}

func (w *WorkingSet) evict() {
	for w.usedTokens > w.maxTokens && len(w.entries)>0 {
		var oldest string; var tick uint64 = ^uint64(0)
		for address,entry:=range w.entries{if entry.Tick<tick{oldest=address;tick=entry.Tick}}
		entry:=w.entries[oldest];delete(w.entries,oldest);w.usedTokens-=entry.Item.TokenCost
	}
}
