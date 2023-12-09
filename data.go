package main

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Data struct {
	Project string `json:"project"`
	Distro  string `json:"distro"`
	Kernel  string `json:"kernel"`
	Flatpak bool   `json:"flatpak"`
	AVX     bool   `json:"avx"`
	CPU     string `json:"cpu"`
	GPUs    string `json:"gpu"` // list seperated by commas
}

var (
	ErrDataBadLength             = errors.New("data member has invalid length")
	ErrDataFlatpakDistroMismatch = errors.New("flatpak distro has mismatch")
	ErrDataBadGPUs               = errors.New("gpus given is invalid")
	ErrDataBadCPUVendor          = errors.New("cpu given has bad vendor")
	ErrDataBadCPULength          = errors.New("cpu given has bad length")
)

// Generate the CSV Header based on the JSON tags in the Data struct at run-time
var CSVHeader = func() (h []string) {
	var d Data
	t := reflect.TypeOf(d)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		h = append(h, f.Tag.Get("json"))
	}

	return
}()

func (d Data) CSV() []string {
	return []string{
		d.Project, d.Distro, d.Kernel,
		strconv.FormatBool(d.AVX), strconv.FormatBool(d.Flatpak),
		d.CPU, d.GPUs,
	}
}

func (d *Data) Sanitize() {
	reg := regexp.MustCompile(`[^a-zA-Z0-9@*().\-\_ ]+`)
	r := reflect.ValueOf(d).Elem()

	// Loops over all string members in Data and sanitizes them
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		// GPU has it's own filter in Validate(), it's index in the struct
		// is the last one
		if i+1 != r.NumField() && f.Kind() == reflect.String {
			f.SetString(reg.ReplaceAllString(f.String(), ""))
		}
	}
}

func (d Data) Validate() error {
	// Length of CPU is a different value
	for _, m := range []string{d.Project, d.Distro, d.Kernel, d.GPUs} {
		// Reasonable limit
		if m == "" || len(m) > 256 {
			return ErrDataBadLength
		}

		// Sanitize all fields by making
	}

	// The flatpak runtime always modifies the os-release file
	if d.Flatpak && !strings.Contains(d.Distro, "Flatpak runtime") {
		return fmt.Errorf("%w: %s", ErrDataFlatpakDistroMismatch, d.Distro)
	}

	// GPU driver list must be comma seperated with alphanumerical characters
	// edit: turns out some people have simple-framebuffer as a GPU option, so adding hyphen to filter
	if mgpu, _ := regexp.MatchString(`^[a-zA-Z0-9,-]+$`, d.GPUs); !mgpu {
		return fmt.Errorf("%w: %s", ErrDataBadGPUs, d.GPUs)
	}

	// limit is 64 characters as defined by the Linux kernel.
	if len(d.CPU) > 63 {
		return fmt.Errorf("%w: %s", ErrDataBadCPULength, d.CPU)
	}

	vcpu := false
	// Pretty much the only X86 CPUs that Roblox can run on under WINE
	for _, v := range []string{"Intel", "AMD"} {
		if strings.Contains(d.CPU, v) {
			vcpu = true
		}
	}

	if !vcpu {
		return fmt.Errorf("%w: %s", ErrDataBadCPUVendor, d.CPU)
	}

	return nil
}
