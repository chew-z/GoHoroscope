# GoHoroscope

A bit of astrology in terminal with Go.

Prints Ascendant, Houses, Zodiac signs, Aspects, Positions, Solar and Lunar eclipses, Retrograde movements, nearest New Moon and Full Moon time and date.

For htttp server with charts check [GoHoroscopeChart](https://github.com/chew-z/GoHoroscopeChart)

## Installation

```
- git clone https://github.com/chew-z/GoHoroscope.git
- go build -o ./bin/horoscope .
```

This project is using Swiss Ephemeris with [swephgo](https://github.com/mshafiee/swephgo) as a wrapper. You will also need Swiss Ephemem library and ephemeris files.

Download the Swiss Ephemeris Library [here](https://www.astro.com/ftp/swisseph/). After compiling the library, copy the libswe.so file to /usr/local/lib/

```
$ cp libswe.so /usr/local/lib/
```

[Download some ephemeris files](https://www.astro.com/ftp/swisseph/ephe/) and puth them in designated folder. _sepl_18.se1_ is a good place to start.

Finally you will need some shell variables (timezone, geographical position, location of Swiss Ephemeris files on your system)

```
export CITY="Europe/London"
export LATITUDE="51.5072"¬
export LONGITUDE="0.1276"¬
export HOUSE_SYSTEM="Placidus"
export SWISSPATH="/usr/local/share/sweph/ephe"
```

or put variables in `.env` file looking like `.env.example`

## Usage

```
horoscope --horoscope [date]

20 Apr 22 10:48 CEST - lat: 52.43, lon: 20.89
Ascendant: 139.43 MC: 33.71, House system: Placidus

+-------+----------+-------+-------------+
| House | Position | Cusp  | Sign        |
+-------+----------+-------+-------------+
| I     | 139.43   | 19.43 | Leo         |
| II    | 157.49   | 7.49  | Virgo       |
| III   | 181.35   | 1.35  | Libra       |
| IV    | 213.71   | 3.71  | Scorpio     |
| V     | 253.68   | 13.68 | Sagittarius |
| VI    | 290.72   | 20.72 | Capricorn   |
| VII   | 319.43   | 19.43 | Aquarius    |
| VIII  | 337.49   | 7.49  | Pisces      |
| IX    | 1.35     | 1.35  | Aries       |
| X     | 33.71    | 3.71  | Taurus      |
| XI    | 73.68    | 13.68 | Gemini      |
| XII   | 110.72   | 20.72 | Cancer      |
+-------+----------+-------+-------------+

+---------+----------+-------+-------------+--------------------------------+
| Planet  | Position | House | Sign        | Aspects                        |
+---------+----------+-------+-------------+--------------------------------+
| Sun     | 30.34    | IX    | Taurus      |
|         |          |       |             | Sextile Mars in Pisces         |
|         |          |       |             | Square Pluto in Capricorn      |
| Moon    | 259.69   | V     | Sagittarius |
|         |          |       |             | Square Venus in Pisces         |
|         |          |       |             | Sextile Saturn in Aquarius     |
|         |          |       |             | Square Neptune in Pisces       |
| Mercury | 47.68    | X     | Taurus      |
|         |          |       |             | Sextile Venus in Pisces        |
|         |          |       |             | Square Saturn in Aquarius      |
|         |          |       |             | Conjunction Uranus in Taurus   |
| Venus   | 346.21   | I     | Pisces      |
|         |          |       |             | Conjunction Jupiter in Pisces  |
|         |          |       |             | Sextile Uranus in Taurus       |
|         |          |       |             | Conjunction Neptune in Pisces  |
| Mars    | 334.02   | VII   | Pisces      |
| Jupiter | 355.73   | I     | Pisces      |
|         |          |       |             | Conjunction Neptune in Pisces  |
|         |          |       |             | Sextile Pluto in Capricorn     |
| Saturn  | 323.61   | VII   | Aquarius    |
|         |          |       |             | Semi-sextile Neptune in Pisces |
| Uranus  | 43.94    | X     | Taurus      |
| Neptune | 354.24   | I     | Pisces      |
| Pluto   | 298.58   | VI    | Capricorn   |
+---------+----------+-------+-------------+--------------------------------+

```

```
horoscope --eclipse [date]

+--------------------------------+
| Lunar Eclipse                  |
+--------------------------------+
| 2022-05-16 06:11:00 +0200 CEST |
| 2022-11-08 11:59:00 +0100 CET  |
| 2023-05-05 19:23:00 +0200 CEST |
+--------------------------------+

+--------------------------------+
| Solar Eclipse                  |
+--------------------------------+
| 2022-04-30 22:41:00 +0200 CEST |
| 2022-10-25 13:00:00 +0200 CEST |
| 2023-04-20 06:16:00 +0200 CEST |
+--------------------------------+

```

```
horoscope --retrograde [date]

+---------+----------------------+----------------------+
| Planet  | Starts               | Ends                 |
+---------+----------------------+----------------------+
| Mercury | 10 May 22 13:49 CEST | 03 Jun 22 09:57 CEST |
| Mercury | 10 Sep 22 05:27 CEST | 02 Oct 22 11:07 CEST |
| Mercury | 29 Dec 22 10:24 CET  | 18 Jan 23 14:13 CET  |
| Mars    | 30 Oct 22 14:25 CET  | 12 Jan 23 21:53 CET  |
| Jupiter | 28 Jul 22 22:38 CEST | 24 Nov 22 00:02 CET  |
| Saturn  | 04 Jun 22 23:46 CEST | 23 Oct 22 06:06 CEST |
| Uranus  | 24 Aug 22 15:53 CEST | 22 Jan 23 23:56 CET  |
| Neptune | 28 Jun 22 09:53 CEST | 04 Dec 22 01:12 CET  |
| Pluto   | 29 Apr 22 20:33 CEST | 09 Oct 22 00:02 CEST |
+---------+----------------------+----------------------+

```

```
horoscope --moon [date]

+----------------------+-----+----------------------+-----+
| New Moon             |     | Full Moon            |     |
+----------------------+-----+----------------------+-----+
| 02 Jan 22 18:09 CET  | 🌑  | 17 Jan 22 23:21 CET  | 🌕  |
| 01 Feb 22 05:33 CET  | 🌑  | 16 Feb 22 16:58 CET  | 🌕  |
| 02 Mar 22 17:50 CET  | 🌑  | 18 Mar 22 07:44 CET  | 🌕  |
| 01 Apr 22 06:55 CEST | 🌑  | 16 Apr 22 19:18 CEST | 🌕  |
| 30 Apr 22 20:43 CEST | 🌑  | 16 May 22 04:13 CEST | 🌕  |
| 30 May 22 11:12 CEST | 🌑  | 14 Jun 22 11:30 CEST | 🌕  |
| 29 Jun 22 02:20 CEST | 🌑  | 13 Jul 22 18:18 CEST | 🌕  |
| 28 Jul 22 17:42 CEST | 🌑  | 12 Aug 22 01:38 CEST | 🌕  |
| 27 Aug 22 08:36 CEST | 🌑  | 10 Sep 22 10:24 CEST | 🌕  |
| 25 Sep 22 22:24 CEST | 🌑  | 09 Oct 22 22:24 CEST | 🌕  |
| 25 Oct 22 11:01 CEST | 🌑  | 08 Nov 22 11:02 CET  | 🌕  |
| 23 Nov 22 22:42 CET  | 🌑  | 08 Dec 22 03:39 CET  | 🌕  |
| 23 Dec 22 09:53 CET  | 🌑  | 06 Jan 23 22:45 CET  | 🌕  |
| 21 Jan 23 20:48 CET  | 🌑  | 05 Feb 23 18:39 CET  | 🌕  |
| 20 Feb 23 07:26 CET  | 🌑  | 07 Mar 23 13:13 CET  | 🌕  |
| 21 Mar 23 17:49 CET  | 🌑  | 06 Apr 23 04:54 CEST | 🌕  |
| 20 Apr 23 04:18 CEST | 🌑  | 05 May 23 17:24 CEST | 🌕  |
| 19 May 23 15:32 CEST | 🌑  | 04 Jun 23 03:16 CEST | 🌕  |
| 18 Jun 23 04:09 CEST | 🌑  | 03 Jul 23 11:24 CEST | 🌕  |
| 17 Jul 23 18:27 CEST | 🌑  | 01 Aug 23 18:41 CEST | 🌕  |
| 16 Aug 23 10:04 CEST | 🌑  | 31 Aug 23 02:01 CEST | 🌕  |
+----------------------+-----+----------------------+-----+

```

## Interpretation

As for interpretation of various aspects etc. I find [Astrology King](https://astrologyking.com/) quite informative.

[Aspects](https://astrologyking.com/aspects/)

[Transits](https://astrologyking.com/transits/)

[Retrograde planets](https://astrologyking.com/retrograde/)

[Moon phases](https://astrologyking.com/2022-moon-phases-calendar/)
