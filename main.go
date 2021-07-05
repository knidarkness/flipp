package main

import (
	"errors"
	"fmt"
)

// FlyerDatabase I assume that we use UUID for the flyers IDs. If we would use integers, than
// instead of using map for the storage, we could simply use array (probably sorted)
// (and later binary search for the retrieval operations)
type FlyerDatabase struct {
	flyers map[string]*Flyer
}

// Flyer spamInterval & spamRateLimit are set perFlyer bc they can be hosted
// on platforms with different policies we want to enforce regarding "spam" clicks
type Flyer struct {
	id            string
	clicks        []int
	spamInterval  int
	spamRateLimit int
}

func (f *Flyer) AddClick(timestamp int) error {
	spamIntervalStart := timestamp - f.spamInterval
	if spamIntervalStart < 0 {
		spamIntervalStart = 0
	}
	spamIntervalClicks := f.ClicksDuringInterval(spamIntervalStart, timestamp)
	if spamIntervalClicks >= f.spamRateLimit {
		return errors.New(fmt.Sprintf(
			"too many clicks: attempting to add %d clicks during the interval of %d seconds",
			spamIntervalClicks+1,
			f.spamInterval))
	}
	f.clicks = append(f.clicks, timestamp)
	return nil
}

func (f *Flyer) ClicksDuringInterval(start int, end int) int {
	clicks := 0
	for _, timestamp := range f.clicks {
		if timestamp < start || timestamp > end {
			continue
		}
		clicks += 1
	}
	return clicks
}

func (db *FlyerDatabase) AddFlyer(flyer Flyer) (bool, error) {
	if _, contains := (db.flyers)[flyer.id]; contains {
		return false, errors.New(fmt.Sprintf("Flyer with id %s already stored in the database", flyer.id))
	}
	(db.flyers)[flyer.id] = &flyer
	return true, nil
}

func (db *FlyerDatabase) GetFlyer(id string) (*Flyer, error) {
	if flyer, contains := (db.flyers)[id]; contains {
		return flyer, nil
	}
	return nil, errors.New(fmt.Sprintf("Flyer with id: %s was not found", id))
}

func (db *FlyerDatabase) MostClicked(start int, end int) (*Flyer, error) {
	mostClicked := 0
	var mostClickedId string

	for _, flyer := range db.flyers {
		currentClicks := flyer.ClicksDuringInterval(start, end)
		if currentClicks > mostClicked {
			mostClicked = currentClicks
			mostClickedId = flyer.id
		}
	}

	if mostClickedId == "" {
		return nil, errors.New("there's not flyers in the database")
	}

	return db.flyers[mostClickedId], nil
}

func NewFlyer(id string, spamInterval int, spamLimit int) Flyer {
	return Flyer{
		id,
		[]int{},
		spamInterval,
		spamLimit,
	}
}

func main() {
	db := &FlyerDatabase{flyers: map[string]*Flyer{}}
	db.AddFlyer(NewFlyer("1", 5, 3))
	db.AddFlyer(NewFlyer("2", 5, 3))
	db.AddFlyer(NewFlyer("3", 5, 3))
	_, err := db.AddFlyer(NewFlyer("1", 5, 3))
	println(err.Error())
	f, err := db.GetFlyer("1")
	f.AddClick(1)
	f.AddClick(5)

	f, _ = db.GetFlyer("2")
	f.AddClick(1)
	f.AddClick(3)
	f.AddClick(4)
	f.AddClick(4)
	err = f.AddClick(5)
	if err != nil {
		println(err.Error())
	}

	mc, err := db.MostClicked(1, 4)
	if err != nil {
		println(err.Error())
	}
	println(mc.id)
}
