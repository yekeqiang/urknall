package urknall

import (
	"testing"

	"github.com/dynport/urknall/cmd"
	. "github.com/smartystreets/goconvey/convey"
)

type customCommand struct {
	Content string
}

func (cc *customCommand) Shell() string {
	return "cc: " + cc.Content
}
func (cc *customCommand) Logging() string {
	return ""
}

type somePackage struct {
	SField string
	IField int
}

func (sp *somePackage) Package(rl *Package) {
}

func TestAddCommand(t *testing.T) {
	Convey("Given a runlist for a certain package", t, func() {
		rl := &Package{pkg: &somePackage{SField: "something", IField: 1}}

		Convey("When a string is added", func() {
			rl.Add(`string with "{{ .SField }}" and "{{ .IField }}"`)
			Convey("Then the string is turned to a command and appended to the list of commands", func() {
				c := rl.commands[len(rl.commands)-1]
				sc, ok := c.(*cmd.ShellCommand)

				Convey("And the command is a ShellCommand", func() {
					So(ok, ShouldBeTrue)
				})
				Convey("And the command template itself was expanded", func() {
					So(sc.Command, ShouldEqual, `string with "something" and "1"`)
				})
			})
		})

		Convey("Given a ShellCommand", func() {
			baseCommand := cmd.ShellCommand{Command: `string with "{{ .SField }}" and "{{ .IField }}"`}

			Convey("When it is added to the runlist by value", func() {
				f := func() { rl.Add(baseCommand) }

				Convey("Then Add will panic", func() {
					So(f, ShouldPanic)
				})
			})

			Convey("When it is added by reference", func() {
				rl.Add(&baseCommand)
				c := rl.commands[len(rl.commands)-1]
				sc, ok := c.(*cmd.ShellCommand)

				Convey("Then the command is a ShellCommand", func() {
					So(ok, ShouldBeTrue)
				})
				Convey("And the command template itself was expanded", func() {
					So(sc.Command, ShouldEqual, `string with "something" and "1"`)
				})
			})
		})

		Convey("Given a FileCommand", func() {
			baseCommand := cmd.FileCommand{Path: "/tmp/foo", Content: "{{ .SField }} = {{ .IField }}"}

			Convey("When it is added to the runlist by value", func() {
				f := func() { rl.Add(baseCommand) }

				Convey("Then Add will panic", func() {
					So(f, ShouldPanic)
				})
			})

			Convey("When it is added by reference", func() {
				rl.Add(&baseCommand)
				c := rl.commands[len(rl.commands)-1]
				fc, ok := c.(*cmd.FileCommand)

				Convey("Then the command is a ShellCommand", func() {
					So(ok, ShouldBeTrue)
				})
				Convey("And the command template itself was expanded", func() {
					So(fc.Content, ShouldEqual, `something = 1`)
				})
			})
		})

		Convey("Given a DownloadCommand is added", func() {
			baseCommand := cmd.DownloadCommand{}
			Convey("with neither Url nor destination", func() {
				So(func() { rl.Add(&baseCommand) }, ShouldPanic)
			})

			Convey("only with Url set", func() {
				baseCommand.Destination = ""
				baseCommand.Url = "some.url"
				So(func() { rl.Add(&baseCommand) }, ShouldPanic)
			})

			Convey("only with Destination set", func() {
				baseCommand.Url = ""
				baseCommand.Destination = "/opt"
				So(func() { rl.Add(&baseCommand) }, ShouldPanic)
			})

			Convey("only with Destination set", func() {
				baseCommand.Url = "some.url/{{ .SField }}"
				baseCommand.Destination = "/opt/{{ .IField }}"
				So(func() { rl.Add(&baseCommand) }, ShouldNotPanic)
				So(baseCommand.Url, ShouldEqual, "some.url/something")
				So(baseCommand.Destination, ShouldEqual, "/opt/1")
				So(rl.commands[len(rl.commands)-1], ShouldEqual, &baseCommand)
			})
		})

		Convey("Given a custom command", func() {
			baseCommand := customCommand{Content: "something {{ .NotExpanded }}"}
			Convey("When the custom command is added by value", func() {
				f := func() { rl.Add(baseCommand) }

				Convey("Then Add will panic", func() {
					So(f, ShouldPanic)
				})
			})

			Convey("When the custom command is added by reference", func() {
				rl.Add(&baseCommand)
				c := rl.commands[len(rl.commands)-1]
				cc, ok := c.(*customCommand)

				Convey("Then the command is a customCommand", func() {
					So(ok, ShouldBeTrue)
				})
				Convey("And the content wasn't touched", func() {
					So(cc.Content, ShouldEqual, `something {{ .NotExpanded }}`)
				})
			})
		})
	})
}
