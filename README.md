# csv2midi

[![GitHub release](https://img.shields.io/github/release/moutend/csv2midi.svg?style=flat-square)][release]
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]
[![CircleCI](https://circleci.com/gh/moutend/csv2midi.svg?style=svg&circle-token=e7748578056ded93a5532904c047fc0f23db3bba)](https://circleci.com/gh/moutend/csv2midi)

[release]: https://github.com/moutend/csv2midi/releases
[license]: https://github.com/moutend/csv2midi/blob/master/LICENSE
[status]: https://circleci.com/gh/moutend/csv2midi

`csv2midi` converts CSV to standard MIDI file.

# Download

You can download `csv2midi` from GitHub releases page.

# Usage

```console
csv2midi music.csv
```

And then the standard MIDI file named `music.mid` will be generated.

# CSV specification

The CSV file must have 3 columns at least. The first column is delta time which defines when the event play. The second column is event type and  the third column is value for that event.

# Supported events

| Event Type | Values |
|:---|:---|
| `on` | note name | velocity |

# Example

The following CSV is 8 beat.

## Contributing

1. Fork ([https://github.com/moutend/csv2midi/fork](https://github.com/moutend/csv2midi/fork))
1. Create a feature branch
1. Add changes
1. Run `go fmt`
1. Commit your changes
1. Open a new Pull Request

## Author

[Yoshiyuki Koyanagi](https://github.com/moutend)
