package main

import (
	"reflect"
	"testing"
)

// 総合的なテスト。
func TestNewColor(t *testing.T) {
	tests := []struct {
		str  string
		want Color
	}{
		{"#F00", Color{255, 0, 0}},
		{"#03c", Color{0, 51, 204}},
		{"#AAB0", Color{170, 170, 187}},
		{"#5G5", Color{0, 0, 0}},

		{"#AF0532", Color{175, 5, 50}},
		{"#01975233", Color{1, 151, 82}},

		{"RGB(20,30,40)", Color{20, 30, 40}},
		{"rgba(160,0,0,50%)", Color{160, 0, 0}},

		{"RGB(20%,30%,40%)", Color{51, 76, 102}},
		{"rgba(90%,100%,0%,50%)", Color{229, 255, 0}},

		{"hsl(169,57%,75%)", Color{154, 227, 214}},
		{"HSLA(345,23%,52%,67%)", Color{160, 104, 118}},
		{"hsla(7,34%,40%)", Color{136, 75, 67}},

		{"hsl(10%,54%,32%)", Color{125, 90, 37}},
		{"HSLA(57%,43%,56%,10%)", Color{94, 150, 191}},
		{"hsla(80%,78%,67%/97%)", Color{210, 105, 236}},
	}

	for _, tt := range tests {
		if got, _ := NewColor(tt.str); got != tt.want {
			t.Errorf("want = %s, got = %s", tt.want, got)
		}
	}
}

func TestHex2Decimal(t *testing.T) {
	tests := []struct {
		str  string
		want int
	}{
		{"00", 0},
		{"3c", 60},
		{"FF", 255},
		{"e", 238},
		{"FFA", 0},
		{"G", 0},
	}

	for _, tt := range tests {
		if got := hex2Decimal(tt.str); got != tt.want {
			t.Errorf("want = %d, got = %d", tt.want, got)
		}
	}
}

func TestHsl2Rgb(t *testing.T) {
	tests := []struct {
		hsl  []int
		want []int
	}{
		{[]int{100, 50, 0}, []int{0, 0, 0}},
		{[]int{21, 72, 30}, []int{131, 59, 21}},
		{[]int{90, 20, 55}, []int{140, 163, 117}},
		{[]int{169, 57, 75}, []int{154, 227, 214}},
		{[]int{212, 32, 66}, []int{140, 166, 196}},
		{[]int{300, 100, 21}, []int{107, 0, 107}},
		{[]int{345, 23, 52}, []int{160, 104, 118}},
		{[]int{385, 38, 24}, []int{84, 57, 37}},
	}

	for _, tt := range tests {
		r, g, b, _ := hsl2Rgb(tt.hsl[0], tt.hsl[1], tt.hsl[2])
		got := []int{r, g, b}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("want = %v, got = %v", tt.want, got)
		}
	}
}
