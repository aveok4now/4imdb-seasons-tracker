## Overview
Lightweight microservice designed to monitor and notify about announcements for upcoming TV series seasons on IMDb. It periodically scrapes episode pages to detect when placeholder entries (e.g., "Episode #2.1") are updated with real details like episode titles, release dates, and plot summaries.

## Acknowledgments

*   [goquery](https://github.com/PuerkitoBio/goquery) - For HTML parsing.
*   [robfig/cron](https://github.com/robfig/cron/v3) - For scheduling tasks.
