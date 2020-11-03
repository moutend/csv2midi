csv2midi
========

`csv2midi` converts CSV to standard MIDI file.

## Install

```console
go get -u github.com/moutend/csv2midi/cmd/csv2midi
```

## Usage

```console
csv2midi music.csv
```

Then you'll get the standard MIDI file named `music.mid`.

## CSV specification

The CSV file must have 3 columns at least. The first column is delta time which defines when the event play. The second column is event type and  the third column is value for that event.

## Supported events

| Event Type | Values |
|:---|:---|
| `on` | note name | velocity |

## Example

The following CSV is 8 beat.

## Contributing

1. Fork ([https://github.com/moutend/csv2midi/fork](https://github.com/moutend/csv2midi/fork))
1. Create a feature branch
1. Add changes
1. Run `go fmt`
1. Commit your changes
1. Open a new Pull Request

## LICENSE

MIT

## Author

[Yoshiyuki Koyanagi](https://github.com/moutend)
