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
	"bytes"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/bhojpur/kernel/cmd/server/assets"
	"github.com/bhojpur/kernel/cmd/server/build"
	"github.com/spf13/cobra"
)

const (
	grubDockerImage = "bhojpur/grub:0.2.0"
)

var (
	packFormat     string
	packOutFile    string
	packKernelFile string

	withoutDocker bool
	keepTmpdir    bool
)

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "pack Bhojpur Kernel to release format, e.g., iso",
	Run: func(cmd *cobra.Command, args []string) {
		err := runPackage()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func runPackage() error {
	base, err := ioutil.TempDir("", "bhojpur-package")
	if err != nil {
		return err
	}
	if !keepTmpdir {
		defer os.RemoveAll(base)
	} else {
		log.Println(base)
	}

	isoBase := filepath.Join(base, "iso")

	err = os.MkdirAll(isoBase, 0755)
	if err != nil {
		return err
	}

	err = extractBootDir(isoBase)
	if err != nil {
		return err
	}

	kfile, err := getKernelFile(base)
	if err != nil {
		return err
	}

	err = copyfile(
		filepath.Join(isoBase, "boot", "kernel.elf"),
		kfile,
	)
	if err != nil {
		return err
	}

	tmpOutFile := filepath.Join(base, "bhojpur-kernel.iso")
	err = mkiso(tmpOutFile, isoBase, base)
	if err != nil {
		return err
	}

	if packOutFile == "" {
		packOutFile = "bhojpur-kernel.iso"
	}
	return copyfile(packOutFile, tmpOutFile)
}

func dockerImageExists(imageName string) error {
	var stderr bytes.Buffer
	cmd := exec.Command("docker", "inspect", imageName)
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("inspect docker image error:%s", stderr.String())
	}
	return nil
}

func dockerPullImage(imageName string) error {
	cmd := exec.Command("docker", "pull", imageName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("pull docker image error:%s", err)
	}
	return nil
}

func mkiso(outfile, isobase, moutbase string) error {
	var stderr bytes.Buffer
	if withoutDocker {
		cmd := exec.Command("grub-mkrescue", "-o", outfile, isobase)
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			log.Print(stderr.String())
		}
		return err
	}

	err := dockerImageExists(grubDockerImage)
	if err != nil {
		log.Print(err)
		log.Printf("trying to pull docker image `%s`", grubDockerImage)
		err = dockerPullImage(grubDockerImage)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command(
		"docker", "run", "--rm",
		"-v", moutbase+":"+moutbase,
		"-w", moutbase,
		grubDockerImage,
		"grub-mkrescue", "-o", outfile, isobase,
	)
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		log.Print(stderr.String())
	}
	return err
}

func extractBootDir(base string) error {
	bootfs := assets.Boot
	err := fs.WalkDir(assets.Boot, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		outFullPath := filepath.Join(base, path)
		if d.IsDir() {
			err = os.MkdirAll(outFullPath, 0755)
			if err != nil {
				return err
			}
			return nil
		}

		content, err := bootfs.ReadFile(path)
		if err != nil {
			return err
		}
		err = os.WriteFile(outFullPath, content, 0644)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func getKernelFile(base string) (string, error) {
	if packKernelFile != "" {
		return packKernelFile, nil
	}
	outputFile := filepath.Join(base, "kernel.elf")
	b := build.NewBuilder(build.Config{
		GoRoot:  goroot,
		Basedir: base,
		GoArgs: []string{
			"-o", outputFile,
		},
	})
	err := b.Build()
	if err != nil {
		return "", err
	}

	return outputFile, nil
}

func copyfile(dst, src string) error {
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, content, 0644)
}

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringVarP(&packFormat, "format", "f", "iso", "package format, values `iso`")
	packCmd.Flags().StringVarP(&packKernelFile, "kernel", "k", "", "the kernel file, if empty current package will be built as kernel")
	packCmd.Flags().StringVarP(&packOutFile, "output", "o", "bhojpur-kernel.iso", "file name of output")
	packCmd.Flags().BoolVar(&keepTmpdir, "keep-tmp", false, "keep temp dir, for debugging")
	packCmd.Flags().BoolVarP(&withoutDocker, "without-docker", "d", false, "using docker for grub tools")
}
