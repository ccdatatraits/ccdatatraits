package main

import (
	"fmt"
	"os"
	"strings"
	"strconv"
	"text/scanner"
)

type Header []string
type Postcode int
type PostcodeGroup int

type PostcodeMapping struct {
	postcode Postcode
	group PostcodeGroup
}

func (pm PostcodeMapping) String() string {
	return fmt.Sprintf("%v\t%v", pm.postcode, pm.group)
}

type PostcodeGroups struct {
	header Header
	mappings []PostcodeMapping
}

func (header Header) String() string {
	return fmt.Sprintf("%v\t%v", header[0], header[1])
}

func (pg PostcodeGroups) String() string {
	tmp := make([]string, len(pg.mappings))
	for i, j := range pg.mappings {
		tmp[i] = j.String()
	}
	return fmt.Sprintf("%v\n%v", pg.header, strings.Join(tmp, "\n"))
}

func main() {
	/*const src = `
	"Postcode","Inclusions"
	"820","820,870,829,830,832"
	"828","828,2157,2765,2754,2757"
	"829","829,870,830,832,820,836"
	"830","830,829,870,832,820"
	"832","832,829,830,870,836,837"
	"835","835,847,3806,846,3807,3805,3976,3809,3804"
	`*/
	var (
		s scanner.Scanner
		tok rune
		toggle1stColumn = true
		postcode_group, postcode int
		col1, col2 string
		postcodeGroups *PostcodeGroups
		err error
	)
	//s.Init(strings.NewReader(src))
	s.Init(os.Stdin)
	for tok != scanner.EOF {
		tok = s.Scan()
		nextToken := strings.Replace(s.TokenText(), "\"", "", -1)
		if "," != nextToken && "" != nextToken && s.Position.Line > 1 {
			if s.Position.Line == 2 {
				if "Postcode" == nextToken {
					nextToken = "postcode"
				} else {
					nextToken = "postcode_group"
				}
				if toggle1stColumn {
					col1 = nextToken
				} else {
					col2 = nextToken
					postcodeGroups = &PostcodeGroups{Header{col1, col2},
						[]PostcodeMapping{}}
				}
			} else if toggle1stColumn {
				//fmt.Println("group(string):", nextToken)
				if postcode_group, err = strconv.Atoi(nextToken); err == nil {
					//fmt.Println("group:", postcode_group)
				}
			} else {
				//fmt.Println("postcodes:", nextToken)
				for _, postcode_str := range strings.Split(nextToken, ",") {
					//fmt.Println("postcode(string):", postcode_str)
					if postcode, err = strconv.Atoi(postcode_str); err == nil {
						//fmt.Println("postcode:", postcode)
						postcodeGroups.mappings = append(postcodeGroups.mappings, PostcodeMapping{Postcode(postcode),
							PostcodeGroup(postcode_group)})
						//fmt.Println(postcode, ":", postcode_group)
					}
				}
			}
			toggle1stColumn = !toggle1stColumn
		}
	}
	fmt.Println(postcodeGroups)
}
