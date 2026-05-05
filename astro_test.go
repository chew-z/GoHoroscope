package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mshafiee/swephgo"
)

func requireEphemeris(t *testing.T) {
	t.Helper()
	if os.Getenv("SWISSPATH") == "" {
		t.Skip("SWISSPATH not set; ephemeris data required for retrograde tests")
	}
}

// Test 1: Pluto on 2026-05-05 must classify the next station as a retrograde
// start (idir == -1) and resolve to ~2026-05-07. This directly exercises the
// numerical-noise bug that produced the spurious row in `horoscope -r`.
func TestRetro_PlutoStation_2026May(t *testing.T) {
	requireEphemeris(t)

	loc, _ := time.LoadLocation("UTC")
	start := time.Date(2026, time.May, 5, 12, 0, 0, 0, loc)

	var tx float64
	var idir int
	serr := make([]byte, 256)

	if rv := RetroUt(start, swephgo.SePluto, swephgo.SeflgSwieph, &tx, &idir, &serr); rv < 0 {
		t.Fatalf("RetroUt returned error: %s", strings.TrimRight(string(serr), "\x00"))
	}

	if idir != -1 {
		t.Errorf("Pluto next station from 2026-05-05: got idir=%d, want -1 (retrograde start)", idir)
	}

	got := jdToUTC(&tx)
	want := time.Date(2026, time.May, 7, 0, 0, 0, 0, time.UTC)
	if diff := got.Sub(want); diff < -2*24*time.Hour || diff > 2*24*time.Hour {
		t.Errorf("Pluto station time: got %s, want within ±2 days of %s", got, want)
	}
}

// Test 2: Mercury sanity — fast-planet behavior must remain unchanged.
func TestRetro_MercuryStation_Plausible(t *testing.T) {
	requireEphemeris(t)

	loc, _ := time.LoadLocation("UTC")
	start := time.Date(2026, time.May, 5, 12, 0, 0, 0, loc)

	var tx float64
	var idir int
	serr := make([]byte, 256)

	if rv := RetroUt(start, swephgo.SeMercury, swephgo.SeflgSwieph, &tx, &idir, &serr); rv < 0 {
		t.Fatalf("RetroUt returned error: %s", strings.TrimRight(string(serr), "\x00"))
	}

	got := jdToUTC(&tx)
	if got.Before(start) || got.After(start.AddDate(0, 4, 0)) {
		t.Errorf("Mercury station: got %s, want within 4 months after %s", got, start)
	}
	if idir != -1 && idir != 1 {
		t.Errorf("Mercury idir: got %d, want -1 or +1", idir)
	}
}

// Test 3: PrintRetro must not emit a Pluto row whose start AND end both fall
// within 7 days of `start` — that's the artifact the old code produced.
func TestPrintRetro_NoSpuriousEarlyPlutoRow(t *testing.T) {
	requireEphemeris(t)

	if location == nil {
		var err error
		location, err = time.LoadLocation("UTC")
		if err != nil {
			t.Fatalf("failed to load UTC: %v", err)
		}
	}

	start := time.Date(2026, time.May, 5, 13, 2, 0, 0, location)
	end := start.AddDate(1, 0, 1)

	var buf bytes.Buffer
	printRetroTo(&buf, start, end)

	out := buf.String()
	if !strings.Contains(out, "Pluto") {
		t.Fatalf("expected Pluto row in output, got:\n%s", out)
	}

	cutoff := start.AddDate(0, 0, 7)
	for _, line := range strings.Split(out, "\n") {
		if !strings.Contains(line, "Pluto") {
			continue
		}
		s, e, ok := parseRetroRow(line)
		if !ok {
			continue
		}
		if s.Before(cutoff) && e.Before(cutoff) {
			if s.Equal(start) {
				continue
			}
			t.Errorf("spurious Pluto row: both endpoints within 7 days of start: %s", line)
		}
	}
}

// Test 4: Strictly non-decreasing Starts and Starts <= Ends per row.
func TestPrintRetro_MonotonicProgress(t *testing.T) {
	requireEphemeris(t)

	if location == nil {
		var err error
		location, err = time.LoadLocation("UTC")
		if err != nil {
			t.Fatalf("failed to load UTC: %v", err)
		}
	}

	start := time.Date(2026, time.May, 5, 13, 2, 0, 0, location)
	end := start.AddDate(1, 0, 1)

	var buf bytes.Buffer
	printRetroTo(&buf, start, end)

	perPlanetLast := map[string]time.Time{}
	for _, line := range strings.Split(buf.String(), "\n") {
		planet := planetFromRow(line)
		if planet == "" {
			continue
		}
		s, e, ok := parseRetroRow(line)
		if !ok {
			continue
		}
		if s.After(e) {
			t.Errorf("%s: Starts > Ends: %s", planet, line)
		}
		if last, seen := perPlanetLast[planet]; seen && s.Before(last) {
			t.Errorf("%s: Starts went backward (last=%s, this=%s) in row: %s", planet, last, s, line)
		}
		perPlanetLast[planet] = s
	}
}

// parseRetroRow extracts the two timestamps from a termtables-rendered row of
// the form "| Pluto | 06 May 26 17:35 CEST | 15 Oct 26 ... |".
func parseRetroRow(line string) (time.Time, time.Time, bool) {
	cells := strings.Split(line, "|")
	if len(cells) < 4 {
		return time.Time{}, time.Time{}, false
	}
	s := strings.TrimSpace(cells[2])
	e := strings.TrimSpace(cells[3])
	st, err1 := time.Parse(time.RFC822, s)
	et, err2 := time.Parse(time.RFC822, e)
	if err1 != nil || err2 != nil {
		return time.Time{}, time.Time{}, false
	}
	return st, et, true
}

func planetFromRow(line string) string {
	cells := strings.Split(line, "|")
	if len(cells) < 4 {
		return ""
	}
	planet := strings.TrimSpace(cells[1])
	switch planet {
	case "Mercury", "Venus", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune", "Pluto":
		return planet
	}
	return ""
}
