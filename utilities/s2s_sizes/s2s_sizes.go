// http://weblog.shank.in/input-template-go-for-algorithmic-competitions/

package main

import (
	"bufio"
	"io"
	//"os"
	"strconv"
	"strings"
	"fmt"
	"sort"
	"os"
)

////////////////////////////////////////////////////////////////////////////////

// INPUT TEMPLATE START

type MyInput struct {
	rdr         io.Reader
	lineChan    chan string
	initialized bool
}

func (mi *MyInput) start(done chan struct{}) {
	r := bufio.NewReader(mi.rdr)
	defer func() { close(mi.lineChan) }()
	for {
		line, err := r.ReadString('\n')
		if !mi.initialized {
			mi.initialized = true
			done <- struct{}{}
		}
		mi.lineChan <- strings.TrimSpace(line)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}
}

func (mi *MyInput) readLine() string {
	// if this is the first call, initialize
	if !mi.initialized {
		mi.lineChan = make(chan string)
		done := make(chan struct{})
		go mi.start(done)
		<-done
	}

	res, ok := <-mi.lineChan
	if !ok {
		panic("trying to read from a closed channel")
	}
	return res
}

func (mi *MyInput) readInt() int {
	line := mi.readLine()
	i, err := strconv.Atoi(line)
	if err != nil {
		panic(err)
	}
	return i
}

func (mi *MyInput) readInt64() int64 {
	line := mi.readLine()
	i, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func (mi *MyInput) readInts() []int {
	line := mi.readLine()
	parts := strings.Split(line, " ")
	res := []int{}
	for _, s := range parts {
		tmp, err := strconv.Atoi(s)
		if err != nil {
			panic(err)
		}
		res = append(res, tmp)
	}
	return res
}

func (mi *MyInput) readInt64s() []int64 {
	line := mi.readLine()
	parts := strings.Split(line, " ")
	res := []int64{}
	for _, s := range parts {
		tmp, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			panic(err)
		}
		res = append(res, tmp)
	}
	return res
}

func (mi *MyInput) readWords() []string {
	line := mi.readLine()
	return strings.Split(line, " ")
}

// INPUT TEMPLATE END

////////////////////////////////////////////////////////////////////////////////

type Stats struct {
	st map[int]int
}

func NewStats() *Stats {
	return &Stats{make(map[int]int)}
}

func (stats *Stats) add(segment int) {
	if val, ok := stats.st[segment]; ok {
		stats.st[segment] = val + 1
	} else {
		stats.st[segment] = 1
	}
}

func (stats Stats) String() string {
	output := make([]string, len(stats.st))
	i := 0
	for k, v := range stats.st {
		output[i] = fmt.Sprintf("%v:%v", k, v)
		i++
	}
	sort.Strings(output)
	return strings.Join(output, "\n")
}

func main() {
	/*const src  = `Version: 3
FileIdentifier: target_converters_201706211000
DateCreated: 1498004123
UserNamespace: mm
SegmentNamespace: cs
Mobile: 0
HashSegments: 0

00005790-41fe-4400-abc0-5a1541998ac9 3618 3621 3620 3622
000057a7-36ee-4200-98dd-c6afaa08ee78 2530 2576 2527 2574
000058f8-ebe3-4f00-b020-52abbaa58224 2132 2131 2130 2034 2035 2037 2038 2039 2040 2042 2041 2044 2043 2000 2204 2203 2206 2205 2025 2026 2023 2024 2029 2027 2028 2033 2032 2031 2133 2134 2017 2016 2018 2110 2015 2193 2065 2060 2061 2020 2021 2022 2008 2007 2006 2090 2049 2047 2048 2045 2009 2046 2055 2050 2089 2010 2011
00005909-98d3-4500-9f8e-9d90bcb1210b 2821 2827 2823 2830
00005911-14e2-4d00-af06-030192656b89 2622 2455 2454 2452 2448 2447 2449 2540 2536 2431 2539 2537 2538
0000593f-7006-4500-8202-5d088786f41a 3162 3161 3166 3165 3163 862 3188 3187 3186 3185 3184 3183 3182 3181 3141 3142 3143 3144 3148 3147 3146 3145 3191 3190 3204 3189
000157da-201f-4000-bd2a-8b8730d9970e 3220 5161 3215 5162 5163 5164 5165 3211 5166 3212 5167 5168 3214 5160 4128 2880 3331 3175 2776 4127 3177 2778 3911 3910 2779 5169 4131 4132 3806 3804 3805 3802 3803 4022 2784 2782 2783 4113 4114 4116 4117 4118 4119 3197 3199 3198 4112 3931 3930 4509 4508 3200 4506 3201 3976 4505 4504
0001585c-645e-4400-af32-fb6b05a8edd0 3125 3129 3128 3127 3132 3133 3130 3131 3149 3109 3108 3104 3111 3151
00015880-168f-4500-8b25-276f147f5fda 3067 3066 3065 3057 3058 3121 3070 3078 3072 3071 3073 3031 3068 3081 3044 3039 3000 3079 3050 3005 3052 3003 3051 3002 3054 3053 3008 3101 3056 3055 3006 3010
00015901-2ac5-4e00-b860-509dc494dc18 3161 3067 3066 3065 3057 3121 3122 3126 3124 3070 3123 3078 3071 3068 3184 3183 3182 3181 3000 3141 3142 3143 3144 3079 3146 3050 3005 3052 3003 3051 3002 3103 3054 3053 3008 3101 3056 3102 3006 3206 3207 3205 3010
00015919-26eb-4500-aadc-dcdc4ecce707 4025 4021 4022 4020 4019 4509 4501 4504 4503 4502
00015919-874a-4e00-8c6b-d33e783587a6 2763 2765 2769 2768 2126 2125 2153 2154 2151 2147 2156 2155 2158
	`*/

	//f := strings.NewReader(src)
	//f, _ := os.Open("input_file.in")
	//mi := MyInput{rdr: f}
	 mi := MyInput{rdr: os.Stdin}
	var line string
	for line = mi.readLine(); line != ""; line = mi.readLine() {
		//fmt.Println(line)
	}

	s2sStats := NewStats()

	for line = mi.readLine(); line != ""; line = mi.readLine() {
		s2sRow := strings.Fields(line)
		s2sSegments := s2sRow[1:]
		//fmt.Println(s2sSegments)
		for _, segment := range s2sSegments {
			i, err := strconv.Atoi(segment)
			if err != nil {
				panic(err)
			}
			s2sStats.add(i)
		}
	}

	fmt.Println(s2sStats)
}