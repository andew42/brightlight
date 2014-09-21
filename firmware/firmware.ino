/*  OctoWS2811 BrLights.ino - Elms Bedroom Lights Controller
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
    pin 2:  LED Strip #1    OctoWS2811 drives 8 LED Strips.
    pin 14: LED strip #2    All 8 are the same length.
    pin 7:  LED strip #3
    pin 8:  LED strip #4    A 100 ohm resistor should used
    pin 6:  LED strip #5    between each Teensy pin and the
    pin 20: LED strip #6    wire to the LED strip, to minimize
    pin 21: LED strip #7    high frequency ringining & noise.
    pin 5:  LED strip #8
    pin 15 & 16 - Connect together, but do not use
    pin 4 - Do not use
    pin 3 - Do not use as PWM.  Normal use is ok.

  This test is useful for checking if your LED strips work, and which
  color config (WS2811_RGB, WS2811_GRB, etc) they require.
*/

#include <OctoWS2811.h>

const int ledsPerStrip = 3;

DMAMEM int displayMemory[ledsPerStrip*6];
int drawingMemory[ledsPerStrip*6];

const int config = WS2811_GRB | WS2811_800kHz;

OctoWS2811 leds(ledsPerStrip, displayMemory, drawingMemory, config);

void setup() {
  leds.begin();
  leds.show();
  Serial.begin(9600);
}

#define RED    0xFF0000
#define GREEN  0x00FF00
#define BLUE   0x0000FF
#define YELLOW 0xFFFF00
#define PINK   0xFF1088
#define ORANGE 0xE05800
#define WHITE  0xFFFFFF

int r = 0;
int g = 0;
int b = 0;
int rgb = 0x791a8e;

// Reads a byte from the serial port, waits forever
int serial_read()
{
  if (Serial.available() > 0)
    return Serial.read();
}

void sync()
{
  // Wait for 4 0x20
  int count = 0;
  while (count < 4)
     count = (serial_read() == 0x20) ? count + 1 : 0;
}

int read_colour()
{
  return (serial_read() << 16) +
         (serial_read() << 8) +
         serial_read();
}

void loop()
{
#if 0
  Serial.println("Hello World...");
  delay(1000);  // do not print too fast!
#endif

#if 1
  sync();
  for (int i = 0; i < (ledsPerStrip * 8); i++)
    leds.setPixel(i, read_colour());
  leds.show();
#endif

#if 0 
  leds.setPixel(0, rgb);
  leds.setPixel(1, rgb);
  leds.setPixel(2, rgb);

  leds.setPixel(180, rgb);
  leds.setPixel(181, rgb);
  leds.setPixel(182, rgb);

  leds.show();
  
  delayMicroseconds(10000);
#endif  

#if 0
  leds.setPixel(1, 0);
  leds.setPixel(2, RED);
  leds.show();
  delayMicroseconds(1000000);
  leds.setPixel(1, YELLOW);
  leds.show();
  delayMicroseconds(1000000);
  leds.setPixel(2, 0);
  leds.setPixel(1, 0);
  leds.setPixel(0, GREEN);
  leds.show();
  delayMicroseconds(1000000);
  leds.setPixel(0, 0);
  leds.setPixel(1, YELLOW);
  leds.show();
  delayMicroseconds(1000000);
#endif  
}

void colorWipe(int color, int wait)
{
  for (int i=0; i < leds.numPixels(); i++) {
    leds.setPixel(i, color);
    leds.show();
    delayMicroseconds(wait);
  }
}
