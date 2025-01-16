# Lab
**Lab** is a quick way to spin up throwaway files for rapid experimenting. Reduces the overhead of trying out code snippets. No more fuss about filenames or folders—just type `lab js` (or any extension) and start coding instantly in your editor. All files are auto-named (e.g., `250112a.js`) and live in a `lab` folder. By default, they vanish after 7 days, keeping things tidy.

## Install
### macOS
```bash
curl -L https://github.com/lugenx/lab/releases/latest/download/lab-darwin-amd64 -o /tmp/lab && chmod +x /tmp/lab && sudo mv /tmp/lab /usr/local/bin/lab
```
### Linux
```bash
curl -L https://github.com/lugenx/lab/releases/latest/download/lab-linux-amd64 -o /tmp/lab && chmod +x /tmp/lab && sudo mv /tmp/lab /usr/local/bin/lab
```
### Windows
```powershell
curl -L -o lab-windows-amd64.exe https://github.com/lugenx/lab/releases/latest/download/lab-windows-amd64.exe
```
Move-Item .\lab-windows-amd64.exe "C:\Windows\System32\lab.exe"

## Usage
Create a new file:
```bash
lab js      # creates JavaScript file
lab py      # creates Python file
lab any     # creates file with any extension
```

List your files:
```bash
lab

To open, use: lab <number>
To create: lab <extension>

[1] 250112c.js     [6d]
[2] 250112b.py     [6d]
[3] 250112a.go     [6d]
```

Open a file:
```bash
lab 0       # opens config
lab 1       # opens most recent file
lab 2       # opens second file
```

## Features
- **Instant Start**: `lab <extension>` creates and opens a fresh file.
- **Quick Access**: `lab <number>` reopens recent files.
- **Auto-Cleanup**: Files expire automatically after 7 days (configurable).
- **Simple Listing**: Just run `lab` to see your files, newest first.

## Configuration
Lab creates a config file at `~/lab/.lab`:
```
editor=nvim          # your preferred editor
lifedays=7          # how long to keep files
prefix=lab          # file prefix
<!--
show_tips=true      # show random tips
show_instructions=true
show_filepath=true
-->
```

You can also set `LABPATH` environment variable to change the lab directory location from the default `~/lab`.

Focus on experimenting and trying out ideas without distractions—Lab simplifies the process so you can start coding immediately.
