#include <RunningMedian.h>

int pin = 12;

void setup()
{
  pinMode(pin, INPUT);
  Serial.begin(38400);
}

RunningMedian samples = RunningMedian(5);
int lookup[] = {665,940,1220,1500,1785,2070,2324,2574,2845,3120,3400};

void loop()
{
  // Collect samples
  for (int i = 0; i < 5; i++)
    samples.add(pulseIn(pin, HIGH));
  
  Serial.println(Lookup(samples.getMedian()));
}

int Lookup(int sample)
{
  // Lookup 1..12
  for (int i = 0; i < 11; i++)
    if (sample < lookup[i])
      return i;
   return 12;
}

