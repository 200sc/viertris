package viertris

import (
	"image"
	"image/color"
	"image/draw"
	"strconv"

	"github.com/oakmound/oak/v3/alg/intgeom"
	"github.com/oakmound/oak/v3/render"
)

func drawRect(buff *image.RGBA, pos, dims intgeom.Point2, c color.RGBA) {
	render.DrawLine(buff,
		pos.X(), pos.Y(),
		pos.X()+dims.X(), pos.Y(),
		c)
	render.DrawLine(buff,
		pos.X()+dims.X(), pos.Y(),
		pos.X()+dims.X(), pos.Y()+dims.Y(),
		c)
	render.DrawLine(buff,
		pos.X()+dims.X(), pos.Y()+dims.Y(),
		pos.X(), pos.Y()+dims.Y(),
		c)
	render.DrawLine(buff,
		pos.X(), pos.Y()+dims.Y(),
		pos.X(), pos.Y(),
		c)
}

func drawFilledRect(buff *image.RGBA, pos, dims intgeom.Point2, c color.RGBA) {
	draw.Draw(buff, image.Rect(pos.X(), pos.Y(), pos.X()+dims.X(), pos.Y()+dims.Y()),
		image.NewUniform(c), image.Point{}, draw.Over)
}

type stringerUint64Pointer struct {
	prefix string
	v      *uint64
}

func (sip stringerUint64Pointer) String() string {
	return sip.prefix + strconv.FormatUint(*sip.v, 10)
}

func NewUint64Text(prefix string, f *render.Font, val *uint64) *render.Text {
	return f.NewStringerText(stringerUint64Pointer{prefix: prefix, v: val}, 0, 0)
}
