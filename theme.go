package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Exact colors pulled from the TypeScript project's Tailwind classes
// (slate/emerald/amber/indigo palette) rather than eyeballed from screenshots.
var (
	// Backgrounds
	colorAppBG    = color.NRGBA{R: 0x07, G: 0x09, B: 0x13, A: 0xFF} // #070913
	colorHeaderBG = color.NRGBA{R: 0x0D, G: 0x11, B: 0x24, A: 0xFF} // #0d1124
	colorPanelBG  = color.NRGBA{R: 0x09, G: 0x0D, B: 0x1A, A: 0xFF} // #090d1a
	colorKeycapBG = color.NRGBA{R: 0x0A, G: 0x0D, B: 0x1A, A: 0xFF} // #0a0d1a

	// Slate scale
	colorSlate950 = color.NRGBA{R: 0x02, G: 0x06, B: 0x17, A: 0xFF}
	colorSlate900 = color.NRGBA{R: 0x0F, G: 0x17, B: 0x2A, A: 0xFF}
	colorSlate800 = color.NRGBA{R: 0x1E, G: 0x29, B: 0x3B, A: 0xFF}
	colorSlate700 = color.NRGBA{R: 0x33, G: 0x41, B: 0x55, A: 0xFF}
	colorSlate600 = color.NRGBA{R: 0x47, G: 0x55, B: 0x69, A: 0xFF}
	colorMuted    = color.NRGBA{R: 0x64, G: 0x74, B: 0x8B, A: 0xFF} // slate-500
	colorSlate400 = color.NRGBA{R: 0x94, G: 0xA3, B: 0xB8, A: 0xFF}
	colorSlate200 = color.NRGBA{R: 0xE2, G: 0xE8, B: 0xF0, A: 0xFF}
	colorWhite    = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}

	// Accents
	colorGreen        = color.NRGBA{R: 0x34, G: 0xD3, B: 0x99, A: 0xFF} // emerald-400
	colorGreenStrong  = color.NRGBA{R: 0x10, G: 0xB9, B: 0x81, A: 0xFF} // emerald-500
	colorOrange       = color.NRGBA{R: 0xFB, G: 0xBF, B: 0x24, A: 0xFF} // amber-400
	colorOrangeStrong = color.NRGBA{R: 0xF5, G: 0x9E, B: 0x0B, A: 0xFF} // amber-500
	colorIndigo       = color.NRGBA{R: 0x81, G: 0x8C, B: 0xF8, A: 0xFF} // indigo-400
	colorIndigoStrong = color.NRGBA{R: 0x63, G: 0x66, B: 0xF1, A: 0xFF} // indigo-500

	// Card backgrounds
	colorCardBG           = color.NRGBA{R: 0x10, G: 0x14, B: 0x26, A: 0xFF} // #101426 completed entry
	colorCardBGActive     = color.NRGBA{R: 0x0F, G: 0x1D, B: 0x19, A: 0xFF} // #0f1d19 in-progress entry
	colorCardBorder       = colorSlate800
	colorCardBorderActive = color.NRGBA{R: 0x05, G: 0x96, B: 0x69, A: 0x4D} // emerald-600/30
	colorDivider          = color.NRGBA{R: 0x1E, G: 0x29, B: 0x3B, A: 0x80} // slate-800/50
)

// appTheme is a custom Fyne theme matching the source TypeScript project's
// dark navy / emerald / amber / indigo palette and its Inter + JetBrains
// Mono font pairing.
type appTheme struct{}

var _ fyne.Theme = (*appTheme)(nil)

func (t *appTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNamePrimary:
		return colorGreenStrong
	case theme.ColorNameBackground:
		return colorAppBG
	case theme.ColorNameButton:
		return colorSlate900
	case theme.ColorNameHover:
		return colorSlate800
	case theme.ColorNameSeparator:
		return colorDivider
	case theme.ColorNameForeground:
		return colorSlate200
	case theme.ColorNameDisabled:
		return colorSlate600
	case theme.ColorNameInputBackground:
		return colorSlate900
	case theme.ColorNameInputBorder:
		return colorSlate800
	case theme.ColorNameSuccess:
		return colorGreenStrong
	case theme.ColorNameWarning:
		return colorOrangeStrong
	}
	return theme.DarkTheme().Color(name, variant)
}

// Font maps to the embedded JetBrains Mono (for TextStyle.Monospace) or
// Inter (default sans) faces, in Regular or Bold depending on TextStyle.Bold.
// Note: Fyne 2.4's theme API only distinguishes Bold/Monospace/Italic, not
// arbitrary weights - so the source project's Medium/SemiBold weights can't
// be represented at the theme level in this Fyne version.
func (t *appTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		if style.Bold {
			return fontMonoBold
		}
		return fontMonoRegular
	}
	if style.Bold {
		return fontSansBold
	}
	return fontSansRegular
}

func (t *appTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DarkTheme().Icon(name)
}

func (t *appTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameInputRadius, theme.SizeNameSelectionRadius:
		return 10
	}
	return theme.DarkTheme().Size(name)
}
