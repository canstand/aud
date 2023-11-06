package aud

import (
	"time"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astisub"
)

var (
	// SSADefaultStyle is the default style when export ssa subtitle
	SSADefaultStyle = &astisub.Style{
		ID: "Default",
		InlineStyle: &astisub.StyleAttributes{
			SSAFontName:        "Arial",
			SSAFontSize:        astikit.Float64Ptr(20),
			SSAPrimaryColour:   &astisub.Color{Red: 255, Green: 252, Blue: 3},
			SSASecondaryColour: astisub.ColorBlack,
			SSAOutlineColour:   astisub.ColorBlack,
			SSABackColour:      astisub.ColorBlack,
			SSABold:            astikit.BoolPtr(false),
			SSAItalic:          astikit.BoolPtr(false),
			SSAUnderline:       astikit.BoolPtr(false),
			SSAStrikeout:       astikit.BoolPtr(false),
			SSAScaleX:          astikit.Float64Ptr(100),
			SSAScaleY:          astikit.Float64Ptr(100),
			SSAAngle:           astikit.Float64Ptr(0),
			SSAShadow:          astikit.Float64Ptr(0),
			SSASpacing:         astikit.Float64Ptr(1),
			SSABorderStyle:     astikit.IntPtr(1),
			SSAOutline:         astikit.Float64Ptr(0.5),
			SSAAlignment:       astikit.IntPtr(2),
			SSAMarginLeft:      astikit.IntPtr(80),
			SSAMarginRight:     astikit.IntPtr(80),
			SSAMarginVertical:  astikit.IntPtr(16),
			SSAEncoding:        astikit.IntPtr(1),
		},
	}
	// SSASecondaryStyle is default style of secondary language when export bilingual ssa subtitle
	SSASecondaryStyle = &astisub.Style{
		ID: "Secondary",
		InlineStyle: &astisub.StyleAttributes{
			SSAFontName:        "Arial",
			SSAFontSize:        astikit.Float64Ptr(12),
			SSAPrimaryColour:   astisub.ColorWhite,
			SSASecondaryColour: astisub.ColorBlack,
			SSAOutlineColour:   astisub.ColorBlack,
			SSABackColour:      astisub.ColorBlack,
			SSABold:            astikit.BoolPtr(false),
			SSAItalic:          astikit.BoolPtr(false),
			SSAUnderline:       astikit.BoolPtr(false),
			SSAStrikeout:       astikit.BoolPtr(false),
			SSAScaleX:          astikit.Float64Ptr(100),
			SSAScaleY:          astikit.Float64Ptr(100),
			SSAAngle:           astikit.Float64Ptr(0),
			SSAShadow:          astikit.Float64Ptr(0),
			SSASpacing:         astikit.Float64Ptr(0),
			SSABorderStyle:     astikit.IntPtr(1),
			SSAOutline:         astikit.Float64Ptr(0.5),
			SSAAlignment:       astikit.IntPtr(2),
			SSAMarginLeft:      astikit.IntPtr(80),
			SSAMarginRight:     astikit.IntPtr(80),
			SSAMarginVertical:  astikit.IntPtr(16),
			SSAEncoding:        astikit.IntPtr(1),
		},
	}
	ssaDefaultMetadata = &astisub.Metadata{
		SSACollisions: "Reverse",
		SSAWrapStyle:  "0",
		// 16:9: 640x360, marginLeft 80, marginRight 80, marginVertical 16
		// 4:3: 480x360, marginLeft 60, marginRight 60, marginVertical 16
		SSAPlayResX:   astikit.IntPtr(640),
		SSAPlayResY:   astikit.IntPtr(360),
		SSAScriptType: "v4.00+",
	}
)

// optimizeIntervals optimize subtitle's intervals
func optimizeIntervals(subtitle *astisub.Subtitles) {
	for index, item := range subtitle.Items {
		if index == 0 || item.StartAt == subtitle.Items[index-1].EndAt {
			continue
		}
		diff := item.StartAt - subtitle.Items[index-1].EndAt
		if diff <= 500*time.Millisecond {
			subtitle.Items[index-1].EndAt = item.StartAt
		} else if diff <= 1800*time.Millisecond {
			item.StartAt = item.StartAt - 300*time.Millisecond
			subtitle.Items[index-1].EndAt = item.StartAt - 300*time.Millisecond
		} else if diff > 1800*time.Millisecond {
			item.StartAt = item.StartAt - 300*time.Millisecond
			subtitle.Items[index-1].EndAt = subtitle.Items[index-1].EndAt + 1500*time.Millisecond
		}
	}
}
