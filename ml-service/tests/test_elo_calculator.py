from __future__ import annotations

from datetime import date

import pandas as pd

from app.metrics.elo_calculator import calculate_elo_ratings


def test_elo_updates_chronologically_and_rewards_winner() -> None:
    matches = pd.DataFrame(
        [
            {
                "date": date(2024, 1, 1),
                "home_team_id": "a",
                "away_team_id": "b",
                "home_score": 1,
                "away_score": 0,
                "tournament": "Friendly",
            },
            {
                "date": date(2024, 1, 2),
                "home_team_id": "b",
                "away_team_id": "a",
                "home_score": 0,
                "away_score": 2,
                "tournament": "FIFA World Cup",
            },
        ]
    )

    ratings = calculate_elo_ratings(matches, date(2024, 1, 2))

    assert ratings["a"] > 1500
    assert ratings["b"] < 1500
    assert ratings["a"] - ratings["b"] > 40
