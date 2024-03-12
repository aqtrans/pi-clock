[![builds.sr.ht status](https://builds.sr.ht/~aqtrans/pi-clock.svg)](https://builds.sr.ht/~aqtrans/pi-clock?)

# Description  
This turns a Raspberry Pi and it's official 7-inch display into a simple clock,  
with a Simpsons' Steamed Hams related-background, a silly little  
"Days Since Last Seizure:" counter, based off the old factory  
"Days Since Last Injury" counters.  


It utilizes SDL2 and runs on the framebuffer, avoiding any Xorg or Wayland shenanigans. 


If connected to a Pi with a SenseHat installed and running the included `sense-hat-server.py`,  
or to another server running the  `sense-hat-server.py` and the code in both `main.go`  
and `sense-hat-server.py`, configured properly (one of the TODOs is to make these values configurable)  
it will show the Temperature, Pressure, and Humidity values from the SenseHat.

## Screenshot:  
![Screenshot of the clock in action](/screenshot.jpg)


## Todos:  
- There currently seems to be a slow memory leak that I have been unable to track down,  
likely related to my SDL2 code. This is my first GUI application, so this was a very  
fun little learning experience.  

- Expose the SenseHat server address values, both in the Python and Go code.  

- Expose the "Days Since Last Seizures" value, allowing non-epileptics to count whatever days since whatever they want.  