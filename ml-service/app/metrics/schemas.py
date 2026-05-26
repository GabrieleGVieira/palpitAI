from __future__ import annotations

from dataclasses import dataclass
from datetime import date
from typing import Any


@dataclass(frozen=True)
class Team:
    id: str
    name: str
    country_code: str | None = None


@dataclass(frozen=True)
class TeamMetric:
    team_id: str
    metric_date: date
    elo_score: float | None
    attack_score: float | None
    defense_score: float | None
    recent_form_score: float | None
    world_cup_history_score: float | None
    knockout_score: float | None
    group_stage_score: float | None
    avg_goals_scored: float | None
    avg_goals_conceded: float | None
    win_rate: float | None
    draw_rate: float | None
    loss_rate: float | None
    matches_played: int
    source: str = "metrics-engine-v1"


@dataclass(frozen=True)
class MatchTarget:
    match_id: str | None
    match_date: date
    home_team_id: str
    away_team_id: str
    tournament: str | None
    stage: str | None
    neutral: bool


@dataclass(frozen=True)
class Snapshot:
    team_id: str
    snapshot_type: str
    payload_json: dict[str, Any]
