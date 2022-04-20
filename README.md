# GoHoroscope

A bit of astrology in terminal with Go. Prints Ascendant, Houses, Zodiac signs, Aspects, Positions, Solar and Lunar eclipses, Retrograde movements.

## Installation

```
- git clone
- go build -o ./bin/horoscope .
```
This project is using Swiss Ephemeris with [swephgo](https://github.com/mshafiee/swephgo) as wrapper. You will also need Swiss Ephemem library and ephemeris files.

- 

Download the Swiss Ephemeris Library [here](https://www.astro.com/ftp/swisseph/). After compiling the library, copy the libswe.so file to /usr/local/lib/

````
$ cp libswe.so /usr/local/lib/
````

Finally you will need some shell variables

```
export LATITUDE="51.5072"¬
export LONGITUDE="0.1276"¬
export CITY="Europe/London"
```

or put variables in `.env` file looking like `.env.example`

