# Barely In Charge

An AI-powered calendar block planner that helps you organize your workday with focus blocks and breaks.

## What It Does

Barely In Charge reads your tasks and existing meetings, then uses AI to create an optimized schedule with focus blocks and breaks in your Google Calendar. You can choose between different planning modes:

- **crunch** - Maximum productivity, minimal breaks
- **normal** - Balanced work and rest periods
- **saver** - Energy-conscious with longer breaks

### How it works

1. **Fetches meetings** - Reads existing meetings from your meetings calendar
2. **AI planning** - Generates optimal schedule with focus blocks and breaks
3. **Creates blocks** - Adds all planned blocks to your blocks calendar
4. **Done!** - Check your calendar to see the schedule

## Prerequisites

- Go 1.21 or higher
- Google Calendar API credentials
- OpenAI API key

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
  "default_mode": "normal",
  "openai_api_key": "sk-your-api-key-here"
}
```

**Configuration Options:**

- `work_hours` - Your working hours (24-hour format)
- `lunch_time` - Your lunch break (24-hour format)
- `meetings_calendar` - Calendar ID to read meetings from (use "primary" for your main calendar)
- `blocks_calendar` - Calendar ID to create focus blocks in (can be different from meetings calendar)
- `default_mode` - Default planning mode: `crunch`, `normal`, or `saver`
- `openai_api_key` - Your OpenAI API key (get one from https://platform.openai.com/api-keys)

## Usage

### View Current Configuration

```bash
./barely-incharge config
```

### Plan Your Day

```bash
./barely-incharge plan --tasks "Write documentation:L, Review PRs:S, Team meeting prep:M"
```

**Flags:**

- `-t, --tasks` - Comma-separated list of tasks with optional size (required)
- `-m, --mode` - Override the default planning mode (optional)

**Task Sizes (T-Shirt Sizing):**

Tasks can include an optional size suffix using the format `Task Title:SIZE`:

- `XS` - 10 minutes (quick fixes, small updates)
- `S` - 15 minutes (code reviews, short meetings)
- `M` - 30 minutes (default if no size specified)
- `L` - 60 minutes (feature development, deep work)
- `XL` - 90 minutes (complex features, major refactoring)

**Examples:**

```bash
# Tasks with sizes
./barely-incharge plan -t "Write docs:L, Review PRs:S, Quick fix:XS"

# Tasks without sizes default to M (30 min)
./barely-incharge plan -t "Code feature X, Write tests, Deploy"

# Mixed sizes and modes
./barely-incharge plan -t "Deep work:XL, Coffee break:XS" -m saver

# Override with crunch mode
./barely-incharge plan -t "Urgent bug fix:M, Code review:S" -m crunch
```

### First Run

On first run, the app will open your browser to authenticate with Google Calendar. After authorization:
- A `token.json` file is created with your credentials
- Future runs will use this token (no browser popup)
- The token is valid until you revoke access
