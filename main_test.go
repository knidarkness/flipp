package main

import (
	"fmt"
	"testing"
)

func TestAddingOneFlyer(t *testing.T) {
	db := &FlyerDatabase{flyers: map[string]*Flyer{}}
	db.AddFlyer(NewFlyer("1", 5, 3))
	if len(db.flyers) != 1 {
		t.Errorf("Database should have 1 flyer, has %d instead", len(db.flyers))
	}
}

func TestAddingFlyerSameId(t *testing.T) {
	db := &FlyerDatabase{flyers: map[string]*Flyer{}}
	db.AddFlyer(NewFlyer("1", 5, 3))
	ok, err := db.AddFlyer(NewFlyer("1", 5, 3))
	if len(db.flyers) != 1 {
		t.Errorf("Database should have 1 flyer, has %d instead", len(db.flyers))
	}

	if ok {
		t.Errorf("Successfully added flyer with duplicate id")
	}

	errMsg := "Flyer with id 1 already stored in the database"
	if err.Error() != errMsg {
		t.Errorf("Error mismatch for the overrriding existing flyer. Got %s, should be: %s", err.Error(), errMsg)
	}
}

func TestAddingSingleClick(t *testing.T) {
	db := &FlyerDatabase{flyers: map[string]*Flyer{}}
	db.AddFlyer(NewFlyer("1", 5, 3))
	flyer, _ := db.GetFlyer("1")
	flyer.AddClick(1)
	if len(flyer.clicks) != 1 {
		t.Errorf("Expected to add 1 click to the flyer")
	}
}

func TestAddingClicks(t *testing.T) {
	db := &FlyerDatabase{flyers: map[string]*Flyer{}}
	db.AddFlyer(NewFlyer("1", 5, 3))
	flyer, _ := db.GetFlyer("1")
	flyer.AddClick(1)
	flyer.AddClick(1)
	flyer.AddClick(2)
	if len(flyer.clicks) != 3 {
		t.Errorf("Expected to add 3 clicks to the flyer")
	}
}

func TestAddingClicksSpamDetection(t *testing.T) {
	db := &FlyerDatabase{flyers: map[string]*Flyer{}}
	spamLimit := 3
	db.AddFlyer(NewFlyer("1", 5, spamLimit))
	flyer, _ := db.GetFlyer("1")
	flyer.AddClick(1)
	flyer.AddClick(2)
	err := flyer.AddClick(3)
	if err != nil {
		t.Errorf("Should have added click under spam limits")
	}
	err = flyer.AddClick(4)
	if err == nil {
		t.Errorf("Should have failed to add click over spam limit")
	}
	flyer.AddClick(5)
	if len(flyer.clicks) != 3 {
		t.Errorf(fmt.Sprintf("Expected to have %d clicks in the flyer storage", spamLimit))
	}
}
