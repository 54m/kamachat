package main

import (
	"fmt"
	"log"

	"github.com/54m/kamachat/generator"
	"github.com/docopt/docopt-go"
)

var appVersion = `Okama Nanchatte (kamachat) command version 0.1.0
Copyright (c) 2020 54m
Released under the MIT License.
https://github.com/54m/kamachat`

var usage = `Usage:
  kamachat [options] [<name>]

Options:
  -h, --help      ヘルプを表示.
  -V, --version   バージョンを表示.
  -e <number>     絵文字/顔文字の最大連続数 [default: 4].
  -p <level>      句読点挿入頻度レベル [min:0, max:3] [default: 0].`

// TODO: --type おかまタイプ (絵文字乱用, 顔文字乱用, 句読点, 若作り)

func main() {
	parser := &docopt.Parser{
		OptionsFirst: true,
	}
	args, _ := parser.ParseArgs(usage, nil, appVersion)
	config := generator.Config{}
	err := args.Bind(&config)
	if err != nil {
		log.Fatal(err)
	}

	result, err := generator.Start(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", result)
}
