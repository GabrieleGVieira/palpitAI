from __future__ import annotations

from collections import defaultdict
from datetime import date
from math import log

import pandas as pd


def _expected(rating_a: float, rating_b: float) -> float:
    return 1.0 / (1.0 + 10 ** ((rating_b - rating_a) / 400.0))


def _actual(home_score: int, away_score: int) -> tuple[float, float]:
    if home_score > away_score:
        return 1.0, 0.0
    if home_score < away_score:
        return 0.0, 1.0
    return 0.5, 0.5


def _k_factor(tournament: str | None) -> float:
    value = (tournament or "").lower()
    if "fifa world cup" in value or value == "world cup":
        return 40.0
    continental_markers = [
        "uefa euro",
        "copa america",
        "african cup",
        "asian cup",
        "gold cup",
        "nations cup",
        "confederations cup",
        "oceania",
    ]
    if any(marker in value for marker in continental_markers):
        return 30.0
    return 20.0


def _goal_multiplier(goal_diff: int) -> float:
    diff = abs(goal_diff)
    if diff <= 1:
        return 1.0
    return min(1.75, 1.0 + log(diff) / 3.0)


def calculate_elo_ratings(matches_df: pd.DataFrame, until_date: date, initial_rating: float = 1500.0) -> dict[str, float]:
    ratings: defaultdict[str, float] = defaultdict(lambda: initial_rating)
    df = matches_df[matches_df["date"] <= until_date].sort_values("date")

    for row in df.itertuples(index=False):
        home = getattr(row, "home_team_id", None) or getattr(row, "home_team")
        away = getattr(row, "away_team_id", None) or getattr(row, "away_team")
        home_score = int(row.home_score)
        away_score = int(row.away_score)
        home_rating = ratings[home]
        away_rating = ratings[away]
        actual_home, actual_away = _actual(home_score, away_score)
        expected_home = _expected(home_rating, away_rating)
        expected_away = _expected(away_rating, home_rating)
        multiplier = _goal_multiplier(home_score - away_score)
        k = _k_factor(getattr(row, "tournament", None)) * multiplier

        ratings[home] = home_rating + k * (actual_home - expected_home)
        ratings[away] = away_rating + k * (actual_away - expected_away)

    return dict(ratings)
