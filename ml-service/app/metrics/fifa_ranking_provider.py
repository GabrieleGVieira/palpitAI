from __future__ import annotations

from datetime import date
from pathlib import Path

from .local_data import load_fifa_ranking_rows


class InMemoryFifaRankingProvider:
    def __init__(self, rows: list[dict]) -> None:
        self.rows_by_team: dict[str, list[dict]] = {}
        for row in rows:
            self.rows_by_team.setdefault(row["team_id"], []).append(row)
        for team_rows in self.rows_by_team.values():
            team_rows.sort(key=lambda row: row["ranking_date"], reverse=True)

    def get_latest_ranking_before(self, team_id: str, ranking_date: date) -> int | None:
        for row in self.rows_by_team.get(team_id, []):
            if row["ranking_date"] < ranking_date:
                return row["rank"]
        return None


class CsvFifaRankingProvider(InMemoryFifaRankingProvider):
    def __init__(self, normalizer: object, data_dir: str | Path | None = None) -> None:
        super().__init__(load_fifa_ranking_rows(normalizer, data_dir))
