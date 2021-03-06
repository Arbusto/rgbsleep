package leds

/*
#include "AURA_SDK.h"

static EnumerateMbControllerFunc EnumerateMbController;
static SetMbModeFunc SetMbMode;
static SetMbColorFunc SetMbColor;
static GetMbLedCountFunc GetMbLedCount;

static MbLightControl handle;
static DWORD ledCount;
static BYTE *colors;
static BYTE ready;


static void initDLL() {
	HMODULE dll = LoadLibrary("AURA_SDK.dll");
	EnumerateMbController = (EnumerateMbControllerFunc)GetProcAddress(dll, "EnumerateMbController");
	SetMbMode = (SetMbModeFunc)GetProcAddress(dll, "SetMbMode");
	SetMbColor = (SetMbColorFunc)GetProcAddress(dll, "SetMbColor");
	GetMbLedCount = (GetMbLedCountFunc)GetProcAddress(dll, "GetMbLedCount");
}

void InitAura() {
	initDLL();
	DWORD count = EnumerateMbController(NULL, 0);
	if (count != 1) {
		return;
	}
	MbLightControl *handles = calloc(count, sizeof(MbLightControl));
	EnumerateMbController(handles, count);
	handle = handles[0];
	SetMbMode(handle, 1);
	ledCount = GetMbLedCount(handle);
	colors = calloc(3 * ledCount, sizeof(BYTE));
	ready = 1;
}

void SetAuraColors(BYTE boardR, BYTE boardG, BYTE boardB, BYTE caseR, BYTE caseG, BYTE caseB) {
	if (!ready) {
		return;
	}
	for (int i = 0; i < ledCount; i++) {
		BYTE r, g, b;
		if (i == ledCount - 1) {
			// rgb header for front panel controller
			r = caseR;
			g = caseG;
			b = caseB;
		} else if (i == ledCount - 2) {
			// unused rgb header
			r = 0;
			g = 0;
			b = 0;
		} else {
			// onboard led
			r = boardR;
			g = boardG;
			b = boardB;
		}
		colors[i*3 + 0] = r;
		colors[i*3 + 1] = b;
		colors[i*3 + 2] = g;
	}
	SetMbColor(handle, colors, ledCount * 3 * sizeof(BYTE));
}
*/
import "C"

import (
	"image/color"
	"time"
)

var (
	colorOff     = color.RGBA{}
	colorOnCase  = color.RGBA{R: 255, G: 50, B: 0}
	colorOnBoard = color.RGBA{R: 255, G: 185, B: 15}
)

func fade(from, to color.RGBA, progress uint8, reverse bool) color.RGBA {
	if reverse {
		progress = 100 - progress
	}
	return color.RGBA{
		R: uint8(float32(from.R) + float32(progress)/100*float32(to.R-from.R)),
		G: uint8(float32(from.G) + float32(progress)/100*float32(to.G-from.G)),
		B: uint8(float32(from.B) + float32(progress)/100*float32(to.B-from.B)),
	}
}

func InitAura() {
	C.InitAura()
}

func TurnAuraOn() {
	C.SetAuraColors(
		C.uchar(colorOnBoard.R), C.uchar(colorOnBoard.G), C.uchar(colorOnBoard.B),
		C.uchar(colorOnCase.R), C.uchar(colorOnCase.G), C.uchar(colorOnCase.B),
	)
}

func RunAuraFader(stateChan chan bool) {
	for {
		var reverse bool
		if <-stateChan {
			reverse = false
		} else {
			reverse = true
		}
		for i := uint8(0); i <= 100; i++ {
			colorCase := fade(colorOff, colorOnCase, i, reverse)
			colorBoard := fade(colorOff, colorOnBoard, i, reverse)
			C.SetAuraColors(
				C.uchar(colorBoard.R), C.uchar(colorBoard.G), C.uchar(colorBoard.B),
				C.uchar(colorCase.R), C.uchar(colorCase.G), C.uchar(colorCase.B),
			)
			time.Sleep(10 * time.Millisecond)
		}
	}
}
