package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"agent-plan/internal/chapter01"
	"agent-plan/internal/chapter02"
	"agent-plan/internal/chapter03"
	"agent-plan/internal/chapter04"
	"agent-plan/internal/chapter05"
	"agent-plan/internal/chapter06"
	"agent-plan/internal/chapter07"
	"agent-plan/internal/chapter08"
	"agent-plan/internal/chapter09"
	"agent-plan/internal/chapter10"
	"agent-plan/internal/chapter11"
	"agent-plan/internal/chapter12"
	"agent-plan/internal/chapter13"
	"agent-plan/internal/chapter14"
	"agent-plan/internal/shared"
)

func main() {
	chapterFlag := flag.String("chapter", "", "chapter number, for example 01 or 13")
	flag.Parse()

	chapters := []shared.Chapter{
		chapter01.Chapter(),
		chapter02.Chapter(),
		chapter03.Chapter(),
		chapter04.Chapter(),
		chapter05.Chapter(),
		chapter06.Chapter(),
		chapter07.Chapter(),
		chapter08.Chapter(),
		chapter09.Chapter(),
		chapter10.Chapter(),
		chapter11.Chapter(),
		chapter12.Chapter(),
		chapter13.Chapter(),
		chapter14.Chapter(),
	}

	if *chapterFlag == "" {
		for _, chapter := range chapters {
			fmt.Printf("%s  %s\n", chapter.Number, chapter.Title)
		}
		return
	}

	want := strings.TrimSpace(*chapterFlag)
	for _, chapter := range chapters {
		if chapter.Number == want {
			fmt.Printf("Chapter %s: %s\nSource: %s\nPackage: internal/chapter%s\n", chapter.Number, chapter.Title, chapter.Source, chapter.Number)
			if chapter.Number == "01" {
				if err := chapter01.RunMockDemo(context.Background(), os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "run chapter 01 demo: %v\n", err)
					os.Exit(1)
				}
			}
			return
		}
	}

	fmt.Fprintf(os.Stderr, "chapter %q not found\n", want)
	os.Exit(1)
}
