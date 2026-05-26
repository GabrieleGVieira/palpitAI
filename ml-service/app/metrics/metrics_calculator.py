from __future__ import annotations

from datetime import date
from typing import Any

import pandas as pd

from .elo_calculator import calculate_elo_ratings
from .schemas import TeamMetric


def normalize_score(value: float | None, min_value: float, max_value: float) -> float | None:
    if value is None or max_value == min_value:
        return None
    return max(0.0, min(100.0, (value - min_value) / (max_value - min_value) * 100.0))


def team_matches(df: pd.DataFrame, team_id: str, metric_date: date) -> pd.DataFrame:
    filtered = df[df["date"] <= metric_date]
    return filtered[(filtered["home_team_id"] == team_id) | (filtered["away_team_id"] == team_id)].copy()


def _goals_for_against(row: Any, team_id: str) -> tuple[int, int]:
    if row.home_team_id == team_id:
        return int(row.home_score), int(row.away_score)
    return int(row.away_score), int(row.home_score)


def calculate_basic_stats(df: pd.DataFrame, team_id: str) -> dict[str, float | int | None]:
    matches = df[(df["home_team_id"] == team_id) | (df["away_team_id"] == team_id)]
    if matches.empty:
        return {
            "matches_played": 0,
            "avg_goals_scored": None,
            "avg_goals_conceded": None,
            "win_rate": None,
            "draw_rate": None,
            "loss_rate": None,
            "clean_sheet_rate": None,
        }

    wins = draws = losses = clean_sheets = goals_for = goals_against = 0
    for row in matches.itertuples(index=False):
        scored, conceded = _goals_for_against(row, team_id)
        goals_for += scored
        goals_against += conceded
        clean_sheets += int(conceded == 0)
        if scored > conceded:
            wins += 1
        elif scored == conceded:
            draws += 1
        else:
            losses += 1

    total = len(matches)
    return {
        "matches_played": total,
        "avg_goals_scored": goals_for / total,
        "avg_goals_conceded": goals_against / total,
        "win_rate": wins / total,
        "draw_rate": draws / total,
        "loss_rate": losses / total,
        "clean_sheet_rate": clean_sheets / total,
    }


def calculate_recent_form(df: pd.DataFrame, team_id: str, metric_date: date, last_n: int = 10) -> float | None:
    matches = team_matches(df, team_id, metric_date).sort_values("date").tail(last_n)
    if matches.empty:
        return None

    weighted_points = 0.0
    max_points = 0.0
    for index, row in enumerate(matches.itertuples(index=False), start=1):
        scored, conceded = _goals_for_against(row, team_id)
        points = 3 if scored > conceded else 1 if scored == conceded else 0
        weight = index / len(matches)
        weighted_points += points * weight
        max_points += 3 * weight

    return weighted_points / max_points * 100.0


def calculate_attack_score(
    avg_goals_scored: float | None,
    recent_avg_goals_scored: float | None,
    world_cup_history_score: float | None,
    scorer_depth_score: float | None = None,
    penalty_independence_score: float | None = None,
) -> float | None:
    if avg_goals_scored is None:
        return None
    historical = normalize_score(avg_goals_scored, 0.0, 3.0) or 0.0
    recent = normalize_score(recent_avg_goals_scored if recent_avg_goals_scored is not None else avg_goals_scored, 0.0, 3.0) or 0.0
    world_cup = world_cup_history_score if world_cup_history_score is not None else historical
    scorer_depth = scorer_depth_score if scorer_depth_score is not None else historical
    penalty_independence = penalty_independence_score if penalty_independence_score is not None else historical
    return 0.30 * historical + 0.30 * recent + 0.15 * world_cup + 0.15 * scorer_depth + 0.10 * penalty_independence


def calculate_defense_score(
    avg_goals_conceded: float | None,
    clean_sheet_rate: float | None,
    recent_avg_goals_conceded: float | None,
) -> float | None:
    if avg_goals_conceded is None:
        return None
    conceded_score = 100.0 - (normalize_score(avg_goals_conceded, 0.0, 3.0) or 0.0)
    recent_score = 100.0 - (normalize_score(recent_avg_goals_conceded if recent_avg_goals_conceded is not None else avg_goals_conceded, 0.0, 3.0) or 0.0)
    clean_sheet_score = (clean_sheet_rate or 0.0) * 100.0
    return 0.45 * conceded_score + 0.35 * recent_score + 0.20 * clean_sheet_score


def calculate_world_cup_history_score(df: pd.DataFrame, team_id: str, metric_date: date) -> float | None:
    matches = team_matches(df, team_id, metric_date)
    matches = matches[matches["tournament"].str.lower().fillna("").str.contains("world cup")]
    if matches.empty:
        return None

    weighted_points = 0.0
    max_points = 0.0
    max_year = max(match_date.year for match_date in matches["date"])
    for row in matches.itertuples(index=False):
        scored, conceded = _goals_for_against(row, team_id)
        points = 3 if scored > conceded else 1 if scored == conceded else 0
        recency_weight = 1.0 + max(0, row.date.year - (max_year - 16)) / 16.0
        weighted_points += points * recency_weight
        max_points += 3 * recency_weight

    return weighted_points / max_points * 100.0


def calculate_stage_score(df: pd.DataFrame, team_id: str, metric_date: date, stage_markers: list[str]) -> float | None:
    matches = team_matches(df, team_id, metric_date)
    if "stage" not in matches.columns:
        return None
    stage_text = matches["stage"].fillna("").str.lower()
    matches = matches[stage_text.apply(lambda value: any(marker in value for marker in stage_markers))]
    if matches.empty:
        return None
    stats = calculate_basic_stats(matches, team_id)
    return (stats["win_rate"] or 0.0) * 100.0


def calculate_goalscorer_stats(goalscorers: pd.DataFrame | None, team_id: str, metric_date: date) -> dict[str, float | int | None]:
    if goalscorers is None or goalscorers.empty:
        return {
            "goals": 0,
            "unique_scorers": 0,
            "penalty_goals": 0,
            "own_goals": 0,
            "scorer_depth_score": None,
            "penalty_independence_score": None,
        }

    team_goals = goalscorers[(goalscorers["date"] <= metric_date) & (goalscorers["team_id"] == team_id)]
    if team_goals.empty:
        return {
            "goals": 0,
            "unique_scorers": 0,
            "penalty_goals": 0,
            "own_goals": 0,
            "scorer_depth_score": None,
            "penalty_independence_score": None,
        }

    own_goals = int(team_goals["own_goal"].sum())
    valid_goals = team_goals[~team_goals["own_goal"]]
    goals = len(valid_goals)
    penalty_goals = int(valid_goals["penalty"].sum())
    unique_scorers = valid_goals["scorer"].dropna().nunique()
    scorer_depth_score = normalize_score(unique_scorers / max(1, goals), 0.05, 0.45)
    penalty_share = penalty_goals / goals if goals else 0.0
    penalty_independence_score = 100.0 - (normalize_score(penalty_share, 0.0, 0.35) or 0.0)

    return {
        "goals": goals,
        "unique_scorers": int(unique_scorers),
        "penalty_goals": penalty_goals,
        "own_goals": own_goals,
        "scorer_depth_score": scorer_depth_score,
        "penalty_independence_score": penalty_independence_score,
    }


def calculate_shootout_stats(shootouts: pd.DataFrame | None, team_id: str, metric_date: date) -> dict[str, float | int | None]:
    if shootouts is None or shootouts.empty:
        return {"shootouts_played": 0, "shootouts_won": 0, "shootout_win_rate": None, "shootout_score": None}

    team_shootouts = shootouts[
        (shootouts["date"] <= metric_date)
        & ((shootouts["home_team_id"] == team_id) | (shootouts["away_team_id"] == team_id))
    ]
    if team_shootouts.empty:
        return {"shootouts_played": 0, "shootouts_won": 0, "shootout_win_rate": None, "shootout_score": None}

    played = len(team_shootouts)
    won = int((team_shootouts["winner_team_id"] == team_id).sum())
    win_rate = won / played
    return {
        "shootouts_played": played,
        "shootouts_won": won,
        "shootout_win_rate": win_rate,
        "shootout_score": win_rate * 100.0,
    }


def combine_knockout_score(stage_score: float | None, shootout_score: float | None) -> float | None:
    if stage_score is None:
        return shootout_score
    if shootout_score is None:
        return stage_score
    return 0.75 * stage_score + 0.25 * shootout_score


def _recent_goal_averages(df: pd.DataFrame, team_id: str, metric_date: date, last_n: int = 10) -> tuple[float | None, float | None]:
    matches = team_matches(df, team_id, metric_date).sort_values("date").tail(last_n)
    if matches.empty:
        return None, None
    scored_total = conceded_total = 0
    for row in matches.itertuples(index=False):
        scored, conceded = _goals_for_against(row, team_id)
        scored_total += scored
        conceded_total += conceded
    return scored_total / len(matches), conceded_total / len(matches)


def calculate_team_metrics(
    df: pd.DataFrame,
    team_ids: list[str],
    metric_date: date,
    goalscorers: pd.DataFrame | None = None,
    shootouts: pd.DataFrame | None = None,
) -> tuple[list[TeamMetric], dict[str, dict[str, Any]]]:
    filtered = df[df["date"] <= metric_date].copy()
    elo = calculate_elo_ratings(filtered, metric_date)
    metrics: list[TeamMetric] = []
    snapshots: dict[str, dict[str, Any]] = {}

    for team_id in team_ids:
        stats = calculate_basic_stats(filtered, team_id)
        recent_form = calculate_recent_form(filtered, team_id, metric_date)
        world_cup = calculate_world_cup_history_score(filtered, team_id, metric_date)
        recent_scored, recent_conceded = _recent_goal_averages(filtered, team_id, metric_date)
        goalscorer_stats = calculate_goalscorer_stats(goalscorers, team_id, metric_date)
        shootout_stats = calculate_shootout_stats(shootouts, team_id, metric_date)
        attack = calculate_attack_score(
            stats["avg_goals_scored"],
            recent_scored,
            world_cup,
            goalscorer_stats["scorer_depth_score"],
            goalscorer_stats["penalty_independence_score"],
        )
        defense = calculate_defense_score(stats["avg_goals_conceded"], stats["clean_sheet_rate"], recent_conceded)
        knockout_stage = calculate_stage_score(filtered, team_id, metric_date, ["final", "semi", "quarter", "round of", "knockout"])
        knockout = combine_knockout_score(knockout_stage, shootout_stats["shootout_score"])
        group_stage = calculate_stage_score(filtered, team_id, metric_date, ["group"])

        metrics.append(
            TeamMetric(
                team_id=team_id,
                metric_date=metric_date,
                elo_score=elo.get(team_id),
                attack_score=attack,
                defense_score=defense,
                recent_form_score=recent_form,
                world_cup_history_score=world_cup,
                knockout_score=knockout,
                group_stage_score=group_stage,
                avg_goals_scored=stats["avg_goals_scored"],
                avg_goals_conceded=stats["avg_goals_conceded"],
                win_rate=stats["win_rate"],
                draw_rate=stats["draw_rate"],
                loss_rate=stats["loss_rate"],
                matches_played=int(stats["matches_played"]),
            )
        )
        snapshots[team_id] = {
            "metric_date": metric_date.isoformat(),
            "basic_stats": stats,
            "recent_form_score": recent_form,
            "recent_avg_goals_scored": recent_scored,
            "recent_avg_goals_conceded": recent_conceded,
            "world_cup_history_score": world_cup,
            "goalscorer_stats": goalscorer_stats,
            "shootout_stats": shootout_stats,
            "knockout_stage_score": knockout_stage,
            "elo_score": elo.get(team_id),
        }

    return metrics, snapshots
