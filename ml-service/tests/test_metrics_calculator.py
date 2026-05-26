from __future__ import annotations

from datetime import date

import pandas as pd

from app.metrics.metrics_calculator import (
    calculate_attack_score,
    calculate_basic_stats,
    calculate_defense_score,
    calculate_goalscorer_stats,
    calculate_recent_form,
    calculate_shootout_stats,
)


def sample_matches() -> pd.DataFrame:
    return pd.DataFrame(
        [
            {
                "date": date(2024, 1, 1),
                "home_team_id": "brazil",
                "away_team_id": "argentina",
                "home_score": 2,
                "away_score": 0,
                "tournament": "Friendly",
            },
            {
                "date": date(2024, 1, 5),
                "home_team_id": "brazil",
                "away_team_id": "france",
                "home_score": 1,
                "away_score": 1,
                "tournament": "Friendly",
            },
            {
                "date": date(2024, 1, 9),
                "home_team_id": "england",
                "away_team_id": "brazil",
                "home_score": 3,
                "away_score": 1,
                "tournament": "Friendly",
            },
        ]
    )


def test_basic_stats_win_draw_loss_rates() -> None:
    stats = calculate_basic_stats(sample_matches(), "brazil")

    assert stats["matches_played"] == 3
    assert stats["win_rate"] == 1 / 3
    assert stats["draw_rate"] == 1 / 3
    assert stats["loss_rate"] == 1 / 3
    assert stats["avg_goals_scored"] == 4 / 3
    assert stats["avg_goals_conceded"] == 4 / 3


def test_recent_form_weights_more_recent_matches() -> None:
    score = calculate_recent_form(sample_matches(), "brazil", date(2024, 1, 10), last_n=3)

    assert round(score, 2) == 27.78


def test_attack_score_is_normalized_and_uses_world_cup_component() -> None:
    score = calculate_attack_score(2.0, 3.0, 80.0)

    assert round(score, 2) == 78.67


def test_defense_score_rewards_low_conceded_and_clean_sheets() -> None:
    score = calculate_defense_score(0.5, 0.6, 0.25)

    assert round(score, 2) == 81.58


def test_goalscorer_stats_measure_depth_and_penalty_dependency() -> None:
    goalscorers = pd.DataFrame(
        [
            {"date": date(2024, 1, 1), "team_id": "brazil", "scorer": "A", "own_goal": False, "penalty": False},
            {"date": date(2024, 1, 1), "team_id": "brazil", "scorer": "B", "own_goal": False, "penalty": True},
            {"date": date(2024, 1, 1), "team_id": "brazil", "scorer": "B", "own_goal": False, "penalty": False},
            {"date": date(2024, 1, 1), "team_id": "brazil", "scorer": "C", "own_goal": True, "penalty": False},
        ]
    )

    stats = calculate_goalscorer_stats(goalscorers, "brazil", date(2024, 1, 2))

    assert stats["goals"] == 3
    assert stats["unique_scorers"] == 2
    assert stats["penalty_goals"] == 1
    assert stats["own_goals"] == 1
    assert stats["scorer_depth_score"] == 100.0
    assert round(stats["penalty_independence_score"], 2) == 4.76


def test_shootout_stats_measure_win_rate() -> None:
    shootouts = pd.DataFrame(
        [
            {"date": date(2024, 1, 1), "home_team_id": "brazil", "away_team_id": "france", "winner_team_id": "brazil"},
            {"date": date(2024, 1, 5), "home_team_id": "argentina", "away_team_id": "brazil", "winner_team_id": "argentina"},
        ]
    )

    stats = calculate_shootout_stats(shootouts, "brazil", date(2024, 1, 10))

    assert stats["shootouts_played"] == 2
    assert stats["shootouts_won"] == 1
    assert stats["shootout_win_rate"] == 0.5
    assert stats["shootout_score"] == 50.0
