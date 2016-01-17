/*
Bright Light Teensy 3.x firmware
Copyright (c) 2014-2016 Andrew Hartland, Leading Edge Designs Ltd

Requires library:
http://www.pjrc.com/teensy/td_libs_OctoWS2811.html
Copyright (c) 2013 Paul Stoffregen, PJRC.COM, LLC

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

Required Connections
--------------------
pin 2:  LED Strip #1  OctoWS2811 drives 8 LED Strips.
pin 14: LED strip #2  All 8 are the same length.
pin 7:  LED strip #3
pin 8:  LED strip #4  A 100 ohm resistor should used
pin 6:  LED strip #5  between each Teensy pin and the
pin 20: LED strip #6  wire to the LED strip, to minimize
pin 21: LED strip #7  high frequency ringing & noise.
pin 5:  LED strip #8
pin 15 & 16 - Connect together, but do not use
pin 4 - Do not use
pin 3 - Do not use as PWM.  Normal use is ok.
*/

#include <OctoWS2811.h>

// Double buffering uses twice the memory but allows serial download
// concurrently with LED update so allows for higher frame rates.
// However by carefully arranging the download to fill already sent
// LEDs first we can update a single buffer while it is still being
// sent allowing single buffer concurrency providing timing of serial
// download is carefully tested to ensure serial download does not
// complete while DMA is still in progress.
// #define DOUBLE_BUFFERED

// MUST be synchronised with raspberry pi controller
const int LEDS_PER_STRIP = 300;

// Display buffer(s)
DMAMEM int displayMemory[LEDS_PER_STRIP * 6];
#ifdef DOUBLE_BUFFERED
int drawingMemory[LEDS_PER_STRIP * 6];
#endif // DOUBLE_BUFFERED

// Run at 800Khz, 400kHz has random flicker
const int config = WS2811_GRB | WS2811_800kHz;

// Define the LED controller class
#ifdef DOUBLE_BUFFERED
OctoWS2811 leds(LEDS_PER_STRIP, displayMemory, drawingMemory, config);
#else
OctoWS2811 leds(LEDS_PER_STRIP, displayMemory, NULL, config);
#endif // DOUBLE_BUFFERED

// Sets entire frame buffer to specified colour and updates LEDS
void set_all_leds(int colour) {
  for (int i = 0; i < (LEDS_PER_STRIP * 8); i++)
    leds.setPixel(i, colour);
  leds.show();
}

// Reads a byte from the serial port, fail after a little while (-1)
int serial_read() {
  // Tried using elapsedMillis here but it seems to take an
  // age to return, slowing everything down. Went with count
  // of reties instead. The count should be larger than the
  // frame rate otherwise LEDs will be switched off mid frame
  // 1,000,000 is around 300ms at 96Mhz
  for (int i = 0; i < 1000000; i++) {
    if (Serial.available()) {
      return Serial.read();
    }
  }
  return -1;
}

// Look for four 0xff not possible in frame data
int sync() {
  int count = 0;
  while (count < 4) {
    int c = serial_read();
    if (c == -1) return 0;
    count = (c == 0xff) ? count + 1 : 0;
  }
  return 1;
}

// Colours are transported as four bytes with 0x00 initial value
int read_colour() {
  if (serial_read() != 0) return -1;

  int red = serial_read();
  if (red == -1) return -1;

  int green = serial_read();
  if (green == -1) return -1;

  int blue = serial_read();
  if (blue == -1) return -1;

  return (red << 16) + (green << 8) + blue;
}

// Read int32 MSB first
int read_checksum() {

  int v = 0;
  for (int i = 0; i < 4; i++) {
    int b = serial_read();
    if (b == -1) return -1;
    v = (v << 8) + b;
  }
  return v;
}

// Called on reset
void setup() {

  leds.begin();
  set_all_leds(0);

  // 9600 is arbitrary, always runs at 12Mb on Teensy
  Serial.begin(9600);

  // Teensy status LED (on when we are seeing frames)
  pinMode(13, OUTPUT);

  // Testing: oscilloscope trigger for serial download
  // OctWS2811.cpp is modified to drive pin 19 for dma update
  pinMode(17, OUTPUT);
}

// Called repeatedly from main()
void loop() {

  // Wait for a frame from the Pi
  if (!sync()) {
    // No Pi, turn off all LEDs (including Teensy status)
    digitalWriteFast(13, LOW);
    set_all_leds(0);
    return;
  }

  // Turn on Teensy status LED when we have a frame
  digitalWriteFast(13, HIGH);

  // Drive scope timing trigger (start serial download)
  digitalWriteFast(17, LOW);

  // Download frame and update drawing buffer, the update
  // can start while we are still refreshing the last frame
  // via DMA provided the DMA refresh completes BEFORE the
  // serial download completes
  int checksum = 0;
  for (int i = 0; i < LEDS_PER_STRIP; i++) {
    // Buffer is send 8*LED1, 8*LED2 ... 8*(LEDS_PER_STRIP - 1)
    for (int j = 0; j < 8; j++) {
      int colour = read_colour();
      if (colour == -1) return;
      leds.setPixel(i + (j * LEDS_PER_STRIP), colour);
      checksum += colour;
    }
  }

  // Check checksum, don't update if its bad and flash led
  int cs = read_checksum();
  if (cs != checksum) {
    digitalWriteFast(13, LOW);
    delay(100);
    return;
  }

  // Drive scope timing trigger (end serial download)
  digitalWriteFast(17, HIGH);

  // Update the LEDs
  leds.show();
}
