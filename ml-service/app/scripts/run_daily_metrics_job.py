from __future__ import annotations

import argparse
import logging
import sys
from datetime import date, timedelta
from pathlib import Path

from dotenv import load_dotenv

sys.path.append(str(Path(__file__).resolve().parents[2]))

from app.metrics.database import Database
from app.scripts.calculate_match_features import calculate_and_save_match_features
from app.scripts.calculate_team_metrics import calculate_and_save_team_metrics


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument("--metric-date", type=date.fromisoformat, default=None)
    parser.add_argument("--feature-from-date", type=date.fromisoformat, default=None)
    parser.add_argument("--feature-to-date", required=True, type=date.fromisoformat)
    parser.add_argument("--data-dir", default=None)
    parser.add_argument("--aliases-file", default=None)
    return parser.parse_args()


def main() -> None:
    load_dotenv()
    logging.basicConfig(level=logging.INFO, format="%(levelname)s %(name)s: %(message)s")
    args = parse_args()

    db = Database()
    metric_date = args.metric_date or date.today()
    feature_from_date = args.feature_from_date or metric_date + timedelta(days=1)

    metrics_summary = calculate_and_save_team_metrics(
        db,
        metric_date,
        args.data_dir,
        args.aliases_file,
        include_finished_world_cup=True,
    )
    features_summary = calculate_and_save_match_features(
        db,
        feature_from_date,
        args.feature_to_date,
        args.data_dir,
        args.aliases_file,
    )

    print(
        "daily metrics job completed: "
        f"metric_date={metric_date} "
        f"teams={metrics_summary['teams']} teams_with_matches={metrics_summary['teams_with_matches']} "
        f"matches_used={metrics_summary['matches_used']} "
        f"finished_world_cup_matches={metrics_summary['finished_world_cup_matches']} "
        f"feature_from_date={features_summary['from_date']} feature_to_date={features_summary['to_date']} "
        f"features_saved={features_summary['saved']}"
    )
    unmapped = sorted(set(metrics_summary["unmapped"]) | set(features_summary["unmapped"]))
    if unmapped:
        print("unmapped team names:")
        for name in unmapped:
            print(f"- {name}")


if __name__ == "__main__":
    main()
