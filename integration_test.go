// +build integration

package autotest

import (
	"log"
	"testing"
)

const (
	HelloEvent EventID = iota
)

func TestCase1(t *testing.T) {
	p := NewProcess("bash", "-c", "for ((i = 0; i < 10; i++)); do sleep 1; echo \"hello $i\"; done")
	p.AddTrigger("hello (.*)", HelloEvent)
	if err := p.Start(); err != nil {
		t.Fatalf("Cannot start the process: %v", err)
	}

	if _, _, err := Expect(HelloEvent, "0"); err != nil {
		t.Fatal(err)
	}
	log.Println("HelloEvent")
	for i := 0; i < 5; i++ {
		if _, _, err := Skip(); err != nil {
			t.Fatal(err)
		}
		log.Println("Event skipped")
	}
	if _, _, err := Expect(HelloEvent, "6"); err != nil {
		t.Fatal(err)
	}
	log.Println("HelloEvent")

	if err := p.Wait(); err != nil {
		t.Fatalf("Error when waiting for the process: %v", err)
	}
}
