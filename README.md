# Groupie Tracker

Groupie Tracker is a desktop application built with Go (Golang) and the Fyne GUI toolkit. It interfaces with an external REST API to display information about music artists and bands, including their history, members, and concert locations.

The application focuses on data manipulation, user interface design, and handling asynchronous API requests, including geocoding and map visualization.

## Project Overview

This application retrieves data from a specific API containing information about artists, their concert locations, and dates. It presents this data in a user-friendly graphical interface that allows users to search, filter, and view detailed information. Additionally, it integrates a mapping feature to visualize concert locations using OpenStreetMap.

## Features

- **Artist Directory**: Displays a grid of artists with their names and images.
- **Advanced Search**: Allows users to search by artist name, member name, creation date, first album year, or concert location.
- **Filtering System**: Users can filter the artist list based on:
  - Creation year (range)
  - First album year (range)
  - Number of members
  - Concert locations
- **Detailed View**: Shows comprehensive information for a selected artist, including:
  - Band members
  - Creation date and first album year
  - List of concert dates and locations
- **Geolocation & Mapping**:
  - Converts concert location names into geographic coordinates using the Nominatim API.
  - Displays an interactive map with markers using OpenStreetMap tiles.

## Technical Stack

- **Language**: Go (Golang) version 1.25+
- **GUI Framework**: Fyne v2 (v2.7.2)
- **Data Format**: JSON
- **External APIs**:
  - Artist Data: groupietrackers.herokuapp.com
  - Geocoding: nominatim.openstreetmap.org
  - Map Tiles: tile.openstreetmap.org

## Project Structure

- **main.go**: The entry point of the application. It initializes the application loop.
- **appli/**: Contains the core logic and UI components.
  - **api.go**: Handles all HTTP requests, JSON parsing, and data structures (Artist, Relations, Geocoding).
  - **page.go**: Manages the Fyne UI layout, event handling, filtering logic, and map rendering.
- **go.mod / go.sum**: Manages Go module dependencies.

## Prerequisites

To run this project, you need the following installed on your machine:

1. **Go**: Ensure you have Go installed (version 1.25 or compatible).
2. **C Compiler**: Fyne requires a C compiler (like GCC) for CGO bindings, especially for graphics rendering.
   - **Windows**: TDM-GCC or MinGW-w64.
   - **macOS**: Xcode Command Line Tools.
   - **Linux**: GCC (usually installed by default or via `build-essential`).

## Installation and Execution

1. Clone the repository or download the source code.

2. Open a terminal in the project root directory.

3. Install the dependencies:
   ```bash
   go mod tidy
