//go:generate packer-sdc struct-markdown  KmiFilterOptions

package ksyun

import (
	"fmt"
	"regexp"
	"sort"
	"time"

	"github.com/KscSDK/ksc-sdk-go/service/kec"
	"github.com/mitchellh/mapstructure"
)

type Ks3Image struct {
	ImageId       string `mapstructure:"ImageId"`
	ImageType     string `mapstructure:"Type"`
	Name          string `mapstructure:"Name"`
	CreationDate  string `mapstructure:"CreationDate"`
	Platform      string `mapstructure:"Platform"`
	ImageSource   string `mapstructure:"ImageSource"`
	RealImageId   string `mapstructure:"RealImageId"`
	IsCloudMarket bool   `mapstructure:"IsCloudMarket"`
	SysDisk       int    `mapstructure:"SysDisk"`
	IsPublic      bool   `mapstructure:"IsPublic"`
}

func (ki *Ks3Image) Filler(rawMap interface{}) error {
	return mapstructure.Decode(rawMap, ki)
}

type Ks3Images struct {
	ImagesSet []*Ks3Image `mapstructure:"ImagesSet"`
}

// Filler convert map structure to struct.
func (k *Ks3Images) Filler(rawMap interface{}) error {
	return mapstructure.Decode(rawMap, k)
}

type imageSort []*Ks3Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	itime, _ := time.Parse(time.RFC3339, (*a[i]).CreationDate)
	jtime, _ := time.Parse(time.RFC3339, (*a[j]).CreationDate)
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent AMI out of a slice of images.
func mostRecentKmi(images []*Ks3Image) *Ks3Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}

// GetFilteredImage to get source image with filers
func (d *KmiFilterOptions) GetFilteredImage(params *map[string]interface{}, kecConn *kec.Kec) (image *Ks3Image, err error) {
	// We have filters to apply
	// var imageIds []string

	resp, err := kecConn.DescribeImages(params)
	if err != nil {
		return image, fmt.Errorf("error on reading Image list")
	}

	ks3Images := &Ks3Images{}
	if err := ks3Images.Filler(*resp); err != nil {
		return image, err
	}
	if ks3Images == nil || len(ks3Images.ImagesSet) < 1 {
		return image, fmt.Errorf("cannot match any image")
	}
	data := ks3Images.ImagesSet
	if d.Platform != "" {
		var dataFilter []*Ks3Image
		for _, v := range data {
			if v.Platform == d.Platform {
				dataFilter = append(dataFilter, v)
			}
		}
		data = dataFilter
	}
	if d.ImageSource != "" {
		var dataFilter []*Ks3Image
		for _, v := range data {
			if v.ImageSource == d.ImageSource {
				dataFilter = append(dataFilter, v)
			}
		}
		data = dataFilter
	}

	if d.NameRegex != "" {
		var dataFilter []*Ks3Image
		r := regexp.MustCompile(d.NameRegex)
		for _, v := range data {
			if r == nil || r.MatchString(v.Name) {
				dataFilter = append(dataFilter, v)
			}
		}
		data = dataFilter
	}

	if len(data) < 1 {
		return image, fmt.Errorf("cannot filter an image with the filter option %v", params)
	}

	if d.MostRecent {
		image = mostRecentKmi(data)
	} else {
		image = data[0]
	}
	return image, nil
}
