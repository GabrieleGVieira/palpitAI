from __future__ import annotations

import argparse
import logging
import sys
from datetime import date
from pathlib import Path

from dotenv import load_dotenv

sys.path.append(str(Path(__file__).resolve().parents[2]))

from app.metrics.database import Database
from app.metrics.fifa_ranking_provider import CsvFifaRankingProvider
from app.metrics.local_data import build_aliases_for_existing_teams, load_alias_config
from app.metrics.team_normalizer import TeamNormalizer


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--from-date", required=True, type=date.fromisoformat)
    parser.add_argument("--to-date", required=True, type=date.fromisoformat)
    parser.add_argument("--data-dir", default=None)
    parser.add_argument("--aliases-file", default=None)
    return parser.parse_args()


def build_feature(db: Database, ranking_provider: CsvFifaRankingProvider, target) -> dict:
    metrics = db.load_latest_team_metrics_before(
        [target.home_team_id, target.away_team_id],
        target.match_date,
    )
    home = metrics.get(target.home_team_id, {})
    away = metrics.get(target.away_team_id, {})
    home_rank = ranking_provider.get_latest_ranking_before(target.home_team_id, target.match_date)
    away_rank = ranking_provider.get_latest_ranking_before(target.away_team_id, target.match_date)
    home_elo = home.get("elo_score")
    away_elo = away.get("elo_score")

    return {
        "match_id": target.match_id,
        "match_date": target.match_date,
        "home_team_id": target.home_team_id,
        "away_team_id": target.away_team_id,
        "tournament": target.tournament,
        "stage": target.stage,
        "home_elo_score": home_elo,
        "away_elo_score": away_elo,
        "elo_diff": None if home_elo is None or away_elo is None else float(home_elo) - float(away_elo),
        "home_attack_score": home.get("attack_score"),
        "away_attack_score": away.get("attack_score"),
        "home_defense_score": home.get("defense_score"),
        "away_defense_score": away.get("defense_score"),
        "home_recent_form_score": home.get("recent_form_score"),
        "away_recent_form_score": away.get("recent_form_score"),
        "home_fifa_rank": home_rank,
        "away_fifa_rank": away_rank,
        "fifa_rank_diff": None if home_rank is None or away_rank is None else home_rank - away_rank,
        "home_avg_goals_scored": home.get("avg_goals_scored"),
        "away_avg_goals_scored": away.get("avg_goals_scored"),
        "home_avg_goals_conceded": home.get("avg_goals_conceded"),
        "away_avg_goals_conceded": away.get("avg_goals_conceded"),
        "home_world_cup_history_score": home.get("world_cup_history_score"),
        "away_world_cup_history_score": away.get("world_cup_history_score"),
        "neutral": target.neutral,
    }


def calculate_and_save_match_features(
    db: Database,
    from_date: date,
    to_date: date,
    data_dir: str | None = None,
    aliases_file: str | None = None,
) -> dict:
    teams = db.load_teams()
    alias_config = load_alias_config(aliases_file)
    local_aliases, _ = build_aliases_for_existing_teams(teams, alias_config)
    aliases = db.load_aliases()
    aliases.update(local_aliases)
    normalizer = TeamNormalizer(teams, aliases)
    ranking_provider = CsvFifaRankingProvider(normalizer, data_dir)
    targets = db.load_target_matches_with_normalizer(from_date, to_date, normalizer)
    features = [build_feature(db, ranking_provider, target) for target in targets]
    db.upsert_match_features(features)
    unmapped = normalizer.report_unmapped()
    return {
        "targets": len(targets),
        "saved": len(features),
        "from_date": from_date,
        "to_date": to_date,
        "unmapped": unmapped,
    }


def main() -> None:
    load_dotenv()
    logging.basicConfig(level=logging.INFO, format="%(levelname)s %(name)s: %(message)s")
    args = parse_args()

    db = Database()
    summary = calculate_and_save_match_features(db, args.from_date, args.to_date, args.data_dir, args.aliases_file)

    print(
        "match features calculated: "
        f"targets={summary['targets']} saved={summary['saved']} "
        f"from_date={summary['from_date']} to_date={summary['to_date']}"
    )
    if summary["unmapped"]:
        print("unmapped team names:")
        for name in summary["unmapped"]:
            print(f"- {name}")


if __name__ == "__main__":
    main()
