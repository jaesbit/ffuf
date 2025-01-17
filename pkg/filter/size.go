package filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jaesbit/ffuf/pkg/ffuf"
)

type SizeFilter struct {
	Value []ffuf.ValueRange
}

func NewSizeFilter(value string) (ffuf.FilterProvider, error) {
	var intranges []ffuf.ValueRange
	for _, sv := range strings.Split(value, ",") {
		vr, err := ffuf.ValueRangeFromString(sv)
		if err != nil {
			return &SizeFilter{}, fmt.Errorf("Size filter or matcher (-fs / -ms): invalid value: %s", sv)
		}

		intranges = append(intranges, vr)
	}
	return &SizeFilter{Value: intranges}, nil
}

func (f *SizeFilter) Filter(response *ffuf.Response) (bool, error) {
	for _, iv := range f.Value {
		if iv.Min <= response.ContentLength && response.ContentLength <= iv.Max {
			return true, nil
		}
	}
	return false, nil
}

func (f *SizeFilter) Repr() string {
	var strval []string
	for _, iv := range f.Value {
		if iv.Min == iv.Max {
			strval = append(strval, strconv.Itoa(int(iv.Min)))
		} else {
			strval = append(strval, strconv.Itoa(int(iv.Min))+"-"+strconv.Itoa(int(iv.Max)))
		}
	}
	return fmt.Sprintf("Response size: %s", strings.Join(strval, ","))
}
