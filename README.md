# GoHoroscope

A bit of astrology in terminal with Go. 

Prints Ascendant, Houses, Zodiac signs, Aspects, Positions, Solar and Lunar eclipses, Retrograde movements.

## Installation

```
- git clone
- go build -o ./bin/horoscope .
```
This project is using Swiss Ephemeris with [swephgo](https://github.com/mshafiee/swephgo) as a wrapper. You will also need Swiss Ephemem library and ephemeris files. 

Download the Swiss Ephemeris Library [here](https://www.astro.com/ftp/swisseph/). After compiling the library, copy the libswe.so file to /usr/local/lib/

````
$ cp libswe.so /usr/local/lib/
````

[Download some ephemeris files](https://www.astro.com/ftp/swisseph/ephe/) and puth them in designated folder. _sepl_18.se1_ is a good place to start.

Finally you will need some shell variables (timezone, geographical position, location of Swiss Ephemeris files on your system)

```
export CITY="Europe/London"
export LATITUDE="51.5072"¬
export LONGITUDE="0.1276"¬
export SWISSPATH="/usr/local/share/sweph/ephe"
```

or put variables in `.env` file looking like `.env.example`

## Usage

```
horoscope --horoscope [date]

20 Apr 22 07:08 UTC - lat: 52.43, lon: 20.89
Ascendant: 121.96 MC: 7.01

+-------+----------+-------+-------------+
| House | Position | Cusp  | Sign        |
+-------+----------+-------+-------------+
| I     | 121.96   | 1.96  | Leo         |
| II    | 138.02   | 18.02 | Leo         |
| III   | 158.49   | 8.49  | Virgo       |
| IV    | 187.01   | 7.01  | Libra       |
| V     | 226.61   | 16.61 | Scorpio     |
| VI    | 269.36   | 29.36 | Sagittarius |
| VII   | 301.96   | 1.96  | Aquarius    |
| VIII  | 318.02   | 18.02 | Aquarius    |
| IX    | 338.49   | 8.49  | Pisces      |
| X     | 7.01     | 7.01  | Aries       |
| XI    | 46.61    | 16.61 | Taurus      |
| XII   | 89.36    | 29.36 | Gemini      |
+-------+----------+-------+-------------+


+---------+----------+-------+-------------+-------------------------------------------+
| Planet  | Position | House | Sign        | Aspects                                   |
+---------+----------+-------+-------------+-------------------------------------------+
| Sun     | 30.27    | X     | Taurus      |
|         |          |       |             |    Sextile - Mars - 333.97 in Pisces         |
|         |          |       |             |    Square - Pluto - 298.58 in Capricorn      |
| Moon    | 258.69   | V     | Sagittarius |
|         |          |       |             |    Quincunx - Mercury - 47.56 in Taurus      |
|         |          |       |             |    Square - Venus - 346.13 in Pisces         |
|         |          |       |             |    Square - Neptune - 354.24 in Pisces       |
| Mercury | 47.56    | XI    | Taurus      |
|         |          |       |             |    Sextile - Venus - 346.13 in Pisces        |
|         |          |       |             |    Conjunction - Uranus - 43.94 in Taurus    |
| Venus   | 346.13   | I     | Pisces      |
|         |          |       |             |    Conjunction - Jupiter - 355.71 in Pisces  |
|         |          |       |             |    Sextile - Uranus - 43.94 in Taurus        |
|         |          |       |             |    Conjunction - Neptune - 354.24 in Pisces  |
| Mars    | 333.97   | VIII  | Pisces      |
| Jupiter | 355.71   | I     | Pisces      |
|         |          |       |             |    Conjunction - Neptune - 354.24 in Pisces  |
|         |          |       |             |    Sextile - Pluto - 298.58 in Capricorn     |
| Saturn  | 323.61   | VIII  | Aquarius    |
|         |          |       |             |    Semi-sextile - Neptune - 354.24 in Pisces |
| Uranus  | 43.94    | X     | Taurus      |
| Neptune | 354.24   | I     | Pisces      |
| Pluto   | 298.58   | VI    | Capricorn   |
+---------+----------+-------+-------------+-------------------------------------------+

```

```
horoscope --eclipse

Lunar eclipse: 2022-05-16 04:11:00 +0000 UTC
Lunar eclipse: 2022-11-08 10:59:00 +0000 UTC
Lunar eclipse: 2023-05-05 17:23:00 +0000 UTC
Solar eclipse: 2022-04-30 20:41:00 +0000 UTC
Solar eclipse: 2022-10-25 11:00:00 +0000 UTC
Solar eclipse: 2023-04-20 04:16:00 +0000 UTC

```

```

horoscope --retrograde

Mercury retrograde starts: 10 May 22 11:49 UTC ends: 03 Jun 22 07:57 UTC
Mercury retrograde starts: 10 Sep 22 03:27 UTC ends: 02 Oct 22 09:07 UTC
Mercury retrograde starts: 29 Dec 22 09:24 UTC ends: 18 Jan 23 13:13 UTC
Mars retrograde starts: 30 Oct 22 13:25 UTC ends: 12 Jan 23 20:53 UTC
Jupiter retrograde starts: 28 Jul 22 20:38 UTC ends: 23 Nov 22 23:02 UTC
Saturn retrograde starts: 04 Jun 22 21:46 UTC ends: 23 Oct 22 04:06 UTC
Uranus retrograde starts: 24 Aug 22 13:53 UTC ends: 22 Jan 23 22:56 UTC
Neptune retrograde starts: 28 Jun 22 07:53 UTC ends: 04 Dec 22 00:12 UTC
Pluto retrograde starts: 29 Apr 22 18:33 UTC ends: 08 Oct 22 22:02 UTC

```

