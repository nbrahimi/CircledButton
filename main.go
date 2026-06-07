package main

import (
	"fmt"
	"image/color"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"v14.3/circledbutton"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Custom circled button"))
		w.Option(app.Size(unit.Dp(700), unit.Dp(600)))
		if err := loop(w); err != nil {
			os.Exit(1)
		}
	}()
	app.Main()
}

func loop(w *app.Window) error {
	th := material.NewTheme()

	btn := circledbutton.CircledButton{
		Label:      "Click me",
		TxtColor:   color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		BkColor:    color.NRGBA{R: 12, G: 44, B: 204, A: 63},
		HoverColor: color.NRGBA{R: 52, G: 94, B: 214, A: 103},
		PressColor: color.NRGBA{R: 92, G: 134, B: 254, A: 203},
		ClickColor: color.NRGBA{R: 132, G: 174, B: 254, A: 243},
		Size:       unit.Dp(200),
	}
	str := "Click, Press or Hover the button"
	labelColor := color.NRGBA{B: 1, R: 1, G: 1, A: 255}
	bk := btn.BkColor

	var ops op.Ops
	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			os.Exit(0)
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			switch {
			case btn.IsClicked:
				btn.ClickedAt = time.Now()
				str = fmt.Sprintf("Clicked (%dx)... ", btn.NumClicks)
				labelColor = color.NRGBA{G: 255, A: 255}
				btn.BkColor = btn.ClickColor
				w.Invalidate()
			case btn.IsCancelled:
				btn.ClickedAt = time.Now()
				str = "Cancelled ... "
				labelColor = color.NRGBA{R: 255, B: 255, A: 255}
				w.Invalidate()
			case !btn.ClickedAt.IsZero() && time.Since(btn.ClickedAt) < 2*time.Second:
				w.Invalidate()
			case btn.IsPressed:
				btn.ClickedAt = time.Now()
				btn.BkColor = btn.PressColor
				str = "Pressed ... "
				labelColor = color.NRGBA{R: 255, A: 255}
				w.Invalidate()
			case btn.IsHovered:
				btn.BkColor = btn.HoverColor
				str = "Hovered ... "
				labelColor = color.NRGBA{R: 150, G: 150, B: 155, A: 255}
				w.Invalidate()
			default:
				str = "Click, Press or Hover the button"
				labelColor = color.NRGBA{B: 1, R: 1, G: 1, A: 255}
				btn.ClickedAt = time.Time{}
				btn.BkColor = bk
				w.Invalidate()
			}

			layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						lbl := material.Body1(th, str)
						lbl.Alignment = text.Middle
						lbl.TextSize = unit.Sp(16)
						lbl.Color = labelColor
						return lbl.Layout(gtx)
					})
				}),
				layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
				layout.Rigid(func(gtx C) D {
					return layout.Center.Layout(gtx, func(gtx C) D {
						return btn.Layout(gtx, th)
					})
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}
