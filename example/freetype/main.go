// Copyright 2010 The Freetype-Go Authors. All rights reserved.
// Use of this source code is governed by your choice of either the
// FreeType License or the GNU General Public License version 2 (or
// any later version), both of which can be found in the LICENSE file.

package main

import (
	"bufio"
	"exp/draw"
	"flag"
	"fmt"
	"freetype-go.googlecode.com/hg/freetype"
	"freetype-go.googlecode.com/hg/freetype/raster"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

var (
	dpi      = flag.Int("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "../../luxi-fonts/luxisr.ttf", "filename of the ttf font")
	gamma    = flag.Float("gamma", 1.0, "gamma correction")
	size     = flag.Float("size", 12, "font size in points")
	spacing  = flag.Float("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
)

var text = []string{
	"’Twas brillig, and the slithy toves",
	"Did gyre and gimble in the wabe;",
	"All mimsy were the borogoves,",
	"And the mome raths outgrabe.",
	"",
	"“Beware the Jabberwock, my son!",
	"The jaws that bite, the claws that catch!",
	"Beware the Jubjub bird, and shun",
	"The frumious Bandersnatch!”",
	"",
	"He took his vorpal sword in hand:",
	"Long time the manxome foe he sought—",
	"So rested he by the Tumtum tree,",
	"And stood awhile in thought.",
	"",
	"And as in uffish thought he stood,",
	"The Jabberwock, with eyes of flame,",
	"Came whiffling through the tulgey wood,",
	"And burbled as it came!",
	"",
	"One, two! One, two! and through and through",
	"The vorpal blade went snicker-snack!",
	"He left it dead, and with its head",
	"He went galumphing back.",
	"",
	"“And hast thou slain the Jabberwock?",
	"Come to my arms, my beamish boy!",
	"O frabjous day! Callooh! Callay!”",
	"He chortled in his joy.",
	"",
	"’Twas brillig, and the slithy toves",
	"Did gyre and gimble in the wabe;",
	"All mimsy were the borogoves,",
	"And the mome raths outgrabe.",
}

func main() {
	flag.Parse()

	// Read the font data.
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Stderr(err)
		return
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Stderr(err)
		return
	}

	// Initialize the context.
	fg, bg := image.Black, image.White
	ruler := image.RGBAColor{0xdd, 0xdd, 0xdd, 0xff}
	if *wonb {
		fg, bg = image.White, image.Black
		ruler = image.RGBAColor{0x22, 0x22, 0x22, 0xff}
	}
	rgba := image.NewRGBA(640, 480)
	draw.Draw(rgba, draw.Rect(0, 0, rgba.Width(), rgba.Height()), bg, draw.ZP)
	c := freetype.NewContext()
	c.SetDPI(*dpi)
	c.SetFont(font)
	c.SetFontSize(*size)
	rp := raster.NewRGBAPainter(rgba)
	rp.SetColor(fg)
	gcp := raster.NewGammaCorrectionPainter(rp, *gamma)

	// Draw the guidelines.
	for i := 0; i < 200; i++ {
		rgba.Set(10, 10+i, ruler)
		rgba.Set(10+i, 10, ruler)
	}

	// Draw the text.
	pt := freetype.Pt(10, 10)
	for _, s := range text {
		err = c.DrawText(gcp, pt, s)
		if err != nil {
			log.Stderr(err)
			return
		}
		pt.Y += c.PointToFix32(*size * *spacing)
	}

	// Save that RGBA image to disk.
	f, err := os.Open("out.png", os.O_CREAT|os.O_WRONLY, 0600)
	if err != nil {
		log.Stderr(err)
		os.Exit(1)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Stderr(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Stderr(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}
