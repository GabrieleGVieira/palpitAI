from __future__ import annotations

import argparse
import logging
import sys
from datetime import date
from pathlib import Path

from dotenv import load_dotenv

sys.path.append(str(Path(__file__).resolve().parents[2]))

from app.metrics.database import Database
from app.metrics.local_data import (
    build_aliases_for_existing_teams,
    combine_match_sources,
    load_alias_config,
    load_goalscorers,
    load_results,
    load_shootouts,
    world_cup_match_rows_to_frame,
)
from app.metrics.metrics_calculator import calculate_team_metrics
from app.metrics.schemas import Snapshot
from app.metrics.team_normalizer import TeamNormalizer, attach_goalscorer_team_ids, attach_shootout_team_ids, attach_team_ids


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--metric-date", required=True, type=date.fromisoformat)
    parser.add_argument("--data-dir", default=None)
    parser.add_argument("--aliases-file", default=None)
    parser.add_argument("--include-finished-world-cup", action=argparse.BooleanOptionalAction, default=True)
    return parser.parse_args()


def calculate_and_save_team_metrics(
    db: Database,
    metric_date: date,
    data_dir: str | None = None,
    aliases_file: str | None = None,
    include_finished_world_cup: bool = True,
) -> dict:
    teams = db.load_teams()
    alias_config = load_alias_config(aliases_file)
    local_aliases, target_teams = build_aliases_for_existing_teams(teams, alias_config)
    aliases = db.load_aliases()
    aliases.update(local_aliases)
    normalizer = TeamNormalizer(teams, aliases)

    raw_matches = load_results(data_dir)
    matches = attach_team_ids(raw_matches, normalizer, fallback_unknown=True)
    finished_world_cup_matches = world_cup_match_rows_to_frame(
        db.load_finished_world_cup_matches_until(metric_date)
        if include_finished_world_cup
        else []
    )
    finished_world_cup_matches = attach_team_ids(finished_world_cup_matches, normalizer, fallback_unknown=True)
    matches = combine_match_sources(matches, finished_world_cup_matches)
    goalscorers = attach_goalscorer_team_ids(load_goalscorers(data_dir), normalizer, fallback_unknown=True)
    shootouts = attach_shootout_team_ids(load_shootouts(data_dir), normalizer, fallback_unknown=True)
    filtered = matches[matches["date"] <= metric_date]
    target_team_ids = [team.id for team in target_teams]
    metrics, snapshot_payloads = calculate_team_metrics(filtered, target_team_ids, metric_date, goalscorers, shootouts)

    db.upsert_team_metrics(metrics)
    db.insert_snapshots(
        [
            Snapshot(team_id=team_id, snapshot_type="team_metrics", payload_json=payload)
            for team_id, payload in snapshot_payloads.items()
        ]
    )
    unmapped = normalizer.report_unmapped()
    with_matches = sum(1 for item in metrics if item.matches_played > 0)
    return {
        "teams": len(metrics),
        "target_teams": len(target_team_ids),
        "teams_with_matches": with_matches,
        "matches_used": len(filtered),
        "finished_world_cup_matches": len(finished_world_cup_matches),
        "goalscorer_rows": len(goalscorers),
        "shootouts": len(shootouts),
        "unmapped": unmapped,
        "metric_date": metric_date,
    }


def main() -> None:
    load_dotenv()
    logging.basicConfig(level=logging.INFO, format="%(levelname)s %(name)s: %(message)s")
    args = parse_args()

    db = Database()
    summary = calculate_and_save_team_metrics(
        db,
        args.metric_date,
        args.data_dir,
        args.aliases_file,
        args.include_finished_world_cup,
    )
    print(
        "team metrics calculated: "
        f"teams={summary['teams']} target_teams={summary['target_teams']} "
        f"teams_with_matches={summary['teams_with_matches']} matches_used={summary['matches_used']} "
        f"finished_world_cup_matches={summary['finished_world_cup_matches']} "
        f"goalscorer_rows={summary['goalscorer_rows']} shootouts={summary['shootouts']} "
        f"unmapped_teams={len(summary['unmapped'])} metric_date={summary['metric_date']}"
    )
    if summary["unmapped"]:
        print("unmapped team names:")
        for name in summary["unmapped"]:
            print(f"- {name}")


if __name__ == "__main__":
    main()
