from __future__ import annotations

import json
from datetime import date

import pandas as pd

from app.metrics.local_data import (
    build_aliases_for_existing_teams,
    combine_match_sources,
    load_alias_config,
    load_fifa_ranking_rows,
    world_cup_match_rows_to_frame,
)
from app.metrics.schemas import Team
from app.metrics.team_normalizer import TeamNormalizer


def test_build_aliases_for_existing_teams_uses_mapper_names() -> None:
    config = {
        "teams": [
            {"source_name": "Brazil", "db_name": "Brasil", "aliases": ["BRA"]},
            {"source_name": "France", "db_name": "França", "aliases": ["FRA"]},
        ]
    }

    aliases, target_teams = build_aliases_for_existing_teams([Team(id="br", name="Brasil")], config)

    assert aliases["Brazil"] == "br"
    assert aliases["Brasil"] == "br"
    assert aliases["BRA"] == "br"
    assert [team.id for team in target_teams] == ["br"]


def test_load_fifa_ranking_rows_infers_rank_by_date(tmp_path) -> None:
    data_dir = tmp_path
    (data_dir / "ranking_fifa_historical.csv").write_text(
        "team,total_points,date,id,id_num,team_short\n"
        "Argentina,1855.2,2024-02-15,id14289,14289,ARG\n"
        "France,1845.44,2024-02-15,id14289,14289,FRA\n"
        "Argentina,1867.25,2024-04-04,id14290,14290,ARG\n",
        encoding="utf-8",
    )
    normalizer = TeamNormalizer(
        teams=[Team(id="ar", name="Argentina"), Team(id="fr", name="França")],
        aliases={"France": "fr"},
    )

    rows = load_fifa_ranking_rows(normalizer, data_dir)

    assert rows == [
        {"team_id": "ar", "ranking_date": date(2024, 2, 15), "rank": 1, "total_points": 1855.2},
        {"team_id": "fr", "ranking_date": date(2024, 2, 15), "rank": 2, "total_points": 1845.44},
        {"team_id": "ar", "ranking_date": date(2024, 4, 4), "rank": 1, "total_points": 1867.25},
    ]


def test_load_alias_config_reads_json(tmp_path) -> None:
    aliases_file = tmp_path / "aliases.json"
    aliases_file.write_text(json.dumps({"teams": []}), encoding="utf-8")

    assert load_alias_config(aliases_file) == {"teams": []}


def test_world_cup_match_rows_to_frame_matches_results_shape() -> None:
    frame = world_cup_match_rows_to_frame(
        [
            {
                "date": date(2026, 6, 20),
                "home_team": "Brasil",
                "away_team": "Argentina",
                "home_score": 2,
                "away_score": 1,
                "tournament": "FIFA World Cup",
                "city": "",
                "country": "",
                "neutral": False,
                "stage": "GROUP_STAGE",
            }
        ]
    )

    assert frame.iloc[0]["date"] == date(2026, 6, 20)
    assert frame.iloc[0]["home_score"] == 2
    assert "stage" in frame.columns


def test_combine_match_sources_keeps_finished_world_cup_version_on_duplicates() -> None:
    csv_matches = pd.DataFrame(
        [
            {
                "date": date(2026, 6, 20),
                "home_team_id": "br",
                "away_team_id": "ar",
                "home_score": 0,
                "away_score": 0,
                "tournament": "FIFA World Cup",
                "stage": None,
            }
        ]
    )
    finished_matches = pd.DataFrame(
        [
            {
                "date": date(2026, 6, 20),
                "home_team_id": "br",
                "away_team_id": "ar",
                "home_score": 2,
                "away_score": 1,
                "tournament": "FIFA World Cup",
                "stage": "GROUP_STAGE",
            }
        ]
    )

    combined = combine_match_sources(csv_matches, finished_matches)

    assert len(combined) == 1
    assert combined.iloc[0]["home_score"] == 2
    assert combined.iloc[0]["stage"] == "GROUP_STAGE"
