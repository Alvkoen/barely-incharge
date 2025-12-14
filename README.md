# Barely In Charge

An AI-powered calendar block planner that helps you organize your workday with focus blocks and breaks.

## What It Does

Barely In Charge reads your tasks and existing meetings, then uses AI to create an optimized schedule with focus blocks and breaks in your Google Calendar. You can choose between different planning modes:

- **crunch** - Maximum productivity, minimal breaks
- **normal** - Balanced work and rest periods
- **saver** - Energy-conscious with longer breaks

## Prerequisites

- Go 1.21 or higher
- Google Calendar API credentials

## Setup

### 1. Clone and Build

```bash
git clone https://github.com/Alvkoen/barely-incharge.git
cd barely-incharge
go build -o barely-incharge
```

### 2. Configure Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project (or select existing)
3. Enable the Google Calendar API
4. Create OAuth 2.0 credentials (Desktop app)
5. Download the credentials and save as `credentials.json` in the project root

**Important:** Add these to `.gitignore`:
```
credentials.json
token.json
```

### 3. Create Configuration File

Create a `config.json` file in the project root:

```json
{
  "work_hours": {
    "start": "09:00",
    "end": "17:00"
  },
  "lunch_time": {
    "start": "12:00",
    "end": "13:00"
  },
  "meetings_calendar": "primary",
  "blocks_calendar": "primary",
  "default_mode": "normal"
}
```

**Configuration Options:**

- `work_hours` - Your working hours (24-hour format)
- `lunch_time` - Your lunch break (24-hour format)
- `meetings_calendar` - Calendar ID to read meetings from (use "primary" for your main calendar)
- `blocks_calendar` - Calendar ID to create focus blocks in (can be different from meetings calendar)
- `default_mode` - Default planning mode: `crunch`, `normal`, or `saver`

## Usage

### View Current Configuration

```bash
./barely-incharge config
```

### Plan Your Day

```bash
./barely-incharge plan --tasks "Write documentation, Review PRs, Team meeting prep"
```

**Flags:**

- `-t, --tasks` - Comma-separated list of tasks (required)
- `-m, --mode` - Override the default planning mode (optional)

**Examples:**

```bash
# Use default mode from config
./barely-incharge plan -t "Code feature X, Write tests, Deploy"

# Override with crunch mode
./barely-incharge plan -t "Urgent bug fix, Code review" -m crunch

# Use energy saver mode with longer breaks
./barely-incharge plan -t "Study Go, Workout, Deep work session" -m saver
```

### First Run

On first run, the app will open your browser to authenticate with Google Calendar. After authorization:
- A `token.json` file is created with your credentials
- Future runs will use this token (no browser popup)
- The token is valid until you revoke access

## Project Status

🚧 **In Development** 