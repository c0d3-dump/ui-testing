package main

import (
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/proto"
)

type RodEvent uint8

const (
	Nothing RodEvent = iota
	Wait
	LeftMouseClick
	RightMouseClick
	Enter
	Tab
	Space
	Ctrl
	Alt
)

type TestRod struct {
	Selector string
	Input    string
	Event    RodEvent
}

type Main struct {
	Url         string
	ShowBrowser bool
	Screenshots bool
	Tests       []TestRod
}

func RunTest(page *rod.Page, test *TestRod, idx *int) {
	var el *rod.Element
	if test.Selector != "" {
		el = page.MustElement(test.Selector)
	}

	if test.Input != "" {
		el.MustInput(test.Input)
	}

	switch test.Event {
	case Nothing:
		// do nothing!
	case Wait:
		duration, err := strconv.Atoi(test.Input)
		if err != nil {
			panic(err)
		}
		page.WaitIdle(time.Second * time.Duration(duration))
	case LeftMouseClick:
		el.MustClick()
	case RightMouseClick:
		err := el.Click(proto.InputMouseButtonRight, 1)
		if err != nil {
			panic(err)
		}
	case Enter:
		page.KeyActions().Press(input.Enter).MustDo()
	}

	page.MustWaitStable().MustScreenshot("screenshots/" + strconv.Itoa(*idx) + ".png")
}

func (m *Main) Run() {
	browser := rod.New().MustConnect().Trace(true).Sleeper(rod.NotFoundSleeper)
	defer browser.MustClose()

	page := browser.MustPage(m.Url)
	page.MustWindowMaximize().MustWaitStable().MustScreenshot("screenshots/base.png")

	for i, t := range m.Tests {
		RunTest(page, &t, &i)
	}

	time.Sleep(time.Second)
}

func main() {
	var tests []TestRod
	tests = append(tests, TestRod{
		Selector: "textarea[title='Search']",
		Input:    "github.com",
		Event:    Enter,
	})
	tests = append(tests, TestRod{
		Selector: "a[href='https://github.com/']",
		Input:    "",
		Event:    LeftMouseClick,
	})

	m := Main{
		Url:         "https://www.google.com/",
		ShowBrowser: true,
		Screenshots: false,
		Tests:       tests,
	}

	m.Run()
}
