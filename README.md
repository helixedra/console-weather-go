# Weather CLI

A simple command-line weather application written in Go that shows current and forecasted temperatures for your location.

## Features

- Automatically detects your location using IP geolocation
- Displays current temperature
- Shows daily temperature forecast with min/max values
- Color-coded temperature output for better readability
- Works anywhere with an internet connection

## Installation

1. Make sure you have Go installed (1.16+)
2. Clone this repository
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Build and run:
   ```bash
   go run main.go
   ```

## Dependencies

- [fatih/color](https://github.com/fatih/color) - For colored terminal output

## APIs Used

- [ipwho.is](https://ipwho.is/) - For geolocation
- [Open-Meteo](https://open-meteo.com/) - For weather data

## Example Output

```
Germany / Berlin:  22.5°C

 5.2°C :  8.7°C - (Mon) Jan 02, 2023
 4.1°C :  9.3°C - (Tue) Jan 03, 2023
 3.8°C : 11.2°C - (Wed) Jan 04, 2023
```

## License

MIT
