package ite8291

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Colors", func() {

	Describe("ParseColor", func() {

		var col *Color
		var err error
		var s string

		JustBeforeEach(func() {
			col, err = ParseColor(s)
		})

		DescribeTableSubtree("happy cases",
			func(rgb string, expCol *Color) {

				BeforeEach(func() {
					s = rgb
				})

				It("parsed correctly", func() {
					Ω(col).Should(Equal(expCol))
					Ω(err).Should(Succeed())
				})
			},

			EntryDescription("color as %[1]q"),

			Entry(nil, "0X000000", &Color{Red: 0, Green: 0, Blue: 0}),
			Entry(nil, "0X192837", &Color{Red: 0x19, Green: 0x28, Blue: 0x37}),
			Entry(nil, "0XAFBECD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "0Xafbecd", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "0XaFbeCD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "0Xffffff", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
			Entry(nil, "0XFFFFFF", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),

			Entry(nil, "0x000000", &Color{Red: 0, Green: 0, Blue: 0}),
			Entry(nil, "0x192837", &Color{Red: 0x19, Green: 0x28, Blue: 0x37}),
			Entry(nil, "0xAFBECD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "0xafbecd", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "0xaFbeCD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "0xffffff", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
			Entry(nil, "0xFFFFFF", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),

			Entry(nil, "#X000000", &Color{Red: 0, Green: 0, Blue: 0}),
			Entry(nil, "#X192837", &Color{Red: 0x19, Green: 0x28, Blue: 0x37}),
			Entry(nil, "#XAFBECD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "#Xafbecd", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "#XaFbeCD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "#Xffffff", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
			Entry(nil, "#XFFFFFF", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),

			Entry(nil, "#x000000", &Color{Red: 0, Green: 0, Blue: 0}),
			Entry(nil, "#x192837", &Color{Red: 0x19, Green: 0x28, Blue: 0x37}),
			Entry(nil, "#xAFBECD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "#xafbecd", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "#xaFbeCD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "#xffffff", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
			Entry(nil, "#xFFFFFF", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),

			Entry(nil, "#000000", &Color{Red: 0, Green: 0, Blue: 0}),
			Entry(nil, "#192837", &Color{Red: 0x19, Green: 0x28, Blue: 0x37}),
			Entry(nil, "#AFBECD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "#afbecd", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "#aFbeCD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "#ffffff", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
			Entry(nil, "#FFFFFF", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),

			Entry(nil, "000000", &Color{Red: 0, Green: 0, Blue: 0}),
			Entry(nil, "192837", &Color{Red: 0x19, Green: 0x28, Blue: 0x37}),
			Entry(nil, "AFBECD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "afbecd", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "aFbeCD", &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
			Entry(nil, "ffffff", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
			Entry(nil, "FFFFFF", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),

			Entry(nil, "#000", &Color{Red: 0, Green: 0, Blue: 0}),
			Entry(nil, "#183", &Color{Red: 0x11, Green: 0x88, Blue: 0x33}),
			Entry(nil, "#FBD", &Color{Red: 0xff, Green: 0xbb, Blue: 0xdd}),
			Entry(nil, "#aec", &Color{Red: 0xaa, Green: 0xee, Blue: 0xcc}),
			Entry(nil, "#bCD", &Color{Red: 0xbb, Green: 0xcc, Blue: 0xdd}),
			Entry(nil, "#fff", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
			Entry(nil, "#FFF", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),

			Entry(nil, "000", &Color{Red: 0, Green: 0, Blue: 0}),
			Entry(nil, "183", &Color{Red: 0x11, Green: 0x88, Blue: 0x33}),
			Entry(nil, "FBD", &Color{Red: 0xff, Green: 0xbb, Blue: 0xdd}),
			Entry(nil, "aec", &Color{Red: 0xaa, Green: 0xee, Blue: 0xcc}),
			Entry(nil, "bCD", &Color{Red: 0xbb, Green: 0xcc, Blue: 0xdd}),
			Entry(nil, "fff", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
			Entry(nil, "FFF", &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
		)

		DescribeTableSubtree("validation error",
			func(rgb string) {

				BeforeEach(func() {
					s = rgb
				})

				It("returns a corresponding error", func() {
					Ω(col).Should(BeNil())
					Ω(err).Should(MatchError(ContainSubstring(s)))
				})
			},

			EntryDescription("color as %[1]q"),

			Entry(nil, "0X0000000"),
			Entry(nil, "0X00000"),
			Entry(nil, "0X00000."),
			Entry(nil, "0X0000GH"),
			Entry(nil, "0X000I12"),
			Entry(nil, "0X00OA12"),
			Entry(nil, "0X0-Ac12"),

			Entry(nil, "0xBBBBBBB"),
			Entry(nil, "0xBBBBB"),
			Entry(nil, "0xBBBBB."),
			Entry(nil, "0xBBBBGH"),
			Entry(nil, "0xBBBI12"),
			Entry(nil, "0xBBOA12"),
			Entry(nil, "0x0B-Ac12"),

			Entry(nil, "#xBBBBBBB"),
			Entry(nil, "#xBBBBB"),
			Entry(nil, "#xBBBBB."),
			Entry(nil, "#xBBBBGH"),
			Entry(nil, "#xBBBI12"),
			Entry(nil, "#xBBOA12"),
			Entry(nil, "#x0B-Ac12"),

			Entry(nil, "#BBBBBBB"),
			Entry(nil, "#BBBBB"),
			Entry(nil, "#BBBBB."),
			Entry(nil, "#BBBBGH"),
			Entry(nil, "#BBBI12"),
			Entry(nil, "#BBOA12"),
			Entry(nil, "#0B-Ac12"),

			Entry(nil, "BBBBBBB"),
			Entry(nil, "BBBBB"),
			Entry(nil, "BBBBB."),
			Entry(nil, "BBBBGH"),
			Entry(nil, "BBBI12"),
			Entry(nil, "BBOA12"),
			Entry(nil, "0B-Ac12"),

			Entry(nil, "1F!"),
			Entry(nil, "ay0"),
			Entry(nil, "O00."),

			Entry(nil, "#1F!"),
			Entry(nil, "#ay0"),
			Entry(nil, "#O00."),
		)
	})

	DescribeTableSubtree("String",
		func(col *Color, expS string) {

			It("returns correct string representation", func() {
				Ω(col.String()).Should(Equal(expS))
			})
		},

		EntryDescription("for color %[1]q"),

		Entry(nil, &Color{Red: 0, Green: 0, Blue: 0}, "#000000"),
		Entry(nil, &Color{Red: 0x19, Green: 0x28, Blue: 0x37}, "#192837"),
		Entry(nil, &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}, "#AFBECD"),
		Entry(nil, &Color{Red: 0xff, Green: 0xff, Blue: 0xff}, "#FFFFFF"),
	)

	DescribeTableSubtree("FromRGB",
		func(rgb uint32, expCol *Color) {

			It("returns correct color", func() {
				Ω(FromRGB(rgb)).Should(Equal(expCol))
			})
		},

		EntryDescription("for color %#06[1]x"),

		Entry(nil, uint32(0x0), &Color{Red: 0, Green: 0, Blue: 0}),
		Entry(nil, uint32(0x192837), &Color{Red: 0x19, Green: 0x28, Blue: 0x37}),
		Entry(nil, uint32(0xAFBECD), &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
		Entry(nil, uint32(0xffffff), &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),

		Entry(nil, uint32(0x1000000), &Color{Red: 0, Green: 0, Blue: 0}),
		Entry(nil, uint32(0xA192837), &Color{Red: 0x19, Green: 0x28, Blue: 0x37}),
		Entry(nil, uint32(0x10AFBECD), &Color{Red: 0xaf, Green: 0xbe, Blue: 0xcd}),
		Entry(nil, uint32(0xffffffff), &Color{Red: 0xff, Green: 0xff, Blue: 0xff}),
	)
})
