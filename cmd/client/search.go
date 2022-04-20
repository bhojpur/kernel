package cmd

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/bhojpur/kernel/pkg/types"
	"github.com/layer-x/layerx-commons/lxhttpclient"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search available images in the targeted Bhojpur Kernel Image Repository",
	Long: `
Usage:

kernctl search

  - or -

kernctl search --imageName <imageName>

Requires that you first authenticate to a Bhojpur Kernel image repository with 'kernctl login'`,
	Run: func(cmd *cobra.Command, args []string) {
		c, err := getHubConfig()
		if err != nil {
			logrus.Fatal(err)
		}
		resp, body, err := lxhttpclient.Get(c.URL, "/images", nil)
		if err != nil {
			logrus.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			logrus.Fatal(fmt.Sprintf("failed with status %v: %s", resp.StatusCode, string(body)))
		}
		var images []*types.UserImage
		if err := json.Unmarshal(body, &images); err != nil {
			logrus.Fatal(err)
		}
		filteredImages := images[:0]
		if imageName != "" {
			for _, image := range images {
				if !strings.Contains(image.Name, imageName) {
					filteredImages = append(filteredImages, image)
				}
			}
		} else {
			filteredImages = images
		}
		printUserImages(filteredImages...)
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVar(&imageName, "imageName", "", "<string,optional> search images by names containing this string")
}
