package main

import (
	"errors"
	"strings"
	"testing"
)

func TestDataSanitization(t *testing.T) {
	d := Data{CPU: "meow meow\n\n", GPUs: "meow meow $$$$$"}
	d.Sanitize()

	if d.CPU != "meow meow" {
		t.Fatalf("sanitization did not apply")
	}
	if d.GPUs != "meow meow $$$$$" {
		t.Fatalf("sanitization should not apply to GPUs")
	}
}

func TestDataValidation(t *testing.T) {
	datas := []struct {
		Data  Data
		Error error
	}{
		{
			Error: ErrDataBadLength,
		},
		{
			Data:  Data{Project: strings.Repeat("bad", 24)},
			Error: ErrDataBadLength,
		},
		{
			Data:  Data{Project: "meow", Distro: "meow Linux", Kernel: "meow", CPU: "meow", GPUs: "bad gpu"},
			Error: ErrDataBadGPUs,
		},
		{
			Data:  Data{Project: "meow", Distro: "meow Linux", Kernel: "meow", GPUs: "meow", CPU: "bad cpu"},
			Error: ErrDataBadCPUVendor,
		},
		{
			Data:  Data{Project: "meow", Kernel: "meow", GPUs: "meow", CPU: "meow", Flatpak: true, Distro: "bad distro"},
			Error: ErrDataFlatpakDistroMismatch,
		},
	}

	for i, d := range datas {
		if err := d.Data.Validate(); !errors.Is(err, d.Error) {
			t.Fatalf("data fail test case %d succeeded", i)
		}
	}

	d := Data{
		Project: "vinegar",
		Distro:  "KISS Linux",
		Kernel:  "6.5.9acid",
		CPU:     "13th Gen Intel(R) Core(TM) i5-13600K",
		GPUs:    "i915",
	}

	if err := d.Validate(); err != nil {
		t.Errorf("correct data is bad: %s", err)
	}

	d.GPUs += ",radeonsi"
	if err := d.Validate(); err != nil {
		t.Errorf("correct data is bad: %s", err)
	}
}
